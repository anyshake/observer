package miniseed

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"

	"github.com/anyshake/observer/internal/service"
)

func (s *MiniSeedServiceImpl) getFilePath() string {
	s.mu.Lock()
	defer s.mu.Unlock()

	return s.filePath
}

func (s *MiniSeedServiceImpl) safeFileAccess(baseDir, assetId string) (string, error) {
	absFilePath, err := filepath.Abs(baseDir)
	if err != nil {
		return "", err
	}
	absFilePath = filepath.Clean(absFilePath)

	absAssetIdPath, err := filepath.Abs(assetId)
	if err != nil {
		return "", err
	}
	absAssetIdPath = filepath.Clean(absAssetIdPath)

	rel, err := filepath.Rel(absFilePath, absAssetIdPath)
	if err != nil ||
		rel == ".." ||
		strings.HasPrefix(rel, ".."+string(os.PathSeparator)) {
		return "", fmt.Errorf("asset %s is not available on %s", assetId, ID)
	}

	if !strings.EqualFold(filepath.Ext(absAssetIdPath), ".mseed") {
		return "", fmt.Errorf("asset %s is not available on %s", assetId, ID)
	}

	realFilePath, err := filepath.EvalSymlinks(absAssetIdPath)
	if err != nil {
		return "", err
	}
	realFilePath = filepath.Clean(realFilePath)

	rel, err = filepath.Rel(absFilePath, realFilePath)
	if err != nil ||
		rel == ".." ||
		strings.HasPrefix(rel, ".."+string(os.PathSeparator)) {
		return "", fmt.Errorf("asset %s is not available on %s", assetId, ID)
	}

	return realFilePath, nil
}

func (s *MiniSeedServiceImpl) GetAssetList() ([]service.Asset, error) {
	if !s.status.GetIsRunning() {
		return nil, nil
	}

	var assets []service.Asset
	_ = filepath.Walk(s.getFilePath(), func(path string, info os.FileInfo, err error) error {
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

	absPath, err := s.safeFileAccess(s.getFilePath(), assetId)
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
