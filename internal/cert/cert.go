package cert

import (
	"crypto/tls"
	"encoding/pem"
	"fmt"
	"log"
	"os"
	"os/user"
	"path"
)

func GeneratorCert(serverAddress, configPath string) error {
	serverAddress, err := extractHostAndPort(serverAddress)
	if err != nil {
		return err
	}
	// The HTTPS server's address (replace with the server you want to check)
	// serverAddress := "47.96.96.106:8443"

	// Create a tls.Config with InsecureSkipVerify set to true
	tlsConfig := &tls.Config{
		InsecureSkipVerify: true,
	}

	// Dial the server using TLS with the custom tls.Config
	conn, err := tls.Dial("tcp", serverAddress, tlsConfig)
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	// Retrieve the certificate chain
	certs := conn.ConnectionState().PeerCertificates

	// Print the certificate chain
	for i, cert := range certs {
		fmt.Printf("Certificate %d:\n", i)
		fmt.Printf("  Subject: %s\n", cert.Subject)
		fmt.Printf("  Issuer: %s\n", cert.Issuer)
		fmt.Printf("  Not Before: %s\n", cert.NotBefore)
		fmt.Printf("  Not After: %s\n", cert.NotAfter)
	}

	// The root CA is usually the last certificate in the chain
	rootCert := certs[len(certs)-1]
	fmt.Println("\nRoot CA Certificate:")
	fmt.Printf("  Subject: %s\n", rootCert.Subject)
	fmt.Printf("  Issuer: %s\n", rootCert.Issuer)
	fmt.Printf("  Not Before: %s\n", rootCert.NotBefore)
	fmt.Printf("  Not After: %s\n", rootCert.NotAfter)

	// Save the server certificate to a file
	certOut, err := os.Create(configPath)
	if err != nil {
		return err
	}
	defer certOut.Close()

	pem.Encode(certOut, &pem.Block{Type: "CERTIFICATE", Bytes: certs[0].Raw})
	fmt.Println("Server certificate saved to server.crt")
	return nil
}

func GetCertPath() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}

	configDir := path.Join(usr.HomeDir, ".cli")
	configPath := path.Join(configDir, "server.crt")

	// 确保配置目录存在
	if err := os.MkdirAll(configDir, 0700); err != nil {
		return "", err
	}
	return configPath, nil
}
