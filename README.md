# Gollama

A simple command-line interface (CLI) tool written in Go for interacting with a local Ollama server. It allows you to send a prompt and receive a streamed response directly from a local language model.

---

## Features

* **Easy Interaction**: Send a prompt to your local Ollama server directly from the command line.
* **Streamed Responses**: The tool handles the streaming API from Ollama, providing real-time output as it's generated.
* **Automatic Ollama Management**: Automatically starts the Ollama server if it's not running and attempts to stop it when the command is finished.
* **Default Model**: Uses the `qwen:8b` model by default, but this can be easily modified in the source code.

---

## Prerequisites

Before using this tool, you must have [Ollama](https://ollama.ai/) installed and running on your machine.

---

## Installation

To build and run this tool, you'll need the Go toolchain installed.

1.  **Clone the repository:**
    ```bash
    git clone https://github.com/mikothedog/gollama.git
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

## Customization

You can easily change the default model by modifying the `DEFAULT_MODEL` constant in the `cmd/root.go` file. For example, to use `llama3:8b`:

```go
const (
    DEFAULT_MODEL = "llama3:8b"
    OLLAMA_URL    = "http://localhost:11434/api/generate"
)
```

Remember to rebuild the executable after making any changes.

