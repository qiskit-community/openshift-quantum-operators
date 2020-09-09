package qiskitplayground

import (
	dobtechv1 "github.com/jdob/qiskit-operator/pkg/apis/dobtech/v1"

	appsv1 "k8s.io/api/apps/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	routev1 "github.com/openshift/api/route/v1"
	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/controller-runtime/pkg/controller/controllerutil"
)

const port = 8888
const image = "singhp11/centos-jupyter"

func deploymentName(q *dobtechv1.QiskitPlayground) string {
	return q.Name
}

func serviceName(q *dobtechv1.QiskitPlayground) string {
	return q.Name
}

func routeName(q *dobtechv1.QiskitPlayground) string {
	return q.Name
}

func (r *ReconcileQiskitPlayground) deployment(q *dobtechv1.QiskitPlayground) *appsv1.Deployment {
	labels := labels(q)
	size := int32(1)

	dep := &appsv1.Deployment{
		ObjectMeta: metav1.ObjectMeta{
			Name:		deploymentName(q),
			Namespace: 	q.Namespace,
			Labels:		labels,
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
						Image:	image,
						Name:	"playground",
						Ports:	[]corev1.ContainerPort{{
							ContainerPort: 	port,
							Name:			"playground",
						}},
						VolumeMounts: []corev1.VolumeMount{{
							Name: "qiskit-secret",
							MountPath: "/tmp/secrets/qiskitsecret",
							ReadOnly: true,
	
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

func (r *ReconcileQiskitPlayground) service(q *dobtechv1.QiskitPlayground) *corev1.Service {
	labels := labels(q)

	s := &corev1.Service{
		ObjectMeta: metav1.ObjectMeta{
			Name:		serviceName(q),
			Namespace: 	q.Namespace,
			Labels:		labels,
		},
		Spec: corev1.ServiceSpec{
			Selector: labels,
			Ports: []corev1.ServicePort{{
				Protocol: corev1.ProtocolTCP,
				Port: port,
				TargetPort: intstr.FromInt(port),
			}},
			Type: corev1.ServiceTypeNodePort,
		},
	}

	controllerutil.SetControllerReference(q, s, r.scheme)
	return s
}

func (r *ReconcileQiskitPlayground) route(q *dobtechv1.QiskitPlayground) *routev1.Route {
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