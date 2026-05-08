# Kantin Mardira API

Backend REST API untuk aplikasi POS (Point of Sale) kantin kampus menggunakan Golang dengan Clean Architecture.

## Tech Stack

### Backend

- Golang
- Gin Gonic
- PostgreSQL
- GORM
- JWT Authentication
- Clean Architecture

### Frontend

- Vue.js
- Axios
- Pinia/Vuex (optional)
- Tailwind CSS / Bootstrap

---

# Features

## Authentication

- Login
- Logout
- JWT Authentication
- Role Middleware
- Protected Routes

### Roles

- admin
- cashier
- manager

---

# User Management

- Create User
- Get All Users
- Get User By ID
- Update User
- Delete User

Authorization:

- Admin only for create/update/delete

---

# Category Management

- Create Category
- Get All Categories
- Get Category By ID
- Update Category
- Delete Category

---

# Menu Management

- Create Menu
- Get All Menus
- Get Menu By ID
- Update Menu
- Delete Menu

Features:

- Menu relation with Category
- Menu stock
- Menu availability
- Nested category response

---

# Transaction System

## Features

- Create Transaction
- Multiple Menu Items in One Transaction
- Transaction Detail
- Transaction History
- Role-based Transaction Access
- Automatic Total Calculation
- Automatic Stock Reduction
- Payment Logic
- Database Transaction Rollback

---

# Payment Methods

- cash
- qris

---

# Transaction Features

## Automatic Calculation

- subtotal
- total amount
- change amount

## Stock Validation

- reject insufficient stock
- reject unavailable menu

## Database Safety

- rollback transaction if one item fails
- atomic database transaction

---

# Reports & Analytics

## Reports

- Daily Report
- Weekly Report
- Monthly Report
- Revenue Summary
- Top Selling Menu

## Analytics

- Total Revenue
- Total Transactions
- Total Items Sold
- Average Transaction
- Top Selling Products

---

# PDF Export

## Features

- Daily Report PDF
- Weekly Report PDF
- Monthly Report PDF
- Transaction Invoice PDF

## PDF Includes

- Kop Surat
- Transaction Table
- Revenue Summary
- Footer
- Page Number

---

# Clean Architecture

Request Flow:

Frontend
→ Routes
→ Handler
→ Service
→ Repository
→ PostgreSQL

Response Flow:

PostgreSQL
→ Repository
→ Service
→ Handler
→ JSON Response

---

# Folder Structure

```bash
cmd/
    main.go

internal/
    config/
    routes/
    handler/
    service/
    repository/
    entity/
    dto/
    middleware/
    utils/
    adapter/
```

---

# Environment Variables

```env
APP_PORT=8080

DB_HOST=
DB_PORT=5432
DB_USER=
DB_PASSWORD=
DB_NAME=
DB_SSLMODE=require

JWT_SECRET=
```

---

# Run Project

## Install Dependencies

```bash
go mod tidy
```

## Run Project

```bash
go run cmd/main.go
```

---

# API Base URL

```bash
http://localhost:8080/api/v1
```

---

# Authentication Header

```http
Authorization: Bearer <token>
```

---

# Important API Endpoints

## Authentication

### Login

```http
POST /auth/login
```

---

# Users

### Get All Users

```http
GET /users
```

### Create User

```http
POST /users
```

---

# Categories

### Get Categories

```http
GET /categories
```

---

# Menus

### Get Menus

```http
GET /menus
```

---

# Transactions

### Create Transaction

```http
POST /transactions
```

### Get Transactions

```http
GET /transactions
```

### Get Transaction Detail

```http
GET /transactions/:id
```

---

# Reports

### Daily Report

```http
GET /reports/daily
```

### Weekly Report

```http
GET /reports/weekly
```

### Monthly Report

```http
GET /reports/monthly
```

### Summary Report

```http
GET /reports/summary
```

### Top Selling

```http
GET /reports/top-selling
```

---

# PDF Export

### Daily PDF

```http
GET /reports/daily/pdf
```

### Monthly PDF

```http
GET /reports/monthly/pdf
```

### Invoice PDF

```http
GET /transactions/:id/invoice
```

---

# Notes For Frontend (Vue.js)

## Login Flow

1. User login
2. Save JWT token to:

- localStorage
- cookie
- pinia store

3. Send token using:

```http
Authorization: Bearer <token>
```

---

# Suggested Frontend Pages

## Authentication

- Login Page

## Dashboard

- Revenue Summary
- Daily Revenue Chart
- Top Selling Menu

## Master Data

- User Management
- Category Management
- Menu Management

## Transaction

- POS Page
- Cart System
- Checkout Page
- Transaction History
- Transaction Detail

## Reports

- Daily Report
- Monthly Report
- Download PDF

---

# Suggested Frontend Flow

## POS Transaction Flow

1. Fetch menu list
2. Add menu to cart
3. Update quantity
4. Checkout
5. Send transaction request
6. Show invoice
7. Download PDF invoice

---

# Suggested Frontend Features

## Dashboard

- Revenue Card
- Transaction Card
- Top Selling Chart
- Revenue Chart

## Transaction UI

- Cart Sidebar
- Category Tabs
- Search Menu
- QRIS Payment Modal

---

# Future Improvements

- Upload Menu Image
- Pagination
- Search & Filter
- Real-time Dashboard
- WebSocket Notification
- Refund System
- Audit Logs
- Excel Export
- QR Order Per Table
- AI Analytics

---

# Author

Developed by:
Teguh Bagas Mardiansyah
