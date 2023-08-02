package KubeAPI

import (
	"flag"
	"fmt"
	"path/filepath"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

var KubeClient *kubernetes.Clientset

func ParseKubeConfig() {

	var kubeconfig *string
	if home := homedir.HomeDir(); home != "" {
		kubeconfig = flag.String("kubeconfig", filepath.Join(home, ".kube", "config"), "(optional) absolute path to the kubeconfig file")
	} else {
		fmt.Println("No HomeDir specified")
		kubeconfig = flag.String("kubeconfig", "", "absolute path to the kubeconfig file")
	}
	flag.Parse()

	config, err := clientcmd.BuildConfigFromFlags("", *kubeconfig)
	if err != nil {
		panic(err)
	}

	KubeClient, err = kubernetes.NewForConfig(config)
	if err != nil {
		panic(err)
	}
	fmt.Println("Kubernetes client configured!")

}

func int32Ptr(i int32) *int32 { return &i }
