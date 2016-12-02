package dht

import "net"
import (
    "github.com/zeebo/bencode"
    "fmt"
    "strconv"
    "bytes"
    "encoding/binary"
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
    Conn *net.UDPConn
    IP string;
    Port uint16;
    transMap map[string]*Trans
    currTransID uint32
}

func (this *DHTServer) Init() error {
    this.transMap = make(map[string]*Trans);
    addr, err := net.ResolveUDPAddr("udp", this.IP + ":" + string(this.Port))
    if err != nil {
        fmt.Println("net.ResolveUDPAddr fail.", err)
        return err
    }

    this.Conn, err = net.DialUDP("udp", nil, addr)
    if err != nil {
        fmt.Println("net.DialUDP fail.", err)
        return err
    }
    return nil
}

func (this *DHTServer) sendMsg(address *net.UDPAddr, t string, y string, q string, msg interface{}) error {
    fmt.Println("send  msg", msg)
    if t == "" {
        t = strconv.FormatInt(int64(this.currTransID), 16)
        this.currTransID++
    }
    krcpMsg := make(map[string]interface{})
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
    fmt.Println("send krpc msg", krcpMsg)
    data, err := bencode.EncodeBytes(krcpMsg)
    if err != nil {
        return err
    }
    fmt.Println(data)
    _, err = this.Conn.WriteToUDP(data, address)
    return err
}

func (this *DHTServer) recvMsg(address *net.UDPAddr, msg interface{}) error {
    krcpMsg := msg.(map[string]interface{})
    t := krcpMsg["t"].(string)
    y := krcpMsg["y"].(string)

    if y == "q" {
        q := krcpMsg["q"].(string)
        a := krcpMsg["a"].(map[string]interface{})
        if q == "ping" {

        } else if q == "find_node" {
            this.onFindNodeRequest(address, t, a)
        } else if q == "get_peers" {

        } else if q == "announce_peer" {

        }
    } else if y == "r" {
        tranData := this.transMap[t]
        reqMsg := tranData.Data
        q := reqMsg["q"].(string)
        r := krcpMsg["r"].(map[string]interface{})
        if q == "ping" {

        } else if q == "find_node" {
            this.onFindNodeResponse(address, r)
        } else if q == "get_peers" {

        } else if q == "announce_peer" {

        }
    }
    
    return nil
}

func (this *DHTServer) FindNode(addr *net.UDPAddr, target string) error {
    reqMsg := map[string]interface{} {"id":this.ID, "target": this.ID}
    this.sendMsg(addr, "", "q", "find_node", reqMsg)
    return nil
}

func (this *DHTServer) onFindNodeRequest(addr *net.UDPAddr, t string, msg map[string]interface{})  {
    rspMsg := map[string]interface{} {"id":this.ID, "nodes": ""}
    this.sendMsg(addr, t, "r", "", rspMsg)
}

func inet_ntoa(ipnr uint32) string {
    var bytes [4]byte
    bytes[0] = byte(ipnr & 0xFF)
    bytes[1] = byte((ipnr >> 8) & 0xFF)
    bytes[2] = byte((ipnr >> 16) & 0xFF)
    bytes[3] = byte((ipnr >> 24) & 0xFF)

    return net.IPv4(bytes[3],bytes[2],bytes[1],bytes[0]).String()
}

func (this * DHTServer) onFindNodeResponse(addr *net.UDPAddr, msg map[string]interface{})  {
    nodes := []byte(msg["nodes"].(string))
    if len(nodes) % 26 != 0 {
        return
    }
    reader := bytes.NewReader(nodes)
    for i := 0; i < len(nodes); i += 26 {
        var id [20]byte
        var ip uint32
        var port uint16
        binary.Read(reader, binary.BigEndian, &id)
        binary.Read(reader, binary.BigEndian, &ip)
        binary.Read(reader, binary.BigEndian, &port)
        fmt.Println("new node:", string(id[:]), inet_ntoa(ip), port)
    }
}

func (this *DHTServer) Run() {
    defer this.Conn.Close()
    var Buff [65535]byte;
    for {
        rlen, remote, err := this.Conn.ReadFromUDP(Buff[:])
        if err == nil {
            var msg interface{}
            bencode.DecodeBytes(Buff[:rlen], msg)
            this.recvMsg(remote, msg)
            fmt.Println("recv krpc msg", msg)
        }
    }
}


