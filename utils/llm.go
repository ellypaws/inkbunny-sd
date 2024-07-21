package utils

import (
	"bytes"
	"encoding/json"
	"github.com/ellypaws/inkbunny-sd/entities"
	"io"
	"os"
	"strings"
)

type NameContent map[string][]byte

// ParseDataset takes in a map of text and json files and returns a map of the combined data
// It uses the commonInstruction as a base and appends the input and response to it following completeSample
func ParseDataset(text, json NameContent) map[string][]byte {
	var dataset = make(map[string][]byte)
	for name, input := range text {
		var out bytes.Buffer
		out.WriteString(commonInstruction)
		out.WriteString("### Input:\n")
		out.WriteString(`The file name is: "`)

		// Because some artists already have standardized txt files, opt to split each file separately
		autoSnep := strings.Contains(name, "_AutoSnep_")
		druge := strings.Contains(name, "_druge_")
		aiBean := strings.Contains(name, "_AIBean_")
		artieDragon := strings.Contains(name, "_artiedragon_")
		picker52578 := strings.Contains(name, "_picker52578_")
		fairyGarden := strings.Contains(name, "_fairygarden_")
		if autoSnep || druge || aiBean || artieDragon || picker52578 || fairyGarden {
			var inputResponse map[string]InputResponse
			switch {
			case autoSnep:
				inputResponse = MapParams(AutoSnep, WithBytes(input))
			case druge:
				inputResponse = MapParams(Common, WithBytes(input), UseDruge())
			case aiBean:
				inputResponse = MapParams(Common, WithBytes(input), UseAIBean())
			case artieDragon:
				inputResponse = MapParams(Common, WithBytes(input), UseArtie())
			case picker52578:
				inputResponse = MapParams(
					Common,
					WithBytes(input),
					WithFilename("picker52578_"),
					WithKeyCondition(func(line string) bool { return strings.HasPrefix(line, "File Name") }))
			case fairyGarden:
				inputResponse = MapParams(
					Common,
					// prepend "photo 1" to the input in case it's missing
					WithBytes(bytes.Join([][]byte{[]byte("photo 1"), input}, []byte("\n"))),
					UseFairyGarden())
			}
			if inputResponse != nil {
				out := out.Bytes()
				for name, s := range inputResponse {
					var multi bytes.Buffer
					multi.Write(out)
					multi.WriteString(name)
					multi.WriteString("\"\n\n")

					if s.Input == "" {
						continue
					}

					multi.WriteString(s.Input)

					multi.WriteString("\n\n")
					multi.WriteString("### Response:\n")

					multi.Write(s.Response)
					dataset[name] = multi.Bytes()
				}
				continue
			}
		}

		out.WriteString(name)
		out.WriteString("\"\n\n")

		out.Write(input)

		out.WriteString("\n\n")
		out.WriteString("### Response:\n")
		if j, ok := json[name]; ok {
			out.Write(j)
		}
		dataset[name] = out.Bytes()
	}
	return dataset
}

func FileToRequests(file string, processor Processor, opts ...func(*Config)) (map[string]entities.TextToImageRequest, error) {
	p, err := FileToParams(file, processor, opts...)
	if err != nil {
		return nil, err
	}
	return ParseParams(p), nil
}

// FileToParams reads the file and returns the params using a Processor
func FileToParams(file string, processor Processor, opts ...func(*Config)) (Params, error) {
	f, err := FileToBytes(file)
	if err != nil {
		return nil, err
	}
	opts = append(opts, WithFilename(file))
	opts = append(opts, WithBytes(f))
	return processor(opts...)
}

// FileToBytes reads the file and returns the content as a byte slice
func FileToBytes(file string) ([]byte, error) {
	f, err := os.OpenFile(file, os.O_RDONLY, 0)
	if err != nil {
		return nil, err
	}
	return io.ReadAll(f)
}

// InputResponse is a struct that holds the input and response training data for LLMs
type InputResponse struct {
	Input    string
	Response []byte
}

// MapParams returns the split params files as a map with the corresponding json for LLM training
func MapParams(processor Processor, opts ...func(*Config)) map[string]InputResponse {
	params, err := processor(opts...)
	if err != nil {
		return nil
	}

	if params == nil {
		return nil
	}

	request := ParseParams(params)
	if request == nil {
		return nil
	}

	var out map[string]InputResponse
	for name, r := range request {
		marshal, err := json.MarshalIndent(map[string]entities.TextToImageRequest{name: r}, "", "  ")
		if err != nil {
			continue
		}
		if out == nil {
			out = make(map[string]InputResponse)
		}
		s := InputResponse{
			Response: marshal,
		}
		if chunk, ok := params[name]; ok {
			if p, ok := chunk[Parameters]; ok {
				s.Input = p
			}
		}
		if s.Input != "" {
			out[name] = s
		}
	}
	return out
}

const empty = `###Instruction:
{example['instruction']}

### Input:
{example['input']}

### Response:`

const commonInstruction = `###Instruction: 
You are a backend API that responds to requests in natural language and outputs a raw JSON object.
Process the following description of an image generated with Stable Diffusion.
Output only a raw JSON response and do not include any comments.
IMPORTANT: Do not include comments, only output the JSON object.
Sometimes there's more than one prompt, so intelligently recognize this.
Keep loras as is <lora:MODELNAME:weight>
Use the following JSON format: 
{"filename": {
"steps": <|steps|>,
"width": <|width|>,
"height": <|height|>,
"seed": <|seed|>,
"n_iter": <|n_iter|>, // also known as batch count
"batch_size": <|batch_size|>,
"prompt": <|prompt|>, // look for positive prompt, keep loras as is, e.g. <lora:MODELNAME:float>
"negative_prompt": <|negative_prompt|>, // look for negative prompt, keep loras as is, e.g. <lora:MODELNAME:float>
"sampler_name": <|sampler_name|>, // default is Euler a
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
"comments": {  "description": <|description|>  }, // Find the generator used. Default is Stable Diffusion, ComfyUI, etc.
"denoising_strength": <|denoising_strength|>,
"enable_hr": <|enable_hr|>,
"hr_scale": <|hr_scale|>, // use 2 if not present
"hr_second_pass_steps": <|hr_second_pass_steps|>, // use the same value as steps if not present
"hr_upscaler": <|hr_upscaler|> // default is Latent
}}

`

const completeSample = `###Instruction: 
You are a backend API that responds to requests in natural language and outputs a raw JSON object.
Process the following description of an image generated with Stable Diffusion.
Output only a raw JSON response and do not include any comments.
IMPORTANT: Do not include comments, only output the JSON object.
Sometimes there's more than one prompt, so intelligently recognize this.
Keep loras as is <lora:MODELNAME:weight>
Use the following JSON format: 
{"filename": {
"steps": <|steps|>,
"width": <|width|>,
"height": <|height|>,
"seed": <|seed|>,
"n_iter": <|n_iter|>, // also known as batch count
"batch_size": <|batch_size|>,
"prompt": <|prompt|>, // look for positive prompt, keep loras as is, e.g. <lora:MODELNAME:float>
"negative_prompt": <|negative_prompt|>, // look for negative prompt, keep loras as is, e.g. <lora:MODELNAME:float>
"sampler_name": <|sampler_name|>, // default is Euler a
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
"comments": {  "description": <|description|>  }, // Find the generator used. Default is Stable Diffusion, ComfyUI, etc.
"denoising_strength": <|denoising_strength|>,
"enable_hr": <|enable_hr|>,
"hr_scale": <|hr_scale|>, // use 2 if not present
"hr_second_pass_steps": <|hr_second_pass_steps|>, // use the same value as steps if not present
"hr_upscaler": <|hr_upscaler|> // default is Latent
}}

### Input:
{example['input']}

### Response:
{
"filename": {
 "steps": 20,
 "width": 512,
 "height": 512,
 "seed": 1234,
 "n_iter": 1,
 "batch_size": 1,
 "prompt": "<|prompt|>", 
 "negative_prompt": "<|negative_prompt|>", 
 "sampler_name": "<|sampler_name|>",
 "override_settings": {
   "sd_model_checkpoint": "<|sd_model_checkpoint|>", 
   "sd_checkpoint_hash": "<|sd_checkpoint_hash|>" 
 },
 "alwayson_scripts": {
  "ADetailer": { 
    "args": [] 
  }
 }, 
 "cfg_scale": 7, 
 "comments": {  "description": "<|description|>"  }, 
 "denoising_strength": 0.4,
 "enable_hr": true,
 "hr_scale": 2,
 "hr_second_pass_steps": 20, 
 "hr_upscaler": "<|hr_upscaler|>"
 }
}`
