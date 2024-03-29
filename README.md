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
    image: ghcr.io/j-f1/tailproxy:v1
    environment:
      - TAILPROXY_NAME=myhost
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

Make sure to set a valid `TS_AUTHKEY` environment variable (see below) when running `docker compose up` to ensure that the proxy can join without requiring manual approval. While not required, it’s recommended to mount a volume to `/home/nonroot/data` so that the proxy can persist its state between restarts, including TLS certificates. (Otherwise, you’ll have to wait for a new certificate to be generated every time you restart the proxy.)

## Configuration 

You  are required to provide the following options:

- The machine name (env variable: `TAILPROXY_NAME` or first argument to the CLI) to join your tailnet as. Note that Tailscale will automatically add a `-<number>` suffix to the name if it’s already taken.
- The target to proxy to (env variable: `TAILPROXY_TARGET` or second argument to the CLI). Format it as `host` (to use the default port 80) `host:port` or `host:port/basepath?foo=bar` (in which case `/basepath` will be prepended to all requests to the upstream server and `?foo=bar` will be prepended to the query string of all requests).

Additionally, you can set any environment variables that are supported by Tailscale. You’ll most likely want to set the `TS_AUTHKEY` environment variable to a valid [auth key](https://tailscale.com/kb/1085/auth-keys/) so that you don’t have to click the link to approve the new device every time you restart the proxy. Make sure to configure the auth key to provision pre-approved devices when creating it for the smoothest experience.

You may optionally set `TAILPROXY_DATA_DIR` to a directory where the proxy can store its state. Currently, we’re just storing the Tailscale state (which is placed in the `tailscale` subdirectory of the directory you provide). If you don’t set this, Tailscale will use `/data` if it exists, or a subdirectory named `tsnet-tailproxy` in Go’s `os.UserConfigDir` if `/data` does not exist.

### Raw TCP proxying

If you write your target as `tcp://host:port`, tailproxy will proxy TCP connections to the specified host and port. This is useful if you want to proxy a non-HTTP server, such as a database. The port is required since there isn’t a default port for TCP connections.

Note that HTTPS/TLS and therefore Funnel are not supported because they require TLS termination, which has not been implemented for raw TCP connections. (PRs welcome! You would probably need to rethink the `https` option to make it work for raw TLS as well as HTTPS.)

### HTTPS

You can optionally pass an option to enable HTTPS support (`--https` in the CLI or `TAILPROXY_HTTPS_MODE` as an environment variable). The following values are allowed:

- `off` (default): No HTTPS support. The proxy will only listen on port 80.
- `redirect`: The proxy will listen on both port 80 and port 443. Any HTTP request will be redirected to HTTPS.
- `only`: The proxy will listen for HTTPS requests on port 443. HTTP requests will not be accepted.
- `both`: The proxy will listen for HTTP and HTTPS requests on ports 80 and 443 respectively. It will not redirect HTTP requests to HTTPS.

If HTTPS is enabled, tailproxy will use Tailscale’s API to generate a valid certificate for `<host>.<tailnet name>.ts.net`. It will strip HTTPS and forward the plain HTTP request to the upstream server.

### Funnel

> **Warning**: Tailproxy is relatively safe because it’s only accessible from devices you control. However, Funnel allows anyone from the Internet to talk to your server. That means you have to worry about both the security of Tailproxy and of your server. I don’t know about you, but I don’t really know what I’m doing, so besides the inherent safety of Go and the relative simplicity of my code I can’t guarantee that there aren’t any security issues. Use Funnel at your own risk.

You can optionally make the service behind tailproxy publicly accessible using [Tailscale Funnel](https://tailscale.com/kb/1223/tailscale-funnel/) (`--funnel` in the CLI or `TAILPROXY_FUNNEL_MODE` as an environment variable). The following values are allowed:

- `off` (default): No Funnel support. The proxy will only listen on your tailnet.
- `on`: The proxy will listen on both your tailnet and the public internet.
- `only`: The proxy will listen for requests only on the public internet. Requests on your tailnet will not be accepted.

Note that you’ll have to enable Funnel for your tailnet, and make sure that the tailproxy node has the `funnel` attribute in `nodeAttrs`. Funnel handles the TLS termination, so the HTTPS config option will be ignored for connections coming from Funnel.

### Debugging/Profiling

If you pass `TAILPROXY_PPROF_ENABLED=1` (or `--pprof`), the proxy will expose a pprof server on port 6060 (on your tailnet). You can use this to debug performance issues or to profile the proxy. Note that the pprof server is only available on your tailnet, even if Funnel is enabled.

## Proxy functionality

The proxy will set several headers on the request it sends to the upstream server based on information Tailscale provides about the device making the request:

- `X-Forwarded-For`: the IP address of the device making the request
- `X-Forwarded-Host`: the hostname of the proxy (according to the client making the request)
- `X-Forwarded-Proto`: the protocol used to make the request (either `http` or `https`)
- `X-Tailscale-WhoIs`:
  - `ok` if the request came from a Tailscale device. The below headers identify the device and its owner.
  - `funnel` if the request came from Funnel. The below headers will not be present, since Funnel does not authenticate the request.
  - `error` if the call to Tailscale failed. The below headers will not be present. This is unlikely to happen, and if you rely on authentication, you should probably return an error to the client in this case. We don’t handle this in tailproxy to give you flexibility to continue on if you don’t care about authentication. Check your server logs and report a bug if this happens to you!
- `X-Tailscale-User`: the unique user ID of the owner of the device making the request
- `X-Tailscale-User-LoginName`: the login name (`j-f1@github`) of the user
- `X-Tailscale-User-DisplayName` the display name (`Jed Fox`) of the user
- `X-Tailscale-User-ProfilePicURL` the URL of the user’s profile picture (if available)
- `X-Tailscale-Caps` a comma-separated list of capabilities the user has (if any)
- `X-Tailscale-Node` the unique ID of the device making the request
- `X-Tailscale-Node-Name` the machine name of the device making the request (with the tailnet name appended if the device has been shared from another tailnet)
- `X-Tailscale-Node-Caps` a comma-separated list of capabilities the device has (if any)
- `X-Tailscale-Node-Tags` a comma-separated list of tags the device has (if any)
- `X-Tailscale-Hostinfo` a JSON object containing some info about the device making the request

Any `X-Tailscale-*` headers sent by the client will be stripped before the request is proxied.
