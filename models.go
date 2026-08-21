package main

type User struct {
	Id   int
	Name string
	Age  int
	City string
}

type Item struct {
	Id          int
	Title       string
	UserId      int
	Description string
	Price       string
	Status      string
}

type UserProfile struct {
	User
	Items []Item
}

type GuestbookInfo struct {
	Count int
	Users []User
}

type ValidationRule struct {
	FieldName string
	Value     string
	Validate  func(string) error
}
