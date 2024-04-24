package cert_manager

import (
	"crypto"
	"crypto/ecdsa"
	"crypto/elliptic"
	"crypto/rand"
	"encoding/json"
	"fmt"
	"os"

	"github.com/go-acme/lego/challenge/http01"
	"github.com/go-acme/lego/challenge/tlsalpn01"
	"github.com/go-acme/lego/v4/certcrypto"
	"github.com/go-acme/lego/v4/certificate"
	"github.com/go-acme/lego/v4/lego"
	"github.com/go-acme/lego/v4/registration"
)

var CACHE_DIR = "./stewel_cache"

type AcmeUser struct {
	Email string `json:email`
	Registration *registration.Resource `json:registration`
	Key crypto.PrivateKey `json:key`
}

func (u *AcmeUser) GetEmail() string {
	return u.Email
}
func (u AcmeUser) GetRegistration() *registration.Resource {
	return u.Registration
}
func (u *AcmeUser) GetPrivateKey() crypto.PrivateKey {
	return u.Key
}

func ensureCacheDir() {
	file, err := os.Open(CACHE_DIR)
	if err != nil {
		err = os.Mkdir(CACHE_DIR, os.ModePerm)
		if err != nil {
			panic(err.Error())
		}
	} else {
		file.Close()
	}


	file, err = os.Open(CACHE_DIR + "/certs")
	if err != nil {
		err = os.Mkdir(CACHE_DIR + "/certs", os.ModePerm)
		if err != nil {
			panic(err.Error())
		}
	} else {
		file.Close()
	}
}

func writeUserFile(user AcmeUser) AcmeUser {
	ensureCacheDir()

	data, err := json.Marshal(user)
	if err != nil {
		panic(err.Error())
	}
	os.WriteFile(CACHE_DIR + "/" + "acme_user.json", data, os.ModeAppend)

	return user
}

func readUserFile () (AcmeUser) {
	ensureCacheDir()
	file, err := os.Open(CACHE_DIR + "/" + "acme_user.json")
	if err != nil {
		return writeUserFile(AcmeUser{})
	}
	defer file.Close()

	var user AcmeUser
	err = json.NewDecoder(file).Decode(&user)
	if err != nil {
		panic(err.Error())
	}
	return user
}

func Generate(email string, hostName string, target string) (cert string, key string) {

	user := readUserFile()
	if user.Email == "" {
		user.Email = email
		key, err := ecdsa.GenerateKey(elliptic.P256(), rand.Reader)
		if err != nil {
			panic(err.Error())
		}
		user.Key = key
	}

	legoConfig := lego.NewConfig(&user)
	legoConfig.CADirURL = string(target)
	legoConfig.Certificate.KeyType = certcrypto.RSA2048

	client, err := lego.NewClient(legoConfig)
	if err != nil {
		panic(err.Error())
	}
	
	err = client.Challenge.SetHTTP01Provider(http01.NewProviderServer("", "80"))
	if err != nil {
		panic(err.Error())
	}
	err = client.Challenge.SetTLSALPN01Provider(tlsalpn01.NewProviderServer("", "443"))
	if err != nil {
		panic(err.Error())
	}

	if user.Registration.Body.Status != "valid" {
		reg, err := client.Registration.Register(registration.RegisterOptions{
			TermsOfServiceAgreed: true,
		})

		if err != nil {
			panic(err.Error())
		}

		user.Registration = reg

		writeUserFile(user)
	}

	request := certificate.ObtainRequest{
		Domains: []string{hostName},
		Bundle: true,
	}

	certificates, err := client.Certificate.Obtain(request)
	if err != nil {
		panic(err.Error())
	}

	fmt.Printf("%#v\n", certificates)

	createAndWriteFile(CACHE_DIR + "/certs" + hostName + "_certificate", certificates.Certificate)
	createAndWriteFile(CACHE_DIR + "/certs" + hostName + "_private_key", certificates.PrivateKey)
	createAndWriteFile(CACHE_DIR + "/certs" + hostName + "_iss_certificate", certificates.IssuerCertificate)
	createAndWriteFile(CACHE_DIR + "/certs" + hostName + "_csr", certificates.CSR)
	
	return CACHE_DIR + "/certs" + hostName + "_certificate", CACHE_DIR + "/certs" + hostName + "_private_key"
}

func createAndWriteFile(name string, data []byte) {
	f, err := os.Create(name)
    if err != nil {
        panic(err)
    }
	defer f.Close()
	f.Write(data)
}