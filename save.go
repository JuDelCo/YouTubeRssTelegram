package main

import (
	"gorsstelegram/lib"
	"os"
	"strconv"
)

func saveRssCheckLog() {
	lib.LogDebug("Saving RSS check log to: " + rssSourcesCheckLogPath)
	records := make([][]string, 0)

	for _, rssSource := range rssSources {
		row := make([]string, 0)
		row = append(row, strconv.FormatInt(rssSource.guid, 10))
		row = append(row, strconv.FormatInt(rssSource.lastUpdate, 10))

		records = append(records, row)
	}

	err := lib.WriteCsvFile(rssSourcesCheckLogNewPath, records)

	if err != nil {
		lib.LogError(err.Error())
		return
	}

	lines, err := lib.LinesInFile(rssSourcesCheckLogNewPath)

	if err != nil {
		lib.LogError(err.Error())
		return
	}

	newCount := len(lines)

	lines, err = lib.LinesInFile(rssSourcesCheckLogPath)

	if err != nil {
		lib.LogError(err.Error())
		return
	}

	oldCount := len(lines)

	if newCount >= oldCount {
		os.Rename(rssSourcesCheckLogNewPath, rssSourcesCheckLogPath)
		lib.LogDebug("    " + strconv.FormatInt(int64(len(records)), 10) + " RSS check logs saved")
	} else {
		lib.LogError("Not saving RSS timestamp check log because processing RSS halted")
	}
}

func saveRssUrlLog() {
	lib.LogDebug("Saving RSS URLs log to: " + rssUrlLogPath)
	records := make([][]string, 0)
	c := 0

	for _, rssUrl := range rssUrlLog {
		if c >= rssUrlLogMaxCount {
			break
		}

		row := make([]string, 0)
		row = append(row, rssUrl)

		records = append(records, row)

		c += 1
	}

	err := lib.WriteCsvFile(rssUrlLogPath, records)

	if err != nil {
		lib.LogError(err.Error())
		return
	}

	lib.LogDebug("    " + strconv.FormatInt(int64(len(records)), 10) + " RSS URLs log saved")
}
