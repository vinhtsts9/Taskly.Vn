package impl

import (
	"crypto/hmac"
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"net/url"
	"sort"
	"strconv"
	"time"

	"Taskly.com/m/global"
)

type vnpayService struct {
	TmnCode    string
	HashSecret string
	PaymentURL string
	ReturnURL  string
}

func NewVNPayService(TmnCode, HashSecret, PaymentURL, ReturnURL string) *vnpayService {
	return &vnpayService{
		TmnCode:    TmnCode,
		HashSecret: HashSecret,
		PaymentURL: PaymentURL,
		ReturnURL:  ReturnURL,
	}
}

func (v *vnpayService) GeneratePaymentURL(orderID string, amount int, ipAddr string) (string, error) {
	t := time.Now()
	vnpParams := map[string]string{
		"vnp_Version":    "2.1.0",
		"vnp_Command":    "pay",
		"vnp_TmnCode":    v.TmnCode,
		"vnp_Amount":     strconv.Itoa(amount * 100), // VNPay uses cents
		"vnp_CurrCode":   "VND",
		"vnp_TxnRef":     orderID,
		"vnp_OrderInfo":  "Thanh toan don hang " + orderID,
		"vnp_OrderType":  "other",
		"vnp_IpAddr":     ipAddr,
		"vnp_Locale":     "vn",
		"vnp_ReturnUrl":  v.ReturnURL,
		"vnp_CreateDate": t.Format("20060102150405"),
	}

	var keys []string
	for k := range vnpParams {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var signData string
	var query string
	for _, k := range keys {
    val := vnpParams[k]
    if signData != "" {
        signData += "&"
        query += "&"
    }
    // ĐÚNG: encode cả signData và query
    encodedVal := url.QueryEscape(val)
    signData += k + "=" + encodedVal
    query += k + "=" + encodedVal
}
	global.Logger.Sugar().Infof("Signdata: %s", signData)
	h := hmac.New(sha512.New, []byte(v.HashSecret))
	h.Write([]byte(signData))
	signature := hex.EncodeToString(h.Sum(nil))

	paymentURL := fmt.Sprintf("%s?%s&vnp_SecureHash=%s", v.PaymentURL, query, signature)
	return paymentURL, nil
}

func (v *vnpayService) VerifySignature(params url.Values) bool {
	receivedHash := params.Get("vnp_SecureHash")
	params.Del("vnp_SecureHash")
	params.Del("vnp_SecureHashType")

	var keys []string
	for k := range params {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	var signData string
	for _, k := range keys {
    if signData != "" {
        signData += "&"
    }
    signData += k + "=" + url.QueryEscape(params.Get(k)) // <-- ĐÚNG!
}

	h := hmac.New(sha512.New, []byte(v.HashSecret))
	h.Write([]byte(signData))
	hash := hex.EncodeToString(h.Sum(nil))

	return hash == receivedHash
}
