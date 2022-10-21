package models

import (
	"context"
	"crud-go/utils"
	"fmt"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
)

type user struct {
	Id       uint   `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	City     string `json:"city"`
	Name     string `json:"name"`
}

func ConnectionMongo() *mongo.Client {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		panic(err)
	}
	return client
}

func SaveUser(user utils.RegisterType, client *mongo.Client) {
	usersCollection := client.Database(utils.Db).Collection("users")
	_, err := usersCollection.InsertOne(context.TODO(), &user)
	if err != nil {
		panic(err)
	}

}

func LoginCheckMongo(client *mongo.Client, input utils.RegisterType) string {

	usersCollection := client.Database(utils.Db).Collection("users")
	u := user{}
	if err := usersCollection.FindOne(context.TODO(), bson.M{"username": input.Username}).Decode(&u); err != nil {
		return "error"
	}

	if ok := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(input.Password)); ok != nil {
		fmt.Print(ok)
		return "error"
	}

	token, _ := utils.GenerateToken(u.Id)
	return token

}

func GetUserMongo(username string) user {
	u := user{}
	client := ConnectionMongo()

	usersCollection := client.Database(utils.Db).Collection("users")
	if err := usersCollection.FindOne(context.TODO(), bson.M{"username": username}).Decode(&u); err != nil {
		panic("error in getting")
	}

	return u
}

func UpdateMongo(name string, input utils.RegisterType) {

	filter := bson.D{
		{Key: "$and",
			Value: bson.A{
				bson.D{
					{Key: "name", Value: input.Name}},
			},
		},
	}

	update := bson.D{
		{Key: "$set",
			Value: bson.D{
				{Key: "city", Value: input.City},
				{Key: "username", Value: input.Username},
			},
		},
	}

	client := ConnectionMongo()
	usersCollection := client.Database(utils.Db).Collection("users")
	_, err := usersCollection.UpdateOne(context.TODO(), filter, update)

	// check for errors in the updating
	if err != nil {
		panic(err)
	}

}

func DeleteMongo(name string) {
	client := ConnectionMongo()
	filter := bson.D{
		{Key: "$and",
			Value: bson.A{
				bson.D{
					{Key: "name", Value: name}},
			},
		},
	}

	usersCollection := client.Database("testing").Collection("users")

	_, err := usersCollection.DeleteOne(context.TODO(), filter)
	// check for errors in the deleting

	if err != nil {
		panic(err)
	}
}
