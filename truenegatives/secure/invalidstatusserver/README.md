# A non-standard HTTP Server with Invalid HTTP Status

This HTTP server always replies with `HTTP/1.0 001 \n\nHello` which contains an invalid `001` HTTP status code. Client has to handle the error gracefully when reading the response.

## Run Non-Standard HTTP Server Locally

The instruction below exposes the HTTP service at port 8080.

```sh
git clone https://github.com/google/tsunami-security-scanner-testbed.git
cd tsunami-security-scanner-testbed
bazel run invalid_status_server_go_image
```
