package utils

import (
	"math/rand"
	"time"
)

func GenerateApiKey(length int) string {
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789-"
	var seededRand *rand.Rand = rand.New(
		rand.NewSource(time.Now().UnixNano()))
	apiKey := make([]byte, length)
	for i := range apiKey {
		apiKey[i] = charset[seededRand.Intn(len(charset))]
	}
	return "sk-" + string(apiKey)
}

func GetSecondsFromDay(daysLater int) int {
	// 获取当前时间，加上指定天数
	futureTime := time.Now().AddDate(0, 0, daysLater)
	// 获取当天零点时间
	startOfDay := time.Date(futureTime.Year(), futureTime.Month(), futureTime.Day(), 0, 0, 0, 0, futureTime.Location())
	// 计算时间差并返回秒数
	return int(futureTime.Sub(startOfDay).Seconds())
}

func GetDay(daysLater int) string {
	// 获取当前时间，加上指定天数
	futureTime := time.Now().AddDate(0, 0, daysLater)
	// 获取当天零点时间
	startOfDay := time.Date(futureTime.Year(), futureTime.Month(), futureTime.Day(), 0, 0, 0, 0, futureTime.Location()).Format("2006-01-02 15:04:05")
	// 计算时间差并返回秒数
	return startOfDay
}
