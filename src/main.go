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
	//	e.Use(middleware.CORS())//AccessControl
	/*
	   e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
	   		AllowOrigins: []string{"http://localhost:8000/", "http://localhost:8080/"},
	   		AllowMethods: []string{http.MethodGet, http.MethodPut, http.MethodPost, http.MethodDelete},
	   	}))
	*/
	/*
		e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
	    AllowOrigins:     []string{"*"},
	    AllowHeaders:     []string{"authorization", "Content-Type"},
	    AllowCredentials: true,
	    AllowMethods:     []string{echo.OPTIONS, echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE},
	}))
	*/
	e.Use(func(h echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			req := c.Request()

			headers := c.Response().Header()
			headers.Set("Access-Control-Allow-Origin", "*")
			headers.Set("Access-Control-Allow-Headers", "Authorization")
			headers.Set("Access-Control-Allow-Methods", "GET, PATCH, PUT, POST, DELETE, OPTIONS")

			if "OPTIONS" != req.Method {
				h(c)
			}

			return nil
		}
	})
	// Routes
	e.GET("/todoList", handler)
	e.DELETE("/todoList", deleteTodo)
	e.POST("/todoList", saveTodo) // "/"にしたけど違いそう
	// Start server
	e.Logger.Fatal(e.Start(":8000"))
}

func handler(c echo.Context) error {
	log.Println("handler") // insert delete ここにきてる
	return c.JSON(http.StatusOK, todoList)
}

func deleteTodo(c echo.Context) error {
	log.Println("delete")
	//	log.Println(c.Param("Id"))
	id2 := c.QueryParam("id") //クエリパラメータで送ってるから受け取り方がc.Paramではない
	log.Println("query id2:", id2)

	id, _ := strconv.Atoi(c.Param("Id"))
	log.Println("id:", id)
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
