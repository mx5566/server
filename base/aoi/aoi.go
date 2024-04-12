package aoi

type AOI struct {
	x        int32    // x方向格子坐标
	y        int32    // y方向格子坐标
	unit     IUnit    // aoi目标对象
	realNode *AoiNode // 实际的节点指针
}

type IUnit interface {
	OnEnterAoi(aoi *AOI)
	OnLeaveAoi(aoi *AOI)
}

func InitAoi(aoi *AOI, u IUnit) {
	aoi.unit = u
}
