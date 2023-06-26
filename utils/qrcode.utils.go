package utils

import "github.com/skip2/go-qrcode"

func GenerateQRCode(data string) (string, error) {
	qrFile := "qrcode.png"
	err := qrcode.WriteFile(data, qrcode.Medium, 256, qrFile)
	if err != nil {
		return "", err
	}
	return qrFile, nil
}
