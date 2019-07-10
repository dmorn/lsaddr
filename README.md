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
Choose one
- `go get -u github.com/booster-proj/lsaddr`
- download the [downloader script](https://raw.githubusercontent.com/booster-proj/lsaddr/master/godownloader.sh) and execute it (you can also specify the version that you want to download as argument)
- download your favourite release from the [releases section](https://github.com/booster-proj/lsaddr/releases)

Big thanks to [goreleaser](https://github.com/goreleaser/goreleaser) and [godownloader](https://github.com/goreleaser/godownloader) which made the releasing process **FUN**! 🤩

## Usage
The idea is to easily filter the list of open network files of a specific application. The list is filtered with a regular expression: only
the lines that match against it are kept, the others discarded. You can pass to `lsaddr` either directly the regex, or the root folder of the
target app (supported only on macOS for now). Check out some examples:

```
lsaddr (master) $ bin/lsaddr Spotify --out=csv
COMMAND,NET,SRC,DST
Spotify,tcp,192.168.0.98:59315,104.199.64.253:80
Spotify,udp,*:57621,
Spotify,tcp,*:57621,
Spotify,tcp,*:61357,
Spotify,tcp,192.168.0.98:61340,2.16.186.91:80
Spotify,udp,*:1900,
Spotify,udp,*:63319,
Spotify,udp,*:55092,
Spotify,tcp,192.168.0.98:61342,2.16.186.11:443
Spotify,tcp,192.168.0.98:61360,192.121.140.177:443
Spotify,tcp,192.168.0.98:61344,151.101.112.246:443
Spotify,tcp,192.168.0.98:61343,151.101.112.246:443
Spotify,tcp,192.168.0.98:61350,151.101.112.246:443
Spotify,tcp,192.168.0.98:61345,151.101.112.246:443
Spotify,tcp,192.168.0.98:61346,151.101.112.246:443
Spotify,tcp,192.168.0.98:61347,151.101.112.246:443
Spotify,tcp,192.168.0.98:59355,35.186.224.53:443
Spotify,tcp,192.168.0.98:59360,35.186.224.47:443
Spotify,tcp,192.168.0.98:61348,151.101.112.246:443
Spotify,tcp,192.168.0.98:61349,151.101.112.246:443
Spotify,tcp,192.168.0.98:61351,151.101.112.246:443
Spotify,tcp,192.168.0.98:61352,151.101.112.246:443
Spotify,tcp,192.168.0.98:61353,151.101.112.246:443
Spotify,tcp,192.168.0.98:61361,192.121.140.177:443
```
Note: "Spotify" is used as a regular expression.
```
lsaddr (master) $ bin/lsaddr /Applications/Spotify.app/ --out=csv
COMMAND,NET,SRC,DST
Spotify,tcp,192.168.0.98:59315,104.199.64.253:80
Spotify,udp,*:57621,
Spotify,tcp,*:57621,
Spotify,tcp,*:61357,
Spotify,tcp,192.168.0.98:61340,2.16.186.91:80
Spotify,udp,*:1900,
Spotify,udp,*:63319,
Spotify,udp,*:55092,
Spotify,tcp,192.168.0.98:61344,151.101.112.246:443
Spotify,tcp,192.168.0.98:61343,151.101.112.246:443
Spotify,tcp,192.168.0.98:61350,151.101.112.246:443
Spotify,tcp,192.168.0.98:61345,151.101.112.246:443
Spotify,tcp,192.168.0.98:61346,151.101.112.246:443
Spotify,tcp,192.168.0.98:61347,151.101.112.246:443
Spotify,tcp,192.168.0.98:59355,35.186.224.53:443
Spotify,tcp,192.168.0.98:59360,35.186.224.47:443
Spotify,tcp,192.168.0.98:61348,151.101.112.246:443
Spotify,tcp,192.168.0.98:61349,151.101.112.246:443
Spotify,tcp,192.168.0.98:61351,151.101.112.246:443
Spotify,tcp,192.168.0.98:61352,151.101.112.246:443
Spotify,tcp,192.168.0.98:61353,151.101.112.246:443
```
Note: "/Applications/Spotify.app" is used to find the application's name, then its
process identifiers are used to build the regular expression.
```
lsaddr (master) $ bin/lsaddr /Applications/Spotify.app/ --out=csv --debug
[lookup] app name: Spotify, path: /Applications/Spotify.app
[lsaddr] # of open files: 21
COMMAND,NET,SRC,DST
Spotify,tcp,192.168.0.98:59315,104.199.64.253:80
Spotify,udp,*:57621,
Spotify,tcp,*:57621,
Spotify,tcp,*:61357,
Spotify,tcp,192.168.0.98:61340,2.16.186.91:80
Spotify,udp,*:1900,
Spotify,udp,*:63319,
Spotify,udp,*:55092,
Spotify,tcp,192.168.0.98:61344,151.101.112.246:443
Spotify,tcp,192.168.0.98:61343,151.101.112.246:443
Spotify,tcp,192.168.0.98:61350,151.101.112.246:443
Spotify,tcp,192.168.0.98:61345,151.101.112.246:443
Spotify,tcp,192.168.0.98:61346,151.101.112.246:443
Spotify,tcp,192.168.0.98:61347,151.101.112.246:443
Spotify,tcp,192.168.0.98:59355,35.186.224.53:443
Spotify,tcp,192.168.0.98:59360,35.186.224.47:443
Spotify,tcp,192.168.0.98:61348,151.101.112.246:443
Spotify,tcp,192.168.0.98:61349,151.101.112.246:443
Spotify,tcp,192.168.0.98:61351,151.101.112.246:443
Spotify,tcp,192.168.0.98:61352,151.101.112.246:443
Spotify,tcp,192.168.0.98:61353,151.101.112.246:443
```
Note: `--debug` information is printed to `stderr`, command's output to `stdout`.
```
lsaddr (master) $ bin/lsaddr /Applications/Spotify.app/ --out=bpf
host 104.199.64.253 or 2.16.186.91 or 151.101.112.246 or 35.186.224.53 or 35.186.224.47
```
Notes:
- you can encode the output either in csv or as a [bpf](https://en.wikipedia.org/wiki/Berkeley_Packet_Filter) (hint: very useful for packet capturing tools). 
- only the unique destination addresses are taken into consideration when building the filter,
ignoring the ports and without specifing if the "direction" (incoming or outgoing) that we want to
filter. This is because the expected behaviour has not yet been defined.
