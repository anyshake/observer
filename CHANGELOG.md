# Changelog

Starting from v2.2.5, all notable changes to this project will be documented in this file.

## v4.3.1

### Release Notes

This is a **minor bugfix release** focused on **data correctness and metadata accuracy**, addressing issues discovered after v4.3.0.

### New Features

- Support **E-C131G** device metadata.

### Bug Fixes

- Fixed **sign extension for 24-bit signed integer data**, correcting decoding errors caused by non-standard integer width.
- Updated **device gain values** to ensure metadata accuracy and correct amplitude scaling.

## v4.3.0

### Release Notes

This release focuses on **reliability, maintainability, and long-term operability**, with a strong emphasis on **automatic upgrade infrastructure**, **time synchronization precision**, and **network robustness**.

The introduction of a **fully integrated automatic upgrade mechanism** significantly reduces maintenance overhead, allowing AnyShake Observer to stay up to date with minimal user intervention. At the same time, continued refinements to **NTP, GNSS time handling, and clock drift compensation** further improve timestamp accuracy in both online and offline deployments.

This version also delivers multiple **stability fixes** across protocol handling, metadata accuracy, and connectivity edge cases, ensuring smoother long-term operation in production environments.

### What’s New

With this release, AnyShake Observer gains the ability to **self-update safely and predictably**, including version compatibility checks and support for pre-release channels. Time synchronization logic has been further refined to better handle **local clock drift**, **network latency**, and **unstable environments**, while new mechanisms improve **real-time data forwarding and observability**.

Several internal refactors improve correctness and reduce subtle failure modes, particularly in **device metadata, file handling, and protocol behavior**.

### New Features

- Implemented **automatic upgrade system**, enabling seamless version updates.
- Added support for **pre-release versions** in the build and upgrade pipeline.
- Enabled treating **Git tags as latest releases** for version resolution.
- Introduced **real-time data forwarding**, allowing external consumers to subscribe to live streams.
- Added **manual delay control** for the built-in NTP server.
- Improved **NTP server accuracy**, including higher-precision reference handling.
- Set **GNSS as NTP server refid**, improving traceability and diagnostics.
- Enhanced **clock drift compensation** to better handle long-running deployments.
- Improved **time synchronization accuracy** under unstable local clocks.
- Added **latest log viewing** support via internal logging refinements (infrastructure side).

### Bug Fixes

- Fixed **connectivity issues with CWA data source**, improving reliability.
- Corrected **gain values for accelerometer metadata**.
- Fixed **file ID changes** caused by list updates.
- Resolved **TCP read blocking** issues.
- Fixed **incorrect gain factors** in device metadata.
- Fixed **Docker build failures** and related build inconsistencies.
- Corrected **versioning and build flag handling** issues.
- Fixed **protocol handling regressions** affecting frame stability.
- Resolved **time drift issues** in NTP mode under certain conditions.
- Fixed **typos in logs and diagnostics output**.
- Improved robustness of **local clock drift handling** to prevent cascading timestamp errors.

## v4.2.0

### Release Notes

This release brings a wide range of **time synchronization enhancements, stability improvements, and new data sources**. The configuration field `ntpclient.endpoint` is now deprecated and replaced by **`ntpclient.pool`**, enabling synchronization with **multiple NTP servers** for significantly improved network time accuracy. At the same time, the synchronization mechanism has been further refined, taking overall **stability and precision** to the next level.

Our **crowdfunding campaign has successfully concluded**, and we sincerely thank every supporter — your trust and encouragement continue to drive us forward. With this milestone achieved, we are now preparing to enter **mass production**, moving ahead with confidence toward the official **shipment of AnyShake**.

### What’s New

This version achieves **complete independence from the system clock** — even if the system time experiences large jumps, the internal clock of AnyShake Observer remains unaffected. In **GNSS mode**, the software can now run **entirely offline** and also act as an **NTP server** for other devices on the local network. For example, if there are devices in your LAN without GNSS capability but still requiring **high-precision time synchronization**, you can simply point their NTP address to AnyShake Observer.

In addition, this update introduces several **new seismic data sources**, further expanding available monitoring channels.

### New Features

This release introduces **major improvements in time handling, data sources, and usability features** to make AnyShake Observer more reliable and versatile across deployment scenarios.

- Added **Mexican National Seismological Service** data source.
- Improved **tolerance to timestamp jitter**.
- Implemented reliance on **monotonic time**, completely removing dependency on the system clock.
- Introduced an **NTP forwarding service**, allowing the system to act as a **Stratum 1 time source**.
- The **Web panel** now supports viewing the latest logs.
- Achieved **higher-precision timestamps**.
- In **GNSS mode**, NTP synchronization can now be skipped to allow **fully offline operation**.
- Added **Polish translation** ([#33](https://github.com/anyshake/observer/pull/33)).
- Optimized **browser rendering styles on Windows**.
- Added **British Geological Survey (BGS)** data source.
- Changed configuration from `ntpclient.endpoint` to **`ntpclient.pool`**, enabling synchronization with multiple NTP sources.

### Bug Fixes

This version also delivers **critical stability fixes and platform compatibility improvements**, ensuring more robust performance under diverse conditions.

- Fixed **build failure and availability on Windows 7**.
- Program now returns **error code 1** on fatal errors, allowing external shells to capture and handle it.
- Allowed **unlimited systemd restart attempts**, especially useful when USB serial ports are disturbed.
- Fixed **buffer overflow issue in Protocol v3** that could cause a panic.
- Added **dynamic compensation in NTP mode** to correct time drift caused by oscillator frequency deviation.
- In **GNSS mode**, host time is now always synchronized with AnyShake Explorer.
- Implemented **link delay compensation** during transmission.
- Added **DNSCrypt-based resolution** of `scweb.cwa.gov.tw` to avoid **DNS spoofing in China**.
- Fixed **SQLite compatibility issues** on platforms such as Windows 386.
- Resolved a **Protocol v2 defect** that could cause infinite synchronization attempts.

## v4.1.0

### Release Notes

This update continues our mission to make seismic monitoring more powerful, accessible, and secure. It adds key network and time-handling improvements, introduces better fault tolerance, and enhances QuakeSense detection flexibility. With this release, AnyShake Observer becomes even more robust in various deployment environments.

#### Crowdfunding Now Live!

Don’t miss out on the launch of **AnyShake Explorer** — our open-source, high-performance seismic acquisition hardware. Available now on Crowd Supply: [www.crowdsupply.com/senseplex/anyshake-explorer](https://www.crowdsupply.com/senseplex/anyshake-explorer)

### What’s New

This version delivers major reliability enhancements to time synchronization, networking flexibility, and QuakeSense analytics. It also features UI polish, bug fixes, and better compatibility across diverse user environments.

### New Features

- Integrated **FRP (Fast Reverse Proxy)** client service for connecting devices behind NAT/firewall.
- Built-in **public FRP server** support for quick remote access without external configuration.
- Added **support for custom TLS settings** to improve secure deployment scenarios.
- **Random topic assignment** in QuakeSense for enhanced event stream handling in distributed setups.

### Bug Fixes

- Fixed **time offset issues after system hibernation (e.g., on battery powered laptops)**, ensuring consistent timestamps.
- Improved **NTP mode accuracy** under varying system load conditions.
- Resolved a bug causing **FRP connections not to close cleanly**, preventing potential resource leaks.
- Fixed **login error prompts** to show more informative messages on failure.
- Normalized comparison logic between **null values and empty arrays**, reducing edge case errors.
- Moved potentially blocking code into goroutines to **ensure timeout mechanisms work properly**.
- Removed an **invalid protocol entry** to prevent runtime parsing issues.
- Corrected several **minor typos and documentation mismatches**.
- Refined **login error UI** for better user experience.
- Removed **background colors** on images for a cleaner visual presentation.
- Applied minor **style and layout adjustments** for consistent appearance.

## v4.0.3

### Release Notes

This patch release focuses on improving stability, performance, and overall robustness. It resolves several critical issues, including a potential crash in the helicorder module, and introduces subtle improvements to the application's lifecycle management and UI.

### What’s New

This update enhances reliability during seismic clip exports and real-time data handling, ensuring smoother operation in both background services and user interactions.

### New Features

- Added **timeout control** for service modules to prevent indefinite blocking and ensure smoother background operations.
- Introduced a **timeout threshold** for program exit, improving stability during abnormal shutdown scenarios.

### Bug Fixes

- Fixed **FIR filter issues** in the QuakeSense module and WAV audio export to restore accurate signal processing.
- Corrected a UI issue where **buttons did not reset** properly after submitting field settings.
- Fixed a critical **integer divide by zero** error in the helicorder module that could cause application crashes during rendering.
- Validated **helicorder settings** before applying them to prevent runtime configuration errors.
- Improved system resilience by adding **connectivity checks** before reporting status.
- Normalized **string comparisons** to avoid subtle logic errors.
- Updated **broken documentation links** for easier access to support materials.
- Applied **minor style adjustments** for improved visual consistency.

## v4.0.2

### Release Notes

This patch release focuses on reliability improvements and critical bug fixes. It introduces STEIM2 compression support for MiniSEED export and resolves issues affecting helicorder rendering and token expiration.

### What’s New

Enhancements in this version contribute to better performance, interoperability, and system integrity, especially when exporting seismic clips and managing real-time data views.

### New Features

- Added support for exporting **MiniSEED clips with STEIM2 compression**, fixing Waves' compatibility issue on INT32 encoding.

### Bug Fixes

- Fixed a **division by zero** error in the helicorder generator under certain data conditions.
- Corrected an issue where **temporary tokens never expired**, improving session security and compliance.

## v4.0.1

### Release Notes

A minor but essential update focusing on stability, performance, and new regional data support. This release continues our commitment to delivering high-quality, open-source seismic monitoring software.

### What’s New

This version includes critical bug fixes, UI improvements, and a new data integration, further enhancing the reliability and usability of the AnyShake Observer platform.

### New Features

- Added **Thai Meteorological Department** as a new data agency, expanding regional data support.
- Enhanced UI to **highlight unsaved changes**, reducing the risk of accidental data loss.
- Optimized **error handling** for more robust system behavior.

### Bug Fixes

- Fixed a link error with **Go 1.24.3** that could cause build failures on Darwin.
- Prevented unnecessary **message bus re-creation**, improving runtime stability.
- Corrected **translation keys** for more accurate localization.

## v4.0.0

### Release Notes

We are proud to announce **AnyShake Observer v4.0.0**, a major release that marks a transformative step forward in our seismic monitoring software. This update is **incompatible with all previous versions** due to significant architectural changes and feature enhancements.

#### Now crowdfunding on CrowdSupply

Alongside this release, we are launching the **AnyShake Explorer**, a fully open-source, low-cost, and reliable seismic data acquisition device. Designed for researchers, hobbyists, and professionals alike, the Explorer integrates seamlessly with AnyShake Observer and sets a new standard for affordable and powerful seismology tools.

### What’s New

This version represents a leap forward in terms of **usability, security, reliability, and functionality** — setting a new benchmark in open-source seismic software.

### New Features

- Simplified deployment with just **40 lines of configuration** to get started.
- Full support for **legacy v1**, **mainline v2**, and **latest v3 protocols** used by AnyShake Explorer.
- Enhanced coroutine lifecycle management to ensure clean, stable module operations.
- **Multi-channel sampling** support to accommodate new Explorer variants with both 3-axis geophones and 3-axis accelerometers.
- Efficient waveform data archiving and fast querying for reduced CPU overhead.
- Export station metadata in **SeisComp XML** and **FDSNWS StationXML** formats for Explorer devices.
- Upgraded API architecture combining **GraphQL + RESTful** interfaces.
- WebSocket endpoints now support **120 seconds of historical waveform data**, reducing chart fill time.
- Data export now includes both **TXT (for MATLAB)** and **WAV (for audio sharing)** formats.
- Fully redesigned and optimized UI — **minimalist design with smoother performance** even under large data volumes.
- Web interface now supports editing station metadata and managing coroutine services dynamically.
- Customizable waveform panel layout via drag/zoom actions — with memory and locking features.
- Introduced **QuakeSense service**: earthquake detection using **Classic STA/LTA** and **Z-Detect**, with rapid alerts via **MQTT**.
- Helicorder images can now be saved in **PNG** format (previously only SVG).
- Removed deprecated seismic event data APIs.
- Improved **SeedLink protocol** stability.
- Added **multi-language support** with translations in 9 languages.

### Bug Fixes

- Resolved high CPU usage issue caused by "empty-run" behavior after Explorer device disconnection.
- Fixed significant time offset when AnyShake Observer starts before Explorer under NTP mode.
- Resolved decoder thread crash when receiving delayed data via serial-to-TCP device.
- Fixed numerous potential data race issues for improved multithreaded stability.

---

## v3.6.1

### Bug Fixes

- Disable waveform normalization completely.
- Disable compression for exporting MiniSEED in history service.

## v3.6.0

### Bug Fixes

- Fixed incorrect coordinates when legacy mode enabled.

### New Features

- Support customizing helicorder image size, waveform scale and samples per span.

## v3.5.0

### New Features

- Provide an option to disable MiniSEED compression.

## v3.4.6

### Bug Fixes

- Handle coordinates correctly when by hot switching.

### Chore

- Update copyright and email address.

## v3.4.5

### Bug Fixes

- Fixed device ID calculation when it is not set on AnyShake Explorer.
- Fixed the issue of not being able to reconnect to AnyShake Explorer after hard reset.
- Support GNSS mode enabling or disabling by hot switch.

### Chore

- Update logotypes for AnyShake.
- Disable waveform normalization.

## v3.4.4

### CI/CD Changes

- Added precompiled packages for ARMv5 devices.

### Bug Fixes

- Determine if eligible to enable cache to aviod OOM.
- Discard record that sample rate jitters in history.
- Load dummy settings if service configration field not found.

### Chore

- Added the code of conduct.
- Skip empty files when getting file list.
- Set archiver clean interval to 1 hour.

## v3.4.3

### New Features

- Display progress for helicorder tasks in frontend.

### Refactor

- Plot helicorder hourly for some weak devices.
- Split helicorder line into multiple segments when gap occurs.

## v3.4.2

### Refactor

- Reduced memory usage when plotting helicorder images.
- Record sequence numbers to file to reduce initialization time.

## v3.4.1

### New Features

- Support previewing helicorder images in frontend.
- Added earthquake event source API of KNDC.

### Bug Fixes

- Fixed channel names in inventory response.

### Refactor

- Cache parsed results instead of crawler response.
- Return nil directly for API endpoints if miniSEED or helicorder service is disabled.

## v3.4.0

### New Features

- Support downloading helicorder images in frontend.
- Plot helicorder image periodically.
- Added event data source of KRADE.

### Bug Fixes

- Fixed parsing errors on data sources of CENC, CWA.
- Fixed malformed serial number in legacy mode.

### Refactor

- Changed the cache data type to any.
- Made some tricks to make TypeScript happy.

### Chore

- Attach station information when exporting miniSEED data.
- Disable captcha input until captcha image is loaded.
- Clean expired miniSEED files every hour.

## v3.3.3

### New Features

- Added non-official event data sources (CWA, JMA, CENC).
- Added AFAD and DOST event data sources.

### Refactor

- Enhanced support for devices without GNSS support.

## v3.3.2

### Bug Fixes

- Handle auto flag properly for event data sources of China.
- Fixed occasional checksum error in legacy mode.

### Refactor

- Use Web Workers to process seismic data.

### Chore

- chore: Update slgo dependency to v0.0.4.

## v3.3.1

### Bug Fixes

- Fixed the timestamp accuracy issue of event data source of JMA.
- Fixed token refresh logic in Main component.

### Refactor

- Insert timestamp before buffering legacy mode packets.

### Chore

- Create a copy of the configuration file when executing make run.
- Updated retention default value to 120s in global configuration.

## v3.3.0

### New Features

- Added multi-user support and permission management.
- Added earthquake event source API of BMKG.
- Display AnyShake Explorer device serial number on home page.

### Bug Fixes

- Fix timezone issue on earthquake event source API of CEA.

### Chore

- Refined the module names in the log.
- Add more detailed description to API documentation.
- Replace table component with MUI in the frontend.
- Use year provided by Cloudflare to fill the timing of earthquake event source of CWA.

## v3.2.5

### New Features

- Added more earthquake event source APIs.

### Bug Fixes

- Completely fixed the sampling rate jitter issue in legacy mode.

### Chore

- Convert device ID to hexadecimal string format.
- Replace SVG icons with Material Design Icons in the frontend.
- Send error message when server returned no earthquake event.

## v3.2.4

### CI/CD Changes

- Merge API docs Makefile into the main Makefile.
- Add support for OpenBSD builds.

### New Features

- Add SQLite support on all architectures.
- Provide i18n support for earthquake event data source.
- Show flags next to earthquake data sources.

## v3.2.3

### CI/CD Changes

- Deprecated `.env` file based versioning in frontend, use `craco.config.js` to generate version info.
- Support Multi-platform Docker image build.
- Build and release more easily with GitHub Actions.
- Some Dockerfile optimizations to reduce image size.

### Bug Fixes

- Fixed time offset in EQZT data source.
- Fixed Windows 7 compatibility issue.

## v3.2.2

### New Features

- Enable mainline data protocol support for AnyShake Explorer without GNSS module.
- Added earthquake event source API support of EQZT.

## v3.2.1

### Breaking Changes

- **Data Protocol**: The AnyShake Explorer mainline data protocol has been updated again. **If you are using AnyShake Explorer implemented by STM32, please rebuild and burn the firmware to the latest version, if you are using AnyShake Explorer implemented by ESP8266, just ignore this message.**

### New Features

- Print log messages from AnyShake Explorer driver.
- Improved the stability of the legacy data protocol.
- Support Windows ARM32 architecture.

### Bug Fixes

- Fix SQLite "out of memory" issue on Windows release build.
- Fix release names in the GitHub Actions workflow.
- Use `time.NewTimer` instead of `time.After` to avoid memory leak.

## v3.2.0

### New Features

- Support custom NTP query retry times and timeout.
- Retry querying NTP server when the first attempt fails.
- Support rotating log files by day or by size.
- Validate the configuration file before starting the application.

### Bug Fixes

- Fixed potential data race issues.
- Replace concurrent-map with haxmap to avoid OOM.
- Fixed build failure on windows/386, linux/mips\* architectures.

## v3.1.1

### New Features

- Support export history waveform as MiniSEED.

### Bug Fixes

- Discard records with different sample rate for SeedLink history.

## v3.1.0

### New Features

- Added support for forwarding raw messages via TCP.
- Re-support for the SeedLink protocol.
- Updated Go version to v1.23.0.
- Updated Gin version to v1.10.0.

## v3.0.5

### New Features

- Update network time at 00:00:00 every day to avoid time drift.

### Bug Fixes

- Alleviated the problem of frequent jitter in sampling rate.
- Fixed the time offset of up to several hours caused by external compensation of timestamp.

## v3.0.4

### Bug Fixes

- Fixed an issue where the timepicker component in the frontend would not update the selected time when the time range was changed.

## v3.0.3

### Bug Fixes

- Fixed an issue on the frontend where the CPU and memory usage percentages were not displayed correctly.
- Resolved a backend issue where the timestamp, latitude, longitude, and elevation values were always 0 when fallback values were expected.

## v3.0.2

### New Features

- Added ability to automatically fix time jitter when using internet NTP server as time source.

### Bug Fixes

- Fixed "insufficient arguments" error when using PostgreSQL as the database backend (see https://github.com/go-gorm/gorm/issues/6832#issuecomment-1946211186).
- Never check for sample rate consistency in MiniSEED and SAC records when in legacy mode.
- Send history buffer only if client requests in WebSocket API to avoid flooding the client.
- Lowered `MINISEED_ALLOWED_JITTER_MS` constant to 2 ms for better jitter tolerance.

## v3.0.1

### Bug Fixes

- Fixed the issue where MiniSEED recording in legacy mode would be interrupted due to sampling rate jitter.

## v3.0.0

### Breaking Changes

- **Data Protocol**: The AnyShake Explorer data protocol has been entirely refactored. **Please rebuild and burn the firmware of AnyShake Explorer to the latest version.**
- **Configuration File**: The configuration file layout has been completely overhauled. The old configuration file format is no longer supported.
- **SeedLink Server**: The SeedLink service has been temporarily removed and will be re-implemented in a future release.
- **API Endpoints**: Some request and response fields have been modified in API v1. Please refer to the built-in Swagger API documentation for details.

### New Features

- Added support for accessing AnyShake Explorer via a serial-to-Ethernet converter.
- Introduced custom channel prefixes (e.g., HH*, SH*, EH\*).
- Added log dumping functionality with multiple output levels.
- Enhanced data processing and storage efficiency.
- Improved the accuracy of reading time from the Internet NTP server.
- Refined component lifecycle management using dependency injection for better module decoupling.
- Implemented an asynchronous message bus to optimize application execution efficiency.
- Established a GraphQL-based routing endpoint in preparation for API v2.
- Dockerized the application for easier and faster deployment.

### Bug Fixes

- Completely resolved the gap issue in MiniSEED records.

---

## v2.12.5

- Fix gaps in MiniSEED records

## v2.12.4

- Update frontend map tile provider to OpenStreetMap

## v2.12.3

- Format frontend code using ESLint and Prettier
- Use `time.Ticker` to collect geophone counts by second

## v2.12.2

- Always use fallback locale if the preferred locale is not available
- Optimize the serial port reading process

## v2.12.1

- Sort earthquake event source API by name in frontend
- Fix intensity calculation issue of CSIS in frontend
- Removal of unused utility functions

## v2.12.0

- Support lifecycle configuration for records in database

## v2.11.10

- Set minimum TLS version to 1.2 for HTTP request utility
- Update frontend dependencies

## v2.11.9

- Show disk free space in banner instead of station UUID

## v2.11.8

- Add earthquake event source API support of CEA and INGV
- Reuse of int32 array encoding and decoding functions
- Specify the minimun TLS version to 1.0 in HTTP client

## v2.11.7

- Add earthquake event source API support of KMA

## v2.11.6

- Ensure that there is only one Websocket connection after reconnecting

## v2.11.5

- Fix a frontend issue that causes event querying to fail

## v2.11.4

- Set non-zero start time when SeedLink DATA command has no extra argument

## v2.11.3

- Response with OK when SeedLink DATA command has no extra argument

## v2.11.2

- Code style improvements again
- Basic implementation of SeedLink DATA command
- Fix frontend issue where the input component does not update its value

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

---
