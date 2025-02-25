package main

import (
	"encoding/json"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/Bmartin35000/backend-project/todo"
	"github.com/go-chi/cors"

	"github.com/go-chi/chi/v5"

	"sync"

	_ "github.com/lib/pq"
	log "github.com/sirupsen/logrus"
	"github.com/thedevsaddam/renderer"
)

var rnd *renderer.Render

func main() {
	router := chi.NewRouter()
	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Accept", "Authorization", "Content-Type", "X-CSRF-Token"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300, // Maximum value not ignored by any of major browsers
	}))
	router.Mount("/", todoHandlers())

	server := &http.Server{
		Addr:         ":" + os.Getenv("server.port"),
		Handler:      router,
		ReadTimeout:  60 * time.Second,
		WriteTimeout: 60 * time.Second,
	}

	// start the server
	log.WithFields(log.Fields{"port": os.Getenv("server.port")}).Info("Server starting")
	if err := server.ListenAndServe(); err != nil {
		log.WithFields(log.Fields{"port": os.Getenv("server.port"), "details": err}).Panic("failed to launch server")
		panic("failed to launch server")
	}
}

func init() {
	// Setting the log's level depending on environment
	env := os.Getenv("environment")
	switch env {
	case "development":
		log.SetLevel(log.InfoLevel)
	case "production":
		log.SetLevel(log.ErrorLevel)
	default:
		log.SetLevel(log.InfoLevel)
	}

	rnd = renderer.New()
}

func todoHandlers() http.Handler {
	router := chi.NewRouter()
	router.Group(func(r chi.Router) {
		r.Get("/todos", getTodos)
		r.Post("/todo", createTodo)
		r.Put("/todo", updateTodo)
		r.Delete("/todo/{id}", deleteTodo)
		r.Get("/fakeTask", fakeTask)
	})
	return router
}

func getTodos(response http.ResponseWriter, _ *http.Request) {
	log.Info("Request receive to get Todos")

	todoModels, err := getDbTodos()
	if err != nil {
		rnd.JSON(response, http.StatusBadRequest, renderer.M{
			"message": "Could not get the todos",
			"error":   err.Error(),
		})
	}
	todoList := todo.MapTodoListModelToDto(todoModels)

	rnd.JSON(response, http.StatusOK, todo.TodoListResponse{
		Message: "All todos retrieved",
		Data:    todoList,
	})
}

func createTodo(response http.ResponseWriter, request *http.Request) {
	log.WithFields(log.Fields{"body": request.Body}).Info("Request receive to create Todo")

	todoDto := todo.TodoDto{}
	if err := json.NewDecoder(request.Body).Decode(&todoDto); err != nil {
		log.WithFields(log.Fields{"details": err.Error()}).Error("failed to decode json data")
		rnd.JSON(response, http.StatusBadRequest, renderer.M{
			"message": "could not decode data",
		})
		return
	}

	if todoDto.Title == "" {
		log.Error("no title added to request body")
		rnd.JSON(response, http.StatusBadRequest, renderer.M{
			"message": "please add a title",
		})
		return
	}

	err := createDbTodo(todoDto)
	if err != nil {
		rnd.JSON(response, http.StatusBadRequest, renderer.M{
			"message": "Could not create the todo",
			"error":   err.Error(),
		})
	}

	rnd.JSON(response, http.StatusOK, todo.Response{
		Message: "Todo created",
	})
}

func deleteTodo(response http.ResponseWriter, request *http.Request) {
	todoId := strings.TrimSpace(chi.URLParam(request, "id"))

	log.WithFields(log.Fields{"id": todoId}).Info("Request receive to delete Todo")

	if todoId == "" {
		log.Error("no id added to request body")
		rnd.JSON(response, http.StatusBadRequest, renderer.M{
			"message": "please add an id",
		})
		return
	}

	err := deleteDbTodo(todoId)
	if err != nil {
		rnd.JSON(response, http.StatusBadRequest, renderer.M{
			"message": "Could not delete the todo",
			"error":   err.Error(),
		})
	}

	rnd.JSON(response, http.StatusOK, todo.Response{
		Message: "Todo deleted",
	})
}

func updateTodo(response http.ResponseWriter, request *http.Request) {
	log.WithFields(log.Fields{"body": request.Body}).Info("Request receive to update Todo")

	var todoDto todo.TodoDto
	if err := json.NewDecoder(request.Body).Decode(&todoDto); err != nil {
		log.WithFields(log.Fields{"details": err.Error()}).Error("failed to decode json data")
		rnd.JSON(response, http.StatusBadRequest, renderer.M{
			"message": "could not decode data",
		})
		return
	}

	if todoDto.Title == "" {
		log.Error("no title added to response body")
		rnd.JSON(response, http.StatusBadRequest, renderer.M{
			"message": "please add a title",
		})
		return
	}

	err := updateDbTodo(todoDto)
	if err != nil {
		rnd.JSON(response, http.StatusBadRequest, renderer.M{
			"message": "Could not update the todo",
			"error":   err.Error(),
		})
	}

	rnd.JSON(response, http.StatusOK, todo.Response{
		Message: "Todo updated",
	})
}

// Fake task to exercise with goroutines
func fakeTask(response http.ResponseWriter, _ *http.Request) {
	log.Info("Request receive to do a fake task")

	var wg sync.WaitGroup
	wg.Add(2)
	go ExecuteFakeTask(&wg)
	channel := make(chan any, 1) // buffer channel due to waitgroup
	go ExecuteFakeTaskWithReturn(&wg, channel)
	wg.Wait()

	res := <-channel
	switch res := res.(type) { // verify type
	case int:
		rnd.JSON(response, http.StatusOK, "fake task returned the res : "+strconv.Itoa(res))
	default:
		err := res.(error) // type assertion
		rnd.JSON(response, http.StatusBadRequest, "fake task returned the error : "+string(err.Error()))
	}
}
