package models

//商品明细
type GoodsDetail struct {
	//商品的编号(必填)
	Goods_id string `json:"goods_id"`
	//支付宝定义的统一商品编号(不必填)
	Alipay_goods_id string `json:"alipay_goods_id"`
	//商品名称(必填)
	Goods_name string `json:"goods_name"`
	//商品数量(必填)
	Quantity int `json:"quantity"`
	//商品单价，单位为元(必填)
	Price string `json:"price"`
	//商品类目(不必填)
	Goods_category string `json:"goods_category"`
	//商品描述信息(不必填)
	Body string `json:"body"`
	//商品的展示地址(不必填)
	Show_url string `json:"show_url"`
}

//业务扩展参数
type ExtendParams struct {
	//系统商编号 该参数作为系统商返佣数据提取的依据，请填写系统商签约协议的PID(不必填)
	Sys_service_provider_id string `json:"sys_service_provider_id"`
	//使用花呗分期要进行的分期数(不必填)
	Hb_fq_num string `json:"hb_fq_num"`
	//使用花呗分期需要卖家承担的手续费比例的百分值，传入100代表100%(不必填)
	Hb_fq_seller_percent string `json:"hb_fq_seller_percent"`
}

//分账信息
type RoyaltyInfo struct {
	//分账类型 卖家的分账类型，目前只支持传入ROYALTY（普通分账类型）。(不必填)
	Royalty_info string `json:"royalty_info"`
	//分账明细的信息，可以描述多条分账指令，json数组。(不必填)
	Royalty_detail []*RoyaltyDetailInfos `json:"royalty_detail_infos"`
}

//分账明细信息
type RoyaltyDetailInfos struct {
	//分账序列号，表示分账执行的顺序，必须为正整数(不必填)
	Serial_no int `json:"serial_no"`
	//接受分账金额的账户类型： 	userId：支付宝账号对应的支付宝唯一用户号。  bankIndex：分账到银行账户的银行编号。目前暂时只支持分账到一个银行编号。 storeId：分账到门店对应的银行卡编号。 默认值为userId(不必填)
	Trans_in_type string `json:"trans_in_type"`
	//分账批次号 分账批次号。 目前需要和转入账号类型为bankIndex配合使用。(必填)
	Batch_no string `json:"batch_no"`
	//商户分账的外部关联号，用于关联到每一笔分账信息，商户需保证其唯一性。 如果为空，该值则默认为“商户网站唯一订单号+分账序列号”(不必填)
	Out_relation_id string `json:"out_relation_id"`
	//要分账的账户类型。 目前只支持userId：支付宝账号对应的支付宝唯一用户号。 默认值为userId。(必填)
	Trans_out_type string `json:"trans_out_type"`
	//如果转出账号类型为userId，本参数为要分账的支付宝账号对应的支付宝唯一用户号。以2088开头的纯16位数字。(必填)
	Trans_out string `json:"trans_out"`
	//如果转入账号类型为userId，本参数为接受分账金额的支付宝账号对应的支付宝唯一用户号。以2088开头的纯16位数字。 	如果转入账号类型为bankIndex，本参数为28位的银行编号（商户和支付宝签约时确定）。 如果转入账号类型为storeId，本参数为商户的门店ID。(必填)
	Trans_in string `json:"trans_in"`
	//分账的金额，单位为元(必填)
	Amount int `json:"amount"`
	//分账描述信息(不必填)
	Desc string `json:"desc"`
	//分账的比例，值为20代表按20%的比例分账(不必填)
	Amount_percentage string `json:"amount_percentage"`
}

//二级商户信息
type SubMerchant struct {
	//二级商户的支付宝id(必填)
	Merchant_id string `json:"merchant_id"`
}

//公共参数
type PublicRequest struct {
}

//详细参数
type BizContent struct {
	//***************************************************************************************
	//请求参数
	//商户订单号,64个字符以内、只能包含字母、数字、下划线；需保证在商户端不重复(必填)
	Out_trade_no string `json:"out_trade_no"`
	//卖家支付宝用户ID。 如果该值为空，则默认为商户签约账号对应的支付宝用户ID(不必填)
	Seller_id string `json:"seller_id"`
	//订单总金额，单位为元，精确到小数点后两位，取值范围[0.01,100000000] 如果同时传入了【打折金额】，【不可打折金额】，【订单总金额】三者，则必须满足如下条件：【订单总金额】=【打折金额】+【不可打折金额】(必填)
	Total_amount string `json:"total_amount"`
	//可打折金额. 参与优惠计算的金额，单位为元，精确到小数点后两位，取值范围[0.01,100000000] 如果该值未传入，但传入了【订单总金额】，【不可打折金额】则该值默认为【订单总金额】-【不可打折金额】(不必填)
	Discountable_amount string `json:"discountable_amount"`
	//不可打折金额. 不参与优惠计算的金额，单位为元，精确到小数点后两位，取值范围[0.01,100000000] 如果该值未传入，但传入了【订单总金额】,【打折金额】，则该值默认为【订单总金额】-【打折金额】(不必填)
	Undiscountable_amount string `json:"undiscountable_amount"`
	//买家支付宝账号(不必填)
	Buyer_logon_id string `json:"buyer_logon_id"`
	//订单标题(必填)
	Subject string `json:"subject"`
	//对交易或商品的描述(不必填)
	Body string `json:"body"`
	//订单包含的商品列表信息.Json格式. 其它说明详见：“商品明细说明”(不必填)
	Goods_detail []*GoodsDetail `json:"goods_detail"`
	//商户操作员编号(不必填)
	Operator_id string `json:"operator_id"`
	//商户门店编号(不必填)
	Store_id string `json:"store_id"`
	//商户机具终端编号(不必填)
	Terminal_id string `json:"terminal_id"`
	//业务扩展参数(不必填)
	Extend_params *ExtendParams `json:"extend_params"`
	//该笔订单允许的最晚付款时间，逾期将关闭交易。取值范围：1m～15d。m-分钟，h-小时，d-天，1c-当天（1c-当天的情况下，无论交易何时创建，都在0点关闭）。 该参数数值不接受小数点， 如 1.5h，可转换为 90m。(不必填)
	Timeout_express string `json:"timeout_express"`
	//描述分账信息，json格式。(不必填)
	Royalty_info *RoyaltyInfo `json:"royalty_info"`
	//二级商户信息,当前只对特殊银行机构特定场景下使用此字段(不必填)
	Sub_merchant *SubMerchant `json:"sub_merchant"`
	//支付宝店铺的门店ID(不必填)
	Alipay_store_id string `json:"alipay_store_id"`
}

//请求参数
type AllRequest struct {
	//***************************************************************************************
	//公共参数
	Alipay_sdk string `json:"alipay_sdk"`
	//支付宝分配给开发者的应用ID(必填)
	App_id string `json:"app_id"`
	//接口名称(必填)
	Method string `json:"method"`
	//仅支持json(不必填)
	Format string `json:"format"`
	//请求使用的编码格式，如utf-8，gbk，gb2312等(必填)
	Charset string `json:"charset"`
	//商户生成签名字符串所使用的的签名算法类型，目前支持RSA(必填)
	Sign_type string `json:"sign_type"`
	//商户请求参数的签名串(必填)
	Sign string `json:"sign"`
	//发送请求的时间,格式"yyyy-MM-dd HH:mm:ss" (必填)
	Timestamp string `json:"timestamp"`
	//调用的接口版本，固定为:1.0(必填)
	Version string `json:"version"`
	//支付宝服务器主动通知商户服务器里指定的页面http/https路径(不必填)
	Notify_url string `json:"notify_url"`
	//详见应用授权描述(不必填)
	App_auth_token string `json:"app_auth_token"`
	//请求参数的集合，最大长度不限，除公共参数外所有请求参数都必须放在这个参数中传递，具体参照各产品快速接入稳定(必填)
	Biz_content *BizContent `json:"biz_content"`
}
