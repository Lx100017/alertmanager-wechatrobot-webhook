package main

import (
	"flag"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jacobslei/alertmanager-wechatrobot-webhook/model"
	"github.com/jacobslei/alertmanager-wechatrobot-webhook/notifier"
)

var (
	h         bool
	RobotKey  string
	CorpId    string
	AgentId   string
	AppSecret string
)

func init() {
	flag.BoolVar(&h, "h", false, "help")
	flag.StringVar(&RobotKey, "RobotKey", "", "global wechatrobot webhook, you can overwrite by alert rule with annotations wechatRobot")
	flag.StringVar(&CorpId, "CorpId", "", "global wechat corpId, you can overwrite by alert rule with annotations corpId")
	flag.StringVar(&AgentId, "AgentId", "", "global wechat AgentId, you can overwrite by alert rule with annotations AgentId")
	flag.StringVar(&AppSecret, "AppSecret", "", "global wechat AppSecret, you can overwrite by alert rule with annotations AppSecret")

}

func main() {

	flag.Parse()

	if h {
		flag.Usage()
		return
	}

	router := gin.Default()
	router.POST("/webhook", func(c *gin.Context) {
		var notification model.Notification

		err := c.BindJSON(&notification)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		RobotKey := c.DefaultQuery("key", RobotKey)

		err = notifier.Send(notification, RobotKey)

		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return

		}

		c.JSON(http.StatusOK, gin.H{"message": "send to wechatbot successful!"})

	})
	router.POST("/appWebhook", func(c *gin.Context) {
		var notification model.Notification
		err := c.BindJSON(&notification)

		if err != nil {

			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		ToTag := c.DefaultQuery("toTag", "1")
		CorpId := c.DefaultQuery("corpId", CorpId)
		AgentId := c.DefaultQuery("agentId", AgentId)
		AppSecret := c.DefaultQuery("appSecret", AppSecret)
		err = notifier.SendApp(notification, ToTag, CorpId, AgentId, AppSecret)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return

		}

		c.JSON(http.StatusOK, gin.H{"message": "send to wechatbot successful!"})

	})
	router.Run(":8999")
}
