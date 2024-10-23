/*
Copyright 2024.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package controller

import (
	"context"
	"fmt"
	"strings"
	"time"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"

	"istio.io/client-go/pkg/apis/networking/v1beta1"
	versionedclient "istio.io/client-go/pkg/clientset/versioned"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/clientcmd"

	registryv1alpha1 "github.com/btjimerson/consul-registry-sync/api/v1alpha1"

	consul "github.com/hashicorp/consul/api"
)

// ConsulRegistrySyncReconciler reconciles a ConsulRegistrySync object
type ConsulRegistrySyncReconciler struct {
	client.Client
	Scheme *runtime.Scheme
}

//+kubebuilder:rbac:groups=registry.solo.io,resources=consulregistrysyncs,verbs=get;list;watch;create;update;patch;delete
//+kubebuilder:rbac:groups=registry.solo.io,resources=consulregistrysyncs/status,verbs=get;update;patch
//+kubebuilder:rbac:groups=registry.solo.io,resources=consulregistrysyncs/finalizers,verbs=update

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the ConsulRegistrySync object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.17.3/pkg/reconcile
func (r *ConsulRegistrySyncReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = log.FromContext(ctx)

	consulRegistrySync := &registryv1alpha1.ConsulRegistrySync{}

	err := r.Get(ctx, req.NamespacedName, consulRegistrySync)
	if err != nil {
		return ctrl.Result{RequeueAfter: 30 * time.Second}, err
	}

	serviceCatalog, err := GetConsulServices(consulRegistrySync.Spec.ConsulAddress)
	for key, value := range serviceCatalog {
		CreateServiceEntry(consulRegistrySync.Spec.KubeConfig, key, strings.Fields(value), consulRegistrySync.Spec.ServiceEntryNamespace)
		fmt.Printf("%s address = %s\n", key, value)
	}

	return ctrl.Result{RequeueAfter: 30 * time.Second}, err
}

// SetupWithManager sets up the controller with the Manager.
func (r *ConsulRegistrySyncReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&registryv1alpha1.ConsulRegistrySync{}).
		Complete(r)
}

// Gets a map of all Consul service. Key is the service name and value is the service address
func GetConsulServices(consulAddress string) (map[string]string, error) {

	config := consul.DefaultConfig()
	config.Address = consulAddress
	client, err := consul.NewClient(config)
	if err != nil {
		return nil, err
	}

	services, _, err := client.Catalog().Services(nil)
	if err != nil {
		return nil, err
	}

	serviceCatalog := make(map[string]string)
	for serviceName := range services {
		service, _, err := client.Catalog().Service(serviceName, "", nil)
		if err != nil {
			return nil, err
		}
		serviceAddress := service[0].ServiceMeta["k8s-service-name"] + "." + service[0].ServiceMeta["k8s-namespace"] + ".svc"
		serviceCatalog[serviceName] = serviceAddress
	}

	return serviceCatalog, err
}

// Creates a ServiceEntry for the name and host
func CreateServiceEntry(kubeconfig string, name string, hosts []string, namespace string) error {
	restConfig, err := clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return err
	}

	ic, err := versionedclient.NewForConfig(restConfig)
	if err != nil {
		return err
	}

	serviceEntry := &v1beta1.ServiceEntry{}
	serviceEntry.Name = name
	serviceEntry.Spec.Hosts = hosts
	se, err := ic.NetworkingV1beta1().ServiceEntries(namespace).Create(context.TODO(), serviceEntry, v1.CreateOptions{})

	if err != nil {
		panic(err)
	}

	fmt.Printf("service entry created %s", se.Name)

	return nil
}
