package resource

type ListResponse struct {
	Schemas 	[]string	`json:"schemas"`
	TotalResults 	int		`json:"totalResults"`
	ItemsPerPage 	int		`json:"itemsPerPage"`
	StartIndex 	int 		`json:"startIndex"`
	Resources 	[]ScimObject	`json:"Resources"`
}

func (l *ListResponse) GetId() string {
	return ""
}

func (l *ListResponse) Data() map[string]interface{} {
	data := map[string]interface{}{
		"schemas": l.Schemas,
		"totalResults": l.TotalResults,
		"itemsPerPage": l.ItemsPerPage,
		"startIndex": l.StartIndex,
	}
	resourceData := make([]map[string]interface{}, 0, len(l.Resources))
	for _, r := range l.Resources {
		resourceData = append(resourceData, r.Data())
	}
	data["Resources"] = resourceData
	return data
}

func NewListResponse(resources []ScimObject, startIndex, itemsPerPage int) *ListResponse {
	return &ListResponse{
		Schemas:[]string{ListResponseUrn},
		TotalResults:len(resources),
		ItemsPerPage:itemsPerPage,
		StartIndex:startIndex,
		Resources:resources,
	}
}