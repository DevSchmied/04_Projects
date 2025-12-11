# BookManager-Go

A **containerized book management web application** built with **Go**.  
Start it instantly with **Docker** — no manual setup required.

---

## Quick Start with Docker

This project is fully containerized.  
**No manual environment setup is required.**  
Docker automatically creates all required configuration files inside the container.

### 1. Clone the repository

> Note: This project is part of a multi-project repository (`04_Projects`).
> After cloning, navigate into the project directory before running Docker commands (the required `cd` command is shown below).

```bash
git clone https://github.com/DevSchmied/04_Projects.git
cd 04_Projects/bookmanager-go
```

### 2. Build and start the container

```bash
docker compose up --build
```

Then open your browser:

```
http://localhost:8080/books
```

### How it works (important for reviewers)

- Docker uses a **multi-stage build**
- During the build, it automatically creates:

```
internal/config/app.env
```

based on:

```
internal/config/app.env.example
```

- No manual steps are required for the user
- The application runs in **release mode** and uses **SQLite** inside the container
- This ensures that the project runs **on any machine**, independent of local configuration

---

## Optional: Enable Redis Caching

Redis is **optional**.
The application runs fully without Redis.

To test the caching layer, start a Redis instance locally.

If Redis is unavailable or slow, the application **automatically falls back** to SQLite.

---

## Features
- Add, edit, delete, and list books 
- User authentication system: Register, Login, Logout (JWT + httpOnly cookies)
- Book search (via **Strategy Pattern**: ID or Title)
- Simple, responsive web interface (**HTML + Bootstrap**)  
- Persistent storage via **SQLite** (with optional MySQL support)
- Modular architecture with **SOLID** and **Clean Architecture** principles  
- Unit tests following **AAA (Arrange–Act–Assert)** and **FIRST** principles
- **Redis caching layer** for performance demonstration (see below)

---

## Tech Stack
- **Go** (Gin framework)
- **HTML / Bootstrap 5**
- **GORM** ORM
- **SQLite** (or MySQL)
- **JWT Authentication**
- **Redis** (cache layer)
- **Docker / Docker Compose**
- **Go Testing** with mocks and table-driven tests

---

## Architecture Overview
The project follows a layered, modular structure inspired by Clean Architecture.

bookmanager-go/
├── cmd/                     # Application entry point
│   └── main.go
├── internal/
│   ├── auth/                # JWT logic, authentication middleware
│   ├── cache/               # Redis caching layer (BookList example)
│   ├── config/              # Environment loader with retry mechanism
│   ├── controller/          # HTTP handlers (HTML views)
│   ├── model/               # GORM entities
│   ├── server/              # Gin server setup, routing, DI
│   └── service/             # DB initialization (SQLite connector)
├── internal/view/
│   ├── templates/           # HTML templates
│   └── static/              # CSS, images, logos


---

## Redis Caching Layer (Demonstration Feature)

Caching is implemented for the **BookList endpoint**:

```bash
GET /books/list
```

BookManager-Go uses the **Cache-Aside pattern**:

This is a **demonstration implementation** designed to show how Redis can be integrated into a Go project.
Currently, caching is implemented **only for the BookList**, but the structure allows extending caching to other parts of the project in the future.

BookManager-Go uses the **Cache-Aside pattern**.

### Read Flow

1. Check Redis for cached list
2. **HIT**: return cached response
3. **MISS**: load from DB → write to Redis → return result

### Write Flow

On creation, update, or deletion of a book:
- Redis cache is invalidated
- Next read regenerates it

### TTL (Time-To-Live)

Book list expires after **60 seconds**.

### Fault Tolerance
If Redis is:
- slow
- unreachable
- misconfigured
the operation is aborted after **300 ms**, logs are written, and SQLite is used instead.
**The system remains fully functional**.

---

### Design Patterns Implemented
- **Strategy Pattern**: For dynamic book search (by ID or title) 
- **Interface-based architecture**
- **Dependency Injection** via structured server initialization

---

### Visual Assets
The visual assets (logos, icons, background watermark) were generated using **AI tools (ChatGPT / DALL·E by OpenAI)**.  
They were created and used solely for **learning and demonstration purposes** within this **personal pet project**.

---

### License
This project is open-source and available for educational and portfolio purposes.
