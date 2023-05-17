package repository_test

import (
	"context"
	"net/http"
	"testing"

	"github.com/cezarmathe/semweb/backend/repository"
	"github.com/cezarmathe/semweb/backend/types"
)

func TestRDF4J(t *testing.T) {
	ctx := context.Background()

	client := repository.NewRDF4J(
		"http://localhost:8080/rdf4j-server/repositories/grafexamen",
		http.DefaultClient,
	)

	expected := types.BlogPost{
		Date:  "1970-01-01",
		Title: "Hello, world!",
		URL:   "https://example.com",
	}

	actual, err := client.SaveOne(ctx, expected)
	if err != nil {
		t.Fatalf("SaveOne failed: %v", err)
	}

	if actual != expected {
		t.Fatalf("SaveOne failed: expected %v, got %v", expected, actual)
	}
}
