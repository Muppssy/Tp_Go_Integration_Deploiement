package apitests

import "bytes"
import "encoding/json"
import "net/http"
import "time"

type CatModel struct {
	Name      string `json:"name"`
	ID        string `json:"id,omitempty"`
	BirthDate string `json:"birthDate,omitempty"`
	Color     string `json:"color,omitempty"`
}

var baseUrl = "http://localhost:8080/api"

// Global client with a proper timeout
var client = &http.Client{Timeout: 10 * time.Second}

// Wrapper to HTTP API calls, does the error handling and JSON decoding
func call(method, path string, reqBody any, code *int, result any) error {

	jsonBody, err := json.Marshal(reqBody)
	if err != nil {
		return err
	}

	req, _ := http.NewRequest(method, baseUrl + path,  bytes.NewBuffer(jsonBody))

    // send the request
    res, err := client.Do(req)
	if err != nil {
		return err
	}
	defer res.Body.Close()

	if code != nil {
		*code = res.StatusCode
	}

	if result != nil {
		err = json.NewDecoder(res.Body).Decode(result)
	}

	return err
}
