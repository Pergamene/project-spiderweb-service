package version

// Version keeps track of the version of a particular object.
type Version struct {
	ID         int64  `json:"-"`
	GUID       string `json:"id"`
	Name       string `json:"name"`
	ParentGUID string `json:"parentId"`
}
