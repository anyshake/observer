package miniseed

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/anyshake/observer/internal/service"
)

func (s *MiniSeedServiceImpl) safeFileAccess(assetId string) (string, error) {
	absFilePath, err := filepath.Abs(s.filePath)
	if err != nil {
		return "", err
	}

	absAssetIdPath, err := filepath.Abs(assetId)
	if err != nil {
		return "", err
	}
	absAssetIdPath = filepath.Clean(absAssetIdPath)

	if !strings.HasPrefix(absAssetIdPath, absFilePath) {
		return "", fmt.Errorf("asset %s is not available on %s", assetId, ID)
	}
	if !strings.HasSuffix(absAssetIdPath, ".mseed") {
		return "", fmt.Errorf("asset %s is not available on %s", assetId, ID)
	}

	return absAssetIdPath, nil
}

func (s *MiniSeedServiceImpl) GetAssetList() ([]service.Asset, error) {
	if !s.status.GetIsRunning() {
		return nil, nil
	}

	var assets []service.Asset
	_ = filepath.Walk(s.filePath, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if !info.IsDir() && strings.HasSuffix(info.Name(), ".mseed") {
			assets = append(assets, service.Asset{
				FilePath:   path,
				FileName:   info.Name(),
				Size:       info.Size(),
				ModifiedAt: info.ModTime().UnixMilli(),
			})
		}

		return nil
	})

	return assets, nil
}

func (s *MiniSeedServiceImpl) GetAssetData(assetId string) (*service.AssetData, error) {
	if !s.status.GetIsRunning() {
		return nil, fmt.Errorf("assets ID %s is not available on %s when service is not running", assetId, ID)
	}

	absPath, err := s.safeFileAccess(assetId)
	if err != nil {
		return nil, fmt.Errorf("failed to get safe file access for asset %s: %v", assetId, err)
	}

	if _, err := os.Stat(absPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("asset file %s does not exist", assetId)
	}

	data, err := os.ReadFile(absPath)
	if err != nil {
		return nil, fmt.Errorf("failed to read asset data from file %s: %v", assetId, err)
	}

	return &service.AssetData{
		ContentType: "application/octet-stream",
		FileName:    filepath.Base(assetId),
		Data:        data,
	}, nil
}
