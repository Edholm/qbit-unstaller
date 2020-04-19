package main

import (
	"edholm.dev/qbit-unstaller/api"
	"edholm.dev/qbit-unstaller/qbit"
	"flag"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"log"
	"net/http"
	"net/http/cookiejar"
	"time"
)

var (
	settings api.Settings
	client   http.Client

	loopsMade = promauto.NewCounter(
		prometheus.CounterOpts{
			Name: "qbit_unstaller_loops_made",
			Help: "The number of unstaller loops made",
		})

	stalledDownloads = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "qbit_unstaller_stalled_downloads",
			Help: "The number of stalled downloads seen",
		},
		[]string{"working_tracker"})

	trackerStatus = promauto.NewCounterVec(
		prometheus.CounterOpts{
			Name: "qbit_unstaller_tracker_status",
			Help: "The status of the trackers",
		},
		[]string{"state"})
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
	go startUnstallerLoop()

	http.Handle("/metrics", promhttp.Handler())
	err := http.ListenAndServe(":2112", nil)
	if err != nil {
		log.Panic(err)
	}
}

func startUnstallerLoop() {
	log.Printf("Starting unstaller loop with interval %s", settings.Interval)
	for range time.Tick(settings.Interval) {
		loopsMade.Inc()
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
			stalledDownloads.WithLabelValues("false").Inc()
			hashes = append(hashes, info.Hash)
		} else {
			stalledDownloads.WithLabelValues("true").Inc()
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

	var nonWorking = false
	for _, t := range trackerInfo {
		trackerStatus.WithLabelValues(string(t.Status)).Inc()
		if t.Status != api.TrackerWorking && t.Status != api.TrackerDisabled {
			log.Printf("\t%s - %s has a non-working tracker", info.Name, info.Hash)
			nonWorking = true
		}
	}
	return nonWorking
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
