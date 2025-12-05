package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"time"
)

var (
	outputDir = flag.String("out", "./bin", "Output directory")
	panelBuild = flag.Bool("panel", false, "Build panel (Windows AMD64)")
	stubBuild = flag.Bool("stub", false, "Build stub")
	allBuild = flag.Bool("all", false, "Build all binaries")
	vanityPattern = flag.String("vanity", "", "Tor vanity pattern (regex)")
	version = flag.String("version", "0.1.0", "Version string")
	verbose = flag.Bool("v", false, "Verbose output")
)

func main() {
	flag.Parse()

	if !*allBuild && !*panelBuild && !*stubBuild {
		fmt.Println("Usage: builder [flags]")
		flag.PrintDefaults()
		os.Exit(1)
	}

	if *verbose {
		log.SetFlags(log.LstdFlags | log.Lshortfile)
	} else {
		log.SetFlags(0)
	}

	log.Println("Remote Management Builder v0.1.0")

	// Create output directory
	if err := os.MkdirAll(*outputDir, 0755); err != nil {
		log.Fatalf("Error creating output directory: %v", err)
	}

	if *allBuild || *panelBuild {
		buildPanel()
	}

	if *allBuild || *stubBuild {
		buildStub("windows", "amd64", "stub-windows-amd64.exe")
		buildStub("linux", "amd64", "stub-linux-amd64")
		buildStub("linux", "arm", "stub-linux-arm")
	}

	log.Println("Build completed successfully!")
}

func buildPanel() {
	log.Println("Building panel (Windows AMD64)...")

	outputPath := filepath.Join(*outputDir, "panel-windows-amd64.exe")

	cmd := exec.Command("go", "build",
		"-o", outputPath,
		"-ldflags", fmt.Sprintf("-X main.version=%s", *version),
		"./cmd/panel",
	)

	cmd.Env = append(os.Environ(),
		"GOOS=windows",
		"GOARCH=amd64",
		"CGO_ENABLED=1",
	)

	if *verbose {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}

	log.Printf("Executing: GOOS=windows GOARCH=amd64 go build -o %s ./cmd/panel", outputPath)

	if err := cmd.Run(); err != nil {
		log.Fatalf("Build failed: %v", err)
	}

	log.Printf("✓ Panel built: %s", outputPath)
	printBinaryInfo(outputPath)
}

func buildStub(goos, goarch, binaryName string) {
	log.Printf("Building stub (%s %s)...", goos, goarch)

	outputPath := filepath.Join(*outputDir, binaryName)

	ldflags := fmt.Sprintf("-X main.version=%s", *version)
	if *vanityPattern != "" {
		ldflags += fmt.Sprintf(" -X main.vanityPattern=%s", *vanityPattern)
	}

	cmd := exec.Command("go", "build",
		"-o", outputPath,
		"-ldflags", ldflags,
		"./cmd/stub",
	)

	cmd.Env = append(os.Environ(),
		fmt.Sprintf("GOOS=%s", goos),
		fmt.Sprintf("GOARCH=%s", goarch),
		"CGO_ENABLED=0",
	)

	if *verbose {
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}

	log.Printf("Executing: GOOS=%s GOARCH=%s go build -o %s ./cmd/stub", goos, goarch, outputPath)

	if err := cmd.Run(); err != nil {
		log.Fatalf("Build failed: %v", err)
	}

	log.Printf("✓ Stub built: %s", outputPath)
	printBinaryInfo(outputPath)
}

func printBinaryInfo(path string) {
	stat, err := os.Stat(path)
	if err != nil {
		return
	}

	size := stat.Size()
	modTime := stat.ModTime()

	var sizeStr string
	if size > 1024*1024 {
		sizeStr = fmt.Sprintf("%.2f MB", float64(size)/(1024*1024))
	} else if size > 1024 {
		sizeStr = fmt.Sprintf("%.2f KB", float64(size)/1024)
	} else {
		sizeStr = fmt.Sprintf("%d B", size)
	}

	log.Printf("  Size: %s, Modified: %s", sizeStr, modTime.Format(time.RFC3339))
}
