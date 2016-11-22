package main

import (
	"alipaydev/models"
	"alipaydev/utils"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"time"
)

var targetUrl string = "https://openapi.alipaydev.com/gateway.do"

//拼装http请求
func packetHttpRequest() (url string) {
	var allRequest *models.AllRequest = new(models.AllRequest)
	//这里填写你的app_id
	allRequest.App_id = "2016101000651727"
	allRequest.Format = "json"
	allRequest.Method = "alipay.trade.precreate"
	allRequest.Charset = "utf-8"
	//这里填写你的异步回调地址
	allRequest.Notify_url = "http://120.26.164.8:8080/gateway.do"
	allRequest.Sign_type = "RSA"
	allRequest.Timestamp = time.Unix(time.Now().Unix(), 0).Format("2006-01-02 15:04:05")
	allRequest.Version = "1.0"

	var bizContent *models.BizContent = new(models.BizContent)
	bizContent.Out_trade_no = "tradeprecreate14794641232588204158"
	bizContent.Total_amount = "0.01"
	bizContent.Subject = "xxx品牌xxx门店当面付扫码消费"
	bizContent.Store_id = "test_store_id"
	bizContent.Seller_id = "2088102178907643"

	var goodDetail []*models.GoodsDetail = make([]*models.GoodsDetail, 0, 1)
	var extendParams *models.ExtendParams = new(models.ExtendParams)
	var royaltyInfo *models.RoyaltyInfo = new(models.RoyaltyInfo)
	var royaltyDetail []*models.RoyaltyDetailInfos = make([]*models.RoyaltyDetailInfos, 0, 1)

	var subMerchant *models.SubMerchant = new(models.SubMerchant)

	bizContent.Goods_detail = goodDetail
	bizContent.Extend_params = extendParams
	bizContent.Royalty_info = royaltyInfo
	bizContent.Royalty_info.Royalty_detail = royaltyDetail
	bizContent.Sub_merchant = subMerchant

	allRequest.Biz_content = bizContent

	//生成签名
	sign := utils.AlipaySign(allRequest)

	allRequest.Sign = sign

	targetUrl = utils.HttpUrl(targetUrl, allRequest)
	//测试签名
	//	utils.TestSign()

	return targetUrl
}

//发送get请求 获取数据
func Get(url string) (content string, statusCode int) {
	resp, err := http.Get(url)
	if err != nil {
		statusCode = -100
	}
	defer resp.Body.Close()
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		statusCode = -200
		return
	}
	statusCode = resp.StatusCode
	content = string(data)
	fmt.Println("返回的statusCode:", statusCode)
	fmt.Println("返回的内容:", content)
	return
}

//处理get请求的方法
func dealResponse(content string, statusCode int) {

	if statusCode == 200 {
		fmt.Println("生成二维码成功")
	} else {
		fmt.Println("生成二维码失败")
	}

	var data map[string]interface{}

	json.Unmarshal([]byte(content), &data)
	for _, v := range data {
		//		fmt.Println(v)
		if value, ok := v.(map[string]interface{}); ok {
			fmt.Println(value["qr_code"])
		} else {
			//			fmt.Println("false")
		}
	}
}
func main() {
	content, statusCode := Get(packetHttpRequest())
	dealResponse(content, statusCode)
	// time.Sleep(5 * time.Second)
}
