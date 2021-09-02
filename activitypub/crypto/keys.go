package crypto

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"fmt"
	"net/url"

	"github.com/owncast/owncast/core/data"
	log "github.com/sirupsen/logrus"
)

func GetPublicKey(actorIRI *url.URL) PublicKey {
	key := data.GetPublicKey()
	idURL, err := url.Parse(actorIRI.String() + "#main-key")
	if err != nil {
		log.Errorln("unable to parse actor iri string", idURL, err)
	}

	return PublicKey{
		Id:           idURL,
		Owner:        actorIRI,
		PublicKeyPem: key,
	}
}

func GetPrivateKey() *rsa.PrivateKey {
	key := data.GetPrivateKey()

	block, _ := pem.Decode([]byte(key))
	if block == nil {
		log.Errorln(errors.New("failed to parse PEM block containing the key"))
		return nil
	}

	priv, err := x509.ParsePKCS1PrivateKey(block.Bytes)
	if err != nil {
		log.Errorln("unable to parse private key", err)
		return nil
	}

	return priv
}

func GenerateKeys() ([]byte, []byte, error) {
	// generate key
	privatekey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		fmt.Printf("Cannot generate RSA key\n")
		return nil, nil, err
	}
	publickey := &privatekey.PublicKey

	var privateKeyBytes []byte = x509.MarshalPKCS1PrivateKey(privatekey)
	privateKeyBlock := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyBytes,
	}
	privatePem := pem.EncodeToMemory(privateKeyBlock)

	publicKeyBytes, err := x509.MarshalPKIXPublicKey(publickey)
	if err != nil {
		fmt.Printf("error when dumping publickey: %s \n", err)
		return nil, nil, err
	}
	publicKeyBlock := &pem.Block{
		Type:  "PUBLIC KEY",
		Bytes: publicKeyBytes,
	}
	publicPem := pem.EncodeToMemory(publicKeyBlock)

	return privatePem, publicPem, nil
}
