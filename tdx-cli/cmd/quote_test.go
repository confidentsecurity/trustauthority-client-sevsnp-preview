/*
 *   Copyright (c) 2022 Intel Corporation
 *   All rights reserved.
 *   SPDX-License-Identifier: BSD-3-Clause
 */

package cmd

import (
	"testing"

	"confidentsecurity/trustauthority-client-sevsnp-preview/tdx-cli/constants"
	"github.com/stretchr/testify/assert"
)

func TestQuoteCmd(t *testing.T) {

	tt := []struct {
		args        []string
		wantErr     bool
		description string
	}{
		{
			args:        []string{constants.QuoteCmd},
			wantErr:     false,
			description: "Test without inputs",
		},
		{
			args: []string{constants.QuoteCmd, "--" + constants.UserDataOption, "dGVzdHVzZXJkYXRh",
				"--" + constants.NonceOption, "dGVzdHVzZXJkYXRh"},
			wantErr:     false,
			description: "Test with all valid inputs",
		},
		{
			args:        []string{constants.QuoteCmd, "--" + constants.UserDataOption, "u$erd@t@"},
			wantErr:     true,
			description: "Test with malformed userdata",
		},
		{
			args:        []string{constants.QuoteCmd, "--" + constants.NonceOption, "n@nce"},
			wantErr:     true,
			description: "Test with malformed nonce",
		},
	}

	for _, tc := range tt {
		_, err := execute(t, rootCmd, tc.args...)

		if tc.wantErr == true {
			assert.Error(t, err)
		} else {
			assert.NoError(t, err)
		}
	}
}
