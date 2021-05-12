# Job fail notifier for Kubernetes

## Idea & Croncept

Basically, this tool will fetch the list of jobs from Kubernetes cluster and send alert about the job status

- The criteria for alert job(s) will be defined by:
  
  + The last execution time of a job, by default the application will detect job(s) that executed 5 mins ago.
  
  + Job status (success or fail), for the fail status notify we can define the number of failed times on job before sending the notification. 

- The tool support sending alert through Slack, in two way
  
  + Using Slack webhook
  
  + Slack Bot Token
  
- This also supports using `InClusterConfig` to access Kubernetes API when deploying to Kubernetes cluster

> **Limitation**: Since this tool can only get all jobs in a specified namespace. Therefore, we may need to deploy this tool to multiple times if we have multiple namespaces

## Requirement

- Golang: 1.16.2

## Variables

- `IS_NOTIFY_FAIL_ONLY`: Set this to **true** if we only need to alert the failed job (Default: false)

- `IS_IN_CLUSTER`: For deploying to K8s cluster only, set to **true** and the tool will read config from cluster instead of "**.kube/config**" (Default: ~/.kube/config)

- `JOB_NAMESPACE`: The namespace of checking job(s)/cronjob(s). (Default: default)

- `CLUSTER_NAME`: For Slack alert only, this help us know the alert for cluster

- `JOB_FAIL_THRESHOLD`: Job failed time(s) before notify to slack (Default: 5)

- `JOB_FROM_LAST_MIN`: Only pick and notify job executed from last xxx minutes (Default: 5)

    > If we change the value of `JOB_FROM_LAST_MIN`, make sure to check the cronjob execution time on [cronjob.yaml](./deployment/cronjob.yaml)

- `SLACK_WEBHOOK_URL`: The webhook URL for sending alert to slack

- `SLACK_BOT_TOKEN`: The slack bot token for sending alert to slack, the bot will need **chat:write** permission

    > If we specified both `SLACK_WEBHOOK_URL` and `SLACK_BOT_TOKEN`, we will pick Webhook instead of bot token

- `SLACK_DISPLAY_USER`: Display user on Slack alert message  (Default: Slack bot)

- `SLACK_NOTIFY_CHANNEL`: Channel that alert message will be posted to (Default: general)

## Usage

### Local Testing

#### (Optional) Deploy testing materials

I have created a sample that contain 2 cronjobs (*success and failure)*, this helps create a material for us to test when we don't have the cronjob yet. 

1. Navigate to folder [deployment](./deployment)

1. Deploy [job-tester.yaml](./deployment/job-tester.yaml) file to your kubernetes cluster

    ```bash
      kubectl apply -f job-tester.yaml
    ```

#### Build and execution on local

1. Prepare the connection to Kubernetes cluster

1. Export the necessary variables

1. Build and execute the tool

    ```bash
      go build -o ./app .

      ./app
    ```

### Kubernetes Deployment

Navigate to the deployment folder, deploy the template to delicated namespace. It's working as a cronjob on the cluster 

> **Note**: If we want to change the cronjob execution time, please make sure to update `SERVICE_FAILED_FROM_LAST_MIN` variable on configmap
