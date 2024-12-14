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

package upgo

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"html/template"
	"log/slog"
	"net/http"
	"os"
	"time"

	"github.com/porjo/upgo/oapi"
	"golang.org/x/oauth2"
)

const (
	ServerURL       = "https://api.up.com.au/api/v1"
	TransactionsURL = ServerURL + "/accounts/{{ .accountID }}/transactions"
)

type Client struct {
	client   *http.Client
	upClient *oapi.ClientWithResponses

	logger *slog.Logger
}

// NewClientWithLogger is a wrapper for [NewClient] and takes a custom logger.
func NewClientWithLogger(logger *slog.Logger) (*Client, error) {
	c, err := NewClient()
	if err != nil {
		return nil, err
	}

	c.logger = logger

	return c, nil

}

// NewClient returns a Client
// It expects environment variable 'API_TOKEN' to be set.
// It will use a default [slog.Logger] log handler. Use [NewClientWithLogger] to pass in your own log handler.
func NewClient() (*Client, error) {

	token, ok := os.LookupEnv("API_TOKEN")
	if !ok {
		return nil, fmt.Errorf("environment variable API_TOKEN not set")
	}

	var err error
	c := &Client{}

	lvl := new(slog.LevelVar)
	lvl.Set(slog.LevelInfo)

	c.logger = slog.New(slog.NewTextHandler(os.Stderr, &slog.HandlerOptions{
		Level: lvl,
	}))

	// setting bearer token via roundtripper is a bit tricky
	// let oauth2 package take care of that for us
	// see: https://stackoverflow.com/a/51326483/202311
	c.client = oauth2.NewClient(context.Background(), oauth2.StaticTokenSource(&oauth2.Token{
		AccessToken: token,
		TokenType:   "Bearer",
	}))

	c.upClient, err = oapi.NewClientWithResponses(ServerURL, oapi.WithHTTPClient(c.client))
	if err != nil {
		return nil, fmt.Errorf("error getting client: %w", err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	_, err = c.upClient.GetUtilPing(ctx)
	if err != nil {
		return nil, fmt.Errorf("error pinging API: %w", err)
	}

	return c, nil
}

func (c *Client) GetAccounts(ctx context.Context) ([]oapi.AccountResource, error) {
	c.logger.Info("GetAccounts")
	resp, err := c.upClient.GetAccountsWithResponse(ctx, &oapi.GetAccountsParams{})
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("error getting accounts. Expected HTTP 200 but received %d", resp.StatusCode())
	}

	return resp.JSON200.Data, nil
}

func (c *Client) GetTransactions(ctx context.Context, accountID string) ([]oapi.TransactionResource, error) {

	var tpl bytes.Buffer

	templ := template.Must(template.New("getTransactions").Parse(TransactionsURL))
	err := templ.Execute(&tpl, map[string]interface{}{
		"accountID": accountID,
	})
	if err != nil {
		return nil, err
	}

	c.logger.Info("GetTransactions", "accountID", accountID)

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

	c.logger.Debug("getTransactions", "url", url, "transaction_count", len(transactions))

	if ltr.Links.Next != nil {
		t, err := c.getTransactions(ctx, *ltr.Links.Next)
		if err != nil {
			return nil, err
		}
		transactions = append(transactions, t...)
	}

	return transactions, nil

}
