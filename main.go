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

    // Rule 1: One point for every alphanumeric character in the retailer name.
    for _, char := range receipt.Retailer {
        if (char >= 'a' && char <= 'z') || (char >= 'A' && char <= 'Z') || (char >= '0' && char <= '9') {
            points++
        }
    }

    // Convert Total to float for further calculations
    total, _ := strconv.ParseFloat(receipt.Total, 64)

    // Rule 2: 50 points if the total is a round dollar amount with no cents.
    if total == float64(int(total)) {
        points += 50
    }

    // Rule 3: 25 points if the total is a multiple of 0.25.
    if math.Mod(total, 0.25) == 0 {
        points += 25
    }

    // Rule 4: 5 points for every two items on the receipt.
    points += (len(receipt.Items) / 2) * 5

    // Rule 5: Points based on the trimmed length of the item description.
    for _, item := range receipt.Items {
        trimmedLen := len(strings.TrimSpace(item.ShortDescription))
        if trimmedLen%3 == 0 {
            price, _ := strconv.ParseFloat(item.Price, 64)
            points += int(math.Ceil(price * 0.2))
        }
    }

    // Rule 6: 6 points if the day in the purchase date is odd.
    date, _ := time.Parse("2006-01-02", receipt.PurchaseDate)
    if date.Day()%2 != 0 {
        points += 6
    }

    // Rule 7: 10 points if the time of purchase is after 2:00pm and before 4:00pm.
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
