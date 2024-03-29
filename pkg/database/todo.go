package database
 
import "fmt"

type Todo struct {
    Id int
    Title string
    Body string 
}

func GetTodo(id int) (*Todo, error) {
    row := Db.QueryRow("SELECT * FROM todos WHERE id = ?", id)

    var title, body string

    err := row.Scan(&id, &title, &body)
    if err != nil {
        return nil, fmt.Errorf("couldn't scan row: %w", err)
    }

    return &Todo{
        Id: id,
        Title: title,
        Body: body,
    }, nil
}

func GetTodos() ([]Todo, error) {
    rows, err := Db.Query("SELECT * FROM todos")
    if err != nil {
        return nil, fmt.Errorf("couldn't get todos: %w", err)
    }

    defer rows.Close()

    todos := []Todo{}

    for rows.Next() {
       var id int
       var title string
       var body string

       err := rows.Scan(&id, &title, &body)
       if err != nil {
           return nil, fmt.Errorf("error scanning row: %w", err)
       }

       todos = append(todos, Todo{
           Id: id,
           Title: title,
           Body: body,
       })
    }

    return todos, nil
}

func SaveTodo(title, body string) (int, error) {
    res, err := Db.Exec("INSERT INTO todos (title, body) VALUES (?, ?) RETURNING id", title, body)
    if err != nil {
        return 0, fmt.Errorf("couldn't save todo: %w", err)
    }

    id, err := res.LastInsertId()
    if err != nil {
        return 0, fmt.Errorf("couldn't get id: %w", err)
    }

    return int(id), nil
}

func DeleteTodo(id int) error {
    _, err := Db.Exec("DELETE FROM todos WHERE id = ?", id)
    if err != nil {
        return fmt.Errorf("error deleting todo: %w", err)
    }
    
    return nil
}
