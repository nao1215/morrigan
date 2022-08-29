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
	"fmt"
	"net"
	"testing"
	"time"
)

type myConn struct {
	net.IPConn
	testRun string
}

func (M *myConn) Read(b []byte) (int, error) {
	if M.testRun == "error in read" {
		return 0, fmt.Errorf("err")
	}
	if M.testRun == "rmsg[0] != ICMP_TYPE_ECHO_REPLY" {
		b[20] = 0xff
	}
	if M.testRun == "!net6 && cks != cksum(rmsg)" {
		return 0, nil
	}
	b[22] = 0xff
	b[23] = 0xff
	return 0, nil
}

func (M *myConn) Write(b []byte) (int, error) {
	if M.testRun == "error in write" {
		return 0, fmt.Errorf("err")
	}
	return 0, nil
}

func (M *myConn) SetDeadline(t time.Time) error {
	return nil
}

// Test cksum
func TestCkSum(t *testing.T) {
	for _, tt := range []struct {
		name  string
		input []byte
		want  uint16
	}{
		{
			name:  "ultimate test, triggers the sum == 0xffff and len(bs)%2 !=0",
			input: []byte{0xff, 0xff, 0x00},
			want:  65535,
		},
		{
			name:  "another input",
			input: []byte{0xfe, 0xfe, 0xfe, 0xfe},
			want:  514,
		},
		{
			name:  "empty input",
			input: []byte{},
			want:  65535,
		},
	} {
		if got := cksum(tt.input); got != tt.want {
			t.Errorf("cksum() = '%d', want: '%d'", got, tt.want)
		}
	}

}

// Test Ping1
func TestPing1(t *testing.T) {
	for _, tt := range []struct {
		name           string
		p              Ping
		net6           bool
		host           string
		i              uint64
		waitFor        time.Duration
		packetDataSize int
		want           error
	}{
		{
			name: "ping1 without error",
			p: Ping{
				dial: func(s1, s2 string) (net.Conn, error) {
					return &myConn{testRun: "ping1 without error"}, nil
				},
			},
			net6:           false,
			host:           "test.com",
			i:              0,
			waitFor:        time.Duration(1000) * time.Millisecond,
			packetDataSize: 20,
			want:           fmt.Errorf(""),
		},
		{
			name: "error in dial",
			p: Ping{
				dial: func(s1, s2 string) (net.Conn, error) {
					return nil, fmt.Errorf("err")
				},
			},
			net6:           false,
			host:           "test.com",
			i:              1,
			waitFor:        time.Duration(1000) * time.Millisecond,
			packetDataSize: 20,
			want:           fmt.Errorf("net.Dial(%v %v) failed: %v", "ip4:icmp", "test.com", fmt.Errorf("err")),
		},
		{
			name: "error in write",
			p: Ping{
				dial: func(s1, s2 string) (net.Conn, error) {
					return &myConn{testRun: "error in write"}, nil
				},
			},
			net6:           false,
			host:           "test.com",
			i:              0,
			waitFor:        time.Duration(1000) * time.Millisecond,
			packetDataSize: 20,
			want:           fmt.Errorf("write failed: err"),
		},
		{
			name: "error in read",
			p: Ping{
				dial: func(s1, s2 string) (net.Conn, error) {
					return &myConn{testRun: "error in read"}, nil
				},
			},
			net6:           false,
			host:           "test.com",
			i:              0,
			waitFor:        time.Duration(1000) * time.Millisecond,
			packetDataSize: 20,
			want:           fmt.Errorf("read failed: err"),
		},
		{
			name: "rmsg[0] != ICMP_TYPE_ECHO_REPLY",
			p: Ping{
				dial: func(s1, s2 string) (net.Conn, error) {
					return &myConn{testRun: "rmsg[0] != ICMP_TYPE_ECHO_REPLY"}, nil
				},
			},
			net6:           false,
			host:           "test.com",
			i:              0,
			waitFor:        time.Duration(1000) * time.Millisecond,
			packetDataSize: 20,
			want:           fmt.Errorf("bad ICMP echo reply type, got %d, want %d", 0xff, icmpTypeEchoReply),
		},
		{
			name: "!net6 && cks != cksum(rmsg)",
			p: Ping{
				dial: func(s1, s2 string) (net.Conn, error) {
					return &myConn{testRun: "!net6 && cks != cksum(rmsg)"}, nil
				},
			},
			net6:           false,
			host:           "test.com",
			i:              0,
			waitFor:        time.Duration(1000) * time.Millisecond,
			packetDataSize: 20,
			want:           fmt.Errorf("bad ICMP checksum: %v (expected %v)", 0, 65535),
		},
		{
			name: "rseq != i",
			p: Ping{
				dial: func(s1, s2 string) (net.Conn, error) {
					return &myConn{testRun: "rseq != i"}, nil
				},
			},
			net6:           false,
			host:           "test.com",
			i:              1,
			waitFor:        time.Duration(1000) * time.Millisecond,
			packetDataSize: 20,
			want:           fmt.Errorf("wrong sequence number %v (expected %v)", 0, 1),
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			if _, got := tt.p.ping1(tt.net6, tt.host, tt.i, tt.waitFor, tt.packetDataSize); got != nil {
				if got.Error() != tt.want.Error() {
					t.Errorf("ping1() = '%s', want: '%s'", got, tt.want)
				}
			}
		})
	}
}

// Test refactored ping()
func TestPing(t *testing.T) {
	for _, tt := range []struct {
		name       string
		packetSize int
		audible    bool
		host       string
		waitFor    time.Duration
		want       error
	}{
		{
			name:       "packetSize < 8",
			packetSize: 7,
			host:       "test.com",
			audible:    true,
			waitFor:    time.Minute,
			want:       fmt.Errorf("packet size too small (must be >= 8): %v", 7),
		},
		{
			name:       "ping with error",
			packetSize: 8,
			host:       "",
			audible:    true,
			waitFor:    time.Minute,
			want:       fmt.Errorf("ping failed: net.Dial(ip4:icmp ) failed: dial ip4:icmp: missing address"),
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			flag := &pingFlag{
				packetDataSize: tt.packetSize,
				audible:        tt.audible,
			}
			if got := ping(flag); got != nil {
				if got.Error() != tt.want.Error() {
					t.Errorf("ping() = '%s', want: '%s'", got, tt.want)
				}
			}
		})
	}
}

// This test gets the coverage higher and does not test any functionality.
func TestNew(t *testing.T) {
	_ = New()
}
