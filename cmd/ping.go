package cmd

// This package is forked from https://github.com/u-root/u-root/blob/main/cmds/core/ping/ping.go
// (BSD 3-Clause "New" or "Revised" License)
//
// BSD 3-Clause License
//
// Copyright (c) 2012-2019, u-root Authors
// All rights reserved.
//
// Redistribution and use in source and binary forms, with or without
// modification, are permitted provided that the following conditions are met:
//
// * Redistributions of source code must retain the above copyright notice, this
//  list of conditions and the following disclaimer.
//
// * Redistributions in binary form must reproduce the above copyright notice,
//  this list of conditions and the following disclaimer in the documentation
//  and/or other materials provided with the distribution.
//
// * Neither the name of the copyright holder nor the names of its
//  contributors may be used to endorse or promote products derived from
//  this software without specific prior written permission.
//
// THIS SOFTWARE IS PROVIDED BY THE COPYRIGHT HOLDERS AND CONTRIBUTORS "AS IS"
// AND ANY EXPRESS OR IMPLIED WARRANTIES, INCLUDING, BUT NOT LIMITED TO, THE
// IMPLIED WARRANTIES OF MERCHANTABILITY AND FITNESS FOR A PARTICULAR PURPOSE ARE
// DISCLAIMED. IN NO EVENT SHALL THE COPYRIGHT HOLDER OR CONTRIBUTORS BE LIABLE
// FOR ANY DIRECT, INDIRECT, INCIDENTAL, SPECIAL, EXEMPLARY, OR CONSEQUENTIAL
// DAMAGES (INCLUDING, BUT NOT LIMITED TO, PROCUREMENT OF SUBSTITUTE GOODS OR
// SERVICES; LOSS OF USE, DATA, OR PROFITS; OR BUSINESS INTERRUPTION) HOWEVER
// CAUSED AND ON ANY THEORY OF LIABILITY, WHETHER IN CONTRACT, STRICT LIABILITY,
// OR TORT (INCLUDING NEGLIGENCE OR OTHERWISE) ARISING IN ANY WAY OUT OF THE USE
// OF THIS SOFTWARE, EVEN IF ADVISED OF THE POSSIBILITY OF SUCH DAMAGE.

import (
	"encoding/binary"
	"errors"
	"fmt"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/nao1215/morrigan/internal/print"
	"github.com/spf13/cobra"
)

var pingCmd = &cobra.Command{
	Use:   "ping DESTINATION",
	Short: "Send ICMP ECHO_REQUEST to network hosts",
	Long: `ping subcommand send ICMP ECHO_REQUEST to network hosts.

ping uses the ICMP protocol's mandatory ECHO_REQUEST datagram to elicit
an ICMP ECHO_RESPONSE from a host or gateway. ECHO_REQUEST datagrams
(“pings”) have an IP and ICMP header, followed by a struct timeval and
then an arbitrary number of “pad” bytes used to fill out the packet.
`,
	Example: `  morrigan ping google.com`,
	Run: func(cmd *cobra.Command, args []string) {
		if err := pingRun(cmd, args); err != nil {
			print.Err(err)
			os.Exit(1)
		}
	},
}

func init() {
	pingCmd.Flags().BoolP("ipv6", "6", false, "use ipv4 (means ip4:icmp) or 6 (ip6:ipv6-icmp)")
	pingCmd.Flags().StringP("size", "s", "64", "packet data Size")
	pingCmd.Flags().StringP("count", "c", "5", "iterations")
	pingCmd.Flags().StringP("interval", "i", "1000", "interval in milliseconds")
	pingCmd.Flags().StringP("wait-time", "w", "100", "wait time in milliseconds")
	pingCmd.Flags().BoolP("audible", "a", false, "audible rings a bell when a packet is received")
	rootCmd.AddCommand(pingCmd)
}

type pingFlag struct {
	net6           bool
	packetDataSize int
	iteration      uint64
	interval       int
	waitTime       int
	audible        bool
	host           string
}

const (
	icmpTypeEchoRequest           = 8
	icmpTypeEchoReply             = 0
	icmpEchoReplyHeaderIPv4Offset = 20
)

const (
	icmp6TypeEchoRequest           = 128
	icmp6TypeEchoReply             = 129
	icmp6EchoReplyHeaderIPv6Offset = 40
)

func pingRun(cmd *cobra.Command, args []string) error {
	flag, err := parsePingArgs(cmd, args)
	if err != nil {
		return err
	}

	if err := ping(flag); err != nil {
		return err
	}
	return nil
}

func parsePingArgs(cmd *cobra.Command, args []string) (*pingFlag, error) {
	flag := pingFlag{}

	if len(args) != 1 {
		return nil, errors.New("need DESTINATION. see help")
	}
	flag.host = args[0]

	net6, err := cmd.Flags().GetBool("ipv6")
	if err != nil {
		return nil, fmt.Errorf("%s: %w", "can not parse command line argument (--ipv6)", err)
	}
	flag.net6 = net6

	packetDataSize, err := cmd.Flags().GetString("size")
	if err != nil {
		return nil, fmt.Errorf("%s: %w", "can not parse command line argument (--size)", err)
	}
	size, err := strconv.Atoi(packetDataSize)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", "packet data size is integer value", err)
	}
	flag.packetDataSize = size

	iteration, err := cmd.Flags().GetString("count")
	if err != nil {
		return nil, fmt.Errorf("%s: %w", "can not parse command line argument (--count)", err)
	}
	size, err = strconv.Atoi(iteration)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", "iteration count is integer value", err)
	}
	flag.iteration = uint64(size)

	interval, err := cmd.Flags().GetString("interval")
	if err != nil {
		return nil, fmt.Errorf("%s: %w", "can not parse command line argument (--interval)", err)
	}
	size, err = strconv.Atoi(interval)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", "interval is integer value", err)
	}
	flag.interval = size

	waitTime, err := cmd.Flags().GetString("wait-time")
	if err != nil {
		return nil, fmt.Errorf("%s: %w", "can not parse command line argument (--wait-time)", err)
	}
	size, err = strconv.Atoi(waitTime)
	if err != nil {
		return nil, fmt.Errorf("%s: %w", "wait time is integer value", err)
	}
	flag.waitTime = size

	audible, err := cmd.Flags().GetBool("audible")
	if err != nil {
		return nil, fmt.Errorf("%s: %w", "can not parse command line argument (--audible)", err)
	}
	flag.audible = audible

	return &flag, nil
}

// Ping interface
type Ping struct {
	dial func(string, string) (net.Conn, error)
}

// New return Ping interface
func New() *Ping {
	return &Ping{
		dial: net.Dial,
	}
}

func cksum(bs []byte) uint16 {
	sum := uint32(0)

	for k := 0; k < len(bs)/2; k++ {
		sum += uint32(bs[k*2]) << 8
		sum += uint32(bs[k*2+1])
	}
	if len(bs)%2 != 0 {
		sum += uint32(bs[len(bs)-1]) << 8
	}
	sum = (sum >> 16) + (sum & 0xffff)
	sum = (sum >> 16) + (sum & 0xffff)
	if sum == 0xffff {
		sum = 0
	}

	return ^uint16(sum)
}

func (p *Ping) ping1(net6 bool, host string, i uint64, waitFor time.Duration, packetDataSize int) (string, error) {
	netname := "ip4:icmp"
	// todo: just figure out if it's an ip6 address and go from there.
	if net6 {
		netname = "ip6:ipv6-icmp"
	}
	c, err := p.dial(netname, host)
	if err != nil {
		return "", fmt.Errorf("net.Dial(%v %v) failed: %v", netname, host, err)
	}
	defer c.Close()

	if net6 {
		ipc := c.(*net.IPConn)
		if err := setupICMPv6Socket(ipc); err != nil {
			return "", fmt.Errorf("failed to set up the ICMPv6 connection: %w", err)
		}
	}

	// Send ICMP Echo Request
	if err := c.SetDeadline(time.Now().Add(waitFor)); err != nil {
		return "", err
	}

	msg := make([]byte, packetDataSize)
	if net6 {
		msg[0] = icmp6TypeEchoRequest
	} else {
		msg[0] = icmpTypeEchoRequest
	}
	msg[1] = 0
	binary.BigEndian.PutUint16(msg[6:], uint16(i))
	binary.BigEndian.PutUint16(msg[4:], uint16(i>>16))
	binary.BigEndian.PutUint16(msg[2:], cksum(msg))
	if _, err := c.Write(msg[:]); err != nil {
		return "", fmt.Errorf("write failed: %v", err)
	}

	// Get ICMP Echo Reply
	if err := c.SetDeadline(time.Now().Add(waitFor)); err != nil {
		return "", err
	}
	rmsg := make([]byte, packetDataSize+256)
	before := time.Now()
	amt, err := c.Read(rmsg[:])
	if err != nil {
		return "", fmt.Errorf("read failed: %v", err)
	}
	latency := time.Since(before)
	if !net6 {
		rmsg = rmsg[icmpEchoReplyHeaderIPv4Offset:]
	}
	if net6 {
		if rmsg[0] != icmp6TypeEchoReply {
			return "", fmt.Errorf("bad ICMPv6 echo reply type, got %d, want %d", rmsg[0], icmp6TypeEchoReply)
		}
	} else {
		if rmsg[0] != icmpTypeEchoReply {
			return "", fmt.Errorf("bad ICMP echo reply type, got %d, want %d", rmsg[0], icmpTypeEchoReply)
		}
	}
	cks := binary.BigEndian.Uint16(rmsg[2:])
	binary.BigEndian.PutUint16(rmsg[2:], 0)
	// only validate the checksum for IPv4. For IPv6 this *should* be done by the
	// TCP stack (and do we need to validate the checksum anyway?)
	if !net6 && cks != cksum(rmsg) {
		return "", fmt.Errorf("bad ICMP checksum: %v (expected %v)", cks, cksum(rmsg))
	}
	id := binary.BigEndian.Uint16(rmsg[4:])
	seq := binary.BigEndian.Uint16(rmsg[6:])
	rseq := uint64(id)<<16 + uint64(seq)
	if rseq != i {
		return "", fmt.Errorf("wrong sequence number %v (expected %v)", rseq, i)
	}

	return fmt.Sprintf("%d bytes from %v: icmp_seq=%v, time=%v", amt, host, i, latency), nil
}

func ping(flag *pingFlag) error {
	if flag.packetDataSize < 8 {
		return fmt.Errorf("packet size too small (must be >= 8): %v", flag.packetDataSize)
	}

	interval := time.Duration(flag.interval)
	p := New()

	// ping needs to run forever, except if '*iter' is not zero
	waitFor := time.Duration(flag.waitTime) * time.Millisecond
	for i := uint64(0); i <= flag.iteration; i++ {
		msg, err := p.ping1(flag.net6, flag.host, i+1, waitFor, flag.packetDataSize)
		if err != nil {
			return fmt.Errorf("ping failed: %v", err)
		}
		if flag.audible {
			msg = "\a" + msg
		}
		fmt.Println(msg)
		time.Sleep(time.Millisecond * interval)
	}
	return nil
}
