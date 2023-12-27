package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"sync"
	"syscall"
	"time"
	"unsafe"

	"github.com/fatih/color"
	"github.com/google/uuid"
)

const (
	MaxGoroutines = 500 // Max number of goroutines to run at once | This should be based on your machine's specs
)

type App struct {
	headers   map[string]string
	client    *http.Client
	wg        sync.WaitGroup
	sem       chan struct{}
	successes int
	fails     int
	errors    int
	startTime time.Time
}

func proxyFunc(req *http.Request) (*url.URL, error) {
	u, err := url.Parse("http://rp.proxyscrape.com:6060") // Replace with your proxy URL
	if err != nil {
		return nil, err
	}
	u.User = url.UserPassword("proxyuser", "proxypass") // Replace with your proxy username and password
	return u, nil
}

func NewApp() *App {
	rand.Seed(time.Now().Unix()) // initialize global pseudo random generator

	return &App{
		headers: map[string]string{
			"authority":          "api.discord.gx.games",
			"accept":             "*/*",
			"accept-language":    "en-US,en;q=0.9",
			"content-type":       "application/json",
			"dnt":                "1",
			"origin":             "https://www.opera.com",
			"referer":            "https://www.opera.com/",
			"sec-ch-ua":          `"Chromium";v="118", "Opera GX";v="104", "Not=A?Brand";v="99"`,
			"sec-ch-ua-mobile":   "?0",
			"sec-ch-ua-platform": `"Windows"`,
			"sec-fetch-dest":     "empty",
			"sec-fetch-mode":     "cors",
			"sec-fetch-site":     "cross-site",
			"user-agent":         "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/118.0.0.0 Safari/537.36 OPR/104.0.0.0",
		},
		client: &http.Client{
			Transport: &http.Transport{
				Proxy: proxyFunc,
			},
		},
		sem:       make(chan struct{}, MaxGoroutines),
		startTime: time.Now(),
	}
}

func (app *App) create() {
	defer app.wg.Done()
	defer func() { <-app.sem }() // release a spot in the semaphore when done

	partnerUserID := uuid.New().String()

	payload, _ := json.Marshal(map[string]string{
		"partnerUserId": partnerUserID,
	})

	req, err := http.NewRequest("POST", "https://api.discord.gx.games/v1/direct-fulfillment", bytes.NewBuffer(payload))
	if err != nil {
		color.Red("Error creating request: %v", err)
		app.errors++
		return
	}

	for key, value := range app.headers {
		req.Header.Set(key, value)
	}

	resp, err := app.client.Do(req)
	if err != nil {
		color.Red("Error making request: %v", err)
		app.errors++
		return
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		color.Red("Error reading response body: %v", err)
		app.errors++
		return
	}

	var jsonResponse map[string]interface{}
	if err := json.Unmarshal(body, &jsonResponse); err != nil {
		color.Red("Error unmarshaling JSON response: %v", err)
		app.errors++
		return
	}

	token, ok := jsonResponse["token"].(string)
	if !ok {
		color.Red("Error getting token from response")
		app.fails++
		return
	}

	url := fmt.Sprintf("https://discord.com/billing/partner-promotions/1180231712274387115/%s", token)
	pink := color.New(color.FgHiMagenta)
	pink.Printf("Promo Created. URL: %s\n", url)
	app.successes++

	appendToFile("promos.txt", url+"\n")
	app.updateTitle()
}

func appendToFile(filename, content string) error {
	file, err := os.OpenFile(filename, os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0644)
	if err != nil {
		return err
	}
	defer file.Close()

	if _, err := file.WriteString(content); err != nil {
		return err
	}
	return nil
}

func (app *App) updateTitle() {
	elapsedTime := time.Since(app.startTime).Round(time.Second)
	unlocksPerMinute := float64(app.successes) / elapsedTime.Minutes()
	successRate := float64(app.successes) / float64(app.successes+app.fails+app.errors) * 100

	title := fmt.Sprintf("Elapsed Time: %v, Successes: %d, Fails: %d, Errors: %d, UPM: %.2f @ Success Rate: %.2f%%",
		elapsedTime, app.successes, app.fails, app.errors, unlocksPerMinute, successRate)

	titlePtr, _ := syscall.UTF16PtrFromString(title)
	syscall.Syscall(syscall.MustLoadDLL("kernel32.dll").MustFindProc("SetConsoleTitleW").Addr(), 1, uintptr(unsafe.Pointer(titlePtr)), 0, 0)
}

func (app *App) Run() {
	for {
		for i := 0; i < MaxGoroutines; i++ {
			app.wg.Add(1)
			app.sem <- struct{}{} // acquire a spot in the semaphore
			go app.create()
		}
		app.wg.Wait()
	}
}

func main() {
	app := NewApp()
	app.Run()
}
