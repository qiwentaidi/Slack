package webscan

import (
	"crypto/tls"
	"net"
	"strings"
	"time"

	zasn1 "github.com/zmap/zcrypto/encoding/asn1"
	zpkix "github.com/zmap/zcrypto/x509/pkix"
)

type CertResponse struct {
	SubjectCN string
	SubjectDN string
	IssuerCN  string
	IssuerDN  string
	IssuerOrg []string
}

func GetTLSString(protocol, host string) string {
	TLSData := getCertResponse(protocol, host)
	var result strings.Builder
	// 预分配一个中等大小的缓冲区，以避免频繁的内存重新分配
	result.Grow(512)

	if TLSData == nil {
		return ""
	}
	result.WriteString("SubjectCN: " + TLSData.SubjectCN + "\n")
	result.WriteString("SubjectDN: " + TLSData.SubjectDN + "\n")
	result.WriteString("IssuerCN: " + TLSData.IssuerCN + "\n")
	result.WriteString("IssuerDN: " + TLSData.IssuerDN + "\n")
	result.WriteString("IssuerOrg: \n")

	for _, v := range TLSData.IssuerOrg {
		result.WriteString("    - " + v + "\n")
	}

	return result.String()
}

func getCertResponse(protocol, host string) *CertResponse {
	if protocol == "https" || protocol == "tls" {
		conn, err := tls.DialWithDialer(&net.Dialer{Timeout: time.Duration(3) * time.Second}, "tcp", host, &tls.Config{InsecureSkipVerify: true})
		if err != nil {
			return nil
		}
		defer conn.Close()
		cert := conn.ConnectionState().PeerCertificates[0]
		return &CertResponse{
			IssuerCN:  cert.Issuer.CommonName,
			IssuerDN:  ParseASN1DNSequenceWithZpkixOrDefault(cert.RawIssuer, cert.Issuer.String()),
			SubjectCN: cert.Subject.CommonName,
			SubjectDN: ParseASN1DNSequenceWithZpkixOrDefault(cert.RawSubject, cert.Subject.String()),
			IssuerOrg: cert.Issuer.Organization,
		}
	}
	return nil
}

// ParseASN1DNSequenceWithZpkixOrDefault return the parsed value of ASN1DNSequence or a default string value
func ParseASN1DNSequenceWithZpkixOrDefault(data []byte, defaultValue string) string {
	if value := ParseASN1DNSequenceWithZpkix(data); value != "" {
		return value
	}
	return defaultValue
}

// ParseASN1DNSequenceWithZpkix tries to parse raw ASN1 of a TLS DN with zpkix and
// zasn1 library which includes additional information not parsed by go standard
// library which may be useful.
//
// If the parsing fails, a blank string is returned and the standard library data is used.
func ParseASN1DNSequenceWithZpkix(data []byte) string {
	var rdnSequence zpkix.RDNSequence
	var name zpkix.Name
	if _, err := zasn1.Unmarshal(data, &rdnSequence); err != nil {
		return ""
	}
	name.FillFromRDNSequence(&rdnSequence)
	dnParsedString := name.String()
	return dnParsedString
}
