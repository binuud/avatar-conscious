package apitools

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/binuud/avatar-conscious/pkg/tools/toolMeta"
	"github.com/openai/openai-go/v3"
	"github.com/sirupsen/logrus"
)

// LLM returns a json string
// which has to be parsed into a equivalent structure
// this provides the input structure
type ApiGetTool_Input struct {
	// endpoint of the query
	Endpoint string `json:"endpoint"`
	// query string, this can be a string or a map of string, string
	QueryString string `json:"querystring"`
}

// implements the CommonToolInterface
type ApiGetTool struct {
	toolMeta.CommonToolInterface
}

func NewApiGetTool() *ApiGetTool {
	return &ApiGetTool{}
}

// Name() returns the tool name
// this is the name which is identified by LLM during tool call decision
func (t ApiGetTool) Name() string {
	return "get_http_request"
}

// Description()
// returns the description
// which is used by LLM, to decide if this tool has to be called,
// for the given prompt
func (t ApiGetTool) Description() string {
	return "Get the response from a http server, using get request to an endpoint"
}

func (t ApiGetTool) InputSchema() string {
	return "Get the contents of the file, from the given path"
}

func (t ApiGetTool) InputSchemaAsOpenAiParam() openai.FunctionParameters {
	return openai.FunctionParameters{
		"type": "object",
		"properties": map[string]any{
			"endpoint": map[string]string{
				"type":        "string",
				"description": "The endpoint of the server with http or https schema, eg: https://www.example.com, http://localhost:9080",
			},
			"querystring": map[string]string{
				"type":        "string",
				"description": "The query string used for the api get request",
			},
		},
		"required": []string{"endpoint", "querystring"},
	}
}

func (t ApiGetTool) Handler(args string) (string, error) {
	logrus.Infof("Arguments received %v", args)
	var in ApiGetTool_Input
	if err := json.Unmarshal([]byte(args), &in); err != nil {
		return "", err
	}
	logrus.Infof("Arguments parsed into input structure %v", in)
	return GetRequest(in.Endpoint, in.QueryString)
}

func GetRequest(endpoint string, queryString string) (string, error) {

	queryMap := make(map[string]string, 0)
	if err := json.Unmarshal([]byte(queryString), &queryMap); err != nil {
		// cannot unmarshall the query to a map, probably its a plain string
		// add same to url
		if !strings.HasSuffix(endpoint, "/") {
			endpoint = endpoint + "/"
		}
		endpoint = endpoint + queryString
	}

	// Parse the base URL safely
	u, err := url.Parse(endpoint)
	if err != nil {
		return "", fmt.Errorf("endpoint url malformed %v", err)
	}
	// Construct the query string parameters
	params := url.Values{}
	if len(queryMap) > 0 {
		for k, item := range queryMap {
			params.Add(k, item)
		}
	}

	// Encode and append the query parameters to the URL
	u.RawQuery = params.Encode()
	logrus.Printf("Executing request to: %s\n", u.String())

	// Initialize an HTTP client with a timeout of 10 seconds
	// #TODO make this programmable in settings proto
	client := &http.Client{
		Timeout: 10 * time.Second,
	}

	// Execute the GET request
	resp, err := client.Get(u.String())
	if err != nil {
		return "", fmt.Errorf("HTTP request failed: %v\n", err)
	}
	// Always defer closing the response body to avoid resource leaks
	defer resp.Body.Close()

	// Verify the HTTP status code
	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("Server returned non-200 status: %d %s\n", resp.StatusCode, resp.Status)
	}

	// Read the response body
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", fmt.Errorf("Failed to read response body: %v\n", err)
	}
	logrus.Infof("GetRequest result got: %s", string(body))
	return string(body), nil
}
