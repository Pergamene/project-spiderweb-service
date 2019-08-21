package property

import "github.com/pkg/errors"

// Property is a key/value pair with a specified type.
type Property struct {
	ID    int64       `json:"-"`
	Key   string      `json:"key"`
	Type  Type        `json:"type"`
	Value interface{} `json:"value"`
}

// DBProperty is the Property struct as it comes out of the DB
type DBProperty struct {
	ID          int64
	Key         string
	Type        string
	StringValue string
	NumberValue float64
}

// Type is a valid property type.
type Type string

// All the valid values for Type
const (
	TypeNumber   Type = "number"
	TypeString   Type = "string"
	DBTypeNumber Type = "NU"
	DBTypeString Type = "ST"
)

// GetPropertyType returns the correct property type for the given json or DB stringified version.
func GetPropertyType(propertyTypeString string) (Type, error) {
	switch propertyTypeString {
	case string(TypeNumber), string(DBTypeNumber):
		return TypeNumber, nil
	case string(TypeString), string(DBTypeString):
		return TypeString, nil
	default:
		return TypeString, errors.Errorf("invalid property type %v", propertyTypeString)
	}
}

// GetDBPropertyType returns the property type as it is stored in the db.
func GetDBPropertyType(propertyType Type) (string, error) {
	switch propertyType {
	case TypeNumber:
		return string(DBTypeNumber), nil
	case TypeString:
		return string(DBTypeString), nil
	default:
		return string(DBTypeString), errors.Errorf("invalid property type %v", string(propertyType))
	}
}

// GetProperty returns a Property struct from the given reference DB Property
func (dp DBProperty) GetProperty() (Property, error) {
	pt, err := GetPropertyType(dp.Type)
	if err != nil {
		return Property{}, err
	}
	value, err := dp.getPropertyValue()
	if err != nil {
		return Property{}, err
	}
	return Property{
		ID:    dp.ID,
		Key:   dp.Key,
		Type:  pt,
		Value: value,
	}, nil
}

func (dp DBProperty) getPropertyValue() (interface{}, error) {
	switch dp.Type {
	case string(DBTypeNumber):
		return dp.NumberValue, nil
	case string(DBTypeString):
		return dp.StringValue, nil
	default:
		return dp.StringValue, errors.Errorf("invalid DB property type %v", dp.Type)
	}
}
