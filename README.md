
How to build:
=====================

```
go mod tidy
go build
strip youtube_rss_telegram
upx --brute youtube_rss_telegram
```

Alternative (for **Linux AMD64 Server**)

```
env GOOS=linux GOARCH=amd64 go build -ldflags="-s -w"
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
sudo pip3 install yt-dlp
sudo pip3 install --upgrade yt-dlp
```

- Copy data_example files to data folder
- Update config CSV files in data folder
- Setup cron job

3rd-party dependencies docs:
=====================

- https://pkg.go.dev/github.com/mmcdole/gofeed@v1.1.3
- https://github.com/tucnak/telebot
