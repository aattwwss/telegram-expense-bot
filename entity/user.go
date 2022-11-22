package entity

type User struct {
	Id        int64
	IsBot     bool
	FirstName string
	LastName  *string
	Username  *string
}
