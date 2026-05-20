package main

import (
	"context"
	"regexp"
)

type App struct {
	ctx context.Context
}

type RunResult struct {
	Stdout   string `json:"stdout"`
	Stderr   string `json:"stderr"`
	ExitCode int    `json:"exitCode"`
}

type PackageInfo struct {
	Name    string `json:"name"`
	Version string `json:"version"`
}

var (
	packageNameRx = regexp.MustCompile(`^[a-zA-Z0-9@/_\.\-]+$`)
	pkgLineRx     = regexp.MustCompile(`[├└]──\s+(@?[^@\s]+(?:/[^@\s]+)?)@(.+)$`)
)
