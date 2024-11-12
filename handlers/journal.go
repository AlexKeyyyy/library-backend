package handlers

import (
	"database/sql"
	"encoding/json"
	"library-backend/db"
	"log"
	"net/http"
	"time"
)

type JournalEntry struct {
	ID       int    `json:"id"`
	BookID   int    `json:"book_id"`
	ClientID int    `json:"client_id"`
	DateBeg  string `json:"date_beg"`
	DateEnd  string `json:"date_end"`
	DateRet  string `json:"date_ret"`
	Fine     int    `json:"fine"`
}

func GetJournalEntries(w http.ResponseWriter, r *http.Request) {
	rows, err := db.DB.Query("SELECT id, book_id, client_id, date_beg, date_end, date_ret, fine FROM journal")
	if err != nil {
		http.Error(w, "Error fetching journal entries", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var entries []JournalEntry
	for rows.Next() {
		var entry JournalEntry
		var dateRet sql.NullString
		err := rows.Scan(&entry.ID, &entry.BookID, &entry.ClientID, &entry.DateBeg, &entry.DateEnd, &dateRet, &entry.Fine)
		if err != nil {
			http.Error(w, "Error scanning journal entry", http.StatusInternalServerError)
			return
		}

		// Преобразуем dateRet в строку
		if dateRet.Valid {
			entry.DateRet = dateRet.String
		} else {
			entry.DateRet = "Не возвращена" // Устанавливаем строку, если значение не валидно
		}
		entries = append(entries, entry)
	}

	json.NewEncoder(w).Encode(entries)
}

type IssueRequest struct {
	BookID   int    `json:"book_id"`
	ClientID int    `json:"client_id"`
	DateEnd  string `json:"date_end"`
}

func IssueBook(w http.ResponseWriter, r *http.Request) {
	var request IssueRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		log.Println("Ошибка декодирования JSON:", err)
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	dateEnd, err := time.Parse("2006-01-02", request.DateEnd)
	if err != nil {
		log.Println("Ошибка преобразования даты:", err)
		http.Error(w, "Invalid date format. Use YYYY-MM-DD", http.StatusBadRequest)
		return
	}

	query := "INSERT INTO journal (book_id, client_id, date_beg, date_end) VALUES ($1, $2, $3, $4)"
	_, err = db.DB.Exec(query, request.BookID, request.ClientID, time.Now(), dateEnd)
	if err != nil {
		log.Println("Ошибка добавления записи в журнал:", err)
		http.Error(w, "Error issuing book", http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Book issued successfully"))
}

type ReturnRequest struct {
	JournalID int `json:"journal_id"` // ID записи в журнале
}

func ReturnBook(w http.ResponseWriter, r *http.Request) {
	var request ReturnRequest
	err := json.NewDecoder(r.Body).Decode(&request)
	if err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Получаем информацию о книге и дате возврата
	var dateEnd, dateRet sql.NullTime
	var finePerDay int
	query := `
        SELECT j.date_end, bt.fine
        FROM journal j
        JOIN books b ON j.book_id = b.id
		JOIN book_types bt ON b.type_id = bt.id
        WHERE j.id = $1`
	err = db.DB.QueryRow(query, request.JournalID).Scan(&dateEnd, &finePerDay)
	if err != nil {
		http.Error(w, "Journal entry not found", http.StatusNotFound)
		return
	}

	dateRet.Time = time.Now()

	// Рассчитываем штраф
	var totalFine int
	if dateEnd.Valid && dateRet.Time.After(dateEnd.Time) {
		daysLate := int(dateRet.Time.Sub(dateEnd.Time).Hours() / 24)
		totalFine = finePerDay * daysLate
	}

	// Обновляем запись о возврате и фиксируем штраф
	_, err = db.DB.Exec("UPDATE journal SET date_ret = $1, fine = $2 WHERE id = $3", dateRet.Time, totalFine, request.JournalID)
	if err != nil {
		http.Error(w, "Error updating return date", http.StatusInternalServerError)
		return
	}

	// Возвращаем итоговый штраф (если требуется)
	response := struct {
		Fine int `json:"fine"`
	}{
		Fine: totalFine,
	}
	json.NewEncoder(w).Encode(response)
}

// В контроллере для получения штрафа
func GetFine(w http.ResponseWriter, r *http.Request) {
	var request struct {
		JournalID int `json:"journal_id"`
	}
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, "Invalid request", http.StatusBadRequest)
		return
	}

	// Логика для получения штрафа за просрочку
	var fine int
	query := `
        SELECT fine
        FROM journal j
        JOIN books b ON j.book_id = b.id
        JOIN book_types bt ON b.type_id = bt.id
        WHERE j.id = $1`
	err := db.DB.QueryRow(query, request.JournalID).Scan(&fine)
	if err != nil {
		http.Error(w, "Error retrieving fine", http.StatusInternalServerError)
		return
	}

	// Отправляем штраф обратно на фронтенд
	json.NewEncoder(w).Encode(struct {
		Fine int `json:"fine"`
	}{Fine: fine})
}
