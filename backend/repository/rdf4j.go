package repository

import (
	"context"

	"github.com/cezarmathe/semweb/backend/types"
)

type RDF4J struct {
	base       string
	httpClient HTTPClient
}

func NewRDF4J(base string, httpClient HTTPClient) *RDF4J {
	return &RDF4J{
		base:       base,
		httpClient: httpClient,
	}
}

func (r *RDF4J) DeleteOneByURL(ctx context.Context, url string) (types.BlogPost, error) {
	panic("not implemented")
}

func (r *RDF4J) FindAll(ctx context.Context) ([]types.BlogPost, error) {
	panic("not implemented")
}

func (r *RDF4J) SaveOne(ctx context.Context, blogPost types.BlogPost) (types.BlogPost, error) {
	panic("not implemented")
}

// graph: http://example.com
// <subject> <predicate> <object> format for triples
// <subject> is the URL
// <predicate> is the data type (always omit URL!)
// <object> is the value

// INSERTING:

// prefix articlepredicate: <http://example.com/article>

// insert data {
//   graph <http://example.com/articles> {
//     <https://drewdevault.com/2023/05/01/2023-05-01-Burnout.html>
//       articlepredicate:author "Drew DeVault";
//       articlepredicate:date "2023-05-01";
//       articlepredicate:title "Burnout".
//   }
// }

// SELECTING:

// select ?subject ?predicate ?object from <http://example.com/articles> where { ?subject ?predicate ?object }
