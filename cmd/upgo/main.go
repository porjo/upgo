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

}
