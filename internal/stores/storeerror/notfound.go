package storeerror

import "fmt"

// NotFound is an error that signifies that the item is not found in the store.
type NotFound struct {
	ID  string
	Err error
}

func (e *NotFound) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("Could not find: %v\n%v", e.ID, e.Err)
	}
	return fmt.Sprintf("Could not find: %v", e.ID)
}
