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

func main() {
	mux := &http.ServeMux{}

	mux.HandleFunc("/api/scrape", scrapeServeHTTP)
	mux.HandleFunc("/api/persist-first", persistFirstServeHTTP)
	mux.HandleFunc("/api/persist-second", persistSecondServeHTTP)
	mux.HandleFunc("/api/delete", deleteServeHTTP)

	handler := cors.Default().Handler(mux)

	log.Printf("Server listening on %s:%s", HOST, PORT)

	err := http.ListenAndServe(HOST+":"+PORT, handler)
	if err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server failed to start: %v", err)
	}

	log.Printf("bye bye")
}

const SCRAPE_WEBSITE_URL = "https://drewdevault.com"

var scrapeBlogPostsXPath = xpath.MustCompile("//section[@class='article-list']/div[@class='article']/a[@href]")

var (
	firstRepository  repository.Interface = repository.NewJSONServer("http://localhost:3000", http.DefaultClient)
	secondRepository repository.Interface = repository.NewJSONServer("http://localhost:3001", http.DefaultClient)
)

func scrapeServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

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

	nodes := htmlquery.QuerySelectorAll(doc, scrapeBlogPostsXPath)
	if err != nil {
		log.Printf("Failed to query HTML: %v", err)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	blogPosts := make([]types.BlogPost, 0, 10)
	for i := 0; i < cap(blogPosts); i++ {
		blogPosts = append(blogPosts, types.BlogPost{
			URL:   htmlquery.SelectAttr(nodes[i], "href"),
			Title: htmlquery.InnerText(nodes[i]),
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

func persistFirstServeHTTP(w http.ResponseWriter, r *http.Request) {
	defer r.Body.Close()

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
		_, err = firstRepository.SaveOne(r.Context(), v)
		if err != nil {
			log.Printf("Failed to save blog post: %v", err)
			w.WriteHeader(http.StatusInternalServerError)
			return
		}
	}

	resp, err := firstRepository.FindAll(r.Context())
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

func persistSecondServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func deleteServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
