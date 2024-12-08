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

func GetBooksOnHand(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ClientID int `json:"client_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	var count int
	query := `
		   SELECT COUNT(*)
		   FROM journal
		   WHERE client_id = $1 AND date_ret IS NULL`
	err := db.DB.QueryRow(query, req.ClientID).Scan(&count)
	if err != nil {
		http.Error(w, "Error fetching books on hand", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(struct {
		ClientID    int `json:"client_id"`
		BooksOnHand int `json:"books_on_hand"`
	}{
		ClientID:    req.ClientID,
		BooksOnHand: count,
	})
}
func GetClientFine(w http.ResponseWriter, r *http.Request) {
	var req struct {
		ClientID int `json:"client_id"`
	}

	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	var totalFine int
	query := `
		   SELECT SUM(fine_today)
		   FROM journal
		   WHERE client_id = $1`
	err := db.DB.QueryRow(query, req.ClientID).Scan(&totalFine)
	if err != nil {
		http.Error(w, "Error fetching client fine", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(struct {
		ClientID  int `json:"client_id"`
		TotalFine int `json:"total_fine"`
	}{
		ClientID:  req.ClientID,
		TotalFine: totalFine,
	})
}
