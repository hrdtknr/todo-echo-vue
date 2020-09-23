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

	getTodos()

	// Middleware
	e.Use(middleware.Logger())
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())

	// Routes
	e.GET("/todoList", handler)
	e.DELETE("/todoList", deleteTodo)
	e.POST("/todoList", saveTodo)
	e.PUT("/todoList", updateTodo)

	// Start server
	e.Logger.Fatal(e.Start(":8000"))
}

func handler(c echo.Context) error {
	log.Println("handler") // insert ここにきてる
	return c.JSON(http.StatusOK, todoList)
}

func deleteTodo(c echo.Context) error {
	log.Println("delete")
	id, _ := strconv.Atoi(c.QueryParam("id"))
	del(id)
	return c.NoContent(http.StatusNoContent)
}

func saveTodo(c echo.Context) error {
	log.Println("save")
	log.Println("\n\nc", c)
	t := new(Todo)
	if err := c.Bind(t); err != nil {
		return err
	}
	return c.JSON(http.StatusOK, t)
}

func updateTodo(c echo.Context) error {
	todo := new(Todo)
	if err := c.Bind(todo); err != nil {
		return err
	}
	log.Println("todo", todo)
	upd(todo)
	return nil // 仮
}

func getTodos() {
	var err error
	engine, err := xorm.NewEngine("mysql", "root:1234@tcp(127.0.0.1:3306)/go")
	if err != nil {
		log.Println("sippai", err)
		return
	}
	defer engine.Close()

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
}

func del(id int) {
	t := Todo{}
	log.Println("engine", engine) //渡ってない
	affected, err := engine.Where("id=?", id).Delete(&t)
	log.Println("del3")
	if err != nil {
		log.Println(affected, err)
	}
	log.Println("del4")
}

func upd(todo *Todo) {
	log.Println("engine", engine)
	affected, err := engine.Where("id=?", todo.Id).Update(&todo)
	if err != nil {
		log.Println(affected, err)
	}
}
