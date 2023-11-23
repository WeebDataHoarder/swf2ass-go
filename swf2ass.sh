#!/bin/bash

INPUT="${1}"

go run git.gammaspectra.live/WeebDataHoarder/swf2ass-go -input "${INPUT}" -output "${INPUT}.ass" -audio "${INPUT}.mp3"

ffmpeg -y \
-f lavfi -i "color=size=$(grep PlayResX ${INPUT}.ass | head -n 1 | awk '{ print $2 }')x$(grep PlayResY ${INPUT}.ass | head -n 1 | awk '{ print $2 }'):rate=$(grep '?dummy' ${INPUT}.ass | head -n 1 | awk -F: '{ print $3 }'):color=black" \
-i "${INPUT}.mp3" \
-map 0:v -map 1:a \
-vf "drawtext=box=1:boxcolor=black@0.5:boxborderw=20:x=(w-text_w)/2:y=(h-text_h)/2:fontcolor=white:fontsize=32:text='Enable ASS subtitles'" \
-c:v libx264 -pix_fmt yuv420p -crf 1 -tune stillimage -preset slow -x264-params keyint=240 \
-c:a copy \
-shortest "${INPUT}.video.mkv"

mkvmerge --title "${INPUT}" -o "${INPUT}.swf2ass.mkv" \
--attachment-mime-type "application/x-shockwave-flash" --attachment-description "SWF Source" --attach-file "${INPUT}" \
--language 0:zxx --track-name 0:"bogus video" "${INPUT}.video.mkv" \
--forced-track 0:1 --default-track 0:1 --compression 0:zlib --language 0:zxx --track-name 0:"Vector from ${INPUT}" "${INPUT}.ass"


rm "${INPUT}.ass" "${INPUT}.video.mkv" "${INPUT}.mp3"