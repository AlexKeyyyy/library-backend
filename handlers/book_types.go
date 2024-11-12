package handlers

import (
	"database/sql"
	"encoding/json"
	"library-backend/db"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type BookType struct {
	ID      int     `json:"id"`
	Type    string  `json:"type"`
	Fine    float64 `json:"fine"`
	MaxDays int     `json:"day_count"`
}

func GetBookTypes(w http.ResponseWriter, r *http.Request) {
	rows, err := db.DB.Query("SELECT id, type, fine, day_count FROM book_types")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var bookTypes []BookType
	for rows.Next() {
		var bookType BookType
		err := rows.Scan(&bookType.ID, &bookType.Type, &bookType.Fine, &bookType.MaxDays)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		bookTypes = append(bookTypes, bookType)
	}

	json.NewEncoder(w).Encode(bookTypes)
}

func AddBookType(w http.ResponseWriter, r *http.Request) {
	var bookType BookType
	err := json.NewDecoder(r.Body).Decode(&bookType)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	query := "INSERT INTO book_types (type, fine, day_count) VALUES ($1, $2, $3)"
	_, err = db.DB.Exec(query, bookType.Type, bookType.Fine, bookType.MaxDays)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Book type added successfully"))
}

func UpdateBookType(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bookTypeID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid book type ID", http.StatusBadRequest)
		return
	}

	var bookType BookType
	err = json.NewDecoder(r.Body).Decode(&bookType)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	query := "UPDATE book_types SET type=$1, fine=$2, day_count=$3 WHERE id=$4"
	res, err := db.DB.Exec(query, bookType.Type, bookType.Fine, bookType.MaxDays, bookTypeID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		http.Error(w, "Book type not found", http.StatusNotFound)
		return
	}

	w.Write([]byte("Book type updated successfully"))
}

func DeleteBookType(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	bookTypeID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid book type ID", http.StatusBadRequest)
		return
	}

	query := "DELETE FROM book_types WHERE id=$1"
	res, err := db.DB.Exec(query, bookTypeID)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	rowsAffected, err := res.RowsAffected()
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	if rowsAffected == 0 {
		http.Error(w, "Book type not found", http.StatusNotFound)
		return
	}

	w.Write([]byte("Book type deleted successfully"))
}

// Обработчик для получения информации о типе книги по ID
func GetBookType(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	typeID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid type ID", http.StatusBadRequest)
		return
	}

	var bookType BookType
	query := "SELECT id, type, fine, day_count FROM book_types WHERE id = $1"
	err = db.DB.QueryRow(query, typeID).Scan(&bookType.ID, &bookType.Type, &bookType.Fine, &bookType.MaxDays)
	if err == sql.ErrNoRows {
		http.Error(w, "Book type not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Error fetching book type", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(bookType)
}
