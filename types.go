package main

type UrlType string

const (
	YoutubeChannel  UrlType = "yt_ch"
	YoutubePlaylist         = "yt_pl"
	RSS                     = "rss"
)

type ChannelType string

const (
	Learning      ChannelType = "learning"
	GameDev                   = "gamedev"
	Blender                   = "blender"
	Music                     = "music"
	Entertainment             = "entertainment"
	RssNews                   = "rss_news"
)

type RssSource struct {
	guid        int64
	urlType     UrlType
	url         string
	channelType ChannelType
	desc        string
	enabled     bool
	lastUpdate  int64
}

type RssItem struct {
	title         string
	link          string
	author        string
	desc          string
	publishedDate string
	thumbnailUrl  string
	videoDuration string
}
