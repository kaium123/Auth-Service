package models

import (
	"errors"

	v "github.com/go-ozzo/ozzo-validation/v4"
	"github.com/go-ozzo/ozzo-validation/v4/is"
)

type Attachment struct {
	Path string `json:"path"`
	Name string `json:"name"`
}

type User struct {
	ID             int      `json:"id"`
	Email          string   `json:"email" `
	Password       string   `json:"password,omitempty"`
	Name           string   `json:"name"`
	UserName       string   `json:"user_name" `
	Phone          string   `json:"phone"`
	Websites       []string `json:"websites"`
	Bio            string   `json:"bio"`
	Gender         string   `json:"gender"`
	ProfilePicName string   `json:"profile_pic_name"`
	ProfilePicPath string   `json:"profile_pic_path"`
}

type SignInData struct {
	Email    string `json:"email"`
	Password string `json:"password,omitempty"`
}

type JWTTokenResponse struct {
	ExpiredAt int64  `json:"expired_at"`
	Token     string `json:"access_token"`
	Refresh   string `json:"refresh_token"`
	Email     string `json:"email,omitempty"`
}

func (u *User) Validate() error {
	return v.ValidateStruct(u,
		v.Field(&u.Email,
			v.Length(0, 20).Error("length must be 0 tot 20"),
			is.EmailFormat.Error("invalid email"),
		),

		v.Field(&u.Name,
			v.By(func(value interface{}) error {
				if len(u.Name) > 20 {
					return errors.New("length must be 0 tot 20")
				}
				return nil
			}),
		),

		v.Field(&u.UserName,
			v.Length(0, 20).Error("length must be 0 tot 20"),
		),
	)
}

func (u *SignInData) Validate() error {
	return v.ValidateStruct(u,
		v.Field(&u.Email,
			v.Length(0, 20).Error("length must be 0 tot 1000"),
			is.EmailFormat.Error("invalid email"),
		),
	)
}
