// Package cmd manages the entry points for the subcommands that morrigan has.
package cmd

// netcat subcommand is forked from https://github.com/vfedoroff/go-netcat
//
// The MIT License (MIT)
//
// Copyright (c) 2015 Vadym Fedorov (original author)
//
// Permission is hereby granted, free of charge, to any person obtaining a copy
// of this software and associated documentation files (the "Software"), to deal
// in the Software without restriction, including without limitation the rights
// to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
// copies of the Software, and to permit persons to whom the Software is
// furnished to do so, subject to the following conditions:
//
// The above copyright notice and this permission notice shall be included in all
// copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
// IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
// FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
// AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
// LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
// OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
// SOFTWARE.

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"strconv"

	"github.com/nao1215/morrigan/internal/print"
	"github.com/spf13/cobra"
)

var netcatCmd = &cobra.Command{
	Use:   "netcat [hostname] [port]",
	Short: "Arbitrary TCP and UDP connections and listens",
	Long: `netcat subcommand allows to listen TCP/UDP ports and send data to remote ports over TCP/UDP.
`,
	Example: `  morrigan netcat hostname 12`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := netcat(cmd, args); err != nil {
			print.Err(err)
			os.Exit(1)
		}
	},
}

func init() {
	netcatCmd.Flags().BoolP("udp", "u", false, "use UDP instead of the default option of TCP")
	netcatCmd.Flags().BoolP("listen", "l", false, "listen for an incoming connection rather than initiate a connection to a remote host")
	netcatCmd.Flags().StringP("port", "p", "", "specifies the source port")
	netcatCmd.Flags().StringP("source-ip-address", "s", "", "specify the IP of the interface which is used to send the packets.")
	rootCmd.AddCommand(netcatCmd)
}

type netcatFlag struct {
	isUDP    bool
	isListen bool
	srcPort  string
	srcIP    string
	hostname string
	destPort string
}

func netcat(cmd *cobra.Command, args []string) error {
	flag, err := parseNetCatArgs(cmd, args)
	if err != nil {
		return err
	}

	if flag.isUDP {
		if err := udpHandler(flag); err != nil {
			return err
		}
	} else {
		if err := tcpHandler(flag); err != nil {
			return err
		}
	}
	return nil
}

func tcpHandler(flag *netcatFlag) error {
	if flag.isListen {
		listener, err := net.Listen("tcp", ":"+flag.destPort)
		if err != nil {
			return err
		}
		con, err := listener.Accept()
		if err != nil {
			return err
		}
		print.Info("connect from " + con.RemoteAddr().String())
		tcpConnection(con)
	} else if flag.hostname != "" {
		con, err := net.Dial("tcp", flag.hostname+":"+flag.destPort)
		if err != nil {
			return err
		}
		print.Info("connect to " + flag.hostname + ":" + flag.destPort)
		tcpConnection(con)
	} else {
		return errors.New("see help: $ morrigan netcat --help")
	}
	return nil
}

func udpHandler(flag *netcatFlag) error {
	if flag.isListen {
		addr, err := net.ResolveUDPAddr("udp", ":"+flag.destPort)
		if err != nil {
			return err
		}
		con, err := net.ListenUDP("udp", addr)
		if err != nil {
			return err
		}
		print.Info("has been resolved UDP address: " + addr.String())
		print.Info("Listening on " + flag.destPort)
		udpConnection(con)
	} else if flag.hostname != "" {
		addr, err := net.ResolveUDPAddr("udp", flag.hostname+":"+flag.destPort)
		if err != nil {
			log.Fatalln(err)
		}
		log.Println("Has been resolved UDP address:", addr)
		con, err := net.DialUDP("udp", nil, addr)
		if err != nil {
			log.Fatalln(err)
		}
		udpConnection(con)
	}
	return nil
}

func parseNetCatArgs(cmd *cobra.Command, args []string) (*netcatFlag, error) {
	flag := netcatFlag{}

	isUDP, err := cmd.Flags().GetBool("udp")
	if err != nil {
		return nil, fmt.Errorf("%s: %w", "can not parse command line argument (--udp)", err)
	}
	flag.isUDP = isUDP

	isListen, err := cmd.Flags().GetBool("listen")
	if err != nil {
		return nil, fmt.Errorf("%s: %w", "can not parse command line argument (--listen)", err)
	}
	flag.isListen = isListen

	srcPort, err := cmd.Flags().GetString("port")
	if err != nil {
		return nil, fmt.Errorf("%s: %w", "can not parse command line argument (--port)", err)
	}
	flag.srcPort = srcPort

	srcIP, err := cmd.Flags().GetString("source-ip-address")
	if err != nil {
		return nil, fmt.Errorf("%s: %w", "can not parse command line argument (--source-ip-address)", err)
	}
	flag.srcIP = srcIP

	if len(args) == 0 && !flag.isListen && !flag.isUDP && flag.srcIP == "" && flag.srcPort == "" {
		return nil, errors.New("see help: $ morrigan netcat --help")
	}

	if flag.srcPort != "" {
		if _, err := strconv.Atoi(flag.srcPort); err != nil {
			return nil, errors.New("source port is integer value. you set " + flag.srcPort)
		}
	}

	if !flag.isListen {
		if len(args) < 2 {
			return nil, errors.New("if you don't use listen mode, you specify [hostname] and [port]")
		}

		if _, err := strconv.Atoi(args[1]); err != nil {
			return nil, errors.New("source port is integer value. you set " + args[1])

		}
		flag.hostname = args[0]
		flag.destPort = args[1]
	} else {
		if len(args) < 1 {
			return nil, errors.New("if you use listen mode, you specify [port]")
		}
		if _, err := strconv.Atoi(args[0]); err != nil {
			return nil, errors.New("destination port is integer value. you set " + args[0])
		}
		flag.destPort = args[0]
	}
	return &flag, nil
}

// tcpConnection handles TCP connection and perform synchorinization:
// TCP -> Stdout and Stdin -> TCP
func tcpConnection(con net.Conn) {
	chanToStdout := streamCopy(con, os.Stdout)
	chanToRemote := streamCopy(os.Stdin, con)
	select {
	case <-chanToStdout:
		print.Info("remote connection is closed")
	case <-chanToRemote:
		print.Info("local program is terminated")
	}
}

// streamCopy performs copy operation between streams:
// os and tcp streams
func streamCopy(src io.Reader, dst io.Writer) <-chan int {
	buf := make([]byte, 1024)
	syncChannel := make(chan int)
	go func() {
		defer func() {
			if con, ok := dst.(net.Conn); ok {
				con.Close()
				print.Info("connection fron " + con.RemoteAddr().String() + " is closed")
			}
			syncChannel <- 0 // Notify that processing is finished
		}()
		for {
			var nBytes int
			var err error
			nBytes, err = src.Read(buf)
			if err != nil {
				if err != io.EOF {
					print.Err("read error: " + err.Error())
				}
				break
			}
			_, err = dst.Write(buf[0:nBytes])
			if err != nil {
				print.Fatal("write error: " + err.Error())
			}
		}
	}()
	return syncChannel
}

// udpConnection handle UDP connection
func udpConnection(con net.Conn) {
	inChannel := acceptFromUDPToStream(con, os.Stdout)
	print.Info("waiting for remote connection")

	remoteAddr := <-inChannel
	print.Info("connected from " + remoteAddr.String())
	outChannel := putFromStreamToUDP(os.Stdin, con, remoteAddr)

	select {
	case <-inChannel:
		print.Info("remote connection is closed")
	case <-outChannel:
		print.Info("local program is terminated")
	}
}

// acceptFromUDPToStream accept data from UPD connection and copy it to the stream
func acceptFromUDPToStream(src net.Conn, dst io.Writer) <-chan net.Addr {
	buf := make([]byte, 1024)
	syncChannel := make(chan net.Addr)
	con, ok := src.(*net.UDPConn)
	if !ok {
		print.Info("input must be UDP connection")
		return syncChannel
	}
	go func() {
		var remoteAddr net.Addr
		for {
			var nBytes int
			var err error
			var addr net.Addr
			nBytes, addr, err = con.ReadFromUDP(buf)
			if err != nil {
				if err != io.EOF {
					print.Err("read error: " + err.Error())
				}
				break
			}

			if remoteAddr == nil && remoteAddr != addr {
				remoteAddr = addr
				syncChannel <- remoteAddr
			}
			_, err = dst.Write(buf[0:nBytes])
			if err != nil {
				print.Fatal("write error: " + err.Error())
			}
		}
	}()
	return syncChannel
}

// Put input date from the stream to UDP connection
func putFromStreamToUDP(src io.Reader, dst net.Conn, remoteAddr net.Addr) <-chan net.Addr {
	buf := make([]byte, 1024)
	syncChannel := make(chan net.Addr)
	go func() {
		for {
			var nBytes int
			var err error
			nBytes, err = src.Read(buf)
			if err != nil {
				if err != io.EOF {
					print.Err("read error: " + err.Error())
				}
				break
			}

			print.Info("write to the remote address: " + remoteAddr.String())
			if con, ok := dst.(*net.UDPConn); ok && remoteAddr != nil {
				_, err = con.WriteTo(buf[0:nBytes], remoteAddr)
			}
			if err != nil {
				print.Fatal("write error: " + err.Error())
			}
		}
	}()
	return syncChannel
}
