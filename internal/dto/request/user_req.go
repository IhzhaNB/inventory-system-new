package request

type CreateUserRequest struct {
	Name     string `json:"name" validate:"required,min=3" example:"Staff Satu"`
	Email    string `json:"email" validate:"required,email" example:"staff@gmail.com"`
	Password string `json:"password" validate:"required,min=6" example:"password123"`
	Role     string `json:"role" validate:"required" example:"staff"`
}

type UpdateUserRequest struct {
	Name string `json:"name" validate:"required" example:"Staff Satu Update"`
	Role string `json:"role" validate:"required,oneof=admin staff" example:"admin"`
}
