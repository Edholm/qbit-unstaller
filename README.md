# qBittorrent Unstaller

Sometimes when adding a torrent from an announcement channel really fast, the first tracker announce fails
due to the torrent not being registered with the tracker yet. qBittorrent will then not
do a another announce until next interval (default 30 min) and thus you've "lost the race"

This simple program will find stalled torrents and force a tracker re-announce as a workaround

## How to use
````shell script
go get -v -t -d ./...
go build
./qbit-unstaller -help
````

## Metrics
Metrics are exposed under `:2112/metrics`

### Example prometheus configuration
```
  - job_name: 'qbit-unstaller'
    metrics_path: '/metrics'
    static_configs:
      - targets: ['172.17.0.1:2112']
    relabel_configs:
      - target_label: __address__
        replacement: my-domain.example.com:2112
```

## Affects versions
The "bug" in qBittorrent exists in all versions at least up to `4.2.3`

## See also
[Bug report @ qBittorrent](https://github.com/qbittorrent/qBittorrent/issues/11320)