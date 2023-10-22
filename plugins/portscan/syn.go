package portscan

// import (
// 	"net"
// 	"time"

// 	"github.com/google/gopacket"
// 	"github.com/google/gopacket/layers"
// )

// func TcpSYN(srcIp net.IP, srcPort int, dstIp string, dstPort int) {
// 	dstAddrs, err := net.LookupIP(dstIp)
// 	if err != nil {
// 		return
// 	}
// 	dstip := dstAddrs[0].To4()
// 	srcip := srcIp

// 	// Our port
// 	var (
// 		dstport layers.TCPPort = layers.TCPPort(dstPort)
// 		srcport layers.TCPPort = layers.TCPPort(srcPort)
// 	)

// 	// Our TCP header
// 	tcp := &layers.TCP{
// 		SrcPort: srcport,
// 		DstPort: dstport,
// 		SYN:     true,
// 	}

// 	// Our IP header
// 	// not used, but necessary for TCP checksumming.
// 	ip := &layers.IPv4{
// 		SrcIP:    srcip,
// 		DstIP:    dstip,
// 		Protocol: layers.IPProtocolTCP,
// 	}

// 	if err := tcp.SetNetworkLayerForChecksum(ip); err != nil {
// 		return
// 	}

// 	buf := gopacket.NewSerializeBuffer()
// 	opts := gopacket.SerializeOptions{
// 		FixLengths:       true,
// 		ComputeChecksums: true,
// 	}

// 	if err := gopacket.SerializeLayers(buf, opts, tcp); err != nil {
// 		return
// 	}

// 	// listen on local TCP connection
// 	conn, err := net.ListenPacket("ip4:tcp", "0.0.0.0")
// 	if err != nil {
// 		return
// 	}
// 	defer conn.Close()

// 	// send TCP SYN packet
// 	if _, err := conn.WriteTo(buf.Bytes(), &net.IPAddr{IP: dstip}); err != nil {
// 		return
// 	}
// 	// Set deadline so we do not wait forever.
// 	if err := conn.SetDeadline(time.Now().Add(time.Duration(10) * time.Second)); err != nil {
// 		return
// 	}

// 	for {
// 		b := make([]byte, 4096)
// 		// func (c *IPConn) ReadFrom(b []byte) (int, Addr, error)
// 		n, addr, err := conn.ReadFrom(b)
// 		if err != nil {
// 			return
// 		} else if addr.String() == dstip.String() {
// 			// Decode a packet
// 			packet := gopacket.NewPacket(b[:n], layers.LayerTypeTCP, gopacket.Default)
// 			// Get the TCP layer from this packet
// 			if tcpLayer := packet.Layer(layers.LayerTypeTCP); tcpLayer != nil {
// 				tcp, _ := tcpLayer.(*layers.TCP)

// 				if tcp.SrcPort == dstport {
// 					if tcp.SYN && tcp.ACK {
// 						//PortCheck(dstIp, dstPort,)
// 						//return dstPort, err
// 					} else {
// 						return
// 					}
// 				}
// 			}
// 		}
// 	}

// }
