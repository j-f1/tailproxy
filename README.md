# tailproxy

A proxy server that makes it easy to connect a local HTTP server to a [Tailscale](https://tailscale.com) network. Optimized for containerized usage but can also be used as a command line tool.

## Usage

### Command Line

```bash
tailproxy myhost localhost:3000
```

This will prompt you to approve a new device on your tailnet called `myhost`. Once approved (by clicking the printed link), you can visit `http://myhost` in a browser and it will proxy all requests to the the server listening at `localhost:3000`.

### Docker

Use the package `ghcr.io/j-f1/tailproxy` in your `docker-compose.yml`:

```yaml
version: '3'
services:
  tailproxy:
    image: ghcr.io/j-f1/tailproxy:edge
    environment:
      - TAILPROXY_TAILNET_HOST=myhost
      - TAILPROXY_TARGET=server:8080
      - TS_AUTHKEY=${TS_AUTHKEY}
    volumes:
      - ./tailproxy-data:/home/nonroot/data
    links:
      - server
  server:
    container_name: server
    image: my-server-image
```

Make sure to set a valid `TS_AUTHKEY` environment variable (see below) when running `docker compose up` to ensure that the proxy can join without requiring manual approval. While not required, it’s recommended to mount a volume to `/home/nonroot/data` so that the proxy can persist its state between restarts, including SSL certificates. (Otherwise, you’ll have to wait for a new certificate to be generated every time you restart the proxy.)

## Configuration 

You  are required to provide the following options:

- The machine name (env variable: `TAILPROXY_TAILNET_HOST` or first argument to the CLI) to join your tailnet as
- The target to proxy to (env variable: `TAILPROXY_TARGET` or second argument to the CLI). Format it as `host` (to use the default port 80) `host:port` or `host:port/basepath` (in which case `/basepath` will be prepended to all requests to the upstream server)

Additionally, you can set any environment variables that are supported by Tailscale. You’ll most likely want to set the `TS_AUTHKEY` environment variable to a valid [auth key](https://tailscale.com/kb/1085/auth-keys/) so that you don’t have to click the link to approve the new device every time you restart the proxy. Make sure to configure the auth key to provision ephemeral and pre-approved devices when creating it for the smoothest experience.

You may optionally set `TAILPROXY_DATA_DIR` to a directory where the proxy can store its state. Currently, the is just the Tailscale state (which is placed in the `tailscale` subdirectory of the directory you provide). If you don’t set this, Tailscale will use a subdirectory named `tsnet-tailproxy` in Go’s `os.UserConfigDir`.

### HTTPS

You can optionally  pass an option to enable HTTPS support (`--https` in the CLI or `TAILPROXY_HTTPS_MODE` as an environment variable). The following values are allowed:

- `off` (default): No HTTPS support. The proxy will only listen on port 80.
- `redirect`: The proxy will listen on both port 80 and port 443. Any HTTP request will be redirected to HTTPS.
- `only`: The proxy will listen for HTTPS requests on port 443. HTTP requests will not be accepted.
- `both`: The proxy will listen for HTTP and HTTPS requests on ports 80 and 443 respectively. It will not redirect HTTP requests to HTTPS.

If HTTPS is enabled, tailproxy will use Tailscale’s API to generate a valid certificate for `<host>.<tailnet name>.ts.net`. It will strip HTTPS and forward the plain HTTP request to the upstream server.

### Debugging/Profiling

If you pass `TAILPROXY_PPROF_ENABLED=1`, the proxy will expose a pprof server on port 6060 (on your tailnet). You can use this to debug performance issues or to profile the proxy.
