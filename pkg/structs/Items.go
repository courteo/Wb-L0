package structs

type Items struct {
	ChrtID      uint64 	`json:"chrt_id"`
	TrackNumber string	`json:"track_number"`
	Price       uint64	`json:"price"`
	RID         string	`json:"rid"`
	Name        string	`json:"name"`
	Sale        uint64	`json:"sale"`
	Size        string	`json:"size"`
	TotalPrice  uint64	`json:"tota_price"`
	NmID        uint64	`json:"nm_id"`
	Brand       string	`json:"brand"`
	Status      uint8	`json:"status"`
}