package snowflake

import (
	"time"

	sf "github.com/bwmarrin/snowflake"
)

// 定义一个全局节点
var node *sf.Node

// Init 初始化雪花算法节点
// startTime: 项目上线时间 (格式: "2006-01-02")，比如 "2026-01-01"
// machineID: 机器ID (0-1023)，在分布式部署时，每台服务器必须不同！
func Init(startTime string, machineID int64) (err error) {
	// 1. 定制开始时间 (这一步很重要！)
	// 如果不设置，默认是从 2010 年开始算，你就会损失 16 年的有效期
	var st time.Time
	st, err = time.Parse("2006-01-02", startTime)
	if err != nil {
		return
	}

	// 设置库的 Epoch (起始时间) 为你指定的时间
	sf.Epoch = st.UnixNano() / 1000000

	// 2. 创建节点
	node, err = sf.NewNode(machineID)
	return
}

// GenID 生成一个 int64 的分布式唯一 ID
func GenID() int64 {
	return node.Generate().Int64()
}
