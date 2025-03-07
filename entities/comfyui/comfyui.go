package comfyui

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"strconv"
	"strings"

	"github.com/ellypaws/inkbunny-sd/entities"
)

func UnmarshalComfyUI(data []byte) (ComfyUI, error) {
	var r ComfyUI
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *ComfyUI) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type ComfyUI struct {
	LastNodeID int64           `json:"last_node_id"`
	LastLinkID int64           `json:"last_link_id"`
	Nodes      []Node          `json:"nodes"`
	Links      [][]LinkElement `json:"links"`
	Groups     []Group         `json:"groups"`
	Config     ConfigValue     `json:"config"`
	Extra      Extra           `json:"extra"`
	Version    float64         `json:"version"`
}

type ConfigValue struct{}

type Extra struct {
	WorkspaceInfo *WorkspaceInfo `json:"workspace_info,omitempty"`
	Ds            *Ds            `json:"ds,omitempty"`
	GroupNodes    *GroupNodes    `json:"groupNodes,omitempty"`
}

type WorkspaceInfo struct {
	ID string `json:"id"`
}

type Ds struct {
	Scale  float64 `json:"scale"`
	Offset *Pos    `json:"offset"`
}

type GroupNodes struct {
	Bus Bus `json:"Bus"`
}

type Bus struct {
	Nodes    []BusNode              `json:"nodes"`
	Links    [][]LinkElement        `json:"links"`
	External []interface{}          `json:"external"`
	Config   map[string]ConfigValue `json:"config"`
}

type BusNode struct {
	Type       string      `json:"type"`
	Pos        []int64     `json:"pos"`
	Size       *Size       `json:"size"`
	Flags      ConfigValue `json:"flags"`
	Order      int64       `json:"order"`
	Mode       int64       `json:"mode"`
	Inputs     []Input     `json:"inputs"`
	Outputs    []Output    `json:"outputs"`
	Properties Properties  `json:"properties"`
	Index      int64       `json:"index"`
}

func UnmarshalComfyUIBasic(data []byte) (Basic, error) {
	var r Basic
	err := json.Unmarshal(data, &r)
	return r, err
}

func (r *Basic) Marshal() ([]byte, error) {
	return json.Marshal(r)
}

type Basic struct {
	Nodes   []Node  `json:"nodes"`
	Version float64 `json:"version"`
}

type Group struct {
	Title    string  `json:"title"`
	Bounding []int64 `json:"bounding"`
	Color    string  `json:"color"`
	FontSize int64   `json:"font_size"`
	Locked   bool    `json:"locked"`
}

type Node struct {
	ID            int64               `json:"id"`
	Type          NodeType            `json:"type"`
	Pos           *Pos                `json:"pos"`
	Size          *Pos                `json:"size"`
	Flags         Flags               `json:"flags"`
	Order         int64               `json:"order"`
	Mode          Mode                `json:"mode"`
	Inputs        []Input             `json:"inputs,omitempty"`
	Outputs       []Output            `json:"outputs,omitempty"`
	Properties    Properties          `json:"properties"`
	WidgetsValues *WidgetsValuesUnion `json:"widgets_values,omitempty"`
	Color         *string             `json:"color,omitempty"`
	BGColor       *string             `json:"bgcolor,omitempty"`
	Title         *Title              `json:"title,omitempty"`
	Shape         *int64              `json:"shape,omitempty"`
}

type Mode int64

const (
	ModeNormal Mode = iota << 1
	ModeMuted
	ModeBypass
)

type Flags struct {
	Collapsed *bool `json:"collapsed,omitempty"`
}

type Input struct {
	Name      string   `json:"name"`
	Type      LinkEnum `json:"type"`
	Link      *int64   `json:"link"`
	SlotIndex *int64   `json:"slot_index,omitempty"`
	Widget    *Widget  `json:"widget,omitempty"`
	Dir       *int64   `json:"dir,omitempty"`
}

type Widget struct {
	Name   WidgetName      `json:"name"`
	Config []ConfigElement `json:"config,omitempty"`
}

type ConfigConfig struct {
	Multiline *bool    `json:"multiline,omitempty"`
	Default   *int64   `json:"default,omitempty"`
	Min       *int64   `json:"min,omitempty"`
	Max       *float64 `json:"max,omitempty"`
	Step      *int64   `json:"step,omitempty"`
}

type Output struct {
	Name      string   `json:"name"`
	Type      LinkEnum `json:"type"`
	Links     []int64  `json:"links"`
	SlotIndex *int64   `json:"slot_index,omitempty"`
	Shape     *int64   `json:"shape,omitempty"`
	Dir       *int64   `json:"dir,omitempty"`
	Label     *string  `json:"label,omitempty"`
	Widget    *Widget  `json:"widget,omitempty"`
}

type Properties struct {
	NodeNameForSR      *string `json:"Node name for S&R,omitempty"`
	Text               *string `json:"text,omitempty"`
	MatchColors        *string `json:"matchColors,omitempty"`
	MatchTitle         *string `json:"matchTitle,omitempty"`
	ShowNav            *bool   `json:"showNav,omitempty"`
	Sort               *string `json:"sort,omitempty"`
	CustomSortAlphabet *string `json:"customSortAlphabet,omitempty"`
	ToggleRestriction  *string `json:"toggleRestriction,omitempty"`
	ShowOutputText     *bool   `json:"showOutputText,omitempty"`
	Horizontal         *bool   `json:"horizontal,omitempty"`
	WidgetReplace      *bool   `json:"Run widget replace on values,omitempty"`
	ComparerMode       *string `json:"comparer_mode,omitempty"`
}

type WidgetsValueClass struct {
	Filename         *string  `json:"filename,omitempty"`
	Subfolder        *string  `json:"subfolder,omitempty"`
	Type             *string  `json:"type,omitempty"`
	ImageHash        *float64 `json:"image_hash,omitempty"`
	ForwardFilename  *string  `json:"forward_filename,omitempty"`
	ForwardSubfolder *string  `json:"forward_subfolder,omitempty"`
	ForwardType      *string  `json:"forward_type,omitempty"`

	// Content for LoraLoaderPys
	Content *string `json:"content,omitempty"`
	Image   any     `json:"image,omitempty"`
}

type WidgetsValuesClass struct {
	UpscaleBy         float64 `json:"upscale_by"`
	Seed              int64   `json:"seed"`
	Steps             int64   `json:"steps"`
	CFG               float64 `json:"cfg"`
	SamplerName       string  `json:"sampler_name"`
	Scheduler         string  `json:"scheduler"`
	Denoise           float64 `json:"denoise"`
	ModeType          string  `json:"mode_type"`
	TileWidth         int64   `json:"tile_width"`
	TileHeight        int64   `json:"tile_height"`
	MaskBlur          int64   `json:"mask_blur"`
	TilePadding       int64   `json:"tile_padding"`
	SeamFixMode       string  `json:"seam_fix_mode"`
	SeamFixDenoise    int64   `json:"seam_fix_denoise"`
	SeamFixWidth      int64   `json:"seam_fix_width"`
	SeamFixMaskBlur   int64   `json:"seam_fix_mask_blur"`
	SeamFixPadding    int64   `json:"seam_fix_padding"`
	ForceUniformTiles bool    `json:"force_uniform_tiles"`
	TiledDecode       bool    `json:"tiled_decode"`
}

type LinkEnum string

const (
	LinkASCII           LinkEnum = "ASCII"
	LinkBboxDetector    LinkEnum = "BBOX_DETECTOR"
	LinkBus             LinkEnum = "BUS"
	LinkClip            LinkEnum = "CLIP"
	LinkConditioning    LinkEnum = "CONDITIONING"
	LinkControlNet      LinkEnum = "CONTROL_NET"
	LinkControlNetStack LinkEnum = "CONTROL_NET_STACK"
	LinkDependencies    LinkEnum = "DEPENDENCIES"
	LinkDetailerHook    LinkEnum = "DETAILER_HOOK"
	LinkDetailerPipe    LinkEnum = "DETAILER_PIPE"
	LinkEmpty           LinkEnum = "*"
	LinkFloat           LinkEnum = "FLOAT"
	LinkGuider          LinkEnum = "GUIDER"
	LinkImage           LinkEnum = "IMAGE"
	LinkImagePath       LinkEnum = "IMAGE_PATH"
	LinkInt             LinkEnum = "INT"
	LinkLatent          LinkEnum = "LATENT"
	LinkLoraStack       LinkEnum = "LORA_STACK"
	LinkMask            LinkEnum = "MASK"
	LinkModel           LinkEnum = "MODEL"
	LinkModelStack      LinkEnum = "MODEL_STACK"
	LinkNoise           LinkEnum = "NOISE"
	LinkPipeLine        LinkEnum = "PIPE_LINE"
	LinkSamModel        LinkEnum = "SAM_MODEL"
	LinkSampler         LinkEnum = "SAMPLER"
	LinkScript          LinkEnum = "SCRIPT"
	LinkSdxlTuple       LinkEnum = "SDXL_TUPLE"
	LinkSegmDetector    LinkEnum = "SEGM_DETECTOR"
	LinkSigmas          LinkEnum = "SIGMAS"
	LinkString          LinkEnum = "STRING"
	LinkUpscaleModel    LinkEnum = "UPSCALE_MODEL"
	LinkVae             LinkEnum = "VAE"
)

type WidgetName string

const (
	WidgetSeed  WidgetName = "seed"
	WidgetText  WidgetName = "text"
	WidgetValue WidgetName = "value"
)

type Title string

const (
	Furry          Title = "Furry"
	IncludePrompt  Title = "Include Prompt"
	NegativePrompt Title = "Negative Prompt"
)

type Size struct {
	DoubleMap    map[string]float64
	IntegerArray []int64
}

func (x *Size) UnmarshalJSON(data []byte) error {
	x.IntegerArray = nil
	x.DoubleMap = nil
	object, err := unmarshalUnion(data, nil, nil, nil, nil, true, &x.IntegerArray, false, nil, true, &x.DoubleMap, false, nil, false)
	if err != nil {
		return err
	}
	if object {
	}
	return nil
}

func (x *Size) MarshalJSON() ([]byte, error) {
	return marshalUnion(nil, nil, nil, nil, x.IntegerArray != nil, x.IntegerArray, false, nil, x.DoubleMap != nil, x.DoubleMap, false, nil, false)
}

type LinkElement struct {
	Enum    *LinkEnum
	Integer *int64
}

func (x *LinkElement) UnmarshalJSON(data []byte) error {
	x.Enum = nil
	object, err := unmarshalUnion(data, &x.Integer, nil, nil, nil, false, nil, false, nil, false, nil, true, &x.Enum, false)
	if err != nil {
		return err
	}
	if object {
	}
	return nil
}

func (x *LinkElement) MarshalJSON() ([]byte, error) {
	return marshalUnion(x.Integer, nil, nil, nil, false, nil, false, nil, false, nil, x.Enum != nil, x.Enum, false)
}

type ConfigElement struct {
	ConfigConfig *ConfigConfig
	Enum         *LinkEnum
}

func (x *ConfigElement) UnmarshalJSON(data []byte) error {
	x.ConfigConfig = nil
	x.Enum = nil
	var c ConfigConfig
	object, err := unmarshalUnion(data, nil, nil, nil, nil, false, nil, true, &c, false, nil, true, &x.Enum, false)
	if err != nil {
		return err
	}
	if object {
		x.ConfigConfig = &c
	}
	return nil
}

func (x *ConfigElement) MarshalJSON() ([]byte, error) {
	return marshalUnion(nil, nil, nil, nil, false, nil, x.ConfigConfig != nil, x.ConfigConfig, false, nil, x.Enum != nil, x.Enum, false)
}

type Pos struct {
	DoubleArray []float64
	DoubleMap   map[string]float64
}

func (x *Pos) UnmarshalJSON(data []byte) error {
	x.DoubleArray = nil
	x.DoubleMap = nil
	object, err := unmarshalUnion(data, nil, nil, nil, nil, true, &x.DoubleArray, false, nil, true, &x.DoubleMap, false, nil, false)
	if err != nil {
		return err
	}
	if object {
	}
	return nil
}

func (x *Pos) MarshalJSON() ([]byte, error) {
	return marshalUnion(nil, nil, nil, nil, x.DoubleArray != nil, x.DoubleArray, false, nil, x.DoubleMap != nil, x.DoubleMap, false, nil, false)
}

type WidgetsValuesUnion struct {
	UnionArray         []WidgetsValueElement
	WidgetsValuesClass *WidgetsValuesClass
}

func (x *WidgetsValuesUnion) UnmarshalJSON(data []byte) error {
	x.UnionArray = nil
	x.WidgetsValuesClass = nil
	var c WidgetsValuesClass
	object, err := unmarshalUnion(data, nil, nil, nil, nil, true, &x.UnionArray, true, &c, false, nil, false, nil, false)
	if err != nil {
		return err
	}
	if object {
		x.WidgetsValuesClass = &c
	}
	return nil
}

func (x *WidgetsValuesUnion) MarshalJSON() ([]byte, error) {
	return marshalUnion(nil, nil, nil, nil, x.UnionArray != nil, x.UnionArray, x.WidgetsValuesClass != nil, x.WidgetsValuesClass, false, nil, false, nil, false)
}

type WidgetsValueElement struct {
	Bool              *bool
	Double            *float64
	String            *string
	StringArray       []string
	WidgetsValueClass *WidgetsValueClass
}

func (x *WidgetsValueElement) UnmarshalJSON(data []byte) error {
	x.StringArray = nil
	x.WidgetsValueClass = nil
	var c WidgetsValueClass
	object, err := unmarshalUnion(data, nil, &x.Double, &x.Bool, &x.String, true, &x.StringArray, true, &c, false, nil, false, nil, true)
	if err != nil {
		return err
	}
	if object {
		x.WidgetsValueClass = &c
	}
	return nil
}

func (x *WidgetsValueElement) MarshalJSON() ([]byte, error) {
	return marshalUnion(nil, x.Double, x.Bool, x.String, x.StringArray != nil, x.StringArray, x.WidgetsValueClass != nil, x.WidgetsValueClass, false, nil, false, nil, true)
}

func unmarshalUnion(data []byte, pi **int64, pf **float64, pb **bool, ps **string, haveArray bool, pa interface{}, haveObject bool, pc interface{}, haveMap bool, pm interface{}, haveEnum bool, pe interface{}, nullable bool) (bool, error) {
	if pi != nil {
		*pi = nil
	}
	if pf != nil {
		*pf = nil
	}
	if pb != nil {
		*pb = nil
	}
	if ps != nil {
		*ps = nil
	}

	dec := json.NewDecoder(bytes.NewReader(data))
	dec.UseNumber()
	tok, err := dec.Token()
	if err != nil {
		return false, err
	}

	switch v := tok.(type) {
	case json.Number:
		if pi != nil {
			i, err := v.Int64()
			if err == nil {
				*pi = &i
				return false, nil
			}
		}
		if pf != nil {
			f, err := v.Float64()
			if err == nil {
				*pf = &f
				return false, nil
			}
			return false, errors.New("Unparsable number")
		}
		return false, errors.New("Union does not contain number")
	case float64:
		return false, errors.New("Decoder should not return float64")
	case bool:
		if pb != nil {
			*pb = &v
			return false, nil
		}
		return false, errors.New("Union does not contain bool")
	case string:
		if haveEnum {
			return false, json.Unmarshal(data, pe)
		}
		if ps != nil {
			*ps = &v
			return false, nil
		}
		return false, errors.New("Union does not contain string")
	case nil:
		if nullable {
			return false, nil
		}
		return false, errors.New("Union does not contain null")
	case json.Delim:
		if v == '{' {
			if haveObject {
				return true, json.Unmarshal(data, pc)
			}
			if haveMap {
				return false, json.Unmarshal(data, pm)
			}
			return false, errors.New("Union does not contain object")
		}
		if v == '[' {
			if haveArray {
				return false, json.Unmarshal(data, pa)
			}
			return false, errors.New("Union does not contain array")
		}
		return false, errors.New("Cannot handle delimiter")
	}
	return false, errors.New("Cannot unmarshal union")

}

func marshalUnion(pi *int64, pf *float64, pb *bool, ps *string, haveArray bool, pa interface{}, haveObject bool, pc interface{}, haveMap bool, pm interface{}, haveEnum bool, pe interface{}, nullable bool) ([]byte, error) {
	if pi != nil {
		return json.Marshal(*pi)
	}
	if pf != nil {
		return json.Marshal(*pf)
	}
	if pb != nil {
		return json.Marshal(*pb)
	}
	if ps != nil {
		return json.Marshal(*ps)
	}
	if haveArray {
		return json.Marshal(pa)
	}
	if haveObject {
		return json.Marshal(pc)
	}
	if haveMap {
		return json.Marshal(pm)
	}
	if haveEnum {
		return json.Marshal(pe)
	}
	if nullable {
		return json.Marshal(nil)
	}
	return nil, errors.New("Union must not be null")
}

type NodeType string

const (
	VAEDecode                   NodeType = "VAEDecode"
	VAEEncode                   NodeType = "VAEEncode"
	UpscaleModelLoader          NodeType = "Upscale Model Loader"
	CRModuleInput               NodeType = "CR Module Input"
	ImageUpscaleWithModel       NodeType = "ImageUpscaleWithModel"
	SeedNode                    NodeType = "Seed (rgthree)"
	PreviewImage                NodeType = "PreviewImage"
	VAELoader                   NodeType = "VAELoader"
	CRModulePipeLoader          NodeType = "CR Module Pipe Loader"
	ConditioningConcat          NodeType = "ConditioningConcat"
	SaveImage                   NodeType = "SaveImage"
	CLIPTextEncode              NodeType = "CLIPTextEncode"
	ModelMergeSimple            NodeType = "ModelMergeSimple"
	Note                        NodeType = "Note"
	FreeU_V2                    NodeType = "FreeU_V2"
	CheckpointLoaderSimple      NodeType = "CheckpointLoaderSimple"
	KSamplerAdvanced            NodeType = "KSamplerAdvanced"
	KSamplerCycle               NodeType = "KSampler Cycle"
	CRApplyLoRAStack            NodeType = "CR Apply LoRA Stack"
	CLIPSetLastLayer            NodeType = "CLIPSetLastLayer"
	CRApplyModelMerge           NodeType = "CR Apply Model Merge"
	CRLoRAStack                 NodeType = "CR LoRA Stack"
	LoraLoader                  NodeType = "LoraLoader"
	FastGroupsBypasser          NodeType = "Fast Groups Bypasser (rgthree)"
	EmptyLatentImage            NodeType = "EmptyLatentImage"
	CRModelMergeStack           NodeType = "CR Model Merge Stack"
	FastGroupsMuter             NodeType = "Fast Groups Muter (rgthree)"
	SimpleCounter               NodeType = "Simple Counter"
	KSampler                    NodeType = "KSampler"
	KSamplerEfficient           NodeType = "KSampler (Efficient)"
	UltralyticsDetectorProvider NodeType = "UltralyticsDetectorProvider"
	FaceDetailer                NodeType = "FaceDetailer"
	Reroute                     NodeType = "Reroute"
	UltimateSDUpscale           NodeType = "UltimateSDUpscale"
	ModelSamplingDiscrete       NodeType = "ModelSamplingDiscrete"
	LatentUpscaleBy             NodeType = "LatentUpscaleBy"
	String                      NodeType = "String"
	TextToString                NodeType = "Text to String"
	TextMultiline               NodeType = "Text Multiline"
	ControlNetLoader            NodeType = "ControlNetLoader"
	ControlNetApply             NodeType = "ControlNetApply"
	ImageSender                 NodeType = "ImageSender"
	LatentSender                NodeType = "LatentSender"
	ImageReceiver               NodeType = "ImageReceiver"
	LatentReceiver              NodeType = "LatentReceiver"
	PreviewBridge               NodeType = "PreviewBridge"
	SetLatentNoiseMask          NodeType = "SetLatentNoiseMask"
	VAEEncodeForInpaint         NodeType = "VAEEncodeForInpaint"
	MaskPainter                 NodeType = "MaskPainter"
	LatentUpscale               NodeType = "LatentUpscale"
	PrimitiveNode               NodeType = "PrimitiveNode"
	CLIPTextEncodeSDXL          NodeType = "CLIPTextEncodeSDXL"
	ImageScaleBy                NodeType = "ImageScaleBy"
	SAMLoader                   NodeType = "SAMLoader"
	RebatchImages               NodeType = "RebatchImages"
	CFGGuider                   NodeType = "CFGGuider"
	SamplerDPMPP_3M_SDE         NodeType = "SamplerDPMPP_3M_SDE"
	AlignYourStepsScheduler     NodeType = "AlignYourStepsScheduler"
	SamplerCustomAdvanced       NodeType = "SamplerCustomAdvanced"
	CLIPMergeSimple             NodeType = "CLIPMergeSimple"
	RandomNoise                 NodeType = "RandomNoise"
	LoadImage                   NodeType = "LoadImage"
	VRAM_Debug                  NodeType = "VRAM_Debug"
	VAEDecodeTiled              NodeType = "VAEDecodeTiled"
	DPRandomGenerator           NodeType = "DPRandomGenerator"
	SaveTextFile                NodeType = "Save Text File"
	LoraLoaderStack             NodeType = "Lora Loader Stack (rgthree)"
	StupidSimpleNumber          NodeType = "Stupid Simple Number (INT)"
	JDC_Plasma                  NodeType = "JDC_Plasma"
	StupidSimpleSeed            NodeType = "Stupid Simple Seed (INT)"
	BNK_CLIPTextEncodeAdvanced  NodeType = "BNK_CLIPTextEncodeAdvanced"
	CheckpointLoader            NodeType = "CheckpointLoader"
	LoraLoaderPys               NodeType = "LoraLoader|pysssss"
	ShowTextPys                 NodeType = "ShowText|pysssss"
	PromptWithStyle             NodeType = "Prompt With Style"
	UnpackSDXLTuple             NodeType = "Unpack SDXL Tuple"
	WorkflowBus                 NodeType = "workflow/Bus"
	ImageComparer               NodeType = "Image Comparer (rgthree)"
	LoRAStacker                 NodeType = "LoRA Stacker"
	EfficientLoader             NodeType = "Efficient Loader"
	EffLoaderSDXL               NodeType = "Eff. Loader SDXL"
	KSamplerSDXL                NodeType = "KSampler SDXL (Eff.)"
	CRModuleOutput              NodeType = "CR Module Output"
	ControlNetApplyAdvanced     NodeType = "ControlNetApplyAdvanced"
	workflowDetailer            NodeType = "workflow/Detailer"
	TensorRTLoader              NodeType = "TensorRTLoader"
	PowerLoraLoader             NodeType = "Power Lora Loader (rgthree)"
	WorkflowPrompts             NodeType = "workflow/Prompts"
	HighResFixScript            NodeType = "HighRes-Fix Script"
	Digital2KSampler            NodeType = "CCF_V0.342_Sampler"
)

func fallback[T any](field *T, fallback T) {
	if field == nil {
		panic("fallback called with nil field")
	}
	if reflect.ValueOf(*field).IsZero() {
		*field = fallback
	}
}

var negatives = []string{
	"bad quality",
	"low quality",
	"worst quality",
	"easynegative",
	"embedding:bwu",
	"embedding:dfc",
	"embedding:ubbp",
	"embedding:updn",
	"embedding:bad-artist",
	"embedding:boring_e621",
}

func (r *ComfyUI) Convert() *entities.TextToImageRequest {
	basic := Basic{
		Nodes:   r.Nodes,
		Version: r.Version,
	}
	return basic.Convert()
}

func (r *Basic) Convert() *entities.TextToImageRequest {
	if r == nil {
		return nil
	}
	var req entities.TextToImageRequest
	var prompt PromptWriter
	var loras = make(map[string]float64)
	for _, node := range r.Nodes {
		if node.WidgetsValues == nil {
			continue
		}
		if node.Mode == ModeMuted {
			continue
		}
		if node.Mode == ModeBypass {
			continue
		}
		switch node.Type {
		case CheckpointLoaderSimple:
			for _, input := range node.WidgetsValues.UnionArray {
				if input.String == nil {
					continue
				}
				req.OverrideSettings.SDModelCheckpoint = input.String
			}
		case CheckpointLoader:
			for i, input := range node.WidgetsValues.UnionArray {
				if i%2 == 1 {
					if input.String == nil {
						continue
					}
					req.OverrideSettings.SDModelCheckpoint = input.String
				}
			}
		case VAELoader:
			for _, input := range node.WidgetsValues.UnionArray {
				if input.String == nil {
					continue
				}
				req.OverrideSettings.SDVae = input.String
			}
		case SamplerDPMPP_3M_SDE:
			req.SamplerName = "DPMPP_3M_SDE"
		case CRModelMergeStack:
			if req.OverrideSettings.SDModelCheckpoint != nil {
				continue
			}
			var currentWeight float64
			for i, input := range node.WidgetsValues.UnionArray {
				// check the 2nd input in groups of 4
				if i%4 != 1 {
					continue
				}
				if input.String == nil {
					continue
				}
				if *input.String == "None" {
					continue
				}
				// check if the previous input is "On"
				previous := node.WidgetsValues.UnionArray[i-1]
				if previous.String == nil {
					continue
				}
				if *previous.String != "On" {
					continue
				}
				// check if the next input (weight) is not zero
				if len(node.WidgetsValues.UnionArray) <= i+1 {
					break
				}
				weight := node.WidgetsValues.UnionArray[i+1]
				if weight.Double == nil {
					continue
				}
				if *weight.Double <= 0 {
					continue
				}
				// prefer the model with the highest weight
				if *weight.Double > currentWeight {
					currentWeight = *weight.Double
					req.OverrideSettings.SDModelCheckpoint = input.String
				}
			}
		case CRLoRAStack:
			var lastLora *string
			var enabled bool
			for i, input := range node.WidgetsValues.UnionArray {
				switch i % 4 {
				case 0:
					if input.String == nil {
						continue
					}
					enabled = *input.String == "On"
				case 1:
					if input.String == nil {
						continue
					}
					if *input.String == "None" {
						enabled = false
						continue
					}
					if !enabled {
						continue
					}
					lastLora = input.String
					loras[*lastLora] = 1
				case 2:
					if input.Double == nil {
						continue
					}
					if !enabled {
						continue
					}
					if lastLora == nil {
						continue
					}
					loras[*lastLora] = *input.Double
					enabled = false
					lastLora = nil
				}
			}
		case LoraLoader, LoraLoaderStack:
			var lastLora *string
			for i, input := range node.WidgetsValues.UnionArray {
				switch i % 2 {
				case 0:
					if input.String == nil {
						continue
					}
					lastLora = input.String
					loras[*lastLora] = 1
				case 1:
					if input.Double == nil {
						continue
					}
					if lastLora == nil {
						continue
					}
					loras[*lastLora] = *input.Double
					lastLora = nil
				}
			}
		case LoraLoaderPys:
			var lastLora *string
			for i, input := range node.WidgetsValues.UnionArray {
				switch i {
				case 0:
					if input.WidgetsValueClass == nil {
						continue
					}
					if input.WidgetsValueClass.Content == nil {
						continue
					}
					if *input.WidgetsValueClass.Content == "None" {
						continue
					}
					lastLora = input.WidgetsValueClass.Content
					loras[*lastLora] = 1
				case 1:
					if input.Double == nil {
						continue
					}
					if lastLora == nil {
						continue
					}
					loras[*lastLora] = *input.Double
					lastLora = nil
				}
			}
		case CLIPTextEncode:
			for _, input := range node.WidgetsValues.UnionArray {
				if input.String == nil {
					continue
				}
				if node.Title != nil && strings.Contains(strings.ToLower(string(*node.Title)), "negative") {
					req.NegativePrompt = *input.String
					continue
				}

				// If we already have a negative prompt
				if req.NegativePrompt != "" {
					prompt.WriteString(strings.TrimSpace(*input.String))
					continue
				}

				var foundNegative bool
				for _, negative := range negatives {
					if strings.Contains(*input.String, negative) {
						req.NegativePrompt = *input.String
						foundNegative = true
						break
					}
				}
				if foundNegative {
					continue
				}

				prompt.WriteString(strings.TrimSpace(*input.String))
			}
		case WorkflowPrompts:
			for i, input := range node.WidgetsValues.UnionArray {
				if input.String == nil {
					continue
				}
				switch i {
				case 0, 1:
					prompt.WriteString(strings.TrimSpace(*input.String))
				case 2:
					req.NegativePrompt = *input.String
				}
			}
		case BNK_CLIPTextEncodeAdvanced:
			for _, input := range node.WidgetsValues.UnionArray {
				if input.String == nil {
					continue
				}
				prompt.WriteString(strings.TrimSpace(*input.String))
				break
			}
		case PromptWithStyle:
			for i, input := range node.WidgetsValues.UnionArray {
				switch i {
				case 0:
					if input.String == nil {
						continue
					}
					prompt.WriteString(strings.TrimSpace(*input.String))
				case 1:
					if input.String == nil {
						continue
					}
					req.NegativePrompt = *input.String
				case 3:
					if input.String == nil {
						continue
					}
					dimensions := regexp.MustCompile(`(\d+)x(\d+)`).FindStringSubmatch(*input.String)
					if len(dimensions) < 3 {
						continue
					}

					width, err := strconv.Atoi(dimensions[1])
					if err == nil {
						req.Width = width
					}

					height, err := strconv.Atoi(dimensions[2])
					if err == nil {
						req.Height = height
					}
				case 4:
					if input.Double == nil {
						continue
					}
					req.BatchSize = int(*input.Double)
				case 5:
					if input.Double == nil {
						continue
					}
					req.Seed = int64(*input.Double)
				}
			}
		case RandomNoise:
			for _, input := range node.WidgetsValues.UnionArray {
				if input.Double == nil {
					continue
				}
				req.Seed = int64(*input.Double)
				break
			}
		case SeedNode:
			for _, input := range node.WidgetsValues.UnionArray {
				if input.Double == nil {
					continue
				}
				req.Seed = int64(*input.Double)
				break
			}
		case KSamplerAdvanced:
			for i, input := range node.WidgetsValues.UnionArray {
				switch i {
				case 1:
					if input.Double == nil {
						continue
					}
					req.Seed = int64(*input.Double)
				case 3:
					if input.Double == nil {
						continue
					}
					req.Steps = int(*input.Double)
				case 4:
					if input.Double == nil {
						continue
					}
					req.CFGScale = *input.Double
				case 5:
					if input.String == nil {
						continue
					}
					req.SamplerName = *input.String
				case 6:
					if input.String == nil {
						continue
					}
					req.Scheduler = input.String
				}
			}
		case KSamplerCycle:
			for i, input := range node.WidgetsValues.UnionArray {
				switch i {
				case 0:
					if input.Double == nil {
						continue
					}
					fallback(&req.Seed, int64(*input.Double))
				case 2:
					if input.Double == nil {
						continue
					}
					req.Steps = int(*input.Double)
				case 3:
					if input.Double == nil {
						continue
					}
					req.CFGScale = *input.Double
				case 4:
					if input.String == nil {
						continue
					}
					req.SamplerName = *input.String
				case 8:
					if input.Double == nil {
						continue
					}
					req.HrScale = *input.Double
				}
			}
		case KSampler:
			for i, input := range node.WidgetsValues.UnionArray {
				switch i {
				case 0:
					if input.Double == nil {
						continue
					}
					fallback(&req.Seed, int64(*input.Double))
				case 2:
					if input.Double == nil {
						continue
					}
					fallback(&req.Steps, int(*input.Double))
				case 3:
					if input.Double == nil {
						continue
					}
					fallback(&req.CFGScale, *input.Double)
				case 4:
					if input.String == nil {
						continue
					}
					fallback(&req.SamplerName, *input.String)
				case 5:
					if input.String == nil {
						continue
					}
					fallback(&req.Scheduler, input.String)
				case 6:
					if input.Double == nil {
						continue
					}
					fallback(&req.DenoisingStrength, *input.Double)
				}
			}
		case EfficientLoader:
			var lastLora *string
			for i, input := range node.WidgetsValues.UnionArray {
				switch i {
				case 0:
					if input.String == nil {
						continue
					}
					req.OverrideSettings.SDModelCheckpoint = input.String
				case 1:
					if input.String == nil {
						continue
					}
					req.OverrideSettings.SDVae = input.String
				case 2:
					if input.Double == nil {
						continue
					}
					req.OverrideSettings.CLIPStopAtLastLayers = *input.Double
				case 3:
					if input.String == nil {
						continue
					}
					if *input.String == "None" {
						continue
					}
					loras[*input.String] = 1
					lastLora = input.String
				case 4:
					if input.Double == nil {
						continue
					}
					if lastLora == nil {
						continue
					}
					loras[*lastLora] = *input.Double
					lastLora = nil
				case 6:
					if input.String == nil {
						continue
					}
					prompt.WriteString(*input.String)
				case 7:
					if input.String == nil {
						continue
					}
					req.NegativePrompt = *input.String
				case 10:
					if input.Double == nil {
						continue
					}
					req.Width = int(*input.Double)
				case 11:
					if input.Double == nil {
						continue
					}
					req.Height = int(*input.Double)
				case 12:
					if input.Double == nil {
						continue
					}
					req.BatchSize = int(*input.Double)
				}
			}
		case EffLoaderSDXL:
			for i, input := range node.WidgetsValues.UnionArray {
				switch i {
				case 0:
					if input.String == nil {
						continue
					}
					req.OverrideSettings.SDModelCheckpoint = input.String
				case 1:
					if input.Double == nil {
						continue
					}
					req.OverrideSettings.CLIPStopAtLastLayers = *input.Double
				case 2: // Refiner model
				case 3: // Refiner CLIP skip
				case 4: // Refiner Positive A score
				case 5: // Refiner Negative A score
				case 6:
					if input.String == nil {
						continue
					}
					req.OverrideSettings.SDVae = input.String
				case 7:
					if input.String == nil {
						continue
					}
					prompt.WriteString(*input.String)
				case 8:
					if input.String == nil {
						continue
					}
					req.NegativePrompt = *input.String
				case 9: // Token normalization
				case 10: // Weight interpolation
				case 11:
					if input.Double == nil {
						continue
					}
					req.Width = int(*input.Double)
				case 12:
					if input.Double == nil {
						continue
					}
					req.Height = int(*input.Double)
				case 13:
					if input.Double == nil {
						continue
					}
					req.BatchSize = int(*input.Double)
				}
			}
		case KSamplerSDXL, KSamplerEfficient:
			for i, input := range node.WidgetsValues.UnionArray {
				switch i {
				case 0:
					if input.Double == nil {
						continue
					}
					req.Seed = int64(*input.Double)
				case 1:
					if input.Double == nil {
						continue
					}
					req.Seed = int64(*input.Double)
				case 2:
					if input.Double == nil {
						continue
					}
					req.Steps = int(*input.Double)
				case 3:
					if input.Double == nil {
						continue
					}
					req.CFGScale = *input.Double
				case 4:
					if input.String == nil {
						continue
					}
					req.SamplerName = *input.String
				case 5:
					if input.String == nil {
						continue
					}
					req.Scheduler = input.String
				}
			}
		case CFGGuider:
			if req.CFGScale != 0 {
				continue
			}
			for _, input := range node.WidgetsValues.UnionArray {
				if input.Double == nil {
					continue
				}
				req.CFGScale = *input.Double
				break
			}
		case CRModulePipeLoader:
			for _, input := range node.WidgetsValues.UnionArray {
				if input.Double == nil {
					continue
				}
				req.Seed = int64(*input.Double)
				break
			}
		case AlignYourStepsScheduler:
			for i, input := range node.WidgetsValues.UnionArray {
				if i != 1 {
					continue
				}
				if input.Double == nil {
					continue
				}
				req.Steps = int(*input.Double)
				break
			}
		case DPRandomGenerator:
			for _, input := range node.WidgetsValues.UnionArray {
				if input.String == nil {
					continue
				}
				prompt.WriteString(*input.String)
				break
			}
		}
	}

	for lora, weight := range loras {
		prompt.WriteString(fmt.Sprintf("<lora:%s:%.2f>", lora, weight))
	}

	req.Prompt = prompt.String()

	return &req
}

type PromptWriter strings.Builder

type Prompter interface {
	strings.Builder
	WriteString(string)
	String() string
}

func (p *PromptWriter) WriteString(s string) {
	if (*strings.Builder)(p).Len() > 0 {
		(*strings.Builder)(p).WriteString("\n")
	}
	(*strings.Builder)(p).WriteString(s)
}

func (p *PromptWriter) String() string {
	return (*strings.Builder)(p).String()
}
