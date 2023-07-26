package kubearmor

import (
	kspclient "github.com/kubearmor/KubeArmor/pkg/KubeArmorController/client/clientset/versioned"
	"k8s.io/client-go/kubernetes"
	ctrl "sigs.k8s.io/controller-runtime"
)

func connectK8sClient() (*kubernetes.Clientset, error) {
	config, err := ctrl.GetConfig()
	if err != nil {
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return clientset, nil
}

func connectKubearmorClient() (*kspclient.Clientset, error) {
	config, err := ctrl.GetConfig()
	if err != nil {
		return nil, err
	}

	KSPClient, err := kspclient.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return KSPClient, nil
}
