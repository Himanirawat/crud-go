package utils

type RegisterType struct {
	Username string `json:"username"  validate:"required,min=3,max=12"`
	Password string `json:"password"  validate:"required,min=3,max=12"`
	City     string `json:"city"  validate:"required,min=3,max=12"`
	Name     string `json:"name"  validate:"required,min=3,max=12"`
	Email    string `json:"email"   validate:"required,min=3,max=12"`
}

var (
	Db = "testing"
)
