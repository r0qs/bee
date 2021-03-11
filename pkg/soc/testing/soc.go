// Copyright 2021 The Swarm Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

package testing

import (
	"math/rand"
	"time"

	"github.com/ethersphere/bee/pkg/cac"
	"github.com/ethersphere/bee/pkg/crypto"
	"github.com/ethersphere/bee/pkg/soc"
	"github.com/ethersphere/bee/pkg/swarm"
)

func init() {
	rand.Seed(time.Now().UnixNano())
}

// MockSoc defines a mocked soc with exported fields for easy testing.
type MockSoc struct {
	ID           soc.ID
	Owner        soc.Owner
	Signature    []byte
	WrappedChunk swarm.Chunk
}

// Address returns the soc address of the mocked soc.
func (ms MockSoc) Address() swarm.Address {
	addr, _ := soc.CreateAddress(ms.ID, ms.Owner)
	return addr
}

// Chunk returns the soc chunk of the mocked soc.
func (ms MockSoc) Chunk() swarm.Chunk {
	return swarm.NewChunk(ms.Address(), append(ms.ID, append(ms.Signature, ms.WrappedChunk.Data()...)...))
}

// GenerateMockSoc generates a valid mocked soc from given data.
// If data is nil it generates random data.
func GenerateMockSoc(data []byte) *MockSoc {
	privKey, _ := crypto.GenerateSecp256k1Key()
	signer := crypto.NewDefaultSigner(privKey)
	owner, _ := signer.EthereumAddress()

	if data == nil {
		data = make([]byte, swarm.ChunkSize)
		_, _ = rand.Read(data)
	}
	ch, _ := cac.New(data)

	id := make([]byte, 32)
	hasher := swarm.NewHasher()
	_, _ = hasher.Write(append(id, ch.Address().Bytes()...))
	signature, _ := signer.Sign(hasher.Sum(nil))
	return &MockSoc{
		ID:           id,
		Owner:        owner.Bytes(),
		Signature:    signature,
		WrappedChunk: ch,
	}
}