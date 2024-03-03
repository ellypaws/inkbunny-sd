package llm

import "net/url"

var DefaultSystem = Message{
	Role: SystemRole,
	Content: "You are a backend API that responds to requests in natural language and outputs a raw JSON object. " +
		"Process the following description of an image generated with Stable Diffusion. " +
		"Output only a raw JSON response and do not include any comments. " +
		"IMPORTANT: Do not include comments, only output the JSON object" +
		"Keep loras as is `<lora:MODELNAME:weight>`" +
		"Use the following JSON format: " +
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
"comments": {
  "description": "" // Output everything in the description from the input. Escape characters for JSON
},
"denoising_strength": 0.4,
"enable_hr": false,
"hr_resize_x": 0,
"hr_resize_y": 0,
"hr_scale": 2.0, // use 2.0 if not present
"hr_second_pass_steps": 20, // use the same value as steps if not present
"hr_upscaler": "R-ESRGAN 2x+"
}`,
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
