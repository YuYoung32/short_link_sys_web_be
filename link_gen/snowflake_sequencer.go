/**
 * Created by YuYoung on 2023/4/17
 * Description: 雪花算法自增
 */

package link_gen

import (
	"github.com/bwmarrin/snowflake"
	"short_link_sys_web_be/log"
)

var node *snowflake.Node

func SnowflakeInit() {
	var err error
	node, err = snowflake.NewNode(1)
	if err != nil {
		log.GetLogger().Error("Snowflake init failed")
		return
	}
}

type SnowflakeSequencer struct{}

func (SnowflakeSequencer) GenLink(_ string) string {
	return uint64ToShortLink(uint64(node.Generate().Int64()))
}

func (SnowflakeSequencer) GetType() AlgorithmType {
	return SeqType
}
