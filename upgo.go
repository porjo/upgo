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
	"net/http"
	"os"

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
}

func NewClient() (*Client, error) {

	token, ok := os.LookupEnv("API_TOKEN")
	if !ok {
		return nil, fmt.Errorf("environment variable API_TOKEN not set")
	}

	var err error
	c := &Client{}

	// setting bearer token via roundtripper is a bit tricky
	// let oauth2 package take care of that for us
	// see: https://stackoverflow.com/a/51326483/202311
	ctx := context.Background()
	c.client = oauth2.NewClient(ctx, oauth2.StaticTokenSource(&oauth2.Token{
		AccessToken: token,
		TokenType:   "Bearer",
	}))

	c.upClient, err = oapi.NewClientWithResponses(ServerURL, oapi.WithHTTPClient(c.client))
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

	// get all results in a single request (page)
	q := req.URL.Query()
	q.Add("page[size]", "1")
	req.URL.RawQuery = q.Encode()

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
