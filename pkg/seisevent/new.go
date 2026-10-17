package seisevent

import (
	"time"

	"github.com/anyshake/observer/pkg/cache"
	"github.com/anyshake/observer/pkg/dnsquery"
	"github.com/anyshake/observer/pkg/timesource"
	"github.com/bclswl0827/travel"
)

func New(timeSrc *timesource.Source, cacheTTL time.Duration) (map[string]IDataSource, error) {
	if timeSrc == nil {
		timeSrc = timesource.New(time.Now)
	}

	travelTimeTable, err := travel.NewAK135()
	if err != nil {
		return nil, err
	}

	builtinResolvers, err := dnsquery.NewResolvers()
	if err != nil {
		return nil, err
	}
	dnsClient, err := dnsquery.NewClient(builtinResolvers, len(builtinResolvers), 5*time.Second)
	if err != nil {
		return nil, err
	}

	return map[string]IDataSource{
		AFAD_ID: &AFAD{
			cache:           cache.NewGeneric[[]Event](cacheTTL),
			travelTimeTable: travelTimeTable,
			timeSource:      timeSrc,
		},
		BCSF_ID: &BCSF{
			cache:           cache.NewGeneric[[]Event](cacheTTL),
			travelTimeTable: travelTimeTable,
		},
		BGS_ID: &BGS{
			cache:           cache.NewGeneric[[]Event](cacheTTL),
			travelTimeTable: travelTimeTable,
		},
		BMKG_ID: &BMKG{
			dnsClient:       dnsClient,
			cache:           cache.NewGeneric[[]Event](cacheTTL),
			travelTimeTable: travelTimeTable,
		},
		CEA_ID: &CEA{
			cache:           cache.NewGeneric[[]Event](cacheTTL),
			travelTimeTable: travelTimeTable,
		},
		CENC_APP_ID: &CENC_APP{
			cache:           cache.NewGeneric[[]Event](cacheTTL),
			travelTimeTable: travelTimeTable,
		},
		CENC_WEB_ID: &CENC_WEB{
			cache:           cache.NewGeneric[[]Event](cacheTTL),
			travelTimeTable: travelTimeTable,
			timeSource:      timeSrc,
		},
		CENC_WOLFX_ID: &CENC_WOLFX{
			cache:           cache.NewGeneric[[]Event](cacheTTL),
			travelTimeTable: travelTimeTable,
		},
		CWA_SC_ID: &CWA_SC{
			dnsClient:       dnsClient,
			cache:           cache.NewGeneric[[]Event](cacheTTL),
			travelTimeTable: travelTimeTable,
		},
		DOST_ID: &DOST{
			cache:           cache.NewGeneric[[]Event](cacheTTL),
			travelTimeTable: travelTimeTable,
		},
		EMSC_ID: &EMSC{
			cache:           cache.NewGeneric[[]Event](cacheTTL),
			travelTimeTable: travelTimeTable,
		},
		GA_ID: &GA{
			cache:           cache.NewGeneric[[]Event](cacheTTL),
			travelTimeTable: travelTimeTable,
		},
		GEONET_ID: &GEONET{
			cache:           cache.NewGeneric[[]Event](cacheTTL),
			travelTimeTable: travelTimeTable,
		},
		GFZ_ID: &GFZ{
			cache:           cache.NewGeneric[[]Event](cacheTTL),
			travelTimeTable: travelTimeTable,
		},
		HKO_ID: &HKO{
			cache:           cache.NewGeneric[[]Event](cacheTTL),
			travelTimeTable: travelTimeTable,
		},
		ICL_ID: &ICL{
			cache:           cache.NewGeneric[[]Event](cacheTTL),
			travelTimeTable: travelTimeTable,
		},
		INFP_ID: &INFP{
			cache:           cache.NewGeneric[[]Event](cacheTTL),
			travelTimeTable: travelTimeTable,
		},
		INGV_ID: &INGV{
			cache:           cache.NewGeneric[[]Event](cacheTTL),
			travelTimeTable: travelTimeTable,
		},
		JMA_OFFICIAL_ID: &JMA_OFFICIAL{
			cache:           cache.NewGeneric[[]Event](cacheTTL),
			travelTimeTable: travelTimeTable,
		},
		JMA_P2PQUAKE_ID: &JMA_P2PQUAKE{
			cache:           cache.NewGeneric[[]Event](cacheTTL),
			travelTimeTable: travelTimeTable,
		},
		JMA_WOLFX_ID: &JMA_WOLFX{
			cache:           cache.NewGeneric[[]Event](cacheTTL),
			travelTimeTable: travelTimeTable,
		},
		KMA_ID: &KMA{
			cache:           cache.NewGeneric[[]Event](cacheTTL),
			travelTimeTable: travelTimeTable,
		},
		KNDC_ID: &KNDC{
			cache:           cache.NewGeneric[[]Event](cacheTTL),
			travelTimeTable: travelTimeTable,
		},
		KNMI_ID: &KNMI{
			cache:           cache.NewGeneric[[]Event](cacheTTL),
			travelTimeTable: travelTimeTable,
		},
		KRDAE_ID: &KRDAE{
			cache:           cache.NewGeneric[[]Event](cacheTTL),
			travelTimeTable: travelTimeTable,
		},
		NCS_ID: &NCS{
			dnsClient:       dnsClient,
			cache:           cache.NewGeneric[[]Event](cacheTTL),
			travelTimeTable: travelTimeTable,
		},
		NRCAN_ID: &NRCAN{
			cache:           cache.NewGeneric[[]Event](cacheTTL),
			travelTimeTable: travelTimeTable,
		},
		PALERT_ID: &PALERT{
			cache:           cache.NewGeneric[[]Event](cacheTTL),
			travelTimeTable: travelTimeTable,
			timeSource:      timeSrc,
		},
		SCEA_ID: &SCEA{
			cache:           cache.NewGeneric[[]Event](cacheTTL),
			travelTimeTable: travelTimeTable,
			timeSource:      timeSrc,
		},
		SED_ID: &SED{
			cache:           cache.NewGeneric[[]Event](cacheTTL),
			travelTimeTable: travelTimeTable,
		},
		SSN_ID: &SSN{
			cache:           cache.NewGeneric[[]Event](cacheTTL),
			travelTimeTable: travelTimeTable,
		},
		TMD_ID: &TMD{
			cache:           cache.NewGeneric[[]Event](cacheTTL),
			travelTimeTable: travelTimeTable,
		},
		USGS_ID: &USGS{
			cache:           cache.NewGeneric[[]Event](cacheTTL),
			travelTimeTable: travelTimeTable,
		},
		USP_ID: &USP{
			cache:           cache.NewGeneric[[]Event](cacheTTL),
			travelTimeTable: travelTimeTable,
		},
	}, nil
}
