package connnect

type MrdStateEvent struct {
	Hidden   bool  `json:"hidden"`
	ExpireTs int64 `json:"expire_ts"`
}

type DrStateEvent struct {
	Timeframe int64  `json:"timeframe"`
	At        string `json:"at"`
}
