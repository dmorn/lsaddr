# lsaddr
[![GoDoc](https://godoc.org/github.com/booster-proj/lsaddr?status.svg)](https://godoc.org/github.com/booster-proj/lsaddr)
[![Go Report Card](https://goreportcard.com/badge/github.com/booster-proj/lsaddr)](https://goreportcard.com/report/github.com/booster-proj/lsaddr)
[![Release](https://img.shields.io/github/release/booster-proj/lsaddr.svg)](https://github.com/booster-proj/lsaddr/releases/latest)

## Before we start
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
The idea is to easily filter the list of open network files of a specific application. The list is filtered with a regular expression: only
the lines that match against it are kept, the others discarded. You can pass to `lsaddr` either directly the regex, or the root folder of the
target app (supported only on macOS for now). Check out some examples:

```
lsaddr (master) $ bin/lsaddr Spotify
COMMAND,NET,SRC,DST
Spotify,tcp,192.168.0.61:52439,104.199.64.103:4070
Spotify,tcp,192.168.0.61:54902,35.186.224.53:443
Spotify,tcp,192.168.0.61:52456,35.186.224.47:443
```
Note: "Spotify" is used as a regular expression.
```
lsaddr (master) $ bin/lsaddr /Applications/Spotify.app
COMMAND,NET,SRC,DST
Spotify,tcp,192.168.0.61:52439,104.199.64.103:4070
Spotify,tcp,192.168.0.61:54902,35.186.224.53:443
Spotify,tcp,192.168.0.61:52456,35.186.224.47:443
```
Note: "/Applications/Spotify.app" is used to find the application's name, then its
process identifiers are used to build the regular expression.
```
lsaddr (master) $ bin/lsaddr /Applications/Spotify.app --debug
[lookup] app name: Spotify, path: /Applications/Spotify.app
[lsaddr] # of open files: 9
[lsaddr] skipping open file: *:57621->
[lsaddr] skipping open file: *:57621->
[lsaddr] skipping open file: *:53104->
[lsaddr] skipping open file: *:1900->
[lsaddr] skipping open file: *:57313->
[lsaddr] skipping open file: *:53592->
COMMAND,NET,SRC,DST
Spotify,tcp,192.168.0.61:52439,104.199.64.103:4070
Spotify,tcp,192.168.0.61:54902,35.186.224.53:443
Spotify,tcp,192.168.0.61:52456,35.186.224.47:443
```
Note: `--debug` information is printed to `stderr`, command's output to `stdout`.
