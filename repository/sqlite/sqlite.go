package sqlite

import (
	"database/sql"
	"errors"
	"fmt"
	"log"

	"accu_pls/playlist"
	"accu_pls/repository"
	_ "github.com/mattn/go-sqlite3"
)

type Repository struct {
	db *sql.DB
}

const qCreateTracksTable = `
CREATE TABLE IF NOT EXISTS tracks (
	id INTEGER PRIMARY KEY,
	channel TEXT NOT NULL DEFAULT '',
	artist TEXT NOT NULL DEFAULT '',
	album TEXT NOT NULL DEFAULT '',
	title TEXT NOT NULL DEFAULT '',
	duration INTEGER NOT NULL DEFAULT 0,
	year TEXT NOT NULL DEFAULT '',
	primary_link TEXT NOT NULL DEFAULT '' UNIQUE,
	secondary_link TEXT NOT NULL DEFAULT '' UNIQUE
)
`

func NewSqliteRepo() (*Repository, error) {
	db, err := sql.Open("sqlite3", "file:playlist.sqlite?mode=rwc&cache=shared")
	if err != nil {
		return nil, fmt.Errorf("new sqlite db file %w", err)
	}
	if _, err := db.Exec(qCreateTracksTable); err != nil {
		return nil, fmt.Errorf("new sqlite db file create tracks table %w", err)
	}
	return &Repository{
		db: db,
	}, nil
}

const qSaveTrack = `
INSERT INTO tracks(
	channel,
	artist,
	album,
	title,
	duration,
	year,
	primary_link,
	secondary_link
) VALUES($1, $2, $3, $4, $5, $6, $7, $8)
`

func (r *Repository) Save(channel string, tracks []*playlist.Track) error {
	tx, err := r.db.Begin()
	if err != nil {
		return fmt.Errorf("begin tx %w", err)
	}
	existingTracks := 0
	for _, t := range tracks {
		if exists, err := r.Exists(t); err == nil {
			if exists {
				existingTracks++
				continue
			}
		} else {
			return fmt.Errorf("save track %v from channel %s: %w", t, channel, err)
		}
		if _, err := tx.Exec(qSaveTrack, channel, t.Artist, t.Album, t.Title, t.Duration, t.Year, t.PrimaryLink, t.SecondaryLink); err != nil {
			if err := tx.Rollback(); err != nil {
				return fmt.Errorf("rollback save track %v from channel %s: %w", t, channel, err)
			}
			return fmt.Errorf("save track%v from channel %s: %w", t, channel, err)
		}
	}
	if existingTracks == len(tracks) {
		if err := tx.Rollback(); err != nil {
			return fmt.Errorf("rollback save tracks %w", err)
		}
		return repository.ErrNoTracksSaved
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit tracks %w", err)
	}
	return nil
}

const qTrackExists = `
SELECT 1 FROM tracks WHERE primary_link = $1 OR secondary_link = $2
`

func (r *Repository) Exists(track *playlist.Track) (bool, error) {
	var one int
	if err := r.db.QueryRow(qTrackExists, track.PrimaryLink, track.SecondaryLink).Scan(&one); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return false, nil
		}
		return false, fmt.Errorf("check if track exists %w", err)
	}
	return true, nil
}

func (r *Repository) Close() {
	if err := r.db.Close(); err != nil {
		log.Printf("db close %v", err)
	}
}
