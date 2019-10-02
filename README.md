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
- `windows` (**NEW** 💥)

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
The idea is to easily filter the list of open network files of a specific application. The list is filtered with a regular expression: only
the lines that match against it are kept, the others discarded. You can pass to `lsaddr` either directly the regex, or the root folder of the
target app (supported only on macOS for now). Check out some examples:


#### Example #1
"Spotify" is used as a regular expression.
```
$ bin/lsaddr Spotify
COMMAND,NET,SRC,DST
Spotify,tcp,192.168.0.98:54862,104.199.64.69:4070
Spotify,tcp,*:57621,
Spotify,tcp,*:54850,
Spotify,udp,*:57621,
Spotify,udp,*:1900,
Spotify,udp,*:61152,
Spotify,udp,*:51535,
Spotify,tcp,192.168.0.98:54878,35.186.224.47:443
Spotify,tcp,192.168.0.98:54872,35.186.224.53:443
```

#### Example #2
"/Applications/Spotify.app" is used to find the application's name, then its
process identifiers are used to build the regular expression.
```
$ bin/lsaddr /Applications/Spotify.app/
COMMAND,NET,SRC,DST
Spotify,tcp,192.168.0.98:54862,104.199.64.69:4070
Spotify,tcp,*:57621,
Spotify,tcp,*:54850,
Spotify,udp,*:57621,
Spotify,udp,*:1900,
Spotify,udp,*:61152,
Spotify,udp,*:51535,
Spotify,tcp,192.168.0.98:54878,35.186.224.47:443
Spotify,tcp,192.168.0.98:54872,35.186.224.53:443
```

#### Example #3
`--debug` information is printed to `stderr`, command's output to `stdout`.
```
$ bin/lsaddr /Applications/Spotify.app/ --debug
[lsaddr] 2019/07/12 14:29:50 app name: Spotify, path: /Applications/Spotify.app
[lsaddr] 2019/07/12 14:29:50 regexp built: "48042|48044|48045|48047"
[lsaddr] 2019/07/12 14:29:50 # of open files: 9
COMMAND,NET,SRC,DST
Spotify,tcp,192.168.0.98:54862,104.199.64.69:4070
Spotify,tcp,*:57621,
Spotify,tcp,*:54850,
Spotify,udp,*:57621,
Spotify,udp,*:1900,
Spotify,udp,*:61152,
Spotify,udp,*:51535,
Spotify,tcp,192.168.0.98:54878,35.186.224.47:443
Spotify,tcp,192.168.0.98:54872,35.186.224.53:443
```

#### Example #4
- you can encode the output either in csv or as a [bpf](https://en.wikipedia.org/wiki/Berkeley_Packet_Filter) (hint: very useful for packet capturing tools). 
- only the unique destination addresses are taken into consideration when building the filter,
ignoring the ports and without specifing if the "direction" (incoming or outgoing) that we want to
filter. This is because the expected behaviour has not yet been defined.
```
$ bin/lsaddr /Applications/Mail.app --out=bpf
(tcp and host 192.168.0.98 and port 58100) or (tcp and host 64.233.184.108 and port 993) or (tcp and host 192.168.0.98 and port 58100) or (tcp and host 64.233.184.108 and port 993) or (tcp and host 192.168.0.98 and port 57213) or (tcp and host 10.0.0.1 and port 993) or (tcp and host 192.168.0.98 and port 57213) or (tcp and host 10.0.0.1 and port 993) or (tcp and host 192.168.0.98 and port 57214) or (tcp and host 10.0.0.1 and port 993) or (tcp and host 192.168.0.98 and port 57214) or (tcp and host 10.0.0.1 and port 993) or (tcp and host 192.168.0.98 and port 57216) or (tcp and host 17.56.136.197 and port 993) or (tcp and host 192.168.0.98 and port 57216) or (tcp and host 17.56.136.197 and port 993) or (tcp and host 192.168.0.98 and port 57217) or (tcp and host 17.56.136.197 and port 993) or (tcp and host 192.168.0.98 and port 57217) or (tcp and host 17.56.136.197 and port 993)
```
#### Example #5
At the moment on Windows you can pass the absulute path of the program you want (or straight `<program>.exe`)
to analyze.
```
> lsaddr.exe "chrome.exe"
COMMAND,NET,SRC,DST
chrome.exe,tcp,10.211.55.3:50551,216.58.205.163:443
chrome.exe,tcp,10.211.55.3:50556,216.58.205.195:443
chrome.exe,tcp,10.211.55.3:50558,216.58.205.67:443
chrome.exe,tcp,10.211.55.3:50567,216.58.205.106:443
chrome.exe,udp,0.0.0.0:5353,*:*
chrome.exe,udp,0.0.0.0:5353,*:*
chrome.exe,udp,0.0.0.0:5353,*:*
chrome.exe,udp,[::]:5353,*:*
chrome.exe,udp,[::]:5353,*:*
```
