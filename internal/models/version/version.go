package version

// Version keeps track of the version of a particular object.
type Version struct {
	ID   int64  `json:"-"`
	Name string `json:"name"`
}
