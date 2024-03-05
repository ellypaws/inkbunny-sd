package main

import (
	"fmt"
	"github.com/ellypaws/inkbunny-sd/entities"
	"github.com/ellypaws/inkbunny-sd/llm"
	"github.com/ellypaws/inkbunny-sd/utils"
	"github.com/ellypaws/inkbunny/api"
	"os"
	"testing"
)

func TestInfer(t *testing.T) {
	infer, err := llm.Localhost().Infer(&llm.Request{
		Messages: []llm.Message{
			llm.DefaultSystem,
			llm.UserMessage("Say hello!"),
		},
		Temperature: 1.0,
		MaxTokens:   10,
		Stream:      false,
	})
	if err != nil {
		t.Errorf("Expected no error, got %s", err)
		return
	}

	if infer.Choices[0].Message.Content == "" {
		t.Errorf("Expected content to be non-empty, got empty")
	}

	t.Logf("Infer: %+v", infer)
}

func TestGetSubmissionDetails(t *testing.T) {
	details, err := submissionDetails(t)
	if err != nil {
		t.Errorf("Got an error: %v", err)
	}

	t.Logf("Submission details: %v", details)
}

func submissionDetails(t *testing.T) (api.SubmissionDetailsResponse, error) {
	user := &api.Credentials{Sid: os.Getenv("SID")}
	user, err := api.Guest().Login()
	if err != nil {
		return api.SubmissionDetailsResponse{}, fmt.Errorf("error logging in: %w", err)
	}

	err = user.ChangeRating(api.Ratings{
		General:        true,
		Nudity:         true,
		MildViolence:   true,
		Sexual:         true,
		StrongViolence: true,
	})

	if err != nil {
		return api.SubmissionDetailsResponse{}, fmt.Errorf("error changing rating: %w", err)
	}
	t.Logf("Logged in as %s, sid: %s\n", user.Username, user.Sid)

	var submissionIDs string = os.Getenv("SUBMISSION_IDS")
	if submissionIDs == "" {
		t.Fatal("env var SUBMISSION_IDS is empty")
	}

	t.Logf("Getting submission details for IDs: %s", submissionIDs)
	details, err := user.SubmissionDetails(
		api.SubmissionDetailsRequest{
			SubmissionIDs:   submissionIDs,
			ShowDescription: api.Yes,
		})
	if err != nil {
		t.Errorf("Error getting submission details: %v", err)
		return api.SubmissionDetailsResponse{}, fmt.Errorf("error getting submission details: %w", err)
	}
	return details, nil
}

func TestPositivePrompt(t *testing.T) {
	details, err := submissionDetails(t)
	if err != nil {
		t.Fatalf("Got an error: %v", err)
	}

	const maxRetries = 3
	for i := 0; i < maxRetries; i++ {
		t.Logf("Inferencing text to image (%v/3)\n", i+1)
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
			t.Errorf("Error inferencing, retrying (%v/3): %v", i+1, err)
			continue
		}

		message := utils.ExtactJson(resp.Choices[0].Message.Content)
		textToImage, err := entities.UnmarshalTextToImageRequest([]byte(message))
		if err != nil {
			t.Errorf("Error unmarshalling text to image: %v, retrying (%v/3)", err, i+1)
			continue
		}
		if textToImage.Prompt == "" {
			t.Errorf("Prompt is empty, retrying (%v/3)", i+1)
			continue
		}
	}
}
