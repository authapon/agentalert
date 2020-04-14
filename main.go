package main

import (
	"database/sql"
	"fmt"
	"net"
	"net/http"
	"sync"
	"time"

	_ "github.com/go-sql-driver/mysql"
	ldapserv "github.com/mavricknz/ldap"

	mc "github.com/authapon/mcryptzero"
	"github.com/miekg/dns"
	pp "github.com/sparrc/go-ping"
)

type (
	host struct {
		name   string
		htype  string
		access string
	}
)

var (
	mutex = &sync.Mutex{}
)

func ping(host string) bool {
	pinger, err := pp.NewPinger(host)
	if err != nil {
		return false
	}

	pinger.SetPrivileged(true)
	pinger.Count = 3
	pinger.Timeout = 1 * time.Second
	pinger.Run()
	stats := pinger.Statistics()

	if stats.PacketsRecv == 0 {
		return false
	}

	if stats.PacketsRecv != stats.PacketsSent {
		return false
	}

	return true
}

func checkWeb(access string) bool {
	client := &http.Client{}
	resp, err := client.Get(access)
	if err != nil {
		return false
	}

	resp.Body.Close()
	return true
}

func UDPsend(text string) {
	mutex.Lock()
	defer mutex.Unlock()

	serverAddr, err := net.ResolveUDPAddr("udp", serverAddr)
	if err != nil {
		fmt.Printf("Server Address Error!!!\n")
		panic(err)
	}

	salt := mc.SID(8)
	tcrypt := salt + ":" + string(mc.Encrypt([]byte(text), []byte(salt+secret+salt)))
	conn, err := net.DialUDP("udp", nil, serverAddr)
	if err != nil {
		return
	}

	conn.Write([]byte(tcrypt))
	conn.Close()
}

func checkMySQL(access string) bool {
	con, err := sql.Open("mysql", access)
	if err != nil {
		return false
	}

	defer con.Close()

	err = con.Ping()
	if err != nil {
		return false
	}

	return true
}

func checkDNS(access string) bool {
	c := dns.Client{}
	m := dns.Msg{}
	m.SetQuestion("test.com.", dns.TypeA)
	_, _, err := c.Exchange(&m, access+":53")
	if err != nil {
		return false
	}

	return true
}

func checkLDAP(access string) bool {
	l := ldapserv.NewLDAPConnection(access, 389)
	err := l.Connect()
	if err != nil {
		return false
	}

	l.Close()

	return true
}

func agentRUN(h int) {
	UDPsend("start " + hosts[h].htype + " " + hosts[h].name)
	for {
		check := true
		switch hosts[h].htype {
		case "ping":
			check = ping(hosts[h].access)
		case "web":
			check = checkWeb(hosts[h].access)
		case "mysql":
			check = checkMySQL(hosts[h].access)
		case "dns":
			check = checkDNS(hosts[h].access)
		case "ldap":
			check = checkLDAP(hosts[h].access)
		case "watch":
			check = true
		default:
			check = false
		}

		if check {
			fmt.Printf("%s is successful for %s\n", hosts[h].htype, hosts[h].name)
			UDPsend("up " + hosts[h].htype + " " + hosts[h].name)
		} else {
			fmt.Printf("%s is fail for %s\n", hosts[h].htype, hosts[h].name)
		}

		time.Sleep(intervalDuration)
	}
}

func main() {
	for i := range hosts {
		go agentRUN(i)
	}
	select {}
}
