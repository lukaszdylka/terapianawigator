package main

import (
	"context"
	"crypto/sha256"
	_ "embed"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strconv"
	"strings"
	"time"
)

//go:embed index.html
var indexHTML []byte

const addr = "127.0.0.1:8765"
const appID = "kalendarz-specjalisty-html-v3"
const appVersion = 3

func edgePath() string {
	candidates := []string{}
	if p := os.Getenv("PROGRAMFILES(X86)"); p != "" {
		candidates = append(candidates, filepath.Join(p, "Microsoft", "Edge", "Application", "msedge.exe"))
	}
	if p := os.Getenv("PROGRAMFILES"); p != "" {
		candidates = append(candidates, filepath.Join(p, "Microsoft", "Edge", "Application", "msedge.exe"))
	}
	if p := os.Getenv("LOCALAPPDATA"); p != "" {
		candidates = append(candidates, filepath.Join(p, "Microsoft", "Edge", "Application", "msedge.exe"))
	}
	for _, p := range candidates {
		if st, err := os.Stat(p); err == nil && !st.IsDir() {
			return p
		}
	}
	return ""
}

func stopOldInstance() {
	client := &http.Client{Timeout: 700 * time.Millisecond}
	resp, err := client.Get("http://" + addr + "/__ks_app")
	if err != nil {
		return
	}
	body, _ := io.ReadAll(io.LimitReader(resp.Body, 128))
	resp.Body.Close()
	if strings.TrimSpace(string(body)) != appID {
		return
	}
	req, _ := http.NewRequest(http.MethodPost, "http://"+addr+"/__quit", nil)
	_, _ = client.Do(req)
	for i := 0; i < 20; i++ {
		c, err := net.DialTimeout("tcp", addr, 100*time.Millisecond)
		if err != nil {
			return
		}
		c.Close()
		time.Sleep(100 * time.Millisecond)
	}
}

type updateReq struct { URL string `json:"url"` }
type updateManifest struct {
	Version int `json:"version"`
	Label string `json:"label"`
	URL string `json:"url"`
	SHA256 string `json:"sha256"`
}

func appDataDir() string {
	base := os.Getenv("LOCALAPPDATA")
	if base == "" { base = os.TempDir() }
	d := filepath.Join(base, "KalendarzSpecjalisty")
	_ = os.MkdirAll(d, 0700)
	return d
}

func fetchManifest(url string) (updateManifest, error) {
	var m updateManifest
	if !strings.HasPrefix(strings.ToLower(url), "https://") {
		return m, fmt.Errorf("adres aktualizacji musi używać HTTPS")
	}
	c := &http.Client{Timeout: 12*time.Second}
	resp, err := c.Get(url)
	if err != nil { return m, err }
	defer resp.Body.Close()
	if resp.StatusCode != 200 { return m, fmt.Errorf("HTTP %d", resp.StatusCode) }
	if err := json.NewDecoder(io.LimitReader(resp.Body, 1<<20)).Decode(&m); err != nil { return m, err }
	if m.Version <= 0 || !strings.HasPrefix(strings.ToLower(m.URL), "https://") || len(m.SHA256) != 64 {
		return m, fmt.Errorf("nieprawidłowy manifest")
	}
	return m,nil
}

func downloadUpdate(m updateManifest) (string,error) {
	c := &http.Client{Timeout: 60*time.Second}
	resp,err:=c.Get(m.URL)
	if err!=nil{return "",err}
	defer resp.Body.Close()
	if resp.StatusCode!=200{return "",fmt.Errorf("HTTP %d",resp.StatusCode)}
	dir:=filepath.Join(appDataDir(),"Updates")
	_ = os.MkdirAll(dir,0700)
	p:=filepath.Join(dir,"KalendarzSpecjalisty.new.exe")
	f,err:=os.Create(p); if err!=nil{return "",err}
	h:=sha256.New()
	_,err=io.Copy(io.MultiWriter(f,h),io.LimitReader(resp.Body,100<<20))
	cerr:=f.Close()
	if err!=nil{return "",err}; if cerr!=nil{return "",cerr}
	got:=hex.EncodeToString(h.Sum(nil))
	if !strings.EqualFold(got,m.SHA256){_ = os.Remove(p);return "",fmt.Errorf("suma SHA256 nie pasuje")}
	return p,nil
}

func scheduleReplacement(newExe string) error {
	cur,err:=os.Executable(); if err!=nil{return err}
	cur,_=filepath.Abs(cur)
	helper:=filepath.Join(appDataDir(),"apply_update.cmd")
	script:=fmt.Sprintf("@echo off\r\nping 127.0.0.1 -n 3 >nul\r\ncopy /Y %s %s >nul\r\nstart \"\" %s\r\ndel /Q %s\r\ndel /Q \"%%~f0\"\r\n",
		strconv.Quote(newExe),strconv.Quote(cur),strconv.Quote(cur),strconv.Quote(newExe))
	if err:=os.WriteFile(helper,[]byte(script),0600);err!=nil{return err}
	cmd:=exec.Command("cmd.exe","/C",helper)
	cmd.SysProcAttr = nil
	return cmd.Start()
}

func main() {
	if runtime.GOOS != "windows" {
		return
	}

	stopOldInstance()

	ln, err := net.Listen("tcp", addr)
	if err != nil {
		return
	}

	mux := http.NewServeMux()
	server := &http.Server{
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" && r.URL.Path != "/index.html" {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Header().Set("Cache-Control", "no-store, no-cache, must-revalidate")
		w.Header().Set("Pragma", "no-cache")
		w.Write(indexHTML)
	})

	mux.HandleFunc("/favicon.ico", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "image/svg+xml")
		w.Header().Set("Cache-Control", "public, max-age=86400")
		fmt.Fprint(w, `<svg xmlns="http://www.w3.org/2000/svg" viewBox="0 0 64 64"><rect width="64" height="64" rx="14" fill="#2443B0"/><rect x="14" y="18" width="36" height="32" rx="5" fill="white"/><path d="M14 27h36" stroke="#2443B0" stroke-width="4"/><circle cx="24" cy="35" r="3" fill="#2443B0"/><circle cx="32" cy="35" r="3" fill="#2443B0"/><circle cx="40" cy="35" r="3" fill="#2443B0"/><circle cx="24" cy="43" r="3" fill="#2443B0"/><circle cx="32" cy="43" r="3" fill="#2443B0"/></svg>`)
	})

	mux.HandleFunc("/__backup", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost { http.Error(w,"method",http.StatusMethodNotAllowed); return }
		body,err:=io.ReadAll(io.LimitReader(r.Body,20<<20))
		if err!=nil || len(body)==0 { http.Error(w,"bad body",400); return }
		dir:=filepath.Join(appDataDir(),"Backups")
		_ = os.MkdirAll(dir,0700)
		name:=filepath.Join(dir,"auto_"+time.Now().Format("2006-01-02")+".ksbackup")
		if err:=os.WriteFile(name,body,0600);err!=nil{http.Error(w,err.Error(),500);return}
		w.Header().Set("Content-Type","application/json"); fmt.Fprint(w,`{"ok":true}`)
	})

	mux.HandleFunc("/__update/check", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost { http.Error(w,"method",http.StatusMethodNotAllowed); return }
		var q updateReq
		if err:=json.NewDecoder(io.LimitReader(r.Body,4096)).Decode(&q);err!=nil{http.Error(w,`{"error":"bad request"}`,400);return}
		m,err:=fetchManifest(q.URL)
		w.Header().Set("Content-Type","application/json")
		if err!=nil{w.WriteHeader(502); _=json.NewEncoder(w).Encode(map[string]any{"error":err.Error()});return}
		_ = json.NewEncoder(w).Encode(map[string]any{"available":m.Version>appVersion,"version":m.Version,"label":m.Label,"current":appVersion})
	})

	mux.HandleFunc("/__update/install", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost { http.Error(w,"method",http.StatusMethodNotAllowed); return }
		var q updateReq
		if err:=json.NewDecoder(io.LimitReader(r.Body,4096)).Decode(&q);err!=nil{http.Error(w,`{"error":"bad request"}`,400);return}
		m,err:=fetchManifest(q.URL)
		if err==nil && m.Version<=appVersion { err=fmt.Errorf("brak nowszej wersji") }
		var p string
		if err==nil { p,err=downloadUpdate(m) }
		if err==nil { err=scheduleReplacement(p) }
		w.Header().Set("Content-Type","application/json")
		if err!=nil{w.WriteHeader(502); _=json.NewEncoder(w).Encode(map[string]any{"error":err.Error()});return}
		_ = json.NewEncoder(w).Encode(map[string]any{"ok":true,"version":m.Version})
	})

	mux.HandleFunc("/__ks_app", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/plain; charset=utf-8")
		fmt.Fprint(w, appID)
	})

	mux.HandleFunc("/__quit", func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			http.Error(w, "method", http.StatusMethodNotAllowed)
			return
		}
		fmt.Fprint(w, "ok")
		go func() {
			time.Sleep(100 * time.Millisecond)
			ctx, cancel := context.WithTimeout(context.Background(), time.Second)
			defer cancel()
			_ = server.Shutdown(ctx)
		}()
	})

	go func() {
		_ = server.Serve(ln)
	}()

	url := "http://" + addr + "/?desktop=1"
	edge := edgePath()
	if edge == "" {
		_ = exec.Command("rundll32", "url.dll,FileProtocolHandler", url).Start()
		select {}
	}

	local := os.Getenv("LOCALAPPDATA")
	if local == "" {
		local = os.TempDir()
	}
	profile := filepath.Join(local, "KalendarzSpecjalisty", "EdgeProfile")
	_ = os.MkdirAll(profile, 0700)

	cmd := exec.Command(edge,
		"--app="+url,
		"--user-data-dir="+profile,
		"--no-first-run",
		"--no-default-browser-check",
		"--disable-features=msEdgeSidebarV2",
		"--disable-session-crashed-bubble",
	)
	_ = cmd.Start()
	_ = cmd.Wait()

	ctx, cancel := context.WithTimeout(context.Background(), 1200*time.Millisecond)
	defer cancel()
	_ = server.Shutdown(ctx)
}
