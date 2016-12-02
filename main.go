package main;

import (
	"bt_spider/dht"
	"fmt"
	"net"
)

func main() {
	s := dht.DHTServer{}
	s.ID = "mnopqrstuvwxyz123456"
	s.IP = "0.0.0.0"
	s.Port = 12000
	s.Init()
	ns, _ := net.LookupHost("router.bittorrent.com")
	ip := ns[0]
	fmt.Println(ip)
	addr, err := net.ResolveUDPAddr("udp", ip+":6881")
	if err != nil {
		fmt.Println("net.ResolveUDPAddr fail.", err)
		return
	}
	s.FindNode(addr, s.ID)
	s.Run()
}
