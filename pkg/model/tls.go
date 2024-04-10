package model

import (
	"crypto/tls"
	"crypto/x509"

	"github.com/WangYihang/http-grab/pkg/util"
)

// CertificateWrapper is copied from crypto/x509/x509.go
type CertificateWrapper struct {
	Fingerprint string `json:"fingerprint" bson:"fingerprint"`
	Raw         []byte `json:"raw" bson:"raw"`
	// Raw         []byte // Complete ASN.1 DER content (certificate, signature algorithm and signature).
	// RawTBSCertificate       []byte // Certificate part of raw ASN.1 DER content.
	// RawSubjectPublicKeyInfo []byte // DER encoded SubjectPublicKeyInfo.
	// RawSubject              []byte // DER encoded Subject
	// RawIssuer               []byte // DER encoded Issuer

	// Signature          []byte
	// SignatureAlgorithm x509.SignatureAlgorithm

	// PublicKeyAlgorithm x509.PublicKeyAlgorithm
	// PublicKey          crypto.PublicKey

	// Version             int
	// SerialNumber        *big.Int
	// Issuer              pkix.Name
	// Subject             pkix.Name
	// NotBefore, NotAfter time.Time // Validity bounds.
	// KeyUsage            x509.KeyUsage

	// // Extensions contains raw X.509 extensions. When parsing certificates,
	// // this can be used to extract non-critical extensions that are not
	// // parsed by this package. When marshaling certificates, the Extensions
	// // field is ignored, see ExtraExtensions.
	// Extensions []pkix.Extension

	// // ExtraExtensions contains extensions to be copied, raw, into any
	// // marshaled certificates. Values override any extensions that would
	// // otherwise be produced based on the other fields. The ExtraExtensions
	// // field is not populated when parsing certificates, see Extensions.
	// ExtraExtensions []pkix.Extension

	// // UnhandledCriticalExtensions contains a list of extension IDs that
	// // were not (fully) processed when parsing. Verify will fail if this
	// // slice is non-empty, unless verification is delegated to an OS
	// // library which understands all the critical extensions.
	// //
	// // Users can access these extensions using Extensions and can remove
	// // elements from this slice if they believe that they have been
	// // handled.
	// UnhandledCriticalExtensions []asn1.ObjectIdentifier

	// ExtKeyUsage        []x509.ExtKeyUsage      // Sequence of extended key usages.
	// UnknownExtKeyUsage []asn1.ObjectIdentifier // Encountered extended key usages unknown to this package.

	// // BasicConstraintsValid indicates whether IsCA, MaxPathLen,
	// // and MaxPathLenZero are valid.
	// BasicConstraintsValid bool
	// IsCA                  bool

	// // MaxPathLen and MaxPathLenZero indicate the presence and
	// // value of the BasicConstraints' "pathLenConstraint".
	// //
	// // When parsing a certificate, a positive non-zero MaxPathLen
	// // means that the field was specified, -1 means it was unset,
	// // and MaxPathLenZero being true mean that the field was
	// // explicitly set to zero. The case of MaxPathLen==0 with MaxPathLenZero==false
	// // should be treated equivalent to -1 (unset).
	// //
	// // When generating a certificate, an unset pathLenConstraint
	// // can be requested with either MaxPathLen == -1 or using the
	// // zero value for both MaxPathLen and MaxPathLenZero.
	// MaxPathLen int
	// // MaxPathLenZero indicates that BasicConstraintsValid==true
	// // and MaxPathLen==0 should be interpreted as an actual
	// // maximum path length of zero. Otherwise, that combination is
	// // interpreted as MaxPathLen not being set.
	// MaxPathLenZero bool

	// SubjectKeyId   []byte
	// AuthorityKeyId []byte

	// // RFC 5280, 4.2.2.1 (Authority Information Access)
	// OCSPServer            []string
	// IssuingCertificateURL []string

	// // Subject Alternate Name values. (Note that these values may not be valid
	// // if invalid values were contained within a parsed certificate. For
	// // example, an element of DNSNames may not be a valid DNS domain name.)
	// DNSNames       []string
	// EmailAddresses []string
	// IPAddresses    []net.IP
	// URIs           []*url.URL

	// // Name constraints
	// PermittedDNSDomainsCritical bool // if true then the name constraints are marked critical.
	// PermittedDNSDomains         []string
	// ExcludedDNSDomains          []string
	// PermittedIPRanges           []*net.IPNet
	// ExcludedIPRanges            []*net.IPNet
	// PermittedEmailAddresses     []string
	// ExcludedEmailAddresses      []string
	// PermittedURIDomains         []string
	// ExcludedURIDomains          []string

	// // CRL Distribution Points
	// CRLDistributionPoints []string

	// // PolicyIdentifiers contains asn1.ObjectIdentifiers, the components
	// // of which are limited to int32. If a certificate contains a policy which
	// // cannot be represented by asn1.ObjectIdentifier, it will not be included in
	// // PolicyIdentifiers, but will be present in Policies, which contains all parsed
	// // policy OIDs.
	// PolicyIdentifiers []asn1.ObjectIdentifier

	// // Policies contains all policy identifiers included in the certificate.
	// Policies []x509.OID
}

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

	// PeerCertificate *x509.Certificate `json:"peer_certificate,omitempty"`
	PeerCertificate *CertificateWrapper `json:"peer_certificate,omitempty"`

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

func NewTLS(cs *tls.ConnectionState) *TLS {
	var peerCertificate *x509.Certificate
	if len(cs.PeerCertificates) > 0 {
		peerCertificate = cs.PeerCertificates[0]
	}
	return &TLS{
		Version:            cs.Version,
		HandshakeComplete:  cs.HandshakeComplete,
		DidResume:          cs.DidResume,
		CipherSuite:        cs.CipherSuite,
		NegotiatedProtocol: cs.NegotiatedProtocol,
		ServerName:         cs.ServerName,
		PeerCertificate: &CertificateWrapper{
			Raw:         peerCertificate.Raw,
			Fingerprint: util.Sha256(peerCertificate.Raw),
		},
	}
}
