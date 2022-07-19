# Fake HTTPS Server

This is an example HTTPS server. When scanned by Tsunami Scanner, this service
will trigger a directory traversal finding from its
[Apache HTTP Server 2.4.49 Path traversal and disclosure vulnerability - CVE_2021_41773]

This server terminates connections right after a successful handshake.

## Run Fake HTTPS Server Locally

The instruction below exposes the fake HTTPS service at port 8443.

```sh
git clone https://github.com/google/tsunami-security-scanner-testbed.git
cd tsunami-security-scanner-testbed
bazel run //:fake_https_server_go_image
```
