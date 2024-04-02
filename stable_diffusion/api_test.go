package sd

import (
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
