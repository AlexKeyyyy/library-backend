package handlers

import (
	"database/sql"
	"encoding/json"
	"library-backend/db"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
)

type Client struct {
	ID             int    `json:"id"`
	FirstName      string `json:"first_name"`
	LastName       string `json:"last_name"`
	FatherName     string `json:"father_name"`
	PassportSeria  string `json:"passport_seria"`
	PassportNumber string `json:"passport_number"`
}

func GetClients(w http.ResponseWriter, r *http.Request) {
	rows, err := db.DB.Query("SELECT id, first_name, last_name, father_name, passport_seria, passport_number  FROM clients")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var clients []Client
	for rows.Next() {
		var client Client
		err := rows.Scan(&client.ID, &client.FirstName, &client.LastName, &client.FatherName, &client.PassportSeria, &client.PassportNumber)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		clients = append(clients, client)
	}

	json.NewEncoder(w).Encode(clients)
}

func GetAllClients(w http.ResponseWriter, r *http.Request) {
	rows, err := db.DB.Query("SELECT id, first_name, last_name, father_name FROM clients")
	if err != nil {
		http.Error(w, "Error fetching clients", http.StatusInternalServerError)
		return
	}
	defer rows.Close()

	var clients []Client
	for rows.Next() {
		var client Client
		err := rows.Scan(&client.ID, &client.FirstName, &client.LastName, &client.FatherName)
		if err != nil {
			http.Error(w, "Error scanning clients", http.StatusInternalServerError)
			return
		}
		clients = append(clients, client)
	}

	json.NewEncoder(w).Encode(clients)
}

func GetClientByID(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id := vars["id"]

	var client Client
	query := "SELECT id, first_name, last_name, father_name, passport_seria, passport_number FROM clients WHERE id = $1"
	err := db.DB.QueryRow(query, id).Scan(&client.ID, &client.FirstName, &client.LastName, &client.FatherName, &client.PassportSeria, &client.PassportNumber)
	if err == sql.ErrNoRows {
		http.Error(w, "Client not found", http.StatusNotFound)
		return
	} else if err != nil {
		http.Error(w, "Error fetching client", http.StatusInternalServerError)
		return
	}

	json.NewEncoder(w).Encode(client)
}

func AddClient(w http.ResponseWriter, r *http.Request) {
	var client Client
	// Декодируем JSON из тела запроса
	err := json.NewDecoder(r.Body).Decode(&client)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Вставляем данные в базу
	query := "INSERT INTO clients (first_name, last_name, father_name, passport_seria, passport_number) VALUES ($1, $2, $3, $4, $5)"
	_, err = db.DB.Exec(query, client.FirstName, client.LastName, client.FatherName, client.PassportSeria, client.PassportNumber)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Write([]byte("Client added successfully"))
}

func UpdateClient(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	clientID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid client ID", http.StatusBadRequest)
		return
	}

	var client Client
	err = json.NewDecoder(r.Body).Decode(&client)
	if err != nil {
		http.Error(w, "Invalid input", http.StatusBadRequest)
		return
	}

	// Обновляем данные клиента
	query := "UPDATE clients SET first_name=$1, last_name=$2, father_name=$3, passport_seria=$4, passport_number=$5 WHERE id=$6"
	res, err := db.DB.Exec(query, client.FirstName, client.LastName, client.FatherName, client.PassportSeria, client.PassportNumber, clientID)
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
		http.Error(w, "Client not found", http.StatusNotFound)
		return
	}

	w.Write([]byte("Client updated successfully"))
}

func DeleteClient(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	clientID, err := strconv.Atoi(vars["id"])
	if err != nil {
		http.Error(w, "Invalid client ID", http.StatusBadRequest)
		return
	}

	// Удаляем клиента
	query := "DELETE FROM clients WHERE id=$1"
	res, err := db.DB.Exec(query, clientID)
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
		http.Error(w, "Client not found", http.StatusNotFound)
		return
	}

	w.Write([]byte("Client deleted successfully"))
}
