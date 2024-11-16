package upgo

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/porjo/upgo/oapi"
)

const (
	ServerURL = "https://api.up.com.au/api/v1"
)

type Client struct {
	client *oapi.ClientWithResponses
}

type Account struct {
	AccountType string
}

func NewClient() (*Client, error) {

	c := &Client{}
	hc := http.Client{}
	var err error
	c.client, err = oapi.NewClientWithResponses(ServerURL, oapi.WithHTTPClient(&hc))
	if err != nil {
		log.Fatal(err)
	}

	return c, nil

}

func (c *Client) GetAccounts(ctx context.Context) ([]Account, error) {
	resp, err := c.client.GetAccountsWithResponse(ctx, &oapi.GetAccountsParams{})
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("Expected HTTP 200 but received %d", resp.StatusCode())
	}

	fmt.Printf("resp.JSON200: %v\n", resp.JSON200)

	return nil, nil
}
