package pagedetail

import "github.com/pkg/errors"

// PartitionType is a valid property type.
type PartitionType string

// All the valid values for PartitionType
const (
	PartitionTypeText          PartitionType = "text"
	PartitionTypeBold          PartitionType = "bold"
	PartitionTypeItalics       PartitionType = "italics"
	PartitionTypeLink          PartitionType = "link"
	PartitionTypeRelation      PartitionType = "relation"
	PartitionTypeColor         PartitionType = "color"
	PartitionTypeUnorderedList PartitionType = "ul"
	PartitionTypeOrderedList   PartitionType = "ol"
	PartitionTypeImage         PartitionType = "image"
	PartitionTypeQuote         PartitionType = "quote"
	PartitionTypePageBreak     PartitionType = "hr"
)

// GetPermissionType returns the correct permission type for the given string.
func GetPermissionType(propertyTypeString string) (PartitionType, error) {
	switch propertyTypeString {
	case string(PartitionTypeText):
		return PartitionTypeText, nil
	case string(PartitionTypeBold):
		return PartitionTypeBold, nil
	case string(PartitionTypeItalics):
		return PartitionTypeItalics, nil
	case string(PartitionTypeLink):
		return PartitionTypeLink, nil
	case string(PartitionTypeRelation):
		return PartitionTypeRelation, nil
	case string(PartitionTypeColor):
		return PartitionTypeColor, nil
	case string(PartitionTypeUnorderedList):
		return PartitionTypeUnorderedList, nil
	case string(PartitionTypeOrderedList):
		return PartitionTypeOrderedList, nil
	case string(PartitionTypeImage):
		return PartitionTypeImage, nil
	case string(PartitionTypeQuote):
		return PartitionTypeQuote, nil
	case string(PartitionTypePageBreak):
		return PartitionTypePageBreak, nil
	default:
		return PartitionTypeText, errors.Errorf("invalid property type %v", propertyTypeString)
	}
}

// PageDetail is a single detail for a page.
type PageDetail struct {
	ID         int64       `json:"-"`
	GUID       string      `json:"id"`
	Title      string      `json:"title"`
	Summary    string      `json:"summary"`
	Partitions []Partition `json:"partitions"`
}

// Partition is a single markdown partition for a detail.
type Partition struct {
	Type       PartitionType `json:"type"`
	Value      string        `json:"value,omitempty"`
	Partitions []Partition   `json:"partitions,omitempty"`
	Link       string        `json:"link,omitempty"`
	Relation   string        `json:"relation,omitempty"`
	Color      string        `json:"color,omitempty"`
}
