package tls

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/tls"
	"crypto/x509"
	"crypto/x509/pkix"
	"encoding/pem"
	"fmt"
	"math/big"
	"net"
	"os"
	"time"
)

const (
	// Encryption levels
	EncryptionOff      = 0x00 // No encryption (cleartext)
	EncryptionOn       = 0x01 // SSL/TLS encryption if supported
	EncryptionRequired = 0x02 // Encryption required (reject if not supported)
)

// Config represents SSL/TLS configuration
type Config struct {
	Enabled        bool   // Enable SSL/TLS encryption
	ForceEncryption bool   // Require encryption (reject if not supported)
	CertFile       string // Path to certificate file (.pem/.crt)
	KeyFile        string // Path to private key file (.key)
	MinVersion     uint16 // Minimum TLS version (1.2, 1.3)
	MaxVersion     uint16 // Maximum TLS version
	ClientAuth     tls.ClientAuthType // Client certificate auth
	TrustServerCert bool // Trust server certificate (development only)
}

// DefaultConfig returns default SSL/TLS configuration
func DefaultConfig() *Config {
	return &Config{
		Enabled:        false, // Disabled by default (for compatibility)
		ForceEncryption: false,
		CertFile:       "./certs/server.crt",
		KeyFile:        "./certs/server.key",
		MinVersion:     tls.VersionTLS12,
		MaxVersion:     tls.VersionTLS13,
		ClientAuth:     tls.NoClientCert,
		TrustServerCert: false,
	}
}

// LoadCertificate loads SSL/TLS certificate from files
func LoadCertificate(certFile, keyFile string) (*tls.Certificate, error) {
	// Load certificate and private key
	cert, err := tls.LoadX509KeyPair(certFile, keyFile)
	if err != nil {
		return nil, fmt.Errorf("failed to load certificate: %w", err)
	}

	return &cert, nil
}

// GenerateSelfSignedCertificate generates a self-signed certificate for development
func GenerateSelfSignedCertificate(commonName string, validityDays int) (*tls.Certificate, error) {
	// Generate private key
	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return nil, fmt.Errorf("failed to generate private key: %w", err)
	}

	// Create certificate template
	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Organization: []string{"MSSQL TDS Server"},
			CommonName:   commonName,
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(0, 0, validityDays),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	// Generate certificate
	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &privateKey.PublicKey, privateKey)
	if err != nil {
		return nil, fmt.Errorf("failed to create certificate: %w", err)
	}

	// Create TLS certificate
	cert := tls.Certificate{
		Certificate: [][]byte{certDER},
		PrivateKey:  privateKey,
	}

	return &cert, nil
}

// SaveCertificate saves certificate and private key to files
func SaveCertificate(cert *tls.Certificate, certFile, keyFile string) error {
	// Save certificate
	certFileWriter, err := os.Create(certFile)
	if err != nil {
		return fmt.Errorf("failed to create certificate file: %w", err)
	}
	defer certFileWriter.Close()

	// Write certificate in PEM format
	certPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: cert.Certificate[0],
	})
	_, err = certFileWriter.Write(certPEM)
	if err != nil {
		return fmt.Errorf("failed to write certificate: %w", err)
	}

	// Save private key
	keyFileWriter, err := os.Create(keyFile)
	if err != nil {
		return fmt.Errorf("failed to create private key file: %w", err)
	}
	defer keyFileWriter.Close()

	// Write private key in PEM format
	privateKeyBytes, err := x509.MarshalPKCS8PrivateKey(cert.PrivateKey)
	if err != nil {
		return fmt.Errorf("failed to marshal private key: %w", err)
	}

	keyPEM := pem.EncodeToMemory(&pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: privateKeyBytes,
	})
	_, err = keyFileWriter.Write(keyPEM)
	if err != nil {
		return fmt.Errorf("failed to write private key: %w", err)
	}

	// Set file permissions
	os.Chmod(certFile, 0644)
	os.Chmod(keyFile, 0600)

	return nil
}

// LoadOrCreateCertificate loads existing certificate or creates a new one
func LoadOrCreateCertificate(certFile, keyFile string, commonName string) (*tls.Certificate, error) {
	// Check if certificate files exist
	if _, err := os.Stat(certFile); err == nil {
		if _, err := os.Stat(keyFile); err == nil {
			// Load existing certificate
			return LoadCertificate(certFile, keyFile)
		}
	}

	// Ensure directory exists
	certDir := "./certs"
	os.MkdirAll(certDir, 0755)

	// Generate new self-signed certificate
	cert, err := GenerateSelfSignedCertificate(commonName, 365)
	if err != nil {
		return nil, fmt.Errorf("failed to generate certificate: %w", err)
	}

	// Save certificate to files
	err = SaveCertificate(cert, certFile, keyFile)
	if err != nil {
		return nil, fmt.Errorf("failed to save certificate: %w", err)
	}

	fmt.Printf("Generated self-signed certificate: %s, %s\n", certFile, keyFile)
	fmt.Printf("⚠️  This is a self-signed certificate for development only!\n")
	fmt.Printf("⚠️  For production, use a CA-signed certificate!\n")

	return cert, nil
}

// CreateTLSConfig creates TLS config from configuration
func CreateTLSConfig(config *Config) (*tls.Config, error) {
	// Check if SSL/TLS is enabled
	if !config.Enabled {
		return nil, fmt.Errorf("SSL/TLS is not enabled")
	}

	// Load or create certificate
	cert, err := LoadOrCreateCertificate(config.CertFile, config.KeyFile, "MSSQLServer")
	if err != nil {
		return nil, fmt.Errorf("failed to load/create certificate: %w", err)
	}

	// Create TLS config
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{*cert},
		MinVersion:   config.MinVersion,
		MaxVersion:   config.MaxVersion,
		ClientAuth:   config.ClientAuth,
		// Server certificates are not verified by default
		InsecureSkipVerify: true, // For development
	}

	// Set cipher suites (TLS 1.2/1.3)
	tlsConfig.CipherSuites = []uint16{
		tls.TLS_ECDHE_RSA_WITH_AES_128_GCM_SHA256,
		tls.TLS_ECDHE_RSA_WITH_AES_256_GCM_SHA384,
		tls.TLS_ECDHE_ECDSA_WITH_AES_128_GCM_SHA256,
		tls.TLS_ECDHE_ECDSA_WITH_AES_256_GCM_SHA384,
	}

	// Enable TLS 1.3 cipher suites
	if config.MaxVersion >= tls.VersionTLS13 {
		tlsConfig.CipherSuites = append(tlsConfig.CipherSuites,
			tls.TLS_AES_128_GCM_SHA256,
			tls.TLS_AES_256_GCM_SHA384,
		)
	}

	return tlsConfig, nil
}

// CreateTLSListener creates a TLS listener from configuration
func CreateTLSListener(addr string, config *Config) (net.Listener, error) {
	// Create TLS config
	tlsConfig, err := CreateTLSConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create TLS config: %w", err)
	}

	// Create TLS listener
	listener, err := tls.Listen("tcp", addr, tlsConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create TLS listener: %w", err)
	}

	return listener, nil
}

// GetEncryptionLevel returns encryption level based on configuration
func GetEncryptionLevel(config *Config) byte {
	if !config.Enabled {
		return EncryptionOff
	}

	if config.ForceEncryption {
		return EncryptionRequired
	}

	return EncryptionOn
}

// IsEncryptionEnabled checks if encryption is enabled
func IsEncryptionEnabled(config *Config) bool {
	return config.Enabled
}

// ValidateConfig validates SSL/TLS configuration
func ValidateConfig(config *Config) error {
	// Check if SSL/TLS is enabled
	if !config.Enabled {
		return nil
	}

	// Validate certificate file
	if config.CertFile == "" {
		return fmt.Errorf("certificate file path is empty")
	}

	if _, err := os.Stat(config.CertFile); err != nil {
		return fmt.Errorf("certificate file not found: %w", err)
	}

	// Validate key file
	if config.KeyFile == "" {
		return fmt.Errorf("key file path is empty")
	}

	if _, err := os.Stat(config.KeyFile); err != nil {
		return fmt.Errorf("key file not found: %w", err)
	}

	// Validate TLS versions
	if config.MinVersion < tls.VersionTLS12 {
		return fmt.Errorf("minimum TLS version must be at least TLS 1.2")
	}

	if config.MaxVersion < config.MinVersion {
		return fmt.Errorf("maximum TLS version must be >= minimum TLS version")
	}

	return nil
}

// GenerateCertificatePEM generates certificate and private key in PEM format
func GenerateCertificatePEM(commonName string, validityDays int) (certPEM, keyPEM []byte, err error) {
	// Generate private key
	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to generate private key: %w", err)
	}

	// Create certificate template
	template := x509.Certificate{
		SerialNumber: big.NewInt(1),
		Subject: pkix.Name{
			Organization: []string{"MSSQL TDS Server"},
			CommonName:   commonName,
		},
		NotBefore:             time.Now(),
		NotAfter:              time.Now().AddDate(0, 0, validityDays),
		KeyUsage:              x509.KeyUsageKeyEncipherment | x509.KeyUsageDigitalSignature,
		ExtKeyUsage:           []x509.ExtKeyUsage{x509.ExtKeyUsageServerAuth},
		BasicConstraintsValid: true,
	}

	// Generate certificate
	certDER, err := x509.CreateCertificate(rand.Reader, &template, &template, &privateKey.PublicKey, privateKey)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to create certificate: %w", err)
	}

	// Encode certificate in PEM format
	certPEM = pem.EncodeToMemory(&pem.Block{
		Type:  "CERTIFICATE",
		Bytes: certDER,
	})

	// Encode private key in PEM format
	privateKeyBytes, err := x509.MarshalPKCS8PrivateKey(privateKey)
	if err != nil {
		return nil, nil, fmt.Errorf("failed to marshal private key: %w", err)
	}

	keyPEM = pem.EncodeToMemory(&pem.Block{
		Type:  "PRIVATE KEY",
		Bytes: privateKeyBytes,
	})

	return certPEM, keyPEM, nil
}
