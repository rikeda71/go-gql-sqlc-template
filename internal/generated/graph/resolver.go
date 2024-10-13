package graph

import (
	"github.com/rikeda71/go-gql-sqlc-template/internal/generated/db"
	"github.com/rikeda71/go-gql-sqlc-template/internal/metrics"
)

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct {
	DBClient      *db.Queries
	MetricsClient *metrics.Client
}
