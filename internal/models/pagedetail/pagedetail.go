package pagedetail

import "github.com/pkg/errors"

// PartitionType is a valid property type.
type PartitionType string

// All the valid values for PartitionType
const (
	PartitionTypeHeaderOne     PartitionType = "h1"
	PartitionTypeHeaderTwo     PartitionType = "h2"
	PartitionTypeHeaderThree   PartitionType = "h3"
	PartitionTypeHeaderFour    PartitionType = "h4"
	PartitionTypeHeaderFive    PartitionType = "h5"
	PartitionTypeHeaderSix     PartitionType = "h6"
	PartitionTypeParagraph     PartitionType = "p"
	PartitionTypeUnorderedList PartitionType = "ul"
	PartitionTypeOrderedList   PartitionType = "ol"
	PartitionTypeImage         PartitionType = "image"
	PartitionTypeQuotes        PartitionType = "quotes"
	PartitionTypePageBreak     PartitionType = "hr"
	PartitionTypeText          PartitionType = "text"
	PartitionTypeBold          PartitionType = "bold"
	PartitionTypeItalics       PartitionType = "italics"
	PartitionTypeLink          PartitionType = "link"
	PartitionTypeRelation      PartitionType = "relation"
	PartitionTypeColor         PartitionType = "color"
)

// GetPartitionType returns the correct permission type for the given string.
func GetPartitionType(propertyTypeString string) (PartitionType, error) {
	switch propertyTypeString {
	case string(PartitionTypeHeaderOne):
		return PartitionTypeHeaderOne, nil
	case string(PartitionTypeHeaderTwo):
		return PartitionTypeHeaderTwo, nil
	case string(PartitionTypeHeaderThree):
		return PartitionTypeHeaderThree, nil
	case string(PartitionTypeHeaderFour):
		return PartitionTypeHeaderFour, nil
	case string(PartitionTypeHeaderFive):
		return PartitionTypeHeaderFive, nil
	case string(PartitionTypeHeaderSix):
		return PartitionTypeHeaderSix, nil
	case string(PartitionTypeParagraph):
		return PartitionTypeParagraph, nil
	case string(PartitionTypeUnorderedList):
		return PartitionTypeUnorderedList, nil
	case string(PartitionTypeOrderedList):
		return PartitionTypeOrderedList, nil
	case string(PartitionTypeImage):
		return PartitionTypeImage, nil
	case string(PartitionTypeQuotes):
		return PartitionTypeQuotes, nil
	case string(PartitionTypePageBreak):
		return PartitionTypePageBreak, nil
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
	Items      []Partition   `json:"items,omitempty"`
	AltText    string        `json:"altText,omitempty"`
	Link       string        `json:"link,omitempty"`
	Relation   string        `json:"relation,omitempty"`
	Color      string        `json:"color,omitempty"`
}
