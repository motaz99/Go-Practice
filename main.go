package main

import (
	"database/sql"
	"fmt"

	"project/controllers"

	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
)

func main() {
	db, err := sql.Open("mysql", "root:Admin@123@tcp(localhost:3306)/goPractice?charset=utf8&parseTime=True&loc=Local")
	if err != nil {
		fmt.Print(err.Error())
	}
	defer db.Close()
	// make sure connection is available
	err = db.Ping()
	if err != nil {
		fmt.Print(err.Error())
	}

	router := gin.Default()

	// Add API handlers here

	personController := controllers.NewPersonController(db)

	router.GET("/person/:id", personController.GetPersonByID)
	router.GET("/persons", personController.GetAllPersons)
	router.POST("/person", personController.CreatePerson)
	router.PUT("/person", personController.UpdatePerson)
	router.DELETE("/person", personController.DeletePerson)

	router.Run(":3000")
}
