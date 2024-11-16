package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	client "github.com/porjo/upgo"
)

func main() {
	// custom HTTP client
	hc := http.Client{}

	c, err := client.NewClientWithResponses("http://localhost:1234", client.WithHTTPClient(&hc))
	if err != nil {
		log.Fatal(err)
	}

	resp, err := c.GetAccountsWithResponse(context.TODO(), &client.GetAccountsParams{})
	if err != nil {
		log.Fatal(err)
	}
	if resp.StatusCode() != http.StatusOK {
		log.Fatalf("Expected HTTP 200 but received %d", resp.StatusCode())
	}

	fmt.Printf("resp.JSON200: %v\n", resp.JSON200)

}
