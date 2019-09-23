package usecase

import (
	"encoding/json"
	"fmt"
	"net/http"

	"accu_pls/playlist"
)

const accuURI = "https://www.accuradio.com/playlist/json/"

type HTTPGetter func(uri string) (*http.Response, error)

type AccuPlaylist struct {
	httpGet HTTPGetter
}

func channelToURI(channel string) string {
	return fmt.Sprintf("%s%s/", accuURI, channel)
}

func NewAccuPlaylist(httpGet HTTPGetter) *AccuPlaylist {
	return &AccuPlaylist{
		httpGet: httpGet,
	}
}

func (a *AccuPlaylist) GetTracks(channel string) ([]*playlist.Track, error) {
	resp, err := a.httpGet(channelToURI(channel))
	if err != nil {
		return nil, fmt.Errorf("get playlist for channel %q, error %w", channel, err)
	}
	var rawPL []*rawTrack
	if err := json.NewDecoder(resp.Body).Decode(&rawPL); err != nil {
		return nil, fmt.Errorf("get playlist for channel %q, error %w", channel, err)
	}
	if err := resp.Body.Close(); err != nil {
		return nil, fmt.Errorf("get playlist for channel %q, error %w", channel, err)
	}
	pl := make([]*playlist.Track, 0, len(rawPL))
	for _, t := range rawPL {
		if t.Duration == 0 {
			continue
		}
		pl = append(pl, t.toTrack())
	}
	return pl, nil
}

type rawAlbum struct {
	Title string
	Year  string
}

type rawTrack struct {
	Album       rawAlbum
	TrackArtist string `json:"track_artist"`
	Title,
	Primary,
	Secondary,
	Fn string
	Duration float64
}

func (r *rawTrack) toTrack() *playlist.Track {
	return &playlist.Track{
		Artist:        r.TrackArtist,
		Album:         r.Album.Title,
		Title:         r.Title,
		Duration:      int(r.Duration),
		Year:          r.Album.Year,
		PrimaryLink:   r.Primary + r.Fn + ".m4a",
		SecondaryLink: r.Secondary + r.Fn + ".m4a",
	}
}
