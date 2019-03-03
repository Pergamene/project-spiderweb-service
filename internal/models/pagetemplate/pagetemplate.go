package pagetemplate

// PageTemplate keeps track of the pagetemplate of a particular object.
type PageTemplate struct {
	ID   int64  `json:"-"`
	Name string `json:"name"`
}
