package modify

import (
	"fmt"
	"github.com/go-scim/scimify/resource"
)

type Modifier interface {
	Modify(r *resource.Resource, sch *resource.Schema, mod *Modification) error
}

const (
	opAdd     = "add"
	opRemove  = "remove"
	opReplace = "replace"
)

type Modification struct {
	Schemas    []string  `json:"schemas"`
	Operations []ModUnit `json:"Operations"`
}

type ModUnit struct {
	Op    string      `json:"op"`
	Path  string      `json:"path"`
	Value interface{} `json:"value"`
}

type DefaultModifier struct{}

func (dm *DefaultModifier) Modify(r *resource.Resource, sch *resource.Schema, mod *Modification) (err error) {
	defer func() {
		if r := recover(); r != nil {
			err = &ModificationFailedError{r}
		}
	}()

	if len(mod.Schemas) != 1 || mod.Schemas[0] != resource.PathOpUrn {
		return &InvalidModificationError{fmt.Sprintf("schemas must be [%s]", resource.PathOpUrn)}
	}

	baseContainer := modMap(r.Data())
	for _, unit := range mod.Operations {
		tokens, err := tokenize(unit.Path)
		if err != nil {
			return err
		}

		err = baseContainer.Apply(tokens, sch.AsAttribute(), unit, nil)
		if err != nil {
			return err
		}
	}

	return nil
}
