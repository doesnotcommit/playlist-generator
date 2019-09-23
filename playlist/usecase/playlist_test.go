package usecase

import (
	"bytes"
	"io/ioutil"
	"net/http"
	"testing"

	"accu_pls/playlist"
)

var playlistTests = []struct {
	name         string
	channel      string
	mockResponse *http.Response
	expected     []*playlist.Track
}{
	{
		name:    "test playlist fetcher",
		channel: "deadbeef",
		mockResponse: &http.Response{
			Body: ioutil.NopCloser(bytes.NewReader([]byte(`[
				{
					"album": {
						"title": "album_title",
						"year": "9999"
					},
					"track_artist": "track_artist",
					"title": "track_title",
					"primary": "https://foo.bar/",
					"secondary": "https://foo.baz/",
					"fn": "filename",
					"duration": 222.111
				}
			]`))),
		},
		expected: []*playlist.Track{
			&playlist.Track{
				Artist:        "track_artist",
				Album:         "album_title",
				Title:         "track_title",
				Duration:      222,
				Year:          "9999",
				PrimaryLink:   "https://foo.bar/filename.m4a",
				SecondaryLink: "https://foo.baz/filename.m4a",
			},
		},
	},
}

func NewMockHTTPGetter(r *http.Response) HTTPGetter {
	return func(uri string) (*http.Response, error) {
		return r, nil
	}
}

func Test_Playlist(t *testing.T) {
	for _, test := range playlistTests {
		t.Log(test)
		t.Run(test.name, func(t *testing.T) {
			ap := NewAccuPlaylist(NewMockHTTPGetter(test.mockResponse))
			tracks, err := ap.GetTracks(test.channel)
			if err != nil {
				t.Errorf("get tracks error %w", err)
			}
			for i, track := range tracks {
				if *track != *test.expected[i] {
					t.Errorf("expectation failed:\n%v\n%v", *track, *test.expected[i])
				}
			}
		})
	}
}
