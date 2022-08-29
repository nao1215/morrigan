# ping
ping subcommand send ICMP ECHO_REQUEST to network hosts.  
  
ping uses the ICMP protocol's mandatory ECHO_REQUEST datagram to elicit an ICMP ECHO_RESPONSE from a host or gateway. ECHO_REQUEST datagrams(“pings”) have an IP and ICMP header, followed by a struct timeval and then an arbitrary number of “pad” bytes used to fill out the packet.
  
Common ping commands are run with root privileges via seteuid; the ping provided by morrigan does not use seteuid, so only root users or sudoers can run it.
  
# Synopsis
```
morrigan ping DESTINATION [flags]
```
  
# Flags
```
  -a, --audible            audible rings a bell when a packet is received
  -c, --count string       iterations (default "5")
  -h, --help               help for ping
  -i, --interval string    interval in milliseconds (default "1000")
  -6, --ipv6               use ipv4 (means ip4:icmp) or 6 (ip6:ipv6-icmp)
  -s, --size string        packet data Size (default "64")
  -w, --wait-time string   wait time in milliseconds (default "100")
``` 
  
# Examples
```sh
$ sudo ./morrigan ping google.com
84 bytes from google.com: icmp_seq=1, time=12.99805ms
84 bytes from google.com: icmp_seq=2, time=12.685649ms
84 bytes from google.com: icmp_seq=3, time=12.832015ms
```
