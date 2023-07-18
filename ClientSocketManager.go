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
	"encoding/binary"
	"fmt"
	"io"
	"net"
)

type ClientSocketManager struct {
	conn net.Conn
}

type RecvMsg struct {
	success bool
	header  byte
	message []byte
}

func NewClientSocketManagerDef(ipAddress string, port string) *ClientSocketManager {
	// Start the client and connect to the server.
	connection, err := net.Dial("tcp", ipAddress+":"+port)
	if err != nil {
		fmt.Println("Error connecting:", err.Error())
		return nil
	}

	clientSocketManager := ClientSocketManager{conn: connection}

	fmt.Println("Client connected to Server.")

	return &clientSocketManager
}

func (sock *ClientSocketManager) SendMessage(header byte, message []byte) int {
	// Get the length of the message.
	var toSendLen = len(message)

	if toSendLen == 0 || header == 0x0 {
		fmt.Print("Unable to send message.")
		sock.conn.Close()
		return 0
	}

	// Set the size of the message to Big Endian.
	toSendLenBytes := make([]byte, 4)
	binary.BigEndian.PutUint32(toSendLenBytes, uint32(toSendLen))

	// Send the message size as big-endian.
	res, err := sock.conn.Write(toSendLenBytes)
	if err != nil {
		fmt.Printf("Error sending length %v %v", err, res)
		sock.conn.Close()
		return 0
	}

	// Send the header.
	res, err = sock.conn.Write([]byte{header})
	if err != nil {
		fmt.Printf("Error sending header %v %v", err, res)
		sock.conn.Close()
		return 0
	}

	// Send the actual message.
	res, err = sock.conn.Write(message)
	if err != nil {
		fmt.Printf("Error sending message %v %v", err, res)
		sock.conn.Close()
		return 0
	}

	return toSendLen
}

func (sock *ClientSocketManager) ReceiveMessage() RecvMsg {
	// Create RecvMsg struct.
	var msgStruct RecvMsg
	msgStruct.success = false
	msgStruct.header = 0x00
	msgStruct.message = nil

	// Create array to hold the message size coming in.
	rcvLenBytes := make([]byte, 4)
	res, err := sock.conn.Read(rcvLenBytes)
	if err != nil {
		fmt.Print("Unable to receive message.")
		fmt.Print(err, res)
		sock.conn.Close()
		return msgStruct
	}

	rcvLen := binary.BigEndian.Uint32(rcvLenBytes)

	// Get the header.
	var headerByte = make([]byte, 1)
	res, err = sock.conn.Read(headerByte)
	if err != nil {
		fmt.Print("Unable to receive message.")
		fmt.Print(err, res)
		sock.conn.Close()
		return msgStruct
	}
	msgStruct.header = headerByte[0]

	// Receive the message from the other side.
	msgStruct.message = make([]byte, rcvLen)

	_, err = io.ReadFull(sock.conn, msgStruct.message)
	if err != nil {
		fmt.Println(err.Error())
		sock.conn.Close()
		return msgStruct
	}

	// Set status to true.
	msgStruct.success = true

	return msgStruct
}
