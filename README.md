# README: Example Project for `oapi-codegen` with Chi Router
This project demonstrates how to integrate the OpenAPI code generation tool [](https://github.com/oapi-codegen/oapi-codegen)`oapi-codegen` with the [Chi router](https://github.com/go-chi/chi).
It showcases:
- OpenAPI-driven code generation for request handlers.
- Separation of private and public routes.
- JWT-based route protection for private routes.

The project is minimal and acts as an educational example for utilizing `oapi-codegen` with middleware like JWT authentication in a Chi server.
## Prerequisites
To run this project, you need:
- **Go** (>=1.20)
- Access to `jq`:
``` bash
  # On macOS (using brew)
  brew install jq
  
  # On Linux (Ubuntu/Debian):
  sudo apt install jq
```
## Features
This project defines:
1. **Public route `/health`:** Open to the public (no authorization required).
2. **Private route `/api/private`:** Protected using JWT tokens, requiring a valid authorization token to access.
3. **Authorization endpoint `/api/auth`:** Generates a JWT token for the given user.

## Quick Start
### 1. Clone the repository
Clone this project to your machine:
``` bash
git clone <repository_url>
cd <repository_folder>
```
### 2. Run the project locally
To start the server, use the following Makefile command, which handles code generation and runs the application:
``` bash
make run
```
The server will start and be available at `http://localhost:8080`.
### 3. Check public `/health` route
To verify the server is running, you can query the `/health` route:
``` bash
curl -X GET http://localhost:8080/health -i
```
Expected output:
``` http
HTTP/1.1 200 OK
Content-Length: 0
```
Alternatively, you can use the provided Makefile command:
``` bash
make health
```
It will print:
``` 
/health route returned status 200
```
### 4. Authorization flow
The project includes a JWT-based authorization mechanism. To test:
- Use `/api/auth` to get a JWT token for a specified username.
- Use this token to authenticate requests to the private `/api/private` route.

You can automate this process using the provided `make auth` command. Below is the step-by-step breakdown:
#### **Step 1: Get JWT Token**
Run:
``` bash
curl -X POST http://localhost:8080/api/auth \
  -H "Content-Type: application/json" \
  -d '{"username": "testUserName", "password": "pass"}'
```
Expected output:
``` json
{ "token": "<jwt_token>" }
```
#### **Step 2: Use JWT Token for Private Route**
Use the token from Step 1 in the `Authorization` header when querying the protected `/api/private` route:
``` bash
curl -X GET http://localhost:8080/api/private \
  -H "Authorization: Bearer <jwt_token>"
```
Expected output:
``` json
{ "username": "testUserName" }
```
#### Automate with Makefile
Alternatively, you can automate both steps using the Makefile:
``` bash
make auth
```
Expected output:
``` 
Get /api/auth authorization...
Token received: <jwt_token>
Get private /api/private route...
API username (testUserName) matches expected username (testUserName)
```
### 5. Test All Routes at Once
To test both public and private access (routes `/health` and `/api/private`), run:
``` bash
make check-routes
```
Expected output:
``` 
Check /health route...
/health route returned status 200
Get /api/auth authorization...
Token received: <jwt_token>
Get private /api/private route...
API username (testUserName) matches expected username (testUserName)
```
## File Overview
### `main.go`
- Sets up the Chi router (`chi.NewMux()`).
- Configures `middleware.Logger` and `middleware.Recoverer` for logging and error recovery.
- Adds JWT-based authorization middleware (`CreateAuthMiddleware`).
- Routes implemented using the auto-generated OpenAPI handler (`api.HandlerFromMux`).

### `Makefile`
Provides convenient commands for:
- **Code generation:** `make generate`
- **Run server:** `make run`
- **Test specific routes:**
    - Public `/health` (`make health`)
    - Private `/api/private` with auth (`make auth`)
    - Check both (`make check-routes`).

### `/api` and `/server`
- `api`: Contains OpenAPI-generated code.
- `server`: Implements the `api` interface to handle requests.

## Project Structure

```
.
├── Makefile                    # Automation tasks
├── README.md                   # This README.md
├── api
│   ├── api.yaml                # API description in openapi 3.0 format
│   ├── config.yaml             # Codegenerate config
│   └── doc.go                  # Codegenerate tool
├── go.mod                      # Module definition
├── go.sum                      # Dependencies
├── jwt
│   └── jwt.go                  # JWT tools
├── main.go                     # Entry point
├── middleware      
│   └── middleware.go           # Middleware with auth
├── requests.http               # HTTP tests, same as in Makefile
├── server
│   ├── httputils.go            # HTTP utils for responses
│   └── server.go               # Implementation for generated API server
├── service
│   └── service.go              # Your own service for user handle   
└── tools
   └── tools.go                 # Just tool to import openapi-codegen
```
## License
This project is provided as an open-source example and licensed under the [MIT License](LICENSE)
