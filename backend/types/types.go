package types

type BlogPost struct {
	ID     *int64 `json:"id,omitempty"`
	Author string `json:"author"`
	Date   string `json:"date"`
	Title  string `json:"title"`
	URL    string `json:"url"`
}
