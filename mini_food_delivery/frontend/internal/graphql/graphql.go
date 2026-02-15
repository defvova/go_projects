package graphql

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"time"
)

type GraphQLRequest struct {
	Query         string         `json:"query"`
	Variables     map[string]any `json:"variables,omitempty"`
	OperationName string         `json:"operationName,omitempty"`
}

type GraphQLError struct {
	Message string `json:"message"`
}

type GraphQLResponse[T any] struct {
	Data   T              `json:"data"`
	Errors []GraphQLError `json:"errors,omitempty"`
}

type GraphQLClient struct {
	url    string
	client *http.Client
	token  string
}

func NewClient(url, token string) *GraphQLClient {
	return &GraphQLClient{
		url:   url,
		token: token,
		client: &http.Client{
			Timeout: 5 * time.Second,
		},
	}
}

func (c *GraphQLClient) Do(ctx context.Context, req GraphQLRequest, out any) error {
	buf, err := json.Marshal(req)
	if err != nil {
		return err
	}

	httpReq, err := http.NewRequestWithContext(
		ctx,
		http.MethodPost,
		c.url,
		bytes.NewReader(buf),
	)
	if err != nil {
		return err
	}

	httpReq.Header.Set("Content-Type", "application/json")
	httpReq.Header.Set("Authorization", "Bearer "+c.token)

	resp, err := c.client.Do(httpReq)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	var gqlResp struct {
		Data   json.RawMessage `json:"data"`
		Errors []GraphQLError  `json:"errors"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&gqlResp); err != nil {
		return err
	}

	if len(gqlResp.Errors) > 0 {
		return errors.New(gqlResp.Errors[0].Message)
	}

	return json.Unmarshal(gqlResp.Data, out)
}
