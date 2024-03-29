package k8scall

import (
	"context"
	"fmt"

	"google.golang.org/grpc"

	"k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/api/meta"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	kyaml "k8s.io/apimachinery/pkg/runtime/serializer/yaml"
	"k8s.io/client-go/discovery"
	"k8s.io/client-go/discovery/cached/memory"
	"k8s.io/client-go/dynamic"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/restmapper"

	v1k8scall "github.com/Jooho/integration-framework-server/pkg/api/k8scall/v1"
	"github.com/Jooho/integration-framework-server/pkg/logger"
)

type k8sCallServer struct {
	Scheme *runtime.Scheme
	v1k8scall.K8SCallServer
	clientset *kubernetes.Clientset
	config    *rest.Config
}

func NewK8sCallServer(s grpc.Server, scheme *runtime.Scheme, c *kubernetes.Clientset, r *rest.Config) {
	logger.Log.Info("Adding Storage Service to gRPC server...")

	k8sCallServer := &k8sCallServer{}
	k8sCallServer.Scheme = scheme
	k8sCallServer.clientset = c
	k8sCallServer.config = r

	v1k8scall.RegisterK8SCallServer(&s, k8sCallServer)
}

func (k *k8sCallServer) CreateObjectByJson(ctx context.Context, req *v1k8scall.K8SJson) (*v1k8scall.CreateObjectByFileResponse, error) {
	logger.Log.Debug("Entry k8scall.go - CreateObjectByJson")
	logger.Log.Debug(fmt.Sprintf("File String: '%s'", req.Manifest[:]))
	logger.Log.Debug(fmt.Sprintf("Namespace: '%s'", req.Namespace))
	
	reqByteJson := req.Manifest

	// 1. Prepare a RESTMapper to find GVR
	dc, err := discovery.NewDiscoveryClientForConfig(k.config)
	if err != nil {
		logger.Log.Error(fmt.Sprintf("Fail to create discoveryclient: '%s' ", err.Error()))
		return &v1k8scall.CreateObjectByFileResponse{
			Ok:          false,
			Description: "Internal Server Error",
		}, err
	}
	mapper := restmapper.NewDeferredDiscoveryRESTMapper(memory.NewMemCacheClient(dc))

	// 2. Prepare the dynamic client
	dyn, err := dynamic.NewForConfig(k.config)
	if err != nil {
		logger.Log.Error(fmt.Sprintf("Error creating dynamic client: %s ", err.Error()))
		return &v1k8scall.CreateObjectByFileResponse{
			Ok:          false,
			Description: "Internal Server Error",
		}, err
	}

	// 3. Decode YAML manifest into unstructured.Unstructured
	unstructedObj := &unstructured.Unstructured{}
	decUnstructured := kyaml.NewDecodingSerializer(unstructured.UnstructuredJSONScheme)
	_, gvk, err := decUnstructured.Decode([]byte(reqByteJson), nil, unstructedObj)
	if err != nil {
		logger.Log.Error(fmt.Sprintf("Error getting unstructed data: %v", err))
		return &v1k8scall.CreateObjectByFileResponse{
			Ok:          false,
			Description: "Internal Server Error",
		}, err
	}

	// 4. Find GVR
	mapping, err := mapper.RESTMapping(gvk.GroupKind(), gvk.Version)
	if err != nil {
		logger.Log.Error(fmt.Sprintf("Error finding GVR: %v", err))
		return &v1k8scall.CreateObjectByFileResponse{
			Ok:          false,
			Description: "Internal Server Error",
		}, err
	}

	// 5. Obtain REST interface for the GVR
	var dr dynamic.ResourceInterface
	if mapping.Scope.Name() == meta.RESTScopeNameNamespace {
		// namespaced resources should specify the namespace
		dr = dyn.Resource(mapping.Resource).Namespace(req.Namespace)
	} else {
		// for cluster-wide resources
		dr = dyn.Resource(mapping.Resource)
	}

	createdObject, err := dr.Create(context.Background(), unstructedObj, metav1.CreateOptions{})
	if err != nil {
		logger.Log.Error(fmt.Sprintf("Error creating a resource: %v", err))
		if errors.IsAlreadyExists(err) {
			logger.Log.Error("it exists")
			return &v1k8scall.CreateObjectByFileResponse{
				Ok:          false,
				Description: fmt.Sprintf("The object already exist:(%s)",err.Error()),
			}, err
		}

	}

	return &v1k8scall.CreateObjectByFileResponse{
		Ok:          true,
		Description: fmt.Sprintf("Object %s successfully created in %s namespace", createdObject.GetKind(), createdObject.GetNamespace()),
	}, nil

}
