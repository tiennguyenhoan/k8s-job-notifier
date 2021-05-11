# Job fail notifier for Kubernetes

## Idea & Croncept

Basically, this application will fetch the list jobs from Kubernetes cluster and report failed job(s) to Slack. 

- The citeria for notify job will be defined by number of fail and the last execution time, by default the application will detect job(s) that failed for more than 5 times and last execution is 5 mins ago.

- This application only support send notification through slack bot, which mean we have to setup a slack bot and config the bot token to use this app.

> **Limitation**: Since this tool can only get all jobs in it own namespace so we have to deploy it to all namespaces

## Requirement

- Golang: 1.16.2

## Usage

### Local Development

```bash
  
  cd jobs-fail-notifier

  go build

  export SERVICE_FAIL_THRESHOLD="1"; \
  export SERVICE_FAILED_FROM_LAST_MIN="30" \
  export SLACK_BOT_TOKEN="<Your slack bot token>"; \
  export SLACK_DISPLAY_USER="Kubernetes job notifier"; \
  export SERVICE_CLUSTER_NAME="My Kubernetes cluster"

  ./jobs-fail-notifier

```

### Kubernetes Deployment

Navigate to the kubernetes folder, deploy the template to delicated namespace. It's working as a cronjob on the cluster 

> **Note**: If we want to change the cronjob execution time, please make sure to update `SERVICE_FAILED_FROM_LAST_MIN` variable on configmap
