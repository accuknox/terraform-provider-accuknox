package clienthandler

import (
	dsp "github.com/accuknox/auto-policy-discovery/pkg/discoveredpolicy/client/clientset/versioned/typed/security.kubearmor.com/v1"
	kspclient "github.com/kubearmor/KubeArmor/pkg/KubeArmorController/client/clientset/versioned"
	"k8s.io/client-go/kubernetes"
	ctrl "sigs.k8s.io/controller-runtime"
)

func ConnectK8sClient() (*kubernetes.Clientset, error) {
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

func ConnectKubearmorClient() (*kspclient.Clientset, error) {
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

func ConnectDiscoveryEngineClient() (*dsp.SecurityV1Client, error) {
	config, err := ctrl.GetConfig()
	if err != nil {
		return nil, err
	}

	DSPClient, err := dsp.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return DSPClient, nil
}
