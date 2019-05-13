package pageproperty

import "github.com/pkg/errors"

// Type is a valid property type.
type Type string

// All the valid values for Type
const (
	TypeNumber Type = "number"
	TypeString Type = "string"
)

// GetPermissionType returns the correct permission type for the given string.
func GetPermissionType(propertyTypeString string) (Type, error) {
	switch propertyTypeString {
	case string(TypeNumber):
		return TypeNumber, nil
	case string(TypeString):
		return TypeString, nil
	default:
		return TypeString, errors.Errorf("invalid property type %v", propertyTypeString)
	}
}

// PageProperty is a single property for a page.
type PageProperty struct {
	Key   string      `json:"key"`
	Type  Type        `json:"type"`
	Value interface{} `json:"value"`
}
