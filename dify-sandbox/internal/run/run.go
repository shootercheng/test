package run

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

var fileSuffixToLangMap = map[string]string{
	"py": "python3",
	"js": "nodejs",
}

var testFileDir = []string{"python", "nodejs"}

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

func (c *SandboxClient) SendCode(reqBody *RequestBody) ([]byte, error) {
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

	client := NewSandboxClient(requestHost, apiKey)
	defer client.Close()

	for _, filePath := range testFileDir {
		entries, err := os.ReadDir(filePath)
		if err != nil {
			fmt.Println("Error reading directory:", err)
			continue
		}
		for _, entry := range entries {
			fileName := entry.Name()
			codePath := filePath + "/" + fileName
			fmt.Printf("test file path:%s\n", codePath)
			// Read code from file
			codeContent, err := readCodeFromFile(codePath)
			if err != nil {
				fmt.Printf("Error reading code file:%v\n", err)
				continue
			}
			fileNameArr := strings.Split(fileName, ".")
			var fileSuffix string
			if len(fileNameArr) > 1 {
				fileSuffix = fileNameArr[len(fileNameArr)-1]
			}
			language, ok := fileSuffixToLangMap[fileSuffix]
			if !ok {
				fmt.Printf("file name :%s error\n", fileName)
				continue
			}
			reqBody := &RequestBody{
				Language:      language,
				Code:          codeContent,
				Preload:       "",
				EnableNetwork: true,
			}
			resp, err := httpRequest(client, reqBody)
			if err != nil {
				fmt.Printf("http request error:%v\n", err)
				return
			}
			if resp.Data.Error != "" && strings.Contains(resp.Data.Error, "operation not permitted") {
				fmt.Println("syscall error, please check")
			}
		}
	}
}

func httpRequest(client *SandboxClient, reqBody *RequestBody) (*Response, error) {
	res, err := client.SendCode(reqBody)
	if err != nil {
		return nil, err
	}

	var resp Response
	err = json.Unmarshal(res, &resp)
	if err != nil {
		return nil, err
	}
	fmt.Printf("request result:%s\n", string(res))
	return &resp, nil
}

func readCodeFromFile(filename string) (string, error) {
	content, err := os.ReadFile(filename)
	if err != nil {
		return "", err
	}
	return string(content), nil
}
