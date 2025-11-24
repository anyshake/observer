package seisevent

import (
	"time"

	"github.com/anyshake/observer/pkg/cache"
	"github.com/anyshake/observer/pkg/dnsquery"
	"github.com/bclswl0827/travel"
)

func New(cacheTTL time.Duration) (map[string]IDataSource, error) {
	travelTimeTable, err := travel.NewAK135()
	if err != nil {
		return nil, err
	}

	builtinResolvers := dnsquery.NewResolvers()

	return map[string]IDataSource{
		CWA_SC_ID:       &CWA_SC{resolvers: builtinResolvers, cache: cache.New(cacheTTL), travelTimeTable: travelTimeTable},
		CWA_EXPTECH_ID:  &CWA_EXPTECH{cache: cache.New(cacheTTL), travelTimeTable: travelTimeTable},
		PALERT_ID:       &PALERT{cache: cache.New(cacheTTL), travelTimeTable: travelTimeTable},
		HKO_ID:          &HKO{cache: cache.New(cacheTTL), travelTimeTable: travelTimeTable},
		JMA_OFFICIAL_ID: &JMA_OFFICIAL{cache: cache.New(cacheTTL), travelTimeTable: travelTimeTable},
		JMA_P2PQUAKE_ID: &JMA_P2PQUAKE{cache: cache.New(cacheTTL), travelTimeTable: travelTimeTable},
		JMA_WOLFX_ID:    &JMA_WOLFX{cache: cache.New(cacheTTL), travelTimeTable: travelTimeTable},
		KMA_ID:          &KMA{cache: cache.New(cacheTTL), travelTimeTable: travelTimeTable},
		CENC_WEB_ID:     &CENC_WEB{cache: cache.New(cacheTTL), travelTimeTable: travelTimeTable},
		CENC_APP_ID:     &CENC_APP{cache: cache.New(cacheTTL), travelTimeTable: travelTimeTable},
		CENC_WOLFX_ID:   &CENC_WOLFX{cache: cache.New(cacheTTL), travelTimeTable: travelTimeTable},
		USGS_ID:         &USGS{cache: cache.New(cacheTTL), travelTimeTable: travelTimeTable},
		INGV_ID:         &INGV{cache: cache.New(cacheTTL), travelTimeTable: travelTimeTable},
		SCEA_E_ID:       &SCEA_E{cache: cache.New(cacheTTL), travelTimeTable: travelTimeTable},
		SCEA_B_ID:       &SCEA_B{cache: cache.New(cacheTTL), travelTimeTable: travelTimeTable},
		CEA_ID:          &CEA{cache: cache.New(cacheTTL), travelTimeTable: travelTimeTable},
		EMSC_ID:         &EMSC{cache: cache.New(cacheTTL), travelTimeTable: travelTimeTable},
		GFZ_ID:          &GFZ{cache: cache.New(cacheTTL), travelTimeTable: travelTimeTable},
		USP_ID:          &USP{cache: cache.New(cacheTTL), travelTimeTable: travelTimeTable},
		GA_ID:           &GA{cache: cache.New(cacheTTL), travelTimeTable: travelTimeTable},
		AUSPASS_ID:      &AUSPASS{cache: cache.New(cacheTTL), travelTimeTable: travelTimeTable},
		BCSF_ID:         &BCSF{cache: cache.New(cacheTTL), travelTimeTable: travelTimeTable},
		INFP_ID:         &INFP{cache: cache.New(cacheTTL), travelTimeTable: travelTimeTable},
		SED_ID:          &SED{cache: cache.New(cacheTTL), travelTimeTable: travelTimeTable},
		KNMI_ID:         &KNMI{cache: cache.New(cacheTTL), travelTimeTable: travelTimeTable},
		NCS_ID:          &NCS{resolvers: builtinResolvers, cache: cache.New(cacheTTL), travelTimeTable: travelTimeTable},
		NRCAN_ID:        &NRCAN{cache: cache.New(cacheTTL), travelTimeTable: travelTimeTable},
		GEONET_ID:       &GEONET{cache: cache.New(cacheTTL), travelTimeTable: travelTimeTable},
		ICL_ID:          &ICL{cache: cache.New(cacheTTL), travelTimeTable: travelTimeTable},
		BMKG_ID:         &BMKG{resolvers: builtinResolvers, cache: cache.New(cacheTTL), travelTimeTable: travelTimeTable},
		DOST_ID:         &DOST{cache: cache.New(cacheTTL), travelTimeTable: travelTimeTable},
		AFAD_ID:         &AFAD{cache: cache.New(cacheTTL), travelTimeTable: travelTimeTable},
		KRDAE_ID:        &KRDAE{cache: cache.New(cacheTTL), travelTimeTable: travelTimeTable},
		KNDC_ID:         &KNDC{cache: cache.New(cacheTTL), travelTimeTable: travelTimeTable},
		TMD_ID:          &TMD{cache: cache.New(cacheTTL), travelTimeTable: travelTimeTable},
		SSN_ID:          &SSN{cache: cache.New(cacheTTL), travelTimeTable: travelTimeTable},
		BGS_ID:          &BGS{cache: cache.New(cacheTTL), travelTimeTable: travelTimeTable},
	}, nil
}
