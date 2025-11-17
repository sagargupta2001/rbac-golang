package service

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"

	"rbac/internal/domain"
)

// The URL of the public GraphQL API
const countriesGraphQLAPI = "https://countries.trevorblades.com/"

// GraphQLService defines the interface for our client
type GraphQLService interface {
	GetCountryDetails(ctx context.Context, countryCode string) (*domain.CountryDetails, error)
}

// graphqlService is the implementation
type graphqlService struct {
	client *http.Client
}

// NewGraphQLService creates a new GraphQLService
func NewGraphQLService() GraphQLService {
	return &graphqlService{
		client: &http.Client{},
	}
}

func (s *graphqlService) GetCountryDetails(ctx context.Context, countryCode string) (*domain.CountryDetails, error) {
	// Define the GraphQL query
	query := `
        query GetCountry($code: ID!) {
            country(code: $code) {
                name
                capital
                emoji
				currency
            }
        }
    `

	// Define the variables for the query
	variables := map[string]interface{}{
		"code": countryCode,
	}

	// Create the GraphQL request body
	reqBody := domain.GraphQLRequest{
		Query:     query,
		Variables: variables,
	}

	// Marshal the request body to JSON
	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return nil, fmt.Errorf("failed to marshal request: %w", err)
	}

	// Create the HTTP POST request
	req, err := http.NewRequestWithContext(ctx, "POST", countriesGraphQLAPI, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create http request: %w", err)
	}
	req.Header.Set("Content-Type", "application/json")

	// Send the request
	res, err := s.client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send request: %w", err)
	}
	defer res.Body.Close()

	if res.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %s", res.Status)
	}

	// Decode the JSON response
	var graphqlResponse domain.GraphQLCountryResponse
	if err := json.NewDecoder(res.Body).Decode(&graphqlResponse); err != nil {
		return nil, fmt.Errorf("failed to decode response: %w", err)
	}

	// Return the clean, nested data
	countryData := graphqlResponse.Data.Country
	return &domain.CountryDetails{
		Name:    countryData.Name,
		Capital: countryData.Capital,
		Emoji:   countryData.Emoji,
		Currency: countryData.Currency,
	}, nil
}