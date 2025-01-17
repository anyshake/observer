<p align="center">
  <img src="https://raw.githubusercontent.com/anyshake/logotype/master/banner_observer.png" width="500" alt="banner" />
</p>

[![Codacy Badge](https://app.codacy.com/project/badge/Grade/7b75168a5b03403987122835d74bb448)](https://app.codacy.com/gh/anyshake/observer/dashboard)
[![Downloads](https://img.shields.io/github/downloads/anyshake/observer/total.svg)](https://github.com/anyshake/observer/releases/latest)
[![Go Report Card](https://goreportcard.com/badge/github.com/anyshake/observer)](https://goreportcard.com/report/github.com/anyshake/observer)
[![Build Status](https://github.com/anyshake/observer/actions/workflows/release.yml/badge.svg)](https://github.com/anyshake/observer/actions/workflows/release.yml)
[![Latest Release](https://img.shields.io/github/release/anyshake/observer.svg)](https://github.com/anyshake/observer/releases/latest)

## Overview

AnyShake Observer is an open-source, cross-platform software that can be used to monitor, archive, and export seismic data from [AnyShake Explorer](https://github.com/anyshake/explorer) via serial port. It provides a user-friendly web-based interface to visualize and analyze the seismic data. For more professional users, it supports streaming via SeedLink protocol and exporting the data to SAC or MiniSEED format for further analysis.

This software is written in Go and TypeScript, which means it can easily port to a variety of OS and CPU architectures, even embedded Linux devices, AnyShake Observer also supports PostgreSQL, MariaDB (MySQL) and SQL Server as seismic data archiving engines.

## Documentation

Please visit [anyshake.org/docs/introduction](https://anyshake.org/docs/introduction) for quick start guide and more information.

## Features

- User-friendly web-based interface
- Mobile / Tablet friendly interface
- Query seismic waveform by time range
- Query seismic waveform by known event
- Link to share the seismic waveform
- Real-time seismic waveform display
- Swagger generated API documentation
- Support multiple database engines
- Support multiple languages, detected by browser
- Multiple seismic intensity standards, default to JMA
- Cross-platform, runs on Linux, Windows, macOS
- Ability to stream seismic data via SeedLink protocol
- Ability to export data to SAC or MiniSEED format
- AnyShake Explorer data checksum verification
- Auto reset AnyShake Explorer on error
- Flexible channel packet read length
- Variable serial port baud rate

## Preview

![Preview - Home](https://raw.githubusercontent.com/anyshake/logotype/master/preview_home.gif)
![Preview - Realtime](https://raw.githubusercontent.com/anyshake/logotype/master/preview_realtime.gif)
![Preview - History](https://raw.githubusercontent.com/anyshake/logotype/master/preview_history.gif)
![Preview - Export](https://raw.githubusercontent.com/anyshake/logotype/master/preview_export.gif)
![Preview - Settings](https://raw.githubusercontent.com/anyshake/logotype/master/preview_setting.gif)

## Credits

AnyShake Observer is designed and developed by [@bclswl0827](https://github.com/bclswl0827), test work is done by [@TenkyuChimata](https://github.com/TenkyuChimata).

The success of AnyShake Observer is inseparable from the following core libraries:

- [github.com/bclswl0827/mseedio](https://github.com/bclswl0827/mseedio): Pure Go library for reading and writing MiniSEED data.
- [github.com/bclswl0827/sacio](https://github.com/bclswl0827/sacio): Pure Go library for reading and writing SAC data.
- [github.com/bclswl0827/slgo](https://github.com/bclswl0827/slgo): Pure Go library used to build SeedLink server.

![Star History Chart](https://api.star-history.com/svg?repos=anyshake/observer&type=Date)
