package seisevent

import (
	"time"

	"github.com/anyshake/observer/utils/cache"
)

const (
	CWA_ID      = "cwa"
	HKO_ID      = "hko"
	JMA_ID      = "jma"
	KMA_ID      = "kma"
	CEIC_ID     = "ceic"
	USGS_ID     = "usgs"
	INGV_ID     = "ingv"
	SCEA_E_ID   = "scea_e"
	SCEA_B_ID   = "scea_b"
	CEA_DASE_ID = "cea_dase"
	EQZT_ID     = "eqzt"
	EMSC_ID     = "emsc"
	GFZ_ID      = "gfz"
)

func New(cacheTTL time.Duration) map[string]DataSource {
	return map[string]DataSource{
		CWA_ID:      &CWA{cache: cache.New(cacheTTL)},
		HKO_ID:      &HKO{cache: cache.New(cacheTTL)},
		JMA_ID:      &JMA{cache: cache.New(cacheTTL)},
		KMA_ID:      &KMA{cache: cache.New(cacheTTL)},
		CEIC_ID:     &CEIC{cache: cache.New(cacheTTL)},
		USGS_ID:     &USGS{cache: cache.New(cacheTTL)},
		INGV_ID:     &INGV{cache: cache.New(cacheTTL)},
		SCEA_E_ID:   &SCEA_E{cache: cache.New(cacheTTL)},
		SCEA_B_ID:   &SCEA_B{cache: cache.New(cacheTTL)},
		CEA_DASE_ID: &CEA_DASE{cache: cache.New(cacheTTL)},
		EQZT_ID:     &EQZT{cache: cache.New(cacheTTL)},
		EMSC_ID:     &EMSC{cache: cache.New(cacheTTL)},
		GFZ_ID:      &GFZ{cache: cache.New(cacheTTL)},
	}
}
