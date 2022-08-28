# zip-pwcrack
zip-pwcrack subcommand crack password for zip file by brute force. Golang standard library does not support Zip with password. Therefore, we fork and use [libraries that support Standard Zip Encryption](../,,/../zip/README.md). The reason we forked it was because maintenance had stopped.

# Synopsis
```
morrigan zip-pwcrack FILE_NAME [flags]
```

# Examples
```sh
$ morrigan zip-pwcrack sample.zip 
INFO : detect passwords by brute force...
INFO : zip file's password = test
```