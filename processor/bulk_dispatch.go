package processor

import "encoding/json"

type BulkResponse struct {
	Schemas 	[]string		`json:"schemas"`
	Operations 	[]BulkResponseOperation	`json:"Operations"`
}

type BulkResponseOperation struct {
	BulkOperation
	Location	string			`json:"location"`
	Response 	json.RawMessage		`json:"response,omitempty"`
	Status 		int 			`json:"status"`
}

type BulkRequestSource struct {
	target 		string
	method 		string
	urlParams 	map[string]string
	params 		map[string]string
	body 		[]byte
}
func (rs *BulkRequestSource) Target() string {
	return rs.target
}
func (rs *BulkRequestSource) Method() string {
	return rs.method
}
func (rs *BulkRequestSource) UrlParam(name string) string {
	return rs.urlParams[name]
}
func (rs *BulkRequestSource) Param(name string) string {
	return rs.params[name]
}
func (rs *BulkRequestSource) Body() ([]byte, error) {
	return rs.body, nil
}

