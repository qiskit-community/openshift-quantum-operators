package qiskitplayground

import (
	"context"

	singhp11v1 "github.com/husky-parul/openshift-qiskit-operator/pkg/apis/singhp11/v1"
	routev1 "github.com/openshift/api/route/v1"
	appsv1 "k8s.io/api/apps/v1"
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

// blank assignment to verify that ReconcileQiskitPlayground implements reconcile.Reconciler
var _ reconcile.Reconciler = &ReconcileQiskitPlayground{}

// ReconcileQiskitPlayground reconciles a QiskitPlayground object
type ReconcileQiskitPlayground struct {
	// This client, initialized using mgr.Client() above, is a split client
	// that reads objects from the cache and writes to the apiserver
	client client.Client
	scheme *runtime.Scheme
}

const port = 8888
const image = "singhp11/centos-jupyter"

func deploymentName(q *singhp11v1.QiskitPlayground) string {
	return q.Name
}

func serviceName(q *singhp11v1.QiskitPlayground) string {
	return q.Name
}

func routeName(q *singhp11v1.QiskitPlayground) string {
	return q.Name
}

func (r *ReconcileQiskitPlayground) ensureDeployment(request reconcile.Request, instance *singhp11v1.QiskitPlayground, dep *appsv1.Deployment) (*reconcile.Result, error) {
	// See if deployment already exists and create if it doesn't
	found := &appsv1.Deployment{}
	err := r.client.Get(context.TODO(), types.NamespacedName{
		Name:      dep.Name,
		Namespace: instance.Namespace,
	}, found)
	if err != nil && errors.IsNotFound(err) {

		// Create the deployment
		log.Info("Creating a new Deployment", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
		err = r.client.Create(context.TODO(), dep)

		if err != nil {
			// Deployment failed
			log.Error(err, "Failed to create new Deployment", "Deployment.Namespace", dep.Namespace, "Deployment.Name", dep.Name)
			return &reconcile.Result{}, err
		} else {
			// Deployment was successful
			return nil, nil
		}
	} else if err != nil {
		// Error that isn't due to the deployment not existing
		log.Error(err, "Failed to get Deployment")
		return &reconcile.Result{}, err
	}

	return nil, nil
}

func (r *ReconcileQiskitPlayground) ensureService(request reconcile.Request, instance *singhp11v1.QiskitPlayground, s *corev1.Service) (*reconcile.Result, error) {
	found := &corev1.Service{}
	err := r.client.Get(context.TODO(), types.NamespacedName{
		Name:      s.Name,
		Namespace: instance.Namespace,
	}, found)
	if err != nil && errors.IsNotFound(err) {

		// Create the service
		log.Info("Creating a new Service", "Service.Namespace", s.Namespace, "Service.Name", s.Name)
		err = r.client.Create(context.TODO(), s)

		if err != nil {
			// Creation failed
			log.Error(err, "Failed to create new Service", "Service.Namespace", s.Namespace, "Service.Name", s.Name)
			return &reconcile.Result{}, err
		} else {
			// Creation was successful
			return nil, nil
		}
	} else if err != nil {
		// Error that isn't due to the service not existing
		log.Error(err, "Failed to get Service")
		return &reconcile.Result{}, err
	}

	return nil, nil
}

func (r *ReconcileQiskitPlayground) ensureRoute(request reconcile.Request, instance *singhp11v1.QiskitPlayground, rte *routev1.Route) (*reconcile.Result, error) {
	found := &routev1.Route{}
	err := r.client.Get(context.TODO(), types.NamespacedName{
		Name:      rte.Name,
		Namespace: instance.Namespace,
	}, found)
	if err != nil && errors.IsNotFound(err) {

		// Create the route
		log.Info("Creating a new Route", "Route.Namespace", rte.Namespace, "Route.Name", rte.Name)
		err = r.client.Create(context.TODO(), rte)

		if err != nil {
			// Creation failed
			log.Error(err, "Failed to create new Route")
			return &reconcile.Result{}, err
		} else {
			// Creation successful
			return nil, nil
		}
	} else if err != nil {
		// Error that isn't the route not existing
		log.Error(err, "Failed to create Route")
		return &reconcile.Result{}, err
	}

	return nil, nil
}

func labels(q *singhp11v1.QiskitPlayground) map[string]string {
	return map[string]string{
		"app":                 "playground",
		"qiskitplayground_cr": q.Name,
	}
}

func (r *ReconcileQiskitPlayground) deployment(q *singhp11v1.QiskitPlayground) *appsv1.Deployment {
	labels := labels(q)
	size := int32(1)

	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:      deploymentName(q),
			Namespace: q.Namespace,
			Labels:    labels,
		},
		Spec: appsv1.DeploymentSpec{
			Replicas: &size,
			Selector: &metav1.LabelSelector{
				MatchLabels: labels,
			},
			Template: corev1.PodTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: labels,
				},
				Spec: corev1.PodSpec{
					Containers: []corev1.Container{{
						Image: image,
						Name:  "playground",
						Ports: []corev1.ContainerPort{{
							ContainerPort: port,
							Name:          "playground",
						}},
						VolumeMounts: []corev1.VolumeMount{{
							Name:      "qiskit-secret",
							MountPath: "/tmp/secrets/qiskitsecret",
							ReadOnly:  true,
						}},
					}},
					Volumes: []corev1.Volume{{
						Name: "qiskit-secret",
						VolumeSource: corev1.VolumeSource{
							Secret: &corev1.SecretVolumeSource{
								SecretName: "qiskit-secret",
							},
						},
					}},
				},
			},
		},
	}

	controllerutil.SetControllerReference(q, dep, r.scheme)
	return dep
}

func (r *ReconcileQiskitPlayground) service(q *singhp11v1.QiskitPlayground) *corev1.Service {
	labels := labels(q)

	s := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:      serviceName(q),
			Namespace: q.Namespace,
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

	controllerutil.SetControllerReference(q, s, r.scheme)
	return s
}

func (r *ReconcileQiskitPlayground) route(q *singhp11v1.QiskitPlayground) *routev1.Route {
	labels := labels(q)

	rte := &routev1.Route{
		TypeMeta: metav1.TypeMeta{
			Kind:       "Route",
			APIVersion: "route.openshift.io/v1",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      routeName(q),
			Namespace: q.Namespace,
			Labels:    labels,
		},
		Spec: routev1.RouteSpec{
			To: routev1.RouteTargetReference{
				Kind: "Service",
				Name: serviceName(q),
			},
			Port: &routev1.RoutePort{
				TargetPort: intstr.FromInt(8888),
			},
		},
	}

	return rte
}

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

	// Watch for changes to secondary resource Pods and requeue the owner QiskitPlayground

	err = c.Watch(&source.Kind{Type: &appsv1.Deployment{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &singhp11v1.QiskitPlayground{},
	})
	if err != nil {
		return err
	}
	err = c.Watch(&source.Kind{Type: &corev1.Service{}}, &handler.EnqueueRequestForOwner{
		IsController: true,
		OwnerType:    &singhp11v1.QiskitPlayground{},
	})
	if err != nil {
		return err
	}

	return nil
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
	var result *reconcile.Result
	result, err = r.ensureDeployment(request, instance, r.deployment(instance))
	if err != nil {
		return *result, err
	}

	result, err = r.ensureService(request, instance, r.service(instance))
	if err != nil {
		return *result, err
	}

	result, err = r.ensureRoute(request, instance, r.route(instance))
	if err != nil {
		return *result, err
	}
	return reconcile.Result{}, nil
}
