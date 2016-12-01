package dht

import "net"
import (
    "github.com/zeebo/bencode"
    "fmt"
    "strconv"
)

type DHTNode struct{
    ID string
    Address net.UDPAddr
}

type Trans struct {
    Data map[string]interface{}
    Timeout uint32
}

type DHTServer struct{
    ID string
    Conn net.UDPConn
    IP string;
    Port uint16;
    transMap map[string]*Trans
    currTransID uint32
}

func (this *DHTServer) Init() error {
    this.transMap = make(map[string]Trans);
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
    return nil
}

func (this *DHTServer) sendMsg(address net.UDPAddr, t string, y string, q string, msg interface{}) error {
    if t == nil {
        t = strconv.FormatInt(int64(this.currTransID), 16)
        this.currTransID++
    }
    krcpMsg := msg.(map[string]interface{})
    krcpMsg["t"] = t
    krcpMsg["y"] = y

    if y == "q" {
        krcpMsg["q"] = q
        krcpMsg["a"] = msg
    } else {
        krcpMsg["r"] = msg
    }

    transData := new(Trans)
    transData.Data = krcpMsg
    transData.Timeout = 30
    this.transMap[t] = transData
    data, err := bencode.EncodeBytes(krcpMsg)
    if err != nil {
        return err
    }
    fmt.Println(data)
    _, err = DHTServer.Conn.WriteToUDP(data, address)
    return err
}

func (this *DHTServer) recvMsg(address net.UDPAddr, msg interface{}) error {
    krcpMsg := msg.(map[string]interface{})
    t := krcpMsg["t"].(string)
    y := krcpMsg["y"].(string)

    if y == "q" {
        q := krcpMsg["q"].(string)
        if q == "ping" {

        } else if q == "find_node" {

        } else if q == "get_peers" {

        } else if q == "announce_peer" {

        }
    } else if y == "r" {
        tranData := this.transMap[t]
        reqMsg := tranData.Data
        q := reqMsg["q"].(string)
        if q == "ping" {

        } else if q == "find_node" {

        } else if q == "get_peers" {

        } else if q == "announce_peer" {

        }
    }
    
    return nil
}

func (this *DHTServer) findNode(addr net.UDPAddr, target string) error {
    reqMsg := map[string]interface{} {"id":this.ID, "target": this.ID}
    this.sendMsg(addr, nil, "y", "find_node", reqMsg)
    return nil
}

func (this *DHTServer) onFindNode(addr net.UDPAddr, t string, msg map[string]interface{})  {
    rspMsg := map[string]interface{} {"id":this.ID, "nodes": ""}
    this.sendMsg(addr, t, "r", nil, rspMsg)
}

func (this *DHTServer) Run() {
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


