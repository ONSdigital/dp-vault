package main

import (
	"os"

	vault "github.com/ONSdigital/dp-vault"
	"github.com/ONSdigital/log.go/log"
)

const maxRetries = 3

func main() {

	log.Namespace = "vault-example"
	devAddress := os.Getenv("VAULT_ADDR")
	token := os.Getenv("VAULT_TOKEN")

	client, err := vault.CreateClient(token, devAddress, maxRetries)

	// In production no tokens should be logged
	logData := log.Data{"address": devAddress, "token": token}
	log.Event(nil, "Created vault client", logData)

	if err != nil {
		log.Event(nil, "failed to connect to vault", logData, log.Error(err))
	}

	err = client.WriteKey("secret/shared/datasets/CPIH-0000", "PK-Key", "098980474948463874535354")

	if err != nil {
		log.Event(nil, "failed to write to vault", logData, log.Error(err))
	}

	PKKey, err := client.ReadKey("secret/shared/datasets/CPIH-0000", "PK-Key")

	if err != nil {
		log.Event(nil, "failed to read  PK Key from vault", logData, log.Error(err))
	}
	logData["pk-key"] = PKKey
	log.Event(nil, "successfully  written and read a key from vault", logData)
}
