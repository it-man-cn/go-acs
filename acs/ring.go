package main

import (
	log "github.com/it-man-cn/log4go"
)

const (
	// signal command
	SignalNum   = 1
	ProtoFinish = 0
	ProtoReady  = 1
)

//Ring define
//rp和wp双指针可以实现一读一写的两个goroutine间无锁获取和写入数据
//因为writer和reader只会操作自己的wn和rn对象
//我们利用简单的rp==wp就能知道当前ring是emtpy还是full (wp-rp>r.num)
type Ring struct {
	// read
	rp   uint64
	num  uint64 //可以存放Proto结构的个数，即缓存区大小
	mask uint64
	// TODO split cacheline, many cpu cache line size is 64
	// pad [40]byte
	// write
	wp uint64
	//只有一个对象（数组对象），内存连续效率高，以免每次传递数据都创建一个Proto对象
	//每一个Proto都是一个数据包（包含Header和Body)
	data []Message //buffer for reader/writer
}

//NewRing init ring
func NewRing(num int) *Ring {
	r := new(Ring)
	r.init(uint64(num))
	return r
}

//Init init ring
func (r *Ring) Init(num int) {
	r.init(uint64(num))
}

//????  掩码如何计算
func (r *Ring) init(num uint64) {
	// 2^
	//判断是否为2^N，x为2^N，则x的二进制表示为第1位为1其余0，x-1则相反，相与结果为0。
	if num&(num-1) != 0 {
		for num&(num-1) != 0 {
			num &= (num - 1)
		}
		num = num << 1
	}
	r.data = make([]Message, num)
	r.num = num
	r.mask = r.num - 1
}

//获取reader buffer
func (r *Ring) Get() (proto *Message, err error) {
	if r.rp == r.wp { //读写指针为相同的计数，reader buffer为空
		return nil, ErrRingEmpty //????
	}
	proto = &r.data[r.rp&r.mask] //用掩码计算出的是什么？
	return
}

//reader buffer指针偏移，代表已读一个的proto
func (r *Ring) GetAdv() {
	r.rp++
	if Debug {
		log.Debug("ring rp: %d, idx: %d", r.rp, r.rp&r.mask)
	}
}

//获取writer buffer
func (r *Ring) Set() (proto *Message, err error) {
	if r.wp-r.rp >= r.num {
		return nil, ErrRingFull
	}
	proto = &r.data[r.wp&r.mask]
	return
}

func (r *Ring) SetAdv() {
	r.wp++
	if Debug {
		log.Debug("ring wp: %d, idx: %d\n", r.wp, r.wp&r.mask)
	}
}

func (r *Ring) Reset() {
	r.rp = 0
	r.wp = 0
	// prevent pad compiler optimization
	// r.pad = [40]byte{}
}
