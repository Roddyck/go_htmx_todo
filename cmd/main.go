package main

import (
	"html/template"
	"io"
	"log"
	"os"
	"strconv"

	"github.com/Roddyck/go-htmx-todo/pkg/database"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

type Templates struct {
    templates *template.Template
}

func (t *Templates) Render(w io.Writer, name string, data interface{}, c echo.Context) error {
    return t.templates.ExecuteTemplate(w, name, data)
}

func newTemplate() *Templates {
    return &Templates {
        templates: template.Must(template.ParseGlob("views/*.html")),
    }
}

type Todos = []database.Todo

type IndexPage struct {
    Todos Todos
}

type FormData struct {
    Values map[string]string
    Errors map[string]string
}

func newFormData() FormData {
    return FormData{
        Values: make(map[string]string),
        Errors: make(map[string]string),
    }
}

func baseIndex(c echo.Context) error {
    return c.Render(200, "index", IndexPage{
        Todos {
            {Id: 1, Title: "Item 1", Body: "Item 1"},
            {Id: 2, Title: "Item 2", Body: "Item 2"},
            {Id: 3, Title: "Item 3", Body: "Item 3"},
        },
    })
}

func Index(c echo.Context) error {
    todos, err := database.GetTodos()
    if err != nil {
        return baseIndex(c)
    }

    return c.Render(200, "index", IndexPage{
        Todos: todos,
    })
}

func Create(c echo.Context) error {
    title := c.FormValue("title")
    body := c.FormValue("body")

    id, err := database.SaveTodo(title, body)
    if err != nil {
        log.Fatalf("couldn't save todo: %v", err)
    }

    todo, err := database.GetTodo(id)
    if err != nil {
        log.Fatalf("couldn't get todo: %v", err)
    }

    c.Render(200, "form", newFormData())
    return c.Render(200, "oob-todo", todo) 
}

func DeleteTodo(c echo.Context) error {
    idStr := c.Param("id") 
    id, err := strconv.Atoi(idStr) 
    if err != nil { 
        return c.String(400, "Invalid id")
    }

    err = database.DeleteTodo(id)
    if err != nil {
        return c.String(404, "Couldn't find todo")
    }

    return c.NoContent(200)
}

func main() {
    url := os.Getenv("DB_URL")
    if url == "" {
        url = "todos.db" 
    }

    err := database.InitDB(url)
    if err != nil {
        log.Fatalf("coundn't initialize db: %v", err)
    }

    e := echo.New()
    e.Use(middleware.Logger())

    e.Renderer = newTemplate()

    e.Static("/images", "images")
    e.Static("/css", "css")
    
    e.GET("/", Index)
    e.POST("/create", Create)

    e.DELETE("/todos/:id", DeleteTodo)

    e.Logger.Fatal(e.Start(":8080"))
}

