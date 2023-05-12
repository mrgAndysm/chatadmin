package routes

import (
	"chatgpt-go/db"
	"chatgpt-go/utils"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/go-redis/redis"
	"net/http"
	"time"
)

func AddKey(c *gin.Context) {
	daysLater := 2
	second := utils.GetSecondsFromDay(daysLater)
	gptkey := utils.GenerateApiKey(32)
	info := map[string]interface{}{}
	info["user_name"] = "aaa"
	info["phone"] = "123456789"
	info["exp_date"] = utils.GetDay(daysLater)
	v, _ := json.Marshal(info)
	err := db.RedisDb.Set(gptkey, v, time.Second*time.Duration(second)).Err()
	if err != nil {
		c.JSON(http.StatusInternalServerError, c.Error(err))
	}

	hkey := "gptkey:list"
	info["gptkey"] = gptkey
	vv, _ := json.Marshal(info)
	err = db.RedisDb.ZAdd(hkey, redis.Z{float64(second), vv}).Err()
	if err != nil {
		c.JSON(http.StatusInternalServerError, c.Error(err))
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "Success",
		"message": "Generate successfully",
		"data":    gptkey,
	})
}

func KeyList(c *gin.Context) {

	hkey := "gptkey:list"
	vals, err := db.RedisDb.ZRange(hkey, 0, -1).Result()
	if err != nil {
		panic(err)
	}

	v := map[string]interface{}{}
	re := []map[string]interface{}{}
	for _, val := range vals {
		var tmp = map[string]interface{}{}
		json.Unmarshal([]byte(val), &v)
		gptkey := v["gptkey"].(string)
		tmp["gptkey"] = val[:len(gptkey)-6] + "...."
		tmp["user_name"] = v["user_name"]
		tmp["exp_date"] = v["exp_date"]
		tmp["phone"] = v["phone"]
		re = append(re, tmp)
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "Success",
		"message": "GetList successfully",
		"data":    re,
	})
}
