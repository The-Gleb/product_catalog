package entity

type User struct {
	ID       int64
	Login    string
	Password string
}

type Credentials struct {
	Login    string
	Password string
}
