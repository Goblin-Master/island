package api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"io"
	"tgwp/global"
	"tgwp/log/zlog"
	"tgwp/logic"
	"tgwp/model"
	"tgwp/repo/list"
	"tgwp/response"
	"tgwp/types"
	"tgwp/utils/fileUtils"
	"tgwp/utils/hashUtils"
)

// 上传的图片保存在uploads/images
// 由于这个api难于拆分，所以直接写在api包里面
func UploadImages(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)
	fileHeader, err := c.FormFile("file")
	if err != nil {
		zlog.CtxErrorf(ctx, "上传图片失败: %v", err)
		response.NewResponse(c).Error(response.PARAM_NOT_VALID)
		return
	}
	//图片大小限制
	if fileHeader.Size > global.IMAGE_SIZE {
		zlog.CtxErrorf(ctx, "上传图片失败: 图片大小超过限制")
		response.NewResponse(c).Error(response.IMAGE_OVER_SIZE)
		return
	}
	//图片格式限制
	filename := fileHeader.Filename
	suffix, err := fileUtils.ImageSuffixJudge(filename)
	if err != nil {
		zlog.CtxErrorf(ctx, "上传图片失败: %v", err)
		response.NewResponse(c).Error(response.IMAGE_NOT_SUPPORT)
		return
	}
	//文件hash（不同文件名，内容相同视为同一个文件）
	file, err := fileHeader.Open()
	defer file.Close()
	if err != nil {
		zlog.CtxErrorf(ctx, "上传图片失败: %v", err)
		response.NewResponse(c).Error(response.IMAGE_NOT_OPEN)
		return
	}
	byteData, _ := io.ReadAll(file)
	hash := hashUtils.Md5(byteData)
	// 判断hash是否已经存在
	var m model.Image
	err = global.DB.Take(&m, "hash = ?", hash).Error
	if err == nil {
		zlog.CtxInfof(ctx, "上传图片重复 %s <==> %s %s", filename, m.Filename, m.Hash)
		response.NewResponse(c).Success(m.WebPath())
		return
	}
	//入库
	filePath := fmt.Sprintf("%s/%s.%s", global.IMAGE_PATH, hash, suffix)
	m = model.Image{
		Filename: filename,
		Hash:     hash,
		Path:     filePath,
		Size:     fileHeader.Size,
	}
	err = global.DB.Create(&m).Error
	if err != nil {
		zlog.CtxErrorf(ctx, "上传图片失败: %v", err)
		response.NewResponse(c).Error(response.DATABASE_ERROR)
		return
	}
	err = c.SaveUploadedFile(fileHeader, filePath)
	if err != nil {
		zlog.CtxErrorf(ctx, "上传图片失败: %v", err)
		response.NewResponse(c).Error(response.IMAGE_UPLOAD_ERROR)
		return
	}
	response.NewResponse(c).Success(m.WebPath())
}
func DeleteImages(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)
	var req list.RemoveReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		response.NewResponse(c).Error(response.PARAM_NOT_VALID)
		return
	}
	zlog.CtxInfof(ctx, "DeleteImages request: %v", req)
	resp, err := logic.NewImagesLogic().DeleteImages(ctx, req)
	response.Response(c, resp, err)
}

func GetImages(c *gin.Context) {
	ctx := zlog.GetCtxFromGin(c)
	req, err := types.BindReq[list.PageInfo](c)
	if err != nil {
		return
	}
	zlog.CtxInfof(ctx, "GetImages request: %v", req)
	resp, err := logic.NewImagesLogic().GetImages(ctx, req)
	response.Response(c, resp, err)
}
