package playlist

type Track struct {
	Artist,
	Album,
	Title string
	Duration int
	Year,
	PrimaryLink,
	SecondaryLink string
}

type Playlist interface {
	GetTracks(channel string) ([]*Track, error)
}
