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
	"tutorial-mte-socket-go/ecdh"
)

type MteSetupInfo struct {
	ecdhManager     *ecdh.EcdhP256
	publicKey       []byte
	personalization []byte
	nonce           []byte
	peerKey         []byte
}

// NewMteSetupInfoDef creates MteSetupInfo.
func NewMteSetupInfoDef() *MteSetupInfo {
	mteSetupInfo := new(MteSetupInfo)

	mteSetupInfo.ecdhManager = ecdh.NewEcdhP256Def()

	// Create the private and public keys.
	var res int = 0
	res, mteSetupInfo.publicKey = mteSetupInfo.ecdhManager.CreateKeypair()
	if res < 0 {
		return nil
	}

	return mteSetupInfo
}

// Destroy releases resources.
func (mteSetupInfo *MteSetupInfo) Destroy() {
	// Free the buffers.
	ecdh.Zeroize(mteSetupInfo.publicKey)
}

func (mteSetupInfo *MteSetupInfo) GetPublicKey() []byte {
	// Create temp byte array.
	temp := mteSetupInfo.publicKey
	return temp
}

func (mteSetupInfo *MteSetupInfo) GetSharedSecret() []byte {
	if len(mteSetupInfo.peerKey) == 0 {
		return nil
	}

	// Create shared secret.
	res, temp := mteSetupInfo.ecdhManager.GetSharedSecret(mteSetupInfo.peerKey)
	if res < 0 {
		return nil
	}
	return temp
}

func (mteSetupInfo *MteSetupInfo) SetPersonalization(data []byte) {
	// Check if personalization already set.
	if mteSetupInfo.personalization != nil {
		mteSetupInfo.personalization = nil
	}

	// Create personalization.
	mteSetupInfo.personalization = data
}

func (mteSetupInfo *MteSetupInfo) GetPersonalization() []byte {

	if mteSetupInfo.personalization == nil {
		return nil
	} else {
		temp := mteSetupInfo.personalization
		return temp
	}
}

func (mteSetupInfo *MteSetupInfo) SetNonce(data []byte) {
	// Check if nonce already set.
	if mteSetupInfo.nonce != nil {
		mteSetupInfo.nonce = nil
	}

	// Create nonce.
	mteSetupInfo.nonce = data
}

func (mteSetupInfo *MteSetupInfo) GetNonce() []byte {

	if mteSetupInfo.nonce == nil {
		return nil
	} else {
		temp := mteSetupInfo.nonce
		return temp
	}
}

func (mteSetupInfo *MteSetupInfo) SetPeerPublicKey(data []byte) {
	// Check if peer key already set.
	if mteSetupInfo.peerKey != nil {
		mteSetupInfo.peerKey = nil
	}

	// Create peer key.
	mteSetupInfo.peerKey = data
}

func (mteSetupInfo *MteSetupInfo) GetPeerPublicKey() []byte {

	if mteSetupInfo.peerKey == nil {
		return nil
	} else {
		temp := mteSetupInfo.peerKey
		return temp
	}
}
