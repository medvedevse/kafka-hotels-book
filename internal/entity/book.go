package entity

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

type HotelBook struct {
	Id        string    `json:"id"`
	HotelName string    `json:"hotelName"`
	DateStart time.Time `json:"dateStart"`
	DateEnd   time.Time `json:"dateEnd"`
}

func NewHotelBook(hotelName string, dateEnd time.Time) *HotelBook {
	return &HotelBook{
		Id:        uuid.NewString(),
		HotelName: hotelName,
		DateStart: time.Now(),
		DateEnd:   dateEnd,
	}
}

func (h *HotelBook) ToString() string {
	data, err := json.Marshal(h)
	if err != nil {
		return "failed to marshal HotelBook to JSON"
	}
	return string(data)
}
