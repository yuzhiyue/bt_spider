package dht

import "net"

type DHTNode struct{
    ID string
    Ip string
    Port uint16
}

type DHTServer struct{
    ID string
    Conn net.UDPConn
}

func (this *DHTServer) SendMsg(address string, msg interface{}) error {

    DHTServer.Conn.Write()
}

func (this *DHTServer) findNode(target string) error {

}

func getP()  {
    
}
