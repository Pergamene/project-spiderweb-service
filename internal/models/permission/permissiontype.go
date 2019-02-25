package permission

// Type is a valid game type.
type Type string

// All the valid values for Type
const (
	TypePrivate    Type = "PR"
	TypePublic     Type = "PU"
	TypePublicOnly Type = "PO"
	TypeLinkOnly   Type = "LO"
)
