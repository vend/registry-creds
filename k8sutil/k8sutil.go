package k8sutil

import (
	"github.com/Sirupsen/logrus"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/cache"
	"k8s.io/client-go/tools/clientcmd"
)

// KubeInterface abstracts the k8s api
type KubeInterface interface {
	CreateSecret(namespace string, secret *v1.Secret) error
	GetSecret(namespace, secretname string) (*v1.Secret, error)
	UpdateSecret(namespace string, secret *v1.Secret) error

	GetNamespaces() (*v1.NamespaceList, error)

	GetServiceAccount(namespace, name string) (*v1.ServiceAccount, error)
	UpdateServiceAccount(namespace string, sa *v1.ServiceAccount) error
	// WatchAll(labelSelector string, stopCh <-chan struct{}) (<-chan interface{}, error)
}

type kubeImpl struct {
	secretsController *cache.Controller
	secretsStore      cache.Store

	Clientset *kubernetes.Clientset
}

// New returns a new instance of KubeInterface
func New(kubeCfgFile string) (KubeInterface, error) {

	var client *kubernetes.Clientset

	// Should we use in cluster or out of cluster config
	if len(kubeCfgFile) == 0 {
		logrus.Info("Using InCluster k8s config")
		cfg, err := rest.InClusterConfig()

		if err != nil {
			return nil, err
		}

		client, err = kubernetes.NewForConfig(cfg)

		if err != nil {
			return nil, err
		}
	} else {
		logrus.Infof("Using OutOfCluster k8s config with kubeConfigFile: %s", kubeCfgFile)
		cfg, err := clientcmd.BuildConfigFromFlags("", kubeCfgFile)

		if err != nil {
			logrus.Error("Got error trying to create client: ", err)
			return nil, err
		}

		client, err = kubernetes.NewForConfig(cfg)

		if err != nil {
			return nil, err
		}
	}

	return &kubeImpl{
		Clientset: client,
	}, nil
}

// // WatchSecrets starts the watch of Kubernetes secrets resources and updates the corresponding store
// func (k *kubeImpl) WatchSecrets(labelSelector labels.Selector, watchCh chan<- interface{}, stopCh <-chan struct{}) {
// 	source := NewListWatchFromClient(
// 		k.Kclient.ExtensionsClient,
// 		"ingresses",
// 		api.NamespaceAll,
// 		fields.Everything(),
// 		labelSelector)

// 	k.secretsStore, k.secretsController = cache.NewListWatchFromClient(
// 		source,
// 		&v1beta1.Ingress{},
// 		resyncPeriod,
// 		newResourceEventHandlerFuncs(watchCh))
// 	go c.ingController.Run(stopCh)
// }

// GetNamespaces returns all namespaces
func (k *kubeImpl) GetNamespaces() (*v1.NamespaceList, error) {
	namespaces, err := k.Clientset.Namespaces().List(v1.ListOptions{})
	if err != nil {
		logrus.Error("Error getting namespaces: ", err)
		return nil, err
	}

	return namespaces, nil
}

// GetSecret get a secret
func (k *kubeImpl) GetSecret(namespace, secretname string) (*v1.Secret, error) {
	secret, err := k.Clientset.Secrets(namespace).Get(secretname)
	if err != nil {
		logrus.Error("Error getting secret: ", err)
		return nil, err
	}

	return secret, nil
}

// CreateSecret creates a secret
func (k *kubeImpl) CreateSecret(namespace string, secret *v1.Secret) error {
	_, err := k.Clientset.Secrets(namespace).Create(secret)

	if err != nil {
		logrus.Error("Error creating secret: ", err)
		return err
	}

	return nil
}

// UpdateSecret updates a secret
func (k *kubeImpl) UpdateSecret(namespace string, secret *v1.Secret) error {
	_, err := k.Clientset.Secrets(namespace).Update(secret)

	if err != nil {
		logrus.Error("Error updating secret: ", err)
		return err
	}

	return nil
}

// GetServiceAccount updates a secret
func (k *kubeImpl) GetServiceAccount(namespace, name string) (*v1.ServiceAccount, error) {
	sa, err := k.Clientset.ServiceAccounts(namespace).Get(name)

	if err != nil {
		logrus.Error("Error getting service account: ", err)
		return nil, err
	}

	return sa, nil
}

// UpdateServiceAccount updates a secret
func (k *kubeImpl) UpdateServiceAccount(namespace string, sa *v1.ServiceAccount) error {
	_, err := k.Clientset.ServiceAccounts(namespace).Update(sa)

	if err != nil {
		logrus.Error("Error updating service account: ", err)
		return err
	}

	return nil
}
