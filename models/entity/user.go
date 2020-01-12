package entity

// User represents object user
type User struct {
	ID       int64  `json:"id"`
	Email    string `json:"email"`
	Password string `json:"password,omitempty"`
	Address  string `json:"address"`
}
