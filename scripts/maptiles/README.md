## Build Vector Map Tiles from Natural Earth

[Natural Earth](https://www.naturalearthdata.com) provides free raster and vector datasets (cultural and physical) covering the most commonly used map features. These datasets are in the public domain and can be used for commercial purposes.

This guide explains how to obtain Natural Earth data and convert it into vector map tiles (MBTiles format) suitable for production use.

## Install Dependencies

On Debian/Ubuntu:

```bash
$ sudo apt update
$ sudo apt install -y gdal-bin tippecanoe python3
```

- `ogr2ogr` (from GDAL): converts Shapefile to GeoJSON
- `tippecanoe`: builds vector tiles (MBTiles)

For Windows developers, it's recommended to use the Windows Subsystem for Linux (WSL).

## Get Cultural Layers

Go to [www.naturalearthdata.com/downloads/10m-cultural-vectors/](https://www.naturalearthdata.com/downloads/10m-cultural-vectors/), download the following datasets:

- **Admin 0 - Countries >> Download countries**: The filename will be `ne_10m_admin_0_countries.zip`
- **Admin 1 – States, Provinces >> Download boundary lines**: The filename will be `ne_10m_admin_1_states_provinces_lines.zip`
- **Populated Places >> Download populated places**: The filename will be `ne_10m_populated_places.zip`

## Get Physical Layers

Go to [www.naturalearthdata.com/downloads/10m-physical-vectors/](https://www.naturalearthdata.com/downloads/10m-physical-vectors/), download the following datasets:

- **Coastline >> Download coastline**: The filename will be `ne_10m_coastline.zip`
- **Minor Islands >> Download minor islands**: The filename will be `ne_10m_minor_islands.zip`
- **Reefs >> Download reefs**: The filename will be `ne_10m_reefs.zip`
- **Ocean >> Download ocean**: The filename will be `ne_10m_ocean.zip`
- **Rivers + lake centerlines >> Download rivers and lake centerlines**: The filename will be `ne_10m_rivers_lake_centerlines.zip`
- **Lakes + Reservoirs >> Download lakes**: The filename will be `ne_10m_lakes.zip`
- **Playas >> Download playas**: The filename will be `ne_10m_playas.zip`
- **Antarctic Ice Shelves >> Download Antarctic ice shelves**: The filename will be `ne_10m_antarctic_ice_shelves_polys.zip`
- **Glaciated Areas >> Download glaciated areas**: The filename will be `ne_10m_glaciated_areas.zip`

## Generate GeoJSON

Move all the downloaded archives to a directory named `map_data`:

```bash
$ mkdir map_data
$ mv /path_to_your_downloaded_files/*.zip map_data/
```

Enter the `map_data` directory and unzip all the archives using a shell loop:

```bash
$ cd map_data
$ for f in *.zip; do
  unzip "$f" -d "${f%.zip}"
done
```

Convert the shapefile (`.shp`) to GeoJSON format using `ogr2ogr` command in a shell loop:

```bash
$ for d in */; do
  for f in "$d"/*.shp; do
    ogr2ogr -f GeoJSON "${f%.shp}.geojson" "$f"
  done
done
```

If no errors occured, you will be able to see `.geojson` file for each shapefile in the respective directories, check them using `ls` command:

```bash
$ ls */*.geojson
```

You will see a list of all the generated GeoJSON files, for example:

```bash
ne_10m_admin_0_countries/ne_10m_admin_0_countries.geojson
ne_10m_admin_1_states_provinces_lines/ne_10m_admin_1_states_provinces_lines.geojson
ne_10m_antarctic_ice_shelves_polys/ne_10m_antarctic_ice_shelves_polys.geojson
ne_10m_coastline/ne_10m_coastline.geojson
ne_10m_glaciated_areas/ne_10m_glaciated_areas.geojson
ne_10m_lakes/ne_10m_lakes.geojson
ne_10m_minor_islands/ne_10m_minor_islands.geojson
ne_10m_ocean/ne_10m_ocean.geojson
ne_10m_playas/ne_10m_playas.geojson
ne_10m_populated_places/ne_10m_populated_places.geojson
ne_10m_reefs/ne_10m_reefs.geojson
ne_10m_rivers_lake_centerlines/ne_10m_rivers_lake_centerlines.geojson
```

## Generate MBTiles

Define zoom levels using environment variables:

```bash
$ export MIN_ZOOM=1
$ export MAX_ZOOM=6
```

Generate MBTiles using `tippecanoe`:

```bash
$ tippecanoe -o ne_10m_admin_0_countries.mbtiles ne_10m_admin_0_countries/ne_10m_admin_0_countries.geojson -Z$MIN_ZOOM -z$MAX_ZOOM
$ tippecanoe -o ne_10m_admin_1_states_provinces_lines.mbtiles ne_10m_admin_1_states_provinces_lines/ne_10m_admin_1_states_provinces_lines.geojson -Z3 -z$MAX_ZOOM
$ tippecanoe -o ne_10m_antarctic_ice_shelves_polys.mbtiles ne_10m_antarctic_ice_shelves_polys/ne_10m_antarctic_ice_shelves_polys.geojson -Z$MIN_ZOOM -z$MAX_ZOOM
$ tippecanoe -o ne_10m_coastline.mbtiles ne_10m_coastline/ne_10m_coastline.geojson -Z$MIN_ZOOM -z$MAX_ZOOM
$ tippecanoe -o ne_10m_glaciated_areas.mbtiles ne_10m_glaciated_areas/ne_10m_glaciated_areas.geojson -Z$MIN_ZOOM -z$MAX_ZOOM
$ tippecanoe -o ne_10m_lakes.mbtiles ne_10m_lakes/ne_10m_lakes.geojson -Z$MIN_ZOOM -z$MAX_ZOOM
$ tippecanoe -o ne_10m_minor_islands.mbtiles ne_10m_minor_islands/ne_10m_minor_islands.geojson -Z$MIN_ZOOM -z$MAX_ZOOM
$ tippecanoe -o ne_10m_playas.mbtiles ne_10m_playas/ne_10m_playas.geojson -Z3 -z$MAX_ZOOM
$ tippecanoe -o ne_10m_populated_places.mbtiles ne_10m_populated_places/ne_10m_populated_places.geojson -Z$MIN_ZOOM -z$MAX_ZOOM
$ tippecanoe -o ne_10m_reefs.mbtiles ne_10m_reefs/ne_10m_reefs.geojson -Z3 -z$MAX_ZOOM
$ tippecanoe -o ne_10m_ocean.mbtiles ne_10m_ocean/ne_10m_ocean.geojson -Z$MIN_ZOOM -z$MAX_ZOOM
$ tippecanoe -o ne_10m_rivers_lake_centerlines.mbtiles ne_10m_rivers_lake_centerlines/ne_10m_rivers_lake_centerlines.geojson -Z3 -z$MAX_ZOOM
```

> Some layers use `-Z3` instead of `$MIN_ZOOM`, which can avoid geometry distortion / empty tiles at low zooms.

## Merge MBTiles

With all the individual MBTiles files generated, you can merge them into a single `world.mbtiles` file using `tile-join`:

```bash
$ tile-join -o world.mbtiles --force *.mbtiles
```

- `-f` overwrites the output
- Each input file becomes a separate layer
- Layer names default to their filenames

## Split Into Chunks

MBTiles uses the SQLite format. Since the SQLite backend heavily relies on CGO, this makes cross-compilation, deployment, and portability more complex. Meanwhile, pure Go implementations of SQLite backend have performance limitations.

A more practical approach is to preprocess the MBTiles dataset and split it into smaller static chunks, enabling pure Go access without relying on SQLite.

To split the `world.mbtiles` file into smaller chunks, use the Python script `mbtiles2bin.py`:

```bash
$ ./mbtiles2bin.py world.mbtiles maptile_chunks
```

The output directory is organized by **zoom level (z)** and **tile column (x)**:

```
maptile_chunks/
  z0/
    x0.bin
    x1.bin
    ...
  z1/
    x0.bin
    x1.bin
    ...
  z2/
    ...
```

## Chunk Data Format

Each **.bin** file packs multiple tiles sharing the same (z, x) into a single binary blob, structured as follows:

```
+----------------------+------------------------+----------------------+
| Tile Count (uint32)  | Tile Index Table       | Tile Data Section    |
+----------------------+------------------------+----------------------+
```

### Tile Count

```
uint32 tile_count   # little-endian
```

Number of tiles stored in this file.

### Tile Index Table

Each entry:

```
uint32 y        # tile row (XYZ scheme, already flipped from TMS)
uint32 offset   # byte offset into data section
uint32 size     # tile data size in bytes
```

- Total size: `tile_count * 12 bytes`
- Entries are sorted by `y`

### Tile Data Section

Raw concatenation of tile payloads:

```
[tile_0_data][tile_1_data][tile_2_data]...
```

- No padding
- No compression added (original MBTiles data is preserved, typically gzip-compressed MVT)
