# smart-dns-proxy

A lightweight DNS proxy that transparently forwards queries upstream and, when a response contains a known-blocked IP, replaces it with a clean alternative from the same provider.

## How it works

1. A DNS query arrives at the proxy.
2. The query is forwarded to an upstream resolver (default: `8.8.8.8:53`).
3. The response IPs are checked against a periodically-refreshed blocklist fetched from a remote data source.
4. If a response IP is flagged as **blocked**, the proxy finds a **clean IP** from the same CDN/provider and substitutes it in the response before sending it back to the client.
5. If no IP is blocked, the original response is returned unchanged.

This is useful when a content provider (e.g. a CDN like Cloudflare) serves the same domain from multiple IPs, but some of those IPs are unreachable due to IP-based blocking — while other IPs from the same provider remain accessible.

## Usage

```
./smart-dns-proxy [flags]
```

| Flag | Default | Description |
|------|---------|-------------|
| `-host` | `127.0.0.1` | Address to listen on |
| `-port` | `53` | Port to listen on |
| `-resolver` | `8.8.8.8:53` | Upstream DNS resolver |

### Example

```bash
# Listen on all interfaces, use Cloudflare as upstream resolver
./smart-dns-proxy -host 0.0.0.0 -port 53 -resolver 1.1.1.1:53
```

## Docker

```bash
docker build -t smart-dns-proxy .

# Run with defaults (listens on 0.0.0.0:53, resolves via 8.8.8.8)
docker run -p 53:53/udp smart-dns-proxy

# Custom resolver
docker run -p 53:53/udp smart-dns-proxy -host 0.0.0.0 -port 53 -resolver 1.1.1.1:53
```

## Data source

The blocklist is fetched from a remote JSON endpoint at startup and refreshed every minute. The data includes IP addresses grouped by provider, each with a block state derived from its most recent state change. An IP is considered blocked if its latest recorded state marks it as blocked.

## Acknowledgements

The blocklist data powering this proxy is provided by [**¿Hay ahora fútbol?**](https://hayahora.futbol), an independent transparency initiative that monitors and publishes the IP addresses blocked by Spanish ISPs during football matches. Their work documents the collateral damage these blocks cause to unrelated websites that share CDN infrastructure — bringing much-needed visibility to a process that otherwise has no public registry.

Many thanks to them for making this data openly available.

## Legal disclaimer

This tool is provided for **legitimate and lawful use only**. It is intended to work around IP-based network restrictions in contexts where such access is legally permitted — for example, accessing content that is geographically or technically restricted due to infrastructure routing issues rather than legal prohibitions.

**The authors are not responsible for any misuse, illegal activity, or violations of third-party terms of service arising from the use of this software.** Users are solely responsible for ensuring their usage complies with all applicable laws and regulations in their jurisdiction.
