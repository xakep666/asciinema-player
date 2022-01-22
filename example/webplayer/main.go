package main

import (
	"embed"
	"encoding/json"
	"flag"
	"fmt"
	"io/fs"
	"log"
	"net"
	"net/http"
	"os"
	"path/filepath"
)

var (
	rootDir    = flag.String("rootdir", ".", "Root directory for casts")
	listenAddr = flag.String("listen-addr", "", "HTTP server listen address")
)

//go:embed web
var webpages embed.FS

func main() {
	flag.Parse()

	fsys, fset, err := getFSWithSet(*rootDir)
	if err != nil {
		log.Fatalln("Open root directory failed:", err)
	}

	flist := make([]string, 0, len(fset))
	for fname := range fset {
		flist = append(flist, fname)
	}

	mux := http.NewServeMux()
	mux.Handle("/play", &TermHandler{FS: fsys, FileSet: fset})
	mux.HandleFunc("/files", func(w http.ResponseWriter, r *http.Request) {
		json.NewEncoder(w).Encode(map[string]interface{}{"files": flist})
	})

	sub, _ := fs.Sub(webpages, "web")
	mux.Handle("/", http.FileServer(http.FS(sub)))

	lis, err := net.Listen("tcp", *listenAddr)
	if err != nil {
		log.Fatalln("Listen failed:", err)
	}

	log.Println("Listening on ", lis.Addr())

	if err = http.Serve(lis, mux); err != http.ErrServerClosed {
		log.Fatalln("Serve failed:", err)
	}
}

func getFSWithSet(rootDir string) (fs.FS, map[string]struct{}, error) {
	abs, err := filepath.Abs(rootDir)
	if err != nil {
		return nil, nil, fmt.Errorf("resolve absoulete path failed: %s", err)
	}

	osfs := os.DirFS(abs)
	fset := make(map[string]struct{})

	err = fs.WalkDir(osfs, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}

		if filepath.Ext(path) == ".cast" && d.Type().IsRegular() {
			fset[path] = struct{}{}
		}

		return nil
	})

	return osfs, fset, err
}
