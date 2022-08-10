# Testbed for Tsunami Security Scanner

This project aims to provide a central repository for example services and
applications for testing Tsunami Security Scanner plugins' detection capability.

## Run Test Servers Locally

**One-time setup**: [Install bazel](https://bazel.build/install)

Build different application images with bazel:

```sh
git clone https://github.com/google/tsunami-security-scanner-testbed
cd tsunami-security-scanner-testbed

bazel run fake_ssh_server_go_image
bazel run fake_https_server_go_image
bazel run invalid_status_server_go_image
```

## Contributing

Read how to [contribute to Tsunami](docs/contributing.md).

## Source Code Headers

Every file containing source code must include copyright and license
information. This includes any JS/CSS files that you might be serving out to
browsers. (This is to help well-intentioned people avoid accidental copying that
doesn't comply with the license.)

Apache header:

```
Copyright 2022 Google LLC

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    https://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
```

## Disclaimer

Tsunami Security Scanner testbed is not officially supported Google
products.
