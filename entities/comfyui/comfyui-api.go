package comfyui

import (
	"bytes"
	"cmp"
	"encoding/json"
	"fmt"
	"regexp"
	"slices"
	"strconv"
	"strings"

	"github.com/ellypaws/inkbunny-sd/entities"
)

func UnmarshalIsolatedComfyApi(data []byte) (Api, error) {
	var container map[string]any
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.UseNumber()
	err := decoder.Decode(&container)
	if err != nil {
		return nil, err
	}

	var (
		traversable = make(Api)
		nodeErrors  NodeErrors
	)
	for k, v := range container {
		node, err := assertMarshal[ApiNode](v, true)
		if err != nil {
			nodeErrors = append(nodeErrors, err)
			continue
		}
		traversable[k] = node
	}

	return traversable, nodeErrors
}

func UnmarshalComfyApi(data []byte) (Api, error) {
	var a Api
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.UseNumber()
	err := decoder.Decode(&a)
	return a, err
}

type Api map[string]ApiNode

type ApiNode struct {
	Inputs    map[string]any `json:"inputs"`
	ClassType NodeType       `json:"class_type"`
	Meta      struct {
		Title string `json:"title"`
	} `json:"_meta"`
}

type ApiConverted struct {
	ID        int            `json:"id"`
	Inputs    map[string]any `json:"inputs"`
	ClassType NodeType       `json:"class_type"`
	Meta      struct {
		Title string `json:"title"`
	} `json:"_meta"`
}

var notDigit = regexp.MustCompile(`\D`)

func (a *Api) Convert() *entities.TextToImageRequest {
	if a == nil {
		return nil
	}
	var (
		request entities.TextToImageRequest
		prompt  PromptWriter
		loras   = make(map[string]float64)
	)
	for _, node := range *a {
		if node.ClassType == "normal" {
			node.ClassType = stringAs(NodeType(node.Meta.Title), removeEmojis, strings.TrimSpace)
		}
		switch node.ClassType {
		case CheckpointLoaderSimple, LoadCheckpoint:
			for k, v := range node.Inputs {
				if k == "ckpt_name" {
					Assert(v, SetFieldPointerOnce(&request.OverrideSettings.SDModelCheckpoint))
				}
			}
		case EmptyLatentImage:
			for k, v := range node.Inputs {
				switch k {
				case "width":
					AssertNumber(v, SetField(&request.Width))
				case "height":
					AssertNumber(v, SetField(&request.Height))
				default:
					continue
				}
			}
		case CLIPTextEncode, CLIPTextEncodeSDXL, smZCLIPTextEncode:
			for k, v := range node.Inputs {
				switch {
				case strings.HasPrefix(k, "text"):
					AssertGetter(*a, v, GetTexts, Writer(&prompt))
				case k == "target_width":
					AssertNumber(v, SetField(&request.Width))
				case k == "target_height":
					AssertNumber(v, SetField(&request.Height))
				default:
					continue
				}
			}
		case VAELoader:
			for k, v := range node.Inputs {
				if k == "vae_name" {
					Assert(v, SetFieldPointerOnce(&request.OverrideSettings.SDVae))
				}
			}
		case "ttN text":
			for k, v := range node.Inputs {
				if k == "text" {
					Assert(v, Writer(&prompt))
				}
			}
		case "ttN concat":
			for k, v := range node.Inputs {
				if strings.HasPrefix(k, "text") {
					Assert(v, Writer(&prompt))
				}
			}
		case ShowTextPys:
			for k, v := range node.Inputs {
				if strings.HasPrefix(k, "text") {
					Assert(v, Writer(&prompt))
				}
			}
		case CLIPSetLastLayer:
			for k, v := range node.Inputs {
				if k == "stop_at_clip_layer" {
					AssertNumber(v, SetField(&request.OverrideSettings.CLIPStopAtLastLayers))
				}
			}
		case KSamplerEfficient, KSampler, Digital2KSampler:
			for k, v := range node.Inputs {
				switch k {
				case "seed":
					AssertGetterNumber(*a, v, GetSeed[int64], SetField(&request.Seed))
				case "steps":
					AssertNumber(v, SetField(&request.Steps))
				case "cfg":
					AssertNumber(v, SetField(&request.CFGScale))
				case "sampler_name":
					Assert(v, SetField(&request.SamplerName))
				case "scheduler":
					Assert(v, SetFieldPointer(&request.Scheduler))
				case "denoise":
					AssertNumber(v, SetField(&request.DenoisingStrength))
				}
			}
		case SamplerCustomAdvanced:
			for k, v := range node.Inputs {
				switch k {
				case "noise":
					AssertGetterNumber(*a, v, GetSeed[int64], SetField(&request.Seed))
				case "steps":
					AssertNumber(v, SetField(&request.Steps))
				case "cfg":
					AssertNumber(v, SetField(&request.CFGScale))
				case "sampler_name":
					Assert(v, SetField(&request.SamplerName))
				case "scheduler":
					Assert(v, SetFieldPointer(&request.Scheduler))
				case "denoise":
					AssertNumber(v, SetField(&request.DenoisingStrength))
				}
			}
		case CRModelMergeStack:
			for k, v := range node.Inputs {
				if request.OverrideSettings.SDModelCheckpoint != nil {
					continue
				}
				if strings.HasPrefix(k, "ckpt_name") {
					Assert(v, SetFieldPointer(&request.OverrideSettings.SDModelCheckpoint))
				}
			}
		case CRLoRAStack:
			for _, v := range AsLoraStack(node.Inputs) {
				loras[v.LoraName] = v.ModelWeight
			}
		case RandomNoise:
			for k, v := range node.Inputs {
				if k == "noise_seed" {
					AssertNumber(v, SetField(&request.Seed))
				}
			}
		case SeedNode:
			for k, v := range node.Inputs {
				if k == "seed" {
					AssertNumber(v, SetField(&request.Seed))
				}
			}
		default:
			continue
		}
	}

	for lora, weight := range loras {
		prompt.WriteString(fmt.Sprintf("<lora:%s:%.2f>", lora, weight))
	}
	request.Prompt = prompt.String()

	return &request
}

func stringAs[T ~string](v T, f ...func(string) string) T {
	for _, f := range f {
		v = T(f(string(v)))
	}
	return v
}

var notAscii = regexp.MustCompile(`[^\x00-\x7F]+`)

// Function to remove emojis from a string
func removeEmojis(input string) string {
	return notAscii.ReplaceAllString(input, "")
}

type LoraStack struct {
	Switch      bool
	LoraName    string
	ModelWeight float64
	ClipWeight  float64
}

var lastDigit = regexp.MustCompile(`\d+$`)

func AsLoraStack(inputs map[string]any) map[string]*LoraStack {
	if inputs == nil {
		return nil
	}

	var loras = make(map[string]*LoraStack)
	for k, v := range inputs {
		num := lastDigit.FindString(k)
		if num == "" {
			continue
		}
		if _, ok := loras[num]; !ok {
			loras[num] = new(LoraStack)
		}
		switch {
		case strings.HasPrefix(k, "switch"):
			Assert(v, func(s string) { loras[num].Switch = s == "On" })
		case strings.HasPrefix(k, "lora_name"):
			Assert(v, SetField(&loras[num].LoraName))
		case strings.HasPrefix(k, "model_weight"):
			AssertNumber(v, SetField(&loras[num].ModelWeight))
		case strings.HasPrefix(k, "clip_weight"):
			AssertNumber(v, SetField(&loras[num].ClipWeight))
		}
	}

	for _, lora := range loras {
		switch {
		case lora.LoraName == "None":
			delete(loras, lora.LoraName)
		case !lora.Switch:
			delete(loras, lora.LoraName)
		case lora.ModelWeight == 0 && lora.ClipWeight == 0:
			delete(loras, lora.LoraName)
		}
	}

	return loras
}

type Number interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~uint | ~uint8 | ~uint16 | ~uint32 | ~uint64 | ~uintptr | ~float32 | ~float64
}

func AssertNumber[T Number](val any, setter func(T)) {
	if v, ok := val.(float64); ok {
		setter(T(v))
		return
	}
	if v, ok := val.(json.Number); ok {
		if i, err := v.Int64(); err == nil {
			setter(T(i))
			return
		}
		if f, err := v.Float64(); err == nil {
			setter(T(f))
			return
		}
	}
}

// AssertGetterNumber checks if the value is a link and switches context to that node
// It uses getter to get the value from the new node, then uses setter to set the value
func AssertGetterNumber[T Number](nodes Api, val any, getter func(ApiNode) (T, bool), setter func(T)) {
	if id, ok := isLink(val); ok {
		if node, ok := nodes[id]; ok {
			if v, ok := getter(node); ok {
				setter(v)
				return
			}
		}
	}
	AssertNumber(val, setter)
}

type SignedNumber interface {
	~int | ~int8 | ~int16 | ~int32 | ~int64 | ~float32 | ~float64
}

// GetSeed is a getter that checks if the node is a seed node and returns the seed value
func GetSeed[T SignedNumber](node ApiNode) (T, bool) {
	switch node.ClassType {
	case RandomNoise, SeedNode:
	default:
		return -1, false
	}
	var i T
	for k, v := range node.Inputs {
		switch k {
		case "seed", "noise_seed":
			AssertNumber(v, SetField(&i))
		}
	}
	return i, true
}

type Settable interface {
	cmp.Ordered | ~bool
}

type StringBool interface {
	~string | ~bool
}

func Assert[T StringBool](val any, setter func(T)) {
	if v, ok := val.(T); ok {
		setter(v)
	}
}

// AssertLinked checks if the value is a link and switches context to that node
// Then it uses setter to set the value from the new node.
// If you need to transform the new node's value, use AssertGetter instead
//
// Deprecated: because most linked nodes are different types, use AssertGetter instead
func AssertLinked[T StringBool](nodes Api, val any, access string, setter func(T)) {
	if id, ok := isLink(val); ok {
		if v, ok := Access[T](nodes, id, access); ok {
			setter(v)
			return
		}
	}
	Assert(val, setter)
}

// AssertGetter checks if the value is a link and switches context to that node
// It uses getter to get the value from the new node, then uses setter to set the value
func AssertGetter[T StringBool](nodes Api, val any, getter func(ApiNode) (T, bool), setter func(T)) {
	if id, ok := isLink(val); ok {
		if node, ok := nodes[id]; ok {
			if v, ok := getter(node); ok {
				setter(v)
				return
			}
		}
	}
	Assert(val, setter)
}

// GetTexts is a getter that checks if the node is a text node and returns the text value
func GetTexts(node ApiNode) (string, bool) {
	switch node.ClassType {
	case String, TextString:
	default:
		return "", false
	}
	var prompt PromptWriter
	for k, v := range node.Inputs {
		switch {
		case strings.HasPrefix(k, "text"):
			Assert(v, Writer(&prompt))
		case k == "inStr":
			Assert(v, Writer(&prompt))
		}
	}
	return prompt.String(), prompt.Len() > 0
}

// isLink checks if the value is a link to another node
// We know that a node is a link if it is an array with two elements
// e.g. ["id", 0]
func isLink(val any) (string, bool) {
	vals, ok := val.([]any)
	if !ok {
		return "", false
	}
	if len(vals) < 2 {
		return "", false
	}
	if _, ok := vals[0].(string); !ok {
		return "", false
	}
	if _, ok := vals[1].(float64); !ok {
		if _, ok := vals[1].(json.Number); !ok {
			return "", false
		}
	}
	return vals[0].(string), true
}

// Access retrieves a linked node and input name from an Api.
// Use this if you only have one node type and know which input you are looking for.
// If you need to differentiate between nodes, use AssertGetter instead.
func Access[T Settable](inputs Api, id string, inputName string) (T, bool) {
	var zero T
	if inputs == nil {
		return zero, false
	}

	node, ok := inputs[id]
	if !ok {
		return zero, false
	}

	input, ok := node.Inputs[inputName]
	if !ok {
		return zero, false
	}

	t, ok := input.(T)
	return t, ok
}

func SetField[T Settable](field *T) func(T) {
	return func(v T) {
		*field = v
	}
}

func SetFieldPointer[T Settable](field **T) func(T) {
	return func(v T) {
		*field = &v
	}
}

func SetFieldPointerOnce[T Settable](field **T) func(T) {
	return func(v T) {
		if *field != nil {
			return
		}
		*field = &v
	}
}

func Writer(b interface{ WriteString(string) }) func(string) {
	return func(s string) {
		b.WriteString(s)
	}
}

func (a *Api) ConvertSlice() []ApiConverted {
	if a == nil {
		return nil
	}
	var converted []ApiConverted
	for id, v := range *a {
		id = notDigit.ReplaceAllString(id, "")
		if id == "" {
			continue
		}
		n, err := strconv.Atoi(id)
		if err != nil {
			continue
		}
		converted = append(converted, ApiConverted{
			ID:        n,
			Inputs:    v.Inputs,
			ClassType: v.ClassType,
			Meta:      v.Meta,
		})
	}
	slices.SortFunc(converted, func(a, b ApiConverted) int {
		return cmp.Compare(a.ID, b.ID)
	})
	return converted
}

func ConvertApis(apis []Api) []*entities.TextToImageRequest {
	var converted []*entities.TextToImageRequest
	for _, a := range apis {
		converted = append(converted, a.Convert())
	}
	return converted
}
