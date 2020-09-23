package main

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/xorm"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"log"
	"net/http"
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

var (
	todoList = map[int]*Todo{}
	engine   *xorm.Engine
)

func main() {
	e := echo.New()
	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	// Routes
	e.GET("/todoList", getTodoHandler)
	e.DELETE("/todoList", deleteTodoHandler)
	e.POST("/todoList", saveTodoHandler)
	e.PUT("/todoList", updateTodoHandler)
	// Start server
	e.Logger.Fatal(e.Start(":8000"))
}

// echo method
func getTodoHandler(c echo.Context) error {
	getTodos()
	return c.JSON(http.StatusOK, todoList)
}

func deleteTodoHandler(c echo.Context) error {
	id, _ := strconv.Atoi(c.QueryParam("id"))
	deleteTodo(id)
	return c.NoContent(http.StatusNoContent)
}

func saveTodoHandler(c echo.Context) error {
	todo := new(Todo)
	if err := c.Bind(todo); err != nil {
		return err
	}
	saveTodo(todo)
	return c.JSON(http.StatusOK, todo)
}

func updateTodoHandler(c echo.Context) error {
	todo := new(Todo)
	if err := c.Bind(todo); err != nil {
		return err
	}
	updateTodo(todo)
	return c.JSON(http.StatusOK, todo)
}

// xorm method
func getTodos() error {
	var err error
	engine, err = xorm.NewEngine("mysql", "root:1234@tcp(127.0.0.1:3306)/go")
	if err != nil {
		log.Println("sippai", err)
		return err
	}

	results, err := engine.Query("SELECT * FROM todo")

	for _, result := range results {
		i, _ := strconv.Atoi(string(result["id"]))
		tmp := &Todo{
			Id:   i,
			Name: string(result["name"]),
			Todo: string(result["todo"]),
		}
		todoList[tmp.Id] = tmp
	}

	engine.ShowSQL(true)
	engine.SetMapper(core.GonicMapper{})
	return err
}

func deleteTodo(id int) error {
	t := Todo{}
	log.Println("id;", id)
	log.Println("before\n", todoList)
	affected, err := engine.Where("id=?", id).Delete(t)
	log.Println("after\n", todoList)
	if err != nil {
		log.Println("error",affected, err)
		return err
	}
	return err
}

func saveTodo(todo *Todo) error {
	affected, err := engine.Insert(todo)
	if err != nil {
		log.Println(affected, err)
		return err
	}
	return err
}

func updateTodo(todo *Todo) error {
	affected, err := engine.Where("id=?", todo.Id).Update(todo)
	if err != nil {
		log.Println(affected, err)
		return err
	}
	return err
}
