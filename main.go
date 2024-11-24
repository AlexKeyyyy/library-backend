package main

import (
	"library-backend/db"
	"library-backend/handlers"
	"log"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/rs/cors"
)

func main() {

	// Подключение к базе данных
	db.Connect()

	// Инициализация роутера
	r := mux.NewRouter()

	// Определение маршрутов
	r.HandleFunc("/login", handlers.LoginLibrarian).Methods("POST")
	r.HandleFunc("/register", handlers.RegisterLibrarian).Methods("POST")
	r.HandleFunc("/clients", handlers.GetClients).Methods("GET")
	r.HandleFunc("/clients", handlers.AddClient).Methods("POST")
	r.HandleFunc("/clients/{id}", handlers.UpdateClient).Methods("PUT")
	r.HandleFunc("/clients/{id}", handlers.DeleteClient).Methods("DELETE")
	r.HandleFunc("/clients/all", handlers.GetAllClients).Methods("GET")
	r.HandleFunc("/clients/{id}", handlers.GetClientByID).Methods("GET")

	// Маршруты для книг
	r.HandleFunc("/books", handlers.GetBooks).Methods("GET")
	r.HandleFunc("/books", handlers.AddBook).Methods("POST")
	r.HandleFunc("/books/{id}", handlers.UpdateBook).Methods("PUT")
	r.HandleFunc("/books/{id}", handlers.DeleteBook).Methods("DELETE")
	r.HandleFunc("/books/all", handlers.GetAllBooks).Methods("GET")
	r.HandleFunc("/books/{id}", handlers.GetBookByID).Methods("GET")

	// Маршруты для типов книг
	r.HandleFunc("/book_types", handlers.GetBookTypes).Methods("GET")
	r.HandleFunc("/book_types", handlers.AddBookType).Methods("POST")
	r.HandleFunc("/book_types/{id}", handlers.UpdateBookType).Methods("PUT")
	r.HandleFunc("/book_types/{id}", handlers.DeleteBookType).Methods("DELETE")
	r.HandleFunc("/book_types/{id}", handlers.GetBookType).Methods("GET")

	// Маршруты для журнала
	r.HandleFunc("/journal/issue", handlers.IssueBook).Methods("POST")   // Выдача книги
	r.HandleFunc("/journal/return", handlers.ReturnBook).Methods("POST") // Прием книги
	r.HandleFunc("/journal", handlers.GetJournalEntries).Methods("GET")  // Получение записей журнала
	r.HandleFunc("/journal/fine", handlers.GetFine).Methods("POST")

	r.HandleFunc("/reports/top-books", handlers.GetTopBooks).Methods("GET")
	r.HandleFunc("/reports/top-clients-fines", handlers.GetTopClientsWithFines).Methods("GET")

	// Добавление CORS
	c := cors.New(cors.Options{
		AllowedOrigins: []string{"http://localhost:3000"}, // Разрешаем запросы с этого порта
		AllowedMethods: []string{"GET", "POST", "PUT", "DELETE"},
		AllowedHeaders: []string{"Content-Type"},
	})

	handler := c.Handler(r)
	// Запуск сервера с CORS
	log.Println("Server is running on port 8080")
	log.Fatal(http.ListenAndServe(":8080", handler))
}
