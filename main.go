package main

import (
	"edholm.dev/qbit-service"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	flag "github.com/spf13/pflag"
	"github.com/spf13/viper"
	"log"
	"net/http"
	"time"
)

var (
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

	flag.String("url", "https://luna.lan.elee.cloud/qbittorrent", "The full url including protocol and port to qBittorrent")
	flag.String("username", "admin", "Username to qBittorrent webui")
	flag.String("password", "adminadmin", "Password for the -user")
	flag.Duration("interval", 10000*time.Millisecond, "The duration between checking for stalled torrents")
	flag.Parse()
	err := viper.BindPFlags(flag.CommandLine)
	if err != nil {
		log.Panic(err)
	}

	log.Printf("Using the following settings:\n"+
		"\tUrl: %s\n"+
		"\tUser: %s\n"+
		"\tPass: <redacted>\n"+
		"\tInterval: %s\n",
		viper.GetString("url"), viper.GetString("username"), viper.GetString("interval"))

	printVersion()
	go startUnstallerLoop()

	http.Handle("/metrics", promhttp.Handler())
	err = http.ListenAndServe(":2112", nil)
	if err != nil {
		log.Panic(err)
	}
}

func startUnstallerLoop() {
	interval := viper.GetDuration("interval")
	log.Printf("Starting unstaller loop with interval %s", interval)
	for range time.Tick(interval) {
		loopsMade.Inc()
		reannounceStalledDownloads()
	}
}

func reannounceStalledDownloads() {
	downloads, err := qbit.GetStalledDownloads()
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
		qbit.ForceReannounce(&hashes)
	}
}

func hasNonWorkingTracker(info qbit.TorrentInfo) bool {
	trackerInfo, err := qbit.GetTrackerInfo(&info)
	if err != nil {
		log.Printf("ERROR - %s", err)
		return false
	}

	var nonWorking = false
	for _, t := range trackerInfo {
		trackerStatus.WithLabelValues(string(t.Status)).Inc()
		if t.Status != qbit.TrackerWorking && t.Status != qbit.TrackerDisabled {
			log.Printf("\t%s - %s has a non-working tracker", info.Name, info.Hash)
			nonWorking = true
		}
	}
	return nonWorking
}

func printVersion() {
	version, err := qbit.GetVersion()
	if err != nil {
		log.Fatal(err)
	}

	log.Printf("qBittorrent %v", string(version))
}
