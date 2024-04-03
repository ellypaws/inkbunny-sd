package sd

import (
	_ "embed"
	"encoding/base64"
	"fmt"
	"github.com/ellypaws/inkbunny-sd/entities"
	"os"
	"testing"
)

var h = func() *Host {
	if h := FromString(os.Getenv("URL")); h != nil {
		return h
	}
	return DefaultHost
}()

var slow = func() bool {
	if os.Getenv("SLOW") == "true" {
		return true
	}
	return false
}()

func TestHost_GetConfig(t *testing.T) {
	config, err := h.GetConfig()
	if err != nil {
		t.Errorf("Failed to get config: %v", err)
	}

	if config == nil {
		t.Errorf("Config is nil")
	}

	t.Logf("Config: %v", config)
}

func TestToImages(t *testing.T) {
	if !slow {
		t.Skip("Skipping image generation, set SLOW=true to enable")
	}
	request := &entities.TextToImageRequest{
		Prompt:      "A cat",
		Steps:       20,
		SamplerName: "DDIM",
	}
	response, err := h.TextToImageRequest(request)
	if err != nil {
		t.Errorf("Failed to get response: %v", err)
	}

	images, err := ToImages(response)
	if err != nil {
		t.Errorf("Failed to get images: %v", err)
	}

	if os.Getenv("SAVE_IMAGES") != "true" {
		return
	}
	for i, img := range images {
		if _, err := os.Stat("images"); os.IsNotExist(err) {
			os.Mkdir("images", os.ModePerm)
		}

		file, err := os.Create(fmt.Sprintf("images/%d.png", i))
		if err != nil {
			t.Errorf("Failed to create file: %v", err)
		}

		_, err = file.Write(img)
	}
}

//go:embed images/0.png
var image []byte

func TestHost_Interrogate(t *testing.T) {
	if _, err := os.Stat("images/0.png"); err != nil {
		s := slow
		slow = true
		t.Run("TestToImages", TestToImages)
		slow = s
	}
	if len(image) == 0 {
		t.Fatalf("Image is empty")
	}
	b64 := base64.StdEncoding.EncodeToString(image)
	req := (&entities.TaggerRequest{
		Image: &b64,
		Model: entities.TaggerZ3DE621Convnext,
	}).SetThreshold(0.5)
	response, err := h.Interrogate(req)
	if err != nil {
		t.Fatalf("Failed to interrogate: %v", err)
	}

	if response.Caption == nil {
		t.Fatalf("Caption is nil")
	}

	var foundFeline bool
	for tag, confidence := range response.Captions() {
		t.Logf("Tag: %v, Confidence: %.2f", tag, confidence)
		if tag == "felid" {
			foundFeline = true
		}
	}

	if !foundFeline {
		t.Errorf("Failed to find felid")
	}
}
