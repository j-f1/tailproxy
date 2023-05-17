# v1.1.0

- raw TCP can now be proxied
  - set the target using the `tcp://` URL scheme to enable this
  - Funnel is not currently supported for raw TCP connections because we don’t have code to terminate TLS yet (and I don’t need it for my use case); feel free to open an issue if you need this!
- upgrade Tailscale dep to v1.40.1

# v1.0.1

- `/data` will automatically be used as the data directory if it exists and no other directory is specified
- Containers are now built with [`ko`](https://github.com/ko-build/ko)
- Containers are now available for virtually any architecture, not just arm64. Oops!

# v1.0.0

Initial release!
