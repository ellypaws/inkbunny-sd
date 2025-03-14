package entities

type ControlModeV2 = string

const (
	ControlModeV2Balanced    ControlModeV2 = "balanced"
	ControlModeV2MorePrompt  ControlModeV2 = "more_prompt"
	ControlModeV2MoreControl ControlModeV2 = "more_control"
	ControlModeV2Unbalanced  ControlModeV2 = "unbalanced"
)

type CLIPVisionModelV2 = string

const (
	CLIPVisionModelV2ViTH CLIPVisionModelV2 = "ViT-H"
	CLIPVisionModelV2ViTG CLIPVisionModelV2 = "ViT-G"
	CLIPVisionModelV2ViTL CLIPVisionModelV2 = "ViT-L"
)

type IPMethodV2 = string

const (
	IPMethodV2Full        IPMethodV2 = "full"
	IPMethodV2Style       IPMethodV2 = "style"
	IPMethodV2Composition IPMethodV2 = "composition"
)

type Tool = string

const (
	ToolBrush       Tool = "brush"
	ToolEraser      Tool = "eraser"
	ToolMove        Tool = "move"
	ToolRect        Tool = "rect"
	ToolView        Tool = "view"
	ToolBbox        Tool = "bbox"
	ToolColorPicker Tool = "colorPicker"
)

type FillStyle = string

const (
	FillStyleSolid      FillStyle = "solid"
	FillStyleGrid       FillStyle = "grid"
	FillStyleCrosshatch FillStyle = "crosshatch"
	FillStyleDiagonal   FillStyle = "diagonal"
	FillStyleHorizontal FillStyle = "horizontal"
	FillStyleVertical   FillStyle = "vertical"
)

type BoundingBoxScaleMethod = string

const (
	BoundingBoxScaleMethodNone   BoundingBoxScaleMethod = "none"
	BoundingBoxScaleMethodAuto   BoundingBoxScaleMethod = "auto"
	BoundingBoxScaleMethodManual BoundingBoxScaleMethod = "manual"
)

type AspectRatioID = string

const (
	AspectRatioIDFree AspectRatioID = "Free"
	AspectRatioID16_9 AspectRatioID = "16:9"
	AspectRatioID3_2  AspectRatioID = "3:2"
	AspectRatioID4_3  AspectRatioID = "4:3"
	AspectRatioID1_1  AspectRatioID = "1:1"
	AspectRatioID3_4  AspectRatioID = "3:4"
	AspectRatioID2_3  AspectRatioID = "2:3"
	AspectRatioID9_16 AspectRatioID = "9:16"
)
