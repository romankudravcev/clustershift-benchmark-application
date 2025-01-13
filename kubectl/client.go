package kubectl

import (
	"log"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/metrics/pkg/client/clientset/versioned"
)

var Clientset *kubernetes.Clientset
var MetricsClientset *versioned.Clientset

func ConnectK8s() {
    config, err := rest.InClusterConfig()
    if err != nil {
        log.Fatalf("Failed to create in-cluster config: %v", err)
    }

    Clientset, err = kubernetes.NewForConfig(config)
    if err != nil {
        log.Fatalf("Failed to create clientset: %v", err)
    }

    MetricsClientset, err = versioned.NewForConfig(config)
    if err != nil {
        log.Fatalf("Failed to create metrics clientset: %v", err)
    }
}
