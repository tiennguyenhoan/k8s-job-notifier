package config

import (
  "os"
  "strconv"
)

func IsNotifyFailOnly() bool {
  isNotifyFailOnly, err := strconv.ParseBool(os.Getenv("IS_NOTIFY_FAIL_ONLY"))
  if err != nil {
    return false
  }
  return isNotifyFailOnly
}

func IsInCluster() bool {
  inCluster, err := strconv.ParseBool(os.Getenv("IS_IN_CLUSTER"))
  if err != nil {
    return false
  }
  return inCluster
}

func GetNamespace() (namespace string) {
	if namespace = os.Getenv("JOB_NAMESPACE"); namespace == "" {
		namespace = "default"
	}
	return namespace
}

func GetClusterName() (clusterName string) {
	if clusterName = os.Getenv("CLUSTER_NAME"); clusterName == "" {
		clusterName = ""
	}
	return clusterName
}

func JobFailThreshold() int32 {
  failThreshold, err := strconv.ParseInt(os.Getenv("JOB_FAIL_THRESHOLD"), 10, 32)
  if err != nil {
    return int32(1)
  }
  return int32(failThreshold)
}

func GetJobFromLastMin() int32 {
  jobFromLastMin, err := strconv.ParseInt(os.Getenv("JOB_FROM_LAST_MIN"), 10, 32)
  if err != nil {
    return int32(5)
  }

  if jobFromLastMin < 1 {
    jobFromLastMin = 1
  }
  return int32(jobFromLastMin)
}

func GetSlackWebhook() (webhook string) {
  if webhook = os.Getenv("SLACK_WEBHOOK_URL"); webhook == "" {
    webhook = ""
  }
  return webhook
}

func GetSlackBotToken() (botToken string) {
  if botToken = os.Getenv("SLACK_BOT_TOKEN"); botToken == "" {
    botToken = ""
  }
  return botToken
}

func GetSlackDisplayUser() (displayUser string) {
  if displayUser = os.Getenv("SLACK_DISPLAY_USER"); displayUser == "" {
    displayUser = "Slack bot"
  }
  return displayUser
}

func GetSlackChannel() (slackChannel string) {
  if slackChannel = os.Getenv("SLACK_NOTIFY_CHANNEL"); slackChannel == "" {
    slackChannel = "#general"
  }
  return slackChannel
}

