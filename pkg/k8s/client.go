package k8s

import (
	"context"
	"fmt"

	appsV1 "k8s.io/api/apps/v1"
	coreV1 "k8s.io/api/core/v1"
	metaV1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

type Error string

func (e Error) Error() string { return string(e) }

const (
	NotImplementedYetErr = Error("function / feature not implemented yet")
	NotFoundErr          = Error("resource not found")
)

type K8sClient struct {
	Client *kubernetes.Clientset
}

func NewKubeClient() (*K8sClient, error) {
	var client K8sClient

	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}
	// creates the clientset
	client.Client, err = kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return &client, nil
}

func (kc *K8sClient) GetPods(namespace string, deplName string) (*coreV1.PodList, error) {

	listOpts := metaV1.ListOptions{}

	fmt.Printf("Selecting only pods for deployment %s\n", deplName)

	if deplName != "" {
		// Get deployment by name
		deployment, err := kc.getDeploymentByName(deplName, namespace)
		if err != nil {
			return nil, err
		}
		// read selector labels from deployment
		selector, err := metaV1.LabelSelectorAsSelector(deployment.Spec.Selector)
		if err != nil {
			return nil, err
		}
		listOpts.LabelSelector = selector.String()
	}

	// get pods in namespace and matching it against the label selector
	pods, err := kc.Client.CoreV1().Pods(namespace).List(context.TODO(), listOpts)
	if err != nil {
		fmt.Printf("reading pod info from cluster failed\n")
		return nil, err
	}
	fmt.Printf("There are %d pods in the cluster\n", len(pods.Items))

	return pods, nil
}

func (kc *K8sClient) getDeploymentByName(name, namespace string) (*appsV1.Deployment, error) {
	allDepls, err := kc.Client.AppsV1().Deployments(namespace).List(context.TODO(), metaV1.ListOptions{})
	if err != nil {
		return nil, err
	}
	for _, depl := range allDepls.Items {
		if depl.Name == name {
			return &depl, nil
		}
	}
	return nil, NotFoundErr
}
