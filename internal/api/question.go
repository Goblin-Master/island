package api

import (
	"github.com/gin-gonic/gin"
	"tgwp/log/zlog"
	"tgwp/logic"
	"tgwp/response"
	"tgwp/types"
)

func RunCode(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)
	req, err := types.BindReq[types.RunCodeReq](c)
	if err != nil {
		return
	}
	zlog.CtxInfof(ctx, "RunCode request: %v", req)
	// 默认值
	if req.TimeLimit == 0 {
		req.TimeLimit = 1000
	}
	if req.MemoryLimit == 0 {
		req.MemoryLimit = 512 * 1024 * 1024
	}

	resp, err := logic.NewQuestionLogic().RunCode(ctx, req)
	response.Response(c, resp, err)
}
