# SQLKV - Key-Value Store with SQL Backend

SQLKV is a key-value store implemented in Go, with an SQL backend (using SQLite). The store exposes an HTTP interface for interacting with the key-value pairs, and it's designed to be easily extensible for SDK or CLI interaction.

## Features

- **RESTful API**: Store and retrieve key-value pairs over HTTP.
- **SQLite Backend**: Stores key-value pairs in a SQLite database.
- **Middleware**: Includes middleware for injecting the app configuration into Echo requests.
- **Health Check**: Endpoint to check the health of the server.

## Getting Started

### Prerequisites

- **Go**: You need to have Go installed (version 1.18 or later).
- **SQLite**: SQLite is used as the database backend.

### Clone the Repository

```bash
git clone https://github.com/yourusername/sqlkv.git
cd sqlkv
```

### Configure the Database

Create a new SQLite database file named `sqlkv.db` in the root directory of the project. Or you can run the server and call the seed endpoint to create the database.

### Configure the Server

Create a `.env` file in the root directory of the project and add the following environment variables:

- `SCHEMA_FILE_PATH`: The path to the SQL schema file. Default value is `./schema.sql`.

Example `.env` file:

```bash
SCHEMA_FILE_PATH=./schema.sql
```

### Run the Server

To run the server, use the following command:

```bash
go run main.go
```

The server will start listening on port 8000 by default.

## Usage

### HTTP API

The server exposes the following endpoints:

- `GET /`: Health check endpoint.
- `GET /seed`: Seed the database with some initial data.
- `GET /kv/get/{key}`: Get the value for a given key.
- `POST /kv/set`: Set the value for a given key. Requires a `key` and `value` parameter in the request body as JSON. Returns the value of the key.

### CLI

You can use the CLI to interact with the server. To install the CLI, run the following command:

```bash
go install github.com/yourusername/sqlkv@latest
```

Once the CLI is installed, you can use the following commands:

- `get [key]`: Get the value for a given key.
- `set [key] [value]`: Set the value for a given key.

## Contributing

Contributions are welcome! If you find any issues or have suggestions for improvements, please open an issue or submit a pull request on the GitHub repository.

