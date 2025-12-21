package models

type SearchGameVo struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Info  string `json:"info"`
	Cover string `json:"cover"`
}

type SearchRequest struct {
	Txt  string `json:"txt"`
	Lang string `json:"lang"`
}
