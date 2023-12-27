# Changelog

Starting from v2.2.5, all notable changes to this project will be documented in this file.

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
