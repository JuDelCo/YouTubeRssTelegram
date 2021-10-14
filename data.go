package main

var rssSources = make([]RssSource, 0)
var rssUrlLog = make([]string, 0)
var rssSourcesProcessCounter = 0
var telegramRssItemSendCounter = 0

var intervalInMinutesForEachRun = 5
var timeLocation = "Europe/Madrid"
var rssSourcesProcessMaxCount = 25
var waitMinutesToRefreshAgainRssSource = 120
var ageInDaysToIgnoreRssItem = 7
var rssUrlLogMaxCount = 1024
