package controllers

import (
	"crud-go/models"
	"crud-go/utils"
	"fmt"
	"net/http"
	"encoding/json"
	"context"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
	elastic "github.com/olivere/elastic/v7"
	"regexp"
)


// crud operations for mysql
func Register(c *gin.Context) {

	var input utils.RegisterType
	if err := c.BindJSON(&input); err != nil {
		fmt.Printf("%s", err)

		return
	}
	fmt.Print(input)
	u := models.User{}
	u.Username = input.Username
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return
	}
	u.Password = string(hashedPassword)
	u.City = input.City
	u.Name = input.Name

	db := models.Connection()
	u.Save(db)
	c.JSON(http.StatusOK, gin.H{
		"message": "registered",
	})
	db.Close()
}

func Login(c *gin.Context) {
	var input utils.RegisterType
	db := models.Connection()
	if err := c.BindJSON(&input); err != nil {
		fmt.Printf("%s", err)

		return
	}
	token := models.LoginCheck(db, input)
	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})

}

func GetUser(c *gin.Context) {
	if err := utils.TokenValid(c); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}
	id := c.Param("id")

	u := models.GetUser(id)

	c.JSON(http.StatusOK, gin.H{
		"user": u,
	})

}

func Delete(c *gin.Context) {
	if err := utils.TokenValid(c); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}
	id := c.Param("id")

	models.Delete(id)

	c.JSON(http.StatusOK, gin.H{
		"message": "deleted successfully",
	})

}

func Update(c *gin.Context) {
	if err := utils.TokenValid(c); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}
	id := c.Param("id")
	input := utils.RegisterType{}
	c.BindJSON(&input)
	u := models.Update(id, input)

	c.JSON(http.StatusOK, gin.H{
		"user": u,
	})

}

//crud operations for mongo specific

func RegisterUserMongo(c *gin.Context) {

	var input utils.RegisterType
	if err := c.BindJSON(&input); err != nil {
		fmt.Printf("%s", err)

		return
	}
	fmt.Print(input)
	re := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")
	if  !re.MatchString(input.Email){
		c.JSON(http.StatusOK, gin.H{
			"message": "Please enter valid email",
		})
		return
	}
	client := models.ConnectionMongo()

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return
	}
	input.Password = string(hashedPassword)

	check:=models.SaveUser(input, client)
	if check {
		c.JSON(http.StatusOK, gin.H{
			"message": "already there in db with same email",
		})
		return
	}
	if !check{
		fmt.Print("hii")
		InsertElastic(input)
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "registered",
	})
}



func InsertElastic(user utils.RegisterType) {
	ctx := context.Background()
	esclient, err := GetESClient()
	if err != nil {
		fmt.Println("Error initializing : ", err)
		panic("Client fail ")
	}
	
		doc1 := utils.RegisterType{}
		doc1.Username = user.Username
		doc1.Password = user.Password
		doc1.City = user.City
		doc1.Name = user.Name
		doc1.Email = user.Email
		fmt.Print(444)
		fmt.Print(user)
		dataJSON, _ := json.Marshal(doc1)
		js := string(dataJSON)
		ind, err := esclient.Index().
			Index("users").
			BodyJson(js).
			Do(ctx)

		if err != nil {
			panic(ind)
		}
	


}

func GetESClient() (*elastic.Client, error) {

	client, err := elastic.NewClient(elastic.SetURL("http://localhost:9200"),
		elastic.SetSniff(false),
		elastic.SetHealthcheck(false))

	fmt.Println("ES initialized...")

	return client, err

}

func Get(c *gin.Context) { //for getting documents values
	input := utils.RegisterType{}
	ctx := context.Background()
	esclient, err := GetESClient()
	if err != nil {
		fmt.Println("Error initializing : ", err)
		panic("Client fail ")
	}
	c.BindJSON(&input)
	var users []utils.RegisterType

	searchSource := elastic.NewSearchSource()
	searchSource.Query(elastic.NewMatchQuery("name", input.Name))

	/* this block will basically print out the es query */
	queryStr, err1 := searchSource.Source()
	queryJs, err2 := json.Marshal(queryStr)

	if err1 != nil || err2 != nil {
		fmt.Println("Error1=", err1, err2)
	}
	fmt.Println("query=\n", string(queryJs))
	/* until this block */

	searchService := esclient.Search().Index("users").SearchSource(searchSource)

	searchResult, err := searchService.Do(ctx)
	if err != nil {
		fmt.Println("Error2", err)
		return
	}

	for _, hit := range searchResult.Hits.Hits {
		var user utils.RegisterType
		//fmt.Print(hit.Source)
		err := json.Unmarshal(hit.Source, &user)
		if err != nil {
			fmt.Println("Error3", err)
		}

		users = append(users, user)
	}
	fmt.Print(users)
	if err != nil {
		fmt.Println("Fetching user fail: ", err)
	} else {
		for _, s := range users {
			
			c.JSON(http.StatusOK, gin.H{
				"data": "user found Name: " + s.Username + "and password " + s.Password,
			})
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"data": "user not found",
	})
}

func LoginMongo(c *gin.Context) {
	var input utils.RegisterType
	client := models.ConnectionMongo()
	if err := c.BindJSON(&input); err != nil {
		fmt.Printf("%s", err)

		return
	}
	token := models.LoginCheckMongo(client, input)
	c.JSON(http.StatusOK, gin.H{
		"token": token,
	})

}

func GetUserMongo(c *gin.Context) {
	if _,err := utils.VerifyToken(utils.ExtractToken(c)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}
	id := c.Param("id")

	u := models.GetUserMongo(id)

	c.JSON(http.StatusOK, gin.H{
		"user": u,
	})

}

func UpdateMongo(c *gin.Context) {
	if _,err := utils.VerifyToken(utils.ExtractToken(c)); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err,
		})
		return
	}
	

	id := c.Param("id")
	input := utils.RegisterType{}
	c.BindJSON(&input)
	models.UpdateMongo(id, input)

	//elastic updation
	ctx := context.Background()
	client, _ := GetESClient()
    query := elastic.NewMatchQuery("email", input.Email)
	_, err := client.UpdateByQuery().Index("users").Query(query).Script(elastic.NewScriptInline("ctx._source.username = '" + input.Username + "';ctx._source.city = '" + input.City + "';ctx._source.name = '" + input.Name+"'")).Do(ctx)
    if nil != err {
        fmt.Println("Error",err)
    }

	c.JSON(http.StatusOK, gin.H{
		"message": "user updated",
	})

}

func DeleteMongo(c *gin.Context) {
	
	id := c.Param("id")
	input:=models.GetUserMongo(id)  //getting for deleting from elastic
	models.DeleteMongo(id)

	//for elastic deletion
	ctx := context.Background()
	client, _ := GetESClient()
    query := elastic.NewMatchQuery("email", input.Email)
	_, err := client.DeleteByQuery().Index("users").Query(query).Do(ctx)
    if nil != err {
        fmt.Println("Error",err)
    }

	c.JSON(http.StatusOK, gin.H{
		"message": "deleted successfully",
	})

}
