package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"github.com/gofiber/fiber/v2"
	//"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/joho/godotenv"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options" //options package is used to configure how your Go app connects to MongoDB.
)

type Todo struct {
	ID        primitive.ObjectID `json:"_id,omitempty" bson:"_id,omitempty"` //beacause and mongo db store date in bson format (binary json),omitempty-> “If this field has an empty (zero) value, don’t include it in the output.”
	Completed bool               `json:"completed"`
	Body      string             `json:"body"`
}

var collection *mongo.Collection

func main() {
	fmt.Println("hello world")
	if os.Getenv("ENV") != "production"{
		//load only when ENV == development	
		err := godotenv.Load(".env")
		if err != nil {
			log.Fatal("Error loading .env file", err)
		}
	}

	MONGODB_URI := os.Getenv("MONGODB_URI")

	clientOptions := options.Client().ApplyURI(MONGODB_URI)
	client, err := mongo.Connect(context.Background(), clientOptions) // := use only when we delaring new variable, context.Background() ->timer + control button, “Start the work and don’t stop unless I say so.”
	if err != nil {
		log.Fatal(err)
	}

	defer client.Disconnect(context.Background()) //“Do this at the END”
	err = client.Ping(context.Background(), nil)  //to check if the connection is alive
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB!")
	collection = client.Database("golang_db").Collection("todos")
	app := fiber.New()

	// app.Use(cors.New(cors.Config{ //not req for production because FE and BE must serve under same domain
	// 	AllowHeaders:"Origin, Content-Type,Accept", 
	// 	AllowOrigins: "http://localhost:5173",
	// 	AllowMethods: "GET,POST,PUT,DELETE,PATCH",
	// }))

	app.Get("/api/todos", getTodos)
	app.Post("/api/todos", createTodos)
	app.Patch("/api/todos/:id", updateTodos)
	app.Delete("/api/todos/:id", deleteTodos)

	PORT := os.Getenv("PORT")
	if PORT == "" {
		PORT = "5000" // Default port if not specified in .env
	}
	if os.Getenv("ENV") =="production"{
		app.Static("/","client/dist")
	}
	log.Fatal(app.Listen("0.0.0.0:" + PORT)) //0.0.0.0 - listen on all available interfaces
}

func getTodos(c *fiber.Ctx) error {
	var todos []Todo
	cursor, err := collection.Find(context.Background(), bson.M{}) // cursor is pointer to the result set

	if err != nil {
		return err
	}
	defer cursor.Close(context.Background())
	for cursor.Next(context.Background()) {
		var todo Todo
		err := cursor.Decode(&todo)
		if err != nil {
			return err
		}
		todos = append(todos, todo)
	}
	return c.JSON(todos)
}

func createTodos(c *fiber.Ctx) error {
	todo := new(Todo) //return pointer to Todo struct
	//{id:0,completed:false, body: ""}
	err := c.BodyParser(todo)
	if err != nil {
		return err

	}
	if todo.Body == "" {
		return c.Status(400).JSON(fiber.Map{
			"error": "Body is required",
		})
	}
	insertResult, err := collection.InsertOne(context.Background(), todo) //
	if err != nil {
		return err
	}
	todo.ID = insertResult.InsertedID.(primitive.ObjectID) //type assertion is used to convert interface{} to primitive.ObjectID
	return c.Status(201).JSON(todo)

}

func updateTodos(c *fiber.Ctx) error {
	id := c.Params("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		c.Status(400).JSON(fiber.Map{
			"error": "Invalid ID",
		})
	}
	filter := bson.M{"_id": objectID}
	var todo Todo
	collection.FindOne(context.Background(), filter).Decode(&todo)

	update := bson.M{
		"$set": bson.M{
			"completed": !todo.Completed,
		},
	}
	_, err = collection.UpdateOne(context.Background(), filter, update)
	if err != nil {
		return err
	}
	return c.Status(200).JSON(fiber.Map{
		"message": "Todo updated successfully",
	})
}

func deleteTodos(c *fiber.Ctx) error {
	id := c.Params("id")
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return c.Status(400).JSON(fiber.Map{
			"error": "Invalid ID",
		})
	}
	filter := bson.M{"_id": objectID}
	_, err = collection.DeleteOne(context.Background(), filter)
	if err != nil {
		return err
	}
	return c.Status(200).JSON(fiber.Map{"message": "Deleted successfuly"})

}
