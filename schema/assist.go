// Internally used fields for attributes
package schema

type Assist struct {
	JSONName      string `json:"_jsonName"`      // JSON field name used to render this field
	Path          string `json:"_path"`          // period delimited field names, useful to retrieve nested fields
	ArrayIndexKey string `json:"_arrayIndexKey"` // the field name of the multiValued complex fields that can be used as a search index
}
