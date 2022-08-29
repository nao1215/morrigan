# whris
whris subcommnad is forked from [whris](https://github.com/harakeishi/whris). whris is displays management information for IPs associated with the domain.

# Synopsis
```
morrigan whris DOMAIN_NAMEs [flags]
```

# Examples
If you want to know the administrator of the IP associated with the domain.
```sh
$ morrigan whris example.com 
[Basic Information]
 Target domain: example.com
 Target ip    : 93.184.216.34
 Network Admin: NETBLK-03-EU-93-184-216-0-24
 Network name : EDGECAST-NETBLK-03
 ip range     : 93.184.216.0 - 93.184.216.255
 country      : EU
```

If you want to see more details, use the `--verbose` option.
```sh
$ morrigan whris example.com --verbose
[Basic Information]
 Target domain: example.com
 Target ip    : 93.184.216.34
 Network Admin: NETBLK-03-EU-93-184-216-0-24
 Network name : EDGECAST-NETBLK-03
 ip range     : 93.184.216.0 - 93.184.216.255
 country      : EU
[Detailed Information]
 0:
 Network Admin: RIPE NCC
 Network name : 
 ip range     : 93.0.0.0 - 93.255.255.255
 country      : 
 1:
 Network Admin: NETBLK-03-EU-93-184-216-0-24
 Network name : EDGECAST-NETBLK-03
 ip range     : 93.184.216.0 - 93.184.216.255
 country      : EU
```