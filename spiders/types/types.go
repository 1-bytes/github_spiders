package types

type User struct {
	Token string
}

type GitHubUser struct {
	Users []User
}

// type JsonBody map[string]interface{}
//
// type SliceBody []JsonBody
