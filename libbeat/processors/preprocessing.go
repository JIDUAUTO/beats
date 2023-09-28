package processors

import (
	"github.com/bytedance/sonic"

	"github.com/elastic/beats/v7/libbeat/beat"
)

type LogFormat string

const (
	LogFormatIlogtail LogFormat = "ilogtail"
)

type IlogtailMessage struct {
	Contents struct {
		Content string `json:"content"`
	} `json:"contents"`
	Tags struct {
		ContainerImageName string `json:"container.image.name"`
		ContainerIp        string `json:"container.ip"`
		ContainerName      string `json:"container.name"`
		HostIp             string `json:"host.ip"`
		HostName           string `json:"host.name"`
		K8sNamespaceName   string `json:"k8s.namespace.name"`
		K8sNodeIp          string `json:"k8s.node.ip"`
		K8sNodeName        string `json:"k8s.node.name"`
		K8sPodName         string `json:"k8s.pod.name"`
		K8sPodUid          string `json:"k8s.pod.uid"`
		LogFilePath        string `json:"log.file.path"`
	} `json:"tags"`
	Time int `json:"time"`
}

func LogPreprocessing(event *beat.Event, format LogFormat) error {
	message := event.Fields["message"].(string)

	var err error
	switch format {
	case LogFormatIlogtail:
		var msg IlogtailMessage
		err = sonic.UnmarshalString(message, &msg)
		if err != nil {
			return err
		}

		event.Fields["namespace"] = msg.Tags.K8sNamespaceName
		event.Fields["nodeip"] = msg.Tags.K8sNodeIp
		event.Fields["podip"] = msg.Tags.ContainerIp

		event.Fields["message"] = msg.Contents.Content

		delete(event.Fields, "input")
	default:

	}

	return nil
}
