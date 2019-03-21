package permission

import "github.com/pkg/errors"

// Type is a valid permission type.
type Type string

// All the valid values for Type
const (
	TypePrivate    Type = "PR"
	TypePublic     Type = "PU"
	TypePublicOnly Type = "PO"
	TypeLinkOnly   Type = "LO"
)

// GetPermissionType returns the correct permission type for the given string.
func GetPermissionType(permissionTypeString string) (Type, error) {
	switch permissionTypeString {
	case string(TypePrivate):
		return TypePrivate, nil
	case string(TypePublic):
		return TypePublic, nil
	case string(TypePublicOnly):
		return TypePublicOnly, nil
	case string(TypeLinkOnly):
		return TypeLinkOnly, nil
	default:
		return TypePrivate, errors.Errorf("invalid permisson type %v", permissionTypeString)
	}
}

// IsPublic returns true if the Type is a type that is readable to the public.
func (t Type) IsPublic() bool {
	return t != TypePrivate
}
