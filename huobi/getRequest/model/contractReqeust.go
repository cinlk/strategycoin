package model


type GetContractBasisRequest struct {

	BasisPriceType  string
}

type GetContractAccountRequest struct {
	Symbol string `json:"symbol"`
}


type GetContractCodeRequest struct {
	ContractCode string `json:"contract_code"`
}

