package main

import (
	"flag"
	"fmt"
	"net/http"
	"os"
	"path"
	"slices"
	"strconv"
	"strings"
	"time"
)

func main() {
	servePath := flag.String("path", ".", "Path to serve")
	flag.Parse()

	dirList, err := os.ReadDir(*servePath)
	if err != nil {
		panic(err)
	}

	files := make(map[string]string)
	for _, e := range dirList {
		if !e.IsDir() {
			if path.Ext(e.Name()) == ".swf" {
				files["/"+e.Name()] = path.Join(*servePath, e.Name())
				files["/"+e.Name()+".swf2ass.mkv"] = path.Join(*servePath, e.Name()+".swf2ass.mkv")
				files["/"+e.Name()+".swf2ass.mkv.zstd"] = path.Join(*servePath, e.Name()+".swf2ass.mkv.zstd")
				files["/"+e.Name()+".swf2ass.mkv.br"] = path.Join(*servePath, e.Name()+".swf2ass.mkv.br")
				files["/"+e.Name()+".swf2ass.mkv.gzip"] = path.Join(*servePath, e.Name()+".swf2ass.mkv.gzip")
			}

		}
	}

	server := &http.Server{
		Addr:        "0.0.0.0:8008",
		ReadTimeout: time.Second * 2,
		Handler: http.HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
			if request.URL.Path == "/" {
				writer.Header().Set("Content-Type", "text/html; charset=utf-8")

				var output string
				output += "<html><body><table><tr><th align=\"right\">NAME</th><th>Size</th><th>Links</th></tr>"
				for _, e := range dirList {
					if path.Ext(e.Name()) == ".swf" {
						stat, err := os.Stat(path.Join(*servePath, e.Name()+".swf2ass.mkv"))
						if err != nil {
							continue
						}
						output += fmt.Sprintf("<tr><th align=\"right\">%s</th><td>%d MiB</td><td><a href=\"/%s.swf2ass.mkv\">MKV</a> <a href=\"/%s.swf2ass.mkv.zstd\">Zstandard</a> <a href=\"/%s.swf2ass.mkv.br\">Brotli</a> <a href=\"/%s.swf2ass.mkv.gzip\">gzip</a> <a href=\"/%s\">SWF</a></td></tr>", e.Name(), stat.Size()/(1024*1024), e.Name(), e.Name(), e.Name(), e.Name(), e.Name())
					}
				}
				output += "</table></body></html>"
				writer.Header().Set("Content-Length", strconv.FormatUint(uint64(len([]byte(output))), 10))
				writer.WriteHeader(http.StatusOK)
				writer.Write([]byte(output))
				return
			}

			if fpath, ok := files[request.URL.Path]; !ok {
				writer.WriteHeader(http.StatusNotFound)
				return
			} else {

				encodings := strings.Split(request.Header.Get("Accept-Encoding"), ",")
				for i := range encodings {
					//drop preference
					e := strings.Split(encodings[i], ";")
					encodings[i] = strings.TrimSpace(e[0])
				}
				writer.Header().Set("Content-Type", "application/octet-stream")

				if path.Ext(request.URL.Path) == ".mkv" {
					writer.Header().Set("Content-Type", "video/x-matroska")

					ua := request.UserAgent()
					//some players allow decoding gzip even if it's not on accept-encoding
					//  || strings.Contains(ua, "LibVLC/")
					if (strings.Contains(ua, "libmpv") || strings.Contains(ua, "Lavf/")) && !slices.Contains(encodings, "gzip") {
						encodings = append(encodings, "gzip")
					}

					if slices.Contains(encodings, "zstd") {
						writer.Header().Set("Content-Encoding", "zstd")
						fpath += ".zstd"
						request.Header.Del("Range")
					} else if slices.Contains(encodings, "br") {
						writer.Header().Set("Content-Encoding", "br")
						fpath += ".br"
						request.Header.Del("Range")
					} else if slices.Contains(encodings, "gzip") {
						writer.Header().Set("Content-Encoding", "gzip")
						fpath += ".gz"
						request.Header.Del("Range")
					} else {
						//serve line-compressed mkv instead
						fpath = strings.TrimSuffix(fpath, ".mkv") + ".zlib.mkv"
					}
				}

				stat, err := os.Stat(fpath)
				if err != nil {
					panic(err)
				}
				f, err := os.Open(fpath)
				if err != nil {
					panic(err)
				}
				defer f.Close()
				writer.Header().Set("Content-Length", strconv.FormatUint(uint64(stat.Size()), 10))
				http.ServeContent(writer, request, path.Base(request.URL.Path), stat.ModTime(), f)
			}

		}),
	}

	err = server.ListenAndServe()
	if err != nil {
		panic(err)
	}
}
