package quakesense

import (
	"fmt"

	"github.com/anyshake/observer/internal/service"
)

func (s *QuakeSenseServiceImpl) GetAssetList() ([]service.Asset, error) {
	return nil, fmt.Errorf("assets are not available on %s", ID)
}

func (s *QuakeSenseServiceImpl) GetAssetData(assetId string) (*service.AssetData, error) {
	return nil, fmt.Errorf("assets ID %s is not available on %s", assetId, ID)
}
