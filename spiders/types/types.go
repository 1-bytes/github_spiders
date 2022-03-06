package types

type User struct {
	Token string
}

type GitHubUser struct {
	Users []User
}

const DefaultPerPage = 30
