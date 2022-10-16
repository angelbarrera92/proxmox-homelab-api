package wol

import (
	"bytes"
	"encoding/binary"
	"errors"
	"log"
	"net"

	"github.com/sabhiram/go-wol/wol"
)

// SendMagicPacket to send a magic packet to a given mac address, and optionally
// receives an iface to broadcast on. An iface of "" implies a nil net.UDPAddr
func SendMagicPacket(macAddr, bcastAddr string) error {
	// Construct a MagicPacket for the given MAC Address
	magicPacket, err := wol.New(macAddr)
	if err != nil {
		return err
	}

	// Fill our byte buffer with the bytes in our MagicPacket
	var buf bytes.Buffer
	binary.Write(&buf, binary.BigEndian, magicPacket)
	log.Printf("Attempting to send a magic packet to MAC %s\n", macAddr)
	log.Printf("... Broadcasting to: %s\n", bcastAddr)

	// Get a UDPAddr to send the broadcast to
	udpAddr, err := net.ResolveUDPAddr("udp", bcastAddr)
	if err != nil {
		log.Printf("Unable to get a UDP address for %s\n", bcastAddr)
		return err
	}

	// Open a UDP connection, and defer it's cleanup
	var localAddr *net.UDPAddr
	connection, err := net.DialUDP("udp", localAddr, udpAddr)
	if err != nil {
		log.Printf("ERROR: %s\n", err.Error())
		return errors.New("unable to dial UDP address")
	}
	defer connection.Close()

	// Write the bytes of the MagicPacket to the connection
	bytesWritten, err := connection.Write(buf.Bytes())
	if err != nil {
		log.Printf("Unable to write packet to connection\n")
		return err
	} else if bytesWritten != 102 {
		log.Printf("Warning: %d bytes written, %d expected!\n", bytesWritten, 102)
	}

	return nil
}
