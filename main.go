package main

import (
	"context"
	"fmt"
	"log"

	corev1 "k8s.io/api/core/v1"
	ctrl "sigs.k8s.io/controller-runtime"
	kclient "sigs.k8s.io/controller-runtime/pkg/client"
)

func main() {
	log.Println("Creating the client")

	client, err := kclient.New(ctrl.GetConfigOrDie(), kclient.Options{})
	if err != nil {
	        log.Fatalf("Couldn't create client: %v", err)
	}

	log.Println("Looking for kube-apiserver pods...")
	pods := corev1.PodList{}
	err = client.List(context.TODO(), &pods, &kclient.ListOptions{Namespace: "openshift-kube-apiserver"})
	if err != nil {
		log.Fatalf("Couldn't list pods in 'openshift-kube-apiserver': %v", err)
	}
	log.Println("Pods in 'openshift-kube-apiserver': ")
	for _, pod := range pods.Items {
		fmt.Println("- ", pod.Name)
	}
}
