// Data structures related to whole program
package main

import (
	"bytes"
	"net/http"
	"time"
)

const (
	SERVER = "HLS Probe II"
)

// Kinds of streams
const (
	SAMPLE StreamType = iota // internet resources for monitor self checks
	HTTP                     // "plain" HTTP
	HLS                      // Apple HTTP Live Streaming
	HDS                      // Adobe HTTP Dynamic Streaming
)

// Error codes (ordered by errors importance).
// If several errors detected then only one with the heaviest weight reported.
const (
	SUCCESS        ErrType = iota
	DEBUG_LEVEL            // Internal debug messages follow below:
	HLSPARSER              // HLS parser error (debug)
	BADREQUEST             // Request failed (internal client error)
	WARNING_LEVEL          // Warnings follow below:
	SLOW                   // SlowWarning threshold on reading server response
	VERYSLOW               // VerySlowWarning threshold on reading server response
	ERROR_LEVEL            // Errors follow below:
	RTIMEOUT               // Timeout on read
	CTIMEOUT               // Timeout on connect
	BADLENGTH              // ContentLength value not equal real content length
	BODYREAD               // Response body read error
	CRITICAL_LEVEL         // Permanent errors level
	REFUSED                // Connectin refused
	BADSTATUS              // HTTP Status >= 400
	BADURI                 // Incorret URI format
	LISTEMPTY              // HLS specific (by m3u8 lib)
	BADFORMAT              // HLS specific (by m3u8 lib)
	UNKNOWN
)

// Commands for probers.
const (
	STOP Command = iota
	START
	RELOAD // reread config
)

type StreamType uint // Type of checked streams
type ErrType uint
type Command uint

type Stream struct {
	URI   string
	Type  StreamType
	Name  string
	Group string
}

// Stream checking task
type Task struct {
	Stream
	ReadBody bool
	ReplyTo  chan TaskResult
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
	ErrType           ErrType
	HTTPCode          int    // HTTP status code
	HTTPStatus        string // HTTP status string
	ContentLength     int64
	RealContentLength int64
	Headers           http.Header
	Body              bytes.Buffer
	Started           time.Time
	Elapsed           time.Duration
	TotalErrs         uint
}

type StreamStats struct {
	Stream Stream
	Last   TaskResult
}

type ErrHistoryKey struct {
	Curhour string
	ErrType ErrType
	Group   string
	Name    string
	URI     string
}

type ErrTotalHistoryKey struct {
	Curhour string
	Group   string
	Name    string
}