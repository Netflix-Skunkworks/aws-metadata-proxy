# AWS Metadata Proxy

Example AWS Metadata proxy to protect against attack vectors targetting AWS Credentials 

## Getting Started

Clone the repo

```
git clone https://github.com/Netflix-Skunkworks/aws-metadata-proxy.git
cd aws-metadata-proxy
```

Build the proxy

```golang
go get
go build
```

## Network Setup

Create an `iptable` rule that prevents talking directly to the AWS Metadata Service **except** for a particular user, `proxy_user` in the example below.  This is the user you run the proxy as on your server.

```
/sbin/iptables -t nat -A OUTPUT -m owner ! --uid-owner proxy_user -d 169.254.169.254 -p tcp -m tcp --dport 80 -j DNAT --to-destination 127.0.0.1:9090
```

