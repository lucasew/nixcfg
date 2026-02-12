package media

import (
	"context"
)

type PlaybackStatus string

const (
	StatusPlaying PlaybackStatus = "Playing"
	StatusPaused  PlaybackStatus = "Paused"
	StatusStopped PlaybackStatus = "Stopped"
)

type Metadata struct {
	Title    string
	Artist   string
	ArtUrl   string
	Length   int64 // in microseconds
	Position int64 // in microseconds
	Status   PlaybackStatus
	Player   string // player name/bus name
}

type Driver interface {
	Next(ctx context.Context) error
	Previous(ctx context.Context) error
	PlayPause(ctx context.Context) error
	Stop(ctx context.Context) error
	GetMetadata(ctx context.Context) (*Metadata, error)
}
