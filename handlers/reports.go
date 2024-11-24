package handlers

import (
	"encoding/json"
	"library-backend/db"
	"net/http"
)

// Определение структуры для топ-книг
type TopBook struct {
	BookName    string `json:"name"`
	BorrowCount int    `json:"borrow_count"`
}

func GetTopBooks(w http.ResponseWriter, r *http.Request) {
	query := `
        SELECT 
            b.name AS name,
            COUNT(j.book_id) AS borrow_count
        FROM journal j
        JOIN books b ON j.book_id = b.id
        GROUP BY b.name
        ORDER BY borrow_count DESC
        LIMIT 3;
    `

	rows, err := db.DB.Query(query)
	if err != nil {
		http.Error(w, "Error fetching top books", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var books []TopBook
	for rows.Next() {
		var book TopBook
		if err := rows.Scan(&book.BookName, &book.BorrowCount); err != nil {
			http.Error(w, "Error scanning top books", http.StatusInternalServerError)
			return
		}
		books = append(books, book)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(books)
}

type ClientWithFine struct {
	ClientName string `json:"client_name"`
	TotalFine  int    `json:"total_fine"`
}

func GetTopClientsWithFines(w http.ResponseWriter, r *http.Request) {
	query := `
        SELECT 
            c.last_name || ' ' || c.first_name AS client_name,
            SUM(j.fine_today) AS total_fine
        FROM journal j
        JOIN clients c ON j.client_id = c.id
        WHERE j.fine_today > 0
        GROUP BY c.id, c.last_name, c.first_name
        ORDER BY total_fine DESC;
    `

	rows, err := db.DB.Query(query)
	if err != nil {
		http.Error(w, "Error fetching clients with fines", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var clients []ClientWithFine
	for rows.Next() {
		var client ClientWithFine
		if err := rows.Scan(&client.ClientName, &client.TotalFine); err != nil {
			http.Error(w, "Error scanning clients with fines", http.StatusInternalServerError)
			return
		}
		clients = append(clients, client)
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(clients)
}
