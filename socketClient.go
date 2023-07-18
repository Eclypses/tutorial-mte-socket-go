/*
THIS SOFTWARE MAY NOT BE USED FOR PRODUCTION. Otherwise,
The MIT License (MIT)

Copyright (c) Eclypses, Inc.

All rights reserved.

Permission is hereby granted, free of charge, to any person obtaining a copy
of this software and associated documentation files (the "Software"), to deal
in the Software without restriction, including without limitation the rights
to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
copies of the Software, and to permit persons to whom the Software is
furnished to do so, subject to the following conditions:

The above copyright notice and this permission notice shall be included in
all copies or substantial portions of the Software.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
SOFTWARE.
*/
package main

import (
	"bufio"
	"crypto/rand"
	"fmt"
	"os"
	"strconv"
	"strings"

	//
	// Step 2 add new import items

	mte "tutorial-mte-socket-go/mte"
)

var socketManager *ClientSocketManager = NewClientSocketManagerDef("", "")

var clientEncoderInfo MteSetupInfo
var clientDecoderInfo MteSetupInfo

// ---------------------------------------------------
// MKE and Fixed length add-ons are NOT in all SDK
// MTE versions. If the name of the SDK includes
// "-MKE" then it will contain the MKE add-on. If the
// name of the SDK includes "-FLEN" then it contains
// the Fixed length add-on.
// ---------------------------------------------------
// ---------------------------------------------------
// Uncomment to use MTE core Encoder and Decoder
// ---------------------------------------------------
var encoder *mte.Enc
var decoder *mte.Dec

//---------------------------------------------------
// Uncomment to use MKE Encoder and Decoder
//---------------------------------------------------
//var encoder *mte.MkeEnc
//var decoder *mte.MkeDec
//---------------------------------------------------
// Uncomment to use Flen Encoder and Decoder
//---------------------------------------------------
//var encoder *mte.FlenEnc
//var decoder *mte.Dec

func main() {
	//
	// This tutorial uses Sockets for communication.
	// It should be noted that the MTE can be used with any type of communication. (SOCKETS are not required!)
	//

	fmt.Println("Starting Go Socket Client.")

	// Step 6

	//---------------------------------------------------
	// Uncomment to use MTE core Encoder and Decoder
	//---------------------------------------------------
	encoder = mte.NewEncDef()
	decoder = mte.NewDecDef()
	var mteType = "Core"
	//---------------------------------------------------
	// Uncomment to use MKE Encoder and Decoder
	//---------------------------------------------------
	//encoder = mte.NewMkeEncDef()
	//decoder = mte.NewMkeDecDef()
	//var mteType = "MKE"
	//---------------------------------------------------
	// Uncomment to use Flen Encoder and Decoder
	//---------------------------------------------------
	//encoder = mte.NewFlenEncDef(MAX_INPUT_BYTES)
	//decoder = mte.NewDecDef()
	// var mteType = "FLEN"
	defer encoder.Destroy()
	defer decoder.Destroy()

	// defer the exit so all other defer calls are called
	retcode := 0
	defer func() { os.Exit(retcode) }()

	// Step 4
	// (optional) Check and print out mte Version
	mteVersion := mte.GetVersion()
	fmt.Printf("Using Mte Version %s-%s\n", mteVersion, mteType)

	// Step 5
	// Initialize MTE license. If trial version or non-licensed version
	// any value may be used for the companyName and companyLicense
	if !mte.InitLicense("LicenseCompanyName", "LicenseKey") {
		fmt.Println("There was an error attempting to initialize the MTE License.")
		return
	}

	// Set default port and server IP address.
	var connPort string = DEFAULT_PORT
	var ipAddress = DEFAULT_SERVER_IP

	// Prompt for port from user -- use defaultPort if nothing entered
	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Please enter ip address of Server, press Enter to use default: " + DEFAULT_SERVER_IP)
	ipAddress, _ = reader.ReadString('\n')
	// Trim off the carriage return
	ipAddress = strings.TrimSuffix(ipAddress, "\n")
	ipAddress = strings.TrimSuffix(ipAddress, "\r")

	if ipAddress == "" {
		ipAddress = DEFAULT_SERVER_IP
	}

	fmt.Println("Server is at " + ipAddress)

	isValidPort := false
	fmt.Println("Please enter port to use, press Enter to use default: " + DEFAULT_PORT)

	// Create loop to get port, make sure it is valid
	for !isValidPort {
		connPort, _ = reader.ReadString('\n')

		// Trim off the carriage return
		connPort = strings.TrimSuffix(connPort, "\n")
		connPort = strings.TrimSuffix(connPort, "\r")

		if connPort == "" {
			connPort = DEFAULT_PORT
		}

		// Make sure port is an integer
		if _, err := strconv.ParseInt(connPort, 10, 64); err == nil {
			isValidPort = true
		} else {
			fmt.Println("Port is not valid, please enter valid port number")
		}
	}

	clientEncoderInfo = *NewMteSetupInfoDef()
	clientDecoderInfo = *NewMteSetupInfoDef()

	// Set up socket.
	socketManager = NewClientSocketManagerDef(ipAddress, connPort)

	// Exchange entropy, nonce, and personalization string between the client and server.
	if !ExchangeMteInfo() {
		fmt.Fprint(os.Stderr, "There was an error attempting to exchange information betweeen this and the server.")
		return
	}

	// Create the Encoder.
	if !CreateEncoder() {
		fmt.Fprint(os.Stderr, "There was a problem creating the Encoder.")
		return
	}

	// Create the Decoder.
	if !CreateDecoder() {
		fmt.Fprint(os.Stderr, "There was a problem creating the Decoder.")
		return
	}

	// Run the diagnostic test.
	if !RunDiagnosticTest() {
		fmt.Fprint(os.Stderr, "There was a problem running the diagnostic test.")
		return
	}

	// Run loop forever, until exit.
	for {
		// Prompting message.
		fmt.Printf("Please enter text up to %v bytes to send: (To end please type 'quit')", MAX_INPUT_BYTES)

		// Read in input until newline, Enter key.
		input, _ := reader.ReadString('\n')

		// If entered text is quit we want to close client
		input = strings.TrimSuffix(input, "\n")
		input = strings.TrimSuffix(input, "\r")

		if strings.ToLower(input) == "quit" {
			fmt.Println("Program stopped.")
			retcode = 0
			return
		}

		// Encode and send the input.
		if !EncodeAndSendMessage([]byte(input)) {
			break
		}

		// Receive and decode the returned data.
		decoded, res := ReceiveAndDecodeMessage()
		if !res {
			break
		}

		// Compare the decoded message to the original.
		if string(decoded) == input {
			fmt.Println("The original input and decoded return match.")
		} else {
			fmt.Println("The original input and decoded return DO NOT match.")
			return
		}
	}

	// Close the socket.
	socketManager.conn.Close()
}

func ExchangeMteInfo() bool {
	// The client Encoder and the server Decoder will be paired.
	// The client Decoder and the server Encoder will be paired.

	// Prepare to send client information.

	// Create personalization strings.
	clientEncoderInfo.SetPersonalization(CreateGuid())
	clientDecoderInfo.SetPersonalization(CreateGuid())

	// Send out information to the server.
	// 1 - client Encoder public key (to server Decoder)
	// 2 - client Encoder personalization string (to server Decoder)
	// 3 - client Decoder public key (to server Encoder)
	// 4 - client Decoder personalization string (to server Encoder)
	socketManager.SendMessage('1', clientEncoderInfo.GetPublicKey())
	socketManager.SendMessage('2', clientEncoderInfo.GetPersonalization())
	socketManager.SendMessage('3', clientDecoderInfo.GetPublicKey())
	socketManager.SendMessage('4', clientDecoderInfo.GetPersonalization())

	// Wait for ack from server.
	var recvData RecvMsg = socketManager.ReceiveMessage()
	if recvData.header != 'A' {

		return false
	}

	recvData.message = nil

	// Processing incoming message all 4 will be needed.
	var recvCount int = 0
	for {
		if recvCount >= 4 {
			break
		}

		// Receive the next message from the server.
		recvData = socketManager.ReceiveMessage()

		// Evaluate the header.
		// 1 - client Decoder peer public key (from server Encoder)
		// 2 - client Decoder nonce (from server Encoder)
		// 3 - client Encoder peer public key (from server Decoder)
		// 4 - client Encoder nonce (from server Decoder)
		switch recvData.header {
		case '1':
			if clientDecoderInfo.GetPeerPublicKey() == nil {
				recvCount++
			}
			clientDecoderInfo.SetPeerPublicKey(recvData.message)
		case '2':
			if clientDecoderInfo.GetNonce() == nil {
				recvCount++
			}
			clientDecoderInfo.SetNonce(recvData.message)
		case '3':
			if clientEncoderInfo.GetPeerPublicKey() == nil {
				recvCount++
			}
			clientEncoderInfo.SetPeerPublicKey(recvData.message)
		case '4':
			if clientEncoderInfo.GetNonce() == nil {
				recvCount++
			}
			clientEncoderInfo.SetNonce(recvData.message)

		default:
			// Unknown message, abort here, send an 'E' for error.
			socketManager.SendMessage('E', []byte("ERR"))
		}
	}

	// Now all values from server have been received, send an 'A' for acknowledge to server.
	socketManager.SendMessage('A', []byte("ACK"))

	return true
}

func CreateGuid() []byte {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		fmt.Print(err)
	}

	return []byte(fmt.Sprintf("%x-%x-%x-%x-%x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:]))
}

func CreateEncoder() bool {
	// Display all info related to the client Encoder.
	fmt.Printf("Client Encoder public key:\n %X\n", clientEncoderInfo.GetPublicKey())
	fmt.Printf("Client Encoder peer's key:\n %X\n", clientEncoderInfo.GetPeerPublicKey())
	fmt.Printf("Client Encoder nonce:\n %X\n", clientEncoderInfo.GetNonce())
	fmt.Printf("Client Encoder personalization:\n %v\n", string(clientEncoderInfo.GetPersonalization()))

	// Create shared secret.
	res, secret := clientEncoderInfo.ecdhManager.GetSharedSecret(clientEncoderInfo.GetPeerPublicKey())
	if res < 0 {
		fmt.Print("Failed to get shared secret.")
		return false
	}

	// Set Encoder entropy using this shared secret.
	encoder.SetEntropy(secret)

	// Set Encoder nonce.
	encoder.SetNonce(clientEncoderInfo.GetNonce())

	// Instantiate Encoder.
	var status = encoder.Instantiate(clientEncoderInfo.GetPersonalization())

	if status != mte.Status_mte_status_success {
		fmt.Fprintf(os.Stderr, "Encoder instantiate error %s: %s", mte.GetStatusName(status), mte.GetStatusDescription(status))
	}

	// Destroy client Encoder info.
	clientEncoderInfo.Destroy()

	return true
}

func CreateDecoder() bool {
	// Display all info related to the client Decoder.
	fmt.Printf("Client Decoder public key:\n %X\n", clientDecoderInfo.GetPublicKey())
	fmt.Printf("Client Decoder peer's key:\n %X\n", clientDecoderInfo.GetPeerPublicKey())
	fmt.Printf("Client Decoder nonce:\n %X\n", clientDecoderInfo.GetNonce())
	fmt.Printf("Client Decoder personalization:\n %v\n", string(clientDecoderInfo.GetPersonalization()))

	// Create shared secret.
	res, secret := clientDecoderInfo.ecdhManager.GetSharedSecret(clientDecoderInfo.GetPeerPublicKey())
	if res < 0 {
		fmt.Print("Failed to get shared secret.")
		return false
	}

	// Set Decoder entropy using this shared secret.
	decoder.SetEntropy(secret)

	// Set Decoder nonce.
	decoder.SetNonce(clientDecoderInfo.GetNonce())

	// Instantiate Decoder.
	var status = decoder.Instantiate(clientDecoderInfo.GetPersonalization())
	if status != mte.Status_mte_status_success {
		fmt.Fprintf(os.Stderr, "Decoder instantiate error %s: %s", mte.GetStatusName(status), mte.GetStatusDescription(status))
	}

	// Destroy client Decoder info.
	clientDecoderInfo.Destroy()

	return true
}

func EncodeAndSendMessage(message []byte) bool {
	// Display original message.
	fmt.Printf("Message to be encoded:\n %v\n", string(message))

	// Encode the message.
	encoded, status := encoder.Encode(message)
	if status != mte.Status_mte_status_success {
		fmt.Fprintf(os.Stderr, "Error encoding %s: %s", mte.GetStatusName(status), mte.GetStatusDescription(status))
		return false
	}

	// Send the encoded message.
	res := socketManager.SendMessage('m', encoded)
	if res <= 0 {
		return false
	}

	// Display encoded message.
	fmt.Printf("Encoded message being sent:\n %X\n", encoded)

	return true
}

func ReceiveAndDecodeMessage() ([]byte, bool) {
	// Wait for return message.
	msgStruct := socketManager.ReceiveMessage()

	if !msgStruct.success || msgStruct.message == nil || msgStruct.header != 'm' {
		return nil, false
	}

	// Display encoded message.
	fmt.Printf("Encoded message received:\n %X\n", msgStruct.message)

	// Decode the message.
	message, status := decoder.Decode(msgStruct.message)
	if mte.StatusIsError(status) {
		fmt.Fprintf(os.Stderr, "Error decoding %s: %s", mte.GetStatusName(status), mte.GetStatusDescription(status))
		return nil, false
	}

	// Display decoded message.
	fmt.Printf("Decoded message: \n %v\n", string(message))
	return message, true
}

func RunDiagnosticTest() bool {
	// Create ping message.
	message := []byte("ping")

	// Encode and send message.
	if !EncodeAndSendMessage(message) {
		return false
	}

	// Receive and decode the message.
	decoded, res := ReceiveAndDecodeMessage()
	if !res {
		return false
	}

	// Check that it successfully decoded as "ack".
	if string(decoded) == "ack" {
		fmt.Println("Client Decoder decoded the message from the server Encoder successfully.")
	} else {
		fmt.Println("Client Decoder DID NOT decode the message from the server Encoder successfully.")
		return false
	}
	return true
}
