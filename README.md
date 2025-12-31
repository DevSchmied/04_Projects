# 04_Projects

This repository contains my personal and experimental software projects.  
It serves as a collection of pet projects, prototypes, and learning applications —  
mostly developed with Go, but also with other languages and tools used for practice and experimentation.

Each subfolder inside this repository represents a separate project with its own structure,  
README file, and Git commit history.

## Projects

- **bookmanager-go**  
  Go-based CRUD web application with Clean Architecture, JWT authentication, Redis caching, and Docker.  

- **2025-12-02 — URL Status & PDF Report Service**  
  Go-based web service that checks URL availability using a worker pool and generates PDF reports.  
  Includes concurrent processing with goroutines and channels, JSON-based storage, and graceful shutdown handling.  
  https://github.com/DevSchmied/2025-12-02-URLStatusAndPDFReportService


**Purpose:**  
To document my learning progress, demonstrate practical implementations, and explore different development patterns.
# 04_Projects

This repository serves as an index of my personal pet projects.

It does **not** contain source code itself.  
Instead, it provides an overview and direct links to standalone repositories where each project is developed and maintained independently.

The goal of this repository is to present my projects in a clean, structured way and make navigation easier for reviewers, recruiters, and collaborators.

---

## Projects

### BookManager-Go
Go-based CRUD web application demonstrating Clean Architecture principles.

**Highlights:**
- REST API built with Go
- Clean Architecture & SOLID principles
- JWT authentication
- Redis caching (optional)
- Simple web UI (HTML + Bootstrap)
- Docker & Docker Compose support

Repository:  
https://github.com/DevSchmied/bookmanager-go

---

### URL Status & PDF Report Service (2025-12-02)
Go-based backend service for checking URL availability and generating PDF reports.

**Highlights:**
- Concurrent URL checks (worker pool with goroutines and channels)
- JSON-based persistence
- PDF report generation
- Graceful shutdown (SIGINT/SIGTERM)
- Clean, modular project structure

Repository:  
https://github.com/DevSchmied/2025-12-02-url-status-and-pdf-report-service

---

## Purpose

This repository acts as a **project index / portfolio entry point**.

Each listed project:
- lives in its own GitHub repository
- has a clean commit history
- is developed and maintained independently

This structure reflects real-world development practices and makes individual projects easier to review and evaluate.
