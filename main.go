package main

import (
	"bufio"
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"runtime"

	"github.com/google/generative-ai-go/genai"
	"github.com/joho/godotenv"
	"google.golang.org/api/option"
)

func main() {
	_, filename, _, ok := runtime.Caller(0)
	if !ok {
		panic("failed to get caller information")
	}
	dir := filepath.Dir(filename)
	err := godotenv.Load(filepath.Join(dir, ".env"))
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	ctx := context.Background()
	client, err := genai.NewClient(ctx, option.WithAPIKey(os.Getenv("API_KEY")))
	if err != nil {
		log.Fatal(err)
	}
	defer client.Close()

	model := client.GenerativeModel("gemini-1.5-flash")

	scanner := bufio.NewScanner(os.Stdin)

	for {
		fmt.Printf(">>")
		scanner.Scan()

		fmt.Println("|")
		fmt.Printf(">>")

		prompt := genai.Text(scanner.Text())
		iter := model.GenerateContentStream(ctx, prompt)
		for {
			resp, err := iter.Next()
			if err != nil {
				break
			}
			for _, cand := range resp.Candidates {
				for _, part := range cand.Content.Parts {
					fmt.Print(part)
				}
			}
		}
		fmt.Println()
	}
}
