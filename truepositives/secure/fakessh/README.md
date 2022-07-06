# Fake SSH Server with Weak Password

This example SSH server is configured to use a known weak password (default to
'qwerty') to authenticate root user. When scanned by Tsunami Scanner, this
service will trigger a weak password finding from its [Ncrack Weak Credential Detector](https://github.com/google/tsunami-security-scanner-plugins/tree/master/google/detectors/credentials/ncrack).

This server terminates connections right after successful handshake. It is not
configured to any local bash shell session.

## Run Fake SSH Server Locally

The instruction below exposes the fake SSH service at port 8022.

```sh
git clone https://github.com/google/tsunami-security-scanner-testbed.git
cd tsunami-security-scanner-testbed
bazel run //:fake_ssh_server_go_image
```

