package kubernetes

import (
  "context"
  "os/user"
  "flag"
  batchv1 "k8s.io/api/batch/v1"
  metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
  "k8s.io/client-go/kubernetes"
  clientcmd "k8s.io/client-go/tools/clientcmd"
  rest "k8s.io/client-go/rest"
  c "k8s-job-notifier/config"
)

type Client struct {
  clientset kubernetes.Interface
}

func (c *Client) ListJobs(namespace string) (*batchv1.JobList, error){
  jobs, err := c.clientset.BatchV1().Jobs(namespace).List(context.Background(), metav1.ListOptions{})
  if err != nil {
		return nil, err
	}

	return jobs, nil
}

func ConnectToCluster() (*Client, error) {
  config, err := getConfig()
  if err != nil {
    return nil, err
  }

  clientset, err := kubernetes.NewForConfig(config)
  if err != nil {
    return nil, err
  }

  return &Client {
    clientset: clientset,
  }, nil
}

func getConfig() (config *rest.Config, err error) {
  if c.IsInCluster() {
    config, err = rest.InClusterConfig()
		if err != nil {
			return nil, err
		}
  } else {
    usr, err := user.Current()
    if err != nil {
      return nil, err
    }

    filePath := usr.HomeDir + "/.kube/config"
    kubeconfig := flag.String("kubeconfig", filePath, "absolute path to file")
    flag.Parse()

    config, err = clientcmd.BuildConfigFromFlags("", *kubeconfig)
    if err != nil {
      return nil, err
    }
  }

  return config, nil
}

