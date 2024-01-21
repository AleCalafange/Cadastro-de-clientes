package main

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/go-chi/chi/v5"

	_ "github.com/mattn/go-sqlite3"
)

type Customer struct {
	ID          int    `json:"id"`
	FirstName   string `json:"firstName"`
	LastName    string `json:"lastName"`
	Address     string `json:"address"`
	Phone       string `json:"phone"`
	Email       string `json:"email"`
	DateOfBirth string `json:"dateOfBirth"`
}

type Resource struct {
	db *sql.DB
}

func (res *Resource) CreateCustomer(w http.ResponseWriter, r *http.Request) {
	var customer Customer
	if err := json.NewDecoder(r.Body).Decode(&customer); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	insertQuery := `INSERT INTO customers(firstName, lastName, address, phone, email, dateOfBirth) 
	                VALUES (?, ?, ?, ?, ?, ?);`

	_, err := res.db.Exec(insertQuery, customer.FirstName, customer.LastName, customer.Address, customer.Phone, customer.Email, customer.DateOfBirth)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	fmt.Println("Customer created - ID:", customer.ID, "Name:", customer.FirstName, customer.LastName)
	w.Write([]byte("Customer created successfully"))
}

func (res *Resource) ListCustomers(w http.ResponseWriter, r *http.Request) {
	var customers []Customer

	rows, err := res.db.Query("SELECT id, firstName, lastName, address, phone, email, dateOfBirth FROM customers")
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	defer rows.Close()
	for rows.Next() {
		var customer Customer
		err := rows.Scan(&customer.ID, &customer.FirstName, &customer.LastName, &customer.Address, &customer.Phone, &customer.Email, &customer.DateOfBirth)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		customers = append(customers, customer)
	}

	result, err := json.Marshal(customers)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Write([]byte(result))
	fmt.Println(result)
}

func main() {
	db, err := sql.Open("sqlite3", "customers.db")
	if err != nil {
		fmt.Println(err)
		return
	}
	createCustomersTable(db)

	resource := Resource{db}

	router := chi.NewRouter()
	router.Post("/create-customer", resource.CreateCustomer)
	router.Get("/list-customers", resource.ListCustomers)

	http.ListenAndServe(":3333", router)
}

func createCustomersTable(db *sql.DB) {
	createTableQuery := `CREATE TABLE IF NOT EXISTS customers (
		id INTEGER NOT NULL PRIMARY KEY AUTOINCREMENT,
		firstName TEXT,
		lastName TEXT,
		address TEXT,
		phone TEXT,
		email TEXT,
		dateOfBirth TEXT
	);`

	_, err := db.Exec(createTableQuery)
	if err != nil {
		fmt.Println(err)
		return
	}

	fmt.Println("Customers table created successfully")
}
