//go:build debug
// +build debug

package client

import "log"

func debug(fmt string, args ...interface{}) {
	log.Printf(fmt, args...)
}
