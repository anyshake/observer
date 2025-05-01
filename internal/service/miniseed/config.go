package miniseed

import (
	"errors"
	"fmt"
	"path"

	"github.com/anyshake/observer/config"
	"github.com/anyshake/observer/internal/dao/action"
)

type miniSeedConfigEnabledImpl struct{}

func (s *miniSeedConfigEnabledImpl) GetName() string             { return "Enable" }
func (s *miniSeedConfigEnabledImpl) GetNamespace() string        { return ID }
func (s *miniSeedConfigEnabledImpl) GetKey() string              { return "enabled" }
func (s *miniSeedConfigEnabledImpl) GetType() action.SettingType { return action.Bool }
func (s *miniSeedConfigEnabledImpl) IsRequired() bool            { return true }
func (s *miniSeedConfigEnabledImpl) GetVersion() int             { return 0 }
func (s *miniSeedConfigEnabledImpl) GetOptions() map[string]any  { return nil }
func (s *miniSeedConfigEnabledImpl) GetDefaultValue() any        { return false }
func (s *miniSeedConfigEnabledImpl) GetDescription() string {
	return "Enable MiniSEED service to save daily split records per channel in the specified path, creating it if the path does not exist."
}
func (s *miniSeedConfigEnabledImpl) Init(handler *action.Handler) error {
	if _, err := handler.SettingsInit(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), s.GetDefaultValue()); err != nil {
		return fmt.Errorf("failed to set default MiniSEED service availablity: %w", err)
	}
	return nil
}
func (s *miniSeedConfigEnabledImpl) Set(handler *action.Handler, newVal any) error {
	enabled, err := config.GetConfigValBool(newVal)
	if err != nil {
		return err
	}
	if err := handler.SettingsSet(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), enabled); err != nil {
		return fmt.Errorf("failed to set MiniSEED service availablity: %w", err)
	}
	return nil
}
func (s *miniSeedConfigEnabledImpl) Get(handler *action.Handler) (any, error) {
	val, _, _, err := handler.SettingsGet(s.GetNamespace(), s.GetKey())
	if err != nil {
		return nil, fmt.Errorf("failed to get MiniSEED service availablity: %w", err)
	}
	enabled, ok := val.(bool)
	if !ok {
		return nil, errors.New("boolean expected")
	}
	return enabled, nil
}
func (s *miniSeedConfigEnabledImpl) Restore(handler *action.Handler) error {
	if err := handler.SettingsSet(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), s.GetDefaultValue()); err != nil {
		return fmt.Errorf("failed to reset MiniSEED service availablity: %w", err)
	}
	return nil
}

type miniSeedConfigUseCompressImpl struct{}

func (s *miniSeedConfigUseCompressImpl) GetName() string             { return "Use Compression" }
func (s *miniSeedConfigUseCompressImpl) GetNamespace() string        { return ID }
func (s *miniSeedConfigUseCompressImpl) GetKey() string              { return "use_compression" }
func (s *miniSeedConfigUseCompressImpl) GetType() action.SettingType { return action.Bool }
func (s *miniSeedConfigUseCompressImpl) IsRequired() bool            { return false }
func (s *miniSeedConfigUseCompressImpl) GetVersion() int             { return 0 }
func (s *miniSeedConfigUseCompressImpl) GetOptions() map[string]any  { return nil }
func (s *miniSeedConfigUseCompressImpl) GetDefaultValue() any        { return false }
func (s *miniSeedConfigUseCompressImpl) GetDescription() string {
	return "Whether to compress MiniSEED files. If true, files will use STEIM-2 compression, encoding the count data as 30-bit instead of 32-bit. This may cause data overflow and checksum errors if a 32-bit ADC reaches full scale."
}
func (s *miniSeedConfigUseCompressImpl) Init(handler *action.Handler) error {
	if _, err := handler.SettingsInit(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), s.GetDefaultValue()); err != nil {
		return fmt.Errorf("failed to set default MiniSEED compression flag: %w", err)
	}
	return nil
}
func (s *miniSeedConfigUseCompressImpl) Set(handler *action.Handler, newVal any) error {
	useCompress, err := config.GetConfigValBool(newVal)
	if err != nil {
		return err
	}
	if err := handler.SettingsSet(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), useCompress); err != nil {
		return fmt.Errorf("failed to set MiniSEED compression flag: %w", err)
	}
	return nil
}
func (s *miniSeedConfigUseCompressImpl) Get(handler *action.Handler) (any, error) {
	val, _, _, err := handler.SettingsGet(s.GetNamespace(), s.GetKey())
	if err != nil {
		return nil, fmt.Errorf("failed to get MiniSEED compression flag: %w", err)
	}
	useCompress, ok := val.(bool)
	if !ok {
		return nil, errors.New("boolean expected")
	}
	return useCompress, nil
}
func (s *miniSeedConfigUseCompressImpl) Restore(handler *action.Handler) error {
	if err := handler.SettingsSet(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), s.GetDefaultValue()); err != nil {
		return fmt.Errorf("failed to reset MiniSEED compression flag: %w", err)
	}
	return nil
}

type miniSeedConfigFilePathImpl struct{}

func (s *miniSeedConfigFilePathImpl) GetName() string             { return "File Path" }
func (s *miniSeedConfigFilePathImpl) GetNamespace() string        { return ID }
func (s *miniSeedConfigFilePathImpl) GetKey() string              { return "file_path" }
func (s *miniSeedConfigFilePathImpl) GetType() action.SettingType { return action.String }
func (s *miniSeedConfigFilePathImpl) IsRequired() bool            { return true }
func (s *miniSeedConfigFilePathImpl) GetVersion() int             { return 0 }
func (s *miniSeedConfigFilePathImpl) GetOptions() map[string]any  { return nil }
func (s *miniSeedConfigFilePathImpl) GetDefaultValue() any        { return "./service_data/miniseed" }
func (s *miniSeedConfigFilePathImpl) GetDescription() string {
	return "The path to which MiniSEED files will be written, if the path does not exist, it will be automatically created."
}
func (s *miniSeedConfigFilePathImpl) Init(handler *action.Handler) error {
	if _, err := handler.SettingsInit(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), s.GetDefaultValue()); err != nil {
		return fmt.Errorf("failed to set default storage path for MiniSEED service: %w", err)
	}
	return nil
}
func (s *miniSeedConfigFilePathImpl) Set(handler *action.Handler, newVal any) error {
	filePath, err := config.GetConfigValString(newVal)
	if err != nil {
		return err
	}
	filePath = path.Clean(filePath)
	if err := handler.SettingsSet(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), filePath); err != nil {
		return fmt.Errorf("failed to set MiniSEED storage path: %w", err)
	}
	return nil
}
func (s *miniSeedConfigFilePathImpl) Get(handler *action.Handler) (any, error) {
	val, _, _, err := handler.SettingsGet(s.GetNamespace(), s.GetKey())
	if err != nil {
		return nil, fmt.Errorf("failed to get MiniSEED storage path: %w", err)
	}
	filePath, ok := val.(string)
	if !ok {
		return nil, errors.New("string expected")
	}
	return filePath, nil
}
func (s *miniSeedConfigFilePathImpl) Restore(handler *action.Handler) error {
	if err := handler.SettingsSet(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), s.GetDefaultValue()); err != nil {
		return fmt.Errorf("failed to reset MiniSEED storage path: %w", err)
	}
	return nil
}

type miniSeedConfigLifeCycleImpl struct{}

func (s *miniSeedConfigLifeCycleImpl) GetName() string             { return "Life Cycle" }
func (s *miniSeedConfigLifeCycleImpl) GetNamespace() string        { return ID }
func (s *miniSeedConfigLifeCycleImpl) GetKey() string              { return "life_cycle" }
func (s *miniSeedConfigLifeCycleImpl) GetType() action.SettingType { return action.Int }
func (s *miniSeedConfigLifeCycleImpl) IsRequired() bool            { return true }
func (s *miniSeedConfigLifeCycleImpl) GetVersion() int             { return 0 }
func (s *miniSeedConfigLifeCycleImpl) GetOptions() map[string]any  { return nil }
func (s *miniSeedConfigLifeCycleImpl) GetDefaultValue() any        { return 0 }
func (s *miniSeedConfigLifeCycleImpl) GetDescription() string {
	return "The number of days after which MiniSEED files will be automatically purged, set to 0 to keep files indefinitely."
}
func (s *miniSeedConfigLifeCycleImpl) Init(handler *action.Handler) error {
	if _, err := handler.SettingsInit(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), s.GetDefaultValue()); err != nil {
		return fmt.Errorf("failed to set default life cycle for MiniSEED service: %w", err)
	}
	return nil
}
func (s *miniSeedConfigLifeCycleImpl) Set(handler *action.Handler, newVal any) error {
	lifeCycle, err := config.GetConfigValInt64(newVal)
	if err != nil {
		return err
	}
	if lifeCycle < 0 {
		return errors.New("life cycle cannot be negative")
	}
	if err := handler.SettingsSet(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), lifeCycle); err != nil {
		return fmt.Errorf("failed to set MiniSEED life cycle: %w", err)
	}
	return nil
}
func (s *miniSeedConfigLifeCycleImpl) Get(handler *action.Handler) (any, error) {
	val, _, _, err := handler.SettingsGet(s.GetNamespace(), s.GetKey())
	if err != nil {
		return nil, fmt.Errorf("failed to get MiniSEED life cycle: %w", err)
	}
	lifeCycle, ok := val.(int64)
	if !ok {
		return nil, errors.New("integer expected")
	}
	return int(lifeCycle), nil
}
func (s *miniSeedConfigLifeCycleImpl) Restore(handler *action.Handler) error {
	if err := handler.SettingsSet(s.GetNamespace(), s.GetKey(), s.GetType(), s.GetVersion(), s.GetDefaultValue()); err != nil {
		return fmt.Errorf("failed to reset MiniSEED life cycle: %w", err)
	}
	return nil
}

func (s *MiniSeedServiceImpl) GetConfigConstraint() []config.IConstraint {
	return []config.IConstraint{
		&miniSeedConfigEnabledImpl{},
		&miniSeedConfigUseCompressImpl{},
		&miniSeedConfigLifeCycleImpl{},
		&miniSeedConfigFilePathImpl{},
	}
}
