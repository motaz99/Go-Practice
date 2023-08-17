package controllers

import (
	"bytes"
	"database/sql"
	"fmt"

	"net/http"

	"github.com/gin-gonic/gin"
)

type PersonController struct {
	DB *sql.DB
}

type Person struct {
	Id         int
	First_Name string
	Last_Name  string
}

// Constructor function to create a new PersonController instance
func NewPersonController(db *sql.DB) *PersonController {
	return &PersonController{DB: db}
}

func (pc *PersonController) GetPersonByID(c *gin.Context) {
	var (
		person Person
		result gin.H
	)
	id := c.Param("id")
	row := pc.DB.QueryRow("select id, first_name, last_name from person where id = ?;", id)
	err := row.Scan(&person.Id, &person.First_Name, &person.Last_Name)
	if err != nil {
		// If no results send null
		result = gin.H{
			"result": nil,
			"count":  0,
		}
	} else {
		result = gin.H{
			"result": person,
			"count":  1,
		}
	}
	c.JSON(http.StatusOK, result)
}

func (pc *PersonController) GetAllPersons(c *gin.Context) {
	var (
		person  Person
		persons []Person
	)
	rows, err := pc.DB.Query("select id, first_name, last_name from person;")
	if err != nil {
		fmt.Print(err.Error())
	}
	for rows.Next() {
		err = rows.Scan(&person.Id, &person.First_Name, &person.Last_Name)
		persons = append(persons, person)
		if err != nil {
			fmt.Print(err.Error())
		}
	}
	defer rows.Close()
	c.JSON(http.StatusOK, gin.H{
		"result": persons,
		"count":  len(persons),
	})
}

func (pc *PersonController) CreatePerson(c *gin.Context) {

	var buffer bytes.Buffer
	first_name := c.PostForm("first_name")
	last_name := c.PostForm("last_name")
	stmt, err := pc.DB.Prepare("insert into person (first_name, last_name) values(?,?);")
	if err != nil {
		fmt.Print(err.Error())
	}
	_, err = stmt.Exec(first_name, last_name)

	if err != nil {
		fmt.Print(err.Error())
	}

	// Fastest way to append strings
	buffer.WriteString(first_name)
	buffer.WriteString(" ")
	buffer.WriteString(last_name)
	defer stmt.Close()
	name := buffer.String()
	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf(" %s successfully created", name),
	})
}

func (pc *PersonController) UpdatePerson(c *gin.Context) {
	var buffer bytes.Buffer
	id := c.Query("id")
	first_name := c.PostForm("first_name")
	last_name := c.PostForm("last_name")
	stmt, err := pc.DB.Prepare("update person set first_name= ?, last_name= ? where id= ?;")
	if err != nil {
		fmt.Print(err.Error())
	}
	_, err = stmt.Exec(first_name, last_name, id)
	if err != nil {
		fmt.Print(err.Error())
	}

	// Fastest way to append strings
	buffer.WriteString(first_name)
	buffer.WriteString(" ")
	buffer.WriteString(last_name)
	defer stmt.Close()
	name := buffer.String()
	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("Successfully updated to %s", name),
	})
}

func (pc *PersonController) DeletePerson(c *gin.Context) {

	id := c.Query("id")
	stmt, err := pc.DB.Prepare("delete from person where id= ?;")
	if err != nil {
		fmt.Print(err.Error())
	}
	_, err = stmt.Exec(id)
	if err != nil {
		fmt.Print(err.Error())
	}
	c.JSON(http.StatusOK, gin.H{
		"message": fmt.Sprintf("Successfully deleted user: %s", id),
	})
}
