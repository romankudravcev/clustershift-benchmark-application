package kubectl

import (
	"context"
	"log"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

var Clientset *kubernetes.Clientset

func ConnectK8s() {
	config, err := rest.InClusterConfig()
	if err != nil {
		log.Fatalf("Failed to create in-cluster config: %v", err)
	}

	Clientset, err = kubernetes.NewForConfig(config)
	if err != nil {
		log.Fatalf("Failed to create clientset: %v", err)
	}
}

func GetExternalIP() string {
	services, err := Clientset.CoreV1().Services("kube-system").List(context.TODO(), metav1.ListOptions{
		FieldSelector: "metadata.name=traefik",
	})

	if err != nil {
		log.Fatalf("Failed to get service: %v", err)
		return ""
	}

	if len(services.Items) > 0 &&
		len(services.Items[0].Status.LoadBalancer.Ingress) > 0 {
		ingressIP := services.Items[0].Status.LoadBalancer.Ingress[0].IP
		if ingressIP != "" {
			if ingressIP == "192.168.131.10" {
				return "141.13.5.206"
			}
			return "141.13.5.205"
		}
	}

	return ""
}
