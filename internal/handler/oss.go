package handler

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/qiniu/go-sdk/v7/auth/qbox"
	"github.com/qiniu/go-sdk/v7/storage"
	"github.com/spf13/viper"
	"log"
	"micro-shop/internal/util"
	"net/http"
	"path"
	"strings"
	"time"
)

type Oss struct {
	formUploader *storage.FormUploader
	mac          *qbox.Mac
	putExtra     *storage.PutExtra

	bucket        string
	putPolicy     *storage.PutPolicy
	bucketManager *storage.BucketManager
}

type UpHl struct {
	ctx    context.Context
	oss    Oss
	logger *log.Logger
}

func NewUpHl(ak, sk, bucket string, zone *storage.Zone, logger *log.Logger) *UpHl {
	putpolicy := storage.PutPolicy{
		Scope: bucket,
	}
	mac := qbox.NewMac(ak, sk)
	cfg := storage.Config{
		Zone:          zone,
		UseHTTPS:      false,
		UseCdnDomains: false,
	}
	formUploader := storage.NewFormUploader(&cfg)
	putExtra := storage.PutExtra{
		Params: map[string]string{
			"x:name": "micro-shop",
		},
	}
	bucketManager := storage.NewBucketManager(mac, &cfg)
	return &UpHl{
		oss: Oss{
			formUploader:  formUploader,
			bucketManager: bucketManager,
			putExtra:      &putExtra,
			bucket:        bucket,
			putPolicy:     &putpolicy,
			mac:           mac,
		},
		ctx:    context.Background(),
		logger: logger,
	}
}

func (uh *UpHl) UploadFile(c *gin.Context) {
	multi, err := c.MultipartForm()
	if err != nil {
		uh.logger.Printf("get multipart form error:%s", err.Error())
		c.JSON(http.StatusOK, util.GetResponse(501, "internal server error", ""))
		return
	}
	files := multi.File["files"]
	if len(files) == 0 {
		c.JSON(http.StatusOK, util.GetResponse(200, "ok", "{}"))
		return
	}
	ret := &storage.PutRet{}
	nameList := make([]string, len(files))
	base := viper.GetString("oss.base")
	for i, fh := range files {
		f, err := fh.Open()
		if err != nil {
			uh.logger.Printf("open file error:%s", err.Error())
			continue
		}
		name := EncodeName(fh.Filename)
		err = uh.oss.formUploader.Put(uh.ctx, ret, uh.oss.putPolicy.UploadToken(uh.oss.mac), name, f, fh.Size, uh.oss.putExtra)
		if err != nil {
			uh.logger.Printf("oss put error:%s", err.Error())
			continue
		}
		nameList[i] = fmt.Sprintf("http://%s/%s", base, name)
	}
	c.JSON(http.StatusOK, util.GetResponse(200, "ok", nameList))
}

func EncodeName(name string) string {
	ext := path.Ext(name)
	base := strings.TrimRight(name, ext)
	now := time.Now().Format("2006-01-02:15-04")
	name = fmt.Sprintf("%s-%s%s", base64.StdEncoding.EncodeToString([]byte(base)), now, ext)
	return name

}
