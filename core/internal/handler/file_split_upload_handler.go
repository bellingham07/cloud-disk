package handler

import (
	"cloud-disk/core/helper"
	"net/http"

	"cloud-disk/core/internal/logic"
	"cloud-disk/core/internal/svc"
	"cloud-disk/core/internal/types"
	"github.com/zeromicro/go-zero/rest/httpx"
)

func FileSplitUploadHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.FileSplitUploadRequest
		if err := httpx.Parse(r, &req); err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		helper.FileSplitUpload(r)

		l := logic.NewFileSplitUploadLogic(r.Context(), svcCtx)
		resp, err := l.FileSplitUpload(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
