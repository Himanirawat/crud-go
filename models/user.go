package models

import (
	"crud-go/utils"
	"fmt"
	_ "log"
	_ "os"

	_ "github.com/go-sql-driver/mysql"
	"golang.org/x/crypto/bcrypt"

	"github.com/jinzhu/gorm"
)

type User struct {
	gorm.Model
	Username string `json:"username"`
	Password string `json:"password"`
	City     string `json:"city"`
	Name     string `json:"name"`
}

func Connection() (db *gorm.DB) {
	username := "root"
	password := ""
	dbname := "crud"
	dbHost := "localhost"
	URL := fmt.Sprintf("%s:%s@tcp(%s)/%s", username, password,
		dbHost, dbname)
	db, err1 := gorm.Open("mysql", URL+"?parseTime=true")
	if err1 != nil {
		panic(err1.Error())
	}
	db.AutoMigrate(&User{})
	return db
}

func (u *User) Save(db *gorm.DB) {
	err := db.Create(&u).Error
	if err != nil {
		panic("stop")
	}

}

func LoginCheck(db *gorm.DB, input utils.RegisterType) string {
	u := User{}
	if err := db.Model(User{}).Where("username=?", input.Username).Take(&u).Error; err != nil {
		fmt.Print(err)
		return "error"
	}

	if ok := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(input.Password)); ok != nil {
		fmt.Print(ok)
		return "error"
	}

	token, _ := utils.GenerateToken(u.ID)
	return token

}

func GetUser(id string) User {
	u := User{}
	db := Connection()

	if err := db.Model(User{}).Where("id=?", id).Take(&u).Error; err != nil {
		panic(err)
	}
	return u
}

func Delete(id string) User {
	u := User{}
	db := Connection()

	if err := db.Model(User{}).Where("id=?", id).Delete(&u).Error; err != nil {
		panic(err)
	}
	return u
}

func Update(id string, input utils.RegisterType) User {
	u := User{}
	u.Username = input.Username
	hashed, _ := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	u.Password = string(hashed)
	db := Connection()

	if err := db.Model(User{}).Where("id=?", id).Update(&u).Error; err != nil {
		panic(err)
	}
	return u
}
