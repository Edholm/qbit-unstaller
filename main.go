package main

import (
	"edholm.dev/qbit-unstaller/api"
	"edholm.dev/qbit-unstaller/qbit"
	"flag"
	"log"
	"net/http"
	"net/http/cookiejar"
	"time"
)

var (
	settings api.Settings
	client   http.Client
)

func main() {
	log.SetFlags(log.LstdFlags | log.Lmicroseconds)

	flag.StringVar(&settings.Url, "url", "https://luna.lan.elee.cloud/qbittorrent", "The full url including protocol and port to qBittorrent")
	flag.StringVar(&settings.User, "user", "admin", "Username to qBittorrent webui")
	flag.StringVar(&settings.Pass, "password", "adminadmin", "Password for the -user")
	flag.DurationVar(&settings.Interval, "interval", 10000*time.Millisecond, "The duration between checking for stalled torrents")
	flag.Parse()

	log.Printf("Using the following settings:\n"+
		"\tUrl: %s\n"+
		"\tUser: %s\n"+
		"\tPass: <redacted>\n"+
		"\tInterval: %s\n",
		settings.Url, settings.User, settings.Interval)

	client = setupClient()

	printVersion()
	startUnstallerLoop()
}

func startUnstallerLoop() {
	log.Printf("Starting unstaller loop with interval %s", settings.Interval)
	for range time.Tick(settings.Interval) {
		reannounceStalledDownloads()
	}
}

func reannounceStalledDownloads() {
	downloads, err := qbit.GetStalledDownloads(&client, &settings)
	if err != nil {
		log.Panic(err)
	}

	var hashes []string
	for _, info := range downloads {
		if hasNonWorkingTracker(info) {
			hashes = append(hashes, info.Hash)
		}
	}
	if len(hashes) > 0 {
		qbit.ForceReannounce(&client, &settings, &hashes)
	}
}

func hasNonWorkingTracker(info api.TorrentInfo) bool {
	trackerInfo, err := qbit.GetTrackerInfo(&client, &settings, &info)
	if err != nil {
		log.Printf("ERROR - %s", err)
		return false
	}
	for _, t := range trackerInfo {
		if t.Status != api.TrackerWorking && t.Status != api.TrackerDisabled {
			log.Printf("\t%s - %s has a non-working tracker", info.Name, info.Hash)
			return true
		}
	}
	return false
}

func printVersion() {
	version, err := qbit.GetVersion(&client, &settings)
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("qBittorrent %v", string(version))
}

func setupClient() http.Client {
	jar, err := cookiejar.New(nil)
	if err != nil {
		log.Panic(err)
	}

	var client = http.Client{
		Timeout: 1 * time.Second,
		Jar:     jar,
	}
	return client
}
