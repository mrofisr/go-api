package model

type Person struct {
	ID       int    `json:"id"`
	Name     string `json:"name"`
	Age      int    `json:"age"`
	Username string `json:"username"`
	Password string `json:"password"`
}
