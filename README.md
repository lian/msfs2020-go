# msfs2020-go

simconnect package [msfs2020-go/simconnect](simconnect/) connects to microsoft flight simulator 2020 using golang.

cross-compiles from macos/linux, no other dependencies required. produces a single binary with no other files or configuration required.

## status

[msfs2020-go/simconnect](simconnect/) package currently only implements enough of the simconnect api for [examples](examples/) and [vfrmap](vfrmap).

## download

program zips are uploaded [here](https://github.com/lian/msfs2020-go/releases)

## tools

* [vfrmap](vfrmap/) local web-server that will allow you to view your location, and some information about your trajectory including airspeed and altitude.

## examples

* [examples/request_data](examples/request_data/) port of `MSFS-SDK/Samples/SimConnectSamples/RequestData/RequestData.cpp`

