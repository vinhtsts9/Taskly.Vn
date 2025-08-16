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
)

type vnpayService struct {
	TmnCode    string
	HashSecret string
	PaymentURL string
	ReturnURL  string
}

func NewVNPayService(code, secret, payURL, returnURL string) *vnpayService {
	return &vnpayService{
		TmnCode:    code,
		HashSecret: secret,
		PaymentURL: payURL,
		ReturnURL:  returnURL,
	}
}

func (v *vnpayService) GeneratePaymentURL(orderID string, amount int) (string, error) {
	t := time.Now()
	vnpParams := map[string]string{
		"vnp_Version":    "2.1.0",
		"vnp_Command":    "pay",
		"vnp_TmnCode":    v.TmnCode,
		"vnp_Amount":     strconv.Itoa(amount),
		"vnp_CurrCode":   "VND",
		"vnp_TxnRef":     orderID,
		"vnp_OrderInfo":  "Thanh toan don hang " + orderID,
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
		signData += k + "=" + val
		query += k + "=" + url.QueryEscape(val)
	}

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
		signData += k + "=" + params.Get(k)
	}

	h := hmac.New(sha512.New, []byte(v.HashSecret))
	h.Write([]byte(signData))
	hash := hex.EncodeToString(h.Sum(nil))

	return hash == receivedHash
}
