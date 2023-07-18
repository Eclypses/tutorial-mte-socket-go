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
	"fmt"
	"os"
	"strconv"
	"strings"

	//
	// Step 2 add new import items

	"tutorial-mte-socket-go/ecdh"
	mte "tutorial-mte-socket-go/mte"
)

var socketManager *ServerSocketManager

var serverEncoderInfo MteSetupInfo
var serverDecoderInfo MteSetupInfo

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

	fmt.Println("Starting Go Socket Server.")

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

	// Set default port.
	var connPort string = DEFAULT_PORT

	// Prompt for port from user -- use defaultPort if nothing entered
	reader := bufio.NewReader(os.Stdin)

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

	serverEncoderInfo = *NewMteSetupInfoDef()
	serverDecoderInfo = *NewMteSetupInfoDef()

	// Set up socket.
	socketManager = NewServerSocketManagerDef(connPort)

	// Exchange entropy, nonce, and personalization string between the client and server.
	if !ExchangeMteInfo() {
		fmt.Fprint(os.Stderr, "There was an error attempting to exchange information betweeen this and the client.")
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

	for {
		fmt.Println("Listening for messages from client...")

		// Receive and decode the message from the client.
		decoded, res := ReceiveAndDecodeMessage()
		if !res {
			break
		}

		// Encode and send the input.
		res = EncodeAndSendMessage(decoded)
		if !res {
			break
		}

		// Free the decoded message.
		decoded = nil
	}

	// Close the socket.
	socketManager.conn.Close()

	// This is the program exit point
	fmt.Println("Program stopped.")
	return
}

func ExchangeMteInfo() bool {
	// The client Encoder and the server Decoder will be paired.
	// The client Decoder and the server Encoder will be paired.

	// Loop until all 4 data are received from client, can be in any order.
	var recvCount int = 0
	var recvData RecvMsg
	for {
		if recvCount >= 4 {
			break
		}

		// Receive the next message from the client.
		recvData = socketManager.ReceiveMessage()

		// Evaluate the header.
		// 1 - server Decoder peer public key (from client Encoder)
		// 2 - server Decoder nonce (from client Encoder)
		// 3 - server Encoder peer public key (from client Decoder)
		// 4 - server Encoder nonce (from client Decoder)
		switch recvData.header {
		case '1':
			if serverDecoderInfo.GetPeerPublicKey() == nil {
				recvCount++
			}
			serverDecoderInfo.SetPeerPublicKey(recvData.message)
		case '2':
			if serverDecoderInfo.GetPersonalization() == nil {
				recvCount++
			}
			serverDecoderInfo.SetPersonalization(recvData.message)
		case '3':
			if serverEncoderInfo.GetPeerPublicKey() == nil {
				recvCount++
			}
			serverEncoderInfo.SetPeerPublicKey(recvData.message)
		case '4':
			if serverEncoderInfo.GetPersonalization() == nil {
				recvCount++
			}
			serverEncoderInfo.SetPersonalization(recvData.message)

		default:
			// Unknown message, abort here, send an 'E' for error.
			socketManager.SendMessage('E', []byte("ERR"))
		}
	}

	// Now all values from client have been received, send an 'A' for acknowledge to client.
	socketManager.SendMessage('A', []byte("ACK"))

	// Prepare to send server information now.

	// Create nonces.
	var minNonceBytes int = mte.GetDrbgsNonceMinBytes(encoder.GetDrbg())
	if minNonceBytes <= 0 {
		minNonceBytes = 1
	}

	res, serverEncoderNonce := ecdh.GetRandom(minNonceBytes)
	if res < 0 {
		return false
	}
	serverEncoderInfo.SetNonce(serverEncoderNonce)

	res, serverDecoderNonce := ecdh.GetRandom(minNonceBytes)
	if res < 0 {
		return false
	}
	serverDecoderInfo.SetNonce(serverDecoderNonce)

	// Send out information to the client.
	// 1 - server Encoder public key (to client Decoder)
	// 2 - server Encoder nonce (to client Decoder)
	// 3 - server Decoder public key (to client Encoder)
	// 4 - server Decoder nonce (to client Encoder)
	socketManager.SendMessage('1', serverEncoderInfo.GetPublicKey())
	socketManager.SendMessage('2', serverEncoderInfo.GetNonce())
	socketManager.SendMessage('3', serverDecoderInfo.GetPublicKey())
	socketManager.SendMessage('4', serverDecoderInfo.GetNonce())

	// Wait for ack from client.
	recvData = socketManager.ReceiveMessage()
	if recvData.header != 'A' {
		return false
	}

	recvData.message = nil

	return true
}

func CreateEncoder() bool {
	// Display all info related to the server Encoder.
	fmt.Printf("Server Encoder public key:\n %X\n", serverEncoderInfo.GetPublicKey())
	fmt.Printf("Server Encoder peer's key:\n %X\n", serverEncoderInfo.GetPeerPublicKey())
	fmt.Printf("Server Encoder nonce:\n %X\n", serverEncoderInfo.GetNonce())
	fmt.Printf("Server Encoder personalization:\n %v\n", string(serverEncoderInfo.GetPersonalization()))

	// Create shared secret.
	res, secret := serverEncoderInfo.ecdhManager.GetSharedSecret(serverEncoderInfo.GetPeerPublicKey())
	if res < 0 {
		fmt.Print("Failed to get shared secret.")
		return false
	}

	// Set Encoder entropy using this shared secret.
	encoder.SetEntropy(secret)

	// Set Encoder nonce.
	encoder.SetNonce(serverEncoderInfo.GetNonce())

	// Instantiate Encoder.
	var status = encoder.Instantiate(serverEncoderInfo.GetPersonalization())

	if status != mte.Status_mte_status_success {
		fmt.Fprintf(os.Stderr, "Encoder instantiate error %s: %s", mte.GetStatusName(status), mte.GetStatusDescription(status))
	}

	// Destroy server Encoder info.
	serverEncoderInfo.Destroy()

	return true
}

func CreateDecoder() bool {
	// Display all info related to the server Decoder.
	fmt.Printf("Server Decoder public key:\n %X\n", serverDecoderInfo.GetPublicKey())
	fmt.Printf("Server Decoder peer's key:\n %X\n", serverDecoderInfo.GetPeerPublicKey())
	fmt.Printf("Server Decoder nonce:\n %X\n", serverDecoderInfo.GetNonce())
	fmt.Printf("Server Decoder personalization:\n %v\n", string(serverDecoderInfo.GetPersonalization()))

	// Create shared secret.
	res, secret := serverDecoderInfo.ecdhManager.GetSharedSecret(serverDecoderInfo.GetPeerPublicKey())
	if res < 0 {
		fmt.Print("Failed to get shared secret.")
		return false
	}

	// Set Decoder entropy using this shared secret.
	decoder.SetEntropy(secret)

	// Set Decoder nonce.
	decoder.SetNonce(serverDecoderInfo.GetNonce())

	// Instantiate Decoder.
	var status = decoder.Instantiate(serverDecoderInfo.GetPersonalization())
	if status != mte.Status_mte_status_success {
		fmt.Fprintf(os.Stderr, "Decoder instantiate error %s: %s", mte.GetStatusName(status), mte.GetStatusDescription(status))
	}

	// Destroy server Decoder info.
	serverDecoderInfo.Destroy()

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
	// Receive and decode the message.
	decoded, status := ReceiveAndDecodeMessage()
	if !status {
		return false
	}

	// Check that it successfully decoded as "ping".
	if string(decoded) == "ping" {
		fmt.Println("Server Decoder decoded the message from the client Encoder successfully.")
	} else {
		fmt.Println("Server Decoder DID NOT decode the message from the client Encoder successfully.")
		return false
	}

	// Create the "ack" message.
	message := []byte("ack")

	// Encode and send message.
	if !EncodeAndSendMessage(message) {
		return false
	}

	return true
}
