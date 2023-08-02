package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"regexp"
	"strings"

	"github.com/Kubernetes-API-Go/KubeAPI"
	"github.com/Kubernetes-API-Go/exceptions"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
)

func isValidProtocol(protocol string) bool {
	// A valid protocol can be "TCP", "UDP", or "SCTP".
	validProtocolRegex := regexp.MustCompile(`^(TCP|UDP|SCTP)$`)
	return validProtocolRegex.MatchString(strings.ToUpper(protocol))
}

func isValidName(name string) bool {
	// A valid name should consist of lowercase alphanumeric characters or '-', and it should start and end with an alphanumeric character.
	validNameRegex := regexp.MustCompile(`^[a-z0-9]([-a-z0-9]*[a-z0-9])?$`)
	return validNameRegex.MatchString(name)
}

func CreateDeployment(c *gin.Context) {

	var configDeployment KubeAPI.ConfigDeploymentData

	// Process request params with required field validation
	if err := c.BindJSON(&configDeployment); err != nil {
		var verr validator.ValidationErrors
		if errors.As(err, &verr) {
			c.Error(&exceptions.MissingFields{})
			return
		}
		c.Error(&exceptions.InvalidRequest{})
		return
	}

	if isValidName(configDeployment.ContainerName) == false {
		c.Error(&exceptions.InvalidRequest{})
		return
	}

	for _, portData := range configDeployment.ContainerPorts {
		fmt.Println(portData.Port)
		if isValidProtocol(portData.Protocol) == false {
			c.Error(&exceptions.InvalidRequest{})
			return
		}
	}

	// Check namespace
	if configDeployment.Namespace == "" {
		configDeployment.Namespace = "default"
	} else {
		err := KubeAPI.NamespaceExists(configDeployment.Namespace)
		if err != nil {
			c.Error(&exceptions.InvalidRequest{})
			return
		}
	}

	// Verify if a deployment already exists with that name
	deploymentList, err := KubeAPI.KubeGetDeployments(configDeployment.Namespace)
	if err != nil {
		c.Error(&exceptions.InternalError{})
		return
	}

	for _, deployment := range deploymentList.Items {
		if deployment.Name == configDeployment.Name {
			c.Error(&exceptions.InvalidRequest{})
			return
		}
	}

	// Create deployment
	deployment, err := KubeAPI.KubeCreateDeployment(configDeployment)
	if err != nil {
		c.Error(&exceptions.InternalError{})
		return
	}

	c.JSON(http.StatusOK, gin.H{"deployment": deployment})

}
