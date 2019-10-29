package appuser

// User are the details of a gameuser.
type User struct {
	ID    int64  `json:"-"`
	GUID  string `json:"id"`
	Email string `json:"email" bson:"email"`
}
