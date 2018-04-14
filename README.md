# Webtail

### Yet another browser based remote log viewer written in Golang

Webtail is a web-socket based server that streams log files onto your browser. Written in [Golang](https://golang.org), this application provides a basic and clean material UI to view the logs as well.

### Usage

```
$ webtail --help
usage: webtail [<flags>] [<dir>...]

Flags:
      --help                 Show context-sensitive help (also try --help-long and --help-man).
  -p, --port=8080            Port number to host the server
  -r, --restrict             Enforce PAM authentication (single level)
  -w, --whitelist=WHITELIST  enable whitelisting with users in the provided file
  -c, --cron="1h"            configure cron for re-indexing files (Not supported right now)

Args:
  [<dir>]  Directory path(s) to look for files

```

To view the UI, navigate to *http://server_ip:port* and you will be presented with a UI to view the logs.

#### Examples
```
./webtail
```
This will run the server on port `8080` and look for files in the current Directory
```
./webtail --port 15000 /var/log/tomcat /tmp/
```
This will run the server on port 15000 and recursively look for files in `/var/log/tomcat` and `/tmp` directories (provided the permissions)

```
./webtail /var/log/tomcat /tmp/ --restrict
```
This will add an authentication layer over it. Once you navigate to the home page, it will redirect to the `/login` page and ask for username and password. Since this is supposed to be as generic as possible, hence it uses PAM authentication to authenticate the user. You need to provide the credentials that you would use to login to the host on which the server is hosted. Right now it would only authenticate via PAM if only a single step is required.

For this it uses [CGO](https://github.com/golang/go/wiki/cgo). A basic starting point for this would be [Calling Go functions from C](https://medium.com/using-go-in-mobile-apps/using-go-in-mobile-apps-part-1-calling-go-functions-from-c-be1ecf7dfbc6). This also has information on how to call C funtions from Golang, which is what is being used in this module. It is performing PAM authentication in [pam_auth.go](https://github.com/prateeknischal/webtail/blob/master/util/pam_auth.go).
This contains the code to interact with the system's PAM and check if a username/password combination is valid or not.

```
./webtail /var/log/tomcat --restrict --whitelist ~/allowed-users.txt
```
This will have the authentication layer with another layer of ACL, allowing only those users to authenticate with the server which are in the `~/allowed-users.txt`. It accepts a file that has newline separated usernames.


### Build From Source

To build from source, please refer the Makefile, all you need to do is change the `GOOS` and `GOARCH` variables to your specific distributions. The Makefile has it configured for `linux/amd64` by default.

```
make clean package
```

**Note**    
If you get the following error:
```
security/pam_appl.h: No such file or directory
```
Then you need to install PAM developement libraries.    
*CentOS* : `yum install pam-devel`    
*Debian* : `sudo apt-get install -y libpam0g-dev`   
Reference: [pam_appl.h and pam_misc.h missing](https://stackoverflow.com/questions/15614823/pam-appl-h-and-pam-misc-h-missing-in-rshd-c-source-code)

It uses the `LDFLAG -lpam` to build the C code, so if you get errors related to this, worth looking for.

### Screenshots
Login
![N|Solid](https://raw.githubusercontent.com/prateeknischal/webtail/master/screenshots/webtail_login.png)

Dashboard
![N|Solid](https://raw.githubusercontent.com/prateeknischal/webtail/master/screenshots/webtail_dashboard.png)

SideNav
![N|Solid](https://raw.githubusercontent.com/prateeknischal/webtail/master/screenshots/webtail_filenav.png)

Tail
![N|Solid](https://raw.githubusercontent.com/prateeknischal/webtail/master/screenshots/webtail_tail.png)

#### TODOs
* Add https support
* Add cron support to re-index files in the provided directories
* Add a proper logger
* Any help in UI is most welcome
