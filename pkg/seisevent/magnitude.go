package seisevent

import (
	"strings"
)

func ParseMagnitude(typ string) MagnitudeType {
	_typ := strings.ToUpper(typ)

	switch {
	case _typ == "M": // Magnitude
		return "M"
	case strings.Contains(_typ, "ML"): // Local magnitude (Richter scale)
		return "Ml"
	case strings.Contains(_typ, "MS"): // Surface-wave magnitude
		return "MS"
	case strings.Contains(_typ, "MW"): // Moment magnitude
		return "Mw"
	case strings.Contains(_typ, "MB"): // Body-wave magnitude
		return "Mb"
	case strings.Contains(_typ, "MD"): // Duration magnitude
		return "Md"
	}

	return MagnitudeType(_typ)
}
