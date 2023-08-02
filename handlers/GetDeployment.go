package handlers

import (
	"net/http"

	"github.com/Kubernetes-API-Go/KubeAPI"
	"github.com/Kubernetes-API-Go/exceptions"
	"github.com/gin-gonic/gin"
)

func GetDeployment(c *gin.Context) {

	namespace := c.Query("namespace")

	if namespace == "" {
		namespace = "default"
	} else {
		err := KubeAPI.NamespaceExists(namespace)
		if err != nil {
			c.Error(&exceptions.InvalidRequest{})
			return
		}
	}

	list, err := KubeAPI.KubeGetDeployments(namespace)
	if err != nil {
		c.Error(&exceptions.InternalError{})
		return
	}

	c.JSON(http.StatusOK, gin.H{"deployments": list})
}
