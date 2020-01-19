package service

/**
  分布式ID生成器
*/

type IdWorker struct {
	WorkerId     int //生成器的序号
	DataCenterId int //数据中心的序号
}
