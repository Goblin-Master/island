package global

import (
	"tgwp/utils/snowflake"
	"time"
)

// 所有常量文件读取位置
const (
	DEFAULT_CONFIG_FILE_PATH = "/config.yaml"
	ATOKEN_EFFECTIVE_TIME    = time.Hour * 2
	RTOKEN_EFFECTIVE_TIME    = time.Hour * 24 * 7
	AUTH_ENUMS_ATOKEN        = "atoken"
	AUTH_ENUMS_RTOKEN        = "rtoken"
	DEFAULT_NODE_ID          = 1
	TOKEN_USER_ID            = "UserId"
	IMAGE_PATH               = "uploads"       //自定义图片路径
	IMAGE_SIZE               = 1024 * 1024 * 3 //自定义图片大小

	QUESTION_TYPE_CHOICE     = 1 // 选择题
	QUESTION_TYPE_MULTIPLE   = 2 // 多选题
	QUESTION_TYPE_BLANK      = 3 // 填空题
	QUESTION_TYPE_SUBJECTIVE = 4 // 主观题
)

var Node, _ = snowflake.NewNode(DEFAULT_NODE_ID)
