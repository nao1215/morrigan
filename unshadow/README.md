# unshadow
unshadow package provide function that combine password fields in /etc/passwd and /etc/shadow. unshadow is known for its distribution of John The Ripper, a password cracking tool.

# Example
```go
import (
	"fmt"
	"github.com/nao1215/morrigan/unshadow"
)

func main() {
	unshadowList, err := unshadow.Unshadow(unshadow.PasswdFilePath, unshadow.ShadowFilePath)
	if err != nil {
		fmt.Errorf("error:", err)
		return
	}

	for _, v := range unshadowList {
		fmt.Fprintln(os.Stdout, v)
	}
	return nil
}
```


# What is /etc/passwd
The "/etc/passwd" file contains account information for each user, one user per line. There are a total of seven fields, separated by colons.

|no|field|description|
|:--|:--|:--|
|1|user name| user account name|
|2|password|"x":use /etc/shadow<br>"*":deactivate account<br>not set: no password |
|3|user id|uid|
|4|group id|gid|
|5|comment|comment such as full name and role is described.|
|6|home directory|absolute path to user home directory|
|7|login shell|shell to be applied when a user logs in|
  
If the password field in /etc/passwd is x, unshadow subcommand replaces x with the encrypted password listed in /etc/shadow.
  
# What is /etc/shadow
The "/etc/shadow" file contains the encrypted password and information related to the password, and can only be read and written by the root user.


|no|field|description|
|:--|:--|:--|
|1|user name| user account name|
|2|password|encrypted password|
|3|date of last password change|The date the password was last changed (expressed as the number of days elapsed since January 1, 1970)|
|4|number of days available to change password|The number of days before the password can be changed again|
|5|password validity period| The number of days before the password needs to be changed|
|6|password change period warning notification date|The number of days prior to the password expiration warning|
|7|number of days until account inactivity after password expiration|Number of days before account becomes inactive if password is not changed after expiration date|
|8|account validity period|The number of days (expressed as the number of days elapsed since January 1, 1970) before the account becomes unavailable|
|9|reserved Field| not used|
  
Passwords are encrypted using md5, sha256, sha512, etc. The encryption (hash calculation) is described in [gocrypt](../../../nao1215/morrigan/gocrypt/README.md).