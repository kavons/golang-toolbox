package wx_test

import (
	"net/http"
	"strconv"
	"testing"
	"time"

	"github.com/iGoogle-ink/gopay"
	"reflect"
)

var (
	appId  = "wxa3c60f9442e768dd"
	mchId  = "1494888482"
	apiKey = "688DB4632F17FDAA76FE5C16E8339ECD"
	client *gopay.WeChatClient
)

func init() {
	client = gopay.NewWeChatClient(appId, mchId, apiKey, true)
	client.SetCountry(gopay.China)
}

func TestUnifiedOrderAndGetMiniPaySign(t *testing.T) {
	bm := make(gopay.BodyMap)
	bm.Set("sign_type", gopay.SignType_MD5)
	bm.Set("nonce_str", gopay.GetRandomString(32))
	bm.Set("body", "锁相册") //

	bm.Set("out_trade_no", "192936950645133312") //
	bm.Set("trade_type", gopay.TradeType_JsApi)
	bm.Set("openid", "omMIK4-vQ5Y0xjuFskR3F0wS7t64")
	bm.Set("spbill_create_ip", "127.0.0.1")                       //
	bm.Set("device_info", "bb7d429d-ca34-4dc1-90fe-7ec42dba7e18") //

	bm.Set("fee_type", "CNY")
	bm.Set("total_fee", "8800") //

	bm.Set("notify_url", "https://api.ibanana.club")

	wxRsp, err := client.UnifiedOrder(bm)
	if err != nil {
		t.Logf(err.Error())
	}

	if wxRsp.ResultCode == gopay.SUCCESS {
		pac := "prepay_id=" + wxRsp.PrepayId
		timeStamp := strconv.FormatInt(time.Now().Unix(), 10)
		paySign := gopay.GetMiniPaySign(appId, wxRsp.NonceStr, pac, gopay.SignType_MD5, timeStamp, apiKey)
		t.Logf("paySign:%s", paySign)

		// appId
		// timeStamp
		// wxRsp.NonceStr
		// pac
		// gopay.SignType_MD5
	}
}

func TestQueryOrder(t *testing.T) {
	body := make(gopay.BodyMap)
	body.Set("out_trade_no", "192936950645133312")
	body.Set("nonce_str", gopay.GetRandomString(32))
	body.Set("sign_type", gopay.SignType_MD5)

	wxRsp, err := client.QueryOrder(body)
	if err != nil {
		t.Logf(err.Error())
	}

	t.Log(wxRsp)
	// wxRsp.TradeState
}

type ServiceProcessor interface {
	Do(*gopay.WeChatNotifyRequest) error
}

func ParseWeChatNotifyResult(req *http.Request, processor ServiceProcessor) (resp *gopay.WeChatNotifyResponse) {
	resp = new(gopay.WeChatNotifyResponse)
	notifyReq, err := gopay.ParseWeChatNotifyResult(req)
	if err != nil {
		resp.ReturnCode = gopay.FAIL
		resp.ReturnMsg = "参数格式校验错误"
		return
	}

	ok, err := gopay.VerifyWeChatSign(apiKey, gopay.SignType_MD5, notifyReq)
	if err != nil || !ok {
		resp.ReturnCode = gopay.FAIL
		resp.ReturnMsg = "签名失败"
		return
	}

	if !reflect.ValueOf(processor).IsNil() {
		err = processor.Do(notifyReq)
		if err != nil {
			resp.ReturnCode = gopay.FAIL
			resp.ReturnMsg = "业务处理失败"
			return
		}
	}

	resp.ReturnCode = gopay.SUCCESS
	resp.ReturnMsg = "OK"
	return
}
