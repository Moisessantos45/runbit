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


func NewApp() *App {
	return &App{}
}

func (a *App) registerEvents() {
	runtime.EventsOn(a.ctx, "runner:run", func(data ...any) {
		a.parseData(data...)
	})
}

func (a *App) Greet(name string) string {
	return fmt.Sprintf("Hello %s, It's show time!", name)
}

func (a *App) FindBunPath() (string, error) {
	if bunPath, err := exec.LookPath("bun"); err == nil {
		return bunPath, nil
	}

	homeDir, err := os.UserHomeDir()
	if err == nil {
		bunPath := filepath.Join(homeDir, ".bun", "bin", "bun")
		if runtimepkg.GOOS == "windows" {
			bunPath = filepath.Join(homeDir, ".bun", "bin", "bun.exe")
		}

		if _, err := os.Stat(bunPath); err == nil {
			return bunPath, nil
		}
	}

	return "", fmt.Errorf("bun no está instalado")
}

func (a *App) parseData(data ...any) {
	if len(data) == 0 {
		runtime.EventsEmit(a.ctx, "runner:error", map[string]any{
			"result": "No data received",
		})
		return
	}

	payload, ok := data[0].(map[string]any)
	if !ok {
		runtime.EventsEmit(a.ctx, "runner:error", map[string]any{
			"result": "Invalid data format",
		})
		return
	}

	code, ok := payload["code"].(string)
	if !ok {
		runtime.EventsEmit(a.ctx, "runner:error", map[string]any{
			"result": "code debe ser string",
		})
		return
	}

	lang, ok := payload["lang"].(string)
	if !ok {
		runtime.EventsEmit(a.ctx, "runner:error", map[string]any{
			"result": "lang debe ser string",
		})
		return
	}

	v, exists := payload["tab"]
	if !exists {
		runtime.EventsEmit(a.ctx, "runner:error", map[string]any{
			"result": "tab es obligatorio",
		})
		return
	}

	tabFloat, ok := v.(float64)
	if !ok {
		runtime.EventsEmit(a.ctx, "runner:error", map[string]any{
			"tab":    0,
			"result": "tab debe ser un número",
		})
		return
	}

	if tabFloat != float64(int(tabFloat)) {
		runtime.EventsEmit(a.ctx, "runner:error", map[string]any{
			"tab":    0,
			"result": "tab debe ser un entero",
		})
		return
	}

	tab := int(tabFloat)

	log.Printf("Received code to run: lang=%s, tab=%d", lang, tab)

	switch lang {
	case "js", "ts":
		runtime.EventsEmit(a.ctx, "runner:start", map[string]any{
			"tab": tab,
		})

		go func(tabToRun int, codeToRun string) {
			res := a.RunJS(codeToRun)

			if res.Stdout != "" {
				runtime.EventsEmit(a.ctx, "runner:stdout", map[string]any{
					"tab":    tabToRun,
					"result": res.Stdout,
				})
			}

			if res.Stderr != "" {
				runtime.EventsEmit(a.ctx, "runner:stderr", map[string]any{
					"tab":    tabToRun,
					"result": formatBunError(res.Stderr),
				})
			}

			runtime.EventsEmit(a.ctx, "runner:done", map[string]any{
				"tab":      tabToRun,
				"exitCode": res.ExitCode,
			})
		}(tab, code)

	default:
		runtime.EventsEmit(a.ctx, "runner:error", map[string]any{
			"tab":    tab,
			"result": "Lenguaje no soportado",
		})
	}
}

func (a *App) RunJS(code string) RunResult {
	if containsMaliciousPatterns(code) {
		return RunResult{
			Stderr:   "Ejecución bloqueada: Se detectó código potencialmente peligroso o no permitido por seguridad.",
			ExitCode: -1,
		}
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	appDir, err := AppDir()
	if err != nil {
		return RunResult{
			Stderr:   fmt.Sprintf("Error obteniendo directorio de la aplicación: %v", err),
			ExitCode: -1,
		}
	}

	runnerDir := filepath.Join(appDir, "runner")

	bunPath, err := a.FindBunPath()
	if err != nil {
		return RunResult{
			Stderr:   "Bun no está instalado o no está disponible",
			ExitCode: -1,
		}
	}

	cmd := exec.CommandContext(ctx, bunPath, "run", "-")
	cmd.Dir = runnerDir
	cmd.Stdin = bytes.NewBufferString(code)

	var stdout bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err = cmd.Run()

	result := RunResult{
		Stdout: stdout.String(),
		Stderr: stderr.String(),
	}

	if err != nil {
		if errors.Is(err, context.DeadlineExceeded) || errors.Is(ctx.Err(), context.DeadlineExceeded) {
			result.Stderr = "Execution timed out"
			result.ExitCode = -1
			return result
		}

		if exitErr, ok := err.(*exec.ExitError); ok {
			result.ExitCode = exitErr.ExitCode()
		} else {
			result.ExitCode = -1
		}
		return result
	}

	result.ExitCode = 0
	return result
}

func (a *App) AddPackage(name string) error {
	if err := validatePackageName(name); err != nil {
		return err
	}

	appDir, err := AppDir()
	if err != nil {
		return err
	}

	bunPath, err := a.FindBunPath()
	if err != nil {
		return err
	}

	runnerDir := filepath.Join(appDir, "runner")
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()
	cmd := exec.CommandContext(ctx, bunPath, "add", name)
	cmd.Dir = runnerDir

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error agregando paquete: %s", strings.TrimSpace(out.String()))
	}

	return nil
}

func (a *App) RemovePackage(name string) error {
	if err := validatePackageName(name); err != nil {
		return err
	}

	appDir, err := AppDir()
	if err != nil {
		return err
	}

	bunPath, err := a.FindBunPath()
	if err != nil {
		return err
	}

	runnerDir := filepath.Join(appDir, "runner")
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Minute)
	defer cancel()
	cmd := exec.CommandContext(ctx, bunPath, "remove", name)
	cmd.Dir = runnerDir

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("error removiendo paquete: %s", strings.TrimSpace(out.String()))
	}

	return nil
}

func (a *App) ListPackages(all bool) (string, error) {
	appDir, err := AppDir()
	if err != nil {
		return "", err
	}

	bunPath, err := a.FindBunPath()
	if err != nil {
		return "", err
	}

	runnerDir := filepath.Join(appDir, "runner")

	args := []string{"pm", "ls"}
	if all {
		args = append(args, "--all")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, bunPath, args...)
	cmd.Dir = runnerDir

	var out bytes.Buffer
	var stderr bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &stderr

	if err := cmd.Run(); err != nil {
		stderrStr := stderr.String()
		if strings.Contains(stderrStr, "Lockfile not found") || strings.Contains(stderrStr, "No package.json was found") {
			return "", nil
		}
		return "", fmt.Errorf("%v: %s", err, stderrStr)
	}

	return out.String(), nil
}

func (a *App) ListInstalledPackages() ([]PackageInfo, error) {
	output, err := a.ListPackages(false)
	if err != nil {
		return nil, err
	}

	return ParseBunPmLs(output), nil
}
