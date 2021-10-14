package main

func main() {
	initialize()

	loadRssSources()
	loadRssCheckLog()
	loadRssUrlLog()

	processRssSources()

	saveRssCheckLog()
	saveRssUrlLog()
}
