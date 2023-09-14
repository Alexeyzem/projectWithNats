package delivery

type Delivery struct {
	ID      int    `json:"id"`
	UserId  int    `json:"user_id"`
	Name    string `json:"name"`
	Phone   string `json:"phone"`
	Zip     string `json:"zip"`
	City    string `json:"city"`
	Address string `json:"addres"`
	Region  string `json:"region"`
	Email   string `json:"email"`
}
