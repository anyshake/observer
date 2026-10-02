package tailscale

import (
	"fmt"

	"github.com/anyshake/observer/internal/service"
)

func (s *TailscaleServiceImpl) GetAssetList() ([]service.Asset, error) {
	return nil, fmt.Errorf("assets are not available on %s", ID)
}

func (s *TailscaleServiceImpl) GetAssetData(assetID string) (*service.AssetData, error) {
	return nil, fmt.Errorf("asset ID %s is not available on %s", assetID, ID)
}
