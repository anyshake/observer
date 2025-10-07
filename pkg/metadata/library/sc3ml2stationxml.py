#!/usr/bin/env python3
# -*- coding: utf-8 -*-

def sc3ml_to_stationxml(input_file: str, output_file: str):
    from obspy import read_inventory

    inv = read_inventory(input_file, format="SC3ML")
    inv.write(output_file, format="STATIONXML")


def main():
    import argparse

    parser = argparse.ArgumentParser(
        description="Convert SeisComP SC3ML inventory to FDSN StationXML"
    )
    parser.add_argument("input", help="Path to SeisComP XML (SC3ML) file")
    parser.add_argument("output", help="Output path to FDSN StationXML file")

    args = parser.parse_args()
    sc3ml_to_stationxml(args.input, args.output)


if __name__ == "__main__":
    main()
