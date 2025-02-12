package qr_services

import (
	"bytes"
	"image/png"

	"github.com/boombuler/barcode"
	"github.com/boombuler/barcode/qr"
)

type QRService interface {
	GenerateQRCode(restaurantID string) ([]byte, error)
}

type qrService struct {
}

func NewQRService() QRService {
	return &qrService{}
}

func (q *qrService) GenerateQRCode(restaurantID string) ([]byte, error) {
	url := "https://meal-map-2.vercel.app/restaurant/" + restaurantID + "/menu"

	qrCode, err := qr.Encode(url, qr.L, qr.Auto)
	if err != nil {
		return nil, err
	}

	qrCode, err = barcode.Scale(qrCode, 400, 400)
	if err != nil {
		return nil, err
	}

	var buf bytes.Buffer
	err = png.Encode(&buf, qrCode)
	if err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}
