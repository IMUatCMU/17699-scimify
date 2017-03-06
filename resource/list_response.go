package resource

type ListResponse struct {
	Schemas      []string    `json:"schemas"`
	TotalResults int         `json:"totalResults"`
	ItemsPerPage int         `json:"itemsPerPage"`
	StartIndex   int         `json:"startIndex"`
	Resources    interface{} `json:"Resources"`
}

func (l *ListResponse) GetId() string {
	return ""
}

func (l *ListResponse) Data() map[string]interface{} {
	data := map[string]interface{}{
		"schemas":      l.Schemas,
		"totalResults": l.TotalResults,
		"itemsPerPage": l.ItemsPerPage,
		"startIndex":   l.StartIndex,
	}
	resourceData := make([]map[string]interface{}, 0)
	switch l.Resources.(type) {
	case []ScimObject:
		for _, r := range l.Resources.([]ScimObject) {
			resourceData = append(resourceData, r.Data())
		}
		data["Resources"] = resourceData
	default:
		data["Resources"] = l.Resources
	}

	return data
}

func NewListResponse(resources interface{}, startIndex, itemsPerPage, totalResults int) *ListResponse {
	return &ListResponse{
		Schemas:      []string{ListResponseUrn},
		TotalResults: totalResults,
		ItemsPerPage: itemsPerPage,
		StartIndex:   startIndex,
		Resources:    resources,
	}
}
