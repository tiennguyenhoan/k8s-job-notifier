package main

import (
  "log"
  "strconv"
  "time"
  "os"
)

func main() {
  clientset, err := connectToK8s()
  if err != nil {
    log.Fatalf("Fail to connect to Cluster %v", err)
    os.Exit(1)
  }

  sc := SlackClient{
    WebHookUrl: getSlackWebhook(),
    UserName:   getSlackDisplayUser(),
    Channel:    getSlackChannel(),
    Bottoken:   getSlackBotToken(),
  }

  namespace := getNamespace()

  log.Printf("Fetching jobs from namespace: %s\n", namespace)

  jobs, err := listJobs(clientset, namespace)
  if err != nil {
    log.Fatalf("Failed to list all jobs in the namespace %v", err)
    os.Exit(0)
  }

  currentTime := time.Now().Unix()

	for _, job := range jobs.Items {
    jobStartTime := job.Status.StartTime

    // Only report failed job(s) of the last xxx minutes ago
    if job.Status.Failed >= getJobFailThreshold() && jobStartTime.Add(time.Duration(getJobFailedLastMin())*time.Minute).Unix() > currentTime {

      attachment := Attachment{
        Color: "#f7310a",               // Since this task is only for notify fail job only so we don't need to make a variable for it
        Pretext: "*Cluster Name*: " + getClusterName(),
        Text: "*Start time*: " + jobStartTime.Format("2006.01.02 15:04:05") + // Format time defined by golang, do not change it.
            "\nFail Num (time): " + strconv.FormatInt(int64(job.Status.Failed), 10),
        MarkdownIn: []string{"pretext", "text"},
      }

      sp := SlackPayload {
        Text: "Cron job fail on job `" + job.Name +"`",
        Attachments: []Attachment{attachment},
      }

      notifyErr := sc.sendNotification(sp)
      if notifyErr != nil {
        log.Fatalf("Failed to send notify to slack %v", notifyErr)
        os.Exit(1)
      }
    }
  }
}
