// Copyright (c) 2024 H0llyW00dz All rights reserved.
//
// License: BSD 3-Clause License

package server

import (
	"crypto/x509"
	"encoding/asn1"
	"encoding/base64"
	"encoding/binary"
	"fmt"

	"github.com/bytedance/sonic"
)

var (
	// OIDExtensionCTSCT is the OID for the Certificate Transparency SCT extension.
	OIDExtensionCTSCT = asn1.ObjectIdentifier{1, 3, 6, 1, 4, 1, 11129, 2, 4, 2}
)

// SignedCertificateTimestampList represents a list of Signed Certificate Timestamps (SCTs).
type SignedCertificateTimestampList struct {
	SCTList []asn1.RawValue `asn1:"optional,omitempty"`
}

// CTVerifier represents a Certificate Transparency verifier.
type CTVerifier struct {
	SignedCertificateTimestampList
}

// SCTData represents the SCT data for ASN.1 marshaling.
type SCTData struct {
	Version    uint8
	LogID      []byte
	Timestamp  uint64
	Extensions []byte
	Signature  []byte
}

// VerifyCertificateTransparency verifies the Certificate Transparency information for a given certificate.
//
// TODO: Improve this protocol, as it is currently unstable due to the difficulty of testing certificates in localhost.
func (ct *CTVerifier) VerifyCertificateTransparency(cert *x509.Certificate) error {
	// Check if the certificate has SCTs (Signed Certificate Timestamps)
	scts, err := ct.ExtractSCTsFromCertificate(cert)
	if err != nil {
		return fmt.Errorf("failed to extract SCTs from certificate: %v", err)
	}

	if len(scts) == 0 {
		return fmt.Errorf("certificate does not have any SCTs")
	}

	// Verify each SCT against the CT logs
	for _, sct := range scts {
		if err := ct.VerifySCT(sct, cert); err != nil {
			return err
		}
	}

	return nil
}

// ExtractSCTsFromCertificate extracts the Signed Certificate Timestamps (SCTs) from a certificate.
//
// TODO: Improve this protocol, as it is currently unstable due to the difficulty of testing certificates in localhost.
func (ct *CTVerifier) ExtractSCTsFromCertificate(cert *x509.Certificate) ([]*SCTResponse, error) {
	var scts []*SCTResponse

	for _, ext := range cert.Extensions {
		if ext.Id.Equal(OIDExtensionCTSCT) {
			if len(ext.Value) < 44 {
				return nil, fmt.Errorf("invalid SCT data: insufficient length")
			}

			version := ext.Value[0]
			logID := ext.Value[1:33]
			timestamp := binary.BigEndian.Uint64(ext.Value[33:41])
			extensionsLen := int(ext.Value[41])
			extensions := ext.Value[42 : 42+extensionsLen]
			signature := ext.Value[42+extensionsLen:]

			sct := &SCTResponse{
				SCTVersion: version,
				ID:         string(logID),
				Timestamp:  timestamp,
				Extensions: string(extensions),
				Signature:  base64.StdEncoding.EncodeToString(signature),
			}

			scts = append(scts, sct)
		}
	}

	return scts, nil
}

// VerifySCT verifies a Signed Certificate Timestamp (SCT) against the certificate.
//
// TODO: Improve this protocol, as it is currently unstable due to the difficulty of testing certificates in localhost.
func (ct *CTVerifier) VerifySCT(sct *SCTResponse, cert *x509.Certificate) error {
	jsonConfig := json{
		Marshal:   sonic.Marshal,
		Unmarshal: sonic.Unmarshal,
	}
	sctVerifier := SCTVerifier{
		Response: *sct,
		Cert:     cert,
		json:     jsonConfig,
	}

	// Verify the SCT
	if err := sctVerifier.VerifySCT(); err != nil {
		return fmt.Errorf("invalid SCT: %v", err)
	}

	return nil
}
