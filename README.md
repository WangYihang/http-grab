# HTTP Grab

## Description

`http-grab` is a tool for grabbing HTTP response from a list of IP addresses.

## Installation

```bash
go install github.com/WangYihang/http-grab@latest
```

## Usage

```
$ http-grab --h
Usage:
  http-grab [OPTIONS]

Application Options:
  -i, --input=       input file path
  -o, --output=      output file path
  -n, --num-workers= number of workers (default: 32)
  -t, --timeout=     timeout (default: 8)
  -p, --port=        port (default: 80)
  -P, --path=        path (default: index.html)
  -H, --host=        host

Help Options:
  -h, --help         Show this help message
```

```
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

```
$ http-grab -i input.txt -o output.txt
...
```

```json
$ head -n 1 output.txt
{"index":8,"started_at":1706763874999,"finished_at":1706763875591,"num_tries":1,"timeout":8,"error":"","ip":"104.233.202.168","port":80,"path":"verification.html","host":"104.233.202.168","http":{"request":{"method":"GET","url":"http://104.233.202.168:80/verification.html","host":"104.233.202.168","remote_addr":"","request_uri":"","proto":"HTTP/1.1","proto_major":1,"proto_minor":1,"header":{"User-Agent":["Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/121.0.0.0 Safari/537.36 Edg/121.0.0.0"]},"content_length":0,"transfer_encoding":null,"close":false,"form":null,"post_form":null,"multipart_form":null,"trailer":null},"response":{"status":"403 Forbidden","status_code":403,"proto":"HTTP/1.1","proto_major":1,"proto_minor":1,"header":{"Connection":["keep-alive"],"Content-Length":["548"],"Content-Type":["text/html; charset=utf-8"],"Date":["Thu, 01 Feb 2024 05:04:35 GMT"],"Server":["nginx"]},"raw_body":"PGh0bWw+DQo8aGVhZD48dGl0bGU+NDAzIEZvcmJpZGRlbjwvdGl0bGU+PC9oZWFkPg0KPGJvZHk+DQo8Y2VudGVyPjxoMT40MDMgRm9yYmlkZGVuPC9oMT48L2NlbnRlcj4NCjxocj48Y2VudGVyPm5naW54PC9jZW50ZXI+DQo8L2JvZHk+DQo8L2h0bWw+DQo8IS0tIGEgcGFkZGluZyB0byBkaXNhYmxlIE1TSUUgYW5kIENocm9tZSBmcmllbmRseSBlcnJvciBwYWdlIC0tPg0KPCEtLSBhIHBhZGRpbmcgdG8gZGlzYWJsZSBNU0lFIGFuZCBDaHJvbWUgZnJpZW5kbHkgZXJyb3IgcGFnZSAtLT4NCjwhLS0gYSBwYWRkaW5nIHRvIGRpc2FibGUgTVNJRSBhbmQgQ2hyb21lIGZyaWVuZGx5IGVycm9yIHBhZ2UgLS0+DQo8IS0tIGEgcGFkZGluZyB0byBkaXNhYmxlIE1TSUUgYW5kIENocm9tZSBmcmllbmRseSBlcnJvciBwYWdlIC0tPg0KPCEtLSBhIHBhZGRpbmcgdG8gZGlzYWJsZSBNU0lFIGFuZCBDaHJvbWUgZnJpZW5kbHkgZXJyb3IgcGFnZSAtLT4NCjwhLS0gYSBwYWRkaW5nIHRvIGRpc2FibGUgTVNJRSBhbmQgQ2hyb21lIGZyaWVuZGx5IGVycm9yIHBhZ2UgLS0+DQo=","body":"\u003chtml\u003e\r\n\u003chead\u003e\u003ctitle\u003e403 Forbidden\u003c/title\u003e\u003c/head\u003e\r\n\u003cbody\u003e\r\n\u003ccenter\u003e\u003ch1\u003e403 Forbidden\u003c/h1\u003e\u003c/center\u003e\r\n\u003chr\u003e\u003ccenter\u003enginx\u003c/center\u003e\r\n\u003c/body\u003e\r\n\u003c/html\u003e\r\n\u003c!-- a padding to disable MSIE and Chrome friendly error page --\u003e\r\n\u003c!-- a padding to disable MSIE and Chrome friendly error page --\u003e\r\n\u003c!-- a padding to disable MSIE and Chrome friendly error page --\u003e\r\n\u003c!-- a padding to disable MSIE and Chrome friendly error page --\u003e\r\n\u003c!-- a padding to disable MSIE and Chrome friendly error page --\u003e\r\n\u003c!-- a padding to disable MSIE and Chrome friendly error page --\u003e\r\n","content_length":548,"transfer_encoding":null,"close":false,"uncompressed":false,"trailer":null}}}
```