package main

import (
	"fmt"
	"log"

	"ichibuy/auth/internal/domain"
)

func main() {
	privateKey, err := domain.GenerateRSAKeyPair()
	if err != nil {
		log.Fatal("Failed to generate RSA key pair:", err)
	}

	privateKeyPEM := domain.PrivateKeyToPEM(privateKey)
	publicKeyPEM, err := domain.PublicKeyToPEM(&privateKey.PublicKey)
	if err != nil {
		log.Fatal("Failed to convert public key to PEM:", err)
	}

	fmt.Println("=== RSA PRIVATE KEY (for JWT_PRIVATE_KEY env var) ===")
	fmt.Println(privateKeyPEM)

	fmt.Println("=== RSA PUBLIC KEY (for verification) ===")
	fmt.Println(publicKeyPEM)

	fmt.Println("=== USAGE ===")
	fmt.Println("Copy the private key above and set it as JWT_PRIVATE_KEY environment variable.")
	fmt.Println("The public key will be automatically exposed via the JWKS endpoint.")
}