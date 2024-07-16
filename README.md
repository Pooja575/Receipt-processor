# Receipt Processor API

This project implements a receipt processing API using Go. It calculates points based on various rules defined for a receipt.

## Installation

### Prerequisites

- Docker installed on your machine

### Steps to Run

1. Clone this repository:

   ```bash
   git clone <repository-url>
   cd receipt-processor

2. Build the Docker image:
docker build -t receipt-processor .

3.Run the Docker container:
docker run -p 8081:8080 receipt-processor

4. Access the server at http://localhost:8081.

## Usage

### API Endpoints

#### Process Receipts

Endpoint: `http://localhost:8081/receipts/process`
- Method: `POST`
- Payload: JSON representing the receipt
  ```json
  {
    "retailer": "Example Retailer",
    "purchaseDate": "2024-07-16",
    "purchaseTime": "15:30",
    "total": "100.00",
    "items": [
      {"shortDescription": "Item 1", "price": "50.00"},
      {"shortDescription": "Item 2", "price": "50.00"}
    ]
  }

Response: JSON containing an id for the receipt
{"id":"083e6fa5-f770-478c-9235-667d1af54b53"}


# Get Points

Endpoint: http://localhost:8081/receipts/{id}/points

Method: GET
Path Parameter: id - The ID of the receipt
Response: JSON containing the points awarded for the receipt

{ "points": 100 }

# Testing

# Submit a receipt for processing
curl -X POST -H "Content-Type: application/json" -d '{"retailer": "Example Retailer", "purchaseDate": "2024-07-16", "purchaseTime": "15:30", "total": "100.00", "items": [{"shortDescription": "Item 1", "price": "50.00"}, {"shortDescription": "Item 2", "price": "50.00"}]}' http://localhost:8081/receipts/process

# Retrieve points for a receipt ID (replace {id} with actual ID)
curl http://localhost:8081/receipts/{id}/points

# Additional Information
For more details, refer to the API specification in api.yml.
