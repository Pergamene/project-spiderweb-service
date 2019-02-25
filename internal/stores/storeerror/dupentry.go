package storeerror

import "fmt"

// DupEntry is an error that signifies that the item is a duplicate entry in the store.
type DupEntry struct {
	ID  string
	Err error
}

func (e *DupEntry) Error() string {
	return fmt.Sprintf("Duplicate id: %v\n%v", e.ID, e.Err)
}
