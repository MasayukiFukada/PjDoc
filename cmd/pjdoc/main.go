package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"runtime"
	"syscall"
	"time"

	"github.com/MasayukiFukada/PjDoc/internal/features/watch_changes"
	sharedweb "github.com/MasayukiFukada/PjDoc/internal/shared/web"
)

func main() {
	port := flag.Int("port", 8080, "Port to listen on")
	dir := flag.String("dir", ".", "Root directory to serve Markdown files from")
	openBrowser := flag.Bool("open", true, "Automatically open browser")
	plantumlServer := flag.String("plantuml-server", "https://kroki.io", "PlantUML / Kroki server URL endpoint for diagram rendering")
	flag.Parse()

	absDir, err := filepath.Abs(*dir)
	if err != nil {
		log.Fatalf("Failed to resolve absolute path: %v", err)
	}

	broadcaster := watch_changes.NewBroadcaster()

	// 定期的にファイルの更新時刻を走査して変更検知 (Hot Reload) を行います
	go watchFileChanges(absDir, broadcaster)

	cfg := sharedweb.ServerConfig{
		RootDir:        absDir,
		Port:           *port,
		PlantUMLServer: *plantumlServer,
		Broadcaster:    broadcaster,
	}

	serverMux, err := sharedweb.NewServer(cfg, sharedweb.EmbeddedAssets)
	if err != nil {
		log.Fatalf("Failed to initialize server: %v", err)
	}

	addr := fmt.Sprintf(":%d", *port)
	server := &http.Server{
		Addr:    addr,
		Handler: serverMux,
	}

	go func() {
		url := fmt.Sprintf("http://localhost:%d", *port)
		fmt.Printf("🚀 PjDoc server starting on %s\n", url)
		fmt.Printf("📂 Serving Markdown documents from: %s\n", absDir)

		if *openBrowser {
			openURL(url)
		}

		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	// シグナルハンドリングによるグレースフルシャットダウン処理
	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt, syscall.SIGTERM)
	<-quit

	fmt.Println("\nShutting down PjDoc server gracefully...")
}

func watchFileChanges(rootDir string, broadcaster *watch_changes.Broadcaster) {
	var lastMaxModTime time.Time

	for {
		time.Sleep(1 * time.Second)
		var currentMax time.Time

		_ = filepath.Walk(rootDir, func(path string, info os.FileInfo, err error) error {
			if err != nil || info.IsDir() {
				return nil
			}
			if filepath.Ext(path) == ".md" {
				if info.ModTime().After(currentMax) {
					currentMax = info.ModTime()
				}
			}
			return nil
		})

		if !lastMaxModTime.IsZero() && currentMax.After(lastMaxModTime) {
			broadcaster.NotifyReload()
		}
		lastMaxModTime = currentMax
	}
}

func openURL(url string) {
	var err error
	switch runtime.GOOS {
	case "linux":
		err = exec.Command("xdg-open", url).Start()
	case "windows":
		err = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
	case "darwin":
		err = exec.Command("open", url).Start()
	}
	if err != nil {
		fmt.Printf("Failed to automatically open browser: %v\n", err)
	}
}
