// Model for attributes within a SCIM schema
package resource

import (
	"strings"
)

type AttributeGetter interface {
	GetAttribute(path string) *Attribute
}

type Attribute struct {
	Name            string       `json:"name"`
	Type            string       `json:"type"`
	SubAttributes   []*Attribute `json:"subAttributes"`
	MultiValued     bool         `json:"multiValued"`
	Description     string       `json:"description"`
	Required        bool         `json:"required"`
	CanonicalValues []string     `json:"canonicalValues"`
	CaseExact       bool         `json:"caseExact"`
	Mutability      string       `json:"mutability"`
	Returned        string       `json:"returned"`
	Uniqueness      string       `json:"uniqueness"`
	ReferenceTypes  []string     `json:"referenceTypes"`
	Assist          *Assist      `json:"_assist"`
}

func (a *Attribute) GetAttribute(path string) *Attribute {
	for _, attr := range a.SubAttributes {
		if strings.ToLower(attr.Name) == strings.ToLower(path) {
			return attr
		}
	}
	return nil
}

func (a *Attribute) Clone() *Attribute {
	cloned := &Attribute{
		Name:            a.Name,
		Type:            a.Type,
		MultiValued:     a.MultiValued,
		Description:     a.Description,
		Required:        a.Required,
		CanonicalValues: a.CanonicalValues,
		CaseExact:       a.CaseExact,
		Mutability:      a.Mutability,
		Returned:        a.Returned,
		Uniqueness:      a.Uniqueness,
		ReferenceTypes:  a.ReferenceTypes,
		Assist:          a.Assist,
		SubAttributes:   make([]*Attribute, 0),
	}
	for _, subAttr := range a.SubAttributes {
		cloned.SubAttributes = append(cloned.SubAttributes, subAttr.Clone())
	}
	return cloned
}
