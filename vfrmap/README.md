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

## usage

click on your plane to see gps coordinates, follow or don't follow the plane, or open the current location on google maps in a new tab.

if you click on the map itself a new marker appears. clicking on that marker allows you to teleport to this location or enter your own gps coordinates.

esc key switching between following the plane or freely moving around on the map.

## change visualisation

if you want to change how the webpage looks then copy and change [index.html](html/index.html) to the same folder as `vfrmap.exe` and relaunch the program.

## openstreetmap

earlier versions of this app used google maps directly, but this was too expensive. openstreetmap is free to use and very good as well.

## compile

`GOOS=windows GOARCH=amd64 go build github.com/lian/msfs2020-go/vfrmap` or see [build-vfrmap.sh](https://github.com/lian/msfs2020-go/blob/master/build-vfrmap.sh)

## screenshots

![screenshot](https://i.imgur.com/5PZyKC8.png)
