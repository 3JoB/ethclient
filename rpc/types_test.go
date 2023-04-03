// Copyright 2015 The go-ethereum Authors
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

package rpc

import (
	"testing"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/math"
	"github.com/goccy/go-json"
	"github.com/goccy/go-reflect"
)

func TestBlockNumberJSONUnmarshal(t *testing.T) {
	tests := []struct {
		input    string
		mustFail bool
		expected BlockNumber
	}{
		0:  {input: `"0x"`, mustFail: true, expected: BlockNumber(0)},
		1:  {input: `"0x0"`, mustFail: false, expected: BlockNumber(0)},
		2:  {input: `"0X1"`, mustFail: false, expected: BlockNumber(1)},
		3:  {input: `"0x00"`, mustFail: true, expected: BlockNumber(0)},
		4:  {input: `"0x01"`, mustFail: true, expected: BlockNumber(0)},
		5:  {input: `"0x1"`, mustFail: false, expected: BlockNumber(1)},
		6:  {input: `"0x12"`, mustFail: false, expected: BlockNumber(18)},
		7:  {input: `"0x7fffffffffffffff"`, mustFail: false, expected: BlockNumber(math.MaxInt64)},
		8:  {input: `"0x8000000000000000"`, mustFail: true, expected: BlockNumber(0)},
		9:  {input: "0", mustFail: true, expected: BlockNumber(0)},
		10: {input: `"ff"`, mustFail: true, expected: BlockNumber(0)},
		11: {input: `"pending"`, mustFail: false, expected: PendingBlockNumber},
		12: {input: `"latest"`, mustFail: false, expected: LatestBlockNumber},
		13: {input: `"earliest"`, mustFail: false, expected: EarliestBlockNumber},
		14: {input: `someString`, mustFail: true, expected: BlockNumber(0)},
		15: {input: `""`, mustFail: true, expected: BlockNumber(0)},
		16: {input: ``, mustFail: true, expected: BlockNumber(0)},
	}

	for i, test := range tests {
		var num BlockNumber
		err := json.Unmarshal([]byte(test.input), &num)
		if test.mustFail && err == nil {
			t.Errorf("Test %d should fail", i)
			continue
		}
		if !test.mustFail && err != nil {
			t.Errorf("Test %d should pass but got err: %v", i, err)
			continue
		}
		if num != test.expected {
			t.Errorf("Test %d got unexpected value, want %d, got %d", i, test.expected, num)
		}
	}
}

func TestBlockNumberOrHash_UnmarshalJSON(t *testing.T) {
	tests := []struct {
		input    string
		mustFail bool
		expected BlockNumberOrHash
	}{
		0:  {input: `"0x"`, mustFail: true, expected: BlockNumberOrHash{}},
		1:  {input: `"0x0"`, mustFail: false, expected: BlockNumberOrHashWithNumber(0)},
		2:  {input: `"0X1"`, mustFail: false, expected: BlockNumberOrHashWithNumber(1)},
		3:  {input: `"0x00"`, mustFail: true, expected: BlockNumberOrHash{}},
		4:  {input: `"0x01"`, mustFail: true, expected: BlockNumberOrHash{}},
		5:  {input: `"0x1"`, mustFail: false, expected: BlockNumberOrHashWithNumber(1)},
		6:  {input: `"0x12"`, mustFail: false, expected: BlockNumberOrHashWithNumber(18)},
		7:  {input: `"0x7fffffffffffffff"`, mustFail: false, expected: BlockNumberOrHashWithNumber(math.MaxInt64)},
		8:  {input: `"0x8000000000000000"`, mustFail: true, expected: BlockNumberOrHash{}},
		9:  {input: "0", mustFail: true, expected: BlockNumberOrHash{}},
		10: {input: `"ff"`, mustFail: true, expected: BlockNumberOrHash{}},
		11: {input: `"pending"`, mustFail: false, expected: BlockNumberOrHashWithNumber(PendingBlockNumber)},
		12: {input: `"latest"`, mustFail: false, expected: BlockNumberOrHashWithNumber(LatestBlockNumber)},
		13: {input: `"earliest"`, mustFail: false, expected: BlockNumberOrHashWithNumber(EarliestBlockNumber)},
		14: {input: `someString`, mustFail: true, expected: BlockNumberOrHash{}},
		15: {input: `""`, mustFail: true, expected: BlockNumberOrHash{}},
		16: {input: ``, mustFail: true, expected: BlockNumberOrHash{}},
		17: {input: `"0x0000000000000000000000000000000000000000000000000000000000000000"`, mustFail: false, expected: BlockNumberOrHashWithHash(common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000000"), false)},
		18: {input: `{"blockHash":"0x0000000000000000000000000000000000000000000000000000000000000000"}`, mustFail: false, expected: BlockNumberOrHashWithHash(common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000000"), false)},
		19: {input: `{"blockHash":"0x0000000000000000000000000000000000000000000000000000000000000000","requireCanonical":false}`, mustFail: false, expected: BlockNumberOrHashWithHash(common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000000"), false)},
		20: {input: `{"blockHash":"0x0000000000000000000000000000000000000000000000000000000000000000","requireCanonical":true}`, mustFail: false, expected: BlockNumberOrHashWithHash(common.HexToHash("0x0000000000000000000000000000000000000000000000000000000000000000"), true)},
		21: {input: `{"blockNumber":"0x1"}`, mustFail: false, expected: BlockNumberOrHashWithNumber(1)},
		22: {input: `{"blockNumber":"pending"}`, mustFail: false, expected: BlockNumberOrHashWithNumber(PendingBlockNumber)},
		23: {input: `{"blockNumber":"latest"}`, mustFail: false, expected: BlockNumberOrHashWithNumber(LatestBlockNumber)},
		24: {input: `{"blockNumber":"earliest"}`, mustFail: false, expected: BlockNumberOrHashWithNumber(EarliestBlockNumber)},
		25: {input: `{"blockNumber":"0x1", "blockHash":"0x0000000000000000000000000000000000000000000000000000000000000000"}`, mustFail: true, expected: BlockNumberOrHash{}},
	}

	for i, test := range tests {
		var bnh BlockNumberOrHash
		err := json.Unmarshal([]byte(test.input), &bnh)
		if test.mustFail && err == nil {
			t.Errorf("Test %d should fail", i)
			continue
		}
		if !test.mustFail && err != nil {
			t.Errorf("Test %d should pass but got err: %v", i, err)
			continue
		}
		hash, hashOk := bnh.Hash()
		expectedHash, expectedHashOk := test.expected.Hash()
		num, numOk := bnh.Number()
		expectedNum, expectedNumOk := test.expected.Number()
		if bnh.RequireCanonical != test.expected.RequireCanonical ||
			hash != expectedHash || hashOk != expectedHashOk ||
			num != expectedNum || numOk != expectedNumOk {
			t.Errorf("Test %d got unexpected value, want %v, got %v", i, test.expected, bnh)
		}
	}
}

func TestBlockNumberOrHash_WithNumber_MarshalAndUnmarshal(t *testing.T) {
	tests := []struct {
		name   string
		number int64
	}{
		{name: "max", number: math.MaxInt64},
		{name: "pending", number: int64(PendingBlockNumber)},
		{name: "latest", number: int64(LatestBlockNumber)},
		{name: "earliest", number: int64(EarliestBlockNumber)},
	}
	for _, test := range tests {
		test := test
		t.Run(test.name, func(t *testing.T) {
			bnh := BlockNumberOrHashWithNumber(BlockNumber(test.number))
			marshalled, err := json.Marshal(bnh)
			if err != nil {
				t.Fatal("cannot marshal:", err)
			}
			var unmarshalled BlockNumberOrHash
			err = json.Unmarshal(marshalled, &unmarshalled)
			if err != nil {
				t.Fatal("cannot unmarshal:", err)
			}
			if !reflect.DeepEqual(bnh, unmarshalled) {
				t.Fatalf("wrong result: expected %v, got %v", bnh, unmarshalled)
			}
		})
	}
}
