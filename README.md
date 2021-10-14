
How to build:
=====================

```
go mod tidy
go build
strip youtube_rss_telegram
upx --brute youtube_rss_telegram
```

Alternative (for **Raspberry Pi Zero**)

```
env GOOS=linux GOARCH=arm GOARM=6 go build -ldflags="-s -w"
upx --brute youtube_rss_telegram
```

How to install:
=====================

```
sudo apt install python3-pip
sudo pip3 install youtube-dl
sudo pip3 install --upgrade youtube-dl
```

- Copy data_example files to data folder
- Update config CSV files in data folder
- Setup cron job
