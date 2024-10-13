package internal

import (
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/rikeda71/go-gql-sqlc-template/internal/generated/db"
	"github.com/rikeda71/go-gql-sqlc-template/internal/generated/graph"
	"github.com/rikeda71/go-gql-sqlc-template/internal/metrics"
)

func NewGraphQLHandler(cnf *Config, dbc *db.Queries, m *metrics.Client) (*handler.Server, error) {
	// initialize usecase, service, or repository through selected architecture

	gqlHandler := *handler.NewDefaultServer(
		graph.NewExecutableSchema(graph.Config{Resolvers: &graph.Resolver{
			DBClient:      dbc,
			MetricsClient: m,
		}}),
	)

	return &gqlHandler, nil
}
