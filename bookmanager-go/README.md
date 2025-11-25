# BookManager-Go

A **containerized book management web application** built with **Go**.  
Start it instantly with **Docker** — no manual setup required.

---

### 🐳 Quick Start with Docker

To run the application instantly on any system:

# Clone the repository
git clone https://github.com/DevSchmied/bookmanager-go.git
cd bookmanager-go

# Build and start the container
docker compose up --build
Then open your browser and navigate to:
http://localhost:8080/books

This setup runs the full application (backend, templates, static assets, SQLite DB) inside a single Docker container, ensuring consistent behavior across all environments.

### Features
- Add, edit, delete, and list books 
- User authentication system: Register, Login, Logout (JWT + httpOnly cookies)
- Search functionality 
- Simple, responsive web interface (HTML + Bootstrap)  
- Persistent storage via **SQLite** (optionally MySQL)  
- Modular architecture with **SOLID principles**  
- Unit tests following **AAA (Arrange–Act–Assert)** and **FIRST** principles  

---

### Tech Stack
- **Language:** Go (Gin framework)  
- **Frontend:** HTML / Bootstrap  
- **Database:** SQLite (or MySQL)  
- **ORM:** GORM  
- **Authentication:** JWT tokens (Gin middleware, httpOnly cookies)
- **Testing:** Go `testing` package with mocks and table-driven tests  

---

### Architecture Overview
The project follows a **modular, layered design** inspired by **Clean Architecture**, featuring separation of concerns across the following layers:

- **Controller Layer:** Handles HTTP requests via the Gin framework.  
- **Service Layer:** (prepared for future business logic extensions).  
- **Model Layer:** Defines entities using GORM.  
- **View Layer:** Server-rendered HTML templates with Bootstrap 5.  

A key part of the design is the **use of the Strategy Pattern** to implement flexible and extensible  
search functionality (e.g., search by ID or title). This pattern allows additional strategies —  
such as searching by author or ISBN — to be easily added later without modifying controller logic.

---

### Project Goal
To practice **Clean Architecture** and **Go best practices** by building a modular, testable CRUD web application  
that demonstrates:
- Interface-driven development  
- SOLID design principles  
- Conscious use of **Design Patterns** (e.g., Strategy Pattern)  
- Maintainable and reusable code structure  

---

### Visual Assets
The visual assets (logos, icons, background watermark) were generated using **AI tools (ChatGPT / DALL·E by OpenAI)**.  
They were created and used solely for **learning and demonstration purposes** within this **personal pet project**.

---

### Design Patterns Implemented
- **Strategy Pattern** → For dynamic book search (by ID or title)  
  → Enables easy extension with new search strategies in the future  

---

### License
This project is open-source and available for educational and portfolio purposes.
