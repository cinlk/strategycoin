package contract

type BasisResponse struct {
	Ch     string       `json:"ch"`
	Status string       `json:"status"`
	Ts     int64        `json:"ts"`
	Data   []basisDatas `json:"data"`
}

type basisDatas struct {
	Id            int64  `json:"id"`
	ContractPrice string `json:"contract_price"`
	Basis         string `json:"basis"`
	IndexPrice    string `json:"index_price"`
	BasisRate     string `json:"basis_rate"`
}
