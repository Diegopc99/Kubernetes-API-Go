package KubeAPI

import (
	"context"
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func KubeDeleteDeployment(namespace string, deploymentData string) error {

	deploymentsClient := KubeClient.AppsV1().Deployments(namespace)

	fmt.Println(deploymentData)

	deletePolicy := metav1.DeletePropagationForeground
	if err := deploymentsClient.Delete(context.TODO(), deploymentData, metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	}); err != nil {
		return err
	}
	fmt.Println("Deleted deployment.")

	return nil

}
