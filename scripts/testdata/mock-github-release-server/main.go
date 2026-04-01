package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type releaseAsset struct {
	Name               string `json:"name"`
	BrowserDownloadURL string `json:"browser_download_url"`
}

type releaseResponse struct {
	TagName string         `json:"tag_name"`
	Assets  []releaseAsset `json:"assets"`
}

func main() {
	addr := flag.String("addr", "127.0.0.1:0", "listen address")
	repo := flag.String("repo", "", "repository in owner/repo form")
	tag := flag.String("tag", "", "release tag")
	bundleTarball := flag.String("bundle-tarball", "", "path to tar.gz asset")
	checksumsFile := flag.String("checksums-file", "", "path to checksums asset")
	readyFile := flag.String("ready-file", "", "path to write server base URL")
	flag.Parse()

	if *repo == "" || *tag == "" || *bundleTarball == "" || *readyFile == "" {
		log.Fatal("repo, tag, bundle-tarball, and ready-file are required")
	}

	bundleAssetName := filepath.Base(*bundleTarball)
	checksumsAssetName := ""
	if *checksumsFile != "" {
		checksumsAssetName = filepath.Base(*checksumsFile)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/repos/"+*repo+"/releases/latest", func(w http.ResponseWriter, r *http.Request) {
		writeReleaseJSON(w, r, *repo, *tag, bundleAssetName, checksumsAssetName)
	})
	mux.HandleFunc("/repos/"+*repo+"/releases/tags/"+*tag, func(w http.ResponseWriter, r *http.Request) {
		writeReleaseJSON(w, r, *repo, *tag, bundleAssetName, checksumsAssetName)
	})
	mux.HandleFunc("/downloads/"+*repo+"/releases/download/"+*tag+"/"+bundleAssetName, func(w http.ResponseWriter, r *http.Request) {
		http.ServeFile(w, r, *bundleTarball)
	})
	if *checksumsFile != "" {
		mux.HandleFunc("/downloads/"+*repo+"/releases/download/"+*tag+"/"+checksumsAssetName, func(w http.ResponseWriter, r *http.Request) {
			http.ServeFile(w, r, *checksumsFile)
		})
	}

	ln, err := net.Listen("tcp", *addr)
	if err != nil {
		log.Fatalf("listen: %v", err)
	}

	baseURL := "http://" + ln.Addr().String()
	if err := os.WriteFile(*readyFile, []byte(baseURL), 0644); err != nil {
		log.Fatalf("write ready file: %v", err)
	}

	server := &http.Server{Handler: mux}
	if err := server.Serve(ln); err != nil && err != http.ErrServerClosed {
		log.Fatalf("serve: %v", err)
	}
}

func writeReleaseJSON(w http.ResponseWriter, r *http.Request, repo, tag, bundleAssetName, checksumsAssetName string) {
	baseURL := fmt.Sprintf("http://%s", r.Host)
	assets := []releaseAsset{{
		Name:               bundleAssetName,
		BrowserDownloadURL: fmt.Sprintf("%s/downloads/%s/releases/download/%s/%s", baseURL, repo, tag, bundleAssetName),
	}}
	if checksumsAssetName != "" {
		assets = append(assets, releaseAsset{
			Name:               checksumsAssetName,
			BrowserDownloadURL: fmt.Sprintf("%s/downloads/%s/releases/download/%s/%s", baseURL, repo, tag, checksumsAssetName),
		})
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(releaseResponse{TagName: strings.TrimSpace(tag), Assets: assets})
}
