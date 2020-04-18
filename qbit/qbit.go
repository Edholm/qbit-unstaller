package qbit

import (
	"edholm.dev/qbit-unstaller/api"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strings"
)

type LoginError struct {
	Cause string
}

func (e *LoginError) Error() string {
	return e.Cause
}

type MissingTorrentError struct {
	Message string
}

func (e *MissingTorrentError) Error() string {
	return e.Message
}

func needLogin(c *http.Client, urlToCall string) bool {
	parsedUrl, err := url.Parse(urlToCall)
	if err != nil {
		log.Panic(err)
	}

	cookies := c.Jar.Cookies(parsedUrl)
	return len(cookies) == 0
}

func login(client *http.Client, settings *api.Settings) (err error) {
	var values = url.Values{}
	values.Set("username", settings.User)
	values.Set("password", settings.Pass)

	var loginUrl = settings.Url + "/api/v2/auth/login"
	req, err := http.NewRequest(http.MethodPost, loginUrl, strings.NewReader(values.Encode()))
	if err != nil {
		return
	}
	req.Header.Add("Referer", settings.Url)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return &LoginError{Cause: "Got non-ok status code on login: " + resp.Status}
	}

	log.Printf("%s was successfully logged in", settings.User)
	return nil
}

func loginIfNeeded(c *http.Client, s *api.Settings, url string) {
	if needLogin(c, url) {
		err := login(c, s)
		if err != nil {
			log.Panic(err)
		}
	}
}

func GetStalledDownloads(c *http.Client, s *api.Settings) (downloads []api.TorrentInfo, err error) {
	stalledUrl := s.Url + "/api/v2/torrents/info?filter=stalled_downloading&limit=10&sort=added_on&reverse=true"
	loginIfNeeded(c, s, stalledUrl)

	resp, err := c.Get(stalledUrl)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	err = json.Unmarshal(body, &downloads)
	return
}

func GetVersion(client *http.Client, settings *api.Settings) (version []byte, err error) {
	versionUrl := settings.Url + "/api/v2/app/version"
	loginIfNeeded(client, settings, versionUrl)

	resp, err := client.Get(versionUrl)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	version, err = ioutil.ReadAll(resp.Body)
	return
}

func GetTrackerInfo(c *http.Client, s *api.Settings, torrent *api.TorrentInfo) (trackerInfo []api.TrackerInfo, err error) {
	var trackerInfoUrl = fmt.Sprintf("%s/api/v2/torrents/trackers?hash=%s", s.Url, torrent.Hash)
	loginIfNeeded(c, s, trackerInfoUrl)

	resp, err := c.Get(trackerInfoUrl)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		err = &MissingTorrentError{
			fmt.Sprintf("Cannot find torrent with hash %s - %s", torrent.Hash, resp.Status),
		}
		return
	}

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}

	err = json.Unmarshal(body, &trackerInfo)
	return
}

func ForceReannounce(c *http.Client, s *api.Settings, hashes *[]string) {
	var announceUrl = s.Url + "/api/v2/torrents/reannounce?hashes=" + combineHashes(hashes)
	resp, err := c.Get(announceUrl)
	if err != nil {
		log.Printf("Failed to reannounce %v: %s", hashes, err)
		return
	}
	defer resp.Body.Close()

	log.Printf("Successfully reannounced %v", hashes)
}

func combineHashes(hashes *[]string) string {
	return strings.Join(*hashes, "|")
}
