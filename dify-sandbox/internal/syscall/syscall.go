package syscall

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strings"
	"time"
)

const (
	SYSCALL_NUMS = 500
)

type RequestBody struct {
	Language      string `json:"language"`
	Code          string `json:"code"`
	Preload       string `json:"preload"`
	EnableNetwork bool   `json:"enable_network"`
}

type SandboxClient struct {
	client      *http.Client
	requestHost string
	apiKey      string
	baseURL     string
}

type Response struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    Data   `json:"data"`
}

type Data struct {
	Error  string `json:"error"`
	Stdout string `json:"stdout"`
}

func NewSandboxClient(requestHost string, apiKey string) *SandboxClient {
	transport := &http.Transport{
		MaxIdleConns:        100,
		MaxIdleConnsPerHost: 10,
		IdleConnTimeout:     90 * time.Second,
		TLSHandshakeTimeout: 10 * time.Second,
		DisableKeepAlives:   false,
	}

	client := &http.Client{
		Transport: transport,
		Timeout:   30 * time.Second,
	}

	return &SandboxClient{
		client:      client,
		requestHost: requestHost,
		apiKey:      apiKey,
		baseURL:     fmt.Sprintf("%s/v1/sandbox/run", requestHost),
	}
}

func (c *SandboxClient) SendCode(code string) ([]byte, error) {
	// Prepare request body
	reqBody := RequestBody{
		Language:      "python3",
		Code:          code,
		Preload:       "",
		EnableNetwork: false,
	}

	jsonData, err := json.Marshal(reqBody)
	if err != nil {
		fmt.Printf("Error marshaling JSON: %v\n", err)
		return nil, err
	}

	req, err := http.NewRequest("POST", c.baseURL, bytes.NewBuffer(jsonData))
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("X-Api-Key", c.apiKey)

	resp, err := c.client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	fmt.Printf("request json data:%s, response status: %s\n", jsonData, resp.Status)
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	return body, nil
}

func (c *SandboxClient) Close() error {
	c.client.Transport.(*http.Transport).CloseIdleConnections()
	return nil
}

func InvokeSandBoxRun() {
	// Read environment variables
	requestHost := os.Getenv("REQUEST_HOST")
	apiKey := os.Getenv("X_API_KEY")

	if requestHost == "" || apiKey == "" {
		fmt.Println("Error: REQUEST_HOST and X_API_KEY environment variables must be set")
		return
	}

	pythonFilePath := "python"
	entries, err := os.ReadDir(pythonFilePath)
	if err != nil {
		fmt.Println("Error reading directory:", err)
		return
	}

	client := NewSandboxClient(requestHost, apiKey)
	defer client.Close()

	for _, entry := range entries {
		codePath := pythonFilePath + "/" + entry.Name()
		fmt.Printf("test python file:%s\n", codePath)
		// Read code from file
		codeContent, err := readCodeFromFile(codePath)
		if err != nil {
			fmt.Println("Error reading code file:", err)
			continue
		}
		resp := httpRequest(client, codeContent)
		if resp.Data.Error != "" && strings.Contains(resp.Data.Error, "operation not permitted") {
			fmt.Println("syscall error, please check")
		}
	}
}

func httpRequest(client *SandboxClient, code string) Response {
	res, err := client.SendCode(code)
	if err != nil {
		fmt.Printf("request err:%s\n", err.Error())
	}

	var resp Response
	err = json.Unmarshal(res, &resp)
	if err != nil {
		panic(err)
	}
	fmt.Printf("request result:%s\n", string(res))
	return resp
}

func readCodeFromFile(filename string) (string, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(content), nil
}
