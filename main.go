package main

import (
	"fmt"
	"log"
	"net"

	"github.com/influxdata/go-syslog/rfc5424"
	"gopkg.in/alecthomas/kingpin.v2"
)

var (
	viewRaw = kingpin.Flag("raw", "View raw message").Short('r').Bool()
	tcp     = kingpin.Flag("tcp", "Listen on tcp (601)").Short('t').Bool()
)

func printMessage(i []byte, src net.Addr) {

	p := rfc5424.NewParser()
	bestEffort := true
	m, parseErr := p.Parse(i, &bestEffort)
	if parseErr != nil {
		log.Fatal(parseErr)
	}
	fmt.Println("\n>> MESSAGE RECEIVED:")
	if *viewRaw {
		fmt.Println("   RAW:       ", string(i))
	}

	fmt.Println("   Source IP: ", src.String())

	if m.Timestamp() != nil {
		fmt.Println("   Timestamp: ", m.Timestamp())
	}
	if m.Severity() != nil {
		fmt.Println("   Severity:  ", *m.Severity())
	}
	if m.Facility() != nil {
		fmt.Println("   Facility:  ", *m.Facility())
	}
	if m.Hostname() != nil {
		fmt.Println("   Hostname:  ", *m.Hostname())
	}
	if m.Appname() != nil {
		fmt.Println("   Appname:   ", *m.Appname())
	}
	if m.ProcID() != nil {
		fmt.Println("   PID:       ", *m.ProcID())
	}
	if m.MsgID() != nil {
		fmt.Println("   Msg id:    ", *m.MsgID())
	}
	if m.Message() != nil {
		fmt.Println("   Payload:   ", *m.Message())
	}

}

func main() {
	//parse cli
	kingpin.Version("1.0")
	kingpin.Parse()

	if *tcp {
		listen, err := net.Listen("tcp", ":601")
		if err != nil {
			log.Fatal(err)
		}

		defer listen.Close()
		for {
			buf := make([]byte, 2048)
			conn, err := listen.Accept()
			if err != nil {
				log.Fatal(err)
			}
			_, err = conn.Read(buf)
			if err != nil {
				log.Fatal(err)
			}
			addr := conn.RemoteAddr()
			conn.Close()
			printMessage(buf, addr)
		}
	} else {
		udpServer, err := net.ListenPacket("udp", ":514")
		if err != nil {
			log.Fatal(err)
		}
		defer udpServer.Close()

		for {
			buf := make([]byte, 2048)
			_, addr, err := udpServer.ReadFrom(buf)
			if err != nil {
				log.Fatal(err)
				continue
			}
			printMessage(buf, addr)
		}
	}
}
