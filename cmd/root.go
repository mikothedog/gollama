/*
Copyright © 2025 Miko the dog <gchaimowicz.dev@gmail.com>
*/
package cmd

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"time"

	"github.com/spf13/cobra"
)

const (
	DEFAULT_MODEL = "qwen3:8b"
	OLLAMA_URL    = "http://localhost:11434/api/generate"
)

type GenerateRequest struct {
	Model  string `json:"model"`
	Prompt string `json:"prompt"`
	Stream bool   `json:"stream"`
}

type GenerateResponse struct {
	Response  string `json:"response"`
	CreatedAt string `json:"created_at"`
	Model     string `json:"model"`
	Done      bool   `json:"done"`
}

func NewGenerateRequest(prompt []string) GenerateRequest {
	return GenerateRequest{
		Model:  DEFAULT_MODEL,
		Prompt: fmt.Sprintln(prompt),
		Stream: true,
	}
}

func runApp(prompt []string) error {
	if len(prompt) == 0 {
		return fmt.Errorf("no prompt provided")
	}

	requestBody := NewGenerateRequest(prompt)
	jsonBody, err := json.Marshal(requestBody)
	if err != nil {
		return fmt.Errorf("failed to marshal request body: %w", err)
	}
	req, err := createHttpRequest(jsonBody)
	if err != nil {
		return err
	}

	if err := startOllama(); err != nil {
		return err
	}

	resp, err := sendHttpRequest(req)
	if err != nil {
		return err
	}

	done := make(chan struct{})

	go func() {
		<-done
		if closeErr := resp.Body.Close(); closeErr != nil {
			log.Printf("failed to close response body: %v", closeErr)
		}
	}()

	err = checkResponse(resp, done)
	if err != nil {
		return err
	}

	err = printResponse(resp)
	if err != nil {
		return err
	}

	stopOllama()

	return nil
}

func stopOllama() {
	exec.Command("ollama", "stop", DEFAULT_MODEL)
	exec.Command("killall", "ollama")
}

func startOllama() error {
	cmd := exec.Command("ollama", "serve")

	if err := cmd.Start(); err != nil {
		return fmt.Errorf("failed to start 'ollama serve' command: %w", err)
	}
	time.Sleep(1 * time.Second)
	return nil
}

func printResponse(resp *http.Response) error {
	decoder := json.NewDecoder(resp.Body)
	for {
		var response GenerateResponse
		if err := decoder.Decode(&response); err != nil {
			if err == io.EOF {
				break
			}
			log.Printf("failed to decode JSON from stream: %v", err)
			break
		}

		fmt.Print(response.Response)

		if response.Done {
			break
		}
	}

	if closeErr := resp.Body.Close(); closeErr != nil {
		log.Printf("failed to close response body: %v", closeErr)
	}
	return nil
}

func checkResponse(resp *http.Response, done chan struct{}) error {
	if resp.StatusCode != http.StatusOK {
		body, _ := io.ReadAll(resp.Body)
		close(done)
		return fmt.Errorf("received non-200 response: %d, body: %s", resp.StatusCode, string(body))
	}
	return nil
}

func sendHttpRequest(req *http.Request) (*http.Response, error) {
	client := &http.Client{}
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		return nil, fmt.Errorf("failed to send http request: %w", err)
	}
	return resp, nil
}

func createHttpRequest(jsonBody []byte) (*http.Request, error) {
	req, err := http.NewRequestWithContext(context.Background(), "POST", OLLAMA_URL, bytes.NewBuffer(jsonBody))
	if err != nil {
		return nil, fmt.Errorf("failed to create http request: %w", err)
	}

	req.Header.Set("Content-Type", "application/json")
	return req, nil
}

var rootCmd = &cobra.Command{
	Use:   "gollama",
	Short: "Interact with Ollama using the llama3.1:8b model",
	Long:  `interact with Ollama using the llama3.1:8b model`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := runApp(args); err != nil {
			fmt.Println(err)
		}
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		os.Exit(1)
	}
}

func init() {
	// Here you will define your flags and configuration settings.
	// Cobra supports persistent flags, which, if defined here,
	// will be global for your application.

	// rootCmd.PersistentFlags().StringVar(&cfgFile, "config", "", "config file (default is $HOME/.gollama.yaml)")

	// Cobra also supports local flags, which will only run
	// when this action is called directly.
	rootCmd.Flags().BoolP("toggle", "t", false, "Help message for toggle")
}
