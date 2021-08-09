.DEFAULT_GOAL = all

.PHONY: all clean 


all:: dependencies
all:: build


dependencies:
	go get

build: grabber_config.json
	go build -ldflags="-s -w" -o builds/grabber_linux_x64 main.go

grabber_config.json:
	mkdir --parents builds
	echo -e '{\n  "savepath": "Manga/",\n  "readmanga": {\n    "timeout_image": 500,\n    "timeout_chapter": 1000\n  },\n  "mangalib": {\n    "timeout_image": 700,\n    "timeout_chapter": 2000\n  }\n}' > builds/grabber_config.json

clean:
	git clean -xddff
