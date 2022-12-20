package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"github.com/PullRequestInc/go-gpt3"
	"github.com/spf13/cobra"
)

func main() {
	log.SetOutput(new(NullWriter))
	apiKey, ok := os.LookupEnv("API_KEY")
	if !ok {
		panic("Missing API_KEY")
	}

	client := gpt3.NewClient(apiKey)

	rootCmd := &cobra.Command{
		Use:   "chatgpt",
		Short: "Chat with ChatGPT in console.",
		Run: func(cmd *cobra.Command, args []string) {
			scanner := bufio.NewScanner(os.Stdin)

			for {
				fmt.Print(">> ")

				if !scanner.Scan() {
					break
				}

				question := scanner.Text()

				switch question {
				case "quit", "exit":
					return

				default:
					func() {
						ctx, cancel := context.WithTimeout(ctx, 10*time.Second)
						defer cancel()

						response, err := getResponse(ctx, client, question)
						if err != nil {
							fmt.Println(err)
						}

						fmt.Println(response)
						fmt.Println()
					}()
				}
			}
		},
	}

	rootCmd.Execute()
}

func getResponse(ctx context.Context, client gpt3.Client, question string) (response string, err error) {
	sb := strings.Builder{}

	err = client.CompletionStreamWithEngine(
		ctx,
		gpt3.TextDavinci003Engine,
		gpt3.CompletionRequest{
			Prompt: []string{
				question,
			},
			MaxTokens:   gpt3.IntPtr(3000),
			Temperature: gpt3.Float32Ptr(0),
		},
		func(resp *gpt3.CompletionResponse) {
			text := resp.Choices[0].Text

			sb.WriteString(text)
		},
	)
	if err != nil {
		return "", err
	}

	response = sb.String()
	response = strings.TrimLeft(response, "\n")

	return response, nil
}

type NullWriter int

func (NullWriter) Write([]byte) (int, error) {
	return 0, nil
}
