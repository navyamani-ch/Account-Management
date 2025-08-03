
markdown
Copy
Edit
# 💳 Account Management Service

A basic account and transaction management microservice built using **Go** and **PostgreSQL**.

---

## 📝 Assumptions

- A single user has one account.
- Each transaction is between two valid accounts.
- Balance values are stored in `float64` (for simplicity).
- All accounts and transactions are stored in PostgreSQL.

---

## ⚙️ Tech Stack

- **Language**: Go 1.20+
- **Database**: PostgreSQL
- **Router**: Gorilla Mux
- **Migrations**: golang-migrate

---

## ✅ Prerequisites

### 📦 Install Go

- Download: https://go.dev/dl/
- Verify installation:
  ```bash
  go version
  
### 🐘 Install PostgreSQL
- **For macOS (Homebrew)**:
   ```bash
  brew install postgresql
  brew services start postgresql

- **For Ubuntu/Debian**:
  ```bash
  sudo apt update
  sudo apt install postgresql postgresql-contrib
  sudo service postgresql start
  
### 🛠️ Setup
- Create PostgreSQL Role
   ```bash
   psql postgres
   
- Then inside the shell:
  ```sql
  CREATE ROLE root WITH LOGIN PASSWORD 'password';
  ALTER ROLE root CREATEDB;
  \q;
  CREATE account_service DATABASE;

---

### Run the Application
- Account-Management
  ```bash
   go run main.go
  
-The service will be running at: http://localhost:8080


## Use curl or Postman:

### ➕ Create Account
- POST /accounts

  ```bash
  curl --location 'http://localhost:8080/accounts' \
  --header 'Content-Type: application/json' \
  --data '{
  "account_id":123,
  "initial_balance":"1000.23"
  }'

### 📊 Get Account
- GET /accounts/{account_id}
  
  ```bash
  curl http://localhost:8080/accounts/123

### 📊 POST Transaction
- POST /transactions

  ```bash
  curl --location 'localhost:8080/transactions' \
  --header 'Content-Type: application/json' \
  --data '{
  "source_account_id": 123,
  "destination_account_id": 456,
  "amount": "100.12345"
  }'