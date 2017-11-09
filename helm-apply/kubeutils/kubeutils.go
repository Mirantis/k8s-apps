package kubeutils

import (
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/pkg/api/v1"
	"k8s.io/client-go/pkg/apis/rbac/v1beta1"
)

func EnsureNamespace(client kubernetes.Interface, namespace string) error {
	_, err := client.CoreV1().Namespaces().Get(namespace, metav1.GetOptions{})
	if err != nil && errors.IsNotFound(err) {
		ns := v1.Namespace{
			ObjectMeta: metav1.ObjectMeta{
				Name: namespace,
			},
		}
		_, err := client.CoreV1().Namespaces().Create(&ns)
		if err != nil {
			return err
		}
	}
	return nil
}

func CreateServiceAccount(client kubernetes.Interface, namespace string) error {
	sa := v1.ServiceAccount{
		ObjectMeta: metav1.ObjectMeta{
			Name: "tiller-" + namespace,
		},
	}
	_, err := client.CoreV1().ServiceAccounts(namespace).Create(&sa)
	if err != nil {
		return err
	}
	return nil
}

func CreateClusterRoleBinding(client kubernetes.Interface, namespace string) error {
	crb := v1beta1.ClusterRoleBinding{
		ObjectMeta: metav1.ObjectMeta{
			Name: "tiller-" + namespace,
		},
		RoleRef: v1beta1.RoleRef{
			Kind: "ClusterRole",
			Name: "cluster-admin",
		},
		Subjects: []v1beta1.Subject{{
			Kind:      "ServiceAccount",
			Name:      "tiller",
			Namespace: namespace,
		}},
	}
	_, err := client.RbacV1beta1().ClusterRoleBindings().Create(&crb)
	if err != nil {
		return err
	}
	return nil
}
