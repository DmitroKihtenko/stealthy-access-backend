package api

type AdminUser struct {
	Username string `json:"username"`
}

type User struct {
	Username  string `json:"username" validate:"required,username"`
	Password  string `json:"password" validate:"required,password"`
	FirstName string `json:"first_name" validate:"required,min=1,max=24"`
	LastName  string `json:"last_name" validate:"required,min=1,max=24"`
	Email     string `json:"email" validate:"required,email"`
}
