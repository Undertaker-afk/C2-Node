package main

import (
	"context"
	"flag"
	"log"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/container"
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
