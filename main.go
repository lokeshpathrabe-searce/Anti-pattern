package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

const (
	apiURL   = "https://api.openai.com/v1/chat/completions"
	apiKey   = "API_KEY"
	model    = "gpt-3.5-turbo"
)

func main() {
	message := map[string]interface{}{
		"model":    model,
		"messages": []map[string]string{
			{"role": "system", "content": "You are a bigquery Optimizer, skilled in explaining solution to optimize the query."},
			{"role": "user", "content": "SELECT t1.col1 FROM `project.dataset.table1` t1 WHERE t1.col2 not in (select col2 from `project.dataset.table2`);"},
		},
	}

	requestBody, err := json.Marshal(message)
	if err != nil {
		fmt.Println("Error encoding JSON:", err)
		return
	}

	client := &http.Client{}
	req, err := http.NewRequest("POST", apiURL, bytes.NewBuffer(requestBody))
	if err != nil {
		fmt.Println("Error creating request:", err)
		return
	}

	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+apiKey)

	response, err := client.Do(req)
	if err != nil {
		fmt.Println("Error making request:", err)
		return
	}
	defer response.Body.Close()

	if response.StatusCode == http.StatusOK {
		body, err := io.ReadAll(response.Body)
		if err != nil {
			fmt.Println("Error reading response body:", err)
			return
		}

		var result map[string]interface{}
		if err := json.Unmarshal(body, &result); err != nil {
			fmt.Println("Error decoding JSON:", err)
			return
		}

		generatedText := result["choices"].([]interface{})[0].(map[string]interface{})["message"].(map[string]interface{})["content"].(string)
		fmt.Println(generatedText)
	} else {
		fmt.Println("Error:", response.Status)
		body, _ := io.ReadAll(response.Body)
		fmt.Println(string(body))
	}
}