// Copyright 2024 Colorado School of Mines CSCI 370 FA24 NREL 2 Group

package vwifi

import (
	"fmt"

	log "github.com/sandia-minimega/minimega/v2/pkg/minilog"
)

// io.Writer implementation that logs to minilog
type minilogWriter struct {
	// The log level to log at
	level log.Level

	// The log prefix
	prefix string
}

// write data to the minilog logger
func (writer minilogWriter) Write(data []byte) (n int, err error) {
	switch writer.level {
	case log.DEBUG:
		log.Debug(fmt.Sprintf("%s: %%s", writer.prefix), string(data))
	case log.INFO:
		log.Info(fmt.Sprintf("%s: %%s", writer.prefix), string(data))
	case log.WARN:
		log.Warn(fmt.Sprintf("%s: %%s", writer.prefix), string(data))
	case log.ERROR:
		log.Error(fmt.Sprintf("%s: %%s", writer.prefix), string(data))
	case log.FATAL:
		log.Fatal(fmt.Sprintf("%s: %%s", writer.prefix), string(data))
	}

	return len(data), nil
}
