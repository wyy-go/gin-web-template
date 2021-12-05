package util

import (
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"

	"github.com/wyy-go/go-web-template/internal/config"
	"gorm.io/gorm"
	"math/rand"
	"time"

	"github.com/cespare/xxhash/v2" // xxhash 64位高性能版本
	//"bitbucket.org/StephaneBunel/xxhash-go" // xxhash 32位版
	"strings"
)

const (
	idSecret string = "HpegqL4ZWuRoma7dNzi9jQshUMwPk532XDbA8GxKcFvJtfrEVYC1n6SyBT"
)

func IdEncode(id int64) (code string) {
	if id == 0 {
		return
	}

	base := int64(len(idSecret))

	for id > 0 {
		m := id % base
		id = (id - m) / base
		code = string(idSecret[m]) + code
	}
	return
}

func IdDecode(code string) (id int64) {
	if code == "" {
		return
	}
	if i := strings.LastIndex(code, "0"); i != -1 {
		code = code[i+1:]
	}
	code = reverseString(code)
	base := len(idSecret)
	for i := range code {
		id += int64(strings.Index(idSecret, string(code[i])) * pow(base, i))
	}
	return
}

func reverseBytes(b []byte) []byte {
	for from, to := 0, len(b)-1; from < to; from, to = from+1, to-1 {
		b[from], b[to] = b[to], b[from]
	}
	return b
}

func reverseString(s string) string {
	runes := []rune(s)
	for from, to := 0, len(runes)-1; from < to; from, to = from+1, to-1 {
		runes[from], runes[to] = runes[to], runes[from]
	}
	return string(runes)
}

func pow(x int, y int) int {
	var result = 1
	for i := 0; i < y; i++ {
		result *= x
	}
	return result
}

// 根据ID生成统一编码，生成的编码长度较短，不可逆
// 相近的数字编码之后无相关性
func GenCode(id int64) (code string) {
	buf := &bytes.Buffer{}
	binary.Write(buf, binary.BigEndian, id)

	//n := xxhash.Checksum32(buf.Bytes())
	//base := uint32(len(idSecret))
	n := xxhash.Sum64(buf.Bytes())
	base := uint64(len(idSecret))
	for n > 0 {
		m := n % base
		n = (n - m) / base
		code += string(idSecret[m])
	}
	return
}

func ShuffleString(s string) string {
	rand.Seed(time.Now().UnixNano())
	runes := []rune(s)
	for i := len(runes) - 1; i > 0; i-- {
		num := rand.Intn(i + 1)
		runes[i], runes[num] = runes[num], runes[i]
	}

	return string(runes)
}

func RandNumber(n int) string {
	numberStr := "0123456789"
	numberBytes := []byte(numberStr)
	var result []byte
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < n; i++ {
		result = append(result, numberBytes[r.Intn(len(numberBytes))])
	}
	return string(result)
}

func GetOssUrl(key string) string {
	if key == "" {
		return ""
	}

	return config.GetConfig().OssConfig.CdnDomain + "/" + key
	//return "http://" + config.GetConfig().OssConfig.CdnDomain + "/" + key
	//return "https://" + config.GetOss().BucketName + "." + config.GetOss().Endpoint + "/" + key
}

func TimeFormat(t time.Time) string {
	if t.IsZero() {
		return "0000-00-00 00:00:00"
	}
	return t.Format("2006-01-02 15:04:05")
}

func TimeToUnix(t time.Time) int64 {
	if t.IsZero() {
		return 0
	}

	return t.Unix()
}

func TimeToUnixMs(t time.Time) int64 {
	if t.IsZero() {
		return 0
	}

	return t.UnixNano() / 1e6
}

func RandAvatar() string {
	rand.Seed(time.Now().UnixNano())
	n := rand.Intn(10) + 1
	return GetOssUrl(fmt.Sprintf("avatar/default%02d.png", n))
}

func IsNotFound(err error) bool {
	return errors.Is(err, gorm.ErrRecordNotFound)
}

func HumanTime(secs int64) string {
	var num int64
	var quantifier string
	suffix := "后"
	diff := secs - time.Now().Unix()
	if diff < 0 {
		suffix = "前"
		diff = -diff
	}

	seconds := diff
	minutes := seconds / 60
	hours := minutes / 60
	days := hours / 24
	months := days / 30
	years := months / 12

	switch true {
	case years > 0:
		num = years
		quantifier = "年"
	case months > 0:
		num = months
		quantifier = "月"
	case days > 0:
		num = days
		quantifier = "日"
	case hours > 0:
		num = hours
		quantifier = "时"
	case minutes > 0:
		num = minutes
		quantifier = "分"
		break
	default:
		num = seconds
		quantifier = "秒"
	}

	return fmt.Sprintf("%d%s%s", num, quantifier, suffix)
}
