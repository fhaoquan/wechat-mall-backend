package cms

import (
	"encoding/json"
	"github.com/go-playground/validator"
	"github.com/gorilla/mux"
	"net/http"
	"strconv"
	"wechat-mall-backend/defs"
	"wechat-mall-backend/errs"
	"wechat-mall-backend/model"
)

// 查询-优惠券列表
func (h *Handler) GetCouponList(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	page, _ := strconv.Atoi(vars["page"])
	size, _ := strconv.Atoi(vars["size"])

	couponVOList := []defs.CMSCouponVO{}
	couponList, total := h.service.CouponService.GetCouponList(page, size, defs.ALL)
	for _, v := range *couponList {
		couponVO := defs.CMSCouponVO{}
		couponVO.Id = v.Id
		couponVO.Title = v.Title
		couponVO.FullMoney = v.FullMoney
		couponVO.Minus = v.Minus
		couponVO.Rate = v.Rate
		couponVO.Type = v.Type
		couponVO.GrantNum = v.GrantNum
		couponVO.LimitNum = v.LimitNum
		couponVO.StartTime = v.StartTime
		couponVO.EndTime = v.EndTime
		couponVO.Description = v.Description
		couponVO.Online = v.Online
		couponVOList = append(couponVOList, couponVO)
	}
	resp := make(map[string]interface{})
	resp["list"] = couponVOList
	resp["total"] = total
	defs.SendNormalResponse(w, resp)
}

// 查询优惠券
func (h *Handler) GetCoupon(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	coupon := h.service.CouponService.GetCouponById(id)
	if coupon.Id == defs.ZERO || coupon.Del == defs.DELETE {
		panic(errs.ErrorCoupon)
	}
	couponVO := defs.CMSCouponVO{}
	couponVO.Id = coupon.Id
	couponVO.Title = coupon.Title
	couponVO.FullMoney = coupon.FullMoney
	couponVO.Minus = coupon.Minus
	couponVO.Rate = coupon.Rate
	couponVO.Type = coupon.Type
	couponVO.GrantNum = coupon.GrantNum
	couponVO.LimitNum = coupon.LimitNum
	couponVO.StartTime = coupon.StartTime
	couponVO.EndTime = coupon.EndTime
	couponVO.Description = coupon.Description
	couponVO.Online = coupon.Online
	defs.SendNormalResponse(w, couponVO)
}

// 新增/删除 优惠券
func (h *Handler) DoEditCoupon(w http.ResponseWriter, r *http.Request) {
	req := defs.CMSCouponReq{}
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		panic(errs.ErrorParameterValidate)
	}
	validate := validator.New()
	if err := validate.Struct(req); err != nil {
		panic(errs.NewParameterError(err.Error()))
	}
	if req.Id == defs.ZERO {
		_, total := h.service.CouponService.GetCouponList(1, defs.CouponMax, defs.ALL)
		if total >= defs.CouponMax {
			panic(errs.NewErrorCoupon("最多只能添加" + strconv.Itoa(defs.CouponMax) + "张优惠券"))
		}
		coupon := model.WechatMallCouponDO{}
		coupon.Title = req.Title
		coupon.FullMoney = req.FullMoney
		coupon.Minus = req.Minus
		coupon.Rate = req.Rate
		coupon.Type = req.Type
		coupon.GrantNum = req.GrantNum
		coupon.LimitNum = req.LimitNum
		coupon.StartTime = req.StartTime
		coupon.EndTime = req.EndTime
		coupon.Description = req.Description
		coupon.Online = req.Online
		h.service.CouponService.AddCoupon(&coupon)
	} else {
		coupon := h.service.CouponService.GetCouponById(req.Id)
		if coupon.Id == defs.ZERO || coupon.Del == defs.DELETE {
			panic(errs.ErrorCoupon)
		}
		coupon.Title = req.Title
		coupon.FullMoney = req.FullMoney
		coupon.Minus = req.Minus
		coupon.Rate = req.Rate
		coupon.Type = req.Type
		coupon.GrantNum = req.GrantNum
		coupon.LimitNum = req.LimitNum
		coupon.StartTime = req.StartTime
		coupon.EndTime = req.EndTime
		coupon.Description = req.Description
		coupon.Online = req.Online
		h.service.CouponService.UpdateCouponById(coupon)
	}
	defs.SendNormalResponse(w, "ok")
}

// 删除优惠券
func (h *Handler) DoDeleteCoupon(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	id, _ := strconv.Atoi(vars["id"])
	coupon := h.service.CouponService.GetCouponById(id)
	if coupon.Id == defs.ZERO || coupon.Del == defs.DELETE {
		panic(errs.ErrorCoupon)
	}
	coupon.Del = defs.DELETE
	h.service.CouponService.UpdateCouponById(coupon)
	defs.SendNormalResponse(w, "ok")
}
