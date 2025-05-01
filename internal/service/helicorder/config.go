package helicorder

import (
	"errors"
	"fmt"
	"path"

	"github.com/anyshake/observer/config"
	"github.com/anyshake/observer/internal/dao/action"
	"github.com/samber/lo"
)

type helicorderConfigEnabledImpl struct{}

func (s *helicorderConfigEnabledImpl) GetName() string             { return "Enable" }
func (s *helicorderConfigEnabledImpl) GetNamespace() string        { return ID }
func (s *helicorderConfigEnabledImpl) GetKey() string              { return "enabled" }
func (s *helicorderConfigEnabledImpl) GetType() action.SettingType { return action.Bool }
func (s *helicorderConfigEnabledImpl) IsRequired() bool            { return true }
func (s *helicorderConfigEnabledImpl) GetVersion() int             { return 0 }
func (s *helicorderConfigEnabledImpl) GetOptions() map[string]any  { return nil }
func (s *helicorderConfigEnabledImpl) GetDefaultValue() any        { return false }
func (s *helicorderConfigEnabledImpl) GetDescription() string {
	return "Enable helicorder service to save daily waveforms per channel in the specified path, creating it if the path does not exist."
}
func (s *helicorderConfigEnabledImpl) Init(handler *action.Handler) error {
	if _, err := handler.SettingsInit(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), s.GetDefaultValue()); err != nil {
		return fmt.Errorf("failed to set default helicorder service availablity: %w", err)
	}
	return nil
}
func (s *helicorderConfigEnabledImpl) Set(handler *action.Handler, newVal any) error {
	enabled, err := config.GetConfigValBool(newVal)
	if err != nil {
		return err
	}
	if err := handler.SettingsSet(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), enabled); err != nil {
		return fmt.Errorf("failed to set helicorder service availablity: %w", err)
	}
	return nil
}
func (s *helicorderConfigEnabledImpl) Get(handler *action.Handler) (any, error) {
	val, _, _, err := handler.SettingsGet(s.GetNamespace(), s.GetKey())
	if err != nil {
		return nil, fmt.Errorf("failed to get helicorder service availablity: %w", err)
	}
	enabled, ok := val.(bool)
	if !ok {
		return nil, errors.New("boolean expected")
	}
	return enabled, nil
}
func (s *helicorderConfigEnabledImpl) Restore(handler *action.Handler) error {
	if err := handler.SettingsSet(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), s.GetDefaultValue()); err != nil {
		return fmt.Errorf("failed to reset helicorder service availablity: %w", err)
	}
	return nil
}

type helicorderConfigFilePathImpl struct{}

func (s *helicorderConfigFilePathImpl) GetName() string             { return "File Path" }
func (s *helicorderConfigFilePathImpl) GetNamespace() string        { return ID }
func (s *helicorderConfigFilePathImpl) GetKey() string              { return "file_path" }
func (s *helicorderConfigFilePathImpl) GetType() action.SettingType { return action.String }
func (s *helicorderConfigFilePathImpl) IsRequired() bool            { return true }
func (s *helicorderConfigFilePathImpl) GetVersion() int             { return 0 }
func (s *helicorderConfigFilePathImpl) GetOptions() map[string]any  { return nil }
func (s *helicorderConfigFilePathImpl) GetDefaultValue() any        { return "./service_data/helicorder" }
func (s *helicorderConfigFilePathImpl) GetDescription() string {
	return "The path to which helicorder images will be written, if the path does not exist, it will be automatically created."
}
func (s *helicorderConfigFilePathImpl) Init(handler *action.Handler) error {
	if _, err := handler.SettingsInit(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), s.GetDefaultValue()); err != nil {
		return fmt.Errorf("failed to set default helicorder service availablity: %w", err)
	}
	return nil
}
func (s *helicorderConfigFilePathImpl) Set(handler *action.Handler, newVal any) error {
	filePath, err := config.GetConfigValString(newVal)
	if err != nil {
		return err
	}
	filePath = path.Clean(filePath)
	if err := handler.SettingsSet(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), filePath); err != nil {
		return fmt.Errorf("failed to set helicorder service availablity: %w", err)
	}
	return nil
}
func (s *helicorderConfigFilePathImpl) Get(handler *action.Handler) (any, error) {
	val, _, _, err := handler.SettingsGet(s.GetNamespace(), s.GetKey())
	if err != nil {
		return nil, fmt.Errorf("failed to get helicorder service availablity: %w", err)
	}
	filePath, ok := val.(string)
	if !ok {
		return nil, errors.New("string expected")
	}
	return filePath, nil
}
func (s *helicorderConfigFilePathImpl) Restore(handler *action.Handler) error {
	if err := handler.SettingsSet(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), s.GetDefaultValue()); err != nil {
		return fmt.Errorf("failed to reset helicorder service availablity: %w", err)
	}
	return nil
}

type helicorderConfigImageFormatImpl struct{}

func (s *helicorderConfigImageFormatImpl) GetName() string             { return "Image Format" }
func (s *helicorderConfigImageFormatImpl) GetNamespace() string        { return ID }
func (s *helicorderConfigImageFormatImpl) GetKey() string              { return "image_format" }
func (s *helicorderConfigImageFormatImpl) GetType() action.SettingType { return action.String }
func (s *helicorderConfigImageFormatImpl) IsRequired() bool            { return true }
func (s *helicorderConfigImageFormatImpl) GetVersion() int             { return 0 }
func (s *helicorderConfigImageFormatImpl) GetOptions() map[string]any {
	return map[string]any{
		"PNG": IMAGE_FORMAT_PNG,
		"SVG": IMAGE_FORMAT_SVG,
	}
}
func (s *helicorderConfigImageFormatImpl) GetDefaultValue() any { return IMAGE_FORMAT_SVG }
func (s *helicorderConfigImageFormatImpl) GetDescription() string {
	return "Specify image format for helicorder images, supported formats are PNG and SVG, by default SVG is used."
}
func (s *helicorderConfigImageFormatImpl) Init(handler *action.Handler) error {
	if _, err := handler.SettingsInit(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), s.GetDefaultValue()); err != nil {
		return fmt.Errorf("failed to set default image format for helicorder service: %w", err)
	}
	return nil
}
func (s *helicorderConfigImageFormatImpl) Set(handler *action.Handler, newVal any) error {
	imageFormat, err := config.GetConfigValString(newVal)
	if err != nil {
		return err
	}
	if !lo.Contains([]string{
		IMAGE_FORMAT_PNG,
		IMAGE_FORMAT_SVG,
	}, imageFormat) {
		return errors.New("image format must be one of the given options")
	}
	if err := handler.SettingsSet(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), imageFormat); err != nil {
		return fmt.Errorf("failed to set helicorder image format: %w", err)
	}
	return nil
}
func (s *helicorderConfigImageFormatImpl) Get(handler *action.Handler) (any, error) {
	val, _, _, err := handler.SettingsGet(s.GetNamespace(), s.GetKey())
	if err != nil {
		return nil, fmt.Errorf("failed to get helicorder image format: %w", err)
	}
	imageFormat, ok := val.(string)
	if !ok {
		return nil, errors.New("string expected")
	}
	return imageFormat, nil
}
func (s *helicorderConfigImageFormatImpl) Restore(handler *action.Handler) error {
	if err := handler.SettingsSet(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), s.GetDefaultValue()); err != nil {
		return fmt.Errorf("failed to reset helicorder image format: %w", err)
	}
	return nil
}

type helicorderConfigTimeSpanImpl struct{}

func (s *helicorderConfigTimeSpanImpl) GetName() string             { return "Time Span" }
func (s *helicorderConfigTimeSpanImpl) GetNamespace() string        { return ID }
func (s *helicorderConfigTimeSpanImpl) GetKey() string              { return "time_span" }
func (s *helicorderConfigTimeSpanImpl) GetType() action.SettingType { return action.Int }
func (s *helicorderConfigTimeSpanImpl) IsRequired() bool            { return true }
func (s *helicorderConfigTimeSpanImpl) GetVersion() int             { return 0 }
func (s *helicorderConfigTimeSpanImpl) GetOptions() map[string]any {
	return map[string]any{
		"10 min": TIMESPAN_10_MINUTES,
		"15 min": TIMESPAN_15_MINUTES,
		"30 min": TIMESPAN_30_MINUTES,
		"60 min": TIMESPAN_60_MINUTES,
	}
}
func (s *helicorderConfigTimeSpanImpl) GetDefaultValue() any { return TIMESPAN_15_MINUTES }
func (s *helicorderConfigTimeSpanImpl) GetDescription() string {
	return "Time span defines the duration of the horizontal axis of the halicorder chart."
}
func (s *helicorderConfigTimeSpanImpl) Init(handler *action.Handler) error {
	if _, err := handler.SettingsInit(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), s.GetDefaultValue()); err != nil {
		return fmt.Errorf("failed to set default time span for helicorder service: %w", err)
	}
	return nil
}
func (s *helicorderConfigTimeSpanImpl) Set(handler *action.Handler, newVal any) error {
	timeSpan, err := config.GetConfigValInt64(newVal)
	if err != nil {
		return err
	}
	if !lo.Contains([]int64{
		TIMESPAN_10_MINUTES,
		TIMESPAN_15_MINUTES,
		TIMESPAN_30_MINUTES,
		TIMESPAN_60_MINUTES,
	}, timeSpan) {
		return errors.New("time span must be one of the given options")
	}
	if err := handler.SettingsSet(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), timeSpan); err != nil {
		return fmt.Errorf("failed to set helicorder time span: %w", err)
	}
	return nil
}
func (s *helicorderConfigTimeSpanImpl) Get(handler *action.Handler) (any, error) {
	val, _, _, err := handler.SettingsGet(s.GetNamespace(), s.GetKey())
	if err != nil {
		return nil, fmt.Errorf("failed to get helicorder time span: %w", err)
	}
	path, ok := val.(int64)
	if !ok {
		return nil, errors.New("integer expected")
	}
	return int(path), nil
}
func (s *helicorderConfigTimeSpanImpl) Restore(handler *action.Handler) error {
	if err := handler.SettingsSet(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), s.GetDefaultValue()); err != nil {
		return fmt.Errorf("failed to reset helicorder time span: %w", err)
	}
	return nil
}

type helicorderConfigLifeCycleImpl struct{}

func (s *helicorderConfigLifeCycleImpl) GetName() string             { return "Life Cycle" }
func (s *helicorderConfigLifeCycleImpl) GetNamespace() string        { return ID }
func (s *helicorderConfigLifeCycleImpl) GetKey() string              { return "life_cycle" }
func (s *helicorderConfigLifeCycleImpl) GetType() action.SettingType { return action.Int }
func (s *helicorderConfigLifeCycleImpl) IsRequired() bool            { return true }
func (s *helicorderConfigLifeCycleImpl) GetVersion() int             { return 0 }
func (s *helicorderConfigLifeCycleImpl) GetOptions() map[string]any  { return nil }
func (s *helicorderConfigLifeCycleImpl) GetDefaultValue() any        { return 0 }
func (s *helicorderConfigLifeCycleImpl) GetDescription() string {
	return "The number of days after which helicorder images will be automatically purged, set to 0 to keep files indefinitely."
}
func (s *helicorderConfigLifeCycleImpl) Init(handler *action.Handler) error {
	if _, err := handler.SettingsInit(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), s.GetDefaultValue()); err != nil {
		return fmt.Errorf("failed to set default life cycle for helicorder service: %w", err)
	}
	return nil
}
func (s *helicorderConfigLifeCycleImpl) Set(handler *action.Handler, newVal any) error {
	lifeCycle, err := config.GetConfigValInt64(newVal)
	if err != nil {
		return err
	}
	if lifeCycle < 0 {
		return errors.New("life cycle cannot be negative")
	}
	if err := handler.SettingsSet(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), lifeCycle); err != nil {
		return fmt.Errorf("failed to set helicorder life cycle: %w", err)
	}
	return nil
}
func (s *helicorderConfigLifeCycleImpl) Get(handler *action.Handler) (any, error) {
	val, _, _, err := handler.SettingsGet(s.GetNamespace(), s.GetKey())
	if err != nil {
		return nil, fmt.Errorf("failed to get helicorder life cycle: %w", err)
	}
	lifeCycle, ok := val.(int64)
	if !ok {
		return nil, errors.New("integer expected")
	}
	return int(lifeCycle), nil
}
func (s *helicorderConfigLifeCycleImpl) Restore(handler *action.Handler) error {
	if err := handler.SettingsSet(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), s.GetDefaultValue()); err != nil {
		return fmt.Errorf("failed to reset helicorder life cycle: %w", err)
	}
	return nil
}

type helicorderConfigImageSizeImpl struct{}

func (s *helicorderConfigImageSizeImpl) GetName() string             { return "Image Size" }
func (s *helicorderConfigImageSizeImpl) GetNamespace() string        { return ID }
func (s *helicorderConfigImageSizeImpl) GetKey() string              { return "image_size" }
func (s *helicorderConfigImageSizeImpl) GetType() action.SettingType { return action.Int }
func (s *helicorderConfigImageSizeImpl) IsRequired() bool            { return true }
func (s *helicorderConfigImageSizeImpl) GetVersion() int             { return 0 }
func (s *helicorderConfigImageSizeImpl) GetOptions() map[string]any  { return nil }
func (s *helicorderConfigImageSizeImpl) GetDefaultValue() any        { return 800 }
func (s *helicorderConfigImageSizeImpl) GetDescription() string {
	return "Image size in pixels for helicorder images."
}
func (s *helicorderConfigImageSizeImpl) Init(handler *action.Handler) error {
	if _, err := handler.SettingsInit(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), s.GetDefaultValue()); err != nil {
		return fmt.Errorf("failed to set default image size for helicorder service: %w", err)
	}
	return nil
}
func (s *helicorderConfigImageSizeImpl) Set(handler *action.Handler, newVal any) error {
	imageSize, err := config.GetConfigValInt64(newVal)
	if err != nil {
		return err
	}
	if imageSize <= 0 {
		return errors.New("image size cannot be zero or negative")
	}
	if imageSize > 4096 {
		return errors.New("image size cannot be greater than 4096")
	}
	if err := handler.SettingsSet(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), imageSize); err != nil {
		return fmt.Errorf("failed to set helicorder image size: %w", err)
	}
	return nil
}
func (s *helicorderConfigImageSizeImpl) Get(handler *action.Handler) (any, error) {
	val, _, _, err := handler.SettingsGet(s.GetNamespace(), s.GetKey())
	if err != nil {
		return nil, fmt.Errorf("failed to get helicorder image size: %w", err)
	}
	imageSize, ok := val.(int64)
	if !ok {
		return nil, errors.New("integer expected")
	}
	return int(imageSize), nil
}
func (s *helicorderConfigImageSizeImpl) Restore(handler *action.Handler) error {
	if err := handler.SettingsSet(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), s.GetDefaultValue()); err != nil {
		return fmt.Errorf("failed to reset helicorder image size: %w", err)
	}
	return nil
}

type helicorderConfigSpanSamplesImpl struct{}

func (s *helicorderConfigSpanSamplesImpl) GetName() string             { return "Span Samples" }
func (s *helicorderConfigSpanSamplesImpl) GetNamespace() string        { return ID }
func (s *helicorderConfigSpanSamplesImpl) GetKey() string              { return "span_samples" }
func (s *helicorderConfigSpanSamplesImpl) GetType() action.SettingType { return action.Int }
func (s *helicorderConfigSpanSamplesImpl) IsRequired() bool            { return true }
func (s *helicorderConfigSpanSamplesImpl) GetVersion() int             { return 0 }
func (s *helicorderConfigSpanSamplesImpl) GetOptions() map[string]any  { return nil }
func (s *helicorderConfigSpanSamplesImpl) GetDefaultValue() any        { return 10000 }
func (s *helicorderConfigSpanSamplesImpl) GetDescription() string {
	return "Samples per time span in points for helicorder images."
}
func (s *helicorderConfigSpanSamplesImpl) Init(handler *action.Handler) error {
	if _, err := handler.SettingsInit(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), s.GetDefaultValue()); err != nil {
		return fmt.Errorf("failed to set default image size for helicorder service: %w", err)
	}
	return nil
}
func (s *helicorderConfigSpanSamplesImpl) Set(handler *action.Handler, newVal any) error {
	spanSamples, err := config.GetConfigValInt64(newVal)
	if err != nil {
		return err
	}
	if spanSamples <= 0 {
		return errors.New("span samples cannot be zero or negative")
	}
	if spanSamples > 100000 {
		return errors.New("span samples cannot be greater than 100000")
	}
	if err := handler.SettingsSet(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), spanSamples); err != nil {
		return fmt.Errorf("failed to set helicorder span samples: %w", err)
	}
	return nil
}
func (s *helicorderConfigSpanSamplesImpl) Get(handler *action.Handler) (any, error) {
	val, _, _, err := handler.SettingsGet(s.GetNamespace(), s.GetKey())
	if err != nil {
		return nil, fmt.Errorf("failed to get helicorder span samples: %w", err)
	}
	spanSamples, ok := val.(int64)
	if !ok {
		return nil, errors.New("integer expected")
	}
	return int(spanSamples), nil
}
func (s *helicorderConfigSpanSamplesImpl) Restore(handler *action.Handler) error {
	if err := handler.SettingsSet(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), s.GetDefaultValue()); err != nil {
		return fmt.Errorf("failed to reset helicorder span samples: %w", err)
	}
	return nil
}

type helicorderConfigLineWidthImpl struct{}

func (s *helicorderConfigLineWidthImpl) GetName() string             { return "Line Width" }
func (s *helicorderConfigLineWidthImpl) GetNamespace() string        { return ID }
func (s *helicorderConfigLineWidthImpl) GetKey() string              { return "line_width" }
func (s *helicorderConfigLineWidthImpl) GetType() action.SettingType { return action.Float }
func (s *helicorderConfigLineWidthImpl) IsRequired() bool            { return true }
func (s *helicorderConfigLineWidthImpl) GetVersion() int             { return 0 }
func (s *helicorderConfigLineWidthImpl) GetOptions() map[string]any  { return nil }
func (s *helicorderConfigLineWidthImpl) GetDefaultValue() any        { return 1.0 }
func (s *helicorderConfigLineWidthImpl) GetDescription() string {
	return "Samples per time span in points for helicorder images."
}
func (s *helicorderConfigLineWidthImpl) Init(handler *action.Handler) error {
	if _, err := handler.SettingsInit(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), s.GetDefaultValue()); err != nil {
		return fmt.Errorf("failed to set default plot line width for helicorder service: %w", err)
	}
	return nil
}
func (s *helicorderConfigLineWidthImpl) Set(handler *action.Handler, newVal any) error {
	lineWidth, err := config.GetConfigValFloat64(newVal)
	if err != nil {
		return err
	}
	if lineWidth <= 0 {
		return errors.New("line width cannot be zero or negative")
	}
	if lineWidth > 5 {
		return errors.New("line width cannot be greater than 5")
	}
	if err := handler.SettingsSet(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), lineWidth); err != nil {
		return fmt.Errorf("failed to set helicorder plot line width: %w", err)
	}
	return nil
}
func (s *helicorderConfigLineWidthImpl) Get(handler *action.Handler) (any, error) {
	val, _, _, err := handler.SettingsGet(s.GetNamespace(), s.GetKey())
	if err != nil {
		return nil, fmt.Errorf("failed to get helicorder plot line width: %w", err)
	}
	spanSamples, ok := val.(float64)
	if !ok {
		return nil, errors.New("float expected")
	}
	return spanSamples, nil
}
func (s *helicorderConfigLineWidthImpl) Restore(handler *action.Handler) error {
	if err := handler.SettingsSet(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), s.GetDefaultValue()); err != nil {
		return fmt.Errorf("failed to reset helicorder plot line width: %w", err)
	}
	return nil
}

type helicorderConfigScaleFactorImpl struct{}

func (s *helicorderConfigScaleFactorImpl) GetName() string             { return "Sacle Factor" }
func (s *helicorderConfigScaleFactorImpl) GetNamespace() string        { return ID }
func (s *helicorderConfigScaleFactorImpl) GetKey() string              { return "sacle_factor" }
func (s *helicorderConfigScaleFactorImpl) GetType() action.SettingType { return action.Float }
func (s *helicorderConfigScaleFactorImpl) IsRequired() bool            { return true }
func (s *helicorderConfigScaleFactorImpl) GetVersion() int             { return 0 }
func (s *helicorderConfigScaleFactorImpl) GetOptions() map[string]any  { return nil }
func (s *helicorderConfigScaleFactorImpl) GetDefaultValue() any        { return 2.3 }
func (s *helicorderConfigScaleFactorImpl) GetDescription() string {
	return "Waveform scale factor in helicorder."
}
func (s *helicorderConfigScaleFactorImpl) Init(handler *action.Handler) error {
	if _, err := handler.SettingsInit(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), s.GetDefaultValue()); err != nil {
		return fmt.Errorf("failed to set default waveform scale factor for helicorder service: %w", err)
	}
	return nil
}
func (s *helicorderConfigScaleFactorImpl) Set(handler *action.Handler, newVal any) error {
	scaleFactor, err := config.GetConfigValFloat64(newVal)
	if err != nil {
		return err
	}
	if scaleFactor <= 0 {
		return errors.New("scale factor cannot be zero or negative")
	}
	if scaleFactor > 5 {
		return errors.New("scale factor cannot be greater than 5")
	}
	if err := handler.SettingsSet(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), scaleFactor); err != nil {
		return fmt.Errorf("failed to set helicorder waveform scale factor: %w", err)
	}
	return nil
}
func (s *helicorderConfigScaleFactorImpl) Get(handler *action.Handler) (any, error) {
	val, _, _, err := handler.SettingsGet(s.GetNamespace(), s.GetKey())
	if err != nil {
		return nil, fmt.Errorf("failed to get helicorder waveform scale factor: %w", err)
	}
	scaleFactor, ok := val.(float64)
	if !ok {
		return nil, errors.New("float expected")
	}
	return scaleFactor, nil
}
func (s *helicorderConfigScaleFactorImpl) Restore(handler *action.Handler) error {
	if err := handler.SettingsSet(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), s.GetDefaultValue()); err != nil {
		return fmt.Errorf("failed to reset helicorder waveform scale factor: %w", err)
	}
	return nil
}

func (s *HelicorderServiceImpl) GetConfigConstraint() []config.IConstraint {
	return []config.IConstraint{
		&helicorderConfigEnabledImpl{},
		&helicorderConfigFilePathImpl{},
		&helicorderConfigImageFormatImpl{},
		&helicorderConfigTimeSpanImpl{},
		&helicorderConfigLifeCycleImpl{},
		&helicorderConfigImageSizeImpl{},
		&helicorderConfigSpanSamplesImpl{},
		&helicorderConfigLineWidthImpl{},
		&helicorderConfigScaleFactorImpl{},
	}
}
