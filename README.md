# CHATU - A Chat Application

CHATU is a real-time chat application built using Go, Gorilla WebSockets, MongoDB, Docker, Kubernetes, and JWT tokens for authentication.

## Features

- Real-time messaging using WebSockets
- User authentication with JWT tokens
- Secure communication with self-signed certificates
- Multi-room chat functionality
- Dockerized setup for easy deployment
- Kubernetes configurations for scalable management

## Installation

1. **Clone the repository:**
   ```bash
   git clone https://github.com/Arpit529Srivastava/CHATU.git
   cd CHATU
 ### Generate self-signed certificates:
```
bash certgen.bash
```
### Start the server:

```
go run *.go
```
### Access the application: 
### Open your browser and navigate to 
```
https://localhost:8080.
```

### Login: Use the following credentials to log in:

```
Username: percy
Password: 123
```
### 
Send a message: Type your message in the input field and click send.

## Project Structure
- main.go: Entry point of the application.
- manager.go: Manages WebSocket connections and routes events.
- client.go: Handles client connections and message transmission.
- frontend/: Contains the frontend code for the chat application.
- otp.go: Manages OTP generation and verification.

### Dependencies
- Go
- Gorilla WebSocket
- MongoDB
- Docker
- Kubernetes
- License
### This project is licensed under the MIT License.