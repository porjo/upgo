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

	filterSince := time.Now().AddDate(0, -2, 0)
	filterUntil := time.Now().AddDate(0, -1, 0)

	transInput := &oapi.GetTransactionsParams{
		FilterSince: &filterSince,
		FilterUntil: &filterUntil,
	}

	slog.Info("Fetching transactions", "since", filterSince, "until", filterUntil)

	categories := make(map[string]int64)

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

		//slog.Debug("Trans ", "desc", t.Attributes.Description, "amount", t.Attributes.Amount.ValueInBaseUnits, "category", catStr)

		value := t.Attributes.Amount.ValueInBaseUnits
		if value < 0 {
			value = -value
			if catStr == "" {
				catStr = "income"
			}
		}

		if catStr == "" {
			catStr = "uncategorized"
		}
		categories[catStr] += int64(value)
	}

	fmt.Printf("Category Totals\n")

	sortedCats := upgo.SortMapByValue(categories)
	for _, cat := range sortedCats {
		amount := float64(cat.Value) / upgo.BaseUnitDivisor
		fmt.Printf("%-30s : %.2f\n", cat.Key, amount)
	}

}
