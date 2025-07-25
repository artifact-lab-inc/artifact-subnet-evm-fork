// Copyright (C) 2019-2025, Ava Labs, Inc. All rights reserved.
// See the file LICENSE for licensing terms.

package feemanager

import (
	"math/big"
	"testing"

	"github.com/ava-labs/libevm/common"
	"github.com/ava-labs/subnet-evm/commontype"
	"github.com/ava-labs/subnet-evm/precompile/allowlist/allowlisttest"
	"github.com/ava-labs/subnet-evm/precompile/precompileconfig"
	"github.com/ava-labs/subnet-evm/precompile/precompiletest"
	"github.com/ava-labs/subnet-evm/utils"
	"go.uber.org/mock/gomock"
)

var validFeeConfig = commontype.FeeConfig{
	GasLimit:        big.NewInt(8_000_000),
	TargetBlockRate: 2, // in seconds

	MinBaseFee:               big.NewInt(25_000_000_000),
	TargetGas:                big.NewInt(15_000_000),
	BaseFeeChangeDenominator: big.NewInt(36),

	MinBlockGasCost:  big.NewInt(0),
	MaxBlockGasCost:  big.NewInt(1_000_000),
	BlockGasCostStep: big.NewInt(200_000),
}

func TestVerify(t *testing.T) {
	admins := []common.Address{allowlisttest.TestAdminAddr}
	invalidFeeConfig := validFeeConfig
	invalidFeeConfig.GasLimit = big.NewInt(0)
	tests := map[string]precompiletest.ConfigVerifyTest{
		"invalid initial fee manager config": {
			Config:        NewConfig(utils.NewUint64(3), admins, nil, nil, &invalidFeeConfig),
			ExpectedError: "gasLimit = 0 cannot be less than or equal to 0",
		},
		"nil initial fee manager config": {
			Config:        NewConfig(utils.NewUint64(3), admins, nil, nil, &commontype.FeeConfig{}),
			ExpectedError: "gasLimit cannot be nil",
		},
	}
	allowlisttest.VerifyPrecompileWithAllowListTests(t, Module, tests)
}

func TestEqual(t *testing.T) {
	admins := []common.Address{allowlisttest.TestAdminAddr}
	enableds := []common.Address{allowlisttest.TestEnabledAddr}
	tests := map[string]precompiletest.ConfigEqualTest{
		"non-nil config and nil other": {
			Config:   NewConfig(utils.NewUint64(3), admins, enableds, nil, nil),
			Other:    nil,
			Expected: false,
		},
		"different type": {
			Config:   NewConfig(utils.NewUint64(3), admins, enableds, nil, nil),
			Other:    precompileconfig.NewMockConfig(gomock.NewController(t)),
			Expected: false,
		},
		"different timestamp": {
			Config:   NewConfig(utils.NewUint64(3), admins, nil, nil, nil),
			Other:    NewConfig(utils.NewUint64(4), admins, nil, nil, nil),
			Expected: false,
		},
		"non-nil initial config and nil initial config": {
			Config:   NewConfig(utils.NewUint64(3), admins, nil, nil, &validFeeConfig),
			Other:    NewConfig(utils.NewUint64(3), admins, nil, nil, nil),
			Expected: false,
		},
		"different initial config": {
			Config: NewConfig(utils.NewUint64(3), admins, nil, nil, &validFeeConfig),
			Other: NewConfig(utils.NewUint64(3), admins, nil, nil,
				func() *commontype.FeeConfig {
					c := validFeeConfig
					c.GasLimit = big.NewInt(123)
					return &c
				}()),
			Expected: false,
		},
		"same config": {
			Config:   NewConfig(utils.NewUint64(3), admins, nil, nil, &validFeeConfig),
			Other:    NewConfig(utils.NewUint64(3), admins, nil, nil, &validFeeConfig),
			Expected: true,
		},
	}
	allowlisttest.EqualPrecompileWithAllowListTests(t, Module, tests)
}
