package main

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

func (a *App) startup(ctx context.Context) {
	a.ctx = ctx
	a.registerEvents()

	go a.bootstrapRunner()

	a.StartUpdateChecker()
}

func (a *App) CheckBunInstalled() (string, error) {
	bunPath, err := a.FindBunPath()
	if err != nil {
		return "", err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, bunPath, "--version")
	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	if err := cmd.Run(); err != nil {
		return "", fmt.Errorf("no se pudo verificar bun: %w", err)
	}

	return strings.TrimSpace(out.String()), nil
}

func (a *App) bootstrapRunner() {
	version, err := a.CheckBunInstalled()
	if err != nil {
		runtime.EventsEmit(a.ctx, "runner:bun-missing")
		return
	}

	log.Printf("Bun detected: %s", version)

	appDir, err := AppDir()
	if err != nil {
		runtime.EventsEmit(a.ctx, "runner:error", map[string]any{
			"result": fmt.Sprintf("No se pudo obtener AppDir: %v", err),
		})
		return
	}

	runnerDir := filepath.Join(appDir, "runner")

	if !runnerReady(runnerDir) {
		if err := a.EnsureRunnerProject(); err != nil {
			runtime.EventsEmit(a.ctx, "runner:error", map[string]any{
				"result": fmt.Sprintf("No se pudo preparar runner: %v", err),
			})
			return
		}
	}

	runtime.EventsEmit(a.ctx, "runner:ready")
}

func (a *App) PromptInstallBun() (bool, error) {
	result, err := runtime.MessageDialog(a.ctx, runtime.MessageDialogOptions{
		Type:          runtime.QuestionDialog,
		Title:         "Bun no encontrado",
		Message:       "Bun no está instalado. ¿Deseas instalarlo ahora?",
		DefaultButton: "No",
	})
	if err != nil {
		return false, err
	}

	return result == "Yes" || result == "Ok" || result == "Continue", nil
}
