package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"runtime"
	"time"

	wailsruntime "github.com/wailsapp/wails/v2/pkg/runtime"
)

type versionCheckData struct {
	CodeVersion string `json:"codeVersion"`
	DownloadURL string `json:"downloadUrl"`
}

type versionCheckResponse struct {
	Data *versionCheckData `json:"data"`
}

func platformSuffix() string {
	switch runtime.GOOS {
	case "linux":
		return ".deb"
	case "windows":
		return ".exe"
	case "darwin":
		return ".dmg"
	default:
		return ""
	}
}

func checkForUpdate(currentVersion string) (latest, downloadURL string, hasUpdate bool, err error) {
	const baseURL = "https://api-app-version.vercel.app/v2/api/app"

	client := &http.Client{Timeout: 10 * time.Second}

	url := fmt.Sprintf("%s/version-check?app=runbit-%s", baseURL, currentVersion)
	resp, err := client.Get(url)
	if err != nil {
		return "", "", false, fmt.Errorf("update check request failed: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", "", false, fmt.Errorf("update API returned status %d", resp.StatusCode)
	}

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return "", "", false, fmt.Errorf("reading update response: %w", err)
	}

	var result versionCheckResponse
	if err := json.Unmarshal(body, &result); err != nil {
		return "", "", false, fmt.Errorf("parsing update response: %w", err)
	}

	if result.Data == nil || result.Data.CodeVersion == "" {
		return currentVersion, "", false, nil
	}

	latest = result.Data.CodeVersion
	downloadURL = result.Data.DownloadURL
	hasUpdate = latest != currentVersion
	return latest, downloadURL, hasUpdate, nil
}

func (a *App) StartUpdateChecker() {
	go func() {
		time.Sleep(3 * time.Second)

		latest, downloadURL, hasUpdate, err := checkForUpdate(Version)
		if err != nil {
			log.Printf("[updater] version check error: %v", err)
			return
		}

		if !hasUpdate {
			log.Printf("[updater] already on latest version (%s)", Version)
			return
		}

		log.Printf("[updater] new version available: %s -> %s", Version, latest)

		suffix := platformSuffix()
		body := fmt.Sprintf(
			"Version %s is available (you have %s). Download the %s package.",
			latest, Version, suffix,
		)

		wailsruntime.EventsEmit(a.ctx, "app:update-available", map[string]any{
			"current":     Version,
			"latest":      latest,
			"downloadUrl": downloadURL,
			"platform":    suffix,
		})

		_ = wailsruntime.SendNotification(a.ctx, wailsruntime.NotificationOptions{
			ID:    "runbit-update",
			Title: "RunBit — Update available",
			Body:  body,
		})
	}()
}
