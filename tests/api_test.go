package main

import (
	"fmt"
	"github.com/ellypaws/inkbunny-sd/entities"
	"github.com/ellypaws/inkbunny-sd/llm"
	"github.com/ellypaws/inkbunny-sd/utils"
	"github.com/ellypaws/inkbunny/api"
	apiUtils "github.com/ellypaws/inkbunny/api/utils"
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

	if len(details.Submissions) == 0 {
		t.Fatal("No submissions found")
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
		searchResponse, err := user.SearchSubmissions(api.SubmissionSearchRequest{
			SubmissionIDsOnly:  true,
			SubmissionsPerPage: 5,
			Page:               1,
			Text:               "ai_generated",
			Type:               api.SubmissionTypePicturePinup,
			OrderBy:            "views",
			Random:             true,
			Scraps:             "both",
		})
		if err != nil {
			t.Errorf("Error searching submissions: %v", err)
			return api.SubmissionDetailsResponse{}, fmt.Errorf("error searching submissions: %w", err)
		}

		if len(searchResponse.Submissions) == 0 {
			t.Fatal("No submissions found")
		}

		const maxSubmissions = 1
		for i := 0; i < min(maxSubmissions, len(searchResponse.Submissions)); i++ {
			submissionIDs += searchResponse.Submissions[i].SubmissionID
			if i != min(maxSubmissions-1, len(searchResponse.Submissions)-1) {
				submissionIDs += ","
			}
		}
	}

	if submissionIDs == "" {
		t.Fatal("No submission IDs found")
	}

	t.Logf("Getting submission details for IDs: %s", submissionIDs)
	request := api.SubmissionDetailsRequest{
		SID:             user.Sid,
		SubmissionIDs:   submissionIDs,
		ShowDescription: api.Yes,
	}
	t.Log(api.ApiUrl("submissions", apiUtils.StructToUrlValues(request)))
	details, err := user.SubmissionDetails(request)
	if err != nil {
		t.Errorf("Error getting submission details: %v", err)
		return api.SubmissionDetailsResponse{}, fmt.Errorf("error getting submission details: %w", err)
	}

	return details, nil
}

func TestPositivePrompt(t *testing.T) {
	t.Logf("Getting submission details")
	details, err := submissionDetails(t)
	if err != nil {
		t.Fatalf("Got an error: %v", err)
	}

	if len(details.Submissions) == 0 {
		t.Fatal("No submissions found")
	}

	if details.Submissions[0].Description == "" {
		t.Fatal("No description found")
	}

	t.Logf("Inferencing text to image object of https://inkbunny.net/s/%s", details.Submissions[0].SubmissionID)

	var textToImage entities.TextToImageRequest
	const maxRetries = 3
	var success bool
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
			t.Logf("Error inferencing, retrying (%v/3): %v", i+1, err)
			continue
		}

		message := utils.ExtractJson(resp.Choices[0].Message.Content)
		textToImage, err = entities.UnmarshalTextToImageRequest([]byte(message))
		if err != nil {
			t.Logf("Error unmarshalling text to image: %v, retrying (%v/3)", err, i+1)
			continue
		}
		if textToImage.Prompt == "" {
			t.Logf("Prompt is empty, retrying (%v/3)", i+1)
			continue
		}
		success = true
		break
	}

	if !success {
		t.Fatalf("Failed to infer text to image after %v retries", maxRetries)
	}

	if textToImage.Prompt == "" {
		t.Fatal("Prompt is empty")
	}

	bytes, err := textToImage.Marshal()
	if err != nil {
		t.Fatalf("Error marshalling text to image: %v", err)
	}

	t.Logf("Text to image: %s", bytes)
}
