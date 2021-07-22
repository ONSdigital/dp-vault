package main

import (
	"errors"
	"os"
	"time"

	vault "github.com/ONSdigital/dp-vault"
	"github.com/ONSdigital/log.go/v2/log"
)

const maxRetries = 3

func main() {

	log.Namespace = "vault-example-v2"
	devAddress := os.Getenv("VAULT_ADDR")
	token := os.Getenv("VAULT_TOKEN")
	prefix := "secret/data/"
	path := prefix + "shared/datasets/CPIH-0002"
	key := "Key2"
	val := time.Now().Format("2006-01-02 15:04:05")

	// In production no tokens should be logged
	logData := log.Data{
		"address": devAddress,
		"token":   token,
		"path":    path,
		"key":     key,
		"val":     val,
	}

	client, err := vault.CreateClient(token, devAddress, maxRetries)
	if err != nil {
		log.Fatal(nil, "failed to connect to vault", err, logData)
		os.Exit(1)
	}

	log.Info(nil, "Created vault client", logData)

	if err = client.VWriteKey(path, key, val); err != nil {
		log.Fatal(nil, "failed to write to vault", err, logData)
		os.Exit(1)
	}

	readVal, ver, err := client.VReadKey(path, key)
	if err != nil {
		if err == vault.ErrKeyNotFound {
			log.Error(nil, "key not in vault", err, logData)
		} else {
			log.Error(nil, "failed to read PK-Key from vault", err, logData)
		}
		os.Exit(1)
	}

	logData["read_val"] = readVal
	logData["read_version"] = ver

	if readVal != val {
		err = errors.New("read value differs from expected")
		log.Fatal(nil, "", err, logData)
		os.Exit(1)
	}

	log.Info(nil, "successfully written and read from vault", logData)
}
