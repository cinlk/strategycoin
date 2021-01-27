package contract

type ServerHeatBeatInfo struct {
	Status string `json:"status"`
	Ts int64 `json:"ts"`
	Data *heartBeatData `json:"data"`
}


type heartBeatData struct {

	Heartbeat int `json:"heartbeat"`
	SwapHeartbeat int `json:"swap_heartbeat"`
	EstimatedRecoveryTime int64 `json:"estimated_recovery_time"`
	SwapEstimatedRecovery_time int64 `json:"swap_estimated_recovery_time"`
	OptionHeartbeat int `json:"option_heartbeat"`
	OptionEstimatedRecoveryTime int64 `json:"option_estimated_recovery_time"`
	LinearSwapHeartbeat int64 `json:"linear_swap_heartbeat"`
	LinearSwapEstimatedRecoveryTime int64 `json:"linear_swap_estimated_recovery_time"`
}
