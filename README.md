# HTTP Grab

[![Go Reference](https://pkg.go.dev/badge/github.com/WangYihang/http-grab.svg)](https://pkg.go.dev/github.com/WangYihang/http-grab)
[![Go Report Card](https://goreportcard.com/badge/github.com/WangYihang/http-grab)](https://goreportcard.com/report/github.com/WangYihang/http-grab)

## Description

`http-grab` is a tool for grabbing HTTP response from a list of IP addresses.

## Installation

```bash
go install github.com/WangYihang/http-grab@latest
```

## Usage

```bash
$ http-grab --h
Usage:
  http-grab [OPTIONS]

Application Options:
  -i, --input=          input file path
  -o, --output=         output file path
  -s, --status-updates= status updates file path
  -n, --num-workers=    number of workers (default: 32)
      --seed=           seed (default: 0)
      --num-shards=     number of shards (default: 1)
      --shard=          shard (default: 0)
  -p, --port=           port (default: 80)
      --path=           path (default: index.html)
      --host=           http host header
  -m, --max-tries=      max tries (default: 4)
  -t, --timeout=        timeout (default: 8)

Help Options:
  -h, --help            Show this help message
```

```bash
$ head input.txt
23.63.66.161
3.210.226.220
34.149.9.201
38.91.55.188
70.109.57.175
108.138.69.232
3.144.94.58
142.202.80.211
104.233.202.168
23.202.84.42
```

```bash
$ http-grab -i input.txt -o output.txt
...
```

```bash
$ head -n 1 output.txt
```

```json
{
    "index": 26,
    "started_at": 1706764512210,
    "finished_at": 1706764512755,
    "num_tries": 1,
    "timeout": 8,
    "error": "",
    "ip": "34.149.112.180",
    "port": 80,
    "path": "index.html",
    "host": "34.149.112.180",
    "http": {
        "request": {
            "method": "GET",
            "url": "http://34.149.112.180:80/index.html",
            "host": "34.149.112.180",
            "remote_addr": "",
            "request_uri": "",
            "proto": "HTTP/1.1",
            "proto_major": 1,
            "proto_minor": 1,
            "header": {
                "User-Agent": [
                    "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/121.0.0.0 Safari/537.36 Edg/121.0.0.0"
                ]
            },
            "content_length": 0,
            "transfer_encoding": null,
            "close": false,
            "form": null,
            "post_form": null,
            "multipart_form": null,
            "trailer": null
        },
        "response": {
            "status": "404 Not Found",
            "status_code": 404,
            "proto": "HTTP/1.1",
            "proto_major": 1,
            "proto_minor": 1,
            "header": {
                "Content-Length": [
                    "42"
                ],
                "Content-Type": [
                    "text/plain; charset=UTF-8"
                ],
                "Date": [
                    "Thu, 01 Feb 2024 05:15:12 GMT"
                ],
                "Server": [
                    "akka-http/10.2.7"
                ],
                "Via": [
                    "1.1 google"
                ]
            },
            "raw_body": "VGhlIHJlcXVlc3RlZCByZXNvdXJjZSBjb3VsZCBub3QgYmUgZm91bmQu",
            "body": "The requested resource could not be found.",
            "content_length": 42,
            "transfer_encoding": null,
            "close": false,
            "uncompressed": false,
            "trailer": null
        }
    }
}
```