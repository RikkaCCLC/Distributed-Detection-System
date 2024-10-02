package common

// 单一探测结果
type ProberResultOne struct {
	WorkerName   string  `json:"workerName"`
	MetricName   string  `json:"metricName"`
	TargetAddr   string  `json:"targetAddr"`
	SourceRegion string  `json:"sourceRegion"`
	TargetRegion string  `json:"targetRegion"`
	ProbeType    string  `json:"probeType"`
	TimeStamp    int64   `json:"timeStamp"`
	Value        float32 `json:"value"`
}

// 推送本地的探测结果的请求
type ProberResultPushRequest struct {
	ProberResults []*ProberResultOne `json:"prober_results,omitempty"`
}

// 推送结果的响应
type ProberResultPushResponse struct {
	SuccessNum int32 `json:"success_num,omitempty"`
}

// 获取探测目标的请求
type ProberTargetsGetRequest struct {
	LocalRegion string `json:"localRegion"`
	LocalIp     string `json:"localIp"`
}

// 获取探测的目标的响应
type ProberTargetsGetResponse struct {
	Targets []*ProberTargets `json:"targets"`
}

// Targets
type ProberTargets struct {
	ProberType string   `json:"proberType"`
	Region     string   `json:"region"`
	Target     []string `json:"target"`
}
