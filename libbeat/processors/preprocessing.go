package processors

import (
	"path"

	"github.com/bytedance/sonic"

	"github.com/elastic/beats/v7/libbeat/beat"
)

type LogFormat string

const (
	LogFormatIlogtail LogFormat = "ilogtail"
	LogFormatFilebeat LogFormat = "filebeat"
)

const (
	// LogFilename 日志文件名
	LogFilename = "file"
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

type FilebeatMessage struct {
	Timestamp string `json:"@timestamp"`
	Metadata  struct {
		Beat    string `json:"beat"`
		Type    string `json:"type"`
		Version string `json:"version"`
	} `json:"@metadata"`
	Log struct {
		File struct {
			Path string `json:"path"`
		} `json:"file"`
		Offset int `json:"offset"`
	} `json:"log"`
	Message string `json:"message"`
	Fields  struct {
		Servicetype string `json:"servicetype"`
	} `json:"fields"`
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
	case LogFormatFilebeat:
		var msg FilebeatMessage
		err = sonic.UnmarshalString(message, &msg)
		if err != nil {
			return err
		}
		event.Fields[LogFilename] = path.Base(msg.Log.File.Path)

		event.Fields["message"] = msg.Message

		delete(event.Fields, "input")
	default:

	}

	return nil
}
