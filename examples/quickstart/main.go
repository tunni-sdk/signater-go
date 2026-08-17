// Command quickstart runs a complete Signater signing flow: upload a PDF,
// create an envelope with one signer, publish it, and poll until the status
// changes.
//
// Usage:
//
//	export SIGNATER_API_TOKEN=<your sandbox token>
//	go run ./examples/quickstart -pdf contract.pdf -vault <vault-id> -signer-email jane@example.com
//
// Omit -vault to create a personal vault. Publishing consumes an envelope
// credit outside sandbox.
package main

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"time"

	signater "github.com/tunni-sdk/signater-go"
)

func main() {
	pdfPath := flag.String("pdf", "", "path to the PDF to send for signing (required)")
	vaultID := flag.String("vault", "", "vault id (created automatically when empty)")
	signerName := flag.String("signer-name", "Jane Doe", "signer name")
	signerEmail := flag.String("signer-email", "", "signer email (required)")
	flag.Parse()
	if *pdfPath == "" || *signerEmail == "" {
		flag.Usage()
		os.Exit(2)
	}

	client := signater.NewClient() // reads SIGNATER_API_TOKEN
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Minute)
	defer cancel()

	if *vaultID == "" {
		id, err := client.Vaults.Create(ctx, &signater.VaultParams{
			Name: "Quickstart Vault",
			Type: signater.VaultTypeUserAccount,
		})
		if err != nil {
			log.Fatalf("creating vault: %v", err)
		}
		*vaultID = id
		fmt.Printf("created vault %s\n", id)
	}

	file, err := os.Open(*pdfPath)
	if err != nil {
		log.Fatal(err)
	}
	defer file.Close()
	doc, err := client.Documents.Upload(ctx, filepath.Base(*pdfPath), file)
	if err != nil {
		log.Fatalf("uploading document: %v", err)
	}
	fmt.Printf("uploaded document %s (%d bytes, %d pages)\n", doc.ID, doc.OriginalFileSize, len(doc.PageSizes))

	envID, err := client.Envelopes.Create(ctx, &signater.EnvelopeParams{
		VaultID:                 *vaultID,
		Name:                    "Quickstart Envelope",
		JurisdictionCountryCode: "BR",
		Language:                signater.LanguagePtBR,
		Signers: []signater.EnvelopeSignerParams{{
			Name:                      *signerName,
			Email:                     signater.String(*signerEmail),
			EmailCommunicationMode:    signater.CommunicationModeFull,
			SmsCommunicationMode:      signater.CommunicationModeNone,
			WhatsAppCommunicationMode: signater.CommunicationModeNone,
		}},
		Documents: []signater.EnvelopeDocumentParams{{
			ID:   doc.ID,
			Name: filepath.Base(*pdfPath),
		}},
	})
	if err != nil {
		log.Fatalf("creating envelope: %v", err)
	}
	fmt.Printf("created envelope %s\n", envID)

	if err := client.Envelopes.Publish(ctx, envID); err != nil {
		var payErr *signater.PaymentRequiredError
		if errors.As(err, &payErr) {
			log.Fatalf("publishing requires credits (should buy: %s) — see https://app.signater.com/billing/plans", payErr.ShouldBuy)
		}
		log.Fatalf("publishing envelope: %v", err)
	}
	fmt.Println("published — the signer will receive an email")

	for {
		env, err := client.Envelopes.Get(ctx, envID)
		if err != nil {
			log.Fatalf("polling envelope: %v", err)
		}
		fmt.Printf("status: %s\n", env.Status)
		if env.Status != signater.EnvelopeStatusPublished {
			break
		}
		select {
		case <-ctx.Done():
			fmt.Println("stopping poll; check the envelope in the Signater app")
			return
		case <-time.After(15 * time.Second):
		}
	}
}
