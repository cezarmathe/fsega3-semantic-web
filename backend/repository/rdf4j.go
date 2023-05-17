package repository

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"strconv"
	"strings"

	"github.com/cezarmathe/semweb/backend/types"
)

func must[T any](t T, err error) T {
	if err != nil {
		panic(err)
	}

	return t
}

type RDF4J struct {
	base       *url.URL
	httpClient HTTPClient
}

func NewRDF4J(base string, httpClient HTTPClient) *RDF4J {
	return &RDF4J{
		base:       must(url.Parse(base)),
		httpClient: httpClient,
	}
}

func (r *RDF4J) DeleteOneByURL(ctx context.Context, url string) (types.BlogPost, error) {
	panic("not implemented")
}

func (r *RDF4J) FindAll(ctx context.Context) ([]types.BlogPost, error) {
	endpoint := must(url.Parse(r.base.String()))
	form := url.Values{}
	form.Set("query", `
		select ?subject ?predicate ?object from <http://example.com/articles> where { ?subject ?predicate ?object }
	`)

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint.String(), strings.NewReader(form.Encode()))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Content-Length", strconv.Itoa(len(form.Encode())))
	req.Header.Set("Accept", "application/sparql-results+json")

	resp, err := r.httpClient.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	var respBody rdfJSONResponse
	err = json.NewDecoder(resp.Body).Decode(&respBody)
	if err != nil {
		return nil, err
	}

	blogPosts := make(map[string]types.BlogPost, len(respBody.Results.Bindings)/3)

	for _, v := range respBody.Results.Bindings {
		blogPost := blogPosts[v.Subject.Value]

		blogPost.URL = v.Subject.Value

		switch v.Predicate.Value {
		case "http://example.com/articleauthor":
			blogPost.Author = v.Object.Value
		case "http://example.com/articledate":
			blogPost.Date = v.Object.Value
		case "http://example.com/articletitle":
			blogPost.Title = v.Object.Value
		}

		blogPosts[v.Subject.Value] = blogPost
	}

	result := make([]types.BlogPost, 0, len(blogPosts))

	for _, v := range blogPosts {
		result = append(result, v)
	}

	return result, nil
}

func (r *RDF4J) SaveOne(ctx context.Context, blogPost types.BlogPost) (types.BlogPost, error) {
	endpoint := must(url.Parse(r.base.String() + "/statements"))
	form := url.Values{}
	form.Set("update", fmt.Sprintf(
		`
		prefix articlepredicate: <http://example.com/article>

		insert data {
			graph <http://example.com/articles> {
				<%s>
					articlepredicate:author "%s";
					articlepredicate:date "%s";
					articlepredicate:title "%s".
			}
		}
		`,
		blogPost.URL, blogPost.Author, blogPost.Date, blogPost.Title,
	))
	payload := form.Encode()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint.String(), strings.NewReader(payload))
	if err != nil {
		return types.BlogPost{}, err
	}

	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Set("Content-Length", strconv.Itoa(len(payload)))

	resp, err := r.httpClient.Do(req)
	if err != nil {
		return types.BlogPost{}, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusNoContent {
		return types.BlogPost{}, fmt.Errorf("unexpected status code: %d", resp.StatusCode)
	}

	return blogPost, nil
}

type rdfJSONResponse struct {
	Head struct {
		Vars []string `json:"vars"`
	} `json:"head"`
	Results struct {
		Bindings []rdfBindingTriple `json:"bindings"`
	} `json:"results"`
}

type rdfBinding struct {
	Type  string `json:"type"`
	Value string `json:"value"`
}

type rdfBindingTriple struct {
	Subject   rdfBinding `json:"subject"`
	Predicate rdfBinding `json:"predicate"`
	Object    rdfBinding `json:"object"`
}
