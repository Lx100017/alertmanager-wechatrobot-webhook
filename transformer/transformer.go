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

	if status == "firing" {
		buffer.WriteString(fmt.Sprintf("### 监控告警触发通知\n"))
	} else if status == "resolved" {
		buffer.WriteString(fmt.Sprintf("### 监控告警修复通知\n"))
	}

	for _, alert := range notification.Alerts {
		labels := alert.Labels
		buffer.WriteString(fmt.Sprintf("\n>触发时间: %s\n", alert.StartsAt.Format("2006-01-02 15:04:05")))
		buffer.WriteString(fmt.Sprintf("\n>告警严重级别: %s\n", labels["severity"]))
		buffer.WriteString(fmt.Sprintf("\n>告警类型名称: %s\n", labels["alertname"]))
		buffer.WriteString(fmt.Sprintf("\n>告警实例地址: %s\n", labels["instance"]))

		annotations := alert.Annotations
		buffer.WriteString(fmt.Sprintf("\n>告警简述: %s\n", annotations["summary"]))
		buffer.WriteString(fmt.Sprintf("\n>告警详情: %s\n", annotations["description"]))
	}

	markdown = &model.WeChatMarkdown{
		MsgType: "markdown",
		Markdown: &model.Markdown{
			Content: buffer.String(),
		},
	}

	return
}

func TransformToAppMessageText(notification model.Notification, ToTag string, AgentId string) (markdown *model.WeChatText, err error) {

	status := notification.Status

	var buffer bytes.Buffer

	if status == "firing" {
		buffer.WriteString(fmt.Sprintf("监控告警触发通知\n"))
	} else if status == "resolved" {
		buffer.WriteString(fmt.Sprintf("监控告警修复通知\n"))
	}

	for _, alert := range notification.Alerts {
		labels := alert.Labels
		buffer.WriteString(fmt.Sprintf("\n触发时间: %s\n", alert.StartsAt.Format("2006-01-02 15:04:05")))
		buffer.WriteString(fmt.Sprintf("\n告警严重级别: %s\n", labels["severity"]))
		buffer.WriteString(fmt.Sprintf("\n告警类型名称: %s\n", labels["alertname"]))
		buffer.WriteString(fmt.Sprintf("\n告警实例地址: %s\n", labels["instance"]))

		annotations := alert.Annotations
		buffer.WriteString(fmt.Sprintf("\n告警简述: %s\n", annotations["summary"]))
		buffer.WriteString(fmt.Sprintf("\n告警详情: %s\n", annotations["description"]))
	}

	markdown = &model.WeChatText{
		ToTag:   ToTag,
		AgentId: AgentId,
		MsgType: "text",
		Text: &model.Text{
			Content: buffer.String(),
		},
	}

	return
}
