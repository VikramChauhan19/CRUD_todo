package main

import (
	"fmt"
	"log"
	"github.com/gofiber/fiber/v2"
	"github.com/joho/godotenv"
	"os"
)

type Todo struct {
	ID int `json:"id"`
	Completed bool `json:"completed"`
	Body string `json:"body"`
}

func main(){
	fmt.Println("hello world!")
	app := fiber.New() //creating web server

	err := godotenv.Load(".env") // load .env here
	if err != nil{
		log.Fatal("Error Loading .env file")
	}

	PORT := os.Getenv("PORT") // u can access the env variable like this

	todos := []Todo{} // in memory db for todo



	app.Get("/",func(c *fiber.Ctx) error{//ctx0>context , c *fiber.Ctx = request + response handler in Fiber.
		return c.Status(200).JSON(todos) //200->status ok
	})

	app.Post("/api/todos", func (c *fiber.Ctx) error{
		todo := Todo{} // it will have defaoult values  that is ID:0, Completed:false, Body:""
		err := c.BodyParser(&todo) // parsing the request body into todo struct
		if err != nil{
			return err
		}
		if todo.Body == ""{
			return c.Status(400).JSON(fiber.Map{"error":"body is required"}) // 400->bad request
		}
		todo.ID = len(todos)+1
		todos = append(todos,todo)
 
		return c.Status(201).JSON(todo) // status 201 -> created

	})
	//Update todo
	app.Patch("api/todos/:id",func(c *fiber.Ctx) error{
		id := c.Params("id")
		for i,todo := range todos{
			if fmt.Sprint(todo.ID) == id{ //Sprint conv val to string and return without printing
				todos[i].Completed = !todos[i].Completed
				return c.Status(200).JSON(todos[i])  //200 ->status ok
			}
		}
		return c.Status(404).JSON(fiber.Map{"error":"todo not found"}) // 404->not found
	})
	
	//Delete
	app.Delete("api/todos/:id", func(c *fiber.Ctx) error{
		id := c.Params("id")
		for i,todo := range todos{
			if fmt.Sprint(todo.ID) == id{
				deletedTodo := todos[i]
				todos = append(todos[:i],todos[i+1:]...)
				return c.Status(200).JSON(deletedTodo) //200 ->status ok
			}
		}
		return c.Status(404).JSON(fiber.Map{"error":"todo not found"}) // 404->not found
	})

	log.Fatal(app.Listen(":"+PORT))
}