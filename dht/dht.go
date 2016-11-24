package dht

import "net"
import (
    "github.com/zeebo/bencode"
    "fmt"
)

type DHTNode struct{
    ID string
    Address net.UDPAddr
}

type DHTServer struct{
    ID string
    Conn net.UDPConn
    IP string;
    Port uint16;
}

func (this *DHTServer) sendMsg(address net.UDPAddr, msg interface{}) error {
    data, err := bencode.EncodeBytes(msg)
    if err != nil {
        return err
    }
    fmt.Println(data)
    _, err = DHTServer.Conn.WriteToUDP(data, address)
    return err
}

func (this *DHTServer) recvMsg(address net.UDPAddr, msg interface{}) error {
    krcpMsg := msg.(map[string]interface{})
    krcpMsg["r"]
    return nil
}

func (this *DHTServer) findNode(target string) error {
    return nil
}

func (this *DHTServer) Run() {
    addr, err := net.ResolveUDPAddr("udp", this.IP + ":" + this.Port)
    if err != nil {
        fmt.Println("net.ResolveUDPAddr fail.", err)
        return
    }

    this.Conn, err = net.DialUDP("udp", nil, addr)
    if err != nil {
        fmt.Println("net.DialUDP fail.", err)
        return
    }
    defer this.Conn.Close()
    var Buff [65535]byte;
    for {
        rlen, remote, err := this.Conn.ReadFromUDP(Buff)
        if err == nil {
            var msg interface{}
            bencode.DecodeBytes(Buff[:rlen], msg)
            this.recvMsg(remote, msg)
        }
    }
}


