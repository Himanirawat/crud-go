package controllers

import (
	"context"
	"crud-go/models"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	elastic "github.com/olivere/elastic/v7"
)

type ElasticDocs struct {
	Username string `json:"username"`
	Password string `json:"password"`
	City     string `json:"city"`
	Name     string `json:"name"`
}

// elastic functionality is made with my sql(***please note)
func InsertElastic(c *gin.Context) {
	users := []models.User{}

	ctx := context.Background()
	esclient, err := GetESClient()
	if err != nil {
		fmt.Println("Error initializing : ", err)
		panic("Client fail ")
	}
	db := models.Connection()
	db.Find(&users)

	for _, user := range users {
		doc1 := ElasticDocs{}
		doc1.Username = user.Username
		doc1.Password = user.Password
		doc1.City = user.City
		doc1.Name = user.Name
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
	c.JSON(http.StatusOK, gin.H{

		"message": "[Elastic][InsertProduct]Insertion Successful"})

}

func GetESClient() (*elastic.Client, error) {

	client, err := elastic.NewClient(elastic.SetURL("http://localhost:9200"),
		elastic.SetSniff(false),
		elastic.SetHealthcheck(false))

	fmt.Println("ES initialized...")

	return client, err

}

func Get(c *gin.Context) { //for getting documents values
	input := ElasticDocs{}
	ctx := context.Background()
	esclient, err := GetESClient()
	if err != nil {
		fmt.Println("Error initializing : ", err)
		panic("Client fail ")
	}
	c.BindJSON(&input)
	var users []ElasticDocs

	searchSource := elastic.NewSearchSource()
	searchSource.Query(elastic.NewMatchQuery("name", input.Name))

	/* this block will basically print out the es query */
	queryStr, err1 := searchSource.Source()
	queryJs, err2 := json.Marshal(queryStr)

	if err1 != nil || err2 != nil {
		fmt.Println("[esclient][GetResponse]err during query marshal=", err1, err2)
	}
	fmt.Println("[esclient]Final ESQuery=\n", string(queryJs))
	/* until this block */

	searchService := esclient.Search().Index("users").SearchSource(searchSource)

	searchResult, err := searchService.Do(ctx)
	if err != nil {
		fmt.Println("[ProductsES][GetPIds]Error=", err)
		return
	}

	for _, hit := range searchResult.Hits.Hits {
		var user ElasticDocs
		//fmt.Print(hit.Source)
		err := json.Unmarshal(hit.Source, &user)
		if err != nil {
			fmt.Println("[Getting users][Unmarshal] Err=", err)
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
}
