# msfs2020-go/vfrmap

local web-server using msfs2020-go/simconnect that will allow you to view your location, and some information about your trajectory including airspeed and altitude.

also allows you to quickly teleport your plane to any location.

## install

* download latest release zip [here](https://github.com/lian/msfs2020-go/releases)
* unzip `vfrmap-win64.zip`

## run
* run `vfrmap.exe`
* browse to http://localhost:9000
* or to `http://<computer-ip>:9000`

## arguments

* `-v` show program version
* `-verbose` verbose output
* `-disable-teleport` disables teleport

## usage

* clicking on your plane to see gps coordinates, follow or don't follow the plane, or open the current location on google maps in a new tab.
* clicking on the map itself to create a marker. clicking on that marker allows you to teleport to this location or enter your own gps coordinates.
* dragging the map stops following the plane.
* pressing escape key switches between following the plane or freely moving around on the map.
* clicking on the top right corner hides the HUD

## change visualisation

if you want to change how the webpage looks then copy and change [index.html](html/index.html) to the same folder as `vfrmap.exe` and relaunch the program.

## openstreetmap

earlier versions of this app used google maps directly, but this was too expensive. openstreetmap is free to use and very good as well.

## compile

`GOOS=windows GOARCH=amd64 go build github.com/lian/msfs2020-go/vfrmap` or see [build-vfrmap.sh](https://github.com/lian/msfs2020-go/blob/master/build-vfrmap.sh)

## Why does my virus-scanning software think this program is infected?

From official golang website https://golang.org/doc/faq#virus

"This is a common occurrence, especially on Windows machines, and is almost always a false positive. Commercial virus scanning programs are often confused by the structure of Go binaries, which they don't see as often as those compiled from other languages."

## screenshots

![screenshot](https://i.imgur.com/n9vHln8.png)
![osm with openAIP](https://s3.eu-central-1.amazonaws.com/sh4re/2020-08-26_19_37_05_scrot.png)
![stamen watercolor with openAIP](https://s3.eu-central-1.amazonaws.com/sh4re/2020-08-26_19_36_03_scrot.png)
![stamen toner](https://s3.eu-central-1.amazonaws.com/sh4re/2020-08-26_19_37_32_scrot.png)
![screenshot](https://i.imgur.com/5PZyKC8.png)
