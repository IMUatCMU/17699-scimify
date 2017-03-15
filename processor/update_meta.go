package processor

type updateMetaProcessor struct{}

func (ump *updateMetaProcessor) Process(ctx *ProcessorContext) error {
	r := getResource(ctx, true)

	if meta, ok := r.Attributes["meta"].(map[string]interface{}); !ok {
		panic(&PrerequisiteFailedError{reporter: "meta update", requirement: "meta"})
	} else if id, ok := r.Attributes["id"].(string); !ok {
		panic(&PrerequisiteFailedError{reporter: "meta update", requirement: "id"})
	} else {
		newMeta := make(map[string]interface{})
		for k, v := range meta {
			newMeta[k] = v
		}

		newMeta["version"] = generateNewVersion(id)
		newMeta["lastModified"] = getCurrentTime()
		r.Attributes["meta"] = newMeta

		return nil
	}
}
