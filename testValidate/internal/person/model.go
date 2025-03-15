package person

type Person struct {
	Id    int    `json:"id"`
	Name  string `json:"name" validate:"required,min=3"`
	Email string `json:"email" validate:"required,email"`
}

func NewPerson() *Person {
	return &Person{}
}
