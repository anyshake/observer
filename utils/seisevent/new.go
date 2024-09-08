package seisevent

import (
	"time"

	"github.com/anyshake/observer/utils/cache"
)

func New(cacheTTL time.Duration) map[string]DataSource {
	return map[string]DataSource{
		CWA_ID:      &CWA{cache: cache.New(cacheTTL)},
		HKO_ID:      &HKO{cache: cache.New(cacheTTL)},
		JMA_ID:      &JMA{cache: cache.New(cacheTTL)},
		KMA_ID:      &KMA{cache: cache.New(cacheTTL)},
		CENC_WEB_ID: &CENC_WEB{cache: cache.New(cacheTTL)},
		CENC_APP_ID: &CENC_APP{cache: cache.New(cacheTTL)},
		USGS_ID:     &USGS{cache: cache.New(cacheTTL)},
		INGV_ID:     &INGV{cache: cache.New(cacheTTL)},
		SCEA_E_ID:   &SCEA_E{cache: cache.New(cacheTTL)},
		SCEA_B_ID:   &SCEA_B{cache: cache.New(cacheTTL)},
		CEA_ID:      &CEA{cache: cache.New(cacheTTL)},
		EQZT_ID:     &EQZT{cache: cache.New(cacheTTL)},
		EMSC_ID:     &EMSC{cache: cache.New(cacheTTL)},
		GFZ_ID:      &GFZ{cache: cache.New(cacheTTL)},
		USP_ID:      &USP{cache: cache.New(cacheTTL)},
		GA_ID:       &GA{cache: cache.New(cacheTTL)},
		AUSPASS_ID:  &AUSPASS{cache: cache.New(cacheTTL)},
		BCSF_ID:     &BCSF{cache: cache.New(cacheTTL)},
		INFP_ID:     &INFP{cache: cache.New(cacheTTL)},
		SED_ID:      &SED{cache: cache.New(cacheTTL)},
		KNMI_ID:     &KNMI{cache: cache.New(cacheTTL)},
		NCS_ID:      &NCS{cache: cache.New(cacheTTL)},
		NRCAN_ID:    &NRCAN{cache: cache.New(cacheTTL)},
		GEONET_ID:   &GEONET{cache: cache.New(cacheTTL)},
		FJEA_ID:     &FJEA{cache: cache.New(cacheTTL)},
		ICL_ID:      &ICL{cache: cache.New(cacheTTL)},
	}
}
