package main

import (
	"bytes"
	"context"
	"crypto/md5"
	"fmt"
	"os"
	"time"

	"github.com/ProtonMail/gopenpgp/v3/crypto"
	"github.com/go-telegram/bot"
	"github.com/go-telegram/bot/models"
)

func (bak *Backup) performBackup(ctx context.Context) {
	b, err := bot.New(bak.Config.Backup.Token)
	if err != nil {
		bak.Logger.Error.Printf("Failed to create bot: %v\n", err)
		return
	}

	pgp := crypto.PGP()

	fileData, err := os.ReadFile("./data/quotobot.db")
	if err != nil {
		bak.Logger.Error.Printf("Read file failed: %v\n", err)
		return
	}

	md5Sum := md5.Sum(fileData)
	bak.Logger.Info.Printf("Read backup file, size: %d bytes, MD5: %x\n", len(fileData), md5Sum)

	password := []byte(bak.Config.Backup.EncryptionPassphrase)

	encHandle, err := pgp.Encryption().Password(password).New()
	if err != nil {
		bak.Logger.Error.Printf("Create encryption handle failed: %v\n", err)
		return
	}

	pgpMessage, err := encHandle.Encrypt(fileData)
	if err != nil {
		bak.Logger.Error.Printf("Encryption failed: %v\n", err)
		return
	}

	armored, err := pgpMessage.ArmorBytes()
	if err != nil {
		bak.Logger.Error.Printf("Armor failed: %v\n", err)
		return
	}

	bak.Logger.Info.Println("Backup file encrypted successfully")

	_, err = b.SendDocument(ctx, &bot.SendDocumentParams{
		ChatID: bak.Config.Backup.ChatID,
		Document: &models.InputFileUpload{
			Filename: fmt.Sprintf("quotobot-backup-%s.db.pgp", time.Now().Format("20060102")),
			Data:     bytes.NewReader(armored),
		},
		Caption: fmt.Sprintf("Date: %s\nMD5: %x", time.Now().Format(time.DateOnly), md5Sum),
	})
	if err != nil {
		bak.Logger.Error.Printf("Failed to send backup: %v\n", err)
		return
	}

	bak.Logger.Info.Println("Backup sent successfully")
}
