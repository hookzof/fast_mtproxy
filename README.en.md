[Russian version](README.md)

# fast_mtproxy
Fast deployment of the official MTProxy server with Fake TLS support.  
  
It's also worth paying attention to third-party implementations of MTProxy: [Erlang](https://github.com/seriyps/mtproto_proxy), [Golang](https://github.com/9seconds/mtg), [Python](https://github.com/alexbers/mtprotoproxy)

## Supported OS images
* Ubuntu 22.04/24.04
* Debian 12
* CentOS Stream 9/10

*On other platforms, the work of the script is not guaranteed!*

## Working with the script
### Downloading an executable file
```bash
curl -L -o install https://git.io/JeOSr && chmod +x install
```

### Quick installation
```bash
./install
```

Default: port 443, generates a random secret key, domain www.google.com

### Example of detailed installation
```bash
./install -H=443 -S=25c8dfee81acdadaff3a338a10db8497 -D=www.google.com
```
`-p` - local port for statistics (optional);  
`-H` - port to connect;  
`-S` - secret key;  
`-P` - ad tag (get here @MTProxybot);  
`-D` - TLS 1.3 domain (you can check <a href="https://www.cdn77.com/tls-test">here</a>);  
`-6` - enable ipv6 protocol (must be supported by your hosting provider).

### Actions with server
```bash
./install -restart=443
```

Flags indicating the server port:
* `start` - server start;
* `stop` - server stop;
* `restart` - server restart;
* `enable` - server activation + autorun;
* `disable` - server shutdown (soft deletion);
* `delete` - server removal.

## Recommendations
Due to its peculiarities, the protocol mimics under TLS 1.3, so it is necessary to use the corresponding ports for more convincing (examples of domains are given in brackets):
* 261, 271, 324
* 443 (www.google.com, www.youtube.com)
* 465, 563, 636
* 853 (dns.google, cloudflare-dns.com)
* 989, 990, 992, 993, 994, 995, 4843, 5061, 5085, 5349, 5671, 6513, 6514, 6619, 8883
