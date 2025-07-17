package model

type Book struct {
	CoverURL  string `json:"cover_url"`
	Title     string `json:"title"`
	DetailURL string `json:"detail_url"`
	Author    string `json:"author"`
	Intro     string `json:"intro"`
}

type BookElse struct {
	URL    string `json:"url"`
	Title  string `json:"title"`
	Author string `json:"author"`
}
