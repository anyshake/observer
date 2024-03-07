# Changelog

Starting from v2.2.5, all notable changes to this project will be documented in this file.

## v2.11.1

- Some frontend code style improvements

## v2.11.0

- Using NoSQL database as SeedLink ring buffer backend
- Fix SeedLink command parsing issue which causes some clients to be unable to connect
- Fragmenting SeedLink packets to accommodate higher sampling rates
- Remove redundant `status` fields in `/api/v1/station` response
- Fix timestamp issue in geophone data collecting module

## v2.10.2

- Input component optimization
- Back to limiting waveform query duration to 1 hour

## v2.10.1

- Support download SeisComP3 XML inventory directly from the frontend
- Update swagger docs to include new API endpoints

## v2.10.0

- Update frontend dependencies credits in README
- Fixed the validation error of the datetime picker in the frontend
- New API endpoint `/api/v1/inventory` to get SeisComP3 XML inventory data
- Removed unused configuration fields, change geophone sensitivity unit to `V/m/s`

## v2.9.0

- Frontend refactoring: use functional components and hooks
- Support Butterworth bandpass filter in frontend waveform data processing
- API /api/v1/mseed: use unix timestamp as file modification time response
- Check for remote server error before parsing earthquake event data from SCEA API

## v2.8.1

- Show disk usage of current working directory

## v2.8.0

- Allow setting rate limitation for API endpoints

## v2.7.1

- Update frontend dependencies
- Ensure encoded SeedLink packet length is 512 bytes

## v2.7.0

- Support SeedLink buffer size customization
- Basic implementation of SeedLink buffer file

## v2.6.1

- Fixed a depencency issue in MiniSEED data processing

## v2.6.0

- Simple implementation of SeedLink buffer
- Add PowerShell frontend build script for Windows
- Replace CWB to CWA in earthquake event data source API

## v2.5.5

- Allows querying waveform within 24 hours in JSON format
- Save MiniSEED data by channel to separate files
- Remove support for MIPS64 architecture

## v2.5.4

- Use timestamp from backend for frontend chart

## v2.5.3

- Support channel multi-selection in SeedLink protocol

## v2.5.2

- Fix MiniSEED sample rate calculation error under unstable data link
- Use the timestamp of the first packet as reference time in each MiniSEED record

## v2.5.1

- Some efforts to make SeedLink protocol work properly
- Disable SeedLink by default
- Update README to introduce SeedLink protocol

## v2.5.0

- Basic Go implementation of SeedLink

## v2.4.1

- Skip TLS verification for earthquake event data source API
- Remove SQLite support due to MIPS architecture incompatibility
- Make frontend className conditional rendering logic more predictable
- Correct frontend map anchor point offset

## v2.4.0

- Optimization on CPU usage metrics calculation
- Some changes made to `/api/v1/station` response:
- Add `station` object field, which contains `uuid`, `name`, `station` (string), `network` and `location`
- Move the original `uuid` and `station` string fields to `station` object
- Move current API handlers to `v1` subdirectory to prepare for future API versions
- Regulating MiniSEED file naming rules to `NN.SSSSS.LL.D.yyyy.ddd.mseed`
- Regulating SAC file naming rules to `yyyy.ddd.hh.mm.ss.ffff.NN.SSSSS.LL.CCC.D.sac`
- Move `station`, `network` and `location` fields to `station_settings` in configuration file

## v2.3.1

- Translation optimizations
- Allow collapsing PGA, PGV, intensity data in realtime chart
- Automatically adjusting realtime chart height to fit the screen

## v2.3.0

- Update earthquake event arrival estimation algorithm
- Remove compensation IIR filter
- Change `altitude` field in `station_settings` to `elevation`
- Restart daemon only on-failure
- Add restart delay for daemon
- Follow the SemVer principles
- Update README reference docs
- More backend credits
- Run build workflow on tag creations only

## v2.2.6

- Fixed SAC file waveform lag issue caused by sample rate calculation

## v2.2.5

- Migrating to [anyshake/observer](https://github.com/anyshake/observer)
- Supplement of README, CHANGELOG, build instructions, etc.
- Use templates to standardize ISSUEs and Pull Requests
- Update repository frontend logos
