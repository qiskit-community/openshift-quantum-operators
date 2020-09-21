package qiskitplayground

import (
	"context"

	singhp11v1 "github.com/example-inc/new--/openshift-qiskit-operator/pkg/apis/singhp11/v1"
	routev1 "github.com/openshift/api/route/v1"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/controller"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
	"sigs.k8s.io/controller-runtime/pkg/handler"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/reconcile"
	"sigs.k8s.io/controller-runtime/pkg/source"
)

var log = logf.Log.WithName("controller_qiskitplayground")

const port = 8888
const image = "singhp11/centos-jupyter"

/**
* USER ACTION REQUIRED: This is a scaffold file intended for the user to modify with their own Controller
* business logic.  Delete these comments after modifying this file.*
 */

// Add creates a new QiskitPlayground Controller and adds it to the Manager. The Manager will set fields on the Controller
// and Start it when the Manager is Started.
func Add(mgr manager.Manager) error {
	return add(mgr, newReconciler(mgr))
}

// newReconciler returns a new reconcile.Reconciler
func newReconciler(mgr manager.Manager) reconcile.Reconciler {
	return &ReconcileQiskitPlayground{client: mgr.GetClient(), scheme: mgr.GetScheme()}
}

// add adds a new Controller to mgr with r as the reconcile.Reconciler
func add(mgr manager.Manager, r reconcile.Reconciler) error {
	// Create a new controller
	c, err := controller.New("qiskitplayground-controller", mgr, controller.Options{Reconciler: r})
	if err != nil {
		return err
	}

	// Watch for changes to primary resource QiskitPlayground
	err = c.Watch(&source.Kind{Type: &singhp11v1.QiskitPlayground{}}, &handler.EnqueueRequestForObject{})
	if err != nil {
		return err
	}

	// TODO(user): Modify this to be the types you create that are owned by the primary resource
	// Watch for changes to secondary resource Pods and requeue the owner QiskitPlayground
	err = c.Watch(&source.Kind{Type: &corev1.Pod{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &singhp11v1.QiskitPlayground{},
	})
	if err != nil {
		return err
	}

	return nil
}

// blank assignment to verify that ReconcileQiskitPlayground implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileQiskitPlayground{}

// ReconcileQiskitPlayground reconciles a QiskitPlayground object
type ReconcileQiskitPlayground struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

// Reconcile reads that state of the cluster for a QiskitPlayground object and makes changes based on the state read
// and what is in the QiskitPlayground.Spec
// TODO(user): Modify this Reconcile function to implement your Controller logic.  This example creates
// a Pod as an example
// Note:
// The Controller will requeue the Request to be processed again if the returned error is non-nil or
// Result.Requeue is true, otherwise upon completion it will remove the work from the queue.
func (r *ReconcileQiskitPlayground) Reconcile(request reconcile.Request) (reconcile.Result, error) {
	reqLogger := log.WithValues("Request.Namespace", request.Namespace, "Request.Name", request.Name)
	reqLogger.Info("Reconciling QiskitPlayground")

	// Fetch the QiskitPlayground instance
	instance := &singhp11v1.QiskitPlayground{}
	err := r.client.Get(context.TODO(), request.NamespacedName, instance)
	if err != nil {
		if errors.IsNotFound(err) {
			// Request object not found, could have been deleted after reconcile request.
			// Owned objects are automatically garbage collected. For additional cleanup logic use finalizers.
			// Return and don't requeue
			return reconcile.Result{}, nil
		}
		// Error reading the object - requeue the request.
		return reconcile.Result{}, err
	}

	// Define a new Pod object
	pod := newPodForCR(instance)

	// Set QiskitPlayground instance as the owner and controller
	if err := controllerutil.SetControllerReference(instance, pod, r.scheme); err != nil {
		return reconcile.Result{}, err
	}

	// Check if this Pod already exists
	found := &corev1.Pod{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: pod.Name, Namespace: pod.Namespace}, found)
	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating a new Pod", "Pod.Namespace", pod.Namespace, "Pod.Name", pod.Name)
		err = r.client.Create(context.TODO(), pod)
		if err != nil {
			return reconcile.Result{}, err
		}

		// Pod created successfully - don't requeue
		return reconcile.Result{}, nil
	} else if err != nil {
		return reconcile.Result{}, err
	}

	// Pod already exists - don't requeue
	reqLogger.Info("Skip reconcile: Pod already exists", "Pod.Namespace", found.Namespace, "Pod.Name", found.Name)

	// Define new Service object
	service := newServiceForCR(instance)
	if err := controllerutil.SetControllerReference(instance, service, r.scheme); err != nil {
		return reconcile.Result{}, err
	}

	// Check if service exists
	foundService := &corev1.Service{}

	err = r.client.Get(context.TODO(), types.NamespacedName{Name: service.Name, Namespace: service.Namespace}, foundService)
	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating new service", "Service.Namespace", service.Namespace, "Service.name", service.Name)
		err = r.client.Create(context.TODO(), service)
		if err != nil {
			return reconcile.Result{}, err
		}
	} else if err != nil {
		return reconcile.Result{}, err
	}

	// Service already exists - don't requeue
	reqLogger.Info("Skip reconcile: Service already exists", "Service.Namespace", foundService.Namespace, "Service.Name", foundService.Name)

	// Define new Route object
	route := newRouteForCR(instance)
	if err := controllerutil.SetControllerReference(instance, route, r.scheme); err != nil {
		return reconcile.Result{}, err
	}

	//Check if Route exists
	foundRoute := &routev1.Route{}
	err = r.client.Get(context.TODO(), types.NamespacedName{Name: foundRoute.Name, Namespace: foundRoute.Namespace}, foundRoute)
	if err != nil && errors.IsNotFound(err) {
		reqLogger.Info("Creating new route", "Route.Namespace", route.Namespace, "Route.name", route.Name)
		err = r.client.Create(context.TODO(), route)
		if err != nil {
			return reconcile.Result{}, err
		}
	} else if err != nil {
		return reconcile.Result{}, err
	}

	// Route already exists - don't requeue
	reqLogger.Info("Skip reconcile: Route already exists", "Route.Namespace", foundRoute.Namespace, "Route.Name", foundRoute.Name)

	return reconcile.Result{}, nil
}

// newPodForCR returns a busybox pod with the same name/namespace as the cr
func newPodForCR(cr *singhp11v1.QiskitPlayground) *corev1.Pod {
	labels := map[string]string{
		"app": cr.Name,
	}
	return &corev1.Pod{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name + "-pod",
			Namespace: cr.Namespace,
			Labels:    labels,
		},
		Spec: corev1.PodSpec{
			Containers: []corev1.Container{
				{
					Name:  "openshift-qiskit-operator",
					Image: "quay.io/husky_parul/centos-qiskit-jupyter",
					VolumeMounts: []corev1.VolumeMount{{
						Name:      "qiskit-secret",
						MountPath: "/tmp/secrets/qiskitsecret",
						ReadOnly:  true,
					}},
				},
			},
			Volumes: []corev1.Volume{{
				Name: "qiskit-secret",
				VolumeSource: corev1.VolumeSource{
					Secret: &corev1.SecretVolumeSource{
						SecretName: "qiskit-secret",
					},
				},
			}},
		},
	}
}

// newServiceForCR returns a service pod with the same name/namespace as the cr
func newServiceForCR(cr *singhp11v1.QiskitPlayground) *corev1.Service {
	labels := map[string]string{
		"app": cr.Name,
	}
	return &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name,
			Namespace: cr.Namespace,
			Labels:    labels,
		},
		Spec: corev1.ServiceSpec{
			Selector: labels,
			Ports: []corev1.ServicePort{{
				Protocol:   corev1.ProtocolTCP,
				Port:       port,
				TargetPort: intstr.FromInt(port),
			}},
			Type: corev1.ServiceTypeNodePort,
		},
	}
}

// newRouteForCR returns a service pod with the same name/namespace as the cr
func newRouteForCR(cr *singhp11v1.QiskitPlayground) *routev1.Route {
	labels := map[string]string{
		"app": cr.Name,
	}
	rte := &routev1.Route{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Route",
			APIVersion: "route.openshift.io/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      cr.Name,
			Namespace: cr.Namespace,
			Labels:    labels,
		},
		Spec: routev1.RouteSpec{
			To: routev1.RouteTargetReference{
				Kind: "Service",
				Name: cr.Name,
			},
			Port: &routev1.RoutePort{
				TargetPort: intstr.FromInt(port),
			},
		},
	}
	return rte
}
