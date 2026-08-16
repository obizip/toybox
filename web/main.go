package main

import (
	"bytes"
	"embed"
	"flag"
	"fmt"
	"html/template"
	"io"
	"io/fs"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

//go:embed templates/* static/*
var contentFS embed.FS

type GameMeta struct {
	ID          string
	Title       string
	Description string
}

func main() {
	outFlag := flag.String("out", "", "Output directory path (default: <root>/dist)")
	flag.Parse()

	// Determine repository root directory
	execDir, err := os.Getwd()
	if err != nil {
		log.Fatalf("Failed to get working directory: %v", err)
	}

	rootDir := execDir
	if filepath.Base(execDir) == "web" {
		rootDir = filepath.Dir(execDir)
	}

	distDir := *outFlag
	if distDir == "" {
		distDir = filepath.Join(rootDir, "dist")
	}

	log.Printf("Repository root: %s", rootDir)
	log.Printf("Output directory: %s", distDir)

	// Clean and recreate dist directory
	if err := os.RemoveAll(distDir); err != nil {
		log.Fatalf("Failed to clean dist directory: %v", err)
	}
	if err := os.MkdirAll(distDir, 0755); err != nil {
		log.Fatalf("Failed to create dist directory: %v", err)
	}

	// Copy embedded static assets to dist
	if err := copyEmbeddedDir(contentFS, "static", distDir); err != nil {
		log.Fatalf("Failed to copy static assets: %v", err)
	}

	// Copy wasm_exec.js from Go standard library
	if err := copyWasmExec(distDir); err != nil {
		log.Fatalf("Failed to copy wasm_exec.js: %v", err)
	}

	// Parse templates from embedded FS
	indexTmpl, err := template.ParseFS(contentFS, "templates/index.html")
	if err != nil {
		log.Fatalf("Failed to parse index.html template: %v", err)
	}

	playerTmpl, err := template.ParseFS(contentFS, "templates/player.html")
	if err != nil {
		log.Fatalf("Failed to parse player.html template: %v", err)
	}

	// Discover game directories in root
	entries, err := os.ReadDir(rootDir)
	if err != nil {
		log.Fatalf("Failed to read root directory: %v", err)
	}

	var games []GameMeta
	for _, entry := range entries {
		if !entry.IsDir() {
			continue
		}
		name := entry.Name()
		if strings.HasPrefix(name, ".") || name == "web" || name == "dist" {
			continue
		}

		gameDirPath := filepath.Join(rootDir, name)
		mainGoPath := filepath.Join(gameDirPath, "main.go")
		if _, err := os.Stat(mainGoPath); os.IsNotExist(err) {
			continue
		}

		log.Printf("Discovered game: %s", name)

		meta := parseGameMeta(gameDirPath, name)

		// Create game directory in dist
		gameDistDir := filepath.Join(distDir, name)
		if err := os.MkdirAll(gameDistDir, 0755); err != nil {
			log.Fatalf("Failed to create directory %s: %v", gameDistDir, err)
		}

		// Build Wasm binary
		wasmOut := filepath.Join(gameDistDir, "game.wasm")
		log.Printf("Building Wasm for %s...", name)
		cmd := exec.Command("go", "build", "-o", wasmOut, ".")
		cmd.Dir = gameDirPath
		cmd.Env = append(os.Environ(), "GOOS=js", "GOARCH=wasm")
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		if err := cmd.Run(); err != nil {
			log.Fatalf("Failed to build Wasm for %s: %v", name, err)
		}

		// Render game's index.html
		playerHTMLPath := filepath.Join(gameDistDir, "index.html")
		var playerBuf bytes.Buffer
		if err := playerTmpl.Execute(&playerBuf, meta); err != nil {
			log.Fatalf("Failed to render player HTML for %s: %v", name, err)
		}
		if err := os.WriteFile(playerHTMLPath, playerBuf.Bytes(), 0644); err != nil {
			log.Fatalf("Failed to write %s: %v", playerHTMLPath, err)
		}

		games = append(games, meta)
	}

	// Render root index.html
	rootIndexPath := filepath.Join(distDir, "index.html")
	var indexBuf bytes.Buffer
	if err := indexTmpl.Execute(&indexBuf, games); err != nil {
		log.Fatalf("Failed to render root index HTML: %v", err)
	}
	if err := os.WriteFile(rootIndexPath, indexBuf.Bytes(), 0644); err != nil {
		log.Fatalf("Failed to write root index.html: %v", err)
	}

	log.Printf("Build succeeded! Generated static site for %d games in %s", len(games), distDir)
}

func parseGameMeta(dirPath, dirName string) GameMeta {
	meta := GameMeta{
		ID:          dirName,
		Title:       dirName,
		Description: "",
	}

	readmePath := filepath.Join(dirPath, "README.md")
	content, err := os.ReadFile(readmePath)
	if err != nil {
		return meta
	}

	lines := strings.Split(string(content), "\n")
	var descLines []string
	for _, line := range lines {
		trimmed := strings.TrimSpace(line)
		if trimmed == "" {
			continue
		}
		if strings.HasPrefix(trimmed, "# ") && meta.Title == dirName {
			meta.Title = strings.TrimSpace(strings.TrimPrefix(trimmed, "# "))
		} else if !strings.HasPrefix(trimmed, "#") {
			descLines = append(descLines, trimmed)
		}
	}

	if len(descLines) > 0 {
		meta.Description = strings.Join(descLines, " ")
	}

	return meta
}

func copyWasmExec(destDir string) error {
	gorootCmd := exec.Command("go", "env", "GOROOT")
	out, err := gorootCmd.Output()
	if err != nil {
		return fmt.Errorf("could not find GOROOT: %w", err)
	}
	goroot := strings.TrimSpace(string(out))

	candidates := []string{
		filepath.Join(goroot, "lib", "wasm", "wasm_exec.js"),
		filepath.Join(goroot, "misc", "wasm", "wasm_exec.js"),
	}

	for _, src := range candidates {
		if _, err := os.Stat(src); err == nil {
			dest := filepath.Join(destDir, "wasm_exec.js")
			return copyFile(src, dest)
		}
	}

	return fmt.Errorf("wasm_exec.js not found in GOROOT: %s", goroot)
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	_, err = io.Copy(out, in)
	return err
}

func copyEmbeddedDir(embeddedFS embed.FS, srcDir, dstDir string) error {
	return fs.WalkDir(embeddedFS, srcDir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		relPath, err := filepath.Rel(srcDir, path)
		if err != nil {
			return err
		}

		targetPath := filepath.Join(dstDir, relPath)
		if d.IsDir() {
			return os.MkdirAll(targetPath, 0755)
		}

		data, err := embeddedFS.ReadFile(path)
		if err != nil {
			return err
		}

		return os.WriteFile(targetPath, data, 0644)
	})
}
