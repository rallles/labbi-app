package models

type Dog struct {
	ID     string `json:"id"`
	Name   string `json:"name"`
	Born   string `json:"born"`
	Gender string `json:"gender"`
}

type Buyer struct {
	ID      string `json:"id"`
	Name    string `json:"name"`
	Contact string `json:"contact"`
}

type Purchase struct {
	Date  string `json:"date"`
	Price int    `json:"price"`
}

type Puppy struct {
	Name        string
	Birthdate   string
	Gender      string
	Description string
	Images      []string
}
