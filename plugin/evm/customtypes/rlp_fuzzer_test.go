// Copyright (C) 2019-2025, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.
//
// This file is a derived work, based on the go-ethereum library whose original
// notices appear below.
//
// It is distributed under a license compatible with the licensing terms of the
// original code from which it is derived.
//
// Much love to the original authors for their work.
// **********
// Copyright 2019 The go-ethereum Authors
// This file is part of the go-ethereum library.
//
// The go-ethereum library is free software: you can redistribute it and/or modify
// it under the terms of the GNU Lesser General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// The go-ethereum library is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU Lesser General Public License for more details.
//
// You should have received a copy of the GNU Lesser General Public License
// along with the go-ethereum library. If not, see <http://www.gnu.org/licenses/>.

package customtypes

import (
	"bytes"
	"fmt"
	"math/big"
	"testing"

	"github.com/ava-labs/libevm/rlp"
	"github.com/holiman/uint256"

	// TODO(arr4n) These tests were originally part of the `subnet-evm/core/types`
	// package so assume the presence of identifiers. A dot-import reduces PR
	// noise during the refactoring.
	. "github.com/ava-labs/libevm/core/types"
)

func decodeEncode(input []byte, val interface{}) error {
	if err := rlp.DecodeBytes(input, val); err != nil {
		// not valid rlp, nothing to do
		return nil
	}
	// If it _were_ valid rlp, we can encode it again
	output, err := rlp.EncodeToBytes(val)
	if err != nil {
		return err
	}
	if !bytes.Equal(input, output) {
		return fmt.Errorf("encode-decode is not equal, \ninput : %x\noutput: %x", input, output)
	}
	return nil
}

func FuzzRLP(f *testing.F) {
	f.Fuzz(fuzzRlp)
}

func fuzzRlp(t *testing.T, input []byte) {
	if len(input) == 0 || len(input) > 500*1024 {
		return
	}
	rlp.Split(input)
	if elems, _, err := rlp.SplitList(input); err == nil {
		rlp.CountValues(elems)
	}
	rlp.NewStream(bytes.NewReader(input), 0).Decode(new(interface{}))
	if err := decodeEncode(input, new(interface{})); err != nil {
		t.Fatal(err)
	}
	{
		var v struct {
			Int    uint
			String string
			Bytes  []byte
		}
		if err := decodeEncode(input, &v); err != nil {
			t.Fatal(err)
		}
	}
	{
		type Types struct {
			Bool  bool
			Raw   rlp.RawValue
			Slice []*Types
			Iface []interface{}
		}
		var v Types
		if err := decodeEncode(input, &v); err != nil {
			t.Fatal(err)
		}
	}
	{
		type AllTypes struct {
			Int    uint
			String string
			Bytes  []byte
			Bool   bool
			Raw    rlp.RawValue
			Slice  []*AllTypes
			Array  [3]*AllTypes
			Iface  []interface{}
		}
		var v AllTypes
		if err := decodeEncode(input, &v); err != nil {
			t.Fatal(err)
		}
	}
	{
		if err := decodeEncode(input, [10]byte{}); err != nil {
			t.Fatal(err)
		}
	}
	{
		var v struct {
			Byte [10]byte
			Rool [10]bool
		}
		if err := decodeEncode(input, &v); err != nil {
			t.Fatal(err)
		}
	}
	{
		var h Header
		if err := decodeEncode(input, &h); err != nil {
			t.Fatal(err)
		}
		var b Block
		if err := decodeEncode(input, &b); err != nil {
			t.Fatal(err)
		}
		var tx Transaction
		if err := decodeEncode(input, &tx); err != nil {
			t.Fatal(err)
		}
		var txs Transactions
		if err := decodeEncode(input, &txs); err != nil {
			t.Fatal(err)
		}
		var rs Receipts
		if err := decodeEncode(input, &rs); err != nil {
			t.Fatal(err)
		}
	}
	{
		var v struct {
			AnIntPtr  *big.Int
			AnInt     big.Int
			AnU256Ptr *uint256.Int
			AnU256    uint256.Int
			NotAnU256 [4]uint64
		}
		if err := decodeEncode(input, &v); err != nil {
			t.Fatal(err)
		}
	}
}
