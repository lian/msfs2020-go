package leafletjs

import (
	"net/http"
)

//go:generate go-bindata -pkg leafletjs -o bindata.go -modtime 1 -prefix "../../../_vendor/leafletjs" "../../../_vendor/leafletjs" "../../../_vendor/leafletjs/images"

type FS struct {
}

func (_ FS) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	switch r.URL.Path {
	case "leaflet.css":
		w.Header().Set("Content-Type", "text/css")
		w.Write(MustAsset("leaflet.css"))
	case "leaflet.js":
		w.Header().Set("Content-Type", "text/javascript")
		w.Write(MustAsset("leaflet.js"))
	case "leaflet.rotatedMarker.js":
		w.Header().Set("Content-Type", "text/javascript")
		w.Write(MustAsset("leaflet.rotatedMarker.js"))
	case "images/layers-2x.png":
		w.Header().Set("Content-Type", "image/png")
		w.Write(MustAsset("images/layers-2x.png"))
	case "images/layers.png":
		w.Header().Set("Content-Type", "image/png")
		w.Write(MustAsset("images/layers.png"))
	case "images/marker-icon-2x.png":
		w.Header().Set("Content-Type", "image/png")
		w.Write(MustAsset("images/marker-icon-2x.png"))
	case "images/marker-icon.png":
		w.Header().Set("Content-Type", "image/png")
		w.Write(MustAsset("images/marker-icon.png"))
	case "images/marker-shadow.png":
		w.Header().Set("Content-Type", "image/png")
		w.Write(MustAsset("images/marker-shadow.png"))
	}
}
