package main

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	runtimepkg "runtime"
	"strings"
	"time"

	"github.com/wailsapp/wails/v2/pkg/runtime"
)

func (a *App) InstallBun() error {
	switch runtimepkg.GOOS {
	case "linux", "darwin":
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
		defer cancel()
		cmd := exec.CommandContext(ctx, "bash", "-c", "curl -fsSL https://bun.com/install | bash")
		var out bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &out

		if err := cmd.Run(); err != nil {
			return fmt.Errorf("error instalando bun: %s", strings.TrimSpace(out.String()))
		}
		return nil

	case "windows":
		ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
		defer cancel()
		cmd := exec.CommandContext(ctx, "powershell", "-c", "irm bun.sh/install.ps1 | iex")
		var out bytes.Buffer
		cmd.Stdout = &out
		cmd.Stderr = &out

		if err := cmd.Run(); err != nil {
			return fmt.Errorf("error instalando bun: %s", strings.TrimSpace(out.String()))
		}
		return nil

	default:
		return fmt.Errorf("plataforma no soportada para instalación automática: %s", runtimepkg.GOOS)
	}
}

func (a *App) InstallBunFlow() error {
	ok, err := a.PromptInstallBun()
	if err != nil {
		return err
	}
	if !ok {
		return nil
	}

	runtime.EventsEmit(a.ctx, "runner:installing-bun")

	if err := a.InstallBun(); err != nil {
		runtime.EventsEmit(a.ctx, "runner:error", map[string]any{
			"result": err.Error(),
		})
		return err
	}

	version, err := a.CheckBunInstalled()
	if err != nil {
		runtime.EventsEmit(a.ctx, "runner:error", map[string]any{
			"result": "Bun parece haberse instalado, pero aún no está disponible",
		})
		return err
	}

	log.Printf("Bun installed successfully: %s", version)

	if err := a.EnsureRunnerProject(); err != nil {
		runtime.EventsEmit(a.ctx, "runner:error", map[string]any{
			"result": fmt.Sprintf("No se pudo preparar runner: %v", err),
		})
		return err
	}

	runtime.EventsEmit(a.ctx, "runner:ready")
	return nil
}

func (a *App) EnsureRunnerProject() error {
	appDir, err := AppDir()
	if err != nil {
		return err
	}

	runnerDir := filepath.Join(appDir, "runner")
	if err := os.MkdirAll(runnerDir, 0o755); err != nil {
		return err
	}

	packageJSON := filepath.Join(runnerDir, "package.json")
	if _, err := os.Stat(packageJSON); errors.Is(err, os.ErrNotExist) {
		content := `{
		"name": "runbit-runner",
		"module": "index.ts",
		"type": "module",
		"dependencies": {}
		}`
		if err := os.WriteFile(packageJSON, []byte(content), 0o644); err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	indexTS := filepath.Join(runnerDir, "index.ts")
	if _, err := os.Stat(indexTS); errors.Is(err, os.ErrNotExist) {
		content := `console.log("Runner ready");`
		if err := os.WriteFile(indexTS, []byte(content), 0o644); err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	tsconfig := filepath.Join(runnerDir, "tsconfig.json")
	if _, err := os.Stat(tsconfig); errors.Is(err, os.ErrNotExist) {
		content := `{
		  "compilerOptions": {
		    "target": "ESNext",
		    "module": "ESNext",
		    "moduleResolution": "bundler",
		    "strict": false,
		    "skipLibCheck": true
		  }
		}`
		if err := os.WriteFile(tsconfig, []byte(content), 0o644); err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	bunfig := filepath.Join(runnerDir, "bunfig.toml")
	if _, err := os.Stat(bunfig); errors.Is(err, os.ErrNotExist) {
		content := `
		[install]
		minimumReleaseAge = 86400
		minimumReleaseAgeExcludes = ["@types/node", "typescript"]
		`
		if err := os.WriteFile(bunfig, []byte(content), 0o644); err != nil {
			return err
		}
	} else if err != nil {
		return err
	}

	bunPath, err := a.FindBunPath()
	if err != nil {
		return err
	}

	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()
	cmd := exec.CommandContext(ctx, bunPath, "install")
	cmd.Dir = runnerDir

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error ejecutando bun install: %s", strings.TrimSpace(out.String()))
	}

	return nil
}
