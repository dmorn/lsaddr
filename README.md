# lsaddr
[![GoDoc](https://godoc.org/github.com/booster-proj/lsaddr?status.svg)](https://godoc.org/github.com/booster-proj/lsaddr)
[![Go Report Card](https://goreportcard.com/badge/github.com/booster-proj/lsaddr)](https://goreportcard.com/report/github.com/booster-proj/lsaddr)
[![Release](https://img.shields.io/github/release/booster-proj/lsaddr.svg)](https://github.com/booster-proj/lsaddr/releases/latest)

### Notes
Note that this tool is under development, and we plan to make it as much portable as possible.

#### Supported OS
- `macOS`
- `linux`

#### Dependencies
```
$ lsof -v
lsof version information:
    revision: 4.89
    latest revision: ftp://lsof.itap.purdue.edu/pub/tools/unix/lsof/
    latest FAQ: ftp://lsof.itap.purdue.edu/pub/tools/unix/lsof/FAQ
    latest man page: ftp://lsof.itap.purdue.edu/pub/tools/unix/lsof/lsof_man
    configuration info: libproc-based
    Anyone can list all files.
    /dev warnings are disabled.
    Kernel ID check is disabled.
```
## Installation
- `go get -u github.com/booster-proj/lsaddr`
- downloading the [downloader script](https://raw.githubusercontent.com/booster-proj/lsaddr/master/godownloader.sh) and executing it (you can also specify the version that you want to download as argument)
- downloading your favourite release from the [releases section](https://github.com/booster-proj/lsaddr/releases)

Big thanks to [goreleaser](https://github.com/goreleaser/goreleaser) and [godownloader](https://github.com/goreleaser/godownloader) which made the releasing process **FUN**! 🤩

## Usage
The idea is to drag-and-drop you application to `lsaddr`, and it displays the network addresses that that app is using. We plan to make the output configurable, so it is easy to consume it from other programs, for example by allowing to specify output's encoding and fields.

```
lsaddr (master) $ bin/lsaddr Spotify
192.168.0.61:49973->2.16.106.146:80
192.168.0.61:49501->192.121.140.177:80
192.168.0.61:49235->104.199.64.158:4070
192.168.0.61:49252->35.186.224.53:443
192.168.0.61:49671->35.186.224.47:443
192.168.0.61:49974->2.16.186.11:80
```
```
lsaddr (master) $ bin/lsaddr /Applications/Spotify.app
192.168.0.61:49973->2.16.106.146:80
192.168.0.61:49501->192.121.140.177:80
192.168.0.61:49235->104.199.64.158:4070
192.168.0.61:49252->35.186.224.53:443
192.168.0.61:49671->35.186.224.47:443
192.168.0.61:49974->2.16.186.11:80
```
```
lsaddr (master) $ bin/lsaddr /Applications/Spotify.app --debug
[lookup] app name: Spotify, path: /Applications/Spotify.app
[lsaddr] # of open files: 12
[lsaddr] skipping open file: *:57621->
[lsaddr] skipping open file: *:57621->
[lsaddr] skipping open file: *:50086->
192.168.0.61:49973->2.16.106.146:80
192.168.0.61:49501->192.121.140.177:80
[lsaddr] skipping open file: *:1900->
[lsaddr] skipping open file: *:58304->
[lsaddr] skipping open file: *:62516->
192.168.0.61:49235->104.199.64.158:4070
192.168.0.61:49252->35.186.224.53:443
192.168.0.61:49671->35.186.224.47:443
192.168.0.61:49974->2.16.186.11:80
```
Note: `--debug` information is printed to `stderr`, command's output to `stdout`.
