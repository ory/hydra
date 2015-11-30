package handler

/*
 * Generates a private/public key pair in PEM format (not Certificate)
 *
 * The generated private key can be parsed with openssl as follows:
 * > openssl rsa -in key.pem -text
 *
 * The generated public key can be parsed as follows:
 * > openssl rsa -pubin -in pub.pem -text
 */

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"github.com/codegangsta/cli"
	"io/ioutil"
	"log"
)

func CreatePublicPrivatePEMFiles(ctx *cli.Context) {
	// priv *rsa.PrivateKey;
	// err error;
	priv, err := rsa.GenerateKey(rand.Reader, 1024)
	if err != nil {
		fmt.Println(err)
		return
	}
	err = priv.Validate()
	if err != nil {
		fmt.Println("Validation failed.", err)
	}

	// Get der format. priv_der []byte
	privDer := x509.MarshalPKCS1PrivateKey(priv)

	// pem.Block
	// blk pem.Block
	privBlk := pem.Block{
		Type:    "RSA PRIVATE KEY",
		Headers: nil,
		Bytes:   privDer,
	}

	// Resultant private key in PEM format.
	// priv_pem string
	privPEM := string(pem.EncodeToMemory(&privBlk))

	// Public Key generation

	pub := priv.PublicKey
	pubDer, err := x509.MarshalPKIXPublicKey(&pub)
	if err != nil {
		fmt.Println("Failed to get der format for PublicKey.", err)
		return
	}

	pubBlk := pem.Block{
		Type:    "PUBLIC KEY",
		Headers: nil,
		Bytes:   pubDer,
	}
	pubPEM := string(pem.EncodeToMemory(&pubBlk))

	err = ioutil.WriteFile(ctx.String("public-file-path"), []byte(pubPEM), 0644)
	if err != nil {
		log.Fatalf("%v", err)
	}

	err = ioutil.WriteFile(ctx.String("private-file-path"), []byte(privPEM), 0644)
	if err != nil {
		log.Fatalf("%v", err)
	}
}
