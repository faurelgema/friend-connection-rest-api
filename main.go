package main

import (
	"log"
	"net/http"
	"os"

	handlers "friend_connection_rest_api/controller"
	"friend_connection_rest_api/docs"
	"friend_connection_rest_api/utils"
)

// @in header
// @name Authorization
func main() {
	db := utils.CreateConnection()
	r := handlers.Setup(db)
	docs.SwaggerInfo.Title = "Rest API for friend connection"
	docs.SwaggerInfo.Description = "Restful api for friend connection api made by Go-Language and Gin framework"
	docs.SwaggerInfo.Version = "2.0"
	docs.SwaggerInfo.BasePath = "/"
	docs.SwaggerInfo.Schemes = []string{"http"}

	var port string = os.Getenv("PORT")
	if port == "" {
		port = "3000"
	}
	log.Println("Server started on: http://localhost:" + port)
	http.ListenAndServe(":"+port, r)
}
