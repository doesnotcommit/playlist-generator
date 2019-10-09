package format

import (
	"io"
)

type Format interface {
	GetReader(channel string, minDuration int) io.Reader
}
