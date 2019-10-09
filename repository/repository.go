package repository

import (
	"errors"

	"accu_pls/playlist"
)

var ErrNoTracksSaved = errors.New("no tracks saved")

type Repository interface {
	Save(channel string, tracks []*playlist.Track) error
	Exists(track *playlist.Track) (bool, error)
}
