package rpcserverimpl

// protos server方法的实现
import (
	"context"
	upProto "micro-message-system/uploadserver/protos"
)

// Upload : upload结构体
type Upload struct{}

// UploadEntry : 获取上传入口
func (u *Upload) UploadEntry(
	ctx context.Context,
	req *upProto.ReqEntry,
	res *upProto.RespEntry) error {

	res.Entry = ""
	return nil
}
