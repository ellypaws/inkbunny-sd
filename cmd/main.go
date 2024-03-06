package main

import (
	"encoding/json"
	"fmt"
	"github.com/ellypaws/inkbunny-sd/entities"
	"github.com/ellypaws/inkbunny-sd/llm"
	"github.com/ellypaws/inkbunny-sd/utils"
	"github.com/ellypaws/inkbunny/api"
	"log"
	"os"
	"strings"
)

func main() {
	user := &api.Credentials{Sid: os.Getenv("SID")}
	err := login(user)
	if err != nil {
		log.Fatalf("Error logging in: %v", err)
	}
	log.Printf("Logged in as %s, sid: %s\n", user.Username, user.Sid)
	err = os.Setenv("SID", user.Sid)
	if err != nil {
		log.Printf("Error setting sid: %v", err)
	}

	if user.Username == "guest" {
		log.Println("Changing ratings for guest user")
		err = user.ChangeRating(api.Ratings{
			General:        true,
			Nudity:         true,
			MildViolence:   true,
			Sexual:         true,
			StrongViolence: true,
		})
		if err != nil {
			log.Printf("Could not change ratings for guest user: %v\n", err)
		} else {
			log.Printf("Ratings changed for guest user to all true")
		}
	}

	var submissionIDs string = os.Getenv("SUBMISSION_IDS")
	if submissionIDs == "" {
		fmt.Printf("Enter submission IDs (comma separated) or tags [tag:ai_generated]: ")
		fmt.Scanln(&submissionIDs)
	}

	if submissionIDs == "" {
		submissionIDs = "tag:ai_generated"
	}

	submissionIDs = strings.ReplaceAll(submissionIDs, " ", "")

	if strings.HasPrefix(submissionIDs, "tag:") {
		log.Printf("Searching for (5) submissions with tag: %s\n", strings.TrimPrefix(submissionIDs, "tag:"))
		searchResponse, err := user.SearchSubmissions(api.SubmissionSearchRequest{
			SubmissionIDsOnly:  false,
			SubmissionsPerPage: 5,
			Page:               1,
			Text:               strings.TrimPrefix(submissionIDs, "tag:"),
			Type:               api.SubmissionTypePicturePinup,
			OrderBy:            "views",
			Random:             true,
			Scraps:             "both",
		})
		if err != nil {
			log.Printf("Error searching for submission IDs: %v\n", err)
			return
		}

		if len(searchResponse.Submissions) == 0 {
			log.Println("No submissions found")
			main()
		}

		fmt.Printf("Found %d submissions with tag: %s\n", len(searchResponse.Submissions), strings.TrimPrefix(submissionIDs, "tag:"))
		fmt.Printf("%-*s  %-*s  %-*s  %s\n", 8, "ID", 31, "URL", 24, "Username", "Title")
		for _, submission := range searchResponse.Submissions {
			submissionIDs += submission.SubmissionID
			if submission != searchResponse.Submissions[len(searchResponse.Submissions)-1] {
				submissionIDs += ","
			}
			fmt.Printf(
				"%-*s  https://inkbunny.net/s/%-*s  %-*s  %s\n",
				8, // Adjusted width for Title + 2 for brackets
				submission.SubmissionID,
				8, // Adjusted width for Username + 4 for brackets and colon
				submission.SubmissionID,
				24,
				"["+submission.Username+"]: ",
				"["+submission.Title+"]",
			)
		}
		main()
	}

	log.Println("Getting submission details for IDs:")
	details, err := user.SubmissionDetails(
		api.SubmissionDetailsRequest{
			SubmissionIDs:   submissionIDs,
			ShowDescription: api.Yes,
		})
	if err != nil {
		log.Printf("Error getting submission details: %v\n", err)
		return
	}
	for _, submission := range details.Submissions {
		log.Printf("Submission [%s] by [%s]: https://inkbunny.net/s/%s\n", submission.Title, submission.Username, submission.SubmissionID)
	}

	var results []utils.ExtractResult
	for _, submission := range details.Submissions {
		if len(submission.Description) > 256 {
			log.Printf("Title: %s\nDescription: %s ... |>\n", submission.Title, submission.Description[:256])
		} else {
			log.Printf("Title: %s\nDescription: %s\n", submission.Title, submission.Description)
		}
		results = append(results, utils.ExtractAll(submission.Description, utils.Patterns))
	}

	log.Printf("Results: %#v\n", results)

	var out []byte
	for i, result := range results {
		var request entities.TextToImageRequest

		fieldsToSet := map[string]any{
			"steps":     &request.Steps,
			"sampler":   &request.SamplerName,
			"cfg":       &request.CFGScale,
			"seed":      &request.Seed,
			"width":     &request.Width,
			"height":    &request.Height,
			"hash":      &request.OverrideSettings.SDCheckpointHash,
			"model":     &request.OverrideSettings.SDModelCheckpoint,
			"denoising": &request.DenoisingStrength,
		}

		err := utils.ResultsToFields(result, fieldsToSet)
		if err != nil {
			log.Printf("Error setting fields: %v\n", err)
			continue
		}

		request.Prompt = utils.ExtractPositivePrompt(details.Submissions[i].Description)
		request.NegativePrompt = utils.ExtractNegativePrompt(details.Submissions[i].Description)

		system, err := llm.PrefillSystemDump(request)
		if err != nil {
			log.Printf("Error prefilling system dump: %v\n", err)
			continue
		}

		var useLLM string
		fmt.Print("Use an LLM to infer parameters? (y/[n]): ")
		fmt.Scanln(&useLLM)
		if useLLM == "y" {
			const maxRetries = 3
			host := llm.Localhost()
			log.Printf("Inferencing from %s\n", host.Endpoint.String())

			out = inferPrompt(maxRetries, &system, llm.UserMessage(details.Submissions[i].Description), host)
		} else {
			out, err = json.MarshalIndent(request, "", "  ")
			if err != nil {
				log.Printf("Error marshalling request: %v\n", err)
				continue
			}
		}
	}

	if len(out) > 0 {
		//save to json file
		err = os.WriteFile("text_to_image.json", out, 0644)
		if err != nil {
			log.Printf("Error writing text to image: %v\n", err)
			return
		}
		log.Println("Text to image saved to text_to_image.json")
	}

	var l string
	fmt.Print("Logout? (y/[n]): ")
	fmt.Scanln(&l)
	if l == "y" {
		logout(user)
	} else {
		main()
	}
}

func inferPrompt(maxRetries int, system *llm.Message, user llm.Message, host llm.Config) []byte {
	var out []byte
	if system == nil {
		system = &llm.DefaultSystem
		log.Printf("Using default system prompt: %s\n", system.Content)
	}
	for i := 0; i < maxRetries; i++ {
		log.Printf("Inferencing text to image (%d/%d)\n", i+1, maxRetries)
		accumulatedTokens := make(chan *llm.Response)
		go llm.MonitorTokens(accumulatedTokens)
		resp, err := host.Infer(&llm.Request{
			Messages: []llm.Message{
				*system,
				user,
			},
			Temperature:   1.0,
			MaxTokens:     1024,
			Stream:        true,
			StreamChannel: accumulatedTokens,
		})
		if err != nil {
			log.Printf("Error inferencing, retrying (%d/%d): %v\n", i+1, maxRetries, err)
			continue
		}

		message := utils.ExtractJson(resp.Choices[0].Message.Content)
		textToImage, err := entities.UnmarshalTextToImageRequest([]byte(message))
		if err != nil {
			log.Printf("Error unmarshalling text to image: %v, retrying (%d/%d)\n", err, i+1, maxRetries)
			continue
		}
		if textToImage.Prompt == "" {
			log.Printf("Prompt is empty, retrying (%d/%d)\n", i+1, maxRetries)
			continue
		}
		out, _ = json.MarshalIndent(textToImage, "", "  ")
		log.Printf("Successfully inferenced %d tokens with %d retries\n", resp.Usage.CompletionTokens, i)
		log.Println(strings.ReplaceAll(message, "\n", ""))
		break
	}
	return out
}
