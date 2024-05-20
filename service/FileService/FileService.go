package FileService

import (
	"awesomeProject/model"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"mime/multipart"
	"strings"
)

func UploadPubFile(c *gin.Context, file *multipart.FileHeader) (string, error) {
	fileName, err := uploadPubFile(c, file, model.PublicFSPath)
	if err != nil {
		return "", err
	}
	return "/fs/" + fileName, nil
}

func UploadAuthedFile(c *gin.Context, file *multipart.FileHeader) (string, error) {
	fileName, err := uploadPubFile(c, file, model.AuthedFSPath)
	if err != nil {
		return "", err
	}
	return "/authed/fs/" + fileName, nil
}

func uploadPubFile(c *gin.Context, file *multipart.FileHeader, basePath string) (string, error) {
	s := uuid.New().String()
	split := strings.Split(file.Filename, ".")
	s = s + "." + split[len(split)-1]
	path := basePath + s
	err := c.SaveUploadedFile(file, path)
	if err != nil {
		return "", err
	}
	return s, nil

}
