<p align="center">
  <img src="https://raw.githubusercontent.com/anyshake/logotype/master/banner_observer.png" width="500"/>
</p>

[![Build Status](https://github.com/anyshake/observer/actions/workflows/release.yml/badge.svg)](https://github.com/anyshake/observer/actions/workflows/release.yml)
[![Downloads](https://img.shields.io/github/downloads/anyshake/observer/total.svg)](https://github.com/anyshake/observer/releases/latest)
[![Latest Release](https://img.shields.io/github/release/anyshake/observer.svg?style=flat-square)](https://github.com/anyshake/observer/releases/latest)
[![Go Report Card](https://goreportcard.com/badge/github.com/anyshake/observer?style=flat-square)](https://goreportcard.com/report/github.com/anyshake/observer)

## Overview

AnyShake Observer is an open-source, cross-platform software that can be used to monitor, archive, and export seismic data from [AnyShake Explorer](https://github.com/anyshake/explorer) via serial port. It provides a user-friendly web-based interface to visualize and analyze the seismic data. For more professional users, it supports exporting the data to SAC or MiniSEED format for further analysis.

This software is written in Go and TypeScript, which means it can easily port to a variety of OS and CPU architectures, even embedded Linux devices, AnyShake Observer also supports PostgreSQL, MariaDB (MySQL) and SQL Server as seismic data archiving engines. We're currently trying to integrate the SEEDLink server on this software (Using pure Go implementation).

As of the release of the software documentation, AnyShake has successfully captured more than 40 earthquake events, the furthest captured earthquake event is [M 7.1 - 180 km NNE of Gili Air, Indonesia](https://earthquake.usgs.gov/earthquakes/eventpage/us7000krjx/executive), approximately 4,210 km, by the station located in Chongqing, China.

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

Thanks to the following tools and libraries, AnyShake Observer is made possible!

### Backend

 - [github.com/PuerkitoBio/goquery](https://github.com/PuerkitoBio/goquery)
 - [github.com/bclswl0827/go-serial](https://github.com/bclswl0827/go-serial)
 - [github.com/bclswl0827/sacio](https://github.com/bclswl0827/sacio)
 - [github.com/beevik/ntp](https://github.com/beevik/ntp)
 - [github.com/common-nighthawk/go-figure](https://github.com/common-nighthawk/go-figure)
 - [github.com/gin-contrib/gzip](https://github.com/gin-contrib/gzip)
 - [github.com/gin-gonic/gin](https://github.com/gin-gonic/gin)
 - [github.com/gorilla/websocket](https://github.com/gorilla/websocket)
 - [github.com/juju/ratelimit](https://github.com/juju/ratelimit)
 - [github.com/mackerelio/go-osstat](https://github.com/mackerelio/go-osstat)
 - [github.com/sbabiv/xml2map](https://github.com/sbabiv/xml2map)
 - [github.com/shirou/gopsutil](https://github.com/shirou/gopsutil)
 - [github.com/swaggo/files](https://github.com/swaggo/files)
 - [github.com/swaggo/swag](https://github.com/swaggo/swag)
 - [github.com/wille/osutil](https://github.com/wille/osutil)
 - [gorm.io/driver/mysql](https://github.com/go-gorm/mysql)
 - [gorm.io/driver/postgres](https://github.com/go-gorm/postgres)
 - [gorm.io/driver/sqlite](https://github.com/go-gorm/sqlite)
 - [gorm.io/driver/sqlserver](https://github.com/go-gorm/sqlserver)
 - [gorm.io/gorm](https://gorm.io/)
 - [github.com/bclswl0827/mseedio](https://github.com/bclswl0827/mseedio)
 - [github.com/fatih/color](https://github.com/fatih/color)
 - [github.com/json-iterator/go](https://github.com/json-iterator/go)
 - [github.com/swaggo/gin-swagger](https://github.com/swaggo/gin-swagger)

### Frontend

 - [axios](https://axios-http.com/)
 - [date-fns](https://date-fns.org/)
 - [highcharts](https://www.highcharts.com/)
 - [i18next](https://www.i18next.com/)
 - [i18next-browser-languagedetector](https://github.com/i18next/i18next-browser-languageDetector)
 - [i18next-http-backend](https://github.com/i18next/i18next-http-backend)
 - [js-file-download](https://github.com/kennethjiang/js-file-download)
 - [leaflet](https://leafletjs.com/)
 - [mui](https://mui.com/)
 - [react](https://reactjs.org/)
 - [react-dom](https://reactjs.org/)
 - [react-hot-toast](https://react-hot-toast.com/)
 - [react-i18next](https://react.i18next.com/)
 - [react-leaflet](https://react-leaflet.js.org/)
 - [react-polling](https://github.com/vivek12345/react-polling)
 - [react-redux](https://react-redux.js.org/)
 - [react-router-dom](https://reactrouter.com/)
 - [react-scripts](https://github.com/facebook/create-react-app/tree/main/packages/react-scripts)
 - [redux](https://react-redux.js.org/)
 - [redux-persist](https://github.com/rt2zz/redux-persist)
 - [seisplotjs](https://github.com/crotwell/seisplotjs)
 - [tailwindcss](https://tailwindcss.com/)
 - [typescript](https://www.typescriptlang.org/)

## License

[The MIT License (MIT)](https://raw.githubusercontent.com/anyshake/observer/master/LICENSE)
