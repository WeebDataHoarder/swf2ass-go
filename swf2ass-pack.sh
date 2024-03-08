#!/bin/bash

INPUT="${1}"
OUTPUT="${INPUT}.swf2ass.mkv"
OUTPUT_ZLIB="${INPUT}.swf2ass.zlib.mkv"
OUTPUT_CLEAN="${INPUT}.swf2ass.clean.mkv"

./swf2ass-zlib.sh "${INPUT}"

if [[ ! -f "${OUTPUT}" ]]; then
  echo "Output not found!"
  exit 1
fi

#mkclean --remux --keep-cues --live "${OUTPUT}" "${OUTPUT_CLEAN}"
mkclean --remux --keep-cues "${OUTPUT}" "${OUTPUT_CLEAN}"

rm -v "${OUTPUT}"
mv -v "${OUTPUT_CLEAN}" "${OUTPUT}"

zstd -f -k -19 -T0 "${OUTPUT}" -o "${OUTPUT}.zst"

brotli --keep --no-copy-stat --best --verbose --lgwin=24 --stdout "${OUTPUT}" > "${OUTPUT}.br"

#gzip --stdout --no-name --best "${OUTPUT}" > "${OUTPUT}.gz"
zopfli --gzip --i10 -c "${OUTPUT}" > "${OUTPUT}.gz"

