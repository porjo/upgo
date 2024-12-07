package main

import (
	"context"
	"fmt"
	"log"

	"github.com/porjo/upgo"
)

func main() {

	c, err := upgo.NewClient()
	if err != nil {
		log.Fatal(err)
	}

	accounts, err := c.GetAccounts(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("accounts %+v\n", accounts)

	for _, a := range accounts {
		trans, err := c.GetTransactions(context.TODO(), a.Id)
		if err != nil {
			log.Fatal(err)
		}

		fmt.Printf("trans %+v\n", trans)
	}

}
