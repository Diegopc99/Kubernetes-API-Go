package handlers

import (
	"errors"

	"github.com/Kubernetes-API-Go/KubeAPI"
	"github.com/Kubernetes-API-Go/exceptions"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func DeleteDeployment(c *gin.Context) {

	var data KubeAPI.DeleteDeploymentData

	// Process request params with required field validation
	if err := c.BindJSON(&data); err != nil {
		var verr validator.ValidationErrors
		if errors.As(err, &verr) {
			c.Error(&exceptions.MissingFields{})
			return
		}
		c.Error(&exceptions.InvalidRequest{})
		return
	}

	KubeAPI.KubeDeleteDeployment(data)

}
