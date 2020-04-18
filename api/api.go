package api

import "time"

type Settings struct {
	Url      string
	User     string
	Pass     string
	Interval time.Duration
}

type TorrentInfo struct {
	AddedOn           int64   `json:"added_on"`           // Time (Unix Epoch) when the torrent was added to the client
	AmountLeft        int64   `json:"amount_left"`        // Amount of data left to download (bytes)
	AutoTmm           bool    `json:"auto_tmm"`           // Whether this torrent is managed by Automatic Torrent Management
	Availability      float32 `json:"availability"`       // Percentage of file pieces currently available
	Category          string  `json:"category"`           // Category of the torrent
	Completed         int64   `json:"completed"`          // Amount of transfer data completed (bytes)
	CompletionOn      int64   `json:"completion_on"`      // Time (Unix Epoch) when the torrent completed
	DlLimit           int64   `json:"dl_limit"`           // Torrent download speed limit (bytes/s). -1 if ulimited.
	Dlspeed           int64   `json:"dlspeed"`            // Torrent download speed (bytes/s)
	Downloaded        int64   `json:"downloaded"`         // Amount of data downloaded
	DownloadedSession int64   `json:"downloaded_session"` // Amount of data downloaded this session
	Eta               int32   `json:"eta"`                // Torrent ETA (seconds)
	FLPiecePrio       bool    `json:"f_l_piece_prio"`     // True if first last piece are prioritized
	ForceStart        bool    `json:"force_start"`        // True if force start is enabled for this torrent
	Hash              string  `json:"hash"`               // Torrent hash
	LastActivity      int64   `json:"last_activity"`      // Last time (Unix Epoch) when a chunk was downloaded/uploaded
	MagnetUri         string  `json:"magnet_uri"`         // Magnet URI corresponding to this torrent
	MaxRatio          float32 `json:"max_ratio"`          // Maximum share ratio until torrent is stopped from seeding/uploading
	MaxSeedingTime    int32   `json:"max_seeding_time"`   // Maximum seeding time (seconds) until torrent is stopped from seeding
	Name              string  `json:"name"`               // Torrent name
	NumComplete       int32   `json:"num_complete"`       // Number of seeds in the swarm
	NumIncomplete     int32   `json:"num_incomplete"`     // Number of leechers in the swarm
	NumLeechs         int32   `json:"num_leechs"`         // Number of leechers connected to
	NumSeeds          int32   `json:"num_seeds"`          // Number of seeds connected to
	Priority          int32   `json:"priority"`           // Torrent priority.Returns -1 if queuing is disabled or torrent is in seed mode
	Progress          float32 `json:"progress"`           // Torrent progress (percentage/100)
	Ratio             float32 `json:"ratio"`              // Torrent share ratio.Max ratio value: 9999.
	RatioLimit        float32 `json:"ratio_limit"`        // TODO (what is different from max_ratio?)
	SavePath          string  `json:"save_path"`          // Path where this torrent's data is stored
	SeedingTimeLimit  int32   `json:"seeding_time_limit"` // TODO (what is different from max_seeding_time?)
	SeenComplete      int64   `json:"seen_complete"`      // Time (Unix Epoch) when this torrent was last seen complete
	SeqDl             bool    `json:"seq_dl"`             // True if sequential download is enabled
	Size              int64   `json:"size"`               // Total size (bytes) of files selected for download
	State             string  `json:"state"`              // Torrent state.See table here below for the possible values
	SuperSeeding      bool    `json:"super_seeding"`      // True if super seeding is enabled
	Tags              string  `json:"tags"`               // Comma-concatenated tag list of the torrent
	TimeActive        int32   `json:"time_active"`        // Total active time (seconds)
	TotalSize         int64   `json:"total_size"`         // Total size (bytes) of all file in this torrent (including unselected ones)
	Tracker           string  `json:"tracker"`            // The first tracker with working status.(TODO: what is returned if no tracker is working?)
	UpLimit           int32   `json:"up_limit"`           // Torrent upload speed limit (bytes/s).-1 if ulimited.
	Uploaded          int64   `json:"uploaded"`           // Amount of data uploaded
	UploadedSession   int64   `json:"uploaded_session"`   // Amount of data uploaded this session
	Upspeed           int32   `json:"upspeed"`            // Torrent upload speed (bytes/s)
}

type TrackerInfo struct {
	Url           string `json:"url"`            // Tracker url
	Status        int    `json:"status"`         // Tracker status. See the table below for possible values
	NumPeers      int    `json:"num_peers"`      // Number of peers for current torrent, as reported by the tracker
	NumSeeds      int    `json:"num_seeds"`      // Number of seeds for current torrent, asreported by the tracker
	NumLeeches    int    `json:"num_leeches"`    // Number of leeches for current torrent, as reported by the tracker
	NumDownloaded int    `json:"num_downloaded"` // Number of completed downlods for current torrent, as reported by the tracker
	Msg           string `json:"msg"`            // tracker message (there is no way of knowing what this message is - it's up to tracker admins)
}

const (
	TrackerDisabled     = 0 // Tracker is disabled (used for DHT, PeX, and LSD)
	TrackerNotContacted = 1 // Tracker has not been contacted yet
	TrackerWorking      = 2 // Tracker has been contacted and is working
	TrackerUpdating     = 3 // Tracker is updating
	TrackerNotWorking   = 4 // Tracker has been contacted, but it is not working (or doesn't send proper replies)
)
