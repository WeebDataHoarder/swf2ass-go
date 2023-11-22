# swf2ass ⚡➡️🍑

Converts Flash animations into ASS subtitles with vector drawings where possible. Work in progress, expect broken output.

## Usage

```bash
$ go run git.gammaspectra.live/WeebDataHoarder/swf2ass-go -input [file.swf] -output [file.ass] -audio [file.mp3]
```

The `-audio` parameter is optional and will only be produced if the input file has streaming MP3 audio.

Create a bogus video track with the subtitles and audio embedded:
```bash
$ FILENAME=file
$ ffmpeg -y \
-f lavfi -i "color=size=$(grep PlayResX ${FILENAME}.ass | head -n 1 | awk '{ print $2 }')x$(grep PlayResY ${FILENAME}.ass | head -n 1 | awk '{ print $2 }'):rate=$(grep '?dummy' ${FILENAME}.ass | head -n 1 | awk -F: '{ print $3 }'):color=black" \
-i "${FILENAME}.mp3" \
-i "${FILENAME}.ass" \
-map 0:v -map 1:a -map 2:s \
-c:v libx264 -pix_fmt yuv420p -crf 1 -tune stillimage -preset placebo -x264-params keyint=240 \
-c:a copy \
-c:s copy -disposition:s:0 forced -metadata:s:s:0 language=und \
-shortest "${FILENAME}.mkv"
```