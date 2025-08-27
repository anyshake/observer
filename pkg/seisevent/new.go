package seisevent

import (
	"time"

	"github.com/anyshake/observer/pkg/cache"
)

func New(cacheTTL time.Duration) map[string]IDataSource {
	return map[string]IDataSource{
		CWA_SC_ID:       &CWA_SC{cache: cache.New(cacheTTL)},
		CWA_EXPTECH_ID:  &CWA_EXPTECH{cache: cache.New(cacheTTL)},
		HKO_ID:          &HKO{cache: cache.New(cacheTTL)},
		JMA_OFFICIAL_ID: &JMA_OFFICIAL{cache: cache.New(cacheTTL)},
		JMA_P2PQUAKE_ID: &JMA_P2PQUAKE{cache: cache.New(cacheTTL)},
		JMA_WOLFX_ID:    &JMA_WOLFX{cache: cache.New(cacheTTL)},
		KMA_ID:          &KMA{cache: cache.New(cacheTTL)},
		CENC_WEB_ID:     &CENC_WEB{cache: cache.New(cacheTTL)},
		CENC_APP_ID:     &CENC_APP{cache: cache.New(cacheTTL)},
		CENC_WOLFX_ID:   &CENC_WOLFX{cache: cache.New(cacheTTL)},
		USGS_ID:         &USGS{cache: cache.New(cacheTTL)},
		INGV_ID:         &INGV{cache: cache.New(cacheTTL)},
		SCEA_E_ID:       &SCEA_E{cache: cache.New(cacheTTL)},
		SCEA_B_ID:       &SCEA_B{cache: cache.New(cacheTTL)},
		CEA_ID:          &CEA{cache: cache.New(cacheTTL)},
		// EQZT_ID:         &EQZT{cache: cache.New(cacheTTL)}, // broken
		EMSC_ID:    &EMSC{cache: cache.New(cacheTTL)},
		GFZ_ID:     &GFZ{cache: cache.New(cacheTTL)},
		USP_ID:     &USP{cache: cache.New(cacheTTL)},
		GA_ID:      &GA{cache: cache.New(cacheTTL)},
		AUSPASS_ID: &AUSPASS{cache: cache.New(cacheTTL)},
		BCSF_ID:    &BCSF{cache: cache.New(cacheTTL)},
		INFP_ID:    &INFP{cache: cache.New(cacheTTL)},
		SED_ID:     &SED{cache: cache.New(cacheTTL)},
		KNMI_ID:    &KNMI{cache: cache.New(cacheTTL)},
		NCS_ID:     &NCS{cache: cache.New(cacheTTL)},
		NRCAN_ID:   &NRCAN{cache: cache.New(cacheTTL)},
		GEONET_ID:  &GEONET{cache: cache.New(cacheTTL)},
		FJEA_ID:    &FJEA{cache: cache.New(cacheTTL)},
		ICL_ID:     &ICL{cache: cache.New(cacheTTL)},
		BMKG_ID:    &BMKG{cache: cache.New(cacheTTL)},
		DOST_ID:    &DOST{cache: cache.New(cacheTTL)},
		AFAD_ID:    &AFAD{cache: cache.New(cacheTTL)},
		KRDAE_ID:   &KRDAE{cache: cache.New(cacheTTL)},
		KNDC_ID:    &KNDC{cache: cache.New(cacheTTL)},
		TMD_ID:     &TMD{cache: cache.New(cacheTTL)},
		SSN_ID:     &SSN{cache: cache.New(cacheTTL)},
		BGS_ID:     &BGS{cache: cache.New(cacheTTL)},
	}
}
