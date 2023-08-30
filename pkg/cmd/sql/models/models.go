package models

type SQLFieldsBody struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type SQLResponse struct {
	Records []map[string]interface{} `json:"records"`
	Fields  []SQLFieldsBody          `json:"fields"`
}

type ResponseItem struct {
	CaptionTitle       string `json:"caption_title"`
	CaptionDescription string `json:"caption_description"`
	SQLStatement       string `json:"sql_statement"`
}

type RecommendResponse struct {
	Data []ResponseItem `json:"data"`
}
