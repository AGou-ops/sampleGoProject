//go:build !linux
// +build !linux

package lumberjack

import (
	"fmt"
	"os"
)

func chown(_ string, _ os.FileInfo) error {
	return nil
}

