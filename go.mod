module github.com/apache/airflow-on-k8s-operator

go 1.13

require (
	github.com/go-logr/logr v0.1.0
	github.com/kr/pretty v0.2.0 // indirect
	github.com/kubernetes-sigs/application v0.8.1
	github.com/onsi/ginkgo v1.11.0
	github.com/onsi/gomega v1.8.1
	k8s.io/api v0.17.2
	k8s.io/apimachinery v0.17.2
	k8s.io/client-go v0.17.1
	k8s.io/utils v0.0.0-20200117235808-5f6fbceb4c31 // indirect
	sigs.k8s.io/controller-reconciler v0.0.0-00010101000000-000000000000
	sigs.k8s.io/controller-runtime v0.4.0
)

replace sigs.k8s.io/controller-reconciler => ./vendor/sigs.k8s.io/controller-reconciler
