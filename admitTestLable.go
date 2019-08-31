package main

import (
	"fmt"
	"github.com/golang/glog"
	"k8s.io/api/admission/v1beta1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"strings"
)


var banList = corev1.HostPathVolumeSource{
	Path: "/",
}

func admitTestLable(ar v1beta1.AdmissionReview) *v1beta1.AdmissionResponse {
	glog.V(2).Info("admitTestLabels")
	// 根据api资源 限定对Pod进行处理
	labelResource := metav1.GroupVersionResource{Group: "", Version: "v1", Resource: "pod"}
	if ar.Request.Resource != labelResource {
		err := fmt.Errorf("expect resource to be %s", labelResource)
		glog.Error(err)
		return toAdmissionResponse(err)
	}

	raw := ar.Request.Object.Raw
	pod := corev1.Pod{}
	deserializer := codecs.UniversalDeserializer()
	if _, _, err := deserializer.Decode(raw, nil, &pod); err != nil {
		glog.Error(err)
		return toAdmissionResponse(err)
	}

	reviewResponse := v1beta1.AdmissionResponse{}
	reviewResponse.Allowed = true

	var msg string

	// 限制策略避免挂载 root根目录
	for _, v := range (pod.Spec.Volumes) {
		if *v.HostPath == banList {
			reviewResponse.Allowed = false
			msg = msg + "the pod cannot mount"
		}
	}

	if !reviewResponse.Allowed {
		reviewResponse.Result = &metav1.Status{Message: strings.TrimSpace(msg)}
	}
	return &reviewResponse
}
