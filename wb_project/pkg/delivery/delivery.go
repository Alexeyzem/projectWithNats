package delivery

type Delivery struct {
	Name    string `json: name`
	Phone   string `json: phone`
	Zip     string `json: zip`
	City    string `json: city`
	Address string `json: addres`
	Region  string `json: region`
	Email   string `json: email`
}
