

<img src="Eclypses.png" style="width:50%;margin-right:0;"/>

<div align="center" style="font-size:40pt; font-weight:900; font-family:arial; margin-top:300px; " >
Go Socket Tutorial</div>
<br>
<div align="center" style="font-size:28pt; font-family:arial; " >
MTE Implementation Tutorial (MTE Core, MKE, MTE Fixed Length)</div>
<br>
<div align="center" style="font-size:15pt; font-family:arial; " >
Using MTE version 3.1.x</div>





[Introduction](#introduction)

[Socket Tutorial Server and Client](#socket-tutorial-server-and-client)


<div style="page-break-after: always; break-after: page;"></div>

# Introduction

This tutorial is sending messages via a socket connection. This is only a sample, the MTE does NOT require the usage of sockets, you can use whatever communication protocol that is needed.

This tutorial demonstrates how to use Mte Core, Mte MKE and Mte Fixed Length. For this application, only one type can be used at a time; however, it is possible to implement any and all at the same time depending on needs.

This tutorial contains two main programs, a client and a server, which may be run on Windows and Linux. It is required that the server is up and running awaiting for a client connection first before the client can be started; The client program will error out if it starts without being able to connect to a server application. Note that any of the available languages can be used for any available platform as long as communication is possible. It is just recommended that a server program is started first and then a client program can be started.

The MTE Encoder and Decoder need several pieces of information to be the same in order to function properly. This includes entropy, nonce, and personalization. If this information must be shared, the entropy MUST be passed securely. One way to do this is with a Diffie-Hellman approach. Each side will then be able to create two shared secrets to use as entropy for each pair of Encoder/Decoder. The two personalization values will be created by the client and shared to the other side. The two nonce values will be created by the server and shared.

The SDK that you received from Eclypses may not include the MKE or MTE FLEN add-ons. If your SDK contains either the MKE or the Fixed Length add-ons, the name of the SDK will contain "-MKE" or "-FLEN". If these add-ons are not there and you need them please work with your sales associate. If there is no need, please just ignore the MKE and FLEN options.

Here is a short explanation of when to use each, but it is encouraged to either speak to a sales associate or read the dev guide if you have additional concerns or questions.

***MTE Core:*** This is the recommended version of the MTE to use. Unless payloads are large or sequencing is needed this is the recommended version of the MTE and the most secure.

***MTE MKE:*** This version of the MTE is recommended when payloads are very large, the MTE Core would, depending on the token byte size, be multiple times larger than the original payload. Because this uses the MTE technology on encryption keys and encrypts the payload, the payload is only enlarged minimally.

***MTE Fixed Length:*** This version of the MTE is very secure and is used when the resulting payload is desired to be the same size for every transmission. The Fixed Length add-on is mainly used when using the sequencing verifier with MTE. In order to skip dropped packets or handle asynchronous packets the sequencing verifier requires that all packets be a predictable size. If you do not wish to handle this with your application then the Fixed Length add-on is a great choice. This is ONLY an encoder change - the decoder that is used is the MTE Core decoder.

***IMPORTANT NOTE***
>If using the fixed length MTE (FLEN), all messages that are sent that are longer than the set fixed length will be trimmed by the MTE. The other side of the MTE will NOT contain the trimmed portion. Also messages that are shorter than the fixed length will be padded by the MTE so each message that is sent will ALWAYS be the same length. When shorter message are "decoded" on the other side the MTE takes off the extra padding when using strings and hands back the original shorter message, BUT if you use the raw interface the padding will be present as all zeros. Please see official MTE Documentation for more information.

In this tutorial, there is an MTE Encoder on the client that is paired with an MTE Decoder on the server. Likewise, there is an MTE Encoder on the server that is paired with an MTE Decoder on the client. Secured messages wil be sent to and from both sides. If a system only needs to secure messages one way, only one pair could be used.

**IMPORTANT**
>Please note the solution provided in this tutorial does NOT include the MTE library or supporting MTE library files. If you have NOT been provided an MTE library and supporting files, please contact Eclypses Inc. The solution will only work AFTER the MTE library and MTE library files have been incorporated.
  

# Socket Tutorial Server and Client

## MTE Directory and File Setup
<ol>
<li>
Navigate to the "tutorial-mte-socket-go" directory.
</li>
<li>
Create a directory named "MTE". This will contain all needed MTE files.
</li>
<li>
Copy the "lib" directory and contents from the MTE SDK into the "MTE" directory.
</li>
<li>
Copy the "include" directory and contents from the MTE SDK into the "MTE" directory.
</li>
<li>
Copy the "src/go/mte" directory and contents from the MTE SDK into the "MTE" directory.
</li>
</ol>

## ECDH Directory and File Setup
<ol>
<li>
Navigate to the "tutorial-mte-socket-go" directory.
</li>
<li>
Create a directory named "ecdh". This will contain all needed ecdh files.
</li>
<li>
Copy the "lib" directory and contents from the ecdh SDK into the "ecdh" directory.
</li>
<li>
Copy the "include" directory and contents from the ecdh SDK into the "ecdh" directory.
</li>

<li>
Copy the "src/go/ecdh" directory and contents from the ecdh SDK into the "ecdh" directory.
</li>
</ol>

## Source Code Key Points

### MTE Setup

<ol>
<li>
Comment/uncomment various code sections to more easily handle the function calls for the MTE Core or the add-on configurations. In the files "socketClinet.go" and "socketServer.go", there are two sections that will need to be considered.
<ul>

<li>

```go
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
```
</li>
<li>

```go
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
```
</li>
</ul>

</li>

<li>
In this application, the Eclypses Elliptic Curve Diffie-Hellman (ECDH) support package is used to create entropy public and private keys. The public keys are then shared between the client and server, and then shared secrets are created to use as matching entropy for the creation of the Encoders and Decoders. The nonces are also created using the randomization feature of the support package.

```go
// Create the private and public keys.
var res int = 0
res, mteSetupInfo.publicKey = mteSetupInfo.ecdhManager.CreateKeypair()
if res < 0 {
	return nil
}
```
The Go ECDHP256 class will keep the private key to itself and not provide access to the calling application.
</li>
<li>
The public keys created by the client will be sent to the server, and vice versa, and will be received as <i>peer public keys</i>. Then the shared secret can be created on each side. These should match as long as the information has been created and shared correctly.

```go
// Create shared secret.
res, temp := mteSetupInfo.ecdhManager.GetSharedSecret(mteSetupInfo.peerKey)
if res < 0 {
	return nil
}
return temp
```
These secrets will then be used to fufill the entropy needed for the Encoders and Decoders.
</li>
<li>
The client will create the personalization strings, in this case a guid-like structure using the rand library.

```go
func CreateGuid() []byte {
	b := make([]byte, 16)
	_, err := rand.Read(b)
	if err != nil {
		fmt.Print(err)
	}

	return []byte(fmt.Sprintf("%x-%x-%x-%x-%x",
		b[0:4], b[4:6], b[6:8], b[8:10], b[10:]))
}
```
</li>
<li>
The two public keys and the two personalization strings will then be sent to the server. The client will wait for an awknowledgment.

```go
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
```
</li>
<li>
The server will wait for the two public keys and the two personalization strings from the client. Once all four pieces of information have been received, it will send an awknowledgment.

```go
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
```
</li>
<li>
The server will create the private and public keypairs, one for the server Encoder and client Decoder, and one for the server Decoder and client Encoder. The server uses the same file "MTESeuptInfo.go" 

```go
// Create the private and public keys.
var res int = 0
res, mteSetupInfo.publicKey = mteSetupInfo.ecdhManager.CreateKeypair()
if res < 0 {
	return nil
}
```

</li>
<li>
The server will create the nonces, using the platform supplied secure RNG.

```go
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
```
</li>
<li>
The two public keys and the two nonces will then be sent to the client. The server will wait for an awknowledgment. 
```go
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
```
</li>

<li>
The client will now wait for information from the server. This includes the two server public keys, and the two nonces. Once all pieces of information have been obtained, the client will send an awknowledgment back to the server.

```go
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
```

</li>
<li>
After the client and server have exchanged their information, the client and server can each create their respective Encoder and Decoder. This is where the personalization string and nonce will be added. Additionally, the entropy will be set by getting the shared secret from ECDH. This sample code showcases the client Encoder. There will be four of each of these that will be very similar. Ensure carefully that each function uses the appropriate client/server, and Encoder/Decoder variables and functions.

```go
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
```

</li>
</ol>

### Diagnostic Test
<ol>
<li>
The application will run a diagnostic test, where the client will encode the word "ping", then send the encoded message to the server. The server will decode the received message to confirm that the original message is "ping". Then the server will encode the word "ack" and send the encoded message to the client. The client then decodes the received message, and confirms that it decodes it to the word "ack". 
</li>
</ol>

### User Interaction
<ol>
<li>
The application will continously prompt the user for an input (until the user types "quit"). That input will be encoded with the client Encoder and sent to the server.

```go
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
```
</li>
<li>
The server will use its Decoder to decode the message.

```go
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
```

</li>
<li>
Then that message will be re-encoded with the server Encoder and sent to the client.The client Decoder will then decode that message, which then will be compared with the original user input.
</li>
</ol>

<div style="page-break-after: always; break-after: page;"></div>

# Contact Eclypses

<img src="Eclypses.png" style="width:8in;"/>

<p align="center" style="font-weight: bold; font-size: 20pt;">Email: <a href="mailto:info@eclypses.com">info@eclypses.com</a></p>
<p align="center" style="font-weight: bold; font-size: 20pt;">Web: <a href="https://www.eclypses.com">www.eclypses.com</a></p>
<p align="center" style="font-weight: bold; font-size: 20pt;">Chat with us: <a href="https://developers.eclypses.com/dashboard">Developer Portal</a></p>
<p style="font-size: 8pt; margin-bottom: 0; margin: 300px 24px 30px 24px; " >

<b>All trademarks of Eclypses Inc.</b> may not be used without Eclypses Inc.'s prior written consent. No license for any use thereof has been granted without express written consent. Any unauthorized use thereof may violate copyright laws, trademark laws, privacy and publicity laws and communications regulations and statutes. The names, images and likeness of the Eclypses logo, along with all representations thereof, are valuable intellectual property assets of Eclypses, Inc. Accordingly, no party or parties, without the prior written consent of Eclypses, Inc., (which may be withheld in Eclypses' sole discretion), use or permit the use of any of the Eclypses trademarked names or logos of Eclypses, Inc. for any purpose other than as part of the address for the Premises, or use or permit the use of, for any purpose whatsoever, any image or rendering of, or any design based on, the exterior appearance or profile of the Eclypses trademarks and or logo(s).
</p>