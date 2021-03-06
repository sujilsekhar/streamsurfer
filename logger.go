// Everything logging
package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"time"
)

const TimeFormat = "2006-01-02 15:04:05"

type LogMessage struct {
	Severity Severity
	Stream
	Result
}

func LogKeeper(verbose bool) {
	var skip error
	var logw *bufio.Writer

	defer func() {
		if r := recover(); r != nil {
			fmt.Println("Logger trace:", r)
		}
	}()

	logq = make(chan LogMessage, 1024)
	logf, skip := os.Create(cfg.ErrorLog)
	if skip == nil {
		logw = bufio.NewWriter(logf)
		fmt.Printf("Error log: %s\n", cfg.ErrorLog)
	} else {
		println("Can't create file for error log. Error logging to file skiped.")
	}

	for {
		timeout := make(chan bool, 1)
		go func() {
			time.Sleep(2 * time.Second)
			timeout <- true
		}()

		select {
		case msg := <-logq:
			if skip == nil {
				logw.WriteString(msg.Started.Format(TimeFormat))
				logw.WriteRune(' ')
				switch msg.Severity {
				case WARNING:
					logw.WriteString("warning")
				case ERROR:
					logw.WriteString("error")
				}
				logw.WriteString(": ")
				logw.WriteString(StreamErr2String(msg.Result.ErrType))
				logw.WriteRune(' ')
				logw.WriteString(strconv.Itoa(msg.HTTPCode))
				logw.WriteRune(' ')
				logw.WriteString(strconv.FormatInt(msg.ContentLength, 10))
				logw.WriteRune(' ')
				logw.WriteString(msg.Elapsed.String())
				logw.WriteRune(' ')
				logw.WriteString(msg.Group)
				logw.WriteString(": ")
				logw.WriteString(msg.Name)
				logw.WriteRune('\n')
			}
		case <-timeout:
			if skip == nil {
				_ = logw.Flush()
			}
		}
	}
}

func Log(severity Severity, stream Stream, taskres Result) {
	logq <- LogMessage{Severity: severity, Stream: stream, Result: taskres}
}
