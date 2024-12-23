package entity

type User struct {
	ID          uint
	PhoneNumber string
	Name        string
	Password    string // the password is hashed
}
