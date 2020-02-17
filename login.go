package main

import (
	"fmt"
	"strconv"

	nex "github.com/PretendoNetwork/nex-go"
	nexproto "github.com/PretendoNetwork/nex-protocols-go"
)

func login(err error, client *nex.Client, callID uint32, username string) {
	var userPID int

	if username == "guest" {
		userPID = 100
	} else {
		userPID, _ = strconv.Atoi(username)
	}

	var serverPID uint32 = 2 // Quazal Rendez-Vous

	encryptedTicket, errorCode := generateKerberosTicket(uint32(userPID), serverPID, nexServer.GetKerberosKeySize())

	if errorCode != 0 {
		fmt.Println(userPID)
		return
	}

	// Build the response body
	stationURL := "prudps:/address=192.168.0.27;port=60001;CID=1;PID=2;sid=1;stream=10;type=2"
	serverName := "Pretendo Friends Auth"

	rvConnectionData := nex.NewRVConnectionData()
	rvConnectionData.SetStationURL(stationURL)
	rvConnectionData.SetSpecialProtocols([]byte{})
	rvConnectionData.SetStationURLSpecialProtocols("")

	rmcResponseStream := nex.NewStreamOut(nexServer)

	// RESULT CODE
	rmcResponseStream.WriteUInt32LE(0x10001) // success

	// USER PID
	rmcResponseStream.WriteUInt32LE(uint32(userPID))

	// KERBEROS TICKET
	rmcResponseStream.WriteBuffer(encryptedTicket)

	// RVCONNECTION DATA
	rmcResponseStream.WriteStructure(rvConnectionData)

	// SERVER BUILD NAME
	rmcResponseStream.WriteString(serverName)

	rmcResponseBody := rmcResponseStream.Bytes()

	// Build response packet
	rmcResponse := nex.NewRMCResponse(nexproto.AuthenticationProtocolID, callID)
	rmcResponse.SetSuccess(nexproto.AuthenticationMethodLogin, rmcResponseBody)

	rmcResponseBytes := rmcResponse.Bytes()

	responsePacket, _ := nex.NewPacketV0(client, nil)

	responsePacket.SetVersion(0)
	responsePacket.SetSource(0xA1)
	responsePacket.SetDestination(0xAF)
	responsePacket.SetType(nex.DataPacket)
	responsePacket.SetPayload(rmcResponseBytes)

	responsePacket.AddFlag(nex.FlagNeedsAck)
	responsePacket.AddFlag(nex.FlagReliable)

	nexServer.Send(responsePacket)
}
