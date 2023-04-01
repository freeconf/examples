package demo

type Contacts struct {
	Me            *User
	Users         []*User
	Size          int
	HarmlessValue string
}

type User struct {
	FullName string
}
