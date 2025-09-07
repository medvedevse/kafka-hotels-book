package handler

import "log"

type Handler struct{}

func NewHandler() *Handler {
	return &Handler{}
}

func (h *Handler) HandlerMessage(msg []byte) error {
	log.Printf("Kafka message: %s", string(msg))
	return nil
}
