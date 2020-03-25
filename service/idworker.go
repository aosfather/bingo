package service

import (
	"fmt"
	"sync"
	"time"
)

/**
  分布式ID生成器
*/
const (
	_twepoch            = 1546272000000                                     //开始时间截 (2019-01-01)
	_workerIdBits       = 5                                                 //机器id所占的位数
	_datacenterIdBits   = 5                                                 //数据中心标识id所占的位数
	_maxWorkerId        = -1 ^ (-1 << _workerIdBits)                        //支持的最大机器id，结果是31,这个移位算法可以很快的计算出几位二进制数所能表示的最大十进制数
	_maxDatacenterId    = -1 ^ (-1 << _datacenterIdBits)                    // 支持的最大数据标识id，结果是31
	_sequenceBits       = 22                                                //序列在id中占的位数
	_workerIdShift      = _sequenceBits                                     //机器ID向左移22位
	_datacenterIdShift  = _sequenceBits + _workerIdBits                     //数据标识id向左移27位(22+5)
	_timestampLeftShift = _sequenceBits + _workerIdBits + _datacenterIdBits //时间截向左移32位(5+5+22)
	_sequenceMask       = -1 ^ (-1 << _sequenceBits)                        //生成序列的掩码，这里为4194303

	_formatMask = "20060102150405"
)

type IdWorker struct {
	WorkerId      int64 //生成器的序号
	DataCenterId  int64 //数据中心的序号
	sequence      int64 //秒内序列(0~4194303)
	lastTimestamp int64 //上次生成ID的时间截
	mutex         sync.Mutex
}

//创建新的worker
func CreateWorker(worker int64, datacenter int64) *IdWorker {
	if worker > _maxWorkerId || worker < 0 {
		//"worker Id can't be greater than %d or less than 0"
		return nil
	}
	if datacenter > _maxDatacenterId || datacenter < 0 {
		//"datacenter Id can't be greater than %d or less than 0"
		return nil
	}

	w := IdWorker{}
	w.WorkerId = worker
	w.DataCenterId = datacenter
	return &w
}

//带前缀的Id
func (this *IdWorker) NextIdWithPrefix(h string) string {
	return fmt.Sprintf("%s%d", h, this.NextId())
}

func (this *IdWorker) NextIdWithTime() string {
	return fmt.Sprintf("%s%d", time.Now().Format(_formatMask), this.NextId())
}

//获取下一个ID，方法需要做到线程安全
func (this *IdWorker) NextId() int64 {
	this.mutex.Lock()
	timestamp := this.timeGen()
	// 如果当前时间小于上一次ID生成的时间戳，说明系统时钟回退过这个时候应当抛出异常
	if timestamp < this.lastTimestamp {
	}

	// 如果是同一时间生成的，则进行秒内序列
	if timestamp == this.lastTimestamp {
		this.sequence = (this.sequence + 1) & _sequenceMask

		if this.sequence == 0 {
			timestamp = this.tilNextMillis(this.lastTimestamp)
		}
	} else {
		this.sequence = 0
	}

	this.lastTimestamp = timestamp
	defer this.mutex.Unlock()

	return ((timestamp - _twepoch) << _timestampLeftShift) | (this.DataCenterId << _datacenterIdShift) | (this.WorkerId << _workerIdShift) | this.sequence

}

//阻塞到下一个秒，直到获得新的时间戳
func (this *IdWorker) tilNextMillis(lastTimestamp int64) int64 {
	timestamp := this.timeGen()
	for {
		if timestamp > lastTimestamp {
			break
		}
		timestamp = this.timeGen()
	}
	return timestamp

}

//当前秒
func (this *IdWorker) timeGen() int64 {
	return time.Now().Unix()
}
