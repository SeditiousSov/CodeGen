package main

import (
	"os"
	"fmt"
	"bytes"
	"strings"
	"context"

	openai "github.com/sashabaranov/go-openai"
)

gptkey = `YOUR_KEY_HERE`

func GPTTTS(tosend string) {
	client := openai.NewClient(gptkey)
	resp, err := client.CreateSpeech(context.Background(), openai.CreateSpeechRequest {
			Model: openai.TTSModel1,
			Input: tosend,
			Voice: openai.VoiceNova,
		},
	)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	buf := new(bytes.Buffer)
	num, err := buf.ReadFrom(resp)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println(num)
	f, err := os.Create("myfile.mp3")
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	_, err = buf.WriteTo(f)
	f.Close()
}


func TextChat(tosend string) {
	client := openai.NewClient(gptkey)
	resp, err := client.CreateChatCompletion(context.Background(), openai.ChatCompletionRequest {
			Model: openai.GPT3Dot5Turbo,
			Messages: []openai.ChatCompletionMessage{
				{
					Role: openai.ChatMessageRoleUser,
					Content: tosend,
				},
			},
		},
	)

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	response := resp.Choices[0].Message.Content

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
