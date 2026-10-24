package domain

type UserRequest struct {
	ID       string `json:"id"`
	Login    string `json:"login"`
	Role     string `json:"role"`
	Password string `json:"password"`
	Name     string `json:"name"`
	Email    string `json:"email"`
}

type UserResponse struct {
	ID    string `json:"id"`
	Login string `json:"login"`
	Name  string `json:"name"`
	Email string `json:"email"`
}