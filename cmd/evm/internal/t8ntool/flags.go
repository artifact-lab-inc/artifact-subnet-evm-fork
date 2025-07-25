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
// Copyright 2020 The go-ethereum Authors
// This file is part of go-ethereum.
//
// go-ethereum is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// go-ethereum is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with go-ethereum. If not, see <http://www.gnu.org/licenses/>.

package t8ntool

import (
	"fmt"
	"strings"

	"github.com/ava-labs/libevm/core/vm"
	"github.com/ava-labs/subnet-evm/tests"
	"github.com/urfave/cli/v2"
)

var (
	TraceFlag = &cli.BoolFlag{
		Name:  "trace",
		Usage: "Configures the use of the JSON opcode tracer. This tracer emits traces to files as trace-<txIndex>-<txHash>.jsonl",
	}
	TraceTracerFlag = &cli.StringFlag{
		Name:  "trace.tracer",
		Usage: "Configures the use of a custom tracer, e.g native or js tracers. Examples are callTracer and 4byteTracer. These tracers emit results into files as trace-<txIndex>-<txHash>.json",
	}
	TraceTracerConfigFlag = &cli.StringFlag{
		Name:  "trace.jsonconfig",
		Usage: "The configurations for the custom tracer specified by --trace.tracer. If provided, must be in JSON format",
	}
	TraceEnableMemoryFlag = &cli.BoolFlag{
		Name:  "trace.memory",
		Usage: "Enable full memory dump in traces",
	}
	TraceDisableStackFlag = &cli.BoolFlag{
		Name:  "trace.nostack",
		Usage: "Disable stack output in traces",
	}
	TraceEnableReturnDataFlag = &cli.BoolFlag{
		Name:  "trace.returndata",
		Usage: "Enable return data output in traces",
	}
	OutputBasedir = &cli.StringFlag{
		Name:  "output.basedir",
		Usage: "Specifies where output files are placed. Will be created if it does not exist.",
		Value: "",
	}
	OutputBodyFlag = &cli.StringFlag{
		Name:  "output.body",
		Usage: "If set, the RLP of the transactions (block body) will be written to this file.",
		Value: "",
	}
	OutputAllocFlag = &cli.StringFlag{
		Name: "output.alloc",
		Usage: "Determines where to put the `alloc` of the post-state.\n" +
			"\t`stdout` - into the stdout output\n" +
			"\t`stderr` - into the stderr output\n" +
			"\t<file> - into the file <file> ",
		Value: "alloc.json",
	}
	OutputResultFlag = &cli.StringFlag{
		Name: "output.result",
		Usage: "Determines where to put the `result` (stateroot, txroot etc) of the post-state.\n" +
			"\t`stdout` - into the stdout output\n" +
			"\t`stderr` - into the stderr output\n" +
			"\t<file> - into the file <file> ",
		Value: "result.json",
	}
	OutputBlockFlag = &cli.StringFlag{
		Name: "output.block",
		Usage: "Determines where to put the `block` after building.\n" +
			"\t`stdout` - into the stdout output\n" +
			"\t`stderr` - into the stderr output\n" +
			"\t<file> - into the file <file> ",
		Value: "block.json",
	}
	InputAllocFlag = &cli.StringFlag{
		Name:  "input.alloc",
		Usage: "`stdin` or file name of where to find the prestate alloc to use.",
		Value: "alloc.json",
	}
	InputEnvFlag = &cli.StringFlag{
		Name:  "input.env",
		Usage: "`stdin` or file name of where to find the prestate env to use.",
		Value: "env.json",
	}
	InputTxsFlag = &cli.StringFlag{
		Name: "input.txs",
		Usage: "`stdin` or file name of where to find the transactions to apply. " +
			"If the file extension is '.rlp', then the data is interpreted as an RLP list of signed transactions." +
			"The '.rlp' format is identical to the output.body format.",
		Value: "txs.json",
	}
	InputHeaderFlag = &cli.StringFlag{
		Name:  "input.header",
		Usage: "`stdin` or file name of where to find the block header to use.",
		Value: "header.json",
	}
	InputOmmersFlag = &cli.StringFlag{
		Name:  "input.ommers",
		Usage: "`stdin` or file name of where to find the list of ommer header RLPs to use.",
	}
	InputTxsRlpFlag = &cli.StringFlag{
		Name:  "input.txs",
		Usage: "`stdin` or file name of where to find the transactions list in RLP form.",
		Value: "txs.rlp",
	}
	SealCliqueFlag = &cli.StringFlag{
		Name:  "seal.clique",
		Usage: "Seal block with Clique. `stdin` or file name of where to find the Clique sealing data.",
	}
	RewardFlag = &cli.Int64Flag{
		Name:  "state.reward",
		Usage: "Mining reward. Set to -1 to disable",
		Value: 0,
	}
	ChainIDFlag = &cli.Int64Flag{
		Name:  "state.chainid",
		Usage: "ChainID to use",
		Value: 1,
	}
	ForknameFlag = &cli.StringFlag{
		Name: "state.fork",
		Usage: fmt.Sprintf("Name of ruleset to use."+
			"\n\tAvailable forknames:"+
			"\n\t    %v"+
			"\n\tAvailable extra eips:"+
			"\n\t    %v"+
			"\n\tSyntax <forkname>(+ExtraEip)",
			strings.Join(tests.AvailableForks(), "\n\t    "),
			strings.Join(vm.ActivateableEips(), ", ")),
		Value: "GrayGlacier",
	}
	VerbosityFlag = &cli.IntFlag{
		Name:  "verbosity",
		Usage: "sets the verbosity level",
		Value: 3,
	}
)
