package main

type libraryKey struct {
	Name    string
	Version string
}

type licenseOverride struct {
	Module      string
	LicenseName string
	LicenseURL  string
	Notice      string
	Files       []string
}

// Reviewed supplements for licenses the collector cannot identify automatically.
// Match an exact package and version; Files are relative to the module directory.
var licenseAllowlist = map[libraryKey]licenseOverride{
	{
		Name:    "github.com/golang/freetype/raster",
		Version: "v0.0.0-20170609003504-e2365dfdc4a0",
	}: {
		Module:      "github.com/golang/freetype",
		LicenseName: "FTL OR GPL-2.0-or-later",
		LicenseURL:  "https://github.com/golang/freetype/blob/e2365dfdc4a0/LICENSE",
		Notice:      "Copyright 2010 The Freetype-Go Authors. All rights reserved.",
		Files:       []string{"LICENSE", "licenses/ftl.txt", "licenses/gpl.txt"},
	},
}
