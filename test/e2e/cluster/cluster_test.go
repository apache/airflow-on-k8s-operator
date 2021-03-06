// Licensed to the Apache Software Foundation (ASF) under one or more
// contributor license agreements.  See the NOTICE file distributed with
// this work for additional information regarding copyright ownership.
// The ASF licenses this file to You under the Apache License, Version 2.0
// (the "License"); you may not use this file except in compliance with
// the License.  You may obtain a copy of the License at
//
//    http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package e2etest

import (
	"testing"

	"github.com/apache/airflow-on-k8s-operator/api/v1alpha1"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
	_ "k8s.io/client-go/plugin/pkg/client/auth/gcp"
	"sigs.k8s.io/controller-reconciler/pkg/test"
)

const (
	CRName    = "AirflowCluster"
	SampleDir = "../../../hack/sample/"
)

var f *test.Framework
var ctx, basectx *test.Context
var deleteBase bool
var deleteCluster = true

func airflowBase(file string) *v1alpha1.AirflowBase {
	cr := &v1alpha1.AirflowBase{}
	if err := f.LoadFromFile(file, cr); err != nil {
		return nil
	}
	return cr
}
func airflowCluster(file string) *v1alpha1.AirflowCluster {
	cr := &v1alpha1.AirflowCluster{}
	if err := f.LoadFromFile(file, cr); err != nil {
		return nil
	}
	return cr
}

func Test(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, CRName+" Suite")
}

var _ = BeforeSuite(func() {
	f = test.New(CRName)
	err := v1alpha1.SchemeBuilder.AddToScheme(f.GetScheme())
	Expect(err).NotTo(HaveOccurred(), "failed to initialize the Framework: %v", err)
	f.Start()
})

var _ = AfterSuite(func() {
	if ctx != nil {
		ctx.DeleteCR()
	}
	if basectx != nil {
		basectx.DeleteCR()
	}
	if f != nil {
		f.Stop()
	}
})

func isBaseReady(cr interface{}) bool {
	stts := cr.(*v1alpha1.AirflowBase).Status
	return stts.IsReady()
}

func isClusterReady(cr interface{}) bool {
	stts := cr.(*v1alpha1.AirflowCluster).Status
	return stts.IsReady()
}

func checkLocal(cr *v1alpha1.AirflowCluster) {
	ctx.WithTimeout(200).CheckStatefulSet(cr.Name+"-airflowui", 1, 1)
	ctx.WithTimeout(200).CheckStatefulSet(cr.Name+"-scheduler", 1, 1)
	ctx.WithTimeout(200).CheckCR(isClusterReady)
}
func checkCelery(cr *v1alpha1.AirflowCluster) {
	ctx.WithTimeout(200).CheckStatefulSet(cr.Name+"-airflowui", 1, 1)
	ctx.WithTimeout(200).CheckStatefulSet(cr.Name+"-flower", 1, 1)
	ctx.WithTimeout(200).CheckStatefulSet(cr.Name+"-redis", 1, 1)
	ctx.WithTimeout(200).CheckStatefulSet(cr.Name+"-scheduler", 1, 1)
	ctx.WithTimeout(200).CheckStatefulSet(cr.Name+"-worker", 2, 1)
	ctx.WithTimeout(10).CheckService(cr.Name+"-redis", map[string]int32{"redis": 6379})
	ctx.WithTimeout(200).CheckCR(isClusterReady)
}

var _ = Describe(CRName+" controller tests", func() {
	AfterEach(func() {
		if deleteCluster {
			ctx.DeleteCR()
			ctx = nil
		}
		deleteCluster = true
		if deleteBase {
			deleteBase = false
			basectx.DeleteCR()
			basectx = nil
		}
	})

	// Postgres
	It("creating a "+CRName+" with postgres, celery executor", func() {
		basectx = f.NewContext().WithCR(airflowBase(SampleDir + "postgres-celery/base.yaml"))
		ctx = f.NewContext().WithCR(airflowCluster(SampleDir + "postgres-celery/cluster.yaml"))
		basecr := basectx.CR.(*v1alpha1.AirflowBase)
		cr := ctx.CR.(*v1alpha1.AirflowCluster)
		By("creating a base " + basecr.Name)
		basectx.CreateCR()
		basectx.WithTimeout(200).CheckCR(isBaseReady)

		By("creating a new " + CRName + ": " + cr.Name)
		ctx.CreateCR()
		checkCelery(cr)
		deleteCluster = false
	})

	It("scaling up workers for "+CRName+" with postgres, celery executor", func() {
		ctx.RefreshCR()
		cr := ctx.CR.(*v1alpha1.AirflowCluster)
		By("scaling up workers: " + cr.Name)
		cr.Spec.Worker.Replicas++
		ctx.UpdateCR()
		ctx.WithTimeout(200).CheckStatefulSet(cr.Name+"-worker", 3, 2)
	})

	It("creating a "+CRName+" with postgres, local executor", func() {
		ctx = f.NewContext().WithCR(airflowCluster(SampleDir + "postgres-local/cluster.yaml"))
		cr := ctx.CR.(*v1alpha1.AirflowCluster)
		By("creating a new " + CRName + ": " + cr.Name)
		ctx.CreateCR()
		checkLocal(cr)
	})

	It("creating a "+CRName+" with postgres, k8s executor", func() {
		ctx = f.NewContext().WithCR(airflowCluster(SampleDir + "postgres-k8s/cluster.yaml"))
		cr := ctx.CR.(*v1alpha1.AirflowCluster)
		By("creating a new " + CRName + ": " + cr.Name)
		ctx.CreateCR()
		checkLocal(cr)
		deleteBase = true
	})

	// Mysql
	It("creating a "+CRName+" with mysql, celery executor", func() {
		basectx = f.NewContext().WithCR(airflowBase(SampleDir + "mysql-celery/base.yaml"))
		ctx = f.NewContext().WithCR(airflowCluster(SampleDir + "mysql-celery/cluster.yaml"))
		basecr := basectx.CR.(*v1alpha1.AirflowBase)
		cr := ctx.CR.(*v1alpha1.AirflowCluster)
		By("creating a base " + basecr.Name)
		basectx.CreateCR()
		basectx.WithTimeout(200).CheckCR(isBaseReady)

		By("creating a new " + CRName + ": " + cr.Name)
		ctx.CreateCR()
		checkCelery(cr)
	})

	It("creating a "+CRName+" with mysql, local executor", func() {
		ctx = f.NewContext().WithCR(airflowCluster(SampleDir + "mysql-local/cluster.yaml"))
		cr := ctx.CR.(*v1alpha1.AirflowCluster)
		By("creating a new " + CRName + ": " + cr.Name)
		ctx.CreateCR()
		checkLocal(cr)
	})

	It("creating a "+CRName+" with mysql, k8s executor", func() {
		ctx = f.NewContext().WithCR(airflowCluster(SampleDir + "mysql-k8s/cluster.yaml"))
		cr := ctx.CR.(*v1alpha1.AirflowCluster)
		By("creating a new " + CRName + ": " + cr.Name)
		ctx.CreateCR()
		checkLocal(cr)
		deleteBase = true
	})

})
