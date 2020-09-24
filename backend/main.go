package main

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"github.com/joho/godotenv"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"log"
	"net/http"
	"os"
	"strconv"
	"xorm.io/core"
)

type (
	Todo struct {
		Id   int    `json:"id"`
		Name string `json:"name"`
		Todo string `json:"todo"`
	}
)

var engine *xorm.Engine

func main() {
	var err error
	err = godotenv.Load(".env")
	if err != nil {
		log.Println(err)
	}

	driverName := os.Getenv("DRIVER_NAME")
	user := os.Getenv("USER_NAME")
	pass := os.Getenv("PASS")
	tcp := os.Getenv("TCP")
	dbName := os.Getenv("DB_NAME")

	engine, err = xorm.NewEngine(driverName, user+":"+pass+"@"+tcp+"/"+dbName)
	if err != nil {
		log.Println(err)
	}
	defer engine.Close()

	engine.ShowSQL(true)
	engine.SetMapper(core.GonicMapper{})

	e := echo.New()
	e.Use(middleware.CORS())
	e.GET("/todoList", getTodo)
	e.DELETE("/todoList", deleteTodoHandler)
	e.POST("/todoList", saveTodoHandler)
	e.PUT("/todoList", updateTodoHandler)

	e.Logger.Fatal(e.Start(":8000"))
}

func getTodo(c echo.Context) error {

	var todoList []Todo
	err := engine.Find(&todoList)

	if err != nil {
		log.Println(err)
		return err
	}
	log.Println(todoList)

	return c.JSON(http.StatusOK, todoList)
}

func deleteTodoHandler(c echo.Context) error {
	id, _ := strconv.Atoi(c.QueryParam("id"))
	if err := deleteTodo(id); err != nil {
		return c.JSON(http.StatusOK, err)
	}
	return c.JSON(http.StatusOK, id)
}

func saveTodoHandler(c echo.Context) error {
	todo := new(Todo)
	if err := c.Bind(todo); err != nil {
		return err
	}
	if err := saveTodo(todo); err != nil {
		return c.JSON(http.StatusOK, err)
	}
	return c.JSON(http.StatusOK, todo)
}

func updateTodoHandler(c echo.Context) error {
	todo := new(Todo)
	if err := c.Bind(todo); err != nil {
		return err
	}
	if err := updateTodo(todo); err != nil {
		return c.JSON(http.StatusOK, err)
	}
	return c.JSON(http.StatusOK, todo)
}

func deleteTodo(id int) error {
	t := Todo{}
	if _, err := engine.Where("id=?", id).Delete(t); err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func saveTodo(todo *Todo) error {
	if _, err := engine.Insert(todo); err != nil {
		log.Println(err)
		return err
	}
	return nil
}

func updateTodo(todo *Todo) error {
	if _, err := engine.Where("id=?", todo.Id).Update(todo); err != nil {
		log.Println(err)
		return err
	}
	return nil
}
