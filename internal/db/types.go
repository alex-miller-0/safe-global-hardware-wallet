package db

type Config struct {
	DbPath string `json:"db_path"`
}

type Db struct {
	Safes []Safe       `json:"safes"`
	Tags  []AddressTag `json:"tags"`
}

type Safe struct {
	ID        AddressTag `json:"id"`
	Network   string     `json:"network"`
	Owners    []string   `json:"owners"`
	Threshold int        `json:"threshold"`
}

type AddressTag struct {
	Address string `json:"address"`
	Tag     string `json:"tag"`
}
