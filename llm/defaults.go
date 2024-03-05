package llm

import (
	"fmt"
	"github.com/ellypaws/inkbunny-sd/entities"
	"net/url"
	"reflect"
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
}`,
}

const template = "You are a backend API that responds to requests in natural language and outputs a raw JSON object. \n" +
	"Process the following description of an image generated with Stable Diffusion. \n" +
	"Output only a raw JSON response and do not include any comments. \n" +
	"IMPORTANT: Do not include comments, only output the JSON object\n" +
	"Keep loras as is `<lora:MODELNAME:weight>`\n" +
	"Use the following JSON format: \n" +
	`{
"steps": <|steps|>,
"width": <|width|>,
"height": <|height|>,
"seed": <|seed|>,
"n_iter": <|n_iter|>, // also known as batch count
"batch_size": <|batch_size|>,
"prompt": <|prompt|>, // look for positive prompt, keep loras as is, e.g. <lora:MODELNAME:float>
"negative_prompt": <|negative_prompt|>, // look for negative prompt, keep loras as is, e.g. <lora:MODELNAME:float>
"sampler_name": <|sampler_name|>,
"override_settings": {
  "sd_model_checkpoint": <|sd_model_checkpoint|>, // also known as model
  "sd_checkpoint_hash": <|sd_checkpoint_hash|> // also known as model hash
},
"alwayson_scripts": {
 "ADetailer": { // ADetailer is only an example
   "args": [] // contains an "args" array with any type inside
 }
}, // "script": OBJECTS. Include any additional information here such as CFG Rescale, Controlnet, ADetailer, RP, etc.
"cfg_scale": <|cfg_scale|>, // not to be confused rescale
"comments": {  "description": <|description|>}, // Output everything in the description from the input. Escape characters for JSON
"denoising_strength": <|denoising_strength|>,
"enable_hr": <|enable_hr|>,
"hr_resize_x": <|hr_resize_x|>,
"hr_resize_y": <|hr_resize_y|>,
"hr_scale": <|hr_scale|>, // use 2 if not present
"hr_second_pass_steps": <|hr_second_pass_steps|>, // use the same value as steps if not present
"hr_upscaler": <|hr_upscaler|>
}`

func PrefillSystemDump(request entities.TextToImageRequest) (Message, error) {
	easyFluff := "EasyFluff"
	r := entities.TextToImageRequest{
		Steps:          20,
		Width:          512,
		Height:         768,
		Seed:           1234567890,
		NIter:          1,
		BatchSize:      1,
		Prompt:         "",
		NegativePrompt: "",
		SamplerName:    "UniPC",
		OverrideSettings: entities.Config{
			SDModelCheckpoint: &easyFluff,
			SDCheckpointHash:  "f80ed3fee940",
		},
		CFGScale: 7,
		Comments: map[string]string{
			"description": "<|description|>",
		},
		DenoisingStrength: 0.4,
		EnableHr:          false,
		HrResizeX:         0,
		HrResizeY:         0,
		HrScale:           2,
		HrSecondPassSteps: 20,
		HrUpscaler:        "R-ESRGAN 2x+",
	}

	data, err := r.Marshal()
	if err != nil {
		return Message{}, fmt.Errorf("failed to marshal request: %w", err)
	}

	request, err = entities.UnmarshalTextToImageRequest(data)
	if err != nil {
		return Message{}, fmt.Errorf("failed to unmarshal request: %w", err)
	}

	out := template

	// use reflect to get json tags
	v := reflect.ValueOf(request)
	out, err = replaceNestedFields(out, v)
	if err != nil {
		return Message{}, err
	}

	return Message{
		Role:    SystemRole,
		Content: out,
	}, nil
}

func replaceNestedFields(out string, v reflect.Value) (string, error) {
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldName := v.Type().Field(i).Name
		tag := v.Type().Field(i).Tag.Get("json")
		if tag == "" {
			continue // Skip if there's no json tag
		}
		tag = strings.Split(tag, ",")[0]

		// Special handling for the Comments field to insert the description
		if fieldName == "Comments" {
			// Assuming "description" is the key for the description in the map
			description, exists := field.Interface().(map[string]string)["description"]
			if exists {
				// Replace the placeholder for the description in the template
				out = strings.Replace(out, "<|description|>", format(description), -1)
			}
			continue
		}

		// Check if the field is a nested struct
		if field.Kind() == reflect.Struct {
			var err error
			out, err = replaceNestedFields(out, field)
			if err != nil {
				return "", err
			}
		} else {
			// Replace the placeholder with the formatted value
			out = strings.Replace(out, fmt.Sprintf("<|%s|>", tag), format(field.Interface()), -1)
		}
	}
	return out, nil
}

// Deprecated: Use PrefillSystemDump instead.
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
