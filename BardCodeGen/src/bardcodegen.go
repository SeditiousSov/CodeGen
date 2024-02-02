package main

import (
	"io"
	"os"
	"fmt"
	"bytes"
	"strings"

	"net/http"
	"encoding/json"
)

bardkey = `YOUR_KEY_HERE`

type Payload struct {
	Contents []struct {
		Parts []struct {
			Text string `json:"text"`
		} `json:"parts"`
	} `json:"contents"`
}

type Reply struct {
	Candidates []struct {
		Content struct {
			Parts []struct {
				Text string `json:"text"`
			} `json:"parts"`
			Role string `json:"role"`
		} `json:"content"`
		FinishReason  string `json:"finishReason"`
		Index         int    `json:"index"`
		SafetyRatings []struct {
			Category    string `json:"category"`
			Probability string `json:"probability"`
		} `json:"safetyRatings"`
	} `json:"candidates"`
	PromptFeedback struct {
		SafetyRatings []struct {
			Category    string `json:"category"`
			Probability string `json:"probability"`
		} `json:"safetyRatings"`
	} `json:"promptFeedback"`
}

func TextChat(tosend string) {
	var data Payload
	reply := Reply{}


	data.Contents = make([]struct {
		Parts []struct {
			Text string `json:"text"`
		} `json:"parts"`
	}, 1)


	data.Contents[0].Parts = make([]struct {
		Text string `json:"text"`	
	}, 1)

	data.Contents[0].Parts[0].Text = tosend

	payloadBytes, err := json.Marshal(data)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	body := bytes.NewReader(payloadBytes)

	req, err := http.NewRequest("POST", "https://generativelanguage.googleapis.com/v1beta/models/gemini-pro:generateContent?key=" + bardkey, body)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	req.Header.Set("Content-Type", "application/json")
	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer resp.Body.Close()

	jsn, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	err = json.Unmarshal(jsn, &reply)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	response := reply.Candidates[0].Content.Parts[0].Text 

	if ! strings.Contains(response, "```") {
		fmt.Println(response)
		return
	}


	flag := false
	parts := strings.Split(response, "\n")
	for _, part := range parts {
		if strings.HasPrefix(part, "```") && ! flag {
			flag = true
			continue
		} else if strings.HasPrefix(part, "```") && flag  {
			flag = false
			continue
		}

		if flag {
			fmt.Println(part)
		}
	}
}


func main() {

	if len(os.Args) != 2 {
		fmt.Println("Usage: " + os.Args[0] + ` "Can you please generate an ansible playbook for tagging a virtual machine on vmware"`)
		return
	}

	tosend := os.Args[1]
	TextChat(tosend)
}
