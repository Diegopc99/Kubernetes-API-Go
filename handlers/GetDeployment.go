package handlers

import (
	"github.com/Kubernetes-API-Go/KubeAPI"
	"github.com/gin-gonic/gin"
)

func GetDeployment(c *gin.Context) {

	KubeAPI.KubeGetDeployment()

}
