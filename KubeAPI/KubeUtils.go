package KubeAPI

import (
	"context"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type CustomError string

func (c CustomError) Error() string {
	return string(c)
}

func isNamespaceInList(kubeNamespaceList *v1.NamespaceList, namespace string) bool {
	for _, kubeNamespace := range kubeNamespaceList.Items {
		if kubeNamespace.Name == namespace {
			return true
		}
	}
	return false
}

func NamespaceExists(namespace string) error {

	kubeNamespaceList, err := KubeClient.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return err
	}

	exists := isNamespaceInList(kubeNamespaceList, namespace)

	if exists != true {
		return CustomError("Selected namespace doesn't exist")
	}

	return nil
}
