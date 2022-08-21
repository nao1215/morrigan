# morrigan - [WIP] Penetration Tool Set
morrigan command is a tool-set to verify the vulnerability of services developed by you. It is not a tool to attack services on the network. As part of the developer's (it's me) study of security, I will be adding subcommands regarding penetration.  
  
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


# morrigan sub-command list
|sub-command | description |
|:--|:--|
|pwcrack| [WIP] crack local user password|
|pwscore| [WIP] check password strength|
|unshadow| [WIP] combine password fields in /etc/passwd and /etc/shadow|

# Contributing
First off, thanks for taking the time to contribute! ❤️  See [CONTRIBUTING.md](./CONTRIBUTING.md) for more information.
Contributions are not only related to development. For example, GitHub Star motivates me to develop!!  

# Contact
If you would like to send comments such as "find a bug" or "request for additional features" to the developer, please use one of the following contacts.

- [GitHub Issue](https://github.com/nao1215/morrigan/issues)

# LICENSE
The morrigan project is licensed under the terms of [MIT LICENSE](./LICENSE).