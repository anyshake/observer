#!/usr/bin/env python3
# -*- coding: utf-8 -*-

from os import makedirs, path
from shutil import rmtree
from struct import pack as struct_pack


def write_tile_pack(z, x, tiles_dict, output_dir):
    if not tiles_dict:
        return

    dir_path = path.join(output_dir, f"z{z}")
    makedirs(dir_path, exist_ok=True)

    tiles = sorted(tiles_dict.items())
    header = []

    offset = 0
    for y, data in tiles:
        size = len(data)
        header.append((y, offset, size))
        offset += size

    with open(path.join(dir_path, f"x{x}.bin"), "wb") as f:
        f.write(struct_pack("<I", len(header)))

        for y, off, size in header:
            f.write(struct_pack("<III", y, off, size))

        for _, data in tiles:
            f.write(data)


def mbtiles_to_bin(input_file: str, output_dir: str):
    from sqlite3 import connect as sqlite3_connect

    conn = sqlite3_connect(input_file)
    cur = conn.cursor()
    query = """
    SELECT zoom_level, tile_column, tile_row, tile_data
    FROM tiles
    ORDER BY zoom_level, tile_column, tile_row
    """
    cur.execute(query)

    current_key = None
    tiles_dict = {}

    for z, x, y, data in cur:
        y = (1 << z) - 1 - y
        key = (z, x)

        if current_key is None:
            current_key = key

        if key != current_key:
            write_tile_pack(current_key[0], current_key[1], tiles_dict, output_dir)
            tiles_dict = {}
            current_key = key

        tiles_dict[y] = data

    if current_key:
        write_tile_pack(current_key[0], current_key[1], tiles_dict, output_dir)


def main():
    from argparse import ArgumentParser

    parser = ArgumentParser(
        description="Script to convert MBTiles to a custom binary format for map tiles"
    )
    parser.add_argument("input", help="Path to the input MBTiles file")
    parser.add_argument("output", help="Path to the output directory for binary tiles")
    args = parser.parse_args()

    if path.exists(args.output):
        rmtree(args.output)
    makedirs(args.output, exist_ok=True)
    mbtiles_to_bin(args.input, args.output)


if __name__ == "__main__":
    main()
