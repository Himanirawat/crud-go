package models

import (
	"context"
	"crud-go/utils"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"golang.org/x/crypto/bcrypt"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type user struct {
	Id       primitive.ObjectID   `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
	City     string `json:"city"`
	Name     string `json:"name"`
	Email    string  `json:"email"`
}

func ConnectionMongo() *mongo.Client {
	client, err := mongo.Connect(context.TODO(), options.Client().ApplyURI("mongodb://localhost:27017"))
	if err != nil {
		panic(err)
	}
	return client
}

func SaveUser(user utils.RegisterType, client *mongo.Client) bool{
	u :=User{}
	usersCollection := client.Database(utils.Db).Collection("users")
	fmt.Print(user)
	if err := usersCollection.FindOne(context.TODO(), bson.M{"email": user.Email}).Decode(&u); err != nil {
		_, err := usersCollection.InsertOne(context.TODO(), &user)
	if err != nil {
		panic(err)
	}
	}

	if u.Name != ""{
		return true
	}
	return false
	

}

func LoginCheckMongo(client *mongo.Client, input utils.RegisterType) string {

	usersCollection := client.Database(utils.Db).Collection("users")
	u := user{}
	if err := usersCollection.FindOne(context.TODO(), bson.M{"email": input.Email}).Decode(&u); err != nil {
		return "error"
	}

	if ok := bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(input.Password)); ok != nil {
		fmt.Print(ok)
		return "error"
	}

	token, _ := utils.CreateJWT(u.Email)
	return token

}

func GetUserMongo(id string) user {
	u := user{}
	fmt.Print("mongo")
	client := ConnectionMongo()
	objId, _ := primitive.ObjectIDFromHex(id)
	usersCollection := client.Database(utils.Db).Collection("users")
	if err := usersCollection.FindOne(context.TODO(), bson.M{"_id": objId}).Decode(&u); err != nil {
		panic("error in getting")
	}

	return u
}

func UpdateMongo(id string, input utils.RegisterType) {

	objId, _ := primitive.ObjectIDFromHex(id)
	filter := bson.D{
		{Key: "$and",
			Value: bson.A{
				bson.D{
					{Key: "_id", Value: objId}},
			},
		},
	}

	update := bson.D{
		{Key: "$set",
			Value: bson.D{
				{Key: "city", Value: input.City},
				{Key: "username", Value: input.Username},
				{Key: "name", Value: input.Name},
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

func DeleteMongo(id string) {
	client := ConnectionMongo()
	objId, _ := primitive.ObjectIDFromHex(id)

	filter := bson.D{
		{Key: "$and",
			Value: bson.A{
				bson.D{
					{Key: "_id", Value: objId}},
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
