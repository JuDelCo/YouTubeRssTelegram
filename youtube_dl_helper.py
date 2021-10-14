
import sys
import time
import youtube_dl

class MyLogger(object):
	def __init__(self, log_function):
		self.log_function = log_function
	def debug(self, msg):
		pass
	def warning(self, msg):
		pass
	def error(self, msg):
		self.log_function('YouTube Downloader ERROR: ' + msg)

def create(log_function):
	logger = MyLogger(log_function)

	options = {
		'quiet': True,
		'noplaylist': True,
		'verbose': False,
		'debug_printtraffic': False,
		'logger': logger
	}

	return youtube_dl.YoutubeDL(options)

def main():
	if(len(sys.argv) != 2):
		print('No YouTube URL provided')
		sys.exit(1)
		return

	youtube_handler = create(log_function=print)
	yt_link = str(sys.argv[1])

	try:
		dictMeta = youtube_handler.extract_info(yt_link, download=False)

		if(dictMeta['duration'] >= 3600):
			print(time.strftime('%-H:%M:%S', time.gmtime(dictMeta['duration'])))
		else:
			print(time.strftime('%-M:%S', time.gmtime(dictMeta['duration'])))
	except youtube_dl.utils.YoutubeDLError as ex:
		print('Cant get video duration with YouTube-DL from video ' + yt_link + ': ' + str(ex))
		sys.exit(2)

if __name__ == '__main__':
    main()
