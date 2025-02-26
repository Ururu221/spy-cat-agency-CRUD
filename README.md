# Spy Cat Agency

## Overview
Spy Cat Agency (SCA) is a system for managing spy cats. It allows the agency to handle spy cats, assign missions, track targets, and collect intelligence data.

## Features
- **Manage Spy Cats**: Add, update, remove, and list spy cats.
- **Missions & Targets**: Assign cats to missions, complete missions, and track targets.
- **Tech Stack**: Built with Go, Echo framework, and PostgreSQL.

## Installation & Setup

### Prerequisites
- **Go** (latest version)
- **Docker & Docker Compose**

### Steps to Run
1. **Clone the repository:**
   ```sh
   git clone https://github.com/Ururu221/spy-cat-agency-CRUD
   ```
2. **Start services:**
   ```sh
   docker compose up -d
   ```
3. **Run the application:**
   ```sh
   go run main.go
   ```

Your API will be available at `http://localhost:1323`.

## API Endpoints
All endpoints are listed in the Postman collection: `spy-cat-agency.postman_collection.json`.

### Examples:
#### **Cats:**
- `POST /cats` — Create a spy cat
- `GET /cats` — List all spy cats
- `GET /cats/:id` — Get a spy cat by ID
- `DELETE /cats/:id` — Delete a spy cat
- `PUT /cats/:id` — Update a spy cat’s salary

#### **Missions:**
- `POST /missions` — Create a mission with targets
- `GET /missions` — List all missions
- `GET /missions/:id` — Get a mission by ID
- `DELETE /missions/:id` — Delete a mission
- `PUT /missions/complete/:id` — Mark a mission as complete
- `PUT /missions/assign-cat` — Assign a cat to a mission

#### **Targets:**
- `PUT /targets/complete/:id` — Mark a target as complete
- `PUT /targets/update-note/:id` — Update target notes
- `PUT /targets/delete-from-mission/:id` — Remove a target from a mission
- `PUT /targets/add-to-mission/:mission_id` — Add a target to a mission

## Testing
For API testing, refer to the Postman collection: `spy-cat-agency.postman_collection.json`.



