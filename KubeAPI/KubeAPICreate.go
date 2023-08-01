package KubeAPI

import (
	"context"
	"fmt"

	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ConfigDeploymentData struct {
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

func KubeCreateDeployment(configData ConfigDeploymentData) {

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
