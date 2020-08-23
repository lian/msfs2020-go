# msfs2020-go/vfrmap

web server that shows your current MSFS2020 plane position in google maps inside the browser

## install

* download latest release zip [here](https://github.com/lian/msfs2020-go/releases)
* unzip `vfrmap-win64.zip`

## run
* run `vfrmap.exe`
* browse to http://localhost:9000

## arguments

* `-v` show program version
* `-api-key` use your own gmap api-key
* `-verbose` verbose output

## change visualisation

if you want to change how the webpage looks then copy and change [index.html](html/index.html) to the same folder as `vfrmap.exe` and relaunch the program.

## compile

`GOOS=windows GOARCH=amd64 go build github.com/lian/msfs2020-go/vfrmap` or see [build-vfrmap.sh](https://github.com/lian/msfs2020-go/blob/master/build-vfrmap.sh)

## screenshots

![screenshot](https://i.imgur.com/YllMEvG.png)
