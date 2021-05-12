package main

import (
  "log"
  "strconv"
  "time"
  "os"
  k8s "k8s-job-notifier/kubernetes"
  c "k8s-job-notifier/config"
  slack "k8s-job-notifier/slack"
)

func main() {
  client, err := k8s.ConnectToCluster()
  if err != nil {
    log.Fatalf("Fail to connect to Cluster %v", err)
    os.Exit(1)
  }

  sc := slack.SlackClient{
    WebHookUrl: c.GetSlackWebhook(),
    UserName:   c.GetSlackDisplayUser(),
    Channel:    c.GetSlackChannel(),
    Bottoken:   c.GetSlackBotToken(),
  }

  namespace := c.GetNamespace()

  log.Printf("Fetching jobs from namespace: %s\n", namespace)

  jobs, err := client.ListJobs(namespace)
  if err != nil {
    log.Fatalf("Failed to list all jobs in the namespace %v", err)
    os.Exit(0)
  }

  jobCheckingTime := time.Now().Unix()

	for _, job := range jobs.Items {
    var sp *slack.SlackPayload

    jobStartTime := job.Status.StartTime
    jobStartTimeFormmated := jobStartTime.Format("2006.01.02 15:04:05") // Format time defined by golang, do not change it.

    // Only pick job form the last xxx minutes
    if jobStartTime.Add(time.Duration(c.GetJobFromLastMin())*time.Minute).Unix() > jobCheckingTime {
      if job.Status.Failed >= c.JobFailThreshold() {

        attachment := slack.Attachment{
          Color: "#f7310a",                                         // #Red for error
          Text: "*Cluster Name*:   " + c.GetClusterName() + "\n" +
              "*Start time*:   " + jobStartTimeFormmated + "\n" +
              "Fail time(s):   " + strconv.FormatInt(int64(job.Status.Failed), 10),
          MarkdownIn: []string{"pretext", "text"},
        }

        sp = &slack.SlackPayload {
          Text: "[FAIL] Job `" + job.Name +"`",
          Attachments: []slack.Attachment{attachment},
        }

      } else if job.Status.Succeeded > 0 && !c.IsNotifyFailOnly() {

        attachment := slack.Attachment{
          Color: "#03fc28",                                         // #Green for success 
          Text: "*Cluster Name*:   " + c.GetClusterName() + "\n" +
              "*Start time*:   " + jobStartTimeFormmated,
          MarkdownIn: []string{"pretext", "text"},
        }

        sp = &slack.SlackPayload {
          Text: "[SUCCESS] Job `" + job.Name +"`",
          Attachments: []slack.Attachment{attachment},
        }
      }

      notifyErr := sc.SendNotification(*sp)
      if notifyErr != nil {
        log.Fatalf("Failed to send notify to slack %v", notifyErr)
        os.Exit(1)
      }
    }
  }
}
