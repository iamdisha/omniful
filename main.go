package main

import (
  "database/sql"
  "encoding/json"
  "fmt"
  "log"
  "net/http"
  "time"

  "github.com/gorilla/mux"
  _ "github.com/go-sql-driver/mysql"
)


type Order struct {
  ID           int       `json:"id"`
  CustomerName string    `json:"customer_name"`
  ProductID    string    `json:"product_id"`
  Quantity     int       `json:"quantity"`
  OrderDate    time.Time `json:"order_date"`
}


var db *sql.DB


func InitDB() {
  var err error

  dsn := "user:password@tcp(127.0.0.1:3306)/orders_db"
  db, err = sql.Open("mysql", dsn)
  if err != nil {
    log.Fatalf("Failed to connect to the database: %v", err)
  }


  if err = db.Ping(); err != nil {
    log.Fatalf("Failed to ping the database: %v", err)
  }




  log.Println("Database initialized and table ensured.")
}

// CreateOrderHandler handles the creation of a single order.
func CreateOrderHandler(w http.ResponseWriter, r *http.Request) {
  
  var order Order
  if err := json.NewDecoder(r.Body).Decode(&order); err != nil {
    http.Error(w, "Invalid request payload", http.StatusBadRequest)
    return
  }

  // Validate input.
  if order.CustomerName == "" || order.ProductID == "" || order.Quantity <= 0 {
    http.Error(w, "Missing or invalid fields", http.StatusBadRequest)
    return
  }

  // Inserting order
  order.OrderDate = time.Now()
  query := "INSERT INTO orders (customer_name, product_id, quantity, order_date) VALUES (?, ?, ?, ?)"
  result, err := db.Exec(query, order.CustomerName, order.ProductID, order.Quantity, order.OrderDate)
  if err != nil {
    http.Error(w, "Failed to create order", http.StatusInternalServerError)
    log.Printf("Database error: %v", err)
    return
  }

  // Retrieve the ID
  id, _ := result.LastInsertId()
  order.ID = int(id)

  
  w.WriteHeader(http.StatusCreated)
  json.NewEncoder(w).Encode(order)
}

func main() {
  InitDB()
  defer db.Close()
  r := mux.NewRouter()


  r.HandleFunc("/create-order, CreateOrderHandler).Methods("POST")


  log.Println("Server is running on port 8080...")
  if err := http.ListenAndServe(":8080", r); err != nil {
    log.Fatalf("Failed to start the server: %v", err)
  }
}
