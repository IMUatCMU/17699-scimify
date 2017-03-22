package modify

import "github.com/go-scim/scimify/resource"

type Modifier interface {
	Modify(r *resource.Resource, sch *resource.Schema, mod *Modification) error
}

const (
	opAdd     = "add"
	opRemove  = "remove"
	opReplace = "replace"
)

type Modification struct {
	schemas    []string  `json:"schemas"`
	Operations []ModUnit `json:"Operations"`
}

type ModUnit struct {
	Op    string      `json:"op"`
	Path  string      `json:"path"`
	Value interface{} `json:"value"`
}

type DefaultModifier struct{}

func (dm *DefaultModifier) Modify(r *resource.Resource, sch *resource.Schema, mod *Modification) error {
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
