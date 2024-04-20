package uuid

import (
	"encoding/base64"
	"github.com/mx5566/logm"
	"github.com/mx5566/server/base"
	"strconv"
	"sync"
	"time"
)

// 雪花算法分布式视同全局唯一的ID生成
// 不依赖第三方系统，依赖系统时间，时间回调会出现id碰撞的问题
// +--------------------------------------------------------------------------+
// | 1 Bit Unused | 41 Bit Timestamp |  10 Bit NodeID  |   12 Bit Sequence ID |
// +--------------------------------------------------------------------------+

// 		-1 8位
// 原码	10000001
// 反码  11111110
// 补码  11111111

// 2024-04-19 21:59:02
const Epoch int64 = 1713535142000 // 算法的开始时间
const WorkIDBits = 10             // 机器ID占用的位数
const SequenceIDBits = 12         // 序列ID占用的位数

const WorkIDMax = -1 ^ (-1 << WorkIDBits)
const SequenceIDMax = -1 ^ (-1 << SequenceIDBits)
const WorkIDShit = SequenceIDBits            // 机器ID左移的位数
const TimeShit = SequenceIDBits + WorkIDBits // 时间左移的位数

type SnowFlake struct {
	mu            sync.Mutex
	workId        int64 // 动态生成机器id
	lastTimestamp int64 // 上一次获取id的时间
	sequence      int64 // 当前到那个序列号了
}

func NewSnowFlake() *SnowFlake {
	return &SnowFlake{}
}

func (s *SnowFlake) Init() {

}

func (s *SnowFlake) UUID() int64 {
	defer func() {
		if err := recover(); err != nil {
			base.TraceCode(err)
		}
	}()

	s.mu.Lock()
	currentTime := time.Now().UnixNano() / int64(time.Millisecond)
	if currentTime < s.lastTimestamp {
		if s.lastTimestamp-currentTime <= 5 {
			// 5毫秒以内等待
			for currentTime < s.lastTimestamp {
				currentTime = time.Now().UnixNano() / int64(time.Millisecond)
			}
		}

		logm.PanicfE("时间出现了回退 当前:%d, 上一个时间:%d", currentTime, s.lastTimestamp)
	}

	if currentTime == s.lastTimestamp {
		//相同毫秒内，序列号自增
		s.sequence = (s.sequence + 1) & SequenceIDMax
		//同一毫秒的序列数已经达到最大
		if s.sequence == 0 {
			// 当前时间还在这个时间，继续等待到下一毫秒
			for currentTime == s.lastTimestamp {
				currentTime = time.Now().UnixNano() / int64(time.Millisecond)
			}
		}
	} else {
		//不同毫秒内，序列号置为0
		s.sequence = 0
	}

	s.lastTimestamp = currentTime

	id := ((currentTime - Epoch) << TimeShit) | (s.workId << WorkIDShit) | s.sequence

	s.mu.Unlock()
	return id
}

func Base64(id int64) string {
	return base64.StdEncoding.EncodeToString(Bytes(id))
}

func ParseBase64(id string) (int64, error) {
	b, err := base64.StdEncoding.DecodeString(id)
	if err != nil {
		return -1, err
	}
	return ParseBytes(b)

}

// ParseBytes converts a byte slice into a snowflake ID
func ParseBytes(id []byte) (int64, error) {
	i, err := strconv.ParseInt(string(id), 10, 64)
	return i, err
}

// Bytes returns a byte slice of the snowflake ID
func Bytes(id int64) []byte {
	return []byte(String(id))
}

// String returns a string of the snowflake ID
func String(id int64) string {
	return strconv.FormatInt(id, 10)
}

var UUID = NewSnowFlake()
