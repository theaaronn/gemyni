/* 
TODO: update procedure to https://ai.google.dev/gemini-api/docs/text-generation#go using genai (including stream response)
TODO: branch to a crypto option so don't have api key exposed in .env (even though is local)
TODO: default to flash 2.0 if none model is specified (flag for model selection)
TODO: add release so can be summoned by just typing "gemyni" after go installing it as a package
? Potentially integrating image generation with flash 2.0
? Add tests
? Add more LLMs options for api calls (maybe would break the api response and error so integrate with standarized format)
*/

package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/charmbracelet/glamour"
	"github.com/joho/godotenv"
)

type Response struct {
	Candidates []Candidate `json:"candidates"`
}
type Candidate struct {
	Content Content `json:"content"`
}
type Content struct {
	Parts []Part `json:"parts"`
}
type Part struct {
	Text string `json:"text"`
}

type ResponseError struct {
	Err ApiError `json:"error"`
}
type ApiError struct {
	Code int `json:"code"`
	Message string `json:"message"`
	Status string `json:"status"`
}

func main() {
	query := os.Args[1]
	jsondata := fmt.Appendf(nil, `{"contents":[{"parts":[{"text":"%s"}]}]}`, query)	
	
	err := godotenv.Load(".env")
	if err != nil {
		fmt.Println("Failed loading .env file")
	}
	apiKey := os.Getenv("LLM_KEY")
	
	if apiKey == "" {
		fmt.Println("Missing api key in .env")
		return
	}
	url := fmt.Sprintf(flash20L, apiKey)
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsondata))
	if err != nil {
		fmt.Println("Failed creating request: ", err.Error())
		return
	}

	req.Header.Set("content-type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Failed doing request: ", err.Error())
	}

	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Failed reading response body: ", err.Error())
	}

	if resp.StatusCode != 200 {
		var respErr ResponseError
		err = json.Unmarshal(body, &respErr)
		if err != nil {
			fmt.Println("Failed unmarshalling json response error: ", err.Error())
			return
		}
		if len(respErr.Err.Message) <= 0 {
			fmt.Println("Response error message blank")
		} else {
			fmt.Println("Status:", respErr.Err.Code, respErr.Err.Status)
			fmt.Println(respErr.Err.Message)
		}
		return
	}

	var apiResponse Response

	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		fmt.Println("Failed unmarshalling json response:", err.Error())
	}
	if len(apiResponse.Candidates) <= 0 {
		fmt.Println("Request response length ", len(apiResponse.Candidates))
		return
	}

	text := apiResponse.Candidates[0].Content.Parts[0].Text

	out, err := glamour.Render(text, "dark")
	if err != nil {
		fmt.Println("Failed rendering glamour:", err.Error())
	} else {
		fmt.Print(out)
	}
}
