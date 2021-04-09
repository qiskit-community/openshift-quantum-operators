/*
Copyright 2021.

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

package controllers

import (
	"context"

	"strings"

	"github.com/go-logr/logr"
	appsv1 "k8s.io/api/apps/v1"
	apiv1 "k8s.io/api/core/v1"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/rest"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"

	"k8s.io/apimachinery/pkg/api/equality"

	routev1 "github.com/openshift/api/route/v1"
	qiskitv1alpha1 "github.io/blublinsky/qiskitplaygrounds/api/v1alpha1"
)

// getDiscoveryClient returns a discovery client for the current reconciler
func getDiscoveryClient(config *rest.Config) (*discovery.DiscoveryClient, error) {
	return discovery.NewDiscoveryClientForConfig(config)
}

// Check where we are running
func GetClusterType(logger logr.Logger) bool {
	// The discovery package is used to discover APIs supported by a Kubernetes API server.
	config, err := ctrl.GetConfig()
	if err == nil && config != nil {
		dclient, err := getDiscoveryClient(config)
		if err == nil && dclient != nil {
			apiGroupList, err := dclient.ServerGroups()
			if err != nil {
				logger.Info("Error while querying ServerGroups, assuming we're on Vanilla Kubernetes")
				return false
			} else {
				for i := 0; i < len(apiGroupList.Groups); i++ {
					if strings.HasSuffix(apiGroupList.Groups[i].Name, ".openshift.io") {
						logger.Info("We detected being on OpenShift! Wouhou!")
						return true
					}
				}
				return false
			}
		} else {
			logger.Info("Cannot retrieve a DiscoveryClient, assuming we're on Vanilla Kubernetes")
			return false
		}
	} else {
		logger.Info("Cannot retrieve config, assuming we're on Vanilla Kubernetes")
		return false
	}
}

// QiskitPlaygroundReconciler reconciles a QiskitPlayground object
type QiskitPlaygroundReconciler struct {
	client.Client
	Log         logr.Logger
	Scheme      *runtime.Scheme
	IsOpenShift bool
}

// +kubebuilder:rbac:groups=qiskit.ibm.com,resources=qiskitplaygrounds,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=qiskit.ibm.com,resources=qiskitplaygrounds/status,verbs=get;update;patch
// +kubebuilder:rbac:groups=qiskit.ibm.com,resources=qiskitplaygrounds/finalizers,verbs=update

// +kubebuilder:rbac:groups=core,resources=namespaces,verbs=get;watch;list
// +kubebuilder:rbac:groups=apps,resources=deployments,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=core,resources=services,verbs=get;list;watch;create;update;patch;delete
// +kubebuilder:rbac:groups=route.openshift.io,resources=routes,verbs=get;list;watch;create;update;patch;delete

// Reconcile is part of the main kubernetes reconciliation loop which aims to
// move the current state of the cluster closer to the desired state.
// TODO(user): Modify the Reconcile function to compare the state specified by
// the QiskitPlayground object against the actual cluster state, and then
// perform operations to make the cluster state reflect the state specified by
// the user.
//
// For more details, check Reconcile and its Result here:
// - https://pkg.go.dev/sigs.k8s.io/controller-runtime@v0.7.0/pkg/reconcile
func (r *QiskitPlaygroundReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	_ = r.Log.WithValues("qiskitplayground", req.NamespacedName)

	// Get the qiskit playground to reconcile.
	var playground qiskitv1alpha1.QiskitPlayground
	if err := r.Get(ctx, req.NamespacedName, &playground); err != nil {
		//r.Log.Error(err, "Unable to fetch Playground")
		return ctrl.Result{}, client.IgnoreNotFound(err)
	}

	//r.Log.Info("New reconciliation",
	//	"Playground.Namespace", playground.Namespace, "Playground.Name", playground.Name)

	// Support variables
	var err error
	labels := map[string]string{
		"app":        "qiskit",
		"playground": playground.Name,
	}

	// Calculate desired deployment
	deployment := &appsv1.Deployment{}
	deployment, err = r.newDeploymentForPlayground(&playground, &labels)
	if err != nil {
		return ctrl.Result{}, err
	}

	// Get current deployment
	currentdeployment := &appsv1.Deployment{}
	err = r.Get(ctx, types.NamespacedName{Name: playground.Name + "-deployment", Namespace: playground.Namespace}, currentdeployment)

	// if the Deployment doesn't exist create it
	if err != nil && errors.IsNotFound(err) {
		r.Log.Info("Creating a new Deployment",
			"Deployment.Namespace", playground.Namespace, "Deployment.Name", playground.Name+"-deployment")

		// Create it
		err = r.Create(ctx, deployment)
		if err != nil {
			return ctrl.Result{}, err
		}

	} else if err != nil {
		return ctrl.Result{}, err
	} else {
		//Deployment already exists
		if !equality.Semantic.DeepDerivative(currentdeployment.Spec.Template.Spec.Containers[0].Resources, deployment.Spec.Template.Spec.Containers[0].Resources) ||
			(currentdeployment.Spec.Template.Spec.Containers[0].Image != deployment.Spec.Template.Spec.Containers[0].Image) ||
			(currentdeployment.Spec.Template.Spec.Containers[0].ImagePullPolicy != deployment.Spec.Template.Spec.Containers[0].ImagePullPolicy) ||
			(len(currentdeployment.Spec.Template.Spec.Volumes) != len(deployment.Spec.Template.Spec.Volumes)) {
			// Current deployment is different from desired. There is no good way of comparing deployment
			// due to the fact that additional thing inserted during deployment. As a result I am checking parameters that are defined in CRs
			r.Log.Info("Updating existing Deployment",
				"Deployment.Namespace", playground.Namespace, "Deployment.Name", playground.Name+"-deployment")

			// Update deployment
			err = r.Update(ctx, deployment)
			if err != nil {
				return ctrl.Result{}, err
			}
		} else {

			//Update status with deployment status
			if len(deployment.Status.Conditions) > 0 {
				playground.Status.Condition = &deployment.Status.Conditions[0]
				err = r.Status().Update(context.Background(), &playground)
				if err != nil {
					return ctrl.Result{}, err
				}
			}
		}
	}

	// Calculate desired service
	service := &apiv1.Service{}
	service, err = r.newServiceForPlayground(&playground, &labels)
	if err != nil {
		return ctrl.Result{}, err
	}

	//check if the service already exists
	currentservice := &apiv1.Service{}
	err = r.Get(ctx, types.NamespacedName{Name: playground.Name + "-service", Namespace: playground.Namespace}, currentservice)

	//If the service does not exist, create it
	if err != nil && errors.IsNotFound(err) {
		r.Log.Info("Creating a new Service", "Service.Namespace", playground.Namespace, "Service.Name", playground.Name+"-service")

		// Create service
		err = r.Create(ctx, service)
		if err != nil {
			return ctrl.Result{}, err
		}
	} else if err != nil {
		return ctrl.Result{}, err
	} else {
		// Service already exists
		if !equality.Semantic.DeepDerivative(currentservice.Spec.Type, service.Spec.Type) {
			r.Log.Info("Updating Service", "Service.Namespace", playground.Namespace, "Service.Name", playground.Name+"-service")
			// You can't update a service to change service type. So we need to remove existing one and recreate
			err = r.Delete(ctx, currentservice)
			if err != nil {
				return ctrl.Result{}, err
			}
			err = r.Create(ctx, service)
			if err != nil {
				return ctrl.Result{}, err
			}
		}
	}

	// If it is OpenShift, create route
	if r.IsOpenShift {
		//Check if Route exists
		route := &routev1.Route{}
		err = r.Get(ctx, types.NamespacedName{Name: playground.Name + "-route", Namespace: playground.Namespace}, route)
		if err != nil && errors.IsNotFound(err) {
			r.Log.Info("Creating new route", "Route.Namespace", playground.Namespace, "Route.name", playground.Name+"-route")
			// Create route resource
			route, err = r.newRouteForPLayground(&playground, &labels)
			if err != nil {
				return ctrl.Result{}, err
			}
			// Create route
			err = r.Create(ctx, route)
			if err != nil {
				return ctrl.Result{}, err
			}
			//r.Log.Info("Route created", "Route.Namespace", playground.Namespace, "Route.Name", playground.Name+"-route")

			// Route created successfully
			//return ctrl.Result{}, nil
		} else if err != nil {
			return ctrl.Result{}, err
		} /*else {

			// Route already exists - don't requeue
			r.Log.Info("Route already exists", "Route.Namespace", playground.Namespace, "Route.Name", playground.Name+"-route")
		}*/
	}

	return ctrl.Result{}, nil
}

// returns a Deployment that will manage a Playground pod
func (r *QiskitPlaygroundReconciler) newDeploymentForPlayground(pg *qiskitv1alpha1.QiskitPlayground, labels *map[string]string) (*appsv1.Deployment, error) {

	// Resources
	var resources *apiv1.ResourceRequirements
	if pg.Spec.Resources != nil {
		resources = pg.Spec.Resources
	} else {
		resources = &apiv1.ResourceRequirements{
			Limits: apiv1.ResourceList{
				"memory": resource.MustParse("2G"),
				"cpu":    resource.MustParse("2"),
			},
			Requests: apiv1.ResourceList{
				"memory": resource.MustParse("1G"),
				"cpu":    resource.MustParse("1"),
			},
		}
	}

	// Always 1 replica
	replicas := int32(1)

	// Build deployment
	deployment := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      pg.Name + "-deployment",
			Namespace: pg.Namespace,
			Labels:    *labels,
		},

		Spec: appsv1.DeploymentSpec{
			Replicas: &replicas,
			Selector: &metav1.LabelSelector{
				MatchLabels: *labels,
			},
			Template: v1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Name:   pg.Name,
					Labels: *labels,
				},
				Spec: v1.PodSpec{
					Containers: []v1.Container{
						{
							Name:            "qiskit",
							Image:           pg.Spec.Image,
							ImagePullPolicy: pg.Spec.ImagePullPolicy,
							Args:            []string{"start-notebook.sh", "--NotebookApp.token=''", "--NotebookApp.password=''"},
							Ports: []v1.ContainerPort{
								{
									Name:          "notebook-port",
									ContainerPort: 8888,
									Protocol:      v1.ProtocolTCP,
								},
							},
							Resources: *resources,
							Env: []v1.EnvVar{
								{
									Name:  "JUPYTER_ENABLE_LAB",
									Value: "yes",
								},
							},
							TerminationMessagePath:   "/dev/termination-log",
							TerminationMessagePolicy: "File",
							VolumeMounts:             []v1.VolumeMount{},
						},
					},
					ServiceAccountName: "default",
					Volumes:            []v1.Volume{},
				},
			},
		},
	}

	// Check if we have persistemce
	if pg.Spec.PVC != "" {
		// Add volume monts
		deployment.Spec.Template.Spec.Containers[0].VolumeMounts = append(deployment.Spec.Template.Spec.Containers[0].VolumeMounts,
			v1.VolumeMount{
				Name:      "notebook-persistence",
				MountPath: "/home/jovyan",
			})
		// Add volumes
		deployment.Spec.Template.Spec.Volumes = append(deployment.Spec.Template.Spec.Volumes,
			v1.Volume{
				Name: "notebook-persistence",
				VolumeSource: v1.VolumeSource{
					PersistentVolumeClaim: &v1.PersistentVolumeClaimVolumeSource{
						ClaimName: pg.Spec.PVC,
					},
				},
			})
	}

	// SetControllerReference sets owner as a Controller OwnerReference on owned.
	// This is used for garbage collection of the owned object and for
	// reconciling the owner object on changes to owned (with a Watch + EnqueueRequestForOwner).
	// Since only one OwnerReference can be a controller, it returns an error if
	// there is another OwnerReference with Controller flag set.
	if err := controllerutil.SetControllerReference(pg, deployment, r.Scheme); err != nil {
		return nil, err
	}
	return deployment, nil
}

// returns a Service that will expose playground functionality
func (r *QiskitPlaygroundReconciler) newServiceForPlayground(pg *qiskitv1alpha1.QiskitPlayground, labels *map[string]string) (*apiv1.Service, error) {

	service := &apiv1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      pg.Name + "-service",
			Namespace: pg.Namespace,
			Labels:    *labels,
		},
		Spec: apiv1.ServiceSpec{
			Ports: []apiv1.ServicePort{
				{
					Name:       "http",
					Port:       80,
					TargetPort: intstr.FromInt(8888),
				},
			},
			Selector: *labels,
			Type:     apiv1.ServiceTypeClusterIP,
		},
	}

	if pg.Spec.LoadBalancer && !r.IsOpenShift {
		// Our environment supports LoadBalancer service type and this is not OpenShift
		service.Spec.Type = apiv1.ServiceTypeLoadBalancer
	}

	// SetControllerReference sets owner as a Controller OwnerReference on owned.
	// This is used for garbage collection of the owned object and for
	// reconciling the owner object on changes to owned (with a Watch + EnqueueRequestForOwner).
	// Since only one OwnerReference can be a controller, it returns an error if
	// there is another OwnerReference with Controller flag set.
	if err := controllerutil.SetControllerReference(pg, service, r.Scheme); err != nil {
		return nil, err
	}
	return service, nil
}

// newRouteForCR returns a service pod with the same name/namespace as the cr
func (r *QiskitPlaygroundReconciler) newRouteForPLayground(pg *qiskitv1alpha1.QiskitPlayground, labels *map[string]string) (*routev1.Route, error) {
	// Build route
	weight := int32(100)
	rte := &routev1.Route{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Route",
			APIVersion: "route.openshift.io/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      pg.Name + "-route",
			Namespace: pg.Namespace,
			Labels:    *labels,
		},
		Spec: routev1.RouteSpec{
			To: routev1.RouteTargetReference{
				Kind:   "Service",
				Name:   pg.Name + "-service",
				Weight: &weight,
			},
			Port: &routev1.RoutePort{
				TargetPort: intstr.FromString("http"),
			},
			WildcardPolicy: "None",
		},
	}
	// SetControllerReference sets owner as a Controller OwnerReference on owned.
	// This is used for garbage collection of the owned object and for
	// reconciling the owner object on changes to owned (with a Watch + EnqueueRequestForOwner).
	// Since only one OwnerReference can be a controller, it returns an error if
	// there is another OwnerReference with Controller flag set.
	if err := controllerutil.SetControllerReference(pg, rte, r.Scheme); err != nil {
		return nil, err
	}
	return rte, nil

}

// SetupWithManager sets up the controller with the Manager.
func (r *QiskitPlaygroundReconciler) SetupWithManager(mgr ctrl.Manager) error {
	if r.IsOpenShift {
		r.Log.Info("Running on OpenShift cluster")
	} else {
		r.Log.Info("Running on Vanila Kubernetes cluster")
	}
	return ctrl.NewControllerManagedBy(mgr).
		For(&qiskitv1alpha1.QiskitPlayground{}).
		Owns(&appsv1.Deployment{}).
		Complete(r)
}
