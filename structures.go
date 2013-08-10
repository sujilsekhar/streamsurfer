// Data structures related to whole program
package main

import (
	"io"
	"net/http"
	"time"
)

const (
	SERVER = "HLS Probe II"
)

// Kinds of streams
const (
	SAMPLE StreamType = iota // internet resources for monitor self checks
	HTTP                     //
	HLS
)

// Error codes (ordered by errors importance).
// If several errors detected then only one with the heaviest weight reported.
const (
	SUCCESS   ErrType = iota
	SLOW              // SlowWarning threshold on reading server response
	VERYSLOW          // VerySlowWarning threshold on reading server response
	BADSTATUS         // HTTP Status >= 400
	BADURI            // Incorret URI format
	LISTEMPTY         // HLS specific (by m3u8 lib)
	BADFORMAT         // HLS specific (by m3u8 lib)
	RTIMEOUT          // Timeout on read
	CTIMEOUT          // Timeout on connect
	HLSPARSER         // HLS parser error (debug)
	UNKNOWN
)

type StreamType uint // Type of checked streams
type ErrType uint

type Stream struct {
	URI   string
	Type  StreamType
	Name  string
	Group string
}

// Stream checking task
type Task struct {
	Stream
	ReplyTo chan TaskResult
}

// Stream group
type GroupTask struct {
	Type    StreamType
	Name    string
	Tasks   *Task
	ReplyTo chan TaskResult
}

// Stream checking result
type TaskResult struct {
	ErrType       ErrType
	HTTPCode      int    // HTTP status code
	HTTPStatus    string // HTTP status string
	ContentLength int64
	Headers       http.Header
	Body          io.ReadCloser
	Started       time.Time
	Elapsed       time.Duration
}

type StreamStats struct {
	Stream   Stream
	Last     TaskResult
	ErrCount map[ErrType]map[time.Time]uint // error log
}