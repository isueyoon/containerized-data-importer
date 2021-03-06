package v1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	operatorv1 "github.com/openshift/api/operator/v1"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ConfigList is a slice of Config objects.
type ConfigList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata" protobuf:"bytes,1,opt,name=metadata"`
	Items           []Config `json:"items" protobuf:"bytes,2,rep,name=items"`
}

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Config is the configuration object for a registry instance managed by
// the registry operator
type Config struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata" protobuf:"bytes,1,opt,name=metadata"`

	Spec ImageRegistrySpec `json:"spec" protobuf:"bytes,2,opt,name=spec"`
	// +optional
	Status ImageRegistryStatus `json:"status" protobuf:"bytes,3,opt,name=status"`
}

// ImageRegistrySpec defines the specs for the running registry.
type ImageRegistrySpec struct {
	// managementState indicates whether the registry instance represented
	// by this config instance is under operator management or not.  Valid
	// values are Managed, Unmanaged, and Removed.
	ManagementState operatorv1.ManagementState `json:"managementState" protobuf:"bytes,1,opt,name=managementState,casttype=github.com/openshift/api/operator/v1.ManagementState"`
	// httpSecret is the value needed by the registry to secure uploads, generated by default.
	// +optional
	HTTPSecret string `json:"httpSecret" protobuf:"bytes,2,opt,name=httpSecret"`
	// proxy defines the proxy to be used when calling master api, upstream
	// registries, etc.
	// +optional
	Proxy ImageRegistryConfigProxy `json:"proxy" protobuf:"bytes,3,opt,name=proxy"`
	// storage details for configuring registry storage, e.g. S3 bucket
	// coordinates.
	// +optional
	Storage ImageRegistryConfigStorage `json:"storage" protobuf:"bytes,4,opt,name=storage"`
	// readOnly indicates whether the registry instance should reject attempts
	// to push new images or delete existing ones.
	// +optional
	ReadOnly bool `json:"readOnly" protobuf:"varint,5,opt,name=readOnly"`
	// disableRedirect controls whether to route all data through the Registry,
	// rather than redirecting to the backend.
	// +optional
	DisableRedirect bool `json:"disableRedirect" protobuf:"varint,6,opt,name=disableRedirect"`
	// requests controls how many parallel requests a given registry instance
	// will handle before queuing additional requests.
	// +optional
	Requests ImageRegistryConfigRequests `json:"requests" protobuf:"bytes,7,opt,name=requests"`
	// defaultRoute indicates whether an external facing route for the registry
	// should be created using the default generated hostname.
	// +optional
	DefaultRoute bool `json:"defaultRoute" protobuf:"varint,8,opt,name=defaultRoute"`
	// routes defines additional external facing routes which should be
	// created for the registry.
	// +optional
	Routes []ImageRegistryConfigRoute `json:"routes,omitempty" protobuf:"bytes,9,rep,name=routes"`
	// replicas determines the number of registry instances to run.
	Replicas int32 `json:"replicas" protobuf:"varint,10,opt,name=replicas"`
	// logging determines the level of logging enabled in the registry.
	LogLevel int64 `json:"logging" protobuf:"varint,11,opt,name=logging"`
	// resources defines the resource requests+limits for the registry pod.
	// +optional
	Resources *corev1.ResourceRequirements `json:"resources,omitempty" protobuf:"bytes,12,opt,name=resources"`
	// nodeSelector defines the node selection constraints for the registry
	// pod.
	// +optional
	NodeSelector map[string]string `json:"nodeSelector,omitempty" protobuf:"bytes,13,rep,name=nodeSelector"`
	// tolerations defines the tolerations for the registry pod.
	// +optional
	Tolerations []corev1.Toleration `json:"tolerations,omitempty" protobuf:"bytes,14,rep,name=tolerations"`
}

// ImageRegistryStatus reports image registry operational status.
type ImageRegistryStatus struct {
	operatorv1.OperatorStatus `json:",inline" protobuf:"bytes,1,opt,name=operatorStatus"`

	// storageManaged is a boolean which denotes whether or not
	// we created the registry storage medium (such as an
	// S3 bucket).
	StorageManaged bool `json:"storageManaged" protobuf:"varint,2,opt,name=storageManaged"`
	// storage indicates the current applied storage configuration of the
	// registry.
	Storage ImageRegistryConfigStorage `json:"storage" protobuf:"bytes,3,opt,name=storage"`
}

// ImageRegistryConfigProxy defines proxy configuration to be used by registry.
type ImageRegistryConfigProxy struct {
	// http defines the proxy to be used by the image registry when
	// accessing HTTP endpoints.
	// +optional
	HTTP string `json:"http" protobuf:"bytes,1,opt,name=http"`
	// https defines the proxy to be used by the image registry when
	// accessing HTTPS endpoints.
	// +optional
	HTTPS string `json:"https" protobuf:"bytes,2,opt,name=https"`
	// noProxy defines a comma-separated list of host names that shouldn't
	// go through any proxy.
	// +optional
	NoProxy string `json:"noProxy" protobuf:"bytes,3,opt,name=noProxy"`
}

// ImageRegistryConfigStorageS3CloudFront holds the configuration
// to use Amazon Cloudfront as the storage middleware in a registry.
// https://docs.docker.com/registry/configuration/#cloudfront
type ImageRegistryConfigStorageS3CloudFront struct {
	// baseURL contains the SCHEME://HOST[/PATH] at which Cloudfront is served.
	BaseURL string `json:"baseURL" protobuf:"bytes,1,opt,name=baseURL"`
	// privateKey points to secret containing the private key, provided by AWS.
	PrivateKey corev1.SecretKeySelector `json:"privateKey" protobuf:"bytes,2,opt,name=privateKey"`
	// keypairID is key pair ID provided by AWS.
	KeypairID string `json:"keypairID" protobuf:"bytes,3,opt,name=keypairID"`
	// duration is the duration of the Cloudfront session.
	// +optional
	Duration metav1.Duration `json:"duration" protobuf:"bytes,4,opt,name=duration"`
}

// ImageRegistryConfigStorageEmptyDir is an place holder to be used when
// when registry is leveraging ephemeral storage.
type ImageRegistryConfigStorageEmptyDir struct {
}

// ImageRegistryConfigStorageS3 holds the information to configure
// the registry to use the AWS S3 service for backend storage
// https://docs.docker.com/registry/storage-drivers/s3/
type ImageRegistryConfigStorageS3 struct {
	// bucket is the bucket name in which you want to store the registry's
	// data.
	// Optional, will be generated if not provided.
	// +optional
	Bucket string `json:"bucket" protobuf:"bytes,1,opt,name=bucket"`
	// region is the AWS region in which your bucket exists.
	// Optional, will be set based on the installed AWS Region.
	// +optional
	Region string `json:"region" protobuf:"bytes,2,opt,name=region"`
	// regionEndpoint is the endpoint for S3 compatible storage services.
	// Optional, defaults based on the Region that is provided.
	// +optional
	RegionEndpoint string `json:"regionEndpoint" protobuf:"bytes,3,opt,name=regionEndpoint"`
	// encrypt specifies whether the registry stores the image in encrypted
	// format or not.
	// Optional, defaults to false.
	// +optional
	Encrypt bool `json:"encrypt" protobuf:"varint,4,opt,name=encrypt"`
	// keyID is the KMS key ID to use for encryption.
	// Optional, Encrypt must be true, or this parameter is ignored.
	// +optional
	KeyID string `json:"keyID" protobuf:"bytes,5,opt,name=keyID"`
	// cloudFront configures Amazon Cloudfront as the storage middleware in a
	// registry.
	// +optional
	CloudFront *ImageRegistryConfigStorageS3CloudFront `json:"cloudFront,omitempty" protobuf:"bytes,6,opt,name=cloudFront"`
}

// ImageRegistryConfigStorageGCS holds GCS configuration.
type ImageRegistryConfigStorageGCS struct {
	// bucket is the bucket name in which you want to store the registry's
	// data.
	// Optional, will be generated if not provided.
	// +optional
	Bucket string `json:"bucket,omitempty" protobuf:"bytes,1,opt,name=bucket"`
	// region is the GCS location in which your bucket exists.
	// Optional, will be set based on the installed GCS Region.
	// +optional
	Region string `json:"region,omitempty" protobuf:"bytes,2,opt,name=region"`
	// projectID is the Project ID of the GCP project that this bucket should
	// be associated with.
	// +optional
	ProjectID string `json:"projectID,omitempty" protobuf:"bytes,3,opt,name=projectID"`
	// keyID is the KMS key ID to use for encryption.
	// Optional, buckets are encrypted by default on GCP.
	// This allows for the use of a custom encryption key.
	// +optional
	KeyID string `json:"keyID,omitempty" protobuf:"bytes,4,opt,name=keyID"`
}

// ImageRegistryConfigStorageSwift holds the information to configure
// the registry to use the OpenStack Swift service for backend storage
// https://docs.docker.com/registry/storage-drivers/swift/
type ImageRegistryConfigStorageSwift struct {
	// authURL defines the URL for obtaining an authentication token.
	// +optional
	AuthURL string `json:"authURL" protobuf:"bytes,1,opt,name=authURL"`
	// authVersion specifies the OpenStack Auth's version.
	// +optional
	AuthVersion string `json:"authVersion" protobuf:"bytes,2,opt,name=authVersion"`
	// container defines the name of Swift container where to store the
	// registry's data.
	// +optional
	Container string `json:"container" protobuf:"bytes,3,opt,name=container"`
	// domain specifies Openstack's domain name for Identity v3 API.
	// +optional
	Domain string `json:"domain" protobuf:"bytes,4,opt,name=domain"`
	// domainID specifies Openstack's domain id for Identity v3 API.
	// +optional
	DomainID string `json:"domainID" protobuf:"bytes,5,opt,name=domainID"`
	// tenant defines Openstack tenant name to be used by registry.
	// +optional
	Tenant string `json:"tenant" protobuf:"bytes,6,opt,name=tenant"`
	// tenant defines Openstack tenant id to be used by registry.
	// +optional
	TenantID string `json:"tenantID" protobuf:"bytes,7,opt,name=tenantID"`
	// regionName defines Openstack's region in which container exists.
	// +optional
	RegionName string `json:"regionName" protobuf:"bytes,8,opt,name=regionName"`
}

// ImageRegistryConfigStoragePVC holds Persistent Volume Claims data to
// be used by the registry.
type ImageRegistryConfigStoragePVC struct {
	// claim defines the Persisent Volume Claim's name to be used.
	// +optional
	Claim string `json:"claim" protobuf:"bytes,1,opt,name=claim"`
}

// ImageRegistryConfigStorageAzure holds the information to configure
// the registry to use Azure Blob Storage for backend storage.
type ImageRegistryConfigStorageAzure struct {
	// accountName defines the account to be used by the registry.
	// +optional
	AccountName string `json:"accountName" protobuf:"bytes,1,opt,name=accountName"`
	// container defines Azure's container to be used by registry.
	// +optional
	Container string `json:"container" protobuf:"bytes,2,opt,name=container"`
}

// ImageRegistryConfigStorage describes how the storage should be configured
// for the image registry.
type ImageRegistryConfigStorage struct {
	// emptyDir represents ephemeral storage on the pod's host node.
	// WARNING: this storage cannot be used with more than 1 replica and
	// is not suitable for production use. When the pod is removed from a
	// node for any reason, the data in the emptyDir is deleted forever.
	// +optional
	EmptyDir *ImageRegistryConfigStorageEmptyDir `json:"emptyDir,omitempty" protobuf:"bytes,1,opt,name=emptyDir"`
	// s3 represents configuration that uses Amazon Simple Storage Service.
	// +optional
	S3 *ImageRegistryConfigStorageS3 `json:"s3,omitempty" protobuf:"bytes,2,opt,name=s3"`
	// gcs represents configuration that uses Google Cloud Storage.
	// +optional
	GCS *ImageRegistryConfigStorageGCS `json:"gcs,omitempty" protobuf:"bytes,3,opt,name=gcs"`
	// swift represents configuration that uses OpenStack Object Storage.
	// +optional
	Swift *ImageRegistryConfigStorageSwift `json:"swift,omitempty" protobuf:"bytes,4,opt,name=swift"`
	// pvc represents configuration that uses a PersistentVolumeClaim.
	// +optional
	PVC *ImageRegistryConfigStoragePVC `json:"pvc,omitempty" protobuf:"bytes,5,opt,name=pvc"`
	// azure represents configuration that uses Azure Blob Storage.
	// +optional
	Azure *ImageRegistryConfigStorageAzure `json:"azure,omitempty" protobuf:"bytes,6,opt,name=azure"`
}

// ImageRegistryConfigRequests defines registry limits on requests read and write.
type ImageRegistryConfigRequests struct {
	// read defines limits for image registry's reads.
	// +optional
	Read ImageRegistryConfigRequestsLimits `json:"read" protobuf:"bytes,1,opt,name=read"`
	// write defines limits for image registry's writes.
	// +optional
	Write ImageRegistryConfigRequestsLimits `json:"write" protobuf:"bytes,2,opt,name=write"`
}

// ImageRegistryConfigRequestsLimits holds configuration on the max, enqueued
// and waiting registry's API requests.
type ImageRegistryConfigRequestsLimits struct {
	// maxRunning sets the maximum in flight api requests to the registry.
	// +optional
	MaxRunning int `json:"maxRunning" protobuf:"varint,1,opt,name=maxRunning"`
	// maxInQueue sets the maximum queued api requests to the registry.
	// +optional
	MaxInQueue int `json:"maxInQueue" protobuf:"varint,2,opt,name=maxInQueue"`
	// maxWaitInQueue sets the maximum time a request can wait in the queue
	// before being rejected.
	// +optional
	MaxWaitInQueue metav1.Duration `json:"maxWaitInQueue" protobuf:"bytes,3,opt,name=maxWaitInQueue"`
}

// ImageRegistryConfigRoute holds information on external route access to image
// registry.
type ImageRegistryConfigRoute struct {
	// name of the route to be created.
	Name string `json:"name" protobuf:"bytes,1,opt,name=name"`
	// hostname for the route.
	// +optional
	Hostname string `json:"hostname,omitempty" protobuf:"bytes,2,opt,name=hostname"`
	// secretName points to secret containing the certificates to be used
	// by the route.
	// +optional
	SecretName string `json:"secretName,omitempty" protobuf:"bytes,3,opt,name=secretName"`
}
