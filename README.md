<p align="center">
  <img src="https://raw.githubusercontent.com/anyshake/observer/master/images/header.png" width="500" alt="banner" />
</p>

[![Codacy Badge](https://app.codacy.com/project/badge/Grade/7b75168a5b03403987122835d74bb448)](https://app.codacy.com/gh/anyshake/observer/dashboard)
[![Downloads](https://img.shields.io/github/downloads/anyshake/observer/total.svg)](https://github.com/anyshake/observer/releases/latest)
[![Go Report Card](https://goreportcard.com/badge/github.com/anyshake/observer)](https://goreportcard.com/report/github.com/anyshake/observer)
[![Build Status](https://github.com/anyshake/observer/actions/workflows/release.yml/badge.svg)](https://github.com/anyshake/observer/actions/workflows/release.yml)
[![Latest Release](https://img.shields.io/github/release/anyshake/observer.svg)](https://github.com/anyshake/observer/releases/latest)

## 🚀 **Join the Open Science Movement!** 🚀

> 🌟 **AnyShake Explorer is currently in pre-launch! We're working hard to bring the **AnyShake Project** to life, but we need your support to reach our goal of 200 subscribers.** 🌟

> 👉 **[Go and subscribe us on Crowd Supply](https://www.crowdsupply.com/senseplex/anyshake-explorer)** and be a part of the first wave of users to get hands-on with this revolutionary open-source seismic monitoring system!

> Also help us spread the word by sharing this page: 👉 **[www.crowdsupply.com/senseplex/anyshake-explorer](https://www.crowdsupply.com/senseplex/anyshake-explorer)**

---

## Overview

**AnyShake Observer** is the companion software for [AnyShake Explorer](https://github.com/anyshake/explorer), the world’s first fully open-source, high-precision seismic monitoring system. It is a cross-platform, web-based application designed to visualize, archive, and export seismic data in real time.

AnyShake Observer is written in **Go** and **TypeScript**, supporting a wide range of OS and CPU architectures, including embedded Linux. It offers real-time waveform display, event analysis, and database archiving, with a strong focus on professional usability and extensibility.

It works seamlessly with **AnyShake Explorer** over a serial connection and supports exporting data in standard seismic formats such as **SAC** and **MiniSEED**, as well as **SeedLink** streaming for networked data systems.

## Features

- 📦 **Single-binary deployment** – fast, simple, and cross-platform
- 🖥️ **Web-based user interface** – no client installation required
- 📱 **Responsive design** – optimized for desktop, tablet, and mobile
- 🌍 **Multi-language support** – auto-detects browser language
- ⏱️ **Real-time waveform display** – view live seismic data from your device
- 🎛️ **Multi-channel support** – visualize multiple sensors (e.g. geophone and accelerometer)
- 📊 **Historical waveform queries** – search by time or global seismic events
- 🌐 **Sharable waveform links** – share analysis results with a single URL
- 📸 **Daily helicorder generation** – auto-generated visual timeline of activity
- 🚨 **QuakeSense service** – built-in earthquake detection engine (STA/LTA and Z-Detect methods)
- 📁 **Data export** – save data as **MiniSEED**, **SAC**, **TXT**, or **WAV**
- 🔁 **Streaming & forwarding** – supports **SeedLink** and **TCP** protocols
- 🧩 **Flexible storage** – compatible with PostgreSQL, MariaDB/MySQL, SQL Server, and SQLite
- 🔗 **Seamless SeisComP integration** – easily connect to professional seismic networks
- 🚀 **... and more!** – with active development and community-driven features

## Documentation

Start here 👉 [https://anyshake.org/docs/introduction](https://anyshake.org/docs/introduction)

## Preview

Here are some screenshots showcasing the key features of **AnyShake Observer**.

### Home Dashboard

The Home Dashboard provides a concise overview of the station and device status, including current location, service module health, link connectivity, and real-time system statistics.

<img src="https://raw.githubusercontent.com/anyshake/observer/master/images/key-features/home.webp" width="600" alt="Home Dashboard" />

### Realtime Waveform

The real-time waveform view displays seismic data and the current sample rate from your AnyShake Explorer device in a web-based interface. It supports zooming, panning, and customizable channel layouts. Layout configurations are persistent and can be locked to prevent accidental changes.

<img src="https://raw.githubusercontent.com/anyshake/observer/master/images/key-features/realtime.webp" width="600" alt="Realtime View" />

### SeedLink Streaming

The AnyShake team independently developed a SeedLink protocol implementation in pure Go ([github.com/bclswl0827/slgo](https://github.com/bclswl0827/slgo)), enabling native SeedLink services without relying on RingServer or SeisComP. This allows seamless integration with tools like Swarm and ObsPy.

<img src="https://raw.githubusercontent.com/anyshake/observer/master/images/key-features/seedlink.webp" width="600" alt="SeedLink View" />

### Historical Data Query

The historical query feature lets users retrieve waveforms from the database by specifying a time range (up to 1 hour). It also integrates global seismic agency data for reverse earthquake lookup. Retrieved waveform parts can be exported in multiple formats, including MiniSEED, SAC, TXT, and WAV audio.

<img src="https://raw.githubusercontent.com/anyshake/observer/master/images/key-features/history.webp" width="600" alt="History Query" />

Exported data formats, such as MiniSEED, are fully compatible with third-party analysis tools like Swarm. The image below shows a seismic event in Myanmar on Mar 28th, 2025, detected from over 2,400 kilometers away using AnyShake Explorer — demonstrating its remarkable sensitivity and real-world performance, on par with many proprietary, closed-source systems.

<img src="https://raw.githubusercontent.com/anyshake/observer/master/images/key-features/miniseed.webp" width="600" alt="MiniSEED View" />

### Data Download

The Data Download page allows users to access and download daily archived MiniSEED files and helicorder images directly from disk for extended archiving or offline analysis.

<img src="https://raw.githubusercontent.com/anyshake/observer/master/images/key-features/download.webp" width="600" alt="Data Download" />

### Station Metadata

Station metadata files for AnyShake Explorer devices — including instrument response (poles and zeros) are available in both SeisComP XML and FDSNWS StationXML formats. Users can easily copy the content to the clipboard or download it as a file.

<img src="https://raw.githubusercontent.com/anyshake/observer/master/images/key-features/metadata.webp" width="600" alt="Station Metadata" />

### Service Control

A rich set of modular service modules are integrated, each running independently to ensure stability and flexibility. Users can enable, disable, or configure these modules directly through the web interface. The system also includes earthquake detection services using STA/LTA and Z-Detect algorithms (QuakeSense Service). Upon detecting seismic activity, it can push alerts in real time via MQTT protocol.

<img src="https://raw.githubusercontent.com/anyshake/observer/master/images/key-features/service.webp" width="600" alt="Service Control" />

### User Management

The system supports both administrator and general user roles, making it adaptable for scenarios where multiple individuals manage and operate the site.

<img src="https://raw.githubusercontent.com/anyshake/observer/master/images/key-features/users.webp" width="600" alt="User Management" />

## Credits

The project is maintained by **SensePlex Limited**, a UK-based company dedicated to developing open-source hardware and software.

All rights, including the right of interpretation, reproduction, distribution, and commercial use, are reserved by **SensePlex Limited**.

For any inquiries, please refer to the contact information provided on the [AnyShake GitHub organization page](https://github.com/anyshake).

![Star History Chart](https://api.star-history.com/svg?repos=anyshake/observer&type=Date)
