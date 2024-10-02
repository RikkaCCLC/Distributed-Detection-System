package common

type ResourceMountReq struct {
	ResourceType string  `json:"resource_type" binding:"required"` // 资源的类型 host rds dcs
	ResourceIds  []int64 `json:"resource_ids" binding:"required"`  // 要操作的资源id列表
	TargetPath   string  `json:"target_path" binding:"required"`   // 目标 g.p.a
}
