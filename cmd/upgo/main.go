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

	c, err := upgo.NewClientWithLogger(logger)
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
