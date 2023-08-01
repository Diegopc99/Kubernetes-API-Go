package KubeAPI

import (
	"context"
	"flag"
	"fmt"
	"path/filepath"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"

	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var KubeClient kubernetes.Clientset

func ParseKubeConfig() {

	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err)
	}
	KubeClient, err := kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}

	fmt.Println(KubeClient)
	fmt.Println("Kubernetes client configured!")

}

func int32Ptr(i int32) *int32 { return &i }

type ConfigDeployment struct {
	Name           string              `json:"name" binding:"required"`
	Replicas       int32               `json:"replicas"`
	ContainerName  string              `json:"containerName" binding:"required"`
	Image          string              `json:"image" binding:"required"`
	ContainerPorts []ContainerPortData `json:"containerPorts"`
	Labels         map[string]string   `json:"labels"`
}

type ContainerPortData struct {
	Name     string
	Protocol string
	Port     int32
}

func KubeCreateDeployment(configData ConfigDeployment) {

	deploymentsClient := KubeClient.AppsV1().Deployments(apiv1.NamespaceDefault)

	// Create Container Ports
	var containerPorts []apiv1.ContainerPort
	for _, portData := range configData.ContainerPorts {

		containerPort := apiv1.ContainerPort{
			Name:          portData.Name,
			Protocol:      apiv1.Protocol(portData.Protocol),
			ContainerPort: portData.Port,
		}

		containerPorts = append(containerPorts, containerPort)
	}

	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name: configData.Name,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &configData.Replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: configData.Labels,
			},
			Template: apiv1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: configData.Labels,
				},
				Spec: apiv1.PodSpec{
					Containers: []apiv1.Container{
						{
							Name:  configData.ContainerName,
							Image: configData.Image,
							Ports: containerPorts,
						},
					},
				},
			},
		},
	}

	// Create Deployment
	fmt.Println("Creating deployment...")
	result, err := deploymentsClient.Create(context.TODO(), deployment, metav1.CreateOptions{})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Created deployment %q.\n", result.GetObjectMeta().GetName())

}

func KubeDeleteDeployment() {

}
