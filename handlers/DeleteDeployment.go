package handlers

import (
	"github.com/Kubernetes-API-Go/KubeAPI"
	"github.com/gin-gonic/gin"
)

func DeleteDeployment(c *gin.Context) {

	KubeAPI.KubeDeleteDeployment()

}
