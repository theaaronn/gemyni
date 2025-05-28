package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"

	"github.com/charmbracelet/glamour"
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

func main() {
	arg := os.Args[1]
	jsondata := []byte(fmt.Sprintf(`{"contents":[{"parts":[{"text":"%s"}]}]}`, arg))

	req, err := http.NewRequest("POST", "https://generativelanguage.googleapis.com/v1beta/models/gemini-2.0-flash:generateContent?key=", bytes.NewBuffer(jsondata))
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

	var apiResponse Response

	err = json.Unmarshal(body, &apiResponse)
	if err != nil {
		fmt.Println("Failed unmarshalling json response:", err.Error())
	}

	text := apiResponse.Candidates[0].Content.Parts[0].Text

	fmt.Println("Status: ", resp.Status)
	// Usar glamour para renderizar el texto Markdown
	renderer, _ := glamour.NewTermRenderer(
		glamour.("NoTTy"), // Usar el tema "NoTTy"
		glamour.WithWordWrap(80),   // Ajusta el ancho según necesites
	)

	out, err := renderer.Render(text)
	if err != nil {
		fmt.Println("Failed rendering glamour:", err.Error())
		return
	} else {
		fmt.Print(out)
	}

}
