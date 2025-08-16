# Task Manager API

This is a simple task management API built using only the Go standard library and Docker.
This is a test task for LO completed by Dmitry Beskaev

---

## API Endpoints

### `/tasks`

- **GET** `/tasks` — Get a list of all tasks.
- **GET** `/tasks?id={id}` — Get a specific task by its ID.
- **POST** `/tasks` — Create a new task.
- **DELETE** `/tasks?id={id}` — Delete a task by its ID.

---

## Running Locally

Make sure you have Go 1.25 installed.

To run the application locally:

`go run cmd/main.go`


By default the server will start on `http://localhost:8080`. You can change it's port by editing config.json file in the root directory

---

## Building and Running with Docker

Run this command in the project root (where your Dockerfile is): 
`docker-compose up --build`

---

## Example Requests

### Create a task
`curl --request POST
--url http://localhost:8080/tasks
--header 'content-type: application/json'
--data '{
"title": "test task title",
"description": "this is a test task"
}'`
### Get a task by id
`curl --request GET
--url 'http://localhost:8080/tasks?id=1'`
### Get all tasks
`curl --request GET
--url 'http://localhost:8080/tasks'`
### Delete a task by id
`curl --request DELETE
--url 'http://localhost:8080/tasks?id=1'`

---

## Notes

- The project uses **only the Go standard library** (`net/http`, etc.)—no external dependencies.
- Docker is used only for containerizing and running the app.
- The root endpoint `/` returns HTTP 200 OK.

---






