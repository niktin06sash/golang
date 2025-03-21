package person

type Person struct {
	Id       int    `json:"id"`
	Name     string `json:"name" validate:"required,min=3"`
	Email    string `json:"email" validate:"required,email"`
	Password string `json:"password" validate:"required,min=8"`
}

func NewPerson() *Person {
	return &Person{}
}
