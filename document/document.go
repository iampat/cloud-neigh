package document

import "github/iampat/cloudy-neigh/vector"

type Document struct {
	Id    string `json:"id"`
	Title string `json:"title"`
	Url   string `json:"url"`
	Text  string `json:"text"`

	TextEmbedding *vector.Vector32 `json:"text_embedding"`
	TextLshHash   string           `json:"text_lsh_hash"`

	TitleEmbedding *vector.Vector32 `json:"title_embedding"`
	TitleLshHash   string           `json:"title_lsh_hash"`
}
