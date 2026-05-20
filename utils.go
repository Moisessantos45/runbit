package main

import (
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strings"
)

func AppDir() (string, error) {
	return os.Getwd()
}

func runnerReady(runnerDir string) bool {
	packageJSON := filepath.Join(runnerDir, "package.json")
	lockfile := filepath.Join(runnerDir, "bun.lock")

	if _, err := os.Stat(packageJSON); err != nil {
		return false
	}

	if _, err := os.Stat(lockfile); err != nil {
		return false
	}

	return true
}

func containsMaliciousPatterns(code string) bool {
	dangerousPatterns := []string{
		"require('child_process')",
		"require(\"child_process\")",
		"exec(",
		"spawn(",
		"fs.rm",
		"fs.unlink",
		"process.exit",
	}
	codeLower := strings.ToLower(code)
	for _, pattern := range dangerousPatterns {
		if strings.Contains(codeLower, pattern) {
			return true
		}
	}
	return false
}

func formatBunError(stderr string) string {
	lines := strings.Split(stderr, "\n")

	var message string
	var location string

	for _, line := range lines {
		line = strings.TrimSpace(line)

		if msg, ok := strings.CutPrefix(line, "error:"); ok {
			message = strings.TrimSpace(msg)
		}

		if strings.Contains(line, "[stdin]:") {
			location = line
		}
	}

	if message == "" {
		return strings.TrimSpace(stderr)
	}

	re := regexp.MustCompile(`\[stdin\]:(\d+):(\d+)`)
	matches := re.FindStringSubmatch(location)

	if len(matches) == 3 {
		return fmt.Sprintf("%s\nLínea %s, columna %s", message, matches[1], matches[2])
	}

	return message
}

func validatePackageName(name string) error {
	name = strings.TrimSpace(name)
	if name == "" {
		return fmt.Errorf("nombre de paquete vacío")
	}
	if !packageNameRx.MatchString(name) {
		return fmt.Errorf("nombre de paquete inválido")
	}
	return nil
}

func ParseBunPmLs(output string) []PackageInfo {
	lines := strings.Split(output, "\n")
	result := make([]PackageInfo, 0)

	for _, line := range lines {
		line = strings.TrimSpace(line)
		m := pkgLineRx.FindStringSubmatch(line)
		if len(m) == 3 {
			result = append(result, PackageInfo{
				Name:    strings.TrimSpace(m[1]),
				Version: strings.TrimSpace(m[2]),
			})
		}
	}

	return result
}
