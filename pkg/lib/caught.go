package lib

import (
	"fmt"
	"github.com/google/gopacket"
	"github.com/google/gopacket/layers"
	"github.com/google/gopacket/pcap"
	"log"
	"time"
)

// 抓包 只支持linux操作系统 不支持window、mac

// 获取网络设备
func GetNetworkDevice()  {
	deviceList, err := pcap.FindAllDevs()
	if err != nil {
		log.Fatal(err)
	}
	for _, v := range deviceList {
		fmt.Println("网卡设备:", v.Name, "   IP:", v.Addresses)
	}
}

// 三次握手
func threeWayHandshake(deviceName string) {
	handle, err := pcap.OpenLive(deviceName, 1024, false, time.Second*5)
	if err != nil {
		log.Fatal(err)
	}
	defer handle.Close()
	packetSource := gopacket.NewPacketSource(handle, handle.LinkType())
	for v := range packetSource.Packets() {
		//fmt.Println(v.String())
		if layer4 := v.TransportLayer(); layer4 != nil {
			if tcpLayer, ok := layer4.(*layers.TCP); ok {
				if tcpLayer.DstPort == 8080 || tcpLayer.SrcPort == 8080 {
					//fmt.Println(string(tcpLayer.Payload))
					s := fmt.Sprintf("%d=>%d,SYN=%v,ACK=%v,Payload length=%d,seq=%v,ackNum=%v",
						tcpLayer.SrcPort,
						tcpLayer.DstPort,
						tcpLayer.SYN,
						tcpLayer.ACK,
						len(tcpLayer.Payload),
						tcpLayer.Seq,
						tcpLayer.Ack,
					)
					fmt.Println(s)
				}
			}
		}
	}
}
