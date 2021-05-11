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

  jobCheckingTime := time.Now().Unix()

	for _, job := range jobs.Items {
    var sp *SlackPayload

    jobStartTime := job.Status.StartTime

    if jobStartTime.Add(time.Duration(getJobFromLastMin())*time.Minute).Unix() > jobCheckingTime {
      if job.Status.Failed >= jobFailThreshold() {

        attachment := Attachment{
          Color: "#f7310a",                                                            // #Red for error
          Text: "*Cluster Name*:   " + getClusterName() + "\n" +
              "*Start time*:   " + jobStartTime.Format("2006.01.02 15:04:05") + "\n" + // Format time defined by golang, do not change it.
              "Fail time(s):   " + strconv.FormatInt(int64(job.Status.Failed), 10),
          MarkdownIn: []string{"pretext", "text"},
        }

        sp = &SlackPayload {
          Text: "[FAIL] job `" + job.Name +"`",
          Attachments: []Attachment{attachment},
        }

      } else if job.Status.Succeeded > 0 && !isNotifyFailOnly() {

        attachment := Attachment{
          Color: "#03fc28",                                                           // #Green for success 
          Text: "*Cluster Name*:   " + getClusterName() + "\n" +
              "*Start time*:   " + jobStartTime.Format("2006.01.02 15:04:05"),        // Format time defined by golang, do not change it.
          MarkdownIn: []string{"pretext", "text"},
        }

        sp = &SlackPayload {
          Text: "[SUCCESS] job `" + job.Name +"`",
          Attachments: []Attachment{attachment},
        }
      }

      notifyErr := sc.sendNotification(*sp)
      if notifyErr != nil {
        log.Fatalf("Failed to send notify to slack %v", notifyErr)
        os.Exit(1)
      }
    }
  }
}
