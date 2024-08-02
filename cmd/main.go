package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"messenger/internal/adapters/handlers"
	"messenger/internal/adapters/repositories"
	"messenger/internal/core/services"
)

var (
	repo                 = flag.String("db", "mongo", "Database for storing messages")
	httpHandlerMessanger *handlers.HTTPHandlerMessanger
	svcMessanger         *services.MessangerService
	HTTPHandlerUser      *handlers.HTTPHandlerUser
	svcUser              *services.UserService
)

func main() {
	flag.Parse()

	fmt.Printf("Application running using %s\n", *repo)
	switch *repo {
	case "mongo":
		storeMessanger := repositories.NewMessangerMongoRepository()
		svcMessanger = services.NewMessangerService(storeMessanger)
		storeUser := repositories.NewUserMongoRepository()
		svcUser = services.NewUserService(storeUser)
	default:
		storeMessanger := repositories.NewMessangerPostgresRepository()
		svcMessanger = services.NewMessangerService(storeMessanger)
		storeUser := repositories.NewUserPostgresRepository()
		svcUser = services.NewUserService(storeUser)
	}

	InitRoutes()
}

func InitRoutes() {
	router := gin.Default()
	handlerMessanger := handlers.NewHTTPHandlerMessanger(*svcMessanger)
	handlerUser := handlers.NewHTTPHandlerUser(*svcUser)

	router.GET("/users/export-data", handlerUser.GetAllUsersByExportData)
	router.GET("/users", handlerUser.GetAllUsers)
	router.GET("/user/:id", handlerUser.GetOneUser)
	router.PUT("/user/:id", handlerUser.UpdateUser)
	router.DELETE("/user/:id", handlerUser.DeleteUser)
	router.POST("/register", handlerUser.RegisterUser)
	router.POST("/login", handlerUser.LoginUser)

	router.GET("/messages", handlerMessanger.GetAllMessages)
	router.GET("/message/:id", handlerMessanger.GetOneMessage)
	router.POST("/messages", handlerMessanger.CreateMessage)
	router.PUT("/message/:id", handlerMessanger.UpdateMessage)
	router.DELETE("/message/:id", handlerMessanger.DeleteMessage)

	port := "5000"

	server := &http.Server{
		Addr:    ":" + port,
		Handler: router,
	}

	log.Printf("Server is starting on ports %s...\n", port)
	if err := server.ListenAndServe(); err != nil {
		log.Fatalf("Server error: %v", err)
	}
}
