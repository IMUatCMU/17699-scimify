// Internally used fields for attributes
package resource

type Assist struct {
	JSONName      string   `json:"_jsonName"`      // JSON field name used to render this field
	Path          string   `json:"_path"`          // period delimited field names, useful to retrieve nested fields
	FullPath      string   `json:"_full_path"`     // Path prefixed with the URN of this resource
	ArrayIndexKey []string `json:"_arrayIndexKey"` // the field names of the multiValued complex fields that can be used as a search index
}
