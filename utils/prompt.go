package utils

import (
	"errors"
	"fmt"
	"github.com/ellypaws/inkbunny-sd/entities"
	"math"
	"strconv"
	"strings"
)

func ExtractPositivePrompt(s string) string {
	s = RemoveBBCode(s)
	result := Extract(s, positivePattern)

	if result == "" {
		result = ExtractPositiveBackwards(s)
	}

	if result == "" {
		result = ExtractPositiveForward(s)
	}

	return trim(result)
}

func ExtractNegativePrompt(s string) string {
	s = RemoveBBCode(s)
	result := Extract(s, negativePattern)

	if result == "" {
		result = ExtractNegativeForward(s)
	}

	if result == "" {
		result = ExtractNegativeBackwards(s)
	}

	return trim(result)
}

func trim(s string) string {
	return strings.Trim(s, " \n|[]")
}

func DescriptionHeuristics(description string) (entities.TextToImageRequest, error) {
	results := ExtractAll(description, Patterns)

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

	err := ResultsToFields(results, fieldsToSet)
	if err != nil {
		return request, err
	}

	request.Prompt = ExtractPositivePrompt(description)
	request.NegativePrompt = ExtractNegativePrompt(description)
	return request, nil
}

// IncompleteParameters is returned when the parameters potentially are not enough to create a request.
var IncompleteParameters = errors.New("incomplete parameters")

// ParameterHeuristics returns a TextToImageRequest from the given parameters.
// It uses the standard parameters embedded in an image from Stable Diffusion.
// The function emulates the behavior of parse_generation_parameters in
// https://github.com/AUTOMATIC1111/stable-diffusion-webui/blob/master/modules/infotext_utils.py#L233
func ParameterHeuristics(parameters string) (entities.TextToImageRequest, error) {
	var request entities.TextToImageRequest

	lines := strings.Split(parameters, "\n")

	if len(lines) < 2 {
		return entities.TextToImageRequest{
			Prompt: parameters,
		}, IncompleteParameters
	}

	positive, negative := getPrompts(lines)

	results := ExtractKeys(lines[len(lines)-1])

	if sizes, ok := results["Size"]; ok {
		for i, size := range strings.Split(sizes, "x") {
			switch i {
			case 0:
				request.Width, _ = strconv.Atoi(size)
			case 1:
				request.Height, _ = strconv.Atoi(size)
			}
		}
	}

	if _, ok := results["Clip skip"]; ok {
		results["Clip skip"] = "1"
	}

	if hypernet, ok := results["Hypernet"]; ok {
		positive.WriteString(fmt.Sprintf("<hypernet:%s:%s>", hypernet, results["Hypernet strength"]))
	}

	if _, ok := results["Hires resize-1"]; !ok {
		results["Hires resize-1"] = "0"
		results["Hires resize-2"] = "0"
	}

	if _, ok := results["Hires sampler"]; !ok {
		results["Hires sampler"] = "Use same sampler"
	}

	if _, ok := results["Hires checkpoint"]; !ok {
		results["Hires checkpoint"] = "Use same checkpoint"
	}

	if _, ok := results["Hires prompt"]; !ok {
		results["Hires prompt"] = ""
	}

	if _, ok := results["Hires negative prompt"]; !ok {
		results["Hires negative prompt"] = ""
	}

	if _, ok := results["Mask mode"]; !ok {
		results["Mask mode"] = "Inpaint masked"
	}

	if _, ok := results["Masked content"]; !ok {
		results["Masked content"] = "original"
	}

	if _, ok := results["Inpaint area"]; !ok {
		results["Inpaint area"] = "Whole picture"
	}

	if _, ok := results["Masked area padding"]; !ok {
		results["Masked area padding"] = "32"
	}

	restoreOldHiresFixParams(results, false)

	if _, ok := results["RNG"]; !ok {
		results["RNG"] = "GPU"
	}

	if _, ok := results["Schedule type"]; !ok {
		results["Schedule type"] = "Automatic"
	}

	if _, ok := results["Schedule max sigma"]; !ok {
		results["Schedule max sigma"] = "0"
	}

	if _, ok := results["Schedule min sigma"]; !ok {
		results["Schedule min sigma"] = "0"
	}

	if _, ok := results["Schedule rho"]; !ok {
		results["Schedule rho"] = "0"
	}

	if _, ok := results["VAE Encoder"]; !ok {
		results["VAE Encoder"] = "Full"
	}

	if _, ok := results["VAE Decoder"]; !ok {
		results["VAE Decoder"] = "Full"
	}

	if _, ok := results["FP8 weight"]; !ok {
		results["FP8 weight"] = "Disable"
	}

	if _, ok := results["Cache FP16 weight for LoRA"]; !ok && results["FP8 weight"] != "Disable" {
		results["Cache FP16 weight for LoRA"] = "False"
	}

	//promptAttention := parsePromptAttention(request.Prompt) + parsePromptAttention(request.NegativePrompt)

	//var promptUsesEmphasis [][]string
	//for _, p := range promptAttention {
	//	if p[1] == 1.0 || p[0] == "BREAK" {
	//		promptUsesEmphasis = append(promptUsesEmphasis, p)
	//	}
	//}

	//if _, ok := results["Emphasis"]; !ok && promptUsesEmphasis {
	//	results["Emphasis"] = "Original"
	//}

	if _, ok := results["Refiner switch by sampling steps"]; !ok {
		results["Refiner switch by sampling steps"] = "False"
	}

	err := ResultsToFields(results, TextToImageFields(&request))
	if err != nil {
		return request, err
	}

	request.Prompt = positive.String()
	request.NegativePrompt = negative.String()

	// Fallback
	if request.Prompt == "" {
		request.Prompt = ExtractPositivePrompt(parameters)
	}

	if request.NegativePrompt == "" {
		request.NegativePrompt = ExtractNegativePrompt(parameters)
	}

	return request, nil
}

// getPrompts returns the positive and negative prompts from the given lines.
// It goes line by line until it finds the negative prompt, then it returns the positive and negative prompts.
// Everything before the negative prompt is considered the positive prompt.
// The last line is not included since it's the extra parameters.
func getPrompts(lines []string) (strings.Builder, strings.Builder) {
	var positive, negative strings.Builder
	var negativeFound bool
	for _, line := range lines[:len(lines)-1] {
		line = strings.TrimSpace(line)
		switch {
		case strings.HasPrefix(line, "Negative prompt:"):
			negativeFound = true
			negative.WriteString(strings.TrimSpace(line[16:]))
		case negativeFound:
			if negative.Len() > 0 {
				negative.WriteString("\n")
			}
			negative.WriteString(line)
		default:
			if positive.Len() > 0 {
				positive.WriteString("\n")
			}
			positive.WriteString(line)
		}
	}
	return positive, negative
}

func TextToImageFields(request *entities.TextToImageRequest) map[string]any {
	if request == nil {
		return nil
	}
	return map[string]any{
		"Steps":                 &request.Steps,
		"Sampler":               &request.SamplerName,
		"CFG scale":             &request.CFGScale,
		"Seed":                  &request.Seed,
		"Denoising strength":    &request.DenoisingStrength,
		"Model":                 &request.OverrideSettings.SDModelCheckpoint,
		"Model hash":            &request.OverrideSettings.SDCheckpointHash,
		"VAE":                   &request.OverrideSettings.SDVae,
		"VAE hash":              &request.OverrideSettings.SDVaeExplanation,
		"Hires upscale":         &request.HrScale,
		"Hires steps":           &request.HrSecondPassSteps,
		"Hires upscaler":        &request.HrUpscaler,
		"Clip skip":             &request.OverrideSettings.CLIPStopAtLastLayers,
		"Hires resize-1":        &request.HrResizeX,
		"Hires resize-2":        &request.HrResizeY,
		"Hires sampler":         &request.HrSamplerName,
		"Hires checkpoint":      &request.HrCheckpointName,
		"Hires prompt":          &request.HrPrompt,
		"Hires negative prompt": &request.HrNegativePrompt,
		"RNG":                   &request.OverrideSettings.RandnSource,
		"Schedule type":         &request.OverrideSettings.KSchedType,
		"Schedule max sigma":    &request.OverrideSettings.SigmaMax,
		"Schedule min sigma":    &request.OverrideSettings.SigmaMin,
		"Schedule rho":          &request.OverrideSettings.Rho,
		"VAE Encoder":           &request.OverrideSettings.SDVaeEncodeMethod,
		"VAE Decoder":           &request.OverrideSettings.SDVaeDecodeMethod,
		//"FP8 weight":                       &request.OverrideSettings.DisableWeightsAutoSwap,       // TODO: this is a bool, but FP8 weight is a string e.g. "Disable"
		//"Cache FP16 weight for LoRA":       &request.OverrideSettings.SDVaeCheckpointCache,         // TODO: this is a float64, but Cache FP16 weight for LoRA is a bool e.g. False
		//"Emphasis":                         &request.OverrideSettings.EnableEmphasis,               // TODO: this is a string, but Emphasis is a bool e.g. "Original"
		//"Emphasis":                         &request.OverrideSettings.UseOldEmphasisImplementation, // TODO: this is a string, but Emphasis is a bool e.g. "Original"
		//"Refiner switch by sampling steps": &request.OverrideSettings.HiresFixRefinerPass,          // TODO: this is a string, but Refiner switch by sampling steps is a bool e.g. False
	}
}

// Deprecated: not yet implemented in callers
// restoreOldHiresFixParams restores the old hires fix parameters if the new hires fix parameters are not present.
// Set use is true if the new hires fix parameters should be used, false if the old hires fix parameters should be used.
func restoreOldHiresFixParams(results ExtractResult, use bool) {
	var firstpassWidth, firstpassHeight int

	fieldsToSet := map[string]any{
		"First pass size-1": &firstpassWidth,
		"First pass size-2": &firstpassHeight,
	}
	if err := ResultsToFields(results, fieldsToSet); err != nil {
		return
	}

	if use {
		var hiresWidth, hiresHeight int
		fieldsToSet = map[string]any{
			"Hires resize-1": &hiresWidth,
			"Hires resize-2": &hiresHeight,
		}
		if err := ResultsToFields(results, fieldsToSet); err != nil {
			return
		}

		if hiresWidth != 0 && hiresHeight != 0 {
			results["Size-1"] = strconv.Itoa(hiresWidth)
			results["Size-2"] = strconv.Itoa(hiresHeight)
			return
		}
	}

	if firstpassWidth == 0 || firstpassHeight == 0 {
		return
	}

	width, _ := strconv.Atoi(results["Size-1"])
	height, _ := strconv.Atoi(results["Size-2"])

	if firstpassWidth == 0 || firstpassHeight == 0 {
		firstpassWidth, firstpassHeight = oldHiresFixFirstPassDimensions(width, height)
	}

	results["Size-1"] = strconv.Itoa(firstpassWidth)
	results["Size-2"] = strconv.Itoa(firstpassHeight)
	results["Hires resize-1"] = strconv.Itoa(width)
	results["Hires resize-2"] = strconv.Itoa(height)
}

func oldHiresFixFirstPassDimensions(width int, height int) (int, int) {
	desiredPixelCount := 512 * 512
	actualPixelCount := width * height
	scale := math.Sqrt(float64(desiredPixelCount) / float64(actualPixelCount))
	width = int(math.Ceil(scale*float64(width/64)) * 64)
	height = int(math.Ceil(scale*float64(height/64)) * 64)
	return width, height
}
