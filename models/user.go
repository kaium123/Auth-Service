package models

type Attachment struct {
	Path string `json:"path"`
	Name string `json:"name"`
}

type User struct {
	ID         int      `json:"id"`
	Email      string   `json:"email" `
	Password   string   `json:"password,omitempty"`
	Name       string   `json:"name"`
	UserName   string   `json:"user_name" `
	Phone      string   `json:"phone"`
	Websites   []string `json:"websites"`
	Bio        string   `json:"bio"`
	Gender     string   `json:"gender"`
	ProfilePic string   `json:"profile_pic"`
}

type SignInData struct {
	Email    string `json:"email"`
	Password string `json:"password,omitempty"`
}
