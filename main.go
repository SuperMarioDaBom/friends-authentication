package main

import (
	"fmt"

	nex "github.com/PretendoNetwork/nex-go"
	nexproto "github.com/PretendoNetwork/nex-protocols-go"
)

var nexServer *nex.Server
var config *ServerConfig

func main() {
	nexServer = nex.NewServer()
	nexServer.SetPrudpVersion(config.PrudpVersion)
	nexServer.SetSignatureVersion(config.SignatureVersion)
	nexServer.SetKerberosKeySize(config.KerberosKeySize)
	nexServer.SetAccessKey(config.AccessKey)

	nexServer.On("Data", func(packet *nex.PacketV0) {
		request := packet.RMCRequest()

		fmt.Printf("==%s==\r\n", config.ServerName)
		fmt.Printf("Protocol ID: %#v\n", request.ProtocolID())
		fmt.Printf("Method ID: %#v\n", request.MethodID())
		fmt.Println("==================")
	})

	authenticationServer := nexproto.NewAuthenticationProtocol(nexServer)

	// Handle Login RMC method
	authenticationServer.Login(login)

	// Handle RequestTicket RMC method
	authenticationServer.RequestTicket(requestTicket)

	nexServer.Listen(fmt.Sprintf(":%s", config.ServerPort))
}
