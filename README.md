# 04_Projects

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

## Purpose

This repository acts as a **project index and portfolio entry point**.

Each listed project:
- lives in its own GitHub repository
- has a clean, focused commit history
- follows real-world backend development practices
