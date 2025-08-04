# Gollama

A simple command-line interface (CLI) tool written in Go for interacting with a local Ollama server. It allows you to send a prompt and receive a streamed response directly from a local language model.

---

## Features

* **Easy Interaction**: Send a prompt to your local Ollama server directly from the command line.
* **Streamed Responses**: The tool handles the streaming API from Ollama, providing real-time output as it's generated.
* **Automatic Ollama Management**: Automatically starts the Ollama server if it's not running and attempts to stop it when the command is finished.
* **Default Model**: Uses the `gemma:2b` model by default, but this can be easily modified in the source code.

---

## Prerequisites

Before using this tool, you must have [Ollama](https://ollama.ai/) installed and running on your machine.

---

## Installation

To build and run this tool, you'll need the Go toolchain installed.

1.  **Clone the repository:**
    ```bash
    git clone [https://github.com/your-username/gollama.git](https://github.com/your-username/gollama.git)
    cd gollama
    ```

2.  **Build the executable:**
    ```bash
    go build -o gollama
    ```

After building, the `gollama` executable will be in your current directory. You may want to move it to a directory in your system's `PATH` to run it from anywhere.

---

## Usage

To use the tool, simply run the executable with your prompt as command-line arguments.

> ```bash
> ./gollama "Write a short story about a robot who discovers music."
> ```

The response will be streamed directly to your terminal.

---

## Code Overview

The core logic is found within the `cmd/` directory, following the standard Go CLI application structure using the `cobra` library.

### `main.go`

This file initializes and executes the root command.

### `cmd/root.go`

* **`rootCmd`**: The main command for the application. It takes the user's input from the command line and passes it to the `runApp` function.

* **`runApp(prompt []string)`**: The main function that orchestrates the entire process. It performs the following steps:
    1.  Checks if a prompt was provided.
    2.  Creates an HTTP request body for the Ollama API, setting the default model and enabling streaming.
    3.  Attempts to start the Ollama server.
    4.  Sends the HTTP request to the Ollama API.
    5.  Handles the streamed response by decoding JSON chunks and printing the output to the console.
    6.  Attempts to stop the Ollama server after the response is complete.

### Key Functions

* `startOllama()`: Starts the `ollama serve` command in a separate process.
* `stopOllama()`: Sends a `killall` command to stop the Ollama process.
* `createHttpRequest()`: Constructs the HTTP `POST` request with the correct headers and URL.
* `sendHttpRequest()`: Executes the HTTP request and returns the response.
* `checkResponse()`: Verifies that the HTTP response status code is `200 OK`.
* `printResponse()`: Decodes the streamed JSON response from the server and prints the generated text to standard output.
* `NewGenerateRequest()`: A constructor for the `GenerateRequest` struct, which formats the user's prompt for the Ollama API.

---

## Customization

You can easily change the default model by modifying the `DEFAULT_MODEL` constant in the `cmd/root.go` file. For example, to use `llama3:8b`:

```go
const (
    DEFAULT_MODEL = "llama3:8b"
    OLLAMA_URL    = "http://localhost:11434/api/generate"
)
```

Remember to rebuild the executable after making any changes.

