package main

import (
	"fmt"

	nex "github.com/PretendoNetwork/nex-go"
	nexproto "github.com/PretendoNetwork/nex-protocols-go"
)

var nexServer *nex.Server

func main() {
	nexServer = nex.NewServer()
	nexServer.SetPrudpVersion(0)
	nexServer.SetSignatureVersion(1)
	nexServer.SetKerberosKeySize(16)
	nexServer.SetAccessKey("ridfebb9")

	nexServer.On("Data", func(packet *nex.PacketV0) {
		request := packet.RMCRequest()

		fmt.Println("==Friends - Auth==")
		fmt.Printf("Protocol ID: %#v\n", request.ProtocolID())
		fmt.Printf("Method ID: %#v\n", request.MethodID())
		fmt.Println("==================")
	})

	authenticationServer := nexproto.NewAuthenticationProtocol(nexServer)

	// Handle Login RMC method
	authenticationServer.Login(login)

	// Handle RequestTicket RMC method
	authenticationServer.RequestTicket(requestTicket)

	nexServer.Listen(":60000")
}
