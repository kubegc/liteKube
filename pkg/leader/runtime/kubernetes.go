package runtime

import (
	"context"
	"path/filepath"

	// link to github.com/Litekube/kine, we have make some addition

	"github.com/litekube/LiteKube/pkg/options/leader/apiserver"
	"github.com/litekube/LiteKube/pkg/options/leader/controllermanager"
	"github.com/litekube/LiteKube/pkg/options/leader/scheduler"
	rbacv1 "k8s.io/api/rbac/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog/v2"
	"sigs.k8s.io/yaml"
)

// kubectl create clusterrolebinding kubelet-bootstrap --clusterrole=system:node-bootstrapper --user=kubelet-bootstrap
var rolebindingBootstrapYAML = `apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: kubelet-bootstrap
subjects:
- kind: User
  name: kubelet-bootstrap
  apiGroup: rbac.authorization.k8s.io
roleRef:
  kind: ClusterRole
  name: system:node-bootstrapper
  apiGroup: rbac.authorization.k8s.io
`

var rolebindingAccessKubeletYAML = `apiVersion: rbac.authorization.k8s.io/v1
kind: ClusterRoleBinding
metadata:
  name: kube-apiserver:kubelet-apis
subjects:
- kind: User
  name: system:kube-apiserver
  apiGroup: rbac.authorization.k8s.io
roleRef:
  kind: ClusterRole
  name: system:kubelet-api-admin
  apiGroup: rbac.authorization.k8s.io
`

type KubernatesServer struct {
	ctx               context.Context
	logPath           string
	apiserver         *Apiserver
	controllerManager *ControllerManager
	scheduler         *Scheduler
	ApiserverOptions  *apiserver.ApiserverOptions
	ControllerOptions *controllermanager.ControllerManagerOptions
	SchedulerOptions  *scheduler.SchedulerOptions
	KubeAdminPath     string
}

func NewKubernatesServer(ctx context.Context, apiserverOptions *apiserver.ApiserverOptions, controllerOptions *controllermanager.ControllerManagerOptions, schedulerOptions *scheduler.SchedulerOptions, kubeAdminPath string, logPath string) *KubernatesServer {
	return &KubernatesServer{
		ctx:               ctx,
		logPath:           logPath,
		apiserver:         NewApiserver(ctx, apiserverOptions, filepath.Join(logPath, "kube-apiserver.log")),
		controllerManager: NewControllerManager(ctx, controllerOptions, filepath.Join(logPath, "kube-controller-manager.log")),
		scheduler:         NewScheduler(ctx, schedulerOptions, filepath.Join(logPath, "kube-scheduler.log")),
		ApiserverOptions:  apiserverOptions,
		ControllerOptions: controllerOptions,
		SchedulerOptions:  schedulerOptions,
		KubeAdminPath:     kubeAdminPath,
	}
}

// start run in routine and no wait
func (s *KubernatesServer) Run() error {
	klog.Info("start to run kubernates server")

	if err := s.apiserver.Run(); err != nil {
		return err
	}

	if err := s.controllerManager.Run(s.KubeAdminPath); err != nil {
		return err
	}

	if err := s.scheduler.Run(); err != nil {
		return err
	}

	if err := s.RunAfter(); err != nil {
		return err
	}

	return nil
}

func (s *KubernatesServer) RunAfter() error {
	// add kubelet-bootstrap role-binding
	kubeconfig, err := clientcmd.BuildConfigFromFlags("", s.KubeAdminPath)
	if err != nil {
		return err
	}

	k8sClient, err := kubernetes.NewForConfig(kubeconfig)
	if err != nil {
		return err
	}

	if _, err := k8sClient.RbacV1().ClusterRoleBindings().Get(s.ctx, "kubelet-bootstrap", metav1.GetOptions{}); err != nil {
		clusterRoleBindings := &rbacv1.ClusterRoleBinding{}
		if err := yaml.Unmarshal([]byte(rolebindingBootstrapYAML), clusterRoleBindings); err != nil {
			klog.Errorf("fail to unmarshal ClusterRoleBinding yaml, maybe version is not valid, you can run: kubectl create clusterrolebinding kubelet-bootstrap --clusterrole=system:node-bootstrapper --user=kubelet-bootstrap instead")
			return nil
		}

		if _, err := k8sClient.RbacV1().ClusterRoleBindings().Create(s.ctx, clusterRoleBindings, metav1.CreateOptions{}); err != nil {
			klog.Errorf("fail to create clusterrolebinding for kubelet-bootstrap")
			return err
		}
	}

	// add kube-apiserver access to kubelet-server role-binding
	if _, err := k8sClient.RbacV1().ClusterRoleBindings().Get(s.ctx, "kube-apiserver:kubelet-apis", metav1.GetOptions{}); err != nil {
		clusterRoleBindings := &rbacv1.ClusterRoleBinding{}
		if err := yaml.Unmarshal([]byte(rolebindingAccessKubeletYAML), clusterRoleBindings); err != nil {
			klog.Errorf("fail to unmarshal ClusterRoleBinding yaml, maybe version is not valid, you can run: kubectl create clusterrolebinding kube-apiserver:kubelet-apis --clusterrole=system:kubelet-api-admin --user system:kube-apiserver instead")
			return nil
		}

		if _, err := k8sClient.RbacV1().ClusterRoleBindings().Create(s.ctx, clusterRoleBindings, metav1.CreateOptions{}); err != nil {
			klog.Errorf("fail to create clusterrolebinding for kube-apiserver")
			return err
		}
	}

	return nil
}

// func (s *KubernatesServer) StartUpConfig() (*http.Handler, *authenticator.Request) {
// 	return s.apiserver.StartUpConfig()
// }
