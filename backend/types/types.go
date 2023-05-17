package types

type BlogPost struct {
	Author string `json:"author"`
	Date   string `json:"date"`
	Title  string `json:"title"`
	URL    string `json:"url"`
}
