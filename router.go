package main
	
import (
	"github.com/julienschmidt/httprouter"
	"net/http"
	"log"
	"crud/controller"
)

func router() {
	router := httprouter.New()

		/********* User routes **********/
	userControllerInstance := controller.UserController{}
	router.GET("/users/", userControllerInstance.GetAllUsers)
	router.GET("/user/:id", userControllerInstance.GetUserById)
	router.POST("/user/create", userControllerInstance.CreateUser)
	router.PUT("/user/:id", userControllerInstance.UpdateUser)
	router.DELETE("/user/:id", userControllerInstance.DeleteUser)

	log.Fatal(http.ListenAndServe(":8080", router))
}