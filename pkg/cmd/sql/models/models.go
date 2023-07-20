package models

type SQLFieldsBody struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type SQLResponse struct {
	Records []map[string]interface{} `json:"records"`
	Fields  []SQLFieldsBody          `json:"fields"`
}
