package err

import (
	"fmt"
)

type UserExists struct {
	Name  string
	Phone string
}

func (u *UserExists) Error() string {
	return fmt.Sprintf("User with username %s already exists", u.Name)
}

type UserNotFound struct {
	Name string
}

func (u *UserNotFound) Error() string {
	return fmt.Sprintf("User with username %s is not found", u.Name)
}

type PasswordMismatch struct {
}

func (p *PasswordMismatch) Error() string {
	return "Passowrd mismatch!"
}
