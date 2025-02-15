//go:build !test

/*
 *   Copyright (c) 2024 Intel Corporation
 *   All rights reserved.
 *   SPDX-License-Identifier: BSD-3-Clause
 */
package sevsnp

import (
	"confidentsecurity/trustauthority-client-sevsnp-preview/go-connector"
)

// sevsnpAdapter manages sevsnp report collection from sevsnp enabled platform
type sevsnpAdapter struct {
	uData []byte
	uVmpl uint32
}

// NewEvidenceAdapter returns a new sevsnp Adapter instance
func NewEvidenceAdapter(udata []byte, uvmpl uint32) (connector.EvidenceAdapter, error) {
	return &sevsnpAdapter{
		uData: udata,
		uVmpl: uvmpl,
	}, nil
}
