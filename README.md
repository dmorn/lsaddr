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
 $ bin/lsaddr /Applications/Spotify.app
{Spotify 192.168.0.61:59053->104.199.65.114:4070}
{Spotify 192.168.0.61:60237->192.121.140.177:80}
{Spotify 192.168.0.61:60566->35.186.224.53:443}
{Spotify 192.168.0.61:59062->35.186.224.47:443}
{Spotify 192.168.0.61:60938->151.101.12.246:443}
{Spotify 192.168.0.61:60939->151.101.12.246:443}
{Spotify 192.168.0.61:60943->151.101.12.246:443}
{Spotify 192.168.0.61:60940->151.101.12.246:443}
{Spotify 192.168.0.61:60941->151.101.12.246:443}
{Spotify 192.168.0.61:60942->151.101.12.246:443}
{Spotify 192.168.0.61:60944->151.101.12.246:443}
{Spotify 192.168.0.61:60945->151.101.12.246:443}
{Spotify 192.168.0.61:60946->151.101.12.246:443}
{Spotify 192.168.0.61:60947->151.101.12.246:443}
{Spotify 192.168.0.61:60949->151.101.12.246:443}
{Spotify 192.168.0.61:60950->151.101.12.246:443}
{Spotify 192.168.0.61:60951->151.101.12.246:443}
{Spotify 192.168.0.61:60952->151.101.12.246:443}
{Spotify 192.168.0.61:60953->151.101.12.246:443}
{Spotify 192.168.0.61:60954->151.101.12.246:443}
{Spotify 192.168.0.61:60955->151.101.12.246:443}
{Spotify 192.168.0.61:60977->151.101.12.246:443}
{Spotify 192.168.0.61:60957->151.101.12.246:443}
{Spotify 192.168.0.61:60958->151.101.12.246:443}
{Spotify 192.168.0.61:60959->151.101.12.246:443}
```
