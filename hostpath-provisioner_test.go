package main

import (
	"github.com/kubernetes-incubator/external-storage/lib/controller"
	"gopkg.in/check.v1"
	"k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"testing"
)

var _ = check.Suite(&S{})

type S struct {
	h hostPathProvisioner
}

func (s *S) SetUpSuite(c *check.C) {
}

func (s *S) TestParseParameters(c *check.C) {
	_, err := s.h.parseParameters(map[string]string{
		"pvdir": "/tmp/pvLocalDir",
	})
	c.Assert(err, check.IsNil)
	_, err = s.h.parseParameters(map[string]string{
		"otherParam": "otherParam",
	})
	c.Assert(err, check.NotNil)
}

func (s *S) TestProvision(c *check.C) {
	expectedPv := &v1.PersistentVolume{
		ObjectMeta: metav1.ObjectMeta{
			Name: "test",
			Annotations: map[string]string{
				provisionerIDAnn: "",
			},
		},
		Spec: v1.PersistentVolumeSpec{
			PersistentVolumeReclaimPolicy: v1.PersistentVolumeReclaimDelete,
			Capacity: v1.ResourceList{
				v1.ResourceName(v1.ResourceStorage): resource.MustParse("10G"),
			},
			PersistentVolumeSource: v1.PersistentVolumeSource{
				HostPath: &v1.HostPathVolumeSource{
					Path: "/tmp/pvLocalDir/test",
				},
			},
		},
	}
	vo := controller.VolumeOptions{
		Parameters: map[string]string{
			"pvdir": "/tmp/pvLocalDir",
		},
		PersistentVolumeReclaimPolicy: v1.PersistentVolumeReclaimDelete,
		PVName: "test",
		PVC: &v1.PersistentVolumeClaim{
			Spec: v1.PersistentVolumeClaimSpec{
				Resources: v1.ResourceRequirements{
					Requests: v1.ResourceList{
						v1.ResourceName(v1.ResourceStorage): resource.MustParse("10G"),
					},
				},
			},
		},
	}
	pv, err := s.h.Provision(vo)
	c.Assert(err, check.IsNil)
	c.Assert(pv, check.DeepEquals, expectedPv)
}

func Test(t *testing.T) {
	check.TestingT(t)
}
