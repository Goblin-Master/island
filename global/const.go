package global

import (
	"tgwp/util/snowflake"
	"time"
)

// 所有常量文件读取位置
const (
	DEFAULT_CONFIG_FILE_PATH = "/config.yaml"
	ATOKEN_EFFECTIVE_TIME    = time.Hour * 12
	RTOKEN_EFFECTIVE_TIME    = time.Hour * 24 * 30
	AUTH_ENUMS_ATOKEN        = "atoken"
	AUTH_ENUMS_RTOKEN        = "rtoken"
	DEFAULT_NODE_ID          = 1
	TOKEN_USER_ID            = "UserId"
)

var Node, _ = snowflake.NewNode(DEFAULT_NODE_ID)