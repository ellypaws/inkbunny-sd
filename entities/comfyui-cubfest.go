package entities

import (
	"encoding/json"
	"fmt"
)

type ComfyUICubFestAITime map[string]ComfyUICubFestAI

func UnmarshalComfyUICubFestAITime(data []byte) (ComfyUICubFestAITime, error) {
	var r ComfyUICubFestAITime
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *ComfyUICubFestAI) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type ComfyUICubFestAI struct {
	FilenamePrefix    string            `json:"filename_prefix"`
	Resolution        string            `json:"resolution"`
	Checkpoint        string            `json:"checkpoint"`
	Loras             string            `json:"loras"`
	Vae               string            `json:"vae"`
	UpscaleModel      string            `json:"upscale_model"`
	SamplerParameters SamplerParameters `json:"sampler_parameters"`
	PositivePrompt    string            `json:"positive_prompt"`
}

type SamplerParameters struct {
	CkptName    string  `json:"ckpt_name"`
	VaeName     string  `json:"vae_name"`
	ModelName   string  `json:"model_name"`
	Loras       string  `json:"loras"`
	Seed        int64   `json:"seed"`
	Steps       int64   `json:"steps"`
	CFG         float64 `json:"cfg"`
	SamplerName string  `json:"sampler_name"`
	Scheduler   string  `json:"scheduler"`
	Denoise     float64 `json:"denoise"`
}

func (r *ComfyUICubFestAI) Convert() TextToImageRequest {
	var width, height int
	if r.Resolution != "" {
		_, _ = fmt.Sscanf(r.Resolution, "%dx%d", &width, &height)
	}

	var checkpoint, vae *string
	if r.Checkpoint != "" {
		checkpoint = &r.Checkpoint
	}
	if r.SamplerParameters.CkptName != "" {
		checkpoint = &r.SamplerParameters.CkptName
	}
	if r.Vae != "" {
		vae = &r.Vae
	}
	if r.SamplerParameters.VaeName != "" {
		vae = &r.SamplerParameters.VaeName
	}

	return TextToImageRequest{
		Prompt:            r.PositivePrompt,
		Width:             width,
		Height:            height,
		SamplerName:       r.SamplerParameters.SamplerName,
		Seed:              r.SamplerParameters.Seed,
		Steps:             int(r.SamplerParameters.Steps),
		CFGScale:          r.SamplerParameters.CFG,
		Scheduler:         &r.SamplerParameters.Scheduler,
		DenoisingStrength: r.SamplerParameters.Denoise,
		HrUpscaler:        r.UpscaleModel,
		OverrideSettings: Config{
			SDModelCheckpoint: checkpoint,
			SDVae:             vae,
		},
	}
}
