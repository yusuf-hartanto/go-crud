package main

import (
	"fmt"
	"net/http"
	TaskController "project/controllers"

	Model "project/models"

	"github.com/julienschmidt/httprouter"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func main() {
	db, err := gorm.Open(sqlite.Open("database.db"), &gorm.Config{})
	if err != nil {
		panic("Failed to conect database")
	}
	db.AutoMigrate(&Model.Tasks{})

	router := httprouter.New()

	router.ServeFiles("/static/*filepath", http.Dir("assets"))

	router.GET("/", TaskController.Index)
	router.GET("/create", TaskController.Create)
	router.POST("/create", TaskController.Create)
	router.GET("/show/:id", TaskController.Show)
	router.GET("/update/:id", TaskController.Update)
	router.POST("/update/:id", TaskController.Update)
	router.GET("/delete/:id", TaskController.Delete)

	fmt.Println("http://localhost:8080")
	http.ListenAndServe(":8080", router)
}
