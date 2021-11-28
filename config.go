package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"time"

	"gopkg.in/yaml.v3"
)

type (
	hostST struct {
		Name   string `yaml:"name"`
		Access string `yaml:"access"`
	}
	configAgentST struct {
		Bot      string   `yaml:"bot"`
		Interval int      `yaml:"interval"`
		Secret   string   `yaml:"secret"`
		Agent    string   `yaml:"agent"`
		Ping     []hostST `yaml:"ping"`
		Web      []hostST `yaml:"web"`
		Mysql    []hostST `yaml:"mysql"`
		Ldap     []hostST `yaml:"ldap"`
		Dns      []hostST `yaml:"dns"`
	}
)

var (
	serverAddr       = ""
	intervalDuration = time.Second
	secret           = ""
	hosts            = []host{}
)

func configProcess() {
	if len(os.Args) > 1 {
		yfile, err := ioutil.ReadFile(os.Args[1])
		if err != nil {
			fmt.Printf("Cannot read file config\n\n---------------\n\n")
			usage()
		}
		configDat := configAgentST{}
		err2 := yaml.Unmarshal(yfile, &configDat)
		if err2 != nil {
			fmt.Printf("Error in file config\n%v\n----------------\n\n", err2)
			usage()
		}
		serverAddr = configDat.Bot
		intervalDuration = time.Duration(configDat.Interval) * time.Second
		secret = configDat.Secret
		for _, v := range configDat.Ping {
			h := host{}
			h.name = v.Name
			h.access = v.Access
			h.htype = "ping"
			hosts = append(hosts, h)
		}
		for _, v := range configDat.Web {
			h := host{}
			h.name = v.Name
			h.access = v.Access
			h.htype = "web"
			hosts = append(hosts, h)
		}
		for _, v := range configDat.Mysql {
			h := host{}
			h.name = v.Name
			h.access = v.Access
			h.htype = "mysql"
			hosts = append(hosts, h)
		}
		for _, v := range configDat.Ldap {
			h := host{}
			h.name = v.Name
			h.access = v.Access
			h.htype = "ldap"
			hosts = append(hosts, h)
		}
		for _, v := range configDat.Dns {
			h := host{}
			h.name = v.Name
			h.access = v.Access
			h.htype = "dns"
			hosts = append(hosts, h)
		}
		if configDat.Agent != "" {
			h := host{}
			h.name = configDat.Agent
			h.htype = "watch"
			hosts = append(hosts, h)
		}
		fmt.Printf("bot = %v\n", serverAddr)
		fmt.Printf("interval = %v\n", intervalDuration)
		fmt.Printf("secret = %v\n", secret)
		for _, v := range hosts {
			fmt.Printf("%v\n", v)
		}
		return
	}
	usage()
}

func usage() {
	fmt.Printf("usage: agentalert <config.yaml>\n\nExample config:\n-------------\n")
	fmt.Printf("bot: \"127.0.0.1:9055\"\n")
	fmt.Printf("interval: 30\n")
	fmt.Printf("secret: \"Secret text that share with agentalert\"\n")
	fmt.Printf("agent: \"agent1\"\n")
	fmt.Printf("ping:\n")
	fmt.Printf("  - name: \"Server1\"\n")
	fmt.Printf("    access: \"10.10.10.10\"\n")
	fmt.Printf("  - name: \"Server2\"\n")
	fmt.Printf("    access: \"10.10.10.11\"\n")
	fmt.Printf("web:\n")
	fmt.Printf("  - name: \"Web1\"\n")
	fmt.Printf("    access: \"https://web1.com\"\n")
	fmt.Printf("  - name: \"Web2\"\n")
	fmt.Printf("    access: \"http://web2.org\"\n")
	fmt.Printf("mysql:\n")
	fmt.Printf("  - name: \"MySQL Serv1\"\n")
	fmt.Printf("    access: \"user1:password@tcp(10.10.10.10:3306)/db\"\n")
	fmt.Printf("  - name: \"MySQL Serv2\"\n")
	fmt.Printf("    access: \"user2:password@tcp(10.10.10.11:3306)/db\"\n")
	fmt.Printf("ldap:\n")
	fmt.Printf("  - name: \"LDAP1\"\n")
	fmt.Printf("    access: \"10.10.10.10\"\n")
	fmt.Printf("  - name: \"LDAP2\"\n")
	fmt.Printf("    access: \"10.10.10.11\"\n")
	fmt.Printf("dns:\n")
	fmt.Printf("  - name: \"DNS1\"\n")
	fmt.Printf("    access: \"10.10.10.10\"\n")
	fmt.Printf("  - name: \"DNS2\"\n")
	fmt.Printf("    access: \"10.10.10.11\"\n\n")
	os.Exit(0)
}
