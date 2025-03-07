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

func UnmarshalComfyApi(data []byte) (Api, error) {
	var a Api
	decoder := json.NewDecoder(bytes.NewReader(data))
	decoder.UseNumber()
	err := decoder.Decode(&a)
	return a, err
}

type Api map[string]struct {
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
		switch node.ClassType {
		case CLIPTextEncodeSDXL:
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
		case VAELoader:
			for k, v := range node.Inputs {
				if k == "vae_name" {
					Assert(v, SetFieldPointer(&request.OverrideSettings.SDVae))
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
		case Digital2KSampler:
			for k, v := range node.Inputs {
				switch k {
				case "seed":
					AssertNumber(v, SetField(&request.Seed))
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
