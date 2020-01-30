# AWS Metadata Proxy

Example AWS Metadata proxy to protect against attack vectors targetting AWS Credentials 

## Preface
AWS updated their metadata service in Nov '20 ([IMDSv2](https://aws.amazon.com/about-aws/whats-new/2019/11/announcing-updates-amazon-ec2-instance-metadata-service/)). With IMDSv2 a token will be sent when requesting the metadata service, therefore the risk of an attacker requesting the metadata service through ssrf has been minimized. There is a few blogposts about implementation details ([getting started with imdsv2](https://blog.appsecco.com/getting-started-with-version-2-of-aws-ec2-instance-metadata-service-imdsv2-2ad03a1f3650)). 

Depending on your usecase and the threat model it could be reasonable to just use the IMDSV2. 


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

