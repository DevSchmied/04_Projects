# 01_projects-index 

This repository serves as an index of my personal pet projects.

It does **not** contain source code itself.  
Instead, it provides an overview and direct links to standalone repositories,  
where each project is developed and maintained independently.

The goal of this repository is to present my projects in a clean, structured way  
and make navigation easier for reviewers, recruiters, and collaborators.

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

### Subscription Aggregation Service
REST API for managing and aggregating users’ online subscriptions.

**Highlights:**
- CRUDL operations for subscriptions
- Aggregation endpoint for total cost calculation by period
- Filtering by user and service
- PostgreSQL integration with migrations
- Swagger / OpenAPI documentation
- Docker-based deployment and CI pipeline

Repository:  
https://github.com/DevSchmied/subscription-aggregation-service

---

## Learning & Training Repositories

### 02_programming-languages
Collection of programming tasks and exercises completed during professional retraining (Umschulung in Germany) and ongoing self-education.

**Focus areas:**
- **Go (basics, concurrency, channels, select, slices, maps, strings and related topics)**
- **Java (primary programming language during professional retraining / Umschulung in Germany)**
- C
- C#
- SQL
- HTML & CSS
- PHP

Repository:  
https://github.com/DevSchmied/02_programming-languages

---

### 03_UML
Collection of UML and related software design diagrams created during professional retraining (Umschulung in Germany) and used to model software systems, business processes, and data structures.

**Diagram types include:**
- Flowcharts
- Nassi–Shneiderman diagrams
- ER diagrams & database schemas
- Use case diagrams
- Activity diagrams
- Sequence diagrams
- State diagrams
- Class diagrams

Repository:  
https://github.com/DevSchmied/03_UML

---

### 04_tools
Repository for studying and practicing developer tools commonly used in development.

The repository contains focused exercises aimed at understanding how specific tools work in practice and how they are integrated into backend systems.

**Current focus areas:**
- Redis (data structures, caching, TTL, Pub/Sub, streams, transactions, rate limiting, Lua scripting)

Repository:  
https://github.com/DevSchmied/04_tools

---

## Purpose

This repository acts as a **project index and portfolio entry point**.

Each listed project:
- lives in its own GitHub repository
- has a clean, focused commit history
- follows real-world backend development practices
