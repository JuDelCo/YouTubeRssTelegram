package main

import (
	"fmt"
	"gorsstelegram/lib"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/mmcdole/gofeed"
)

func processRssSources() {
	for i := range rssSources {
		processRssSource(&rssSources[i])

		if rssSourcesProcessCounter >= rssSourcesProcessMaxCount {
			lib.LogDebug("Exiting loop because max RSS sources processing count reached (" + strconv.FormatInt(int64(rssSourcesProcessCounter), 10) + ")")
			break
		}
	}

	lib.LogDebug("Total RSS feeds processed: " + strconv.FormatInt(int64(rssSourcesProcessCounter), 10))
	lib.LogDebug("Total Telegram RSS item messages sent: " + strconv.FormatInt(int64(telegramRssItemSendCounter), 10))
}

func processRssSource(rssSource *RssSource) {
	if rssSource.enabled == false {
		rssSourceDesc := rssSource.desc

		if rssSourceDesc == "" {
			rssSourceDesc = "[Not Found]"
		}

		lib.LogDebug("Skipping RSS source (disabled): " + strconv.FormatInt(rssSource.guid, 10) + " - " + rssSourceDesc)
		return
	} else if rssSource.enabled == true {
		waitSeconds := int64(60 * waitMinutesToRefreshAgainRssSource)
		timestampDiff := time.Now().UTC().Unix() - rssSource.lastUpdate

		//lib.LogDebug("Compare now <-> lastUpdate: " + time.Now().UTC().Format(time.RFC3339) + " <-> " + time.Unix(rssSource.lastUpdate, 0).UTC().Format(time.RFC3339) + " -> " + strconv.FormatInt(timestampDiff, 10))

		if timestampDiff < waitSeconds {
			min := timestampDiff / 60

			if min < 0 {
				min = 0
			}

			lib.LogDebug("Skipping RSS source (last update was " + strconv.FormatInt(min, 10) + " min ago): " + strconv.FormatInt(rssSource.guid, 10) + " - " + rssSource.desc)
			return
		}
	}

	lib.LogDebug("Processing RSS source # " + strconv.FormatInt(int64((rssSourcesProcessCounter+1)), 10) + " (" + strconv.FormatInt(rssSource.guid, 10) + " - " + rssSource.desc + ") (" + string(rssSource.urlType) + " - " + string(rssSource.channelType) + ")")

	utcNow := time.Now().UTC()
	feed, resp, err := lib.ParseURL(rssSource.url)

	if err != nil {
		lib.LogError("URL: " + rssSource.url + " DESC: " + rssSource.desc)
		lib.LogError(err.Error())
		return
	}

	if resp.StatusCode == 304 {
		lib.LogDebug("    " + "Got 304 status (not modified), setting last update to: " + utcNow.Format(time.RFC3339))
		rssSource.lastUpdate = utcNow.Unix()
		return
	} else if resp.StatusCode < 200 || resp.StatusCode >= 300 {
		lib.LogError("Cant get RSS, HTTP Status: " + strconv.FormatInt(int64(resp.StatusCode), 10) + " URL: " + rssSource.url + " DESC: " + rssSource.desc)
		return
	}

	rssItems := getRssItems(feed, rssSource)
	lib.LogDebug("    " + "Found " + strconv.FormatInt(int64(len(feed.Items)), 10) + " raw RSS items, filtered out to " + strconv.FormatInt(int64(len(rssItems)), 10) + " RSS items")

	for _, rssItem := range rssItems {
		err := notifyEvents(rssSource, rssItem)

		if err != nil {
			rssSourcesProcessCounter += 1
			return
		}
	}

	expires := utcNow

	if _, exists := resp.Header["Expires"]; exists {
		expires, err = http.ParseTime(resp.Header.Get("Expires"))

		if err == nil {
			expires = expires.UTC()

			lib.LogDebug("    " + "Found RSS expires HTTP header: " + expires.Format(time.RFC3339))

			if utcNow.Unix() > expires.Unix() {
				expires = utcNow
			}
		}
	}

	lib.LogDebug("    " + "Update RSS source last update: " + expires.Format(time.RFC3339) + " -> " + time.Unix(rssSource.lastUpdate, 0).UTC().Format(time.RFC3339))

	rssSource.lastUpdate = expires.Unix()

	rssSourcesProcessCounter += 1
}

func getRssItems(feed *gofeed.Feed, rssSource *RssSource) []RssItem {
	rssItems := make([]RssItem, 0)

	utcNowUnix := time.Now().UTC().Unix()
	aDay := int64(60 * 60 * 24)
	maxDaysToIgnoreItem := aDay * int64(ageInDaysToIgnoreRssItem)

	for _, feedItem := range feed.Items {
		if feedItem.UpdatedParsed != nil && utcNowUnix-feedItem.UpdatedParsed.UTC().Unix() > maxDaysToIgnoreItem {
			continue
		}

		if feedItem.PublishedParsed != nil && utcNowUnix-feedItem.PublishedParsed.UTC().Unix() > maxDaysToIgnoreItem {
			continue
		}

		i := sort.SearchStrings(rssUrlLog, strings.TrimSpace(feedItem.Link))

		if i < len(rssUrlLog) && rssUrlLog[i] == strings.TrimSpace(feedItem.Link) {
			continue
		}

		rssItem := RssItem{}
		rssItem.title = strings.TrimSpace(feedItem.Title)
		rssItem.link = strings.TrimSpace(feedItem.Link)
		rssItem.author = strings.TrimSpace(feedItem.Authors[0].Name)
		rssItem.desc = strings.TrimSpace(feedItem.Description)

		if feedItem.PublishedParsed != nil {
			utcPublishedParsed := feedItem.PublishedParsed.UTC().Local()
			rssItem.publishedDate = fmt.Sprintf("%d-%02d-%02d %02d:%02d:%02d",
				utcPublishedParsed.Year(), utcPublishedParsed.Month(), utcPublishedParsed.Day(),
				utcPublishedParsed.Hour(), utcPublishedParsed.Minute(), utcPublishedParsed.Second())
		}

		if rssSource.urlType == YoutubeChannel || rssSource.urlType == YoutubePlaylist {
			rssItem.thumbnailUrl = strings.TrimSpace(feedItem.Extensions["media"]["group"][0].Children["thumbnail"][0].Attrs["url"])
			rssItem.videoDuration = "-:--"

			//output, err := lib.ExecuteCmd("youtube-dl", "--skip-download", "--get-duration", "-i", feedItem.Link)
			output, err := lib.ExecuteCmd("python3", "-m", "youtube_dl", "--skip-download", "--get-duration", "-i", feedItem.Link)
			//output, err := lib.ExecuteCmd("python3", "youtube_dl_helper.py", feedItem.Link)

			if err != nil {
				lib.LogError("Can't get video duration with YouTube-DL from video: " + feedItem.Link + "\n" + strings.TrimSpace(strings.Replace(output, "\n", " ", -1)))
			} else {
				rssItem.videoDuration = strings.TrimSpace(output)

				if !strings.Contains(rssItem.videoDuration, ":") {
					rssItem.videoDuration = "0:" + rssItem.videoDuration
				}
			}
		}

		rssItems = append(rssItems, rssItem)
	}

	return rssItems
}

func notifyEvents(rssSource *RssSource, rssItem RssItem) error {
	telegramChatID, err := lib.TelegramGetChatID(string(rssSource.channelType))

	if err != nil {
		lib.LogError("Telegram chat ID not found for channel ID: " + string(rssSource.channelType) + ", URL: " + rssItem.link + "\n" + err.Error())
		return err
	}

	if rssSource.urlType == YoutubeChannel || rssSource.urlType == YoutubePlaylist {
		msg := rssItem.title + "\n[" + rssItem.videoDuration + "] " + rssSource.desc + "\n" + strings.TrimPrefix(rssItem.link, "https://") + "\n" + rssItem.publishedDate

		lib.LogDebug("    " + "Send Telegram YouTube video message to: " + string(rssSource.channelType) + " (" + strings.TrimPrefix(rssItem.link, "https://") + ")")
		err := lib.TelegramSendImageMessage(telegramChatID, rssItem.thumbnailUrl, msg)

		if err != nil {
			lib.LogError("Unable to send YouTube event to Telegram: " + rssItem.link + "\n" + err.Error())
			return err
		} else {
			rssUrlLog = append(rssUrlLog, rssItem.link)
			telegramRssItemSendCounter += 1
		}
	} else {
		msg := rssItem.title + "\n" + rssSource.desc + "\n" + strings.TrimPrefix(rssItem.link, "https://")

		if rssItem.publishedDate != "" {
			msg = msg + "\n" + rssItem.publishedDate
		}

		lib.LogDebug("    " + "Send Telegram RSS feed message to: " + string(rssSource.channelType) + " (" + strings.TrimPrefix(rssItem.link, "https://") + ")")
		err := lib.TelegramSendMessage(telegramChatID, msg, false)

		if err != nil {
			lib.LogError("Unable to send RSS event to Telegram: " + rssItem.link + "\n" + err.Error())
			return err
		} else {
			rssUrlLog = append(rssUrlLog, rssItem.link)
			telegramRssItemSendCounter += 1
		}
	}

	return nil
}
