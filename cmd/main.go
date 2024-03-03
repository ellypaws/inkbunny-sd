package main

import (
	"encoding/json"
	"fmt"
	"github.com/ellypaws/inkbunny/api"
	"golang.org/x/term"
	"inkbunny-sd/entities"
	"inkbunny-sd/llm"
	"log"
	"os"
	"regexp"
	"strings"
)

func main() {
	user := &api.Credentials{Sid: os.Getenv("SID")}
	err := login(user)
	if err != nil {
		log.Fatalf("Error logging in: %v", err)
	}
	log.Printf("Logged in as %s, sid: %s\n", user.Username, user.Sid)

	var submissionIDs string = os.Getenv("SUBMISSION_IDS")
	if submissionIDs == "" {
		fmt.Printf("Enter submission IDs (comma separated): ")
		fmt.Scanln(&submissionIDs)
	}

	log.Println("Getting submission details for IDs:", submissionIDs)
	details, err := user.SubmissionDetails(
		api.SubmissionDetailsRequest{
			SubmissionIDs:   submissionIDs,
			ShowDescription: api.Yes,
		})
	if err != nil {
		log.Printf("Error getting submission details: %v\n", err)
		return
	}

	log.Printf("Submission details: %v\n", details)

	var out []byte

	const maxRetries = 3
	for i := 0; i < maxRetries; i++ {
		log.Printf("Inferencing text to image (%v/3)\n", i+1)
		resp, err := llm.Localhost().Infer(&llm.Request{
			Messages: []llm.Message{
				llm.DefaultSystem,
				llm.UserMessage(details.Submissions[0].Description),
			},
			Temperature: 1.0,
			MaxTokens:   2048 * 2,
			Stream:      false,
		})
		if err != nil {
			log.Printf("Error inferencing, retrying (%v/3): %v\n", i+1, err)
			continue
		}

		message := fixJson(resp.Choices[0].Message.Content)
		textToImage, err := entities.UnmarshalTextToImageRequest([]byte(message))
		if err != nil {
			log.Printf("Error unmarshalling text to image: %v, retrying (%v/3)\n", err, i+1)
			continue
		}
		if textToImage.Prompt == "" {
			log.Printf("Prompt is empty, retrying (%v/3)\n", i+1)
			continue
		}
		out, _ = json.MarshalIndent(textToImage, "", "  ")
		log.Printf("Successfully inferenced %d tokens with %d retries\n", resp.Usage.CompletionTokens, i)
		log.Println(strings.ReplaceAll(message, "\n", ""))
		break
	}

	//save to json file
	err = os.WriteFile("text_to_image.json", out, 0644)
	if err != nil {
		log.Printf("Error writing text to image: %v\n", err)
		return
	}
	log.Println("Text to image saved to text_to_image.json")

	var l string
	fmt.Print("Logout? (y/[n]): ")
	fmt.Scanln(&l)
	if l == "y" {
		logout(user)
	}
}

func fixJson(content string) string {
	content = extractJson.FindString(content)
	content = removeComments.ReplaceAllString(content, "")
	content = escapeBackslash.ReplaceAllString(content, escapeBackslashReplacement)
	return content
}

var extractJson = regexp.MustCompile(`(?ms){.*}`)
var removeComments = regexp.MustCompile(`(?m)//.*$`)
var escapeBackslash = regexp.MustCompile(`\\+([()])`)

const escapeBackslashReplacement = `\\$1`

func login(u *api.Credentials) error {
	if u == nil {
		return fmt.Errorf("nil user")
	}
	if u.Sid == "" {
		user, err := loginPrompt().Login()
		if err != nil {
			return err
		}
		*u = *user
	}
	return nil
}

func loginPrompt() *api.Credentials {
	var user api.Credentials
	fmt.Print("Enter username: ")
	fmt.Scanln(&user.Username)
	fmt.Print("Enter password: ")
	bytePassword, err := term.ReadPassword(int(os.Stdin.Fd()))
	if err != nil {
		log.Fatal(err)
	}
	user.Password = string(bytePassword)

	return &user
}

func logout(u *api.Credentials) {
	err := u.Logout()
	if err != nil {
		log.Fatalf("error logging out: %v", err)
	}
}
