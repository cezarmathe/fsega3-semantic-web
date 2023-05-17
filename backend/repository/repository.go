package repository

import (
	"context"
	"net/http"

	"github.com/cezarmathe/semweb/backend/types"
)

type stringError string

func (e stringError) Error() string {
	return string(e)
}

const (
	ErrNotFound stringError = "not found"
)

type HTTPClient interface {
	Do(*http.Request) (*http.Response, error)
}

type Interface interface {
	DeleteOneByURL(context.Context, string) (types.BlogPost, error)
	FindAll(context.Context) ([]types.BlogPost, error)
	SaveOne(context.Context, types.BlogPost) (types.BlogPost, error)
}
