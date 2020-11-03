package models

type User struct {
	Id string
	Email string
	Password string
}

type RequestUserInfo struct {
	Email string
	Password string
}

type ResponseUserID struct {
	Id string `json:"id"`
}