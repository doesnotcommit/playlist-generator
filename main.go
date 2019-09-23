package main

import (
	"io"
	"log"
	"net/http"
	"os"

	playlist "accu_pls/playlist/usecase"
	pls "accu_pls/format/pls"
)

const defDuration = 86400

func main() {
	if len(os.Args) != 2 {
		log.Fatal("channel unspecified, bye")
	}
	l := playlist.NewAccuPlaylist(http.Get)
	pls := pls.NewPls(l)
	io.Copy(os.Stdout, pls.GetReader(os.Args[1], defDuration))
}
