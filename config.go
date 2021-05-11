package main

import (
  "os"
  "strconv"
)

func getNamespace() (namespace string) {
	if namespace = os.Getenv("SERVICE_NAMESPACE"); namespace == "" {
		namespace = "default"
	}
	return namespace
}

func getClusterName() (clusterName string) {
	if clusterName = os.Getenv("SERVICE_CLUSTER_NAME"); clusterName == "" {
		clusterName = ""
	}
	return clusterName
}

func getJobFailThreshold() int32 {
  failThreshold, err := strconv.ParseInt(os.Getenv("SERVICE_FAIL_THRESHOLD"), 10, 32)
  if err != nil {
    return int32(1)
  }
  return int32(failThreshold)
}

func getJobFailedLastMin() int32 {
  failedFromLastMin, err := strconv.ParseInt(os.Getenv("SERVICE_FAILED_FROM_LAST_MIN"), 10, 32)
  if err != nil {
    return int32(5)
  }

  if failedFromLastMin < 1 {
    failedFromLastMin = 1
  }
  return int32(failedFromLastMin)
}

func isInCluster() bool {
  inCluster, err := strconv.ParseBool(os.Getenv("IS_IN_CLUSTER"))
  if err != nil {
    return false
  }
  return inCluster
}

func getSlackWebhook() (webhook string) {
  if webhook = os.Getenv("SLACK_WEBHOOK_URL"); webhook == "" {
    webhook = ""
  }
  return webhook
}

func getSlackBotToken() (botToken string) {
  if botToken = os.Getenv("SLACK_BOT_TOKEN"); botToken == "" {
    botToken = ""
  }
  return botToken
}

func getSlackDisplayUser() (displayUser string) {
  if displayUser = os.Getenv("SLACK_DISPLAY_USER"); displayUser == "" {
    displayUser = "Slack bot"
  }
  return displayUser
}

func getSlackChannel() (slackChannel string) {
  if slackChannel = os.Getenv("SLACK_NOTIFY_CHANNEL"); slackChannel == "" {
    slackChannel = "#general"
  }
  return slackChannel
}

