package main

import (
	"gorsstelegram/lib"
	"math"
	"sort"
	"strconv"
	"strings"
)

func initialize() {
	lib.LogDebug("Load App settings from: " + appSettingsPath)
	err := loadAppSettings(appSettingsPath)

	if err != nil {
		lib.LogErrorFatal(err.Error())
	}

	lib.LogDebug("Setup time locale to: " + timeLocation)
	err = lib.SetupTimeLocale(timeLocation)

	if err != nil {
		lib.LogErrorFatal(err.Error())
	}

	lib.LogDebug("Load Telegram chat IDs from: " + telegramIDsPath)
	err = lib.TelegramLoadChatIDs(telegramIDsPath)

	if err != nil {
		lib.LogErrorFatal(err.Error())
	}

	lib.LogDebug("Init Telegram bot API")
	err = lib.InitializeTelegramBot(telegramApiTokenPath)

	if err != nil {
		lib.LogErrorFatal(err.Error())
	}
}

func loadAppSettings(appSettingsPath string) error {
	records, err := lib.ReadCsvFile(appSettingsPath)

	if err != nil {
		return err
	}

	for _, line := range records {
		settingId := strings.TrimSpace(line[0])

		if settingId == "IntervalInMinutesForEachRun" {
			intervalInMinutesForEachRun, err = strconv.Atoi(line[1])
		} else if settingId == "TimeLocation" {
			timeLocation = strings.TrimSpace(line[1])
		} else if settingId == "RssSourcesProcessMaxCount" {
			rssSourcesProcessMaxCount, err = strconv.Atoi(line[1])
		} else if settingId == "WaitMinutesToRefreshAgainRssSource" {
			waitMinutesToRefreshAgainRssSource, err = strconv.Atoi(line[1])
		} else if settingId == "AgeInDaysToIgnoreRssItem" {
			ageInDaysToIgnoreRssItem, err = strconv.Atoi(line[1])
		} else if settingId == "RSSUrlLogMaxCount" {
			rssUrlLogMaxCount, err = strconv.Atoi(line[1])
		}
	}

	lib.LogDebug("    " + "Interval in minutes for each run: " + strconv.FormatInt(int64(intervalInMinutesForEachRun), 10))
	lib.LogDebug("    " + "Time Location: " + timeLocation)
	lib.LogDebug("    " + "RSS sources process max count: " + strconv.FormatInt(int64(rssSourcesProcessMaxCount), 10))
	lib.LogDebug("    " + "Wait minutes to refresh again a RSS source: " + strconv.FormatInt(int64(waitMinutesToRefreshAgainRssSource), 10))
	lib.LogDebug("    " + "Age in days to ignore a RSS item: " + strconv.FormatInt(int64(ageInDaysToIgnoreRssItem), 10))
	lib.LogDebug("    " + "RSS URL log max count: " + strconv.FormatInt(int64(rssUrlLogMaxCount), 10))

	return err
}

func loadRssSources() {
	lib.LogDebug("Loading RSS sources file: " + rssSourcesPath)
	records, err := lib.ReadCsvFile(rssSourcesPath)
	enabledSources := 0

	if err != nil {
		lib.LogErrorFatal(err.Error())
	}

	for _, line := range records {
		rssSource := RssSource{}
		rssSource.guid, err = strconv.ParseInt(line[0], 10, 0)
		rssSource.urlType = UrlType(strings.TrimSpace(line[1]))
		rssSource.url = strings.TrimSpace(line[2])
		rssSource.channelType = ChannelType(strings.TrimSpace(line[3]))
		rssSource.desc = strings.TrimSpace(line[4])
		rssSource.enabled, err = strconv.ParseBool(line[5])
		rssSource.lastUpdate = 0

		if rssSource.urlType == YoutubeChannel {
			rssSource.url = ytChannelUrlPrefix + rssSource.url
		} else if rssSource.urlType == YoutubePlaylist {
			rssSource.url = ytPlaylistUrlPrefix + rssSource.url
		}

		if err != nil {
			lib.LogErrorFatal(err.Error())
		}

		if rssSource.enabled {
			enabledSources += 1
		}

		rssSources = append(rssSources, rssSource)
	}

	lib.LogDebug("    " + strconv.FormatInt(int64(len(rssSources)), 10) + " RSS sources loaded, " + strconv.FormatInt(int64(enabledSources), 10) + " are enabled")

	totalRssSourcesSupport := int64(math.Max(1, math.Floor(float64(waitMinutesToRefreshAgainRssSource)/float64(intervalInMinutesForEachRun)))) * int64(rssSourcesProcessMaxCount)
	lib.LogDebug("    " + strconv.FormatInt(totalRssSourcesSupport, 10) + " total RSS sources support with current settings")

	if totalRssSourcesSupport < int64(enabledSources) {
		lib.LogError("There are currently " + strconv.FormatInt(int64(enabledSources), 10) + " enabled RSS sources, but current settings only support up to " + strconv.FormatInt(int64(totalRssSourcesSupport), 10))
	}
}

func loadRssCheckLog() {
	lib.LogDebug("Loading RSS check log file: " + rssSourcesCheckLogPath)
	records, err := lib.ReadCsvFile(rssSourcesCheckLogPath)
	updateCounter := 0
	missingCounter := 0

	if err != nil {
		lib.LogErrorFatal(err.Error())
	}

	for _, line := range records {
		guid, err := strconv.ParseInt(line[0], 10, 0)

		i := sort.Search(len(rssSources), func(i int) bool { return rssSources[i].guid >= guid })

		if i < len(rssSources) && rssSources[i].guid == guid {
			rssSources[i].lastUpdate, err = strconv.ParseInt(line[1], 10, 0)
			updateCounter += 1
		} else {
			rssSource := RssSource{}
			rssSource.guid = guid
			rssSource.enabled = false
			rssSource.lastUpdate, err = strconv.ParseInt(line[1], 10, 0)

			// Insert rssSource in rssSources in the specified index (i)
			// slice[startIndex:length]
			rssSources = append(rssSources[:i+1], rssSources[i:]...)
			rssSources[i] = rssSource
			missingCounter += 1
		}

		if err != nil {
			lib.LogErrorFatal(err.Error())
		}
	}

	lib.LogDebug("    " + strconv.FormatInt(int64(updateCounter), 10) + " RSS sources updated (last update property)")

	if missingCounter > 0 {
		lib.LogDebug("    " + strconv.FormatInt(int64(missingCounter), 10) + " missing RSS sources found (not important)")
	}
}

func loadRssUrlLog() {
	lib.LogDebug("Loading RSS notified URLs log file: " + rssUrlLogPath)
	records, err := lib.ReadCsvFile(rssUrlLogPath)

	if err != nil {
		lib.LogErrorFatal(err.Error())
	}

	for _, line := range records {
		rssUrlLog = append(rssUrlLog, strings.TrimSpace(line[0]))
	}

	sort.Strings(rssUrlLog)

	lib.LogDebug("    " + strconv.FormatInt(int64(len(rssUrlLog)), 10) + " RSS URLs loaded")
}
