package pls

import (
	"fmt"
	"io"

	"accu_pls/playlist"
)

const header = "[playlist]\n"
const trackTemplate = "File%[1]d=%[2]s\nTitle%[1]d=%[3]s - %[4]s [%[5]s] - %[6]s\nLength%[1]d=%[7]d\n"

type Pls struct {
	l playlist.Playlist
}

func NewPls(l playlist.Playlist) *Pls {
	return &Pls{
		l: l,
	}
}

func (p *Pls) GetReader(channel string, minDuration int) io.Reader {
	r, w := io.Pipe()
	go func() {
		w.Write([]byte(header))
		trackCount := 1
		for totalDuration := 0; totalDuration < minDuration; {
			tracks, err := p.l.GetTracks(channel)
			if err != nil {
				w.CloseWithError(fmt.Errorf("error getting tracks from channel %s: %w", channel, err))
				return
			}
			for _, track := range tracks {
				w.Write([]byte(fmt.Sprintf(trackTemplate,
					trackCount,
					track.PrimaryLink,
					track.Artist,
					track.Album,
					track.Year,
					track.Title,
					track.Duration,
				)))
				totalDuration += track.Duration
				trackCount++
			}
		}
		w.Write([]byte(fmt.Sprintf("NumberOfEntries=%d\nVersion=2\n", trackCount-1)))
		w.Close()
	}()
	return r
}
