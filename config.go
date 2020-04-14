package main

import (
	"time"
)

var (
	serverAddr       = "127.0.0.1:9055"
	intervalDuration = 30 * time.Second
	secret           = "Secret text that share with agentalert"
)

var (
	hosts = []host{
		host{name: "Server", access: "10.10.10.10", htype: "ping"},
		host{name: "Web Server", access: "https://web.com", htype: "web"},
		host{name: "MySQL Server", access: "user:password@tcp(10.10.10.10:3306)/db", htype: "mysql"},
		host{name: "LDAP Server", access: "10.10.10.10", htype: "ldap"},
		host{name: "DNS Server", access: "10.10.10.10", htype: "dns"},
		host{name: "My Host", access: "", htype: "watch"},
	}
)
