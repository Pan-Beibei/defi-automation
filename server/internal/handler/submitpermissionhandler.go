// Code scaffolded by goctl. Safe to edit.
// goctl 1.10.1

package handler

import (
	"fmt"
	"net/http"

	"server/internal/logic"
	"server/internal/svc"
	"server/internal/types"

	"github.com/zeromicro/go-zero/rest/httpx"
)

func SubmitPermissionHandler(svcCtx *svc.ServiceContext) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		var req types.SubmitPermissionReq
		if err := httpx.Parse(r, &req); err != nil {
		fmt.Printf(">>> parse error: %v\n", err)  // 解析错误日志

			httpx.ErrorCtx(r.Context(), w, err)
			return
		}

		l := logic.NewSubmitPermissionLogic(r.Context(), svcCtx)
		resp, err := l.SubmitPermission(&req)
		if err != nil {
			httpx.ErrorCtx(r.Context(), w, err)
		} else {
			httpx.OkJsonCtx(r.Context(), w, resp)
		}
	}
}
