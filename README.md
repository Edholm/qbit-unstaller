# qBittorrent Unstaller

Sometimes when adding a torrent from an announce channel really fast, the first tracker announce fails
due to the torrent not being registered with the tracker yet. qBittorrent will then not
do a another announce until next interval (default 30 min) and thus you've "lost the race"

This simple program will find stalled torrents and force a tracker reannounce as a workaround

## Affects versions
The "bug" in qBittorrent exists in all versions at least up to `4.2.3`

## See also
[Bug report](https://github.com/qbittorrent/qBittorrent/issues/11320)