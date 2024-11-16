package upgo

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"
	"os"

	"github.com/porjo/upgo/oapi"
)

const (
	ServerURL       = "https://api.up.com.au/api/v1"
	TransactionsURL = ServerURL + "/accounts/{{ .accountID }}/transactions"
)

type Client struct {
	client   http.Client
	upClient *oapi.ClientWithResponses
}

func NewClient() (*Client, error) {

	token, ok := os.LookupEnv("API_TOKEN")
	if !ok {
		return nil, fmt.Errorf("environment variable API_TOKEN not set")
	}

	var err error
	c := &Client{}

	c.client = http.Client{}
	wh := getRTWithHeader(c.client.Transport)
	wh.Set("Authorization", "Bearer "+token)
	c.client.Transport = wh

	c.upClient, err = oapi.NewClientWithResponses(ServerURL, oapi.WithHTTPClient(&c.client))
	if err != nil {
		return nil, fmt.Errorf("error getting client: %w", err)
	}

	_, err = c.upClient.GetUtilPing(context.TODO())
	if err != nil {
		return nil, fmt.Errorf("error pinging API: %w", err)
	}

	return c, nil
}

func (c *Client) GetAccounts(ctx context.Context) ([]oapi.AccountResource, error) {
	resp, err := c.upClient.GetAccountsWithResponse(ctx, &oapi.GetAccountsParams{})
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("error getting accounts. Expected HTTP 200 but received %d", resp.StatusCode())
	}

	return resp.JSON200.Data, nil
}

/*
func (c *Client) GetTransactions(ctx context.Context, accountID string) ([]oapi.TransactionResource, error) {
	resp, err := c.upClient.GetTransactionsWithResponse(ctx, &oapi.GetTransactionsParams{})
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("error getting transactions. Expected HTTP 200 but received %d", resp.StatusCode())
	}

	if resp.JSON200 == nil {
		return []oapi.TransactionResource{}, nil
	}

	transactions := resp.JSON200.Data

	if resp.JSON200.Links.Next != nil {
		t, err := c.getTransactions(ctx, *resp.JSON200.Links.Next)
		if err != nil {
			return nil, err
		}

		transactions = append(transactions, t...)

	}

	return transactions, nil
}
*/

func (c *Client) GetTransactions(ctx context.Context, accountID string) ([]oapi.TransactionResource, error) {

	var tpl bytes.Buffer

	templ := template.Must(template.New("getTransactions").Parse(TransactionsURL))
	err := templ.Execute(&tpl, map[string]interface{}{
		"accountID": accountID,
	})
	if err != nil {
		return nil, err
	}

	return c.getTransactions(ctx, tpl.String())
}

func (c *Client) getTransactions(ctx context.Context, url string) ([]oapi.TransactionResource, error) {

	req, err := http.NewRequestWithContext(ctx, "GET", url, nil)
	if err != nil {
		return nil, err
	}

	r, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer r.Body.Close()

	ltr := &oapi.ListTransactionsResponse{}
	err = json.NewDecoder(r.Body).Decode(ltr)
	if err != nil {
		return nil, err
	}

	transactions := ltr.Data

	if ltr.Links.Next != nil {
		t, err := c.getTransactions(ctx, *ltr.Links.Next)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, t...)
	}

	return transactions, nil

}
