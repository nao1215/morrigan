package cmd

import (
	"reflect"
	"testing"
)

const TestRipeWhoisResponse = `% IANA WHOIS server
% for more information on IANA, visit http://www.iana.org
% This query returned 1 object
refer:        whois.ripe.net
inetnum:      93.0.0.0 - 93.255.255.255
organization: RIPE NCC
status:       ALLOCATED
whois:        whois.ripe.net
changed:      2007-03
source:       IANA
% This is the RIPE Database query service.
% The objects are in RPSL format.
%
% The RIPE Database is subject to Terms and Conditions.
% See http://www.ripe.net/db/support/db-terms-conditions.pdf
% Note: this output has been filtered.
%       To receive output for a database update, use the "-B" flag.
% Information related to '93.184.216.0 - 93.184.216.255'
% Abuse contact for '93.184.216.0 - 93.184.216.255' is 'abuse@verizondigitalmedia.com'
inetnum:        93.184.216.0 - 93.184.216.255
netname:        EDGECAST-NETBLK-03
descr:          NETBLK-03-EU-93-184-216-0-24
country:        EU
admin-c:        DS7892-RIPE
tech-c:         DS7892-RIPE
status:         ASSIGNED PA
mnt-by:         MNT-EDGECAST
created:        2012-06-22T21:48:41Z
last-modified:  2012-06-22T21:48:41Z
source:         RIPE # Filtered
person:         Derrick Sawyer
address:        13031 W Jefferson Blvd #900, Los Angeles, CA 90094
phone:          +18773343236
nic-hdl:        DS7892-RIPE
created:        2010-08-25T18:44:19Z
last-modified:  2017-03-03T09:06:18Z
source:         RIPE
mnt-by:         MNT-EDGECAST
% This query was served by the RIPE Database Query Service version 1.102.2 (WAGYU)
;; Query time: 1194 msec
;; WHEN: Tue Feb 01 21:12:31 JST 2022`

func TestResolve(t *testing.T) {
	type fields struct {
		parseResult []networkAdmin
	}
	type args struct {
		domain  string
		verbose bool
	}
	tests := []struct {
		name string
		args args
		want fields
	}{
		{
			name: "the_result_of_example.com_must_be_returned",
			args: args{
				domain:  "example.com",
				verbose: false,
			},
			want: fields{
				parseResult: []networkAdmin{
					{
						ipRange: "93.0.0.0 - 93.255.255.255",
					},
					{
						ipRange: "93.184.216.0 - 93.184.216.255",
						admin:   "NETBLK-03-EU-93-184-216-0-24",
						country: "EU",
						netName: "EDGECAST-NETBLK-03",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := resolve(tt.args.domain, tt.args.verbose)
			if err != nil {
				t.Errorf("Resolve() error: %v", err)
			}
			if !reflect.DeepEqual(got.parseResult, tt.want.parseResult) {
				t.Fatalf("domainSummary.parseWhoisResponse() = %v, want %v", got.parseResult, tt.want.parseResult)
			}
		})
	}
}

func TestSummary_ParseWhoisResponse(t *testing.T) {
	type fields struct {
		whoisResponseServer string
		whoisResponse       string
		parseResult         []networkAdmin
	}
	tests := []struct {
		name   string
		fields fields
		want   fields
	}{
		{
			name: "be_able_to_correctly_parse_the_response_from_ripe",
			fields: fields{
				whoisResponseServer: "whois.ripe.net",
				whoisResponse:       TestRipeWhoisResponse,
			},
			want: fields{
				whoisResponseServer: "whois.ripe.net",
				whoisResponse:       TestRipeWhoisResponse,
				parseResult: []networkAdmin{
					{
						ipRange: "93.0.0.0 - 93.255.255.255",
						admin:   "RIPE NCC",
					},
					{
						ipRange: "93.184.216.0 - 93.184.216.255",
						admin:   "NETBLK-03-EU-93-184-216-0-24",
						country: "EU",
						netName: "EDGECAST-NETBLK-03",
					},
				},
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &domainSummary{
				whoisResponseServer: tt.fields.whoisResponseServer,
				whoisResponse:       tt.fields.whoisResponse,
				parseResult:         tt.fields.parseResult,
			}
			s.ParseWhoisResponse()
			if !reflect.DeepEqual(s.parseResult, tt.want.parseResult) {
				t.Fatalf("domainSummary.parseWhoisResponse() = %v, want %v", s.parseResult, tt.want.parseResult)
			}
		})
	}
}

func TestSummary_SetWhoisResponseServerFromWhoisResponse(t *testing.T) {
	type fields struct {
		whoisResponseServer string
		whoisResponse       string
	}
	tests := []struct {
		name   string
		fields fields
		want   fields
	}{
		{
			name: "the_response_server_can_be_set_correctly_from_the_response_of_ripe",
			fields: fields{
				whoisResponseServer: "whois.apnic.net",
				whoisResponse:       TestRipeWhoisResponse,
			},
			want: fields{
				whoisResponseServer: "whois.ripe.net",
				whoisResponse:       TestRipeWhoisResponse,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &domainSummary{
				whoisResponseServer: tt.fields.whoisResponseServer,
				whoisResponse:       tt.fields.whoisResponse,
			}
			s.setWhoisResponseServerFromWhoisResponse()
			if s.whoisResponseServer != tt.want.whoisResponseServer {
				t.Errorf("domainSummary.setWhoisResponseServerFromWhoisResponse() = %v, want %v",
					s.whoisResponseServer, tt.want.whoisResponseServer)
			}
		})
	}
}

func TestSummary_ParseCheck(t *testing.T) {
	type fields struct {
		targetDomain        string
		targetIP            string
		whoisResponseServer string
		whoisResponse       string
		parseResult         []networkAdmin
	}
	type wants struct {
		whoisResponseServer string
		result              bool
	}
	tests := []struct {
		name   string
		fields fields
		want   wants
	}{
		{
			name: "return_true_in_non-redirected_response",
			fields: fields{
				whoisResponseServer: "whois.ripe.net",
				parseResult: []networkAdmin{
					{
						ipRange: "93.0.0.0 - 93.255.255.255",
						admin:   "RIPE NCC",
					},
					{
						ipRange: "93.184.216.0 - 93.184.216.255",
						admin:   "NETBLK-03-EU-93-184-216-0-24",
						country: "EU",
						netName: "EDGECAST-NETBLK-03",
					},
				},
			},
			want: wants{
				whoisResponseServer: "whois.ripe.net",
				result:              true,
			},
		},
		{
			name: "return_false_in_redirected_response",
			fields: fields{
				whoisResponseServer: "whois.apnic.net",
				parseResult: []networkAdmin{
					{
						ipRange: "157.1.0.0 - 157.14.255.255",
						admin:   "APNIC",
						country: "",
						netName: "",
					},
					{
						ipRange: "93.0.0.0 - 93.255.255.255",
						admin:   "DS7892-RIPE",
						country: "",
						netName: "RIPE NCC",
					},
				},
			},
			want: wants{
				whoisResponseServer: "whois.ripe.net",
				result:              false,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			s := &domainSummary{
				targetDomain:        tt.fields.targetDomain,
				targetIP:            tt.fields.targetIP,
				whoisResponseServer: tt.fields.whoisResponseServer,
				whoisResponse:       tt.fields.whoisResponse,
				parseResult:         tt.fields.parseResult,
			}
			got := s.parseCheck()
			if got != tt.want.result {
				t.Errorf("domainSummary.parseCheck() Return = %v, want %v", got, tt.want.result)
			}
			if s.whoisResponseServer != tt.want.whoisResponseServer {
				t.Errorf("Summary.ParseCheck() WhoisResponseServer = %v, want %v", s.whoisResponseServer, tt.want.whoisResponseServer)
			}
		})
	}
}
