//go:build debug
// +build debug

// Package client contains struct for client third parties
package client

import "log"

func debug(fmt string, args ...interface{}) {
	log.Printf(fmt, args...)
}
