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
	"log"
	"log/slog"
	"os"

	"github.com/porjo/upgo"
)

func main() {
	lvl := new(slog.LevelVar)
	lvl.Set(slog.LevelDebug)

	logger := slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: lvl,
	}))

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

	accounts, err := c.GetAccounts(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	slog.Info("accounts", "accounts", accounts)

	for _, a := range accounts {
		trans, err := c.GetTransactions(context.TODO(), a.Id)
		if err != nil {
			log.Fatal(err)
		}

		slog.Info("transactions", "transactions", trans)
	}

}
