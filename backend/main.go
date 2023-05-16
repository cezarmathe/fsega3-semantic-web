package main

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/antchfx/htmlquery"
	"github.com/antchfx/xpath"
)

const (
	HOST               = "0.0.0.0"
	PORT               = "8000"
	SCRAPE_WEBSITE_URL = "https://drewdevault.com"
)

func main() {
	server := &http.ServeMux{}

	server.HandleFunc("/api/scrape", scrapeServeHTTP)
	server.HandleFunc("/api/persist-first", persistFirstServeHTTP)
	server.HandleFunc("/api/persist-second", persistSecondServeHTTP)
	server.HandleFunc("/api/delete", deleteServeHTTP)

	log.Printf("Server listening on %s:%s", HOST, PORT)
	err := http.ListenAndServe(HOST+":"+PORT, server)
	if err != nil && err != http.ErrServerClosed {
		log.Fatalf("Server failed to start: %v", err)
	}

	log.Printf("bye bye")
}

type scrapedBlogPost struct {
	URL   string `json:"url"`
	Title string `json:"title"`
}

var scrapeBlogPostsXPath = xpath.MustCompile("//section[@class='article-list']/div[@class='article']/a[@href]")

func scrapeServeHTTP(w http.ResponseWriter, r *http.Request) {
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

	scrapedBlogPosts := make([]scrapedBlogPost, 0, 10)
	for i := 0; i < cap(scrapedBlogPosts); i++ {
		scrapedBlogPosts = append(scrapedBlogPosts, scrapedBlogPost{
			URL:   htmlquery.SelectAttr(nodes[i], "href"),
			Title: htmlquery.InnerText(nodes[i]),
		})
	}

	respBody, err := json.Marshal(scrapedBlogPosts)
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
	w.WriteHeader(http.StatusOK)
}

func persistSecondServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}

func deleteServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.WriteHeader(http.StatusOK)
}
