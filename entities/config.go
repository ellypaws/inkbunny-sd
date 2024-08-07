package entities

import "encoding/json"

func UnmarshalConfig(data []byte) (Config, error) {
	var r Config
	err := json.Unmarshal(data, &r)
	return r, err
}

func (config *Config) Marshal() ([]byte, error) {
	return json.Marshal(config)
}

type Config struct {
	SamplesSave                           bool     `json:"samples_save,omitempty"`
	SamplesFormat                         string   `json:"samples_format,omitempty"`
	SamplesFilenamePattern                string   `json:"samples_filename_pattern,omitempty"`
	SaveImagesAddNumber                   bool     `json:"save_images_add_number,omitempty"`
	SaveImagesReplaceAction               string   `json:"save_images_replace_action,omitempty"`
	GridSave                              bool     `json:"grid_save,omitempty"`
	GridFormat                            string   `json:"grid_format,omitempty"`
	GridExtendedFilename                  bool     `json:"grid_extended_filename,omitempty"`
	GridOnlyIfMultiple                    bool     `json:"grid_only_if_multiple,omitempty"`
	GridPreventEmptySpots                 bool     `json:"grid_prevent_empty_spots,omitempty"`
	GridZipFilenamePattern                string   `json:"grid_zip_filename_pattern,omitempty"`
	NRows                                 float64  `json:"n_rows,omitempty"`
	Font                                  string   `json:"font,omitempty"`
	GridTextActiveColor                   string   `json:"grid_text_active_color,omitempty"`
	GridTextInactiveColor                 string   `json:"grid_text_inactive_color,omitempty"`
	GridBackgroundColor                   string   `json:"grid_background_color,omitempty"`
	EnablePnginfo                         bool     `json:"enable_pnginfo,omitempty"`
	SaveTxt                               bool     `json:"save_txt,omitempty"`
	SaveImagesBeforeFaceRestoration       bool     `json:"save_images_before_face_restoration,omitempty"`
	SaveImagesBeforeHighresFix            bool     `json:"save_images_before_highres_fix,omitempty"`
	SaveImagesBeforeColorCorrection       bool     `json:"save_images_before_color_correction,omitempty"`
	SaveMask                              bool     `json:"save_mask,omitempty"`
	SaveMaskComposite                     bool     `json:"save_mask_composite,omitempty"`
	JPEGQuality                           float64  `json:"jpeg_quality,omitempty"`
	WebpLossless                          bool     `json:"webp_lossless,omitempty"`
	ExportFor4Chan                        bool     `json:"export_for_4chan,omitempty"`
	ImgDownscaleThreshold                 float64  `json:"img_downscale_threshold,omitempty"`
	TargetSideLength                      float64  `json:"target_side_length,omitempty"`
	ImgMaxSizeMp                          float64  `json:"img_max_size_mp,omitempty"`
	UseOriginalNameBatch                  bool     `json:"use_original_name_batch,omitempty"`
	UseUpscalerNameAsSuffix               bool     `json:"use_upscaler_name_as_suffix,omitempty"`
	SaveSelectedOnly                      bool     `json:"save_selected_only,omitempty"`
	SaveInitImg                           bool     `json:"save_init_img,omitempty"`
	TempDir                               string   `json:"temp_dir,omitempty"`
	CleanTempDirAtStart                   bool     `json:"clean_temp_dir_at_start,omitempty"`
	SaveIncompleteImages                  bool     `json:"save_incomplete_images,omitempty"`
	NotificationAudio                     bool     `json:"notification_audio,omitempty"`
	NotificationVolume                    float64  `json:"notification_volume,omitempty"`
	OutdirSamples                         string   `json:"outdir_samples,omitempty"`
	OutdirTxt2ImgSamples                  string   `json:"outdir_txt2img_samples,omitempty"`
	OutdirImg2ImgSamples                  string   `json:"outdir_img2img_samples,omitempty"`
	OutdirExtrasSamples                   string   `json:"outdir_extras_samples,omitempty"`
	OutdirGrids                           string   `json:"outdir_grids,omitempty"`
	OutdirTxt2ImgGrids                    string   `json:"outdir_txt2img_grids,omitempty"`
	OutdirImg2ImgGrids                    string   `json:"outdir_img2img_grids,omitempty"`
	OutdirSave                            string   `json:"outdir_save,omitempty"`
	OutdirInitImages                      string   `json:"outdir_init_images,omitempty"`
	SaveToDirs                            bool     `json:"save_to_dirs,omitempty"`
	GridSaveToDirs                        bool     `json:"grid_save_to_dirs,omitempty"`
	UseSaveToDirsForUI                    bool     `json:"use_save_to_dirs_for_ui,omitempty"`
	DirectoriesFilenamePattern            string   `json:"directories_filename_pattern,omitempty"`
	DirectoriesMaxPromptWords             float64  `json:"directories_max_prompt_words,omitempty"`
	ESRGANTile                            float64  `json:"ESRGAN_tile,omitempty"`
	ESRGANTileOverlap                     float64  `json:"ESRGAN_tile_overlap,omitempty"`
	RealesrganEnabledModels               []string `json:"realesrgan_enabled_models,omitempty"`
	UpscalerForImg2Img                    string   `json:"upscaler_for_img2img,omitempty"`
	FaceRestoration                       bool     `json:"face_restoration,omitempty"`
	FaceRestorationModel                  string   `json:"face_restoration_model,omitempty"`
	CodeFormerWeight                      float64  `json:"code_former_weight,omitempty"`
	FaceRestorationUnload                 bool     `json:"face_restoration_unload,omitempty"`
	AutoLaunchBrowser                     string   `json:"auto_launch_browser,omitempty"`
	EnableConsolePrompts                  bool     `json:"enable_console_prompts,omitempty"`
	ShowWarnings                          bool     `json:"show_warnings,omitempty"`
	ShowGradioDeprecationWarnings         bool     `json:"show_gradio_deprecation_warnings,omitempty"`
	MemmonPollRate                        float64  `json:"memmon_poll_rate,omitempty"`
	SamplesLogStdout                      bool     `json:"samples_log_stdout,omitempty"`
	MultipleTqdm                          bool     `json:"multiple_tqdm,omitempty"`
	PrintHypernetExtra                    bool     `json:"print_hypernet_extra,omitempty"`
	ListHiddenFiles                       bool     `json:"list_hidden_files,omitempty"`
	DisableMmapLoadSafetensors            bool     `json:"disable_mmap_load_safetensors,omitempty"`
	HideLdmPrints                         bool     `json:"hide_ldm_prints,omitempty"`
	DumpStacksOnSignal                    bool     `json:"dump_stacks_on_signal,omitempty"`
	APIEnableRequests                     bool     `json:"api_enable_requests,omitempty"`
	APIForbidLocalRequests                bool     `json:"api_forbid_local_requests,omitempty"`
	APIUseragent                          string   `json:"api_useragent,omitempty"`
	UnloadModelsWhenTraining              bool     `json:"unload_models_when_training,omitempty"`
	PinMemory                             bool     `json:"pin_memory,omitempty"`
	SaveOptimizerState                    bool     `json:"save_optimizer_state,omitempty"`
	SaveTrainingSettingsToTxt             bool     `json:"save_training_settings_to_txt,omitempty"`
	DatasetFilenameWordRegex              string   `json:"dataset_filename_word_regex,omitempty"`
	DatasetFilenameJoinString             string   `json:"dataset_filename_join_string,omitempty"`
	TrainingImageRepeatsPerEpoch          float64  `json:"training_image_repeats_per_epoch,omitempty"`
	TrainingWriteCSVEvery                 float64  `json:"training_write_csv_every,omitempty"`
	TrainingXattentionOptimizations       bool     `json:"training_xattention_optimizations,omitempty"`
	TrainingEnableTensorboard             bool     `json:"training_enable_tensorboard,omitempty"`
	TrainingTensorboardSaveImages         bool     `json:"training_tensorboard_save_images,omitempty"`
	TrainingTensorboardFlushEvery         float64  `json:"training_tensorboard_flush_every,omitempty"`
	SDModelCheckpoint                     *string  `json:"sd_model_checkpoint,omitempty"`
	SDCheckpointsLimit                    float64  `json:"sd_checkpoints_limit,omitempty"`
	SDCheckpointsKeepInCPU                bool     `json:"sd_checkpoints_keep_in_cpu,omitempty"`
	SDCheckpointCache                     float64  `json:"sd_checkpoint_cache,omitempty"`
	SDUnet                                string   `json:"sd_unet,omitempty"`
	EnableQuantization                    bool     `json:"enable_quantization,omitempty"`
	EnableEmphasis                        bool     `json:"enable_emphasis,omitempty"`
	EnableBatchSeeds                      bool     `json:"enable_batch_seeds,omitempty"`
	CommaPaddingBacktrack                 float64  `json:"comma_padding_backtrack,omitempty"`
	CLIPStopAtLastLayers                  float64  `json:"CLIP_stop_at_last_layers,omitempty"`
	UpcastAttn                            bool     `json:"upcast_attn,omitempty"`
	RandnSource                           string   `json:"randn_source,omitempty"`
	Tiling                                bool     `json:"tiling,omitempty"`
	HiresFixRefinerPass                   string   `json:"hires_fix_refiner_pass,omitempty"`
	SdxlCropTop                           float64  `json:"sdxl_crop_top,omitempty"`
	SdxlCropLeft                          float64  `json:"sdxl_crop_left,omitempty"`
	SdxlRefinerLowAestheticScore          float64  `json:"sdxl_refiner_low_aesthetic_score,omitempty"`
	SdxlRefinerHighAestheticScore         float64  `json:"sdxl_refiner_high_aesthetic_score,omitempty"`
	SDVaeExplanation                      string   `json:"sd_vae_explanation,omitempty"`
	SDVaeCheckpointCache                  float64  `json:"sd_vae_checkpoint_cache,omitempty"`
	SDVae                                 *string  `json:"sd_vae,omitempty"`
	SDVaeOverridesPerModelPreferences     bool     `json:"sd_vae_overrides_per_model_preferences,omitempty"`
	AutoVaePrecision                      bool     `json:"auto_vae_precision,omitempty"`
	SDVaeEncodeMethod                     string   `json:"sd_vae_encode_method,omitempty"`
	SDVaeDecodeMethod                     string   `json:"sd_vae_decode_method,omitempty"`
	InpaintingMaskWeight                  float64  `json:"inpainting_mask_weight,omitempty"`
	InitialNoiseMultiplier                float64  `json:"initial_noise_multiplier,omitempty"`
	Img2ImgExtraNoise                     float64  `json:"img2img_extra_noise,omitempty"`
	Img2ImgColorCorrection                bool     `json:"img2img_color_correction,omitempty"`
	Img2ImgFixSteps                       bool     `json:"img2img_fix_steps,omitempty"`
	Img2ImgBackgroundColor                string   `json:"img2img_background_color,omitempty"`
	Img2ImgEditorHeight                   float64  `json:"img2img_editor_height,omitempty"`
	Img2ImgSketchDefaultBrushColor        string   `json:"img2img_sketch_default_brush_color,omitempty"`
	Img2ImgInpaintMaskBrushColor          string   `json:"img2img_inpaint_mask_brush_color,omitempty"`
	Img2ImgInpaintSketchDefaultBrushColor string   `json:"img2img_inpaint_sketch_default_brush_color,omitempty"`
	ReturnMask                            bool     `json:"return_mask,omitempty"`
	ReturnMaskComposite                   bool     `json:"return_mask_composite,omitempty"`
	CrossAttentionOptimization            string   `json:"cross_attention_optimization,omitempty"`
	SMinUncond                            float64  `json:"s_min_uncond,omitempty"`
	TokenMergingRatio                     float64  `json:"token_merging_ratio,omitempty"`
	TokenMergingRatioImg2Img              float64  `json:"token_merging_ratio_img2img,omitempty"`
	TokenMergingRatioHr                   float64  `json:"token_merging_ratio_hr,omitempty"`
	PadCondUncond                         bool     `json:"pad_cond_uncond,omitempty"`
	PersistentCondCache                   bool     `json:"persistent_cond_cache,omitempty"`
	BatchCondUncond                       bool     `json:"batch_cond_uncond,omitempty"`
	UseOldEmphasisImplementation          bool     `json:"use_old_emphasis_implementation,omitempty"`
	UseOldKarrasSchedulerSigmas           bool     `json:"use_old_karras_scheduler_sigmas,omitempty"`
	NoDpmppSdeBatchDeterminism            bool     `json:"no_dpmpp_sde_batch_determinism,omitempty"`
	UseOldHiresFixWidthHeight             bool     `json:"use_old_hires_fix_width_height,omitempty"`
	DontFixSecondOrderSamplersSchedule    bool     `json:"dont_fix_second_order_samplers_schedule,omitempty"`
	HiresFixUseFirstpassConds             bool     `json:"hires_fix_use_firstpass_conds,omitempty"`
	UseOldScheduling                      bool     `json:"use_old_scheduling,omitempty"`
	InterrogateKeepModelsInMemory         bool     `json:"interrogate_keep_models_in_memory,omitempty"`
	InterrogateReturnRanks                bool     `json:"interrogate_return_ranks,omitempty"`
	InterrogateClipNumBeams               float64  `json:"interrogate_clip_num_beams,omitempty"`
	InterrogateClipMinLength              float64  `json:"interrogate_clip_min_length,omitempty"`
	InterrogateClipMaxLength              float64  `json:"interrogate_clip_max_length,omitempty"`
	InterrogateClipDictLimit              float64  `json:"interrogate_clip_dict_limit,omitempty"`
	InterrogateClipSkipCategories         []string `json:"interrogate_clip_skip_categories,omitempty"`
	InterrogateDeepbooruScoreThreshold    float64  `json:"interrogate_deepbooru_score_threshold,omitempty"`
	DeepbooruSortAlpha                    bool     `json:"deepbooru_sort_alpha,omitempty"`
	DeepbooruUseSpaces                    bool     `json:"deepbooru_use_spaces,omitempty"`
	DeepbooruEscape                       bool     `json:"deepbooru_escape,omitempty"`
	DeepbooruFilterTags                   string   `json:"deepbooru_filter_tags,omitempty"`
	ExtraNetworksShowHiddenDirectories    bool     `json:"extra_networks_show_hidden_directories,omitempty"`
	ExtraNetworksHiddenModels             string   `json:"extra_networks_hidden_models,omitempty"`
	ExtraNetworksDefaultMultiplier        float64  `json:"extra_networks_default_multiplier,omitempty"`
	ExtraNetworksCardWidth                float64  `json:"extra_networks_card_width,omitempty"`
	ExtraNetworksCardHeight               float64  `json:"extra_networks_card_height,omitempty"`
	ExtraNetworksCardTextScale            float64  `json:"extra_networks_card_text_scale,omitempty"`
	ExtraNetworksCardShowDesc             bool     `json:"extra_networks_card_show_desc,omitempty"`
	ExtraNetworksCardOrderField           string   `json:"extra_networks_card_order_field,omitempty"`
	ExtraNetworksCardOrder                string   `json:"extra_networks_card_order,omitempty"`
	ExtraNetworksAddTextSeparator         string   `json:"extra_networks_add_text_separator,omitempty"`
	UIExtraNetworksTabReorder             string   `json:"ui_extra_networks_tab_reorder,omitempty"`
	TextualInversionPrintAtLoad           bool     `json:"textual_inversion_print_at_load,omitempty"`
	TextualInversionAddHashesToInfotext   bool     `json:"textual_inversion_add_hashes_to_infotext,omitempty"`
	SDHypernetwork                        *string  `json:"sd_hypernetwork,omitempty"`
	Localization                          string   `json:"localization,omitempty"`
	GradioTheme                           string   `json:"gradio_theme,omitempty"`
	GradioThemesCache                     bool     `json:"gradio_themes_cache,omitempty"`
	GalleryHeight                         string   `json:"gallery_height,omitempty"`
	ReturnGrid                            bool     `json:"return_grid,omitempty"`
	DoNotShowImages                       bool     `json:"do_not_show_images,omitempty"`
	SendSeed                              bool     `json:"send_seed,omitempty"`
	SendSize                              bool     `json:"send_size,omitempty"`
	JSModalLightbox                       bool     `json:"js_modal_lightbox,omitempty"`
	JSModalLightboxInitiallyZoomed        bool     `json:"js_modal_lightbox_initially_zoomed,omitempty"`
	JSModalLightboxGamepad                bool     `json:"js_modal_lightbox_gamepad,omitempty"`
	JSModalLightboxGamepadRepeat          float64  `json:"js_modal_lightbox_gamepad_repeat,omitempty"`
	ShowProgressInTitle                   bool     `json:"show_progress_in_title,omitempty"`
	SamplersInDropdown                    bool     `json:"samplers_in_dropdown,omitempty"`
	DimensionsAndBatchTogether            bool     `json:"dimensions_and_batch_together,omitempty"`
	KeyeditPrecisionAttention             float64  `json:"keyedit_precision_attention,omitempty"`
	KeyeditPrecisionExtra                 float64  `json:"keyedit_precision_extra,omitempty"`
	KeyeditDelimiters                     string   `json:"keyedit_delimiters,omitempty"`
	KeyeditDelimitersWhitespace           []string `json:"keyedit_delimiters_whitespace,omitempty"`
	KeyeditMove                           bool     `json:"keyedit_move,omitempty"`
	QuicksettingsList                     []string `json:"quicksettings_list,omitempty"`
	UITabOrder                            []string `json:"ui_tab_order,omitempty"`
	HiddenTabs                            []string `json:"hidden_tabs,omitempty"`
	UIReorderList                         []string `json:"ui_reorder_list,omitempty"`
	SDCheckpointDropdownUseShort          bool     `json:"sd_checkpoint_dropdown_use_short,omitempty"`
	HiresFixShowSampler                   bool     `json:"hires_fix_show_sampler,omitempty"`
	HiresFixShowPrompts                   bool     `json:"hires_fix_show_prompts,omitempty"`
	DisableTokenCounters                  bool     `json:"disable_token_counters,omitempty"`
	CompactPromptBox                      bool     `json:"compact_prompt_box,omitempty"`
	AddModelHashToInfo                    bool     `json:"add_model_hash_to_info,omitempty"`
	AddModelNameToInfo                    bool     `json:"add_model_name_to_info,omitempty"`
	AddUserNameToInfo                     bool     `json:"add_user_name_to_info,omitempty"`
	AddVersionToInfotext                  bool     `json:"add_version_to_infotext,omitempty"`
	DisableWeightsAutoSwap                bool     `json:"disable_weights_auto_swap,omitempty"`
	InfotextStyles                        string   `json:"infotext_styles,omitempty"`
	ShowProgressbar                       bool     `json:"show_progressbar,omitempty"`
	LivePreviewsEnable                    bool     `json:"live_previews_enable,omitempty"`
	LivePreviewsImageFormat               string   `json:"live_previews_image_format,omitempty"`
	ShowProgressGrid                      bool     `json:"show_progress_grid,omitempty"`
	ShowProgressEveryNSteps               float64  `json:"show_progress_every_n_steps,omitempty"`
	ShowProgressType                      string   `json:"show_progress_type,omitempty"`
	LivePreviewAllowLowvramFull           bool     `json:"live_preview_allow_lowvram_full,omitempty"`
	LivePreviewContent                    string   `json:"live_preview_content,omitempty"`
	LivePreviewRefreshPeriod              float64  `json:"live_preview_refresh_period,omitempty"`
	LivePreviewFastInterrupt              bool     `json:"live_preview_fast_interrupt,omitempty"`
	HideSamplers                          []string `json:"hide_samplers,omitempty"`
	EtaDdim                               float64  `json:"eta_ddim,omitempty"`
	EtaAncestral                          float64  `json:"eta_ancestral,omitempty"`
	DdimDiscretize                        string   `json:"ddim_discretize,omitempty"`
	SChurn                                float64  `json:"s_churn,omitempty"`
	STmin                                 float64  `json:"s_tmin,omitempty"`
	STmax                                 float64  `json:"s_tmax,omitempty"`
	SNoise                                float64  `json:"s_noise,omitempty"`
	KSchedType                            string   `json:"k_sched_type,omitempty"`
	SigmaMin                              float64  `json:"sigma_min,omitempty"`
	SigmaMax                              float64  `json:"sigma_max,omitempty"`
	Rho                                   float64  `json:"rho,omitempty"`
	EtaNoiseSeedDelta                     float64  `json:"eta_noise_seed_delta,omitempty"`
	AlwaysDiscardNextToLastSigma          bool     `json:"always_discard_next_to_last_sigma,omitempty"`
	SgmNoiseMultiplier                    bool     `json:"sgm_noise_multiplier,omitempty"`
	UniPCVariant                          string   `json:"uni_pc_variant,omitempty"`
	UniPCSkipType                         string   `json:"uni_pc_skip_type,omitempty"`
	UniPCOrder                            float64  `json:"uni_pc_order,omitempty"`
	UniPCLowerOrderFinal                  bool     `json:"uni_pc_lower_order_final,omitempty"`
	PostprocessingEnableInMainUI          []string `json:"postprocessing_enable_in_main_ui,omitempty"`
	PostprocessingOperationOrder          []string `json:"postprocessing_operation_order,omitempty"`
	UpscalingMaxImagesInCache             float64  `json:"upscaling_max_images_in_cache,omitempty"`
	DisabledExtensions                    []string `json:"disabled_extensions,omitempty"`
	DisableAllExtensions                  string   `json:"disable_all_extensions,omitempty"`
	RestoreConfigStateFile                string   `json:"restore_config_state_file,omitempty"`
	SDCheckpointHash                      string   `json:"sd_checkpoint_hash,omitempty"`
	// Downcast model alphas_cumprod to fp16 before sampling.
	// [For reproducing old seeds].
	// Set to true to use the old behavior.
	//
	// -----
	//
	// 1.8.0 (dev: 1.7.0-225) [2024-01-01] - zero terminal SNR noise schedule option
	//
	// 	Slightly changes all image generation.
	//	The PR changes alphas_cumprod to be never be fp16 unless the backwards compatibility option is enabled.
	//	Backwards compatibility option is "Downcast model alphas_cumprod to fp16 before sampling",
	//	and it's automatically enabled when restoring parameters from old pictures
	//	(as long as they have Version: ... in infotext).
	//
	// [For reproducing old seeds]: https://github.com/AUTOMATIC1111/stable-diffusion-webui/wiki/Seed-breaking-changes#180-dev-170-225-2024-01-01---zero-terminal-snr-noise-schedule-option
	// [2024-01-01]: https://github.com/AUTOMATIC1111/stable-diffusion-webui/pull/14145
	DowncastAlphasCumprodToFP16 bool `json:"use_downcasted_alpha_bar,omitempty"`
}
