package main

import (
	nex "github.com/PretendoNetwork/nex-go"
	"go.mongodb.org/mongo-driver/bson"
)

func generateKerberosTicket(userPID uint32, serverPID uint32, keySize int) ([]byte, int) {
	var userPassword string

	if userPID == 100 { // "guest" account
		userPassword = "MMQea3n!fsik"
	} else {
		user := getUserByPID(userPID)

		if user == nil {
			return []byte{}, 0x80030064 // RendezVous::InvalidUsername
		}

		userPassword = user["nex"].(bson.M)["password"].(string)
	}

	// Create session key and ticket keys
	sessionKey := make([]byte, keySize)
	//rand.Read(sessionKey) // TODO: Create a random session key
	ticketInfoKey := make([]byte, 16)                   // key for encrypting the internal ticket info. Only used by server. TODO: Make this random!
	userKey := deriveKey(userPID, []byte(userPassword)) // Key for encrypting entire ticket. Used by client and server

	////////////////////////////////
	// Build internal ticket info //
	////////////////////////////////

	ticketInfoStream := nex.NewStreamOut(nexServer)

	// Write ticket expiration time
	expiration := nex.NewDateTime(0)
	ticketInfoStream.WriteUInt64LE(expiration.Now())

	// Write user PID
	ticketInfoStream.WriteUInt32LE(userPID)

	// Write session key
	ticketInfoStream.Grow(int64(keySize))
	ticketInfoStream.WriteBytesNext(sessionKey)

	// Encrypt internal ticket info

	ticketInfoEncryption := nex.NewKerberosEncryption(ticketInfoKey)
	encryptedTicketInfo := ticketInfoEncryption.Encrypt(ticketInfoStream.Bytes())

	//////////////////////
	// Build ticket data//
	//////////////////////

	ticketDataStream := nex.NewStreamOut(nexServer)

	// Write ticket info
	ticketDataStream.Grow(int64(len(encryptedTicketInfo)))
	ticketDataStream.WriteBytesNext(encryptedTicketInfo)

	///////////////////////////
	// Build Kerberos Ticket //
	///////////////////////////

	ticketStream := nex.NewStreamOut(nexServer)

	// Write session key
	ticketStream.Grow(int64(keySize))
	ticketStream.WriteBytesNext(sessionKey)

	// Write Server PID
	ticketStream.WriteUInt32LE(serverPID)

	// Write encrypted ticket data
	ticketStream.WriteBuffer(ticketDataStream.Bytes())

	// Encrypt the ticket
	ticketEncryption := nex.NewKerberosEncryption(userKey)
	encryptedTicket := ticketEncryption.Encrypt(ticketStream.Bytes())

	return encryptedTicket, 0
}

func deriveKey(pid uint32, password []byte) []byte {
	for i := 0; i < 65000+int(pid)%1024; i++ {
		password = nex.MD5Hash(password)
	}

	return password
}
