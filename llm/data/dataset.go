// This script is used to generate the dataset files for LLM fine-tuning.
// The fine-tuning goal is to parse non-standard user description into a standard JSON format.
// Text files are used as the input, and the json files are used as the response.
// Output is written to the dataset directory.
//
// Usage:
// go run dataset.go
//
// Before running, make sure to place the text files and json files in the same directory as this script.
// The text files should contain the user descriptions, and the json files should contain the expected responses.
//
// The json files should have the same name as the txt files, but with the .json extension.
//
// Example:
// text file: 1.txt
// json file: 1.json

package main

import (
	"bytes"
	"fmt"
	"log"
	"os"
	"strings"
)

const empty = `###Instruction:
{example['instruction']}

### Input:
{example['input']}

### Response:`

func main() {
	text, json := getFiles()
	dataset := parseDataset(text, json)

	for name, data := range dataset {
		if _, err := os.Stat("dataset"); os.IsNotExist(err) {
			os.Mkdir("dataset", 0755)
		}
		f, err := os.Create(fmt.Sprintf("dataset/%s.txt", name))
		if err != nil {
			log.Fatal(err)
		}
		defer f.Close()

		_, err = f.Write(data)
		if err != nil {
			log.Fatal(err)
		}
	}
}

// parseDataset takes in a map of text and json files and returns a map of the combined data
// It uses the commonInstruction as a base and appends the input and response to it following completeSample
func parseDataset(text, json map[string][]byte) map[string][]byte {
	var dataset = make(map[string][]byte)
	for name, input := range text {
		var out bytes.Buffer
		out.WriteString(commonInstruction)
		out.WriteString("### Input:\n")
		out.WriteString("The file name is: ")
		out.WriteString(name)
		out.WriteString("\n\n")
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

func getFiles() (text, json map[string][]byte) {
	files, _ := os.ReadDir(".")
	for _, f := range files {
		switch {
		case f.IsDir():
			continue
		case strings.HasSuffix(f.Name(), ".txt"):
			b, err := os.ReadFile(f.Name())
			if err != nil {
				continue
			}

			if text == nil {
				text = make(map[string][]byte)
			}

			text[strings.TrimSuffix(f.Name(), ".txt")] = b
		case strings.HasSuffix(f.Name(), ".json"):
			b, err := os.ReadFile(f.Name())
			if err != nil {
				continue
			}

			if json == nil {
				json = make(map[string][]byte)
			}

			json[strings.TrimSuffix(f.Name(), ".json")] = b
		}
	}
	return text, json
}

const commonInstruction = `###Instruction: 
You are a backend API that responds to requests in natural language and outputs a raw JSON object.
Process the following description of an image generated with Stable Diffusion.
Output only a raw JSON response and do not include any comments.
IMPORTANT: Do not include comments, only output the JSON object.
Sometimes there's more than one prompt, so intelligently recognize this.
Keep loras as is <lora:MODELNAME:weight>
Use the following JSON format: 
{
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
}

`

const completeSample = `###Instruction: 
You are a backend API that responds to requests in natural language and outputs a raw JSON object.
Process the following description of an image generated with Stable Diffusion.
Output only a raw JSON response and do not include any comments.
IMPORTANT: Do not include comments, only output the JSON object.
Sometimes there's more than one prompt, so intelligently recognize this.
Keep loras as is <lora:MODELNAME:weight>
Use the following JSON format: 
{
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
}

### Input:
{example['input']}

### Response:
[
 {
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
]`
