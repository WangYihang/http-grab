package model

import (
	"crypto/tls"
	"crypto/x509"
)

type TLS struct {
	// Version is the TLS version used by the connection (e.g. VersionTLS12).
	Version uint16 `json:"version"`

	// HandshakeComplete is true if the handshake has concluded.
	HandshakeComplete bool `json:"handshake_complete"`

	// DidResume is true if this connection was successfully resumed from a
	// previous session with a session ticket or similar mechanism.
	DidResume bool `json:"did_resume"`

	// CipherSuite is the cipher suite negotiated for the connection (e.g.
	// TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256, TLS_AES_128_GCM_SHA256).
	CipherSuite uint16 `json:"cipher_suite"`

	// NegotiatedProtocol is the application protocol negotiated with ALPN.
	NegotiatedProtocol string `json:"negotiated_protocol"`

	// ServerName is the value of the Server Name Indication extension sent by
	// the client. It's available both on the server and on the client side.
	ServerName string `json:"server_name"`

	PeerCertificate *x509.Certificate `json:"peer_certificate,omitempty"`

	// // PeerCertificates are the parsed certificates sent by the peer, in the
	// // order in which they were sent. The first element is the leaf certificate
	// // that the connection is verified against.
	// //
	// // On the client side, it can't be empty. On the server side, it can be
	// // empty if Config.ClientAuth is not RequireAnyClientCert or
	// // RequireAndVerifyClientCert.
	// //
	// // PeerCertificates and its contents should not be modified.
	// PeerCertificates []*x509.Certificate `json:"peer_certificates"`

	// // VerifiedChains is a list of one or more chains where the first element is
	// // PeerCertificates[0] and the last element is from Config.RootCAs (on the
	// // client side) or Config.ClientCAs (on the server side).
	// //
	// // On the client side, it's set if Config.InsecureSkipVerify is false. On
	// // the server side, it's set if Config.ClientAuth is VerifyClientCertIfGiven
	// // (and the peer provided a certificate) or RequireAndVerifyClientCert.
	// //
	// // VerifiedChains and its contents should not be modified.
	// VerifiedChains [][]*x509.Certificate `json:"verified_chains"`

	// // SignedCertificateTimestamps is a list of SCTs provided by the peer
	// // through the TLS handshake for the leaf certificate, if any.
	// SignedCertificateTimestamps [][]byte `json:"signed_certificate_timestamps"`

	// // OCSPResponse is a stapled Online Certificate Status Protocol (OCSP)
	// // response provided by the peer for the leaf certificate, if any.
	// OCSPResponse []byte

	// // TLSUnique contains the "tls-unique" channel binding value (see RFC 5929,
	// // Section 3). This value will be nil for TLS 1.3 connections and for
	// // resumed connections that don't support Extended Master Secret (RFC 7627).
	// TLSUnique []byte

	// // ekm is a closure exposed via ExportKeyingMaterial.
	// ekm func(label string, context []byte, length int) ([]byte, error)
}

func NewTLS(cs *tls.ConnectionState) TLS {
	var peerCertificate *x509.Certificate
	if len(cs.PeerCertificates) > 0 {
		peerCertificate = cs.PeerCertificates[0]
	}
	return TLS{
		Version:            cs.Version,
		HandshakeComplete:  cs.HandshakeComplete,
		DidResume:          cs.DidResume,
		CipherSuite:        cs.CipherSuite,
		NegotiatedProtocol: cs.NegotiatedProtocol,
		ServerName:         cs.ServerName,
		PeerCertificate:    peerCertificate,
	}
}
