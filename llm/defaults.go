package llm

import (
	"fmt"
	"github.com/ellypaws/inkbunny-sd/entities"
	"net/url"
	"strconv"
	"strings"
)

var DefaultSystem = Message{
	Role: SystemRole,
	Content: "You are a backend API that responds to requests in natural language and outputs a raw JSON object. \n" +
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
   "args": [] // any type
 }
}, // "script": {} OBJECTS. Include any additional information here such as CFG Rescale, Controlnet, ADetailer, RP, etc.
"cfg_scale": 7, // not to be confused rescale
"comments": {  "description": ""}, // Output everything in the description from the input. Escape characters for JSON
"denoising_strength": 0.4,
"enable_hr": false,
"hr_resize_x": 0,
"hr_resize_y": 0,
"hr_scale": 2, // use 2 if not present
"hr_second_pass_steps": 20, // use the same value as steps if not present
"hr_upscaler": "R-ESRGAN 2x+"
}`,
}

func PrefillSystem(request entities.TextToImageRequest) Message {
	var content strings.Builder
	content.WriteString("You are a backend API that responds to requests in natural language and outputs a raw JSON object. \n")
	content.WriteString("Process the following description of an image generated with Stable Diffusion. \n")
	content.WriteString("Output only a raw JSON response and do not include any comments. \n")
	content.WriteString("IMPORTANT: Do not include comments, only output the JSON object\n")
	content.WriteString("Keep loras as is `<lora:MODELNAME:weight>`\n")
	content.WriteString("Use the following JSON format: \n")
	content.WriteString("{\n")
	write(&content, `"steps": `, request.Steps, 20, ",\n")
	write(&content, `"width": `, request.Width, 512, ",\n")
	write(&content, `"height": `, request.Height, 768, ",\n")
	write(&content, `"seed": `, request.Seed, int64(1234567890), ",\n")
	write(&content, `"n_iter": `, request.NIter, 1, ", // also known as batch count\n")
	write(&content, `"batch_size": `, request.BatchSize, 1, ",\n")
	write(&content, `"prompt": `, request.Prompt, "", ", // look for positive prompt, keep loras as is, e.g. <lora:MODELNAME:float>\n")
	write(&content, `"negative_prompt": `, request.NegativePrompt, "", ", // look for negative prompt, keep loras as is, e.g. <lora:MODELNAME:float>\n")
	write(&content, `"sampler_name": `, request.SamplerName, "UniPC", ",\n")
	content.WriteString(`"override_settings": {
`)
	write(&content, `  "sd_model_checkpoint": `, request.OverrideSettings.SDModelCheckpoint, "EasyFluff", ", // also known as model\n")
	write(&content, `  "sd_checkpoint_hash": `, request.OverrideSettings.SDCheckpointHash, "f80ed3fee940", " // also known as model hash\n")
	content.WriteString("},\n")
	content.WriteString(`"alwayson_scripts": {
 "ADetailer": { // ADetailer is only an example
   "args": [] // any type
 }
}, // "script": {} OBJECTS. Include any additional information here such as CFG Rescale, Controlnet, ADetailer, RP, etc.
`)
	write(&content, `"cfg_scale": `, request.CFGScale, 7.0, ", // not to be confused rescale\n")
	content.WriteString(`"comments": {`)
	if description, ok := request.Comments["description"]; ok {
		write(&content, `  "description": `, description, "", "")
	} else {
		content.WriteString(`  "description": ""`)
	}
	content.WriteString("}, // Output everything in the description from the input. Escape characters for JSON\n")
	write(&content, `"denoising_strength": `, request.DenoisingStrength, 0.4, ",\n")
	write(&content, `"enable_hr": `, request.EnableHr, false, ",\n")
	write(&content, `"hr_resize_x": `, request.HrResizeX, 0, ",\n")
	write(&content, `"hr_resize_y": `, request.HrResizeY, 0, ",\n")
	write(&content, `"hr_scale": `, request.HrScale, float64(2), ", // use 2 if not present\n")
	write(&content, `"hr_second_pass_steps": `, request.HrSecondPassSteps, int64(20), ", // use the same value as steps if not present\n")
	write(&content, `"hr_upscaler": `, request.HrUpscaler, "R-ESRGAN 2x+", "\n")
	content.WriteString("}")

	return Message{
		Role:    SystemRole,
		Content: content.String(),
	}
}

func write(content *strings.Builder, key string, value, def any, end string) {
	content.WriteString(key)
	if valueIsSet(value) {
		content.WriteString(format(value))
	} else {
		content.WriteString(format(def))
	}
	content.WriteString(end)
}

func valueIsSet(value any) bool {
	if value == nil {
		return false
	}
	switch v := value.(type) {
	case string:
		return v != ""
	case *string:
		return v != nil && *v != ""
	case int:
		return v != 0
	case int64:
		return v != 0
	case float64:
		return v != 0.0
	case bool:
		return v
	default:
		return false
	}
}

func format(value any) string {
	if value == nil {
		return "null"
	}
	switch v := value.(type) {
	case string:
		return fmt.Sprintf(`"%s"`, v)
	case *string:
		if v == nil {
			return ""
		}
		return fmt.Sprintf(`"%s"`, *v)
	case int:
		return strconv.Itoa(v)
	case int64:
		return strconv.FormatInt(v, 10)
	case float64:
		return strconv.FormatFloat(v, 'f', -1, 64)
	case bool:
		return strconv.FormatBool(v)
	case nil:
		return "null"
	default:
		return fmt.Sprintf(`"%s"`, value)
	}
}

func UserMessage(content string) Message {
	return Message{
		Role:    UserRole,
		Content: content,
	}
}

func Localhost() Config {
	return Config{
		Host:   "localhost:7869",
		APIKey: "api-key",
		Endpoint: url.URL{
			Scheme: "http",
			Host:   "localhost:7869",
			Path:   "/v1/chat/completions",
		},
	}
}

func DefaultRequest(content string) *Request {
	return &Request{
		Messages: []Message{
			DefaultSystem,
			UserMessage(content),
		},
		Temperature: 0.7,
		MaxTokens:   2048,
		Stream:      false,
	}
}
