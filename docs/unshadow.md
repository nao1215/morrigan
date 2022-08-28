# unshadow
unshadow subcommand combine password fields in /etc/passwd and /etc/shadow. unshadow is known for its distribution of John The Ripper, a password cracking tool.

# Synopsis
```
morrigan unshadow ETC_PASSWD_FILE ETC_SHADOW_FILE [flags]
```

# Examples
unshadow subcommand reads /etc/unshadow, so you execute it with root privileges.
```
sudo morrigan unshadow ETC_PASSWD_FILE ETC_SHADOW_FILE
```
  
If you execute unshadow subcommand without specifying the file paths, unshadow will attempt to read /etc/passwd and /etc/shadow on the local system.
```
sudo morrigan unshadow
```

The output of the unshadow subcommand is sorted alphabetically.
```
$ sudo morrigan unshadow 
_apt:*:105:65534::/nonexistent:/usr/sbin/nologin
_flatpak:*:129:138:Flatpak system-wide installation helper,,,:/nonexistent:/usr/sbin/nologin
avahi-autoipd:*:109:117:Avahi autoip daemon,,,:/var/lib/avahi-autoipd:/usr/sbin/nologin
avahi:*:113:120:Avahi mDNS daemon,,,:/run/avahi-daemon:/usr/sbin/nologin
backup:*:34:34:backup:/var/backups:/usr/sbin/nologin
bin:*:2:2:bin:/bin:/usr/sbin/nologin
colord:*:121:127:colord colour management daemon,,,:/var/lib/colord:/usr/sbin/nologin
cups-pk-helper:*:114:121:user for cups-pk-helper service,,,:/home/cups-pk-helper:/usr/sbin/nologin
daemon:*:1:1:daemon:/usr/sbin:/usr/sbin/nologin
dnsmasq:*:112:65534:dnsmasq,,,:/var/lib/misc:/usr/sbin/nologin
games:*:5:60:games:/usr/games:/usr/sbin/nologin
gdm:*:126:131:Gnome Display Manager:/var/lib/gdm3:/bin/false
geoclue:*:122:128::/var/lib/geoclue:/usr/sbin/nologin
gnats:*:41:41:Gnats Bug-Reporting System (admin):/var/lib/gnats:/usr/sbin/nologin
gnome-initial-setup:*:125:65534::/run/gnome-initial-setup/:/bin/false
gup:*:134:9:Group Update Program,,,:/var/lib/gup/:/usr/sbin/nologin
hplip:*:124:7:HPLIP system user,,,:/run/hplip:/bin/false
irc:*:39:39:ircd:/run/ircd:/usr/sbin/nologin
kernoops:*:116:65534:Kernel Oops Tracking Daemon,,,:/:/usr/sbin/nologin
lightdm:*:127:134:Light Display Manager:/var/lib/lightdm:/bin/false
list:*:38:38:Mailing List Manager:/var/list:/usr/sbin/nologin
lp:*:7:7:lp:/var/spool/lpd:/usr/sbin/nologin
mail:*:8:8:mail:/var/mail:/usr/sbin/nologin
man:*:6:12:man:/var/cache/man:/usr/sbin/nologin
messagebus:*:103:106::/nonexistent:/usr/sbin/nologin
mimixbox:!:1001:1001::/home/mimixbox:/bin/sh
morrie:$y$j9T$AWBxIYtBpRyJuzyHhC/4M.$LpGFc4mc0F8/f9w150QMsvqku7hofX4r6YIFiUFiEj1:1002:1002::/home/morrie:/bin/sh
mysql:!:132:144:MySQL Server,,,:/nonexistent:/bin/false
news:*:9:9:news:/var/spool/news:/usr/sbin/nologin
nm-openvpn:*:117:122:NetworkManager OpenVPN,,,:/var/lib/openvpn/chroot:/usr/sbin/nologin
nobody:*:65534:65534:nobody:/nonexistent:/usr/sbin/nologin
 :
 :
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