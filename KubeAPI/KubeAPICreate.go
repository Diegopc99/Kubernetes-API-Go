package KubeAPI

import (
	"context"
	"fmt"
	"strings"

	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ConfigDeploymentData struct {
	Namespace      string              `json:"namespace"`
	Name           string              `json:"name" binding:"required"`
	Replicas       int32               `json:"replicas"`
	ContainerName  string              `json:"containerName" binding:"required"`
	Image          string              `json:"image" binding:"required"`
	ContainerPorts []ContainerPortData `json:"containerPorts"`
	Labels         map[string]string   `json:"labels"`
}

type ContainerPortData struct {
	Name     string `json:"name" binding:"required"`
	Protocol string `json:"protocol" binding:"required"`
	Port     int32  `json:"port" binding:"required"`
}

func KubeCreateDeployment(configData ConfigDeploymentData) (*appsv1.Deployment, error) {

	deploymentsClient := KubeClient.AppsV1().Deployments(configData.Namespace)

	// Create Container Ports
	var containerPorts []apiv1.ContainerPort
	for _, portData := range configData.ContainerPorts {

		containerPort := apiv1.ContainerPort{
			Name:          portData.Name,
			Protocol:      apiv1.Protocol(strings.ToUpper(portData.Protocol)),
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
	result, err := deploymentsClient.Create(context.TODO(), deployment, metav1.CreateOptions{})
	if err != nil {
		return nil, err
	}
	fmt.Printf("Created deployment %q.\n", result.GetObjectMeta().GetName())

	return result, nil
}
