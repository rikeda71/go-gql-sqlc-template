//go:build api

package api_test

import (
	"encoding/json"
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/rikeda71/go-gql-sqlc-template/internal/generated/graph"
	api "github.com/rikeda71/go-gql-sqlc-template/test/api/helper"
)

type createUserMutationResponse struct {
	Data struct {
		CreateUserOutput graph.CreateUserOutput `json:"createUser"`
	} `json:"data"`
}

type userQueryResponse struct {
	Data struct {
		User graph.User `json:"user"`
	} `json:"data"`
}

func TestCrudUser(t *testing.T) {

	t.Parallel()

	// create
	/// given
	createUserMutation := api.NewQuery(`
	mutation CreateUser {
		createUser(input: {name: "test", email: "test@example.com"}) {
			status
			errorMessage
			metadata {
				user {
					id
				}
			}
		}
	}
	`)

	/// when
	resBytes1, err := api.PostGraphQLRequest(createUserMutation, Server)
	if err != nil {
		t.Errorf("cause error when post graphql request. error = %v", err)
	}

	/// then
	var userID string
	{
		var actual createUserMutationResponse
		err = json.Unmarshal(resBytes1, &actual)
		if err != nil {
			t.Errorf("cause error when unmarshal response. error = %v", err)
		}
		// get userID for later test
		userID = actual.Data.CreateUserOutput.Metadata.User.ID

		// set fixed ID for comparison
		actual.Data.CreateUserOutput.Metadata.User.ID = "1"
		expected := createUserMutationResponse{
			Data: struct {
				CreateUserOutput graph.CreateUserOutput `json:"createUser"`
			}{
				CreateUserOutput: graph.CreateUserOutput{
					Status: graph.MutationStatusSuccess,
					Metadata: &graph.CreateUserOutputMetadata{
						User: &graph.User{
							ID: "1",
						},
					},
					ErrorMessage: nil,
				},
			},
		}

		if got := cmp.Diff(actual, expected); got != "" {
			t.Errorf("unexpected response: %v", got)
		}
	}

	// query
	/// given
	userQuery := api.NewQuery(fmt.Sprintf(`
	query User {
		user(id: "%s") {
			id
			name
			email
		}
	}
	`,
		userID,
	))

	/// when
	resBytes2, err := api.PostGraphQLRequest(userQuery, Server)
	if err != nil {
		t.Errorf("cause error when post graphql request. error = %v", err)
	}

	/// then
	{
		var actual userQueryResponse
		err = json.Unmarshal(resBytes2, &actual)
		if err != nil {
			t.Errorf("cause error when unmarshal response. error = %v", err)
		}
		expected := userQueryResponse{
			Data: struct {
				User graph.User `json:"user"`
			}{
				User: graph.User{
					ID:    userID,
					Name:  "test",
					Email: "test@example.com",
				},
			},
		}

		if got := cmp.Diff(actual, expected); got != "" {
			t.Errorf("unexpected response: %v", got)
		}
	}
}
