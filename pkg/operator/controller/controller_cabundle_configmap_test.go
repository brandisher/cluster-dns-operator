package controller

import (
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"
	operatorv1 "github.com/openshift/api/operator/v1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func TestDesiredCABundleConfigmap(t *testing.T) {
	sourceConfigmap := corev1.ConfigMap{
		ObjectMeta: metav1.ObjectMeta{
			Name:      "cabundle-config",
			Namespace: GlobalUserSpecifiedConfigNamespace,
		},
		Data: map[string]string{"caBundle": "test-bundle"},
	}

	destName := CABundleConfigMapName(sourceConfigmap.Name)

	dns := &operatorv1.DNS{
		ObjectMeta: metav1.ObjectMeta{
			Name: DefaultDNSController,
		},
		Spec: operatorv1.DNSSpec{},
	}

	desired, cm, err := desiredCABundleConfigMap(dns, true, &sourceConfigmap, destName)
	if err != nil || desired == false {
		t.Errorf("unexpected error : %v", err)
	} else if diff := cmp.Diff(cm.Data, sourceConfigmap.Data); diff != "" {
		t.Errorf("unexpected CA Bundle ConfigMap data;\n%s", diff)
	} else if diff := cmp.Diff(cm.OwnerReferences, []metav1.OwnerReference{dnsOwnerRef(dns)}); diff != "" {
		t.Errorf("unexpected CA Bundle ConfigMap OwnerReference;\n%s", diff)
	}

	desired, cm, err = desiredCABundleConfigMap(dns, false, &sourceConfigmap, destName)
	if desired != false || cm != nil || err != nil {
		t.Errorf("expected return values of false, nil, nil when haveSource is false: %v", err)
	}

	var delTimestamp *metav1.Time
	delTimestamp = &metav1.Time{
		Time: time.Now(),
	}
	dns.DeletionTimestamp = delTimestamp

	desired, cm, err = desiredCABundleConfigMap(dns, true, &sourceConfigmap, destName)
	if desired != false || cm != nil || err != nil {
		t.Errorf("expected return values of false, nil, nil when dns.DeletionTimestamp is not nil: %v", err)
	}
}
