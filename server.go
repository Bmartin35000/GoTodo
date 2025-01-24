package main

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/go-chi/chi/middleware"
	"github.com/go-chi/chi/v5"

	_ "github.com/lib/pq"
	"github.com/thedevsaddam/renderer"
)

var rnd *renderer.Render
var db *sql.DB

func helloHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello!")
}

func aboutHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "About Page")
}

func main() {
	fmt.Println("main")
	router := chi.NewRouter()
	router.Use(middleware.Logger)
	router.Mount("/", todoHandlers())

	server := &http.Server{
		Addr:         ":8080",
		Handler:      router,
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
	}

	// start the server
	fmt.Println("Server started on port", 8080)
	if err := server.ListenAndServe(); err != nil {
		log.Printf("listen:%s\n", err)

	}
}

func init() {
	fmt.Println("init function running")

	rnd = renderer.New()
	var err error

	connStr := "postgresql://postgres:postgres@localhost:55555/todos?sslmode=disable"
	// Connect to database
	db, err = sql.Open("postgres", connStr)
	checkError(err)
}

func todoHandlers() http.Handler {
	fmt.Println("handler")
	router := chi.NewRouter()
	router.Group(func(r chi.Router) {
		r.Get("/", getTodos)
		//r.Post("/", createTodo)
		//r.Put("/{id}", updateTodo)
		//r.Delete("/{id}", deleteTodo)
	})
	return router
}

func getTodos(rw http.ResponseWriter, r *http.Request) {
	fmt.Println("get todos")

	rows, err := db.Query("SELECT * FROM todos")
	checkError(err)
	defer rows.Close()
	var todoListFromDB []TodoModel
	for rows.Next() {
		var todoModel TodoModel
		err := rows.Scan(&todoModel.ID, &todoModel.Title, &todoModel.Completed, &todoModel.CreatedAt)
		checkError(err)
		todoListFromDB = append(todoListFromDB, todoModel)
	}

	fmt.Println("Todos")
	fmt.Println(todoListFromDB)

	if err != nil {
		log.Printf("failed to fetch todo records from the db: %v\n", err.Error())
		rnd.JSON(rw, http.StatusBadRequest, renderer.M{
			"message": "Could not fetch the todo collection",
			"error":   err.Error(),
		})

		return
	}

	todoList := []Todo{}

	for _, td := range todoListFromDB {
		todoList = append(todoList, Todo{
			Title:     td.Title,
			ID:        td.ID,
			CreatedAt: td.CreatedAt,
			Completed: td.Completed,
		})
	}
	rnd.JSON(rw, http.StatusOK, GetTodoResponse{
		Message: "All todos retrieved",
		Data:    todoList,
	})
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
