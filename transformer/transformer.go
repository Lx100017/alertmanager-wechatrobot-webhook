package transformer

import (
	"bytes"
	"fmt"

	"github.com/jacobslei/alertmanager-wechatrobot-webhook/model"
)

// TransformToMarkdown transform alertmanager notification to wechat markdow message
func TransformToMarkdown(notification model.Notification) (markdown *model.WeChatMarkdown, robotURL string, err error) {

	status := notification.Status

	annotations := notification.CommonAnnotations
	robotURL = annotations["wechatRobot"]

	var buffer bytes.Buffer

	buffer.WriteString(fmt.Sprintf("### 当前告警状态:%s \n", status))
	// buffer.WriteString(fmt.Sprintf("#### 告警项:\n"))

	for _, alert := range notification.Alerts {
		labels := alert.Labels
		buffer.WriteString(fmt.Sprintf("\n>告警严重级别: %s\n", labels["severity"]))
		buffer.WriteString(fmt.Sprintf("\n>告警类型名称: %s\n", labels["alertname"]))
		buffer.WriteString(fmt.Sprintf("\n>故障主机地址: %s\n", labels["ip"]))
		buffer.WriteString(fmt.Sprintf("\n>故障所属团队: %s\n", labels["team"]))

		annotations := alert.Annotations
		buffer.WriteString(fmt.Sprintf("\n>告警主题: %s\n", annotations["summary"]))
		buffer.WriteString(fmt.Sprintf("\n>告警详情: %s\n", annotations["description"]))
		buffer.WriteString(fmt.Sprintf("\n>触发时间: %s\n", alert.StartsAt.Format("2006-01-02 15:04:05")))
	}

	markdown = &model.WeChatMarkdown{
		MsgType: "markdown",
		Markdown: &model.Markdown{
			Content: buffer.String(),
		},
	}

	return
}
