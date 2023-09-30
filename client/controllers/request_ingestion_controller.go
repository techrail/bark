package controllers

import (
	"github.com/techrail/bark/client/services/ingestion"
	"github.com/techrail/bark/models"
)

func SendSingleToClientChannel(l models.BarkLog) {

	go ingestion.InsertSingleRequest(l)
}
