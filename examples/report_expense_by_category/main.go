/*
Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"log/slog"
	"os"
	"time"

	"github.com/porjo/upgo"
	"github.com/porjo/upgo/oapi"
)

type payee struct {
	Name  string
	Total int64
}

type categories map[string][]payee

func main() {
	debug := flag.Bool("debug", false, "debug logging")
	flag.Parse()

	lvl := new(slog.LevelVar)
	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: lvl,
	}))
	slog.SetDefault(logger)

	if *debug {
		fmt.Println("Debug logging enabled")
		lvl.Set(slog.LevelDebug)
	}

	token, ok := os.LookupEnv("API_TOKEN")
	if !ok {
		log.Fatal("environment variable API_TOKEN not set")
	}

	c, err := upgo.NewClient(
		upgo.WithLogger(logger),
		upgo.WithToken(token),
	)
	if err != nil {
		log.Fatal(err)
	}

	filterSince := time.Now().AddDate(0, -3, 0)

	pageSize := 100
	transInput := &oapi.GetTransactionsParams{
		FilterSince: &filterSince,
		PageSize:    &pageSize,
	}

	slog.Info("Fetching transactions", "since", filterSince)

	r := make(categories)

	trans, err := c.GetTransactions(context.TODO(), transInput)
	if err != nil {
		log.Fatal(err)
	}

	fmt.Println("Transactions")

	for _, t := range trans {
		catStr := ""
		if t.Relationships.Category.Data != nil && t.Relationships.Category.Data.Id != "" {
			catStr = t.Relationships.Category.Data.Id
		}

		value := t.Attributes.Amount.ValueInBaseUnits
		if value > 0 {
			slog.Debug("Skipping income", "desc", t.Attributes.Description, "amount", value, "category", catStr)
			continue
		}

		if catStr == "" {
			catStr = "uncategorized"
		}

		if _, ok := r[catStr]; !ok {
			r[catStr] = []payee{}
		}

		p := r[catStr]
		found := false
		for i, p := range p {
			if p.Name == t.Attributes.Description {
				p.Total += -value
				r[catStr][i] = p
				found = true
				break
			}
		}
		if !found {
			r[catStr] = append(r[catStr], payee{
				Name:  t.Attributes.Description,
				Total: -value,
			})
		}
	}

	fmt.Printf("Expense Category Totals\n")
	fmt.Println()

	for cat, payees := range r {
		catStr := ""
		total := int64(0)
		for _, p := range payees {
			amount := float64(p.Total) / upgo.BaseUnitDivisor
			catStr += fmt.Sprintf("%-30s : %7.2f\n", p.Name, amount)
			total += p.Total
		}
		amount := float64(total) / upgo.BaseUnitDivisor

		fmt.Println(cat)
		fmt.Printf("-------------------------------- ----------------\n")
		fmt.Print(catStr)
		fmt.Printf("-------------------------------- ----------------\n")
		fmt.Printf("%30s : %7.2f\n", "total", amount)
		fmt.Println()
	}

}
