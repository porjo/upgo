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

	resp, err := c.GetAccounts(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("resp %+v\n", resp)

	trans, err := c.GetTransactions(context.TODO(), "d63c429b-3538-46c0-bafb-2a26d0f34029")
	if err != nil {
		log.Fatal(err)
	}

	fmt.Printf("trans %+v\n", trans)

}
