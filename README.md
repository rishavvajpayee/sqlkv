# SqlKV - Key-Value Store with SQL Backend

SQLKV is a lightweight key-value store implemented in Go, utilizing SQLite as its backend. It provides a RESTful HTTP interface for easy interaction with key-value pairs and is designed for extensibility, making it suitable for SDK or CLI integration.

## Features

- **RESTful API**: Easily store and retrieve key-value pairs over HTTP.
- **SQLite Backend**: Efficiently stores data in a SQLite database.
- **Health Check Endpoint**: Monitor server health with a dedicated endpoint.

## Getting Started

### Prerequisites

Before you begin, ensure you have the following installed:

- **Go**: Version 1.18 or later.
- **SQLite**: Required for the database backend.

### Clone the Repository

To get started, clone the repository and navigate into the project directory:

```bash
git clone https://github.com/yourusername/sqlkv.git
cd sqlkv
```

### Configure the Database

You can create a new SQLite database file named `sqlkv.db` in the root directory of the project. Alternatively, run the server and call the seed endpoint to automatically create the database.

### Configure the Server

Upon the first run, the server will create a new SQLite database file named `sqlkv.db` in the root directory.

For additional configuration, create a `.env` file in the root directory with the following example:

```bash
SCHEMA_FILE_PATH=./schema.sql
```

### Run the Server

To start the server, execute the following command:

```bash
go run main.go
```

The server will listen on port **8000** by default.

## Usage

### HTTP API

The server exposes the following endpoints:

- **Health Check**: `GET /` - Check the health of the server.
- **Get Value**: `GET /kv/get/{key}` - Retrieve the value associated with a given key.
- **Set Value**: `POST /kv/set` - Set the value for a given key. Requires a JSON body with `key` and `value` parameters.

### Command-Line Interface (CLI)

You can interact with the server using the CLI. To install the CLI, run:

```bash
go install github.com/rishavvajpayee/sqlkv@latest
```

Once installed, you can use the following commands:

- `get [key]`: Retrieve the value for the specified key.
- `set [key] [value]`: Set the value for the specified key.

## Contributing

Contributions are welcome! If you encounter any issues or have suggestions for improvements, please open an issue or submit a pull request on the GitHub repository.

## License

This project is licensed under the MIT License. See the [LICENSE](LICENSE) file for details.