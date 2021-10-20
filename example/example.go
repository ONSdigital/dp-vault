package main

import (
	"context"
	"os"

	vault "github.com/ONSdigital/dp-vault/v2"
	"github.com/ONSdigital/log.go/v2/log"
)

const maxRetries = 3

func main() {

	log.Namespace = "vault-example"
	devAddress := os.Getenv("VAULT_ADDR")
	token := os.Getenv("VAULT_TOKEN")
	ctx := context.Background()

	client, err := vault.CreateClient(token, devAddress, maxRetries)

	// In production no tokens should be logged
	logData := log.Data{"address": devAddress, "token": token}
	log.Info(ctx, "Created vault client", logData)

	if err != nil {
		log.Error(ctx, "failed to connect to vault", err, logData)
	}

	err = client.WriteKey("secret/shared/datasets/CPIH-0000", "PK-Key", "098980474948463874535354")

	if err != nil {
		log.Error(ctx, "failed to write to vault", err, logData)
	}

	PKKey, err := client.ReadKey("secret/shared/datasets/CPIH-0000", "PK-Key")

	if err != nil {
		log.Error(ctx, "failed to read  PK Key from vault", err, logData)
	}
	logData["pk-key"] = PKKey
	log.Info(ctx, "successfully  written and read a key from vault", logData)
}
