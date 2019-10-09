package main

import (
	"log"
	"net/http"
	"os"
	"errors"

	playlist "accu_pls/playlist/usecase"
	"accu_pls/repository"
	"accu_pls/repository/sqlite"
)

const (
	defDuration = 86400
	maxDryRuns = 64
)

func main() {
	if len(os.Args) != 2 {
		log.Fatal("channel unspecified, bye")
	}
	channel := os.Args[1]
	l := playlist.NewAccuPlaylist(http.Get)
	r, err := sqlite.NewSqliteRepo()
	if err != nil {
		log.Fatal(err)
	}
	defer r.Close()
	dryRuns := 0
	for {
		tracks, err := l.GetTracks(channel)
		if err != nil {
			log.Printf("error fetching tracks %v", err)
		}
		if err := r.Save(channel, tracks); err != nil {
			if errors.Is(err, repository.ErrNoTracksSaved) {
				dryRuns++
				if dryRuns > maxDryRuns {
					break
				}
				continue
			}
			log.Fatal(err)
		}
		dryRuns = 0
	}
}
