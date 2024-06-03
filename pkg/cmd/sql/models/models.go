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

type FilterItem struct {
	Field     string `json:"field"`
	Value     string `json:"value"`
	Condition string `json:"condition"`
}
type ResponseFilterItem struct {
	CaptionTitle       string     `json:"caption_title"`
	CaptionDescription string     `json:"caption_description"`
	Filter             FilterItem `json:"filter"`
}

type RecommendFilterResponse struct {
	Data struct {
		Recommendations []ResponseFilterItem `json:"recommendations"`
	} `json:"data"`
}

// Write implements io.Writer.
func (RecommendResponse) Write(p []byte) (n int, err error) {
	panic("unimplemented")
}
