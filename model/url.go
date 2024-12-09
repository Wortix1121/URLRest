package model

type UrlModel struct {
	Id    uint   `json:"id"`
	Alias string `json:"alias"`
	Url   string `json:"url"`
}
