package KubeAPI

import (
	"context"
	"fmt"

	apiv1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type DeleteDeploymentData struct {
	DeploymentName string `json:"deploymentName" binding:"required"`
}

func KubeDeleteDeployment(deploymentData DeleteDeploymentData) {

	deploymentsClient := KubeClient.AppsV1().Deployments(apiv1.NamespaceDefault)

	deletePolicy := metav1.DeletePropagationForeground
	if err := deploymentsClient.Delete(context.TODO(), deploymentData.DeploymentName, metav1.DeleteOptions{
		PropagationPolicy: &deletePolicy,
	}); err != nil {
		panic(err)
	}
	fmt.Println("Deleted deployment.")

}
