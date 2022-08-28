# netcat
This netcat subcommnad is forked from [go-netcat](https://github.com/vfedoroff/go-netcat).
  
netcat subcommand is simple implementation of the netcat utility in golang that allows to listen and send data over TCP and UDP protocols. netcat can open TCP connections, send UDP packets, listen on arbitrary TCP and UDP ports, perform port scanning.
  
# Synopsis
```
morrigan netcat [hostname] [port] [flags]
```
  
# Flags
```
  -h, --help                       help for netcat
  -l, --listen                     listen for an incoming connection rather than initiate a connection to a remote host
  -p, --port string                specifies the source port
  -s, --source-ip-address string   specify the IP of the interface which is used to send the packets.
  -u, --udp                        use UDP instead of the default option of TCP
``` 
  
# Examples
Open a TCP connection to port 42 of hostname.
```sh
$ morrigan netcat hostname 42
```

Open a UDP connection to port 53 of hostname.
```sh
$ morrigan netcat -u hostname 53
```

Listen on TCP port 3000, and once there is a connection, send stdin to the remote host, and send data from the remote host to stdout.
```sh
$ morrigan netcat -l 3000
```

Listen on UDP port 3000, and once there is a connection, send stdin to the remote host, and send data from the remote host to stdout.
```sh
$ morrigan netcat -u -l 3000
```