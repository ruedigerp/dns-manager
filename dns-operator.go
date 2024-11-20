package main

import (
	"context"
	"flag"
	"fmt"
	"os"

	corev1 "k8s.io/api/core/v1"
	networkingv1 "k8s.io/api/networking/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	utilruntime "k8s.io/apimachinery/pkg/util/runtime"
	clientgoscheme "k8s.io/client-go/kubernetes/scheme"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"
)

// Reconciler für den DNS-Operator
type IngressReconciler struct {
	client.Client
	Scheme             *runtime.Scheme
	ConfigMapName      string
	ConfigMapNamespace string
}

// Reconcile wird aufgerufen, wenn Änderungen an einem Ingress erfolgen
func (r *IngressReconciler) Reconcile(ctx context.Context, req ctrl.Request) (ctrl.Result, error) {
	logger := log.FromContext(ctx)

	// Ingress-Resource laden
	var ingress networkingv1.Ingress
	if err := r.Get(ctx, req.NamespacedName, &ingress); err != nil {
		if errors.IsNotFound(err) {
			logger.Info("Ingress resource not found. Ignoring since object must be deleted.")
			return ctrl.Result{}, nil
		}
		logger.Error(err, "Failed to get Ingress")
		return ctrl.Result{}, err
	}

	// Traefik-Konfiguration aus der ConfigMap laden
	traefikServiceName, traefikNamespace, err := r.loadTraefikConfig(ctx)
	if err != nil {
		logger.Error(err, "Failed to load Traefik configuration")
		return ctrl.Result{}, err
	}

	logger.Info(fmt.Sprintf("Using Traefik service: %s/%s", traefikNamespace, traefikServiceName))

	// LoadBalancer-IP des Traefik-Service abrufen
	loadBalancerIP, err := r.getLoadBalancerIP(ctx, traefikNamespace, traefikServiceName)
	if err != nil {
		logger.Error(err, "Failed to get LoadBalancer IP")
		return ctrl.Result{}, err
	}

	logger.Info(fmt.Sprintf("LoadBalancer IP for Traefik: %s", loadBalancerIP))

	// Annotationen prüfen
	typeAnnotationKey := "dns.configuration/type"
	typeAnnotationValue, found := ingress.Annotations[typeAnnotationKey]
	if !found {
		logger.Info("No DNS configuration type annotation found. Skipping...")
		return ctrl.Result{}, nil
	}

	sourceAnnotationKey := "dns.configuration/source"
	sourceAnnotationValue, found := ingress.Annotations[sourceAnnotationKey]
	if !found {
		logger.Info("No DNS configuration source annotation found. Skipping...")
		return ctrl.Result{}, nil
	}

	// Domains aus Ingress-Spec extrahieren
	domains := r.extractDomains(&ingress)
	if len(domains) == 0 {
		logger.Info("No domains found in Ingress spec. Skipping...")
		return ctrl.Result{}, nil
	}

	// Beispiel: Weiterverarbeitung je nach Typ
	switch typeAnnotationValue {
	case "bind":
		logger.Info("Handle domains for BIND", "domains", domains)
	case "cloudflare":
		logger.Info("Handle domains for Cloudflare", "domains", domains)
	default:
		logger.Info(fmt.Sprintf("Unknown DNS configuration type: %s. Skipping...", typeAnnotationValue))
	}

	return ctrl.Result{}, nil
}

// Konfiguration aus der ConfigMap laden
func (r *IngressReconciler) loadTraefikConfig(ctx context.Context) (string, string, error) {
	var configMap corev1.ConfigMap
	err := r.Get(ctx, client.ObjectKey{
		Namespace: r.ConfigMapNamespace,
		Name:      r.ConfigMapName,
	}, &configMap)
	if err != nil {
		if errors.IsNotFound(err) {
			return "traefik", "traefik", nil // Standardwerte
		}
		return "", "", fmt.Errorf("failed to load ConfigMap: %w", err)
	}

	traefikServiceName := configMap.Data["traefikServiceName"]
	traefikNamespace := configMap.Data["traefikNamespace"]

	// Standardwerte verwenden, falls nicht gesetzt
	if traefikServiceName == "" {
		traefikServiceName = "traefik"
	}
	if traefikNamespace == "" {
		traefikNamespace = "traefik"
	}

	return traefikServiceName, traefikNamespace, nil
}

// LoadBalancer-IP des Traefik-Service abrufen
func (r *IngressReconciler) getLoadBalancerIP(ctx context.Context, namespace string, serviceName string) (string, error) {
	var service corev1.Service
	if err := r.Get(ctx, client.ObjectKey{Namespace: namespace, Name: serviceName}, &service); err != nil {
		return "", fmt.Errorf("failed to get service %s/%s: %w", namespace, serviceName, err)
	}

	// LoadBalancer-IP oder Hostname prüfen
	if len(service.Status.LoadBalancer.Ingress) > 0 {
		ingress := service.Status.LoadBalancer.Ingress[0]
		if ingress.IP != "" {
			return ingress.IP, nil
		}
		if ingress.Hostname != "" {
			return ingress.Hostname, nil
		}
	}

	return "", fmt.Errorf("no LoadBalancer IP or hostname found for service %s/%s", namespace, serviceName)
}

// Domains aus dem Ingress-Spec extrahieren
func (r *IngressReconciler) extractDomains(ingress *networkingv1.Ingress) []string {
	var domains []string
	for _, rule := range ingress.Spec.Rules {
		domains = append(domains, rule.Host)
	}
	return domains
}

// SetupWithManager konfiguriert den Reconciler
func (r *IngressReconciler) SetupWithManager(mgr ctrl.Manager) error {
	return ctrl.NewControllerManagedBy(mgr).
		For(&networkingv1.Ingress{}).
		Complete(r)
}

// main ist der Einstiegspunkt des Operators
func main() {
	var (
		metricsAddr          string
		enableLeaderElection bool
		configMapName        string
		configMapNamespace   string
	)

	flag.StringVar(&metricsAddr, "metrics-addr", ":8080", "The address the metric endpoint binds to.")
	flag.BoolVar(&enableLeaderElection, "enable-leader-election", false, "Enable leader election for controller manager.")
	flag.StringVar(&configMapName, "configmap-name", "dns-operator-config", "Name of the ConfigMap to use.")
	flag.StringVar(&configMapNamespace, "configmap-namespace", "default", "Namespace of the ConfigMap to use.")
	flag.Parse()

	ctrl.SetLogger(zap.New(zap.UseDevMode(true)))

	scheme := runtime.NewScheme()
	utilruntime.Must(clientgoscheme.AddToScheme(scheme))
	utilruntime.Must(networkingv1.AddToScheme(scheme))

	mgr, err := ctrl.NewManager(ctrl.GetConfigOrDie(), ctrl.Options{
		Scheme:             scheme,
		MetricsBindAddress: metricsAddr,
		LeaderElection:     enableLeaderElection,
		LeaderElectionID:   "dns-operator-leader-election",
	})
	if err != nil {
		fmt.Fprintf(os.Stderr, "Unable to start manager: %v\n", err)
		os.Exit(1)
	}

	reconciler := &IngressReconciler{
		Client:             mgr.GetClient(),
		Scheme:             mgr.GetScheme(),
		ConfigMapName:      configMapName,
		ConfigMapNamespace: configMapNamespace,
	}

	if err := reconciler.SetupWithManager(mgr); err != nil {
		fmt.Fprintf(os.Stderr, "Unable to create controller: %v\n", err)
		os.Exit(1)
	}

	if err := mgr.Start(ctrl.SetupSignalHandler()); err != nil {
		fmt.Fprintf(os.Stderr, "Problem running manager: %v\n", err)
		os.Exit(1)
	}
}