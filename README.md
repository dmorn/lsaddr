# lsaddr
[![GoDoc](https://godoc.org/github.com/booster-proj/lsaddr?status.svg)](https://godoc.org/github.com/booster-proj/lsaddr)
[![Go Report Card](https://goreportcard.com/badge/github.com/booster-proj/lsaddr)](https://goreportcard.com/report/github.com/booster-proj/lsaddr)
[![Release](https://img.shields.io/github/release/booster-proj/lsaddr.svg)](https://github.com/booster-proj/lsaddr/releases/latest)
[![Build Status](https://travis-ci.org/jecoz/lsaddr.svg?branch=master)](https://travis-ci.org/jecoz/lsaddr)
[![Reviewed by Hound](https://img.shields.io/badge/Reviewed_by-Hound-8E64B0.svg)](https://houndci.com)

## Before we start
#### Supported OS
- `macOS`
- `linux`
- `windows

#### External dependencies
OS | Dep | Notes
------|------|------
**macOS** | `lsof` | (tested revision: 4.89)
**macOS** | `pgrep` |
**Linux** | `lsof` |
**Windows** | `netstat` |
**Windows** | `tasklist` |

## Installation
Choose one
* $ `go get -u github.com/booster-proj/lsaddr`
* $ `bash <(curl -s https://raw.githubusercontent.com/booster-proj/lsaddr/master/install.sh)`
* download your favourite release from the [releases section](https://github.com/booster-proj/lsaddr/releases)

Big thanks to [goreleaser](https://github.com/goreleaser/goreleaser) and [godownloader](https://github.com/goreleaser/godownloader) which made the releasing process **FUN**! 🤩

## Usage
The idea is to easily filter the list of open network files of a specific application. The list is filtered with a regular expression: only the lines that match against it are kept, the others discarded.

## Examples
#### Find connections opened by "Spotify"
```
% bin/lsaddr Spotify
PID,CMD,NET,SRC,DST
62822,Spotify,tcp,10.7.152.118:52213,104.199.64.50:80
62822,Spotify,tcp,10.7.152.118:52255,35.186.224.47:443
62826,Spotify,tcp,10.7.152.118:52196,35.186.224.53:443
```

#### Increment verbosity (debugging)
Note: `debug` information is printed to `stderr`, command's output to `stdout`.
```
% bin/lsaddr Spotify --verbose
...
```
(output omitted for readiness)

#### Dump Spotify's network traffic using tcpdump
```
% bin/lsaddr -f bpf Spotify | xargs -0 sudo tcpdump
```
