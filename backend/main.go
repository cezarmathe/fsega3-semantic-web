package main

import (
	"encoding/json"
	"io"
	"log"
	"net/http"

	"github.com/antchfx/htmlquery"
	"github.com/antchfx/xpath"
	"github.com/cezarmathe/semweb/backend/repository"
	"github.com/cezarmathe/semweb/backend/types"
	"github.com/rs/cors"
)

const (
	HOST = "0.0.0.0"
	PORT = "8000"
)

const (
	SCRAPE_WEBSITE_AUTHOR = "Drew DeVault"
	SCRAPE_WEBSITE_URL    = "https://drewdevault.com"
)

var (
	xpathArticleContainer    = xpath.MustCompile("//section[@class='article-list']/div[@class='article']")
	xpathArticleDate         = xpath.MustCompile("//span[@class='date']")
	xparthArticleTitleAndURL = xpath.MustCompile("//a[@href]")

	firstRepository repository.Interface = repository.NewJSONServer(
		"http://localhost:4000",
		http.DefaultClient,
	)
	secondRepository repository.Interface = repository.NewRDF4J(
		"http://localhost:8080/rdf4j-server/repositories/grafexamen",
		http.DefaultClient,
	)
)

func main() {
	mux := &http.ServeMux{}

	mux.HandleFunc("/api/scrape", scrapeServeHTTP)
	mux.HandleFunc("/api/persist-first", persistWithRepository(firstRepository))
	mux.HandleFunc("/api/persist-second", persistWithRepository(secondRepository))
	mux.HandleFunc("/api/delete", deleteWithRepository(firstRepository))

	handler := cors.AllowAll().Handler(mux)

	log.Printf("Server listening on %s:%s", HOST, PORT)

	err := http.ListenAndServe(HOST+":"+PORT, handler)
	if err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server failed to start: %v", err)
	}

	log.Printf("bye bye")
}

func scrapeServeHTTP(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	log.Printf("Scraping website: %s", SCRAPE_WEBSITE_URL)

	resp, err := http.Get(SCRAPE_WEBSITE_URL)
	if err != nil {
		log.Printf("Failed to scrape website: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		log.Printf("Failed to scrape website: %v", resp.Status)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	doc, err := htmlquery.Parse(resp.Body)
	if err != nil {
		log.Printf("Failed to parse HTML: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	nodes := htmlquery.QuerySelectorAll(doc, xpathArticleContainer)
	if err != nil {
		log.Printf("Failed to query HTML: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	blogPosts := make([]types.BlogPost, 0, 10)
	for i := 0; i < cap(blogPosts); i++ {
		node := nodes[i]

		dateNode := htmlquery.QuerySelector(node, xpathArticleDate)
		if dateNode == nil {
			log.Printf("Failed to query HTML: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		titleAndURLNode := htmlquery.QuerySelector(node, xparthArticleTitleAndURL)
		if titleAndURLNode == nil {
			log.Printf("Failed to query HTML: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		title := htmlquery.InnerText(titleAndURLNode)
		urlString := htmlquery.SelectAttr(titleAndURLNode, "href")
		date := htmlquery.InnerText(dateNode)

		blogPosts = append(blogPosts, types.BlogPost{
			Author: SCRAPE_WEBSITE_AUTHOR,
			Title:  title,
			URL:    urlString,
			Date:   date,
		})
	}

	respBody, err := json.Marshal(blogPosts)
	if err != nil {
		log.Printf("Failed to marshal JSON: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write(respBody)
}

func persistWithRepository(repository repository.Interface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		b, err := io.ReadAll(r.Body)
		if err != nil {
			log.Printf("Failed to read request body: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		var reqBody []types.BlogPost
		err = json.Unmarshal(b, &reqBody)
		if err != nil {
			log.Printf("Failed to unmarshal JSON: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		for _, v := range reqBody {
			_, err = repository.SaveOne(r.Context(), v)
			if err != nil {
				log.Printf("Failed to save blog post: %v", err)
				w.WriteHeader(http.StatusInternalServerError)
				return
			}
		}

		resp, err := repository.FindAll(r.Context())
		if err != nil {
			log.Printf("Failed to find all blog posts: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		respBody, err := json.Marshal(resp)
		if err != nil {
			log.Printf("Failed to marshal JSON: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(respBody)
	}
}

func deleteWithRepository(repository repository.Interface) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodDelete {
			w.WriteHeader(http.StatusMethodNotAllowed)
			return
		}

		author := r.URL.Query().Get("author")
		if author == "" {
			log.Printf("Missing \"author\" query parameter")
			w.WriteHeader(http.StatusBadRequest)
			return
		}

		_, err := repository.DeleteByAuthor(r.Context(), author)
		if err != nil {
			log.Printf("Failed to delete blog post: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		resp, err := repository.FindAll(r.Context())
		if err != nil {
			log.Printf("Failed to find all blog posts: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		respBody, err := json.Marshal(resp)
		if err != nil {
			log.Printf("Failed to marshal JSON: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		w.Write(respBody)
	}
}
