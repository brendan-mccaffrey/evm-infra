package gas

const EthProviderUrl = "https://api.gasprice.io/v1/estimates"

type GasProvider interface {
	Start() error
	Instant() *GasStruct
	Fast() *GasStruct
	Eco() *GasStruct
	BaseFee() float64
	NativePrice() float64
}

type GasStruct struct {
	FeeCap          float64 `json:"feeCap"`
	MaxPrioirityFee float64 `json:"maxPrioirityFee"`
}
