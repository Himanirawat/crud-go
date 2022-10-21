package utils

type RegisterType struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	City     string `json:"city"`
	Name     string `json:"name"`
}

var (
	Db = "testing"
)
