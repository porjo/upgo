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
	"context"
	"fmt"
	"io"
	"log/slog"
	"net/http"
	"time"

	"github.com/porjo/upgo/oapi"
	"golang.org/x/oauth2"
)

const (
	ServerURL = "https://api.up.com.au/api/v1"
)

type Client struct {
	client   *http.Client
	upClient *oapi.ClientWithResponses

	token string

	logger *slog.Logger
}

type ClientOption func(*Client)

// WithLogger allows the user to define how log emitted by upgo will be handled.
// Upgo does not emit log otherwise.
func WithLogger(logger *slog.Logger) ClientOption {
	return func(c *Client) {
		c.logger = logger
	}
}

// WithToken is used to supply an API auth token to upgo.
func WithToken(token string) ClientOption {
	return func(c *Client) {
		c.token = token
	}
}

// NewClient returns a Client
// It will use a default [slog.Logger] log handler unless overriden with [WithLogger]
func NewClient(opts ...ClientOption) (*Client, error) {

	var err error
	c := &Client{}

	// by default log is discarded
	c.logger = slog.New(slog.NewTextHandler(io.Discard, nil))

	// Apply each option
	for _, opt := range opts {
		opt(c)
	}

	// setting bearer token via roundtripper is a bit tricky
	// let oauth2 package take care of that for us
	// see: https://stackoverflow.com/a/51326483/202311
	c.client = oauth2.NewClient(context.Background(), oauth2.StaticTokenSource(&oauth2.Token{
		AccessToken: c.token,
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

// GetTransactions returns transactions for all accounts, optionally filtered by [oapi.GetTransactionsParams].
func (c *Client) GetTransactions(ctx context.Context, params *oapi.GetTransactionsParams) ([]oapi.TransactionResource, error) {
	c.logger.Info("GetTransactions")
	resp, err := c.upClient.GetTransactionsWithResponse(ctx, params)
	if err != nil {
		return nil, err
	}
	if resp.StatusCode() != http.StatusOK {
		return nil, fmt.Errorf("error getting transactions. Expected HTTP 200 but received %d", resp.StatusCode())
	}

	return resp.JSON200.Data, nil
}
