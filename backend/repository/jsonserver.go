package repository

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/url"
	"strconv"

	"github.com/cezarmathe/semweb/backend/types"
)

type JSONServer struct {
	base   string
	client HTTPClient
}

func NewJSONServer(base string, httpClient HTTPClient) *JSONServer {
	return &JSONServer{
		base:   base,
		client: httpClient,
	}
}

func (r *JSONServer) findAllByAuthor(ctx context.Context, author string) ([]types.BlogPost, error) {
	endpoint := r.base + "/posts?author=" + url.QueryEscape(author)

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}

	resp, err := r.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode == http.StatusNotFound {
		return nil, ErrNotFound
	}

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("unexpected status code")
	}

	var respBody []types.BlogPost

	err = json.NewDecoder(resp.Body).Decode(&respBody)
	if err != nil {
		return nil, err
	}

	return respBody, nil
}

func (r *JSONServer) DeleteByAuthor(ctx context.Context, author string) ([]types.BlogPost, error) {
	blogPostsToDelete, err := r.findAllByAuthor(ctx, author)
	if err != nil {
		return nil, err
	}

	for _, v := range blogPostsToDelete {
		endpoint := r.base + "/posts/" + url.PathEscape(strconv.FormatInt(*v.ID, 10))

		req, err := http.NewRequestWithContext(ctx, http.MethodDelete, endpoint, nil)
		if err != nil {
			return nil, err
		}

		resp, err := r.client.Do(req)
		if err != nil {
			return nil, err
		}

		defer resp.Body.Close()

		if resp.StatusCode == http.StatusNotFound {
			return nil, ErrNotFound
		}

		if resp.StatusCode != http.StatusOK {
			return nil, errors.New("unexpected status code")
		}
	}

	return blogPostsToDelete, nil
}

func (r *JSONServer) FindAll(ctx context.Context) ([]types.BlogPost, error) {
	endpoint := r.base + "/posts"

	req, err := http.NewRequestWithContext(ctx, http.MethodGet, endpoint, nil)
	if err != nil {
		return nil, err
	}

	resp, err := r.client.Do(req)
	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("unexpected status code")
	}

	var respBody []types.BlogPost

	err = json.NewDecoder(resp.Body).Decode(&respBody)
	if err != nil {
		return nil, err
	}

	return respBody, nil
}

func (r *JSONServer) SaveOne(ctx context.Context, v types.BlogPost) (types.BlogPost, error) {
	b, err := json.Marshal(v)
	if err != nil {
		return types.BlogPost{}, err
	}

	endpoint := r.base + "/posts"

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, endpoint, bytes.NewReader(b))
	if err != nil {
		return types.BlogPost{}, err
	}

	req.Header.Set("Content-Type", "application/json")

	resp, err := r.client.Do(req)
	if err != nil {
		return types.BlogPost{}, err
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusCreated {
		return types.BlogPost{}, errors.New("unexpected status code")
	}

	var respBody types.BlogPost

	err = json.NewDecoder(resp.Body).Decode(&respBody)
	if err != nil {
		return types.BlogPost{}, err
	}

	return respBody, nil
}
