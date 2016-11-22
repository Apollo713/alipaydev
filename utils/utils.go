package utils

import (
	"bytes"
	"container/list"
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	//	"os"

	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"sort"
	"strings"
)

type ParamBytes struct {
	paramBytes []byte
}

type EveryBytes []*ParamBytes

//按照支付宝规则生成签名
func AlipaySign(param interface{}) string {
	//1.筛选 获取所有请求参数，不包括字节类型参数，如文件、字节流，剔除sign、sign_type字段。
	paramBytes, err := json.Marshal(param)
	if err != nil {
		return ""
	}
	//	fmt.Println("解析成的字节数组为:", string(paramBytes))
	//	fmt.Println(paramBytes[2])
	//	fmt.Println("长度是", len(paramBytes))
	//将字节数组变为字符串数组  需要匹配括号  分割逗号
	var multiBytes [][]byte = cutByComma(paramBytes)
	//	for _, v := range multiBytes {
	//		fmt.Println(string(v))
	//	}

	//2.排序 按照第一个字符的键值ASCII码递增排序（字母升序排序）如果遇到相同字符则按照第二个字符的键值ASCII码递增排序，以此类推。
	multiBytes = sortFunction(multiBytes)
	//	fmt.Println("排序后")
	//	for _, v := range multiBytes {
	//		fmt.Println(string(v))
	//	}
	//3.拼接 将排序后的参数与其对应值，组合成“参数=参数值”的格式，并且把这些参数用&字符连接起来，此时生成的字符串为待签名字符串。
	var signString string = stitchingBytes(multiBytes, 1)
	fmt.Println("拼接后的字符串:", signString)
	//4. 调用签名函数
	var sign string
	sign = signFunction(signString)
	fmt.Println("签名为:", sign)
	return sign
}

//分割逗号
func cutByComma(oneBytes []byte) [][]byte {
	var multiBytes [][]byte = make([][]byte, 0, 1)
	var frontIndex int = 1
	var lastIndex int = 1
	var l *list.List = list.New()
	//循环遍历的时候去掉最外面的{}
	for i := 1; i < len(oneBytes)-1; i++ {
		if string(oneBytes[i]) == "{" || string(oneBytes[i]) == "[" {
			l.PushBack(string(oneBytes[i]))
		}
		if string(oneBytes[i]) == "}" || string(oneBytes[i]) == "]" {
			l.Remove(l.Back())
		}
		if judgeByBracket(l) {
			if string(oneBytes[i]) == "," {
				//记录当前逗号坐标
				lastIndex = i
				var tempBytes []byte = make([]byte, 10, 10)
				tempBytes = oneBytes[frontIndex:lastIndex]
				//对[]byte处理完后 再添加到[][]byte中
				tempBytes = dealBytes(tempBytes)

				multiBytes = append(multiBytes, tempBytes)
				frontIndex = lastIndex + 1
			}
		} else {
			continue
		}
	}
	var lastBytes []byte = dealBytes(oneBytes[frontIndex : len(oneBytes)-1])
	//添加最后一个元素
	multiBytes = append(multiBytes, lastBytes)
	return multiBytes
}

//判断{} [] 等
func judgeByBracket(l *list.List) bool {
	if l.Len() != 0 {
		return false
	}
	return true
}

//处理[]byte格式
func dealBytes(tempBytes []byte) []byte {
	//只处理一次  冒号
	var colonFlag bool = true
	//记录 冒号坐标
	var colonIndex int = 0
	for i, v := range tempBytes {
		if string(v) == ":" {
			if colonFlag {
				//第一个冒号
				tempBytes[i] = '='
				colonFlag = false
				//第一个冒号坐标s
				colonIndex = i
			}
		}
		if string(v) == "\"" {
			//处理掉不需要的双引号
			if string(tempBytes[colonIndex+1]) == "{" && i >= colonIndex {
				break
			} else {
				tempBytes[i] = '@'
			}
		}
	}
	tempBytes = []byte(strings.Replace(string(tempBytes), "@", "", -1))
	return tempBytes
}

//排序 需要实现sort.Interface接口
//获取长度
func (every EveryBytes) Len() int {
	return len(every)
}

//排序规则 按照ASCII码的顺序排序
func (every EveryBytes) Less(i, j int) bool {
	var min int = 0
	//	fmt.Println("i index:", i)
	//	fmt.Println("I:", len(every[i].paramBytes))
	//	fmt.Println("i:", string(every[i].paramBytes))
	//	fmt.Println("j index:", j)
	//	fmt.Println("J:", len(every[j].paramBytes))
	//	fmt.Println("j:", string(every[j].paramBytes))
	//	fmt.Println("***")
	if len(every[i].paramBytes) > len(every[j].paramBytes) {
		min = len(every[j].paramBytes)
	} else {
		min = len(every[i].paramBytes)
	}
	for index := 0; index < min; index++ {
		//		fmt.Println(every[i].paramBytes[index])
		//		fmt.Println(every[j].paramBytes[index])
		if every[i].paramBytes[index] < every[j].paramBytes[index] {
			return true
		} else if every[i].paramBytes[index] > every[j].paramBytes[index] {
			return false
		} else {
			continue
		}
	}
	return false
}

//交换
func (every EveryBytes) Swap(i, j int) {
	var temp *ParamBytes = every[i]
	every[i] = every[j]
	every[j] = temp
}

//排序功能
func sortFunction(multiBytes [][]byte) [][]byte {
	everyBytes := make(EveryBytes, 0, 1)
	for i := 0; i < len(multiBytes); i++ {
		param := new(ParamBytes)
		param.paramBytes = multiBytes[i]
		everyBytes = append(everyBytes, param)
	}
	sort.Sort(everyBytes)
	//排序好之后再还回去
	for i, _ := range multiBytes {
		multiBytes[i] = everyBytes[i].paramBytes
	}
	return multiBytes
}

//拼接
func stitchingBytes(multiBytes [][]byte, operate int) string {
	var signBuf bytes.Buffer
	//循环
	for i, v := range multiBytes {
		//拼接
		//分为拼接url(2) 和 拼接待签名字符串(1)
		if operate == 1 && !strings.Contains(string(v), "sign_type") {
			//如果为 sign 剔除
			if strings.Contains(string(v), "sign") && !strings.Contains(string(v), "sign_type") {
				continue
			}
		}

		//如果为空值
		if v[len(v)-1] == '=' && strings.Count(string(v), "=") == 1 {
			//跳过空值
			continue
		} else if i == len(multiBytes)-1 {
			//如果是最后一串字符
			signBuf.WriteString(string(v))
		} else {
			if strings.Contains(string(v), "biz_content") {
				var needString string = string(dealEmptyJson(v))
				//				fmt.Println("处理后的json:", needString)
				signBuf.WriteString(needString)
				signBuf.WriteString("&")
			} else if strings.Contains(string(v), "timestamp") && operate == 2 {
				signBuf.WriteString(strings.Replace(string(v), " ", "+", -1))
				signBuf.WriteString("&")
			} else {
				signBuf.WriteString(string(v))
				signBuf.WriteString("&")
			}
			/*
				//用map清除空值有问题

					//处理biz_content 处理里面的空值
					if strings.Contains(string(v), "biz_content") {
						var bracketIndex int
						for index, value := range v {
							if string(value) == "=" {
								bracketIndex = index + 1
								break
							}
						}
						var jsonStr string = string(v[bracketIndex:])
						//json str 转map
						var dat map[string]interface{}
						if err := json.Unmarshal([]byte(jsonStr), &dat); err == nil {
							//					fmt.Println("==============json str 转map=======================")
							//					fmt.Println(dat)
						}
						//				fmt.Println("=================清除功能=====================")
						for k, v := range dat {

							value, _ := v.(string)
							//					fmt.Println("key:", k, "          ", "value:", value)
							if value == "" || value == "[]" || value == "null" {
								delete(dat, k)
							}
						}
						//map 到json str
						//				fmt.Println("================map 到json str=====================")
						jsonBytes, _ := json.Marshal(dat)
						//				fmt.Println(string(jsonBytes))

						signBuf.WriteString("biz_content=")
						signBuf.WriteString(string(jsonBytes))
						signBuf.WriteString("&")
					} else {
						signBuf.WriteString(string(v))
						signBuf.WriteString("&")
					}
			*/

		}
	}
	return signBuf.String()
}

//处理json 中空字符串
func dealEmptyJson(oneBytes []byte) []byte {
	//	fmt.Println(string(oneBytes))
	var multiBytes [][]byte = make([][]byte, 0, 1)
	var frontIndex int = 1
	var lastIndex int = 1
	var headBytes []byte
	var resultBytes []byte = make([]byte, 0, 1)

	var l *list.List = list.New()

	for i, v := range oneBytes {
		if string(v) == "{" {
			frontIndex = i + 1
			headBytes = oneBytes[:i]
			break
		}
	}
	//	fmt.Println(string(headBytes))
	//	fmt.Println(frontIndex)
	//循环遍历的时候去掉最外面的{}
	for i := frontIndex; i < len(oneBytes)-1; i++ {

		if string(oneBytes[i]) == "{" || string(oneBytes[i]) == "[" {
			l.PushBack(string(oneBytes[i]))
		}
		if string(oneBytes[i]) == "}" || string(oneBytes[i]) == "]" {
			l.Remove(l.Back())
		}
		if judgeByBracket(l) {
			if string(oneBytes[i]) == "," {
				//记录当前逗号坐标
				lastIndex = i
				var tempBytes []byte = make([]byte, 10, 10)
				tempBytes = oneBytes[frontIndex:lastIndex]
				if strings.Contains(string(tempBytes), "{") {
					tempBytes = dealEmptyJson(tempBytes)
				}
				//处理空值
				if strings.Contains(string(tempBytes), "\"\"") || strings.Contains(string(tempBytes), "[]") || strings.Contains(string(tempBytes), "{}") {
					frontIndex = lastIndex + 1
				} else {
					multiBytes = append(multiBytes, tempBytes)
					frontIndex = lastIndex + 1
				}
			}
		} else {
			continue
		}
	}
	resultBytes = headBytes
	resultBytes = append(resultBytes, '{')
	for i, v := range multiBytes {
		for _, value := range v {
			resultBytes = append(resultBytes, value)

		}
		if i != len(multiBytes)-1 {
			resultBytes = append(resultBytes, ',')
		}

	}
	resultBytes = append(resultBytes, '}')

	return resultBytes
}

//签名
func signFunction(signString string) string {
	block, _ := pem.Decode([]byte(`-----BEGIN PRIVATE KEY-----
MIICeAIBADANBgkqhkiG9w0BAQEFAASCAmIwggJeAgEAAoGBAL3zekgIQNpR0NQz
PQybP8ZAkIBRaIdYCxrSfeJsJG39xrzSjpamB3Xv7Z50RwkrCBooJKqxpn4awl8b
td+wAGQwLnojzvrizzwsthamKBo5q498Hc6BX3ZX93KuCFmwvteULeyPOtWssh5k
sLirOiTAYCT6lSx1acoibLcUwQmDAgMBAAECgYAM8+5xvQZXHN8lqTzPgEKwDTUN
Wv/Kwuk28gWdjAxL59NGiwEoKrg1hZ/pfzpc2K9bwUMG1Mhqrv50J9qWH1VXYZEQ
h14WyTZhQcOw+i2N8JbMrFmrWXHNlK/mN+TukWC+BUs5BvVRLq70CfUlGsktWvee
7piNkVP2mtv2ON/t2QJBAOD3+QB9Of16iCTdlSwyeMM5x7lDlOOQh+OTxUDF6ZAR
jM55fvJrqp0nZpgQTJEaf/ld11GVj/Om0wRDg8pPq9cCQQDYJvmjQVyTO+4M8mmn
aYIteu4/QKpJXRuQ45jRCe4Tio4xYkPmV3AHr66GRZ/VQLArJFXlhPU+8EfIfVhc
0no1AkEA3z0SiSq6xc62lKaRJYd8EHYgu7XVZDACuJDlV05NY9oWeLlVgKfYaRQ1
GUZrRD4gqco2JU4dx7FOileYysRehwJBALO6LKaHaY9vLHANfLZcL4bbiZCEl1Mr
HQmrhVyDUjdjZPpBB85Wc+ugM5CoAc+S2yj0LIwMstMjfbyCJOABjuUCQQCHDQGG
h6CPwQtXAkyRCsa3b9BD68y3D0lvbMUBdEnbycDz4tscNJ9OFIGnY88vgjnc//m5
B7q93Kvp2/Bw9O3I
-----END PRIVATE KEY-----
`))
	if block == nil {
		fmt.Println("解密私钥有问题")
		return ""
	}
	privateKey, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		fmt.Println("x509出错")
		return ""
	}
	//	fmt.Println("转换成的PCKS8格式:", privateKey)

	//	result, err := rsaSign(signString, privateKey)

	priKey, ok := (privateKey).(*rsa.PrivateKey)
	var result string = ""
	if ok {
		result, err := rsaSign(signString, priKey)
		if err != nil {
			fmt.Println("签名出错:", err)
		}
		result = strings.Replace(result, "=", "%3D", -1)
		result = strings.Replace(result, "+", "%2B", -1)
		result = strings.Replace(result, "/", "%2F", -1)
		//	result = strings.Replace(result, "", "", -1)
		return result
	} else {
		fmt.Println("私钥出错")
		return ""
	}

	return result
}

//RSA签名
func rsaSign(originData string, privateKey *rsa.PrivateKey) (string, error) {
	h := sha1.New()
	h.Write([]byte(originData))
	digest := h.Sum(nil)

	s, err := rsa.SignPKCS1v15(rand.Reader, privateKey, crypto.SHA1, digest)
	if err != nil {
		fmt.Errorf("rsaSign SignPKCS1v15 error")
		return "", err
	}
	data := base64.StdEncoding.EncodeToString(s)
	return string(data), nil
}

//获取请求的url
func HttpUrl(targetUrl string, param interface{}) string {
	paramBytes, err := json.Marshal(param)
	if err != nil {
		return ""
	}
	//处理
	var multiBytes [][]byte = cutByComma(paramBytes)
	//排序
	multiBytes = sortFunction(multiBytes)
	//参数
	var paramString string = stitchingBytes(multiBytes, 2)

	var bufUrl bytes.Buffer
	bufUrl.WriteString(targetUrl)
	bufUrl.WriteString("?")
	bufUrl.WriteString(paramString)

	fmt.Println("URL:", bufUrl.String())

	return bufUrl.String()
}

//将对a=123进行签名
func TestSign() {
	var signString string = "a=123"
	//	var signString string = `alipay_sdk=alipay-sdk-java-dynamicVersionNo&app_id=2016101000651727&biz_content={"out_trade_no":"tradeprecreate14795227404829949221","total_amount":"0.01","subject":"xxx品牌xxx门店当面付扫码消费","store_id":"test_store_id"}&charset=utf-8&format=json&method=alipay.trade.precreate&notify_url=http://120.26.164.8:8080/gateway.do&sign_type=RSA&timestamp=2016-11-19 10:32:20&version=1.0`
	fmt.Println("签名前的字符串:", signString)
	fmt.Println("签名之后的:", signFunction(signString))
}
