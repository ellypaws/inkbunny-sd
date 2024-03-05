package llm

import (
	"github.com/ellypaws/inkbunny-sd/entities"
	"github.com/ellypaws/inkbunny-sd/utils"
	"testing"
)

const expected = "You are a backend API that responds to requests in natural language and outputs a raw JSON object. \n" +
	"Process the following description of an image generated with Stable Diffusion. \n" +
	"Output only a raw JSON response and do not include any comments. \n" +
	"IMPORTANT: Do not include comments, only output the JSON object\n" +
	"Keep loras as is `<lora:MODELNAME:weight>`\n" +
	"Use the following JSON format: \n" +
	`{
"steps": 20,
"width": 512,
"height": 768,
"seed": 1234567890,
"n_iter": 1, // also known as batch count
"batch_size": 1,
"prompt": "", // look for positive prompt, keep loras as is, e.g. <lora:MODELNAME:float>
"negative_prompt": "", // look for negative prompt, keep loras as is, e.g. <lora:MODELNAME:float>
"sampler_name": "UniPC",
"override_settings": {
  "sd_model_checkpoint": "EasyFluff", // also known as model
  "sd_checkpoint_hash": "f80ed3fee940" // also known as model hash
},
"alwayson_scripts": {
 "ADetailer": { // ADetailer is only an example
   "args": [] // contains an "args" array with any type inside
 }
}, // "script": OBJECTS. Include any additional information here such as CFG Rescale, Controlnet, ADetailer, RP, etc.
"cfg_scale": 7, // not to be confused rescale
"comments": {  "description": "<|description|>"}, // Output everything in the description from the input. Escape characters for JSON
"denoising_strength": 0.4,
"enable_hr": false,
"hr_resize_x": 0,
"hr_resize_y": 0,
"hr_scale": 2, // use 2 if not present
"hr_second_pass_steps": 20, // use the same value as steps if not present
"hr_upscaler": "R-ESRGAN 2x+"
}`

func TestPrefillSystemDump(t *testing.T) {
	message, err := PrefillSystemDump(entities.TextToImageRequest{})
	if err != nil {
		t.Errorf("Expected no error, got %s", err)
	}

	if message.Content != expected {
		t.Errorf("<| Expected |>\n\n```\n%s\n```\n\n<| Got |>\n\n```\n%s\n```", expected, message.Content)
	}

}

func TestPrefillSystemDumpWithDefault(t *testing.T) {
	defaultRequest := DefaultSystem

	if defaultRequest.Content != expected {
		t.Errorf("<| Expected |>\n\n```\n%s\n```\n\n<| Got |>\n\n```\n%s\n```", expected, defaultRequest.Content)
	}

	message, err := PrefillSystemDump(entities.TextToImageRequest{})
	if err != nil {
		t.Errorf("Expected no error, got %s", err)
	}

	if message.Content != expected {
		t.Errorf("<| Expected |>\n\n```\n%s\n```\n\n<| Got |>\n\n```\n%s\n```", expected, message.Content)
	}

	if message.Content != defaultRequest.Content {
		t.Errorf("Expected the same content, got different")
	}
}

func TestPrefillSystem(t *testing.T) {
	message := PrefillSystem(entities.TextToImageRequest{})

	if message.Content != expected {
		t.Errorf("<| Expected |>\n\n```\n%s\n```\n\n<| Got |>\n\n```\n%s\n```", expected, message.Content)
	}

	message.Content = utils.ExtactJson(message.Content)
	expectedCopy := utils.ExtactJson(expected)
	if message.Content != expectedCopy {
		t.Errorf("<| Expected |>\n\n```json\n%s\n```\n\n<| Got |>\n\n```json\n%s\n```", expectedCopy, message.Content)
	}
}

func TestDefaultRequest(t *testing.T) {
	request := DefaultRequest("")

	if request.Messages[0].Content != expected {
		t.Errorf("<| Expected |>\n\n```\n%s\n```\n\n<| Got |>\n\n```\n%s\n```", expected, request.Messages[0].Content)
	}
}

func TestDefaultSystem(t *testing.T) {
	message := DefaultSystem

	if message.Content != expected {
		t.Errorf("<| Expected |>\n\n```\n%s\n```\n\n<| Got |>\n\n```\n%s\n```", expected, message.Content)
	}
}
