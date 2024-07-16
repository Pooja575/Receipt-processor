package main

import (
    "encoding/json"
    "github.com/gorilla/mux"
    "log"
    "math"
    "net/http"
    "strconv"
    "strings"
    "sync"
    "time"
    "github.com/google/uuid"
)

type Item struct {
    ShortDescription string `json:"shortDescription"`
    Price            string `json:"price"`
}

type Receipt struct {
    ID           string `json:"id"`
    Retailer     string `json:"retailer"`
    PurchaseDate string `json:"purchaseDate"`
    PurchaseTime string `json:"purchaseTime"`
    Items        []Item `json:"items"`
    Total        string `json:"total"`
    Points       int    `json:"points"`
}

var (
    receipts = make(map[string]Receipt)
    mu       sync.Mutex
)

func calculatePoints(receipt Receipt) int {
    points := 0

    // Rule 1
    for _, char := range receipt.Retailer {
        if (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') || (char >= '0' && char <= '9') {
            points++
        }
    }

    
    total, _ := strconv.ParseFloat(receipt.Total, 64)

    // Rule 2
    if total == float64(int(total)) {
        points += 50
    }

    // Rule 3
    if math.Mod(total, 0.25) == 0 {
        points += 25
    }

    // Rule 4
    points += (len(receipt.Items) / 2) * 5

    // Rule 5
    for _, item := range receipt.Items {
        trimmedLen := len(strings.TrimSpace(item.ShortDescription))
        if trimmedLen%3 == 0 {
            price, _ := strconv.ParseFloat(item.Price, 64)
            points += int(math.Ceil(price * 0.2))
        }
    }

    // Rule 6
    date, _ := time.Parse("2006-01-02", receipt.PurchaseDate)
    if date.Day()%2 != 0 {
        points += 6
    }

    // Rule 7
    purchaseTime, _ := time.Parse("15:04", receipt.PurchaseTime)
    if purchaseTime.Hour() == 14 {
        points += 10
    }

    return points
}

func processReceipts(w http.ResponseWriter, r *http.Request) {
    var receipt Receipt
    if err := json.NewDecoder(r.Body).Decode(&receipt); err != nil {
        http.Error(w, err.Error(), http.StatusBadRequest)
        return
    }

    receipt.ID = uuid.New().String()
    receipt.Points = calculatePoints(receipt)

    mu.Lock()
    receipts[receipt.ID] = receipt
    mu.Unlock()

    response := map[string]string{"id": receipt.ID}
    json.NewEncoder(w).Encode(response)
}

func getPoints(w http.ResponseWriter, r *http.Request) {
    vars := mux.Vars(r)
    id := vars["id"]

    mu.Lock()
    receipt, found := receipts[id]
    mu.Unlock()

    if !found {
        http.Error(w, "Receipt not found", http.StatusNotFound)
        return
    }

    response := map[string]int{"points": receipt.Points}
    json.NewEncoder(w).Encode(response)
}

func main() {
    r := mux.NewRouter()
    r.HandleFunc("/receipts/process", processReceipts).Methods("POST")
    r.HandleFunc("/receipts/{id}/points", getPoints).Methods("GET")

    log.Println("Starting server on :8080...")
    http.ListenAndServe(":8080", r)
}
