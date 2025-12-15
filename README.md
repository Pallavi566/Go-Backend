# User Management Backend (Go)

This project is a backend service developed using Go to manage user data where the **age is calculated dynamically from the date of birth** instead of being stored.

The implementation focuses on clean project structure, clear separation of responsibilities, and reliable data handling.

---

## Project Overview

The application stores user information, including date of birth, in a MySQL database.  
Rather than persisting age, it is **computed at runtime** whenever user data is retrieved, ensuring accuracy over time.

---

## Technologies Used

- **Go (Fiber)** – backend framework
- **MySQL** – relational database
- **SQLC** – type-safe database access
- **go-playground/validator** – request validation
- **Uber Zap** – structured logging

---

## Project Structure

/cmd/server/main.go
/config/
/db/migrations/
/db/sqlc/<generated>
/internal/
├── handler/
├── repository/
├── service/
├── routes/
├── middleware/
├── models/
└── logger/

The structure is designed to keep responsibilities clearly separated and the codebase easy to maintain.

---

## Database Design

The application uses a single `users` table with the following fields:

- `id` – primary key
- `name` – user name
- `dob` – date of birth

The user’s age is **not stored** in the database.  
It is calculated dynamically using Go’s `time` package.

---

## Setup Summary

For this project:
- Go (version 1.21 or higher) was installed
- MySQL was installed and configured
- Required environment variables were set
- Database schema was applied using migration files

---

Key Highlights:

Dynamic age calculation without storing redundant data
Type-safe database interactions using SQLC
Input validation at the API boundary
Structured logging for better observability
Clean separation between HTTP, service, and data layers

What This Project Demonstrates:

Backend development using Go
MySQL database integration
Clean and maintainable project architecture
Practical implementation of backend best practices

## API functionality was tested locally using Postman.##

The server runs locally and handles all user-related operations.

## Running the Application

After completing the setup, the application was started using:

```bash
go run cmd/server/main.go
