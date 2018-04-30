# Webtail

[![](https://goreportcard.com/badge/github.com/prateeknischal/webtail)](https://goreportcard.com/report/github.com/prateeknischal/webtail)

### Yet another browser based remote log viewer written in Golang

Webtail is a web-socket based server that streams log files onto your browser. Written in [Golang](https://golang.org), this application provides a basic and clean material UI to view the logs as well.

### Usage

```
$ go run main.go --help
usage: main [<flags>] [<dir>...]

Flags:
      --help               Show context-sensitive help (also try --help-long and --help-man).
  -p, --port=8080          Port number to host the server
  -r, --restrict           Enforce PAM authentication (single level)
  -a, --acl=ACL            enable Access Control List with users in the provided file
  -t, --cron="0h"          configure cron for re-indexing files, Supported durations:[h -> hours, d -> days]
  -s, --secure             Run Server with TLS
  -c, --cert="server.crt"  Server Certificate
  -k, --key="server.key"   Server Key File

Args:
  [<dir>]  Directory path(s) to look for files


```

To view the UI, navigate to *http(s)://server_ip:port* and you will be presented with a UI to view the logs.

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

For this it uses [CGO](https://github.com/golang/go/wiki/cgo). A basic starting point for this would be [Calling Go functions from C](https://medium.com/using-go-in-mobile-apps/using-go-in-mobile-apps-part-1-calling-go-functions-from-c-be1ecf7dfbc6). This also has information on how to call C funtions from Golang, which is what is being used in this module. It is performing PAM authentication using the `passwd` service in [pam_auth.go](https://github.com/prateeknischal/webtail/blob/master/util/pam_auth.go).

This contains the code to interact with the system's PAM and check if a username/password combination is valid or not.

Some information about pam authentication : [RedHat - PAM Configuration files](https://access.redhat.com/documentation/en-us/red_hat_enterprise_linux/6/html/managing_smart_cards/pam_configuration_files)


```
./webtail /var/log/tomcat --restrict --acl ~/allowed-users.txt
```
This will have the authentication layer with another layer of ACL, allowing only those users to authenticate with the server which are in the `~/allowed-users.txt`. It accepts a file that has newline separated usernames.


```
./webtail /var/log/tomcat --restrict --cron 5h
```
This will make the server re-index files every 5 hours. And the **New** files will be served only after a page refresh, once it has been indexed.


### Build From Source

To build from source, please refer the Makefile, all you need to do is change the `GOOS` and `GOARCH` variables to your specific distributions. The Makefile has it configured for `linux/amd64` by default.

```
make clean package
```
If you are testing it then:
```
make clean build
```
and then you can run it from the project root

**Note**    
* If you get the following error:
```
security/pam_appl.h: No such file or directory
```
Then you need to install PAM developement libraries.    
*CentOS* : `yum install pam-devel`    
*Debian* : `sudo apt-get install -y libpam0g-dev`   
Reference: [pam_appl.h and pam_misc.h missing](https://stackoverflow.com/questions/15614823/pam-appl-h-and-pam-misc-h-missing-in-rshd-c-source-code)


* If you get the following error, then try building this in a linux box.
(I saw some SO or Github post on this suggesting a fork of golang)
```
/usr/local/go/pkg/tool/darwin_amd64/link: running clang failed: exit status 1
ld: warning: ignoring file /var/folders/30/qpxs8kwj3jzc612r2gsq17zwc2yvq5/T/go-link-136319902/go.o, file was built for unsupported file format ( 0x7F 0x45 0x4C 0x46 0x02 0x01 0x01 0x00 0x00 0x00 0x00 0x00 0x00 0x00 0x00 0x00 ) which is not the architecture being linked (x86_64): /var/folders/30/qpxs8kwj3jzc612r2gsq17zwc2yvq5/T/go-link-136319902/go.o
Undefined symbols for architecture x86_64:
  "__cgo_topofstack", referenced from:
      __cgo_fc9493e99cab_Cfunc_authenticate_system in 000001.o
      __cgo_f7895c2c5a3a_C2func_getnameinfo in 000004.o
      __cgo_f7895c2c5a3a_Cfunc_getnameinfo in 000004.o
      __cgo_f7895c2c5a3a_C2func_getaddrinfo in 000006.o
      __cgo_f7895c2c5a3a_Cfunc_gai_strerror in 000006.o
      __cgo_f7895c2c5a3a_Cfunc_getaddrinfo in 000006.o
  "_main", referenced from:
     implicit entry/start for main executable
ld: symbol(s) not found for architecture x86_64
clang: error: linker command failed with exit code 1 (use -v to see invocation)
```

* It uses the `LDFLAG -lpam` to build the C code, so if you get errors related to this, worth looking for.

### TLS Server

To run the server with TLS enabled use `--secure` flag. It will search for `server.crt` and `server.key` files in the current directory, if not will fail.

```
./webtail /var/log/tomcat --secure --cert /path/to/server.crt --key /path/to/server.key --port 8443 --restrict
```

If you are running it on `--restrict` mode then it is recommended to use `--secure` flag as well to protect the login credentials on the wire.

Server Accepts connections only on `TLSv1.1` and above    
List of CiphersSuites supported by the server:
```
tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384
tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256
tls.TLS_ECDHE_RSA_WITH_AES_128_CBC_SHA256
tls.TLS_RSA_WITH_AES_256_GCM_SHA384
tls.TLS_RSA_WITH_AES_256_CBC_SHA
```

to connect to a particular cipher you can use       
**TLSv1.1**
```
echo "" | openssl s_client -connect localhost:8443 -cipher AES256-SHA -tls1_1 -quiet 2>/dev/null
```
**TLSv1.2**
```
 echo "" | openssl s_client -connect localhost:8443 -cipher ECDHE-RSA-AES128-SHA256 -tls1_2 -quiet 2>/dev/null
```

**Note**: It requires the private key to be in the plain text format. i.e. it should not be passphrase protected. Can be done via
```
openssl rsa -in [file1.key] -out [file2.key]
```

### Cron

Cron option supports only 2 formats of time: `days` and `hours`.      
You can say something like `5h` or `1d` or `100h` or `4d`. Zero prefixed time intervals are not allowed and will fail.    
**Note**: By default cron is not enabled and will not re-index files.

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
* ~~Add https support~~
* ~~Add cron support to re-index files in the provided directories~~
* Add file descriptor pooling while tailing to conserver resources
* Add authentication on websockets (cannot validate CSRF token on GET request in gorilla/csrf, checking Origin header only for now)
* Add a proper logger
* Any help in UI is most welcome
