# Go API Template

This is a template for building a RESTful API using Go. The project is structured to promote clean code and ease of maintenance.

## Project Structure

- **api/**: Contains the API layer, including controllers and middleware.
  - **controller/**: Handles HTTP requests and responses.
  - **middleware/**: Contains middleware functions for the API.
- **bootstrap/**: Initialization and configuration setup.
- **config/**: Configuration files and settings.
- **domain/**: Contains domain models and interfaces.
- **internal/**: Internal application code, including business logic.
  - **token_util/**: Utility functions for token management.

## Getting Started

### Prerequisites

- Go 1.23.5 or higher
- Docker (optional, for containerization)
- Make (optional, for running make commands)

### Installation

1. Clone the repository:
    ```sh
    git clone https://github.com/netojso/go-api-template.git
    cd go-api-template
    ```

2. Install dependencies:
    ```sh
    go mod tidy
    ```

3. Copy the example environment file and update the environment variables:
    ```sh
    cp .env.example .env
    ```

### Running the Application

1. Run the application:
    ```sh
    go run cmd/main.go
    ```

2. Alternatively, you can use Air for live reloading during development:
  ```sh
  air
  ```

### Contributing

Feel free to submit issues, fork the repository and send pull requests!

### License

This project is licensed under the MIT License.
