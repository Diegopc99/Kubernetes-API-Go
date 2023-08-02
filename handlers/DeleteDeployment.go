package handlers

import (
	"net/http"

	"github.com/Kubernetes-API-Go/KubeAPI"
	"github.com/Kubernetes-API-Go/exceptions"
	"github.com/gin-gonic/gin"
)

func DeleteDeployment(c *gin.Context) {

	deploymentName := c.Param("deploymentName")
	namespace := c.Query("namespace")

	if deploymentName == "" {
		c.Error(&exceptions.InvalidRequest{})
		return
	}

	if namespace == "" {
		namespace = "default"
	} else {
		err := KubeAPI.NamespaceExists(namespace)
		if err != nil {
			c.Error(&exceptions.InvalidRequest{})
			return
		}
	}

	deploymentList, err := KubeAPI.KubeGetDeployments(namespace)
	if err != nil {
		c.Error(&exceptions.InternalError{})
		return
	}

	// Check if deployment exists
	for _, deployment := range deploymentList.Items {
		if deployment.Name == deploymentName {
			err = KubeAPI.KubeDeleteDeployment(namespace, deploymentName)
			if err != nil {
				c.Error(&exceptions.InternalError{})
				return
			} else {
				// Successfully deleted deployment
				c.JSON(http.StatusOK, gin.H{"message": "Item deleted successfully"})
				return
			}
		}
	}

	// Deployment doesn't exist
	c.Error(&exceptions.ResourceNotFound{})
	return
}
