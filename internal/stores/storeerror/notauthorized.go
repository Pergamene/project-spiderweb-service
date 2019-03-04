package storeerror

import "fmt"

// NotAuthorized is an error that signifies that the associated db action is not permitted.
type NotAuthorized struct {
	UserID  string
	TableID string
	Err     error
}

func (e *NotAuthorized) Error() string {
	if e.Err != nil {
		return fmt.Sprintf("User %v is not authorized to perform the action on the ID %v\n%v", e.UserID, e.TableID, e.Err)
	}
	return fmt.Sprintf("User %v is not authorized to perform the action on the ID %v", e.UserID, e.TableID)
}
