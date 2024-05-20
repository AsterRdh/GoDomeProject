package FileControllers

import (
	"awesomeProject/model"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"strings"
)

func UploadPub(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.String(500, "上传图片出错")
	}
	s := uuid.New().String()
	split := strings.Split(file.Filename, ".")
	s = s + "." + split[len(split)-1]
	path := model.PublicFSPath + s
	err = c.SaveUploadedFile(file, path)
	if err != nil {
		c.String(500, "上传图片出错")
	}

	resMsg := model.ResMessage{
		OkFlag:     true,
		Message:    "/fs/" + s,
		Data:       nil,
		ErrDetails: err,
	}

	c.JSON(200, resMsg)
}
func UploadAuthed(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.String(500, "上传图片出错")
	}
	s := uuid.New().String()
	split := strings.Split(file.Filename, ".")
	s = s + "." + split[len(split)-1]
	path := model.AuthedFSPath + s

	err = c.SaveUploadedFile(file, path)
	if err != nil {
		c.String(500, "上传图片出错")
	}
	resMsg := model.ResMessage{
		OkFlag:     true,
		Message:    "/authed/fs/" + s,
		Data:       nil,
		ErrDetails: err,
	}

	c.JSON(200, resMsg)
}
