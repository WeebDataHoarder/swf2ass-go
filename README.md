# swf2ass ⚡➡️🍑

Converts Flash animations into ASS subtitles with vector drawings where possible. Work in progress, expect broken output.

## Usage

```bash
$ go run git.gammaspectra.live/WeebDataHoarder/swf2ass-go -input [file.swf] -output [file.ass] -audio [file.mp3]
```

The `-audio` parameter is optional and will only be produced if the input file has streaming MP3 audio.