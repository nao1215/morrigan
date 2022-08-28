[![Build](https://github.com/nao1215/morrigan/actions/workflows/build.yml/badge.svg)](https://github.com/nao1215/morrigan/actions/workflows/build.yml)
[![UnitTest](https://github.com/nao1215/morrigan/actions/workflows/unit_test.yml/badge.svg)](https://github.com/nao1215/morrigan/actions/workflows/unit_test.yml)
[![codecov](https://codecov.io/gh/nao1215/morrigan/branch/main/graph/badge.svg?token=AGqQgVDcL1)](https://codecov.io/gh/nao1215/morrigan)
[![reviewdog](https://github.com/nao1215/morrigan/actions/workflows/reviewdog.yml/badge.svg)](https://github.com/nao1215/morrigan/actions/workflows/reviewdog.yml)
[![Go Report Card](https://goreportcard.com/badge/github.com/nao1215/morrigan)](https://goreportcard.com/report/github.com/nao1215/morrigan)
# morrigan - Penetration Tool Set
morrigan command is a tool-set to verify the vulnerability of services developed by you. It is not a tool to attack services on the network. As part of the developer's (it's me) study of security, I will be adding subcommands regarding penetration.  
  
### **morrigan sub-command list**
Each subcommand is explained on other pages. We plan to add more subcommands. We will be reimplementing existing commands used in penetration in golang. However, we will also implement morrigan's own subcommands.
  
|sub-command | description |orginal author or inspired by|
|:--|:--|:--|
|[netcat](./docs/netcat.md)| listen TCP/UDP ports and send data to remote ports over TCP/UDP|[Vadym Fedorov](https://github.com/vfedoroff) (forked from [go-netcat](https://github.com/vfedoroff/go-netcat))|
|[pwcrack](./docs/pwcrack.md)| crack local user password|Naohiro CHIKAMATSU (inspired by [John the ripper](https://www.openwall.com/john/))|
|[pwscore](./docs/pwscore.md)| [WIP] check password strength|Naohiro CHIKAMATSU (inspired by [libpwquality](https://github.com/libpwquality/libpwquality))|
|[unshadow](./docs/unshadow.md)| combine password fields in /etc/passwd and /etc/shadow|Naohiro CHIKAMATSU (inspired by [John the ripper](https://www.openwall.com/john/))|
|[zip-pwcrack](./docs/zip-pwcrack.md)|crack zip password|[ICHINOSE Shogo](https://github.com/shogo82148)|  
  
### **Legal Warning**
> With great power comes great responsibility.

morrign command is **under development**. There are no features that compromise security in any way. However,I will be adding subcommands for penetration in stages.  

**Please use the morrigan command only on PCs and servers that you control.** Do not use morrigan command in military or secret service organizations, or for illegal purposes. You may be held legally liable for your actions.  


# How to install
### Use "go install"
If you does not have the golang development environment installed on your system, please install golang from the [golang official website](https://go.dev/doc/install).
```
$ go install github.com/nao1215/morrigan@latest
```

### Install from Package or Binary
[The release page](https://github.com/nao1215/morrigan/releases) contains packages in .deb, .rpm, and .apk formats.


# Contributing
First off, thanks for taking the time to contribute! ❤️  See [CONTRIBUTING.md](./CONTRIBUTING.md) for more information.
Contributions are not only related to development. For example, GitHub Star motivates me to develop!!  
[![Star History Chart](https://api.star-history.com/svg?repos=nao1215/morrigan&type=Date)](https://star-history.com/#nao1215/morrigan&Date)

# Contact
If you would like to send comments such as "find a bug" or "request for additional features" to the developer, please use one of the following contacts.

- [GitHub Issue](https://github.com/nao1215/morrigan/issues)

# LICENSE
The morrigan project is licensed under the terms of [MIT LICENSE](./LICENSE).