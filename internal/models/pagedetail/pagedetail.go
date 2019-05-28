package pagedetail

// PageDetail is a single detail for a page.
type PageDetail struct {
	ID         int64       `json:"-"`
	GUID       string      `json:"id"`
	Title      string      `json:"title"`
	Summary    string      `json:"summary"`
	Partitions []Partition `json:"partitions"`
}
