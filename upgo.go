package upgo

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/porjo/upgo/oapi"
)

const (
	ServerURL = "https://api.up.com.au/api/v1"
)

type Client struct {
	client *oapi.ClientWithResponses
}

func NewClient() (*Client, error) {

	token, ok := os.LookupEnv("API_TOKEN")
	if !ok {
		return nil, fmt.Errorf("environment variable API_TOKEN not set")
	}

	hc := http.Client{}
	rt := getRTWithHeader(hc.Transport)
	rt.Set("Authorization", "Bearer "+token)
	hc.Transport = rt

	var err error
	c := &Client{}

	c.client, err = oapi.NewClientWithResponses(ServerURL, oapi.WithHTTPClient(&hc))
	if err != nil {
		return nil, err
	}

	_, err = c.client.GetUtilPing(context.TODO())
	if err != nil {
		return nil, fmt.Errorf("error pinging API %w", err)
	}

	return c, nil
}

func (c *Client) GetAccounts(ctx context.Context) ([]oapi.AccountResource, error) {
	resp, err := c.client.GetAccountsWithResponse(ctx, &oapi.GetAccountsParams{})
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("error getting accounts. Expected HTTP 200 but received %d", resp.StatusCode())
	}

	return resp.JSON200.Data, nil
}
