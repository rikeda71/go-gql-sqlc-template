package api

import (
	"io"
	"strings"
)

type Query struct {
	query string
}

func NewQuery(query string) Query {
	return Query{query: query}
}

func (q Query) RequestBody() io.Reader {
	return strings.NewReader(`{"query": "` + q.escaped() + `"}`)
}

func (q Query) escaped() string {
	return strings.Replace(strings.Replace(strings.Replace(q.query, "\n", "", -1), "\t", " ", -1), "\"", "\\\"", -1)
}
