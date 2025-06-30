
# Blog API

## Environment Variables

Create a `.env` file in the root of the project with the following variables:

```env
POSTGRES_USER=your_postgres_username
POSTGRES_PASSWORD=your_postgres_password
POSTGRES_DB=your_postgres_database
POSTGRES_HOST=localhost
```

These variables configure the connection to your PostgreSQL database.

---

## How the Application Works

This is a RESTful API for managing a simple blogging platform with posts and comments.

### API Endpoints

- **GET /api/posts**  
  Returns a paginated list of blog posts with the count of comments for each post.  
  Optional query parameters:
    - `limit` (default 10)
    - `offset` (default 0)
    - `title` (filter posts by title substring)

- **POST /api/posts**  
  Create a new blog post by providing a JSON body with `title` and `content`.

- **GET /api/posts/{id}**  
  Retrieve a specific blog post by its ID, including full details and all associated comments.

- **POST /api/posts/{id}/comments**  
  Add a new comment to a specific blog post. Provide JSON with the `content` field.

### Database

The app uses PostgreSQL to store blog posts and their comments. It connects automatically using the credentials specified in the `.env` file.

### Logging

The app logs key events and errors to aid in monitoring and troubleshooting.

### Swagger UI

API documentation is available at:  
[http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html)

---

## Makefile Commands

Use the following commands to manage your project easily:

- `make startup`  
  Build and start containers in detached mode.

- `make stop`  
  Stop containers.

- `make down`  
  Stop and remove containers, networks, and volumes.

- `make logs`  
  Show logs. Use `SERVICE=<name>` to filter logs for a specific service.

- `make migrate`  
  Run database migrations locally.

- `make rollback`  
  Rollback last database migration locally.

- `make migrate-docker`  
  Run database migrations using Docker container.

- `make rollback-docker`  
  Rollback last migration using Docker container.

- `make init`  
  Run the Go application.

- `make restart`  
  Restart containers.

- `make test`  
  Run all tests.

---

## Running the Application

1. Ensure you have Docker and Docker Compose installed.

2. Set your environment variables in `.env`.

3. Run `make startup` to start the app.

4. Access the API at `http://localhost:8080`.

5. Access Swagger UI at `http://localhost:8080/swagger/index.html`.

---
