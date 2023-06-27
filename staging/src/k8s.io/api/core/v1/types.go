/*
Copyright 2015 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package v1

import (
	"k8s.io/apimachinery/pkg/api/resource"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/apimachinery/pkg/util/intstr"
)

const (
	// NamespaceDefault means the object is in the default namespace which is applied when not specified by clients
	NamespaceDefault string = "default"
	// NamespaceAll is the default argument to specify on a context when you want to list or filter resources across all namespaces
	NamespaceAll string = ""
	// NamespaceNodeLease is the namespace where we place node lease objects (used for node heartbeats)
	NamespaceNodeLease string = "kube-node-lease"
)

// Volume represents a named volume in a pod that may be accessed by any container in the pod.
type Volume struct {
	// name of the volume.
	// Must be a DNS_LABEL and unique within the pod.
	// More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names
	Name string `json:"name" protobuf:"bytes,1,opt,name=name"`
	// volumeSource represents the location and type of the mounted volume.
	// If not specified, the Volume is implied to be an EmptyDir.
	// This implied behavior is deprecated and will be removed in a future version.
	VolumeSource `json:",inline" protobuf:"bytes,2,opt,name=volumeSource"`
}

// Represents the source of a volume to mount.
// Only one of its members may be specified.
type VolumeSource struct {
	// hostPath represents a pre-existing file or directory on the host
	// machine that is directly exposed to the container. This is generally
	// used for system agents or other privileged things that are allowed
	// to see the host machine. Most containers will NOT need this.
	// More info: https://kubernetes.io/docs/concepts/storage/volumes#hostpath
	// ---
	// TODO(jonesdl) We need to restrict who can use host directory mounts and who can/can not
	// mount host directories as read/write.
	// +optional
	HostPath *HostPathVolumeSource `json:"hostPath,omitempty" protobuf:"bytes,1,opt,name=hostPath"`
	// emptyDir represents a temporary directory that shares a pod's lifetime.
	// More info: https://kubernetes.io/docs/concepts/storage/volumes#emptydir
	// +optional
	EmptyDir *EmptyDirVolumeSource `json:"emptyDir,omitempty" protobuf:"bytes,2,opt,name=emptyDir"`
	// gcePersistentDisk represents a GCE Disk resource that is attached to a
	// kubelet's host machine and then exposed to the pod.
	// More info: https://kubernetes.io/docs/concepts/storage/volumes#gcepersistentdisk
	// +optional
	GCEPersistentDisk *GCEPersistentDiskVolumeSource `json:"gcePersistentDisk,omitempty" protobuf:"bytes,3,opt,name=gcePersistentDisk"`
	// awsElasticBlockStore represents an AWS Disk resource that is attached to a
	// kubelet's host machine and then exposed to the pod.
	// More info: https://kubernetes.io/docs/concepts/storage/volumes#awselasticblockstore
	// +optional
	AWSElasticBlockStore *AWSElasticBlockStoreVolumeSource `json:"awsElasticBlockStore,omitempty" protobuf:"bytes,4,opt,name=awsElasticBlockStore"`
	// gitRepo represents a git repository at a particular revision.
	// DEPRECATED: GitRepo is deprecated. To provision a container with a git repo, mount an
	// EmptyDir into an InitContainer that clones the repo using git, then mount the EmptyDir
	// into the Pod's container.
	// +optional
	GitRepo *GitRepoVolumeSource `json:"gitRepo,omitempty" protobuf:"bytes,5,opt,name=gitRepo"`
	// secret represents a secret that should populate this volume.
	// More info: https://kubernetes.io/docs/concepts/storage/volumes#secret
	// +optional
	Secret *SecretVolumeSource `json:"secret,omitempty" protobuf:"bytes,6,opt,name=secret"`
	// nfs represents an NFS mount on the host that shares a pod's lifetime
	// More info: https://kubernetes.io/docs/concepts/storage/volumes#nfs
	// +optional
	NFS *NFSVolumeSource `json:"nfs,omitempty" protobuf:"bytes,7,opt,name=nfs"`
	// iscsi represents an ISCSI Disk resource that is attached to a
	// kubelet's host machine and then exposed to the pod.
	// More info: https://examples.k8s.io/volumes/iscsi/README.md
	// +optional
	ISCSI *ISCSIVolumeSource `json:"iscsi,omitempty" protobuf:"bytes,8,opt,name=iscsi"`
	// glusterfs represents a Glusterfs mount on the host that shares a pod's lifetime.
	// More info: https://examples.k8s.io/volumes/glusterfs/README.md
	// +optional
	Glusterfs *GlusterfsVolumeSource `json:"glusterfs,omitempty" protobuf:"bytes,9,opt,name=glusterfs"`
	// persistentVolumeClaimVolumeSource represents a reference to a
	// PersistentVolumeClaim in the same namespace.
	// More info: https://kubernetes.io/docs/concepts/storage/persistent-volumes#persistentvolumeclaims
	// +optional
	PersistentVolumeClaim *PersistentVolumeClaimVolumeSource `json:"persistentVolumeClaim,omitempty" protobuf:"bytes,10,opt,name=persistentVolumeClaim"`
	// rbd represents a Rados Block Device mount on the host that shares a pod's lifetime.
	// More info: https://examples.k8s.io/volumes/rbd/README.md
	// +optional
	RBD *RBDVolumeSource `json:"rbd,omitempty" protobuf:"bytes,11,opt,name=rbd"`
	// flexVolume represents a generic volume resource that is
	// provisioned/attached using an exec based plugin.
	// +optional
	FlexVolume *FlexVolumeSource `json:"flexVolume,omitempty" protobuf:"bytes,12,opt,name=flexVolume"`
	// cinder represents a cinder volume attached and mounted on kubelets host machine.
	// More info: https://examples.k8s.io/mysql-cinder-pd/README.md
	// +optional
	Cinder *CinderVolumeSource `json:"cinder,omitempty" protobuf:"bytes,13,opt,name=cinder"`
	// cephFS represents a Ceph FS mount on the host that shares a pod's lifetime
	// +optional
	CephFS *CephFSVolumeSource `json:"cephfs,omitempty" protobuf:"bytes,14,opt,name=cephfs"`
	// flocker represents a Flocker volume attached to a kubelet's host machine. This depends on the Flocker control service being running
	// +optional
	Flocker *FlockerVolumeSource `json:"flocker,omitempty" protobuf:"bytes,15,opt,name=flocker"`
	// downwardAPI represents downward API about the pod that should populate this volume
	// +optional
	DownwardAPI *DownwardAPIVolumeSource `json:"downwardAPI,omitempty" protobuf:"bytes,16,opt,name=downwardAPI"`
	// fc represents a Fibre Channel resource that is attached to a kubelet's host machine and then exposed to the pod.
	// +optional
	FC *FCVolumeSource `json:"fc,omitempty" protobuf:"bytes,17,opt,name=fc"`
	// azureFile represents an Azure File Service mount on the host and bind mount to the pod.
	// +optional
	AzureFile *AzureFileVolumeSource `json:"azureFile,omitempty" protobuf:"bytes,18,opt,name=azureFile"`
	// configMap represents a configMap that should populate this volume
	// +optional
	ConfigMap *ConfigMapVolumeSource `json:"configMap,omitempty" protobuf:"bytes,19,opt,name=configMap"`
	// vsphereVolume represents a vSphere volume attached and mounted on kubelets host machine
	// +optional
	VsphereVolume *VsphereVirtualDiskVolumeSource `json:"vsphereVolume,omitempty" protobuf:"bytes,20,opt,name=vsphereVolume"`
	// quobyte represents a Quobyte mount on the host that shares a pod's lifetime
	// +optional
	Quobyte *QuobyteVolumeSource `json:"quobyte,omitempty" protobuf:"bytes,21,opt,name=quobyte"`
	// azureDisk represents an Azure Data Disk mount on the host and bind mount to the pod.
	// +optional
	AzureDisk *AzureDiskVolumeSource `json:"azureDisk,omitempty" protobuf:"bytes,22,opt,name=azureDisk"`
	// photonPersistentDisk represents a PhotonController persistent disk attached and mounted on kubelets host machine
	PhotonPersistentDisk *PhotonPersistentDiskVolumeSource `json:"photonPersistentDisk,omitempty" protobuf:"bytes,23,opt,name=photonPersistentDisk"`
	// projected items for all in one resources secrets, configmaps, and downward API
	Projected *ProjectedVolumeSource `json:"projected,omitempty" protobuf:"bytes,26,opt,name=projected"`
	// portworxVolume represents a portworx volume attached and mounted on kubelets host machine
	// +optional
	PortworxVolume *PortworxVolumeSource `json:"portworxVolume,omitempty" protobuf:"bytes,24,opt,name=portworxVolume"`
	// scaleIO represents a ScaleIO persistent volume attached and mounted on Kubernetes nodes.
	// +optional
	ScaleIO *ScaleIOVolumeSource `json:"scaleIO,omitempty" protobuf:"bytes,25,opt,name=scaleIO"`
	// storageOS represents a StorageOS volume attached and mounted on Kubernetes nodes.
	// +optional
	StorageOS *StorageOSVolumeSource `json:"storageos,omitempty" protobuf:"bytes,27,opt,name=storageos"`
	// csi (Container Storage Interface) represents ephemeral storage that is handled by certain external CSI drivers (Beta feature).
	// +optional
	CSI *CSIVolumeSource `json:"csi,omitempty" protobuf:"bytes,28,opt,name=csi"`
	// ephemeral represents a volume that is handled by a cluster storage driver.
	// The volume's lifecycle is tied to the pod that defines it - it will be created before the pod starts,
	// and deleted when the pod is removed.
	//
	// Use this if:
	// a) the volume is only needed while the pod runs,
	// b) features of normal volumes like restoring from snapshot or capacity
	//    tracking are needed,
	// c) the storage driver is specified through a storage class, and
	// d) the storage driver supports dynamic volume provisioning through
	//    a PersistentVolumeClaim (see EphemeralVolumeSource for more
	//    information on the connection between this volume type
	//    and PersistentVolumeClaim).
	//
	// Use PersistentVolumeClaim or one of the vendor-specific
	// APIs for volumes that persist for longer than the lifecycle
	// of an individual pod.
	//
	// Use CSI for light-weight local ephemeral volumes if the CSI driver is meant to
	// be used that way - see the documentation of the driver for
	// more information.
	//
	// A pod can use both types of ephemeral volumes and
	// persistent volumes at the same time.
	//
	// +optional
	Ephemeral *EphemeralVolumeSource `json:"ephemeral,omitempty" protobuf:"bytes,29,opt,name=ephemeral"`
}

// PersistentVolumeClaimVolumeSource references the user's PVC in the same namespace.
// This volume finds the bound PV and mounts that volume for the pod. A
// PersistentVolumeClaimVolumeSource is, essentially, a wrapper around another
// type of volume that is owned by someone else (the system).
type PersistentVolumeClaimVolumeSource struct {
	// claimName is the name of a PersistentVolumeClaim in the same namespace as the pod using this volume.
	// More info: https://kubernetes.io/docs/concepts/storage/persistent-volumes#persistentvolumeclaims
	ClaimName string `json:"claimName" protobuf:"bytes,1,opt,name=claimName"`
	// readOnly Will force the ReadOnly setting in VolumeMounts.
	// Default false.
	// +optional
	ReadOnly bool `json:"readOnly,omitempty" protobuf:"varint,2,opt,name=readOnly"`
}

// PersistentVolumeSource is similar to VolumeSource but meant for the
// administrator who creates PVs. Exactly one of its members must be set.
type PersistentVolumeSource struct {
	// gcePersistentDisk represents a GCE Disk resource that is attached to a
	// kubelet's host machine and then exposed to the pod. Provisioned by an admin.
	// More info: https://kubernetes.io/docs/concepts/storage/volumes#gcepersistentdisk
	// +optional
	GCEPersistentDisk *GCEPersistentDiskVolumeSource `json:"gcePersistentDisk,omitempty" protobuf:"bytes,1,opt,name=gcePersistentDisk"`
	// awsElasticBlockStore represents an AWS Disk resource that is attached to a
	// kubelet's host machine and then exposed to the pod.
	// More info: https://kubernetes.io/docs/concepts/storage/volumes#awselasticblockstore
	// +optional
	AWSElasticBlockStore *AWSElasticBlockStoreVolumeSource `json:"awsElasticBlockStore,omitempty" protobuf:"bytes,2,opt,name=awsElasticBlockStore"`
	// hostPath represents a directory on the host.
	// Provisioned by a developer or tester.
	// This is useful for single-node development and testing only!
	// On-host storage is not supported in any way and WILL NOT WORK in a multi-node cluster.
	// More info: https://kubernetes.io/docs/concepts/storage/volumes#hostpath
	// +optional
	HostPath *HostPathVolumeSource `json:"hostPath,omitempty" protobuf:"bytes,3,opt,name=hostPath"`
	// glusterfs represents a Glusterfs volume that is attached to a host and
	// exposed to the pod. Provisioned by an admin.
	// More info: https://examples.k8s.io/volumes/glusterfs/README.md
	// +optional
	Glusterfs *GlusterfsPersistentVolumeSource `json:"glusterfs,omitempty" protobuf:"bytes,4,opt,name=glusterfs"`
	// nfs represents an NFS mount on the host. Provisioned by an admin.
	// More info: https://kubernetes.io/docs/concepts/storage/volumes#nfs
	// +optional
	NFS *NFSVolumeSource `json:"nfs,omitempty" protobuf:"bytes,5,opt,name=nfs"`
	// rbd represents a Rados Block Device mount on the host that shares a pod's lifetime.
	// More info: https://examples.k8s.io/volumes/rbd/README.md
	// +optional
	RBD *RBDPersistentVolumeSource `json:"rbd,omitempty" protobuf:"bytes,6,opt,name=rbd"`
	// iscsi represents an ISCSI Disk resource that is attached to a
	// kubelet's host machine and then exposed to the pod. Provisioned by an admin.
	// +optional
	ISCSI *ISCSIPersistentVolumeSource `json:"iscsi,omitempty" protobuf:"bytes,7,opt,name=iscsi"`
	// cinder represents a cinder volume attached and mounted on kubelets host machine.
	// More info: https://examples.k8s.io/mysql-cinder-pd/README.md
	// +optional
	Cinder *CinderPersistentVolumeSource `json:"cinder,omitempty" protobuf:"bytes,8,opt,name=cinder"`
	// cephFS represents a Ceph FS mount on the host that shares a pod's lifetime
	// +optional
	CephFS *CephFSPersistentVolumeSource `json:"cephfs,omitempty" protobuf:"bytes,9,opt,name=cephfs"`
	// fc represents a Fibre Channel resource that is attached to a kubelet's host machine and then exposed to the pod.
	// +optional
	FC *FCVolumeSource `json:"fc,omitempty" protobuf:"bytes,10,opt,name=fc"`
	// flocker represents a Flocker volume attached to a kubelet's host machine and exposed to the pod for its usage. This depends on the Flocker control service being running
	// +optional
	Flocker *FlockerVolumeSource `json:"flocker,omitempty" protobuf:"bytes,11,opt,name=flocker"`
	// flexVolume represents a generic volume resource that is
	// provisioned/attached using an exec based plugin.
	// +optional
	FlexVolume *FlexPersistentVolumeSource `json:"flexVolume,omitempty" protobuf:"bytes,12,opt,name=flexVolume"`
	// azureFile represents an Azure File Service mount on the host and bind mount to the pod.
	// +optional
	AzureFile *AzureFilePersistentVolumeSource `json:"azureFile,omitempty" protobuf:"bytes,13,opt,name=azureFile"`
	// vsphereVolume represents a vSphere volume attached and mounted on kubelets host machine
	// +optional
	VsphereVolume *VsphereVirtualDiskVolumeSource `json:"vsphereVolume,omitempty" protobuf:"bytes,14,opt,name=vsphereVolume"`
	// quobyte represents a Quobyte mount on the host that shares a pod's lifetime
	// +optional
	Quobyte *QuobyteVolumeSource `json:"quobyte,omitempty" protobuf:"bytes,15,opt,name=quobyte"`
	// azureDisk represents an Azure Data Disk mount on the host and bind mount to the pod.
	// +optional
	AzureDisk *AzureDiskVolumeSource `json:"azureDisk,omitempty" protobuf:"bytes,16,opt,name=azureDisk"`
	// photonPersistentDisk represents a PhotonController persistent disk attached and mounted on kubelets host machine
	PhotonPersistentDisk *PhotonPersistentDiskVolumeSource `json:"photonPersistentDisk,omitempty" protobuf:"bytes,17,opt,name=photonPersistentDisk"`
	// portworxVolume represents a portworx volume attached and mounted on kubelets host machine
	// +optional
	PortworxVolume *PortworxVolumeSource `json:"portworxVolume,omitempty" protobuf:"bytes,18,opt,name=portworxVolume"`
	// scaleIO represents a ScaleIO persistent volume attached and mounted on Kubernetes nodes.
	// +optional
	ScaleIO *ScaleIOPersistentVolumeSource `json:"scaleIO,omitempty" protobuf:"bytes,19,opt,name=scaleIO"`
	// local represents directly-attached storage with node affinity
	// +optional
	Local *LocalVolumeSource `json:"local,omitempty" protobuf:"bytes,20,opt,name=local"`
	// storageOS represents a StorageOS volume that is attached to the kubelet's host machine and mounted into the pod
	// More info: https://examples.k8s.io/volumes/storageos/README.md
	// +optional
	StorageOS *StorageOSPersistentVolumeSource `json:"storageos,omitempty" protobuf:"bytes,21,opt,name=storageos"`
	// csi represents storage that is handled by an external CSI driver (Beta feature).
	// +optional
	CSI *CSIPersistentVolumeSource `json:"csi,omitempty" protobuf:"bytes,22,opt,name=csi"`
}

const (
	// BetaStorageClassAnnotation represents the beta/previous StorageClass annotation.
	// It's currently still used and will be held for backwards compatibility
	BetaStorageClassAnnotation = "volume.beta.kubernetes.io/storage-class"

	// MountOptionAnnotation defines mount option annotation used in PVs
	MountOptionAnnotation = "volume.beta.kubernetes.io/mount-options"
)

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// PersistentVolume (PV) is a storage resource provisioned by an administrator.
// It is analogous to a node.
// More info: https://kubernetes.io/docs/concepts/storage/persistent-volumes
type PersistentVolume struct {
	metav1.TypeMeta `json:",inline"`
	// Standard object's metadata.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// spec defines a specification of a persistent volume owned by the cluster.
	// Provisioned by an administrator.
	// More info: https://kubernetes.io/docs/concepts/storage/persistent-volumes#persistent-volumes
	// +optional
	Spec PersistentVolumeSpec `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`

	// status represents the current information/status for the persistent volume.
	// Populated by the system.
	// Read-only.
	// More info: https://kubernetes.io/docs/concepts/storage/persistent-volumes#persistent-volumes
	// +optional
	Status PersistentVolumeStatus `json:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}

// PersistentVolumeSpec is the specification of a persistent volume.
type PersistentVolumeSpec struct {
	// capacity is the description of the persistent volume's resources and capacity.
	// More info: https://kubernetes.io/docs/concepts/storage/persistent-volumes#capacity
	// +optional
	Capacity ResourceList `json:"capacity,omitempty" protobuf:"bytes,1,rep,name=capacity,casttype=ResourceList,castkey=ResourceName"`
	// persistentVolumeSource is the actual volume backing the persistent volume.
	PersistentVolumeSource `json:",inline" protobuf:"bytes,2,opt,name=persistentVolumeSource"`
	// accessModes contains all ways the volume can be mounted.
	// More info: https://kubernetes.io/docs/concepts/storage/persistent-volumes#access-modes
	// +optional
	AccessModes []PersistentVolumeAccessMode `json:"accessModes,omitempty" protobuf:"bytes,3,rep,name=accessModes,casttype=PersistentVolumeAccessMode"`
	// claimRef is part of a bi-directional binding between PersistentVolume and PersistentVolumeClaim.
	// Expected to be non-nil when bound.
	// claim.VolumeName is the authoritative bind between PV and PVC.
	// More info: https://kubernetes.io/docs/concepts/storage/persistent-volumes#binding
	// +optional
	// +structType=granular
	ClaimRef *ObjectReference `json:"claimRef,omitempty" protobuf:"bytes,4,opt,name=claimRef"`
	// persistentVolumeReclaimPolicy defines what happens to a persistent volume when released from its claim.
	// Valid options are Retain (default for manually created PersistentVolumes), Delete (default
	// for dynamically provisioned PersistentVolumes), and Recycle (deprecated).
	// Recycle must be supported by the volume plugin underlying this PersistentVolume.
	// More info: https://kubernetes.io/docs/concepts/storage/persistent-volumes#reclaiming
	// +optional
	PersistentVolumeReclaimPolicy PersistentVolumeReclaimPolicy `json:"persistentVolumeReclaimPolicy,omitempty" protobuf:"bytes,5,opt,name=persistentVolumeReclaimPolicy,casttype=PersistentVolumeReclaimPolicy"`
	// storageClassName is the name of StorageClass to which this persistent volume belongs. Empty value
	// means that this volume does not belong to any StorageClass.
	// +optional
	StorageClassName string `json:"storageClassName,omitempty" protobuf:"bytes,6,opt,name=storageClassName"`
	// mountOptions is the list of mount options, e.g. ["ro", "soft"]. Not validated - mount will
	// simply fail if one is invalid.
	// More info: https://kubernetes.io/docs/concepts/storage/persistent-volumes/#mount-options
	// +optional
	MountOptions []string `json:"mountOptions,omitempty" protobuf:"bytes,7,opt,name=mountOptions"`
	// volumeMode defines if a volume is intended to be used with a formatted filesystem
	// or to remain in raw block state. Value of Filesystem is implied when not included in spec.
	// +optional
	VolumeMode *PersistentVolumeMode `json:"volumeMode,omitempty" protobuf:"bytes,8,opt,name=volumeMode,casttype=PersistentVolumeMode"`
	// nodeAffinity defines constraints that limit what nodes this volume can be accessed from.
	// This field influences the scheduling of pods that use this volume.
	// +optional
	NodeAffinity *VolumeNodeAffinity `json:"nodeAffinity,omitempty" protobuf:"bytes,9,opt,name=nodeAffinity"`
}

// VolumeNodeAffinity defines constraints that limit what nodes this volume can be accessed from.
type VolumeNodeAffinity struct {
	// required specifies hard node constraints that must be met.
	Required *NodeSelector `json:"required,omitempty" protobuf:"bytes,1,opt,name=required"`
}

// PersistentVolumeReclaimPolicy describes a policy for end-of-life maintenance of persistent volumes.
// +enum
type PersistentVolumeReclaimPolicy string

const (
	// PersistentVolumeReclaimRecycle means the volume will be recycled back into the pool of unbound persistent volumes on release from its claim.
	// The volume plugin must support Recycling.
	PersistentVolumeReclaimRecycle PersistentVolumeReclaimPolicy = "Recycle"
	// PersistentVolumeReclaimDelete means the volume will be deleted from Kubernetes on release from its claim.
	// The volume plugin must support Deletion.
	PersistentVolumeReclaimDelete PersistentVolumeReclaimPolicy = "Delete"
	// PersistentVolumeReclaimRetain means the volume will be left in its current phase (Released) for manual reclamation by the administrator.
	// The default policy is Retain.
	PersistentVolumeReclaimRetain PersistentVolumeReclaimPolicy = "Retain"
)

// PersistentVolumeMode describes how a volume is intended to be consumed, either Block or Filesystem.
// +enum
// PersistentVolumeMode描述如何使用卷，Block或Filesystem。
type PersistentVolumeMode string

const (
	// PersistentVolumeBlock means the volume will not be formatted with a filesystem and will remain a raw block device.
	PersistentVolumeBlock PersistentVolumeMode = "Block"
	// PersistentVolumeFilesystem means the volume will be or is formatted with a filesystem.
	PersistentVolumeFilesystem PersistentVolumeMode = "Filesystem"
)

// PersistentVolumeStatus is the current status of a persistent volume.
type PersistentVolumeStatus struct {
	// phase indicates if a volume is available, bound to a claim, or released by a claim.
	// More info: https://kubernetes.io/docs/concepts/storage/persistent-volumes#phase
	// +optional
	Phase PersistentVolumePhase `json:"phase,omitempty" protobuf:"bytes,1,opt,name=phase,casttype=PersistentVolumePhase"`
	// message is a human-readable message indicating details about why the volume is in this state.
	// +optional
	Message string `json:"message,omitempty" protobuf:"bytes,2,opt,name=message"`
	// reason is a brief CamelCase string that describes any failure and is meant
	// for machine parsing and tidy display in the CLI.
	// +optional
	Reason string `json:"reason,omitempty" protobuf:"bytes,3,opt,name=reason"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// PersistentVolumeList is a list of PersistentVolume items.
type PersistentVolumeList struct {
	metav1.TypeMeta `json:",inline"`
	// Standard list metadata.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds
	// +optional
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	// items is a list of persistent volumes.
	// More info: https://kubernetes.io/docs/concepts/storage/persistent-volumes
	Items []PersistentVolume `json:"items" protobuf:"bytes,2,rep,name=items"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// PersistentVolumeClaim is a user's request for and claim to a persistent volume
// PersistentVolumeClaim是用户对持久卷的请求和声明
type PersistentVolumeClaim struct {
	metav1.TypeMeta `json:",inline"`
	// Standard object's metadata.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
	// +optional
	// 标准对象的元数据
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// spec defines the desired characteristics of a volume requested by a pod author.
	// More info: https://kubernetes.io/docs/concepts/storage/persistent-volumes#persistentvolumeclaims
	// +optional
	// spec定义了pod作者请求的卷的所需特征
	Spec PersistentVolumeClaimSpec `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`

	// status represents the current information/status of a persistent volume claim.
	// Read-only.
	// More info: https://kubernetes.io/docs/concepts/storage/persistent-volumes#persistentvolumeclaims
	// +optional
	// status代表持久卷声明的当前信息/状态
	Status PersistentVolumeClaimStatus `json:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// PersistentVolumeClaimList is a list of PersistentVolumeClaim items.
type PersistentVolumeClaimList struct {
	metav1.TypeMeta `json:",inline"`
	// Standard list metadata.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds
	// +optional
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	// items is a list of persistent volume claims.
	// More info: https://kubernetes.io/docs/concepts/storage/persistent-volumes#persistentvolumeclaims
	Items []PersistentVolumeClaim `json:"items" protobuf:"bytes,2,rep,name=items"`
}

// PersistentVolumeClaimSpec describes the common attributes of storage devices
// and allows a Source for provider-specific attributes
// PersistentVolumeClaimSpec描述存储设备的常见属性，并允许Source用于特定于提供程序的属性
type PersistentVolumeClaimSpec struct {
	// accessModes contains the desired access modes the volume should have.
	// More info: https://kubernetes.io/docs/concepts/storage/persistent-volumes#access-modes-1
	// +optional
	// accessModes包含卷应具有的所需访问模式
	AccessModes []PersistentVolumeAccessMode `json:"accessModes,omitempty" protobuf:"bytes,1,rep,name=accessModes,casttype=PersistentVolumeAccessMode"`
	// selector is a label query over volumes to consider for binding.
	// +optional
	// selector是对要考虑绑定的卷的标签查询
	Selector *metav1.LabelSelector `json:"selector,omitempty" protobuf:"bytes,4,opt,name=selector"`
	// resources represents the minimum resources the volume should have.
	// If RecoverVolumeExpansionFailure feature is enabled users are allowed to specify resource requirements
	// that are lower than previous value but must still be higher than capacity recorded in the
	// status field of the claim.
	// More info: https://kubernetes.io/docs/concepts/storage/persistent-volumes#resources
	// +optional
	// resources代表卷应具有的最小资源, 如果启用了RecoverVolumeExpansionFailure功能，则允许用户指定资源要求，这些资源要求必须高于以前的值，但仍必须高于状态字段中记录的容量
	Resources ResourceRequirements `json:"resources,omitempty" protobuf:"bytes,2,opt,name=resources"`
	// volumeName is the binding reference to the PersistentVolume backing this claim.
	// +optional
	// volumeName是对支持此声明的PersistentVolume的绑定引用
	VolumeName string `json:"volumeName,omitempty" protobuf:"bytes,3,opt,name=volumeName"`
	// storageClassName is the name of the StorageClass required by the claim.
	// More info: https://kubernetes.io/docs/concepts/storage/persistent-volumes#class-1
	// +optional
	// storageClassName是声明所需的StorageClass的名称
	StorageClassName *string `json:"storageClassName,omitempty" protobuf:"bytes,5,opt,name=storageClassName"`
	// volumeMode defines what type of volume is required by the claim.
	// Value of Filesystem is implied when not included in claim spec.
	// +optional
	// volumeMode定义了声明所需的卷的类型, 当未包含在声明规范中时，假定值为Filesystem
	VolumeMode *PersistentVolumeMode `json:"volumeMode,omitempty" protobuf:"bytes,6,opt,name=volumeMode,casttype=PersistentVolumeMode"`
	// dataSource field can be used to specify either:
	// * An existing VolumeSnapshot object (snapshot.storage.k8s.io/VolumeSnapshot)
	// * An existing PVC (PersistentVolumeClaim)
	// If the provisioner or an external controller can support the specified data source,
	// it will create a new volume based on the contents of the specified data source.
	// When the AnyVolumeDataSource feature gate is enabled, dataSource contents will be copied to dataSourceRef,
	// and dataSourceRef contents will be copied to dataSource when dataSourceRef.namespace is not specified.
	// If the namespace is specified, then dataSourceRef will not be copied to dataSource.
	// +optional
	// dataSource字段可用于指定以下内容之一:
	// *现有的VolumeSnapshot对象(snapshot.storage.k8s.io/VolumeSnapshot)
	// *现有的PVC(PersistentVolumeClaim)
	// 如果支持指定的数据源的配置程序或外部控制器，则它将根据指定数据源的内容创建新卷
	// 启用AnyVolumeDataSource功能门时，将复制dataSource内容到dataSourceRef，将复制dataSourceRef内容到dataSource，当未指定dataSourceRef.namespace时
	// 如果指定了命名空间，则不会将dataSourceRef复制到dataSource
	DataSource *TypedLocalObjectReference `json:"dataSource,omitempty" protobuf:"bytes,7,opt,name=dataSource"`
	// dataSourceRef specifies the object from which to populate the volume with data, if a non-empty
	// volume is desired. This may be any object from a non-empty API group (non
	// core object) or a PersistentVolumeClaim object.
	// When this field is specified, volume binding will only succeed if the type of
	// the specified object matches some installed volume populator or dynamic
	// provisioner.
	// This field will replace the functionality of the dataSource field and as such
	// if both fields are non-empty, they must have the same value. For backwards
	// compatibility, when namespace isn't specified in dataSourceRef,
	// both fields (dataSource and dataSourceRef) will be set to the same
	// value automatically if one of them is empty and the other is non-empty.
	// When namespace is specified in dataSourceRef,
	// dataSource isn't set to the same value and must be empty.
	// There are three important differences between dataSource and dataSourceRef:
	// * While dataSource only allows two specific types of objects, dataSourceRef
	//   allows any non-core object, as well as PersistentVolumeClaim objects.
	// * While dataSource ignores disallowed values (dropping them), dataSourceRef
	//   preserves all values, and generates an error if a disallowed value is
	//   specified.
	// * While dataSource only allows local objects, dataSourceRef allows objects
	//   in any namespaces.
	// (Beta) Using this field requires the AnyVolumeDataSource feature gate to be enabled.
	// (Alpha) Using the namespace field of dataSourceRef requires the CrossNamespaceVolumeDataSource feature gate to be enabled.
	// +optional
	// dataSourceRef 指定要从中填充卷数据的对象，如果需要非空卷，则可以指定任何对象
	// 这可以是来自非空API组(非核心对象)或PersistentVolumeClaim对象的任何对象
	// 当指定此字段时，如果指定对象的类型与某些已安装的卷填充器或动态配置程序匹配，则卷绑定将仅成功
	// 此字段将替换dataSource字段的功能，因此，如果两个字段都是非空的，则它们必须具有相同的值
	// 为了向后兼容性，如果dataSourceRef中未指定命名空间，则如果其中一个为空且另一个非空，则两个字段(dataSource和dataSourceRef)将自动设置为相同的值
	// 如果在dataSourceRef中指定了命名空间，则dataSource不会设置为相同的值，必须为空
	// dataSource和dataSourceRef之间有三个重要的区别:
	// *虽然dataSource只允许两种特定类型的对象，但dataSourceRef允许任何非核心对象，以及PersistentVolumeClaim对象
	// *虽然dataSource忽略了不允许的值(删除它们)，但dataSourceRef保留所有值，并且如果指定了不允许的值，则会生成错误
	// *虽然dataSource只允许本地对象，但dataSourceRef允许任何命名空间中的对象
	// (Beta) 启用此字段需要启用AnyVolumeDataSource功能门
	// (Alpha) 使用dataSourceRef的命名空间字段需要启用CrossNamespaceVolumeDataSource功能门
	DataSourceRef *TypedObjectReference `json:"dataSourceRef,omitempty" protobuf:"bytes,8,opt,name=dataSourceRef"`
}

type TypedObjectReference struct {
	// APIGroup is the group for the resource being referenced.
	// If APIGroup is not specified, the specified Kind must be in the core API group.
	// For any other third-party types, APIGroup is required.
	// +optional
	APIGroup *string `json:"apiGroup" protobuf:"bytes,1,opt,name=apiGroup"`
	// Kind is the type of resource being referenced
	Kind string `json:"kind" protobuf:"bytes,2,opt,name=kind"`
	// Name is the name of resource being referenced
	Name string `json:"name" protobuf:"bytes,3,opt,name=name"`
	// Namespace is the namespace of resource being referenced
	// Note that when a namespace is specified, a gateway.networking.k8s.io/ReferenceGrant object is required in the referent namespace to allow that namespace's owner to accept the reference. See the ReferenceGrant documentation for details.
	// (Alpha) This field requires the CrossNamespaceVolumeDataSource feature gate to be enabled.
	// +featureGate=CrossNamespaceVolumeDataSource
	// +optional
	Namespace *string `json:"namespace,omitempty" protobuf:"bytes,4,opt,name=namespace"`
}

// PersistentVolumeClaimConditionType is a valid value of PersistentVolumeClaimCondition.Type
type PersistentVolumeClaimConditionType string

const (
	// PersistentVolumeClaimResizing - a user trigger resize of pvc has been started
	PersistentVolumeClaimResizing PersistentVolumeClaimConditionType = "Resizing"
	// PersistentVolumeClaimFileSystemResizePending - controller resize is finished and a file system resize is pending on node
	PersistentVolumeClaimFileSystemResizePending PersistentVolumeClaimConditionType = "FileSystemResizePending"
)

// +enum
type PersistentVolumeClaimResizeStatus string

const (
	// When expansion is complete, the empty string is set by resize controller or kubelet.
	// 当扩展完成时，resize控制器或kubelet将设置空字符串
	PersistentVolumeClaimNoExpansionInProgress PersistentVolumeClaimResizeStatus = ""
	// State set when resize controller starts expanding the volume in control-plane
	// 状态设置为resize控制器在控制平面中开始扩展卷
	PersistentVolumeClaimControllerExpansionInProgress PersistentVolumeClaimResizeStatus = "ControllerExpansionInProgress"
	// State set when expansion has failed in resize controller with a terminal error.
	// Transient errors such as timeout should not set this status and should leave ResizeStatus
	// unmodified, so as resize controller can resume the volume expansion.
	// 状态设置为resize控制器在终端错误中扩展失败, 瞬态错误(如超时)不应设置此状态, 并且应保持ResizeStatus未修改, 以便resize控制器可以恢复卷扩展
	PersistentVolumeClaimControllerExpansionFailed PersistentVolumeClaimResizeStatus = "ControllerExpansionFailed"
	// State set when resize controller has finished expanding the volume but further expansion is needed on the node.
	PersistentVolumeClaimNodeExpansionPending PersistentVolumeClaimResizeStatus = "NodeExpansionPending"
	// State set when kubelet starts expanding the volume.
	PersistentVolumeClaimNodeExpansionInProgress PersistentVolumeClaimResizeStatus = "NodeExpansionInProgress"
	// State set when expansion has failed in kubelet with a terminal error. Transient errors don't set NodeExpansionFailed.
	// 状态设置为kubelet在终端错误中扩展失败, 瞬态错误不设置NodeExpansionFailed
	PersistentVolumeClaimNodeExpansionFailed PersistentVolumeClaimResizeStatus = "NodeExpansionFailed"
)

// PersistentVolumeClaimCondition contails details about state of pvc
type PersistentVolumeClaimCondition struct {
	Type   PersistentVolumeClaimConditionType `json:"type" protobuf:"bytes,1,opt,name=type,casttype=PersistentVolumeClaimConditionType"`
	Status ConditionStatus                    `json:"status" protobuf:"bytes,2,opt,name=status,casttype=ConditionStatus"`
	// lastProbeTime is the time we probed the condition.
	// +optional
	LastProbeTime metav1.Time `json:"lastProbeTime,omitempty" protobuf:"bytes,3,opt,name=lastProbeTime"`
	// lastTransitionTime is the time the condition transitioned from one status to another.
	// +optional
	LastTransitionTime metav1.Time `json:"lastTransitionTime,omitempty" protobuf:"bytes,4,opt,name=lastTransitionTime"`
	// reason is a unique, this should be a short, machine understandable string that gives the reason
	// for condition's last transition. If it reports "ResizeStarted" that means the underlying
	// persistent volume is being resized.
	// +optional
	Reason string `json:"reason,omitempty" protobuf:"bytes,5,opt,name=reason"`
	// message is the human-readable message indicating details about last transition.
	// +optional
	Message string `json:"message,omitempty" protobuf:"bytes,6,opt,name=message"`
}

// PersistentVolumeClaimStatus is the current status of a persistent volume claim.
// PersistentVolumeClaimStatus是持久卷声明的当前状态
type PersistentVolumeClaimStatus struct {
	// phase represents the current phase of PersistentVolumeClaim.
	// +optional
	// phase 表示 PersistentVolumeClaim 的当前阶段
	Phase PersistentVolumeClaimPhase `json:"phase,omitempty" protobuf:"bytes,1,opt,name=phase,casttype=PersistentVolumeClaimPhase"`
	// accessModes contains the actual access modes the volume backing the PVC has.
	// More info: https://kubernetes.io/docs/concepts/storage/persistent-volumes#access-modes-1
	// +optional
	// accessModes 包含 PVC 后端卷的实际访问模式
	AccessModes []PersistentVolumeAccessMode `json:"accessModes,omitempty" protobuf:"bytes,2,rep,name=accessModes,casttype=PersistentVolumeAccessMode"`
	// capacity represents the actual resources of the underlying volume.
	// +optional
	// capacity 表示后端卷的实际资源
	Capacity ResourceList `json:"capacity,omitempty" protobuf:"bytes,3,rep,name=capacity,casttype=ResourceList,castkey=ResourceName"`
	// conditions is the current Condition of persistent volume claim. If underlying persistent volume is being
	// resized then the Condition will be set to 'ResizeStarted'.
	// +optional
	// +patchMergeKey=type
	// +patchStrategy=merge
	// conditions 是持久卷声明的当前条件。如果底层持久卷正在被扩容，那么条件将被设置为“ResizeStarted”。
	Conditions []PersistentVolumeClaimCondition `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type" protobuf:"bytes,4,rep,name=conditions"`
	// allocatedResources is the storage resource within AllocatedResources tracks the capacity allocated to a PVC. It may
	// be larger than the actual capacity when a volume expansion operation is requested.
	// For storage quota, the larger value from allocatedResources and PVC.spec.resources is used.
	// If allocatedResources is not set, PVC.spec.resources alone is used for quota calculation.
	// If a volume expansion capacity request is lowered, allocatedResources is only
	// lowered if there are no expansion operations in progress and if the actual volume capacity
	// is equal or lower than the requested capacity.
	// This is an alpha field and requires enabling RecoverVolumeExpansionFailure feature.
	// +featureGate=RecoverVolumeExpansionFailure
	// +optional
	// allocatedResources 是 AllocatedResources 中存储资源的跟踪容量。当请求卷扩容操作时，它可能大于实际容量。
	// 对于存储配额，从 allocatedResources 和 PVC.spec.resources 中使用较大的值。
	// 如果未设置 allocatedResources，则仅使用 PVC.spec.resources 用于配额计算。
	// 如果卷扩容容量请求被降低，则仅当没有扩容操作正在进行且实际卷容量等于或低于请求的容量时，才会降低 allocatedResources。
	// 这是一个 alpha 字段，需要启用 RecoverVolumeExpansionFailure 功能。
	AllocatedResources ResourceList `json:"allocatedResources,omitempty" protobuf:"bytes,5,rep,name=allocatedResources,casttype=ResourceList,castkey=ResourceName"`
	// resizeStatus stores status of resize operation.
	// ResizeStatus is not set by default but when expansion is complete resizeStatus is set to empty
	// string by resize controller or kubelet.
	// This is an alpha field and requires enabling RecoverVolumeExpansionFailure feature.
	// +featureGate=RecoverVolumeExpansionFailure
	// +optional
	// resizeStatus 存储扩容操作的状态。
	// 默认情况下不设置 ResizeStatus，但是当扩容完成时，resizeStatus 由扩容控制器或 kubelet 设置为空字符串。
	// 这是一个 alpha 字段，需要启用 RecoverVolumeExpansionFailure 功能。
	ResizeStatus *PersistentVolumeClaimResizeStatus `json:"resizeStatus,omitempty" protobuf:"bytes,6,opt,name=resizeStatus,casttype=PersistentVolumeClaimResizeStatus"`
}

// +enum
type PersistentVolumeAccessMode string

const (
	// can be mounted in read/write mode to exactly 1 host
	// 可以以读/写模式挂载到1个主机
	ReadWriteOnce PersistentVolumeAccessMode = "ReadWriteOnce"
	// can be mounted in read-only mode to many hosts
	// 可以以只读模式挂载到多个主机
	ReadOnlyMany PersistentVolumeAccessMode = "ReadOnlyMany"
	// can be mounted in read/write mode to many hosts
	// 可以以读/写模式挂载到多个主机
	ReadWriteMany PersistentVolumeAccessMode = "ReadWriteMany"
	// can be mounted in read/write mode to exactly 1 pod
	// cannot be used in combination with other access modes
	// 可以以读/写模式挂载到1个pod
	ReadWriteOncePod PersistentVolumeAccessMode = "ReadWriteOncePod"
)

// +enum
type PersistentVolumePhase string

const (
	// used for PersistentVolumes that are not available
	VolumePending PersistentVolumePhase = "Pending"
	// used for PersistentVolumes that are not yet bound
	// Available volumes are held by the binder and matched to PersistentVolumeClaims
	// 可用的卷由绑定器持有，并匹配到PersistentVolumeClaims
	VolumeAvailable PersistentVolumePhase = "Available"
	// used for PersistentVolumes that are bound
	// 用于绑定的persistent卷
	VolumeBound PersistentVolumePhase = "Bound"
	// used for PersistentVolumes where the bound PersistentVolumeClaim was deleted
	// released volumes must be recycled before becoming available again
	// this phase is used by the persistent volume claim binder to signal to another process to reclaim the resource
	VolumeReleased PersistentVolumePhase = "Released"
	// used for PersistentVolumes that failed to be correctly recycled or deleted after being released from a claim
	VolumeFailed PersistentVolumePhase = "Failed"
)

// +enum
type PersistentVolumeClaimPhase string

const (
	// used for PersistentVolumeClaims that are not yet bound
	// 使用未绑定的PersistentVolumeClaims
	ClaimPending PersistentVolumeClaimPhase = "Pending"
	// used for PersistentVolumeClaims that are bound
	// 使用绑定的PersistentVolumeClaims
	ClaimBound PersistentVolumeClaimPhase = "Bound"
	// used for PersistentVolumeClaims that lost their underlying
	// PersistentVolume. The claim was bound to a PersistentVolume and this
	// volume does not exist any longer and all data on it was lost.
	ClaimLost PersistentVolumeClaimPhase = "Lost"
)

// +enum
type HostPathType string

const (
	// For backwards compatible, leave it empty if unset
	// 为了向后兼容，如果未设置，请将其留空
	HostPathUnset HostPathType = ""
	// If nothing exists at the given path, an empty directory will be created there
	// as needed with file mode 0755, having the same group and ownership with Kubelet.
	// 如果给定路径上不存在任何内容，则将在那里创建一个空目录
	HostPathDirectoryOrCreate HostPathType = "DirectoryOrCreate"
	// A directory must exist at the given path
	// 给定路径必须存在一个目录
	HostPathDirectory HostPathType = "Directory"
	// If nothing exists at the given path, an empty file will be created there
	// as needed with file mode 0644, having the same group and ownership with Kubelet.
	// 如果给定路径上不存在任何内容，则将在那里创建一个空文件, 为需要的文件模式 0644, 与 Kubelet 拥有相同的组和所有权。
	HostPathFileOrCreate HostPathType = "FileOrCreate"
	// A file must exist at the given path
	// 给定路径必须存在一个文件
	HostPathFile HostPathType = "File"
	// A UNIX socket must exist at the given path
	// 给定路径必须存在一个 UNIX 套接字
	HostPathSocket HostPathType = "Socket"
	// A character device must exist at the given path
	// 给定路径必须存在一个字符设备
	HostPathCharDev HostPathType = "CharDevice"
	// A block device must exist at the given path
	// 给定路径必须存在一个块设备
	HostPathBlockDev HostPathType = "BlockDevice"
)

// Represents a host path mapped into a pod.
// Host path volumes do not support ownership management or SELinux relabeling.
// 表示映射到pod的主机路径。 主机路径卷不支持所有权管理或SELinux重新标记。
type HostPathVolumeSource struct {
	// path of the directory on the host.
	// If the path is a symlink, it will follow the link to the real path.
	// More info: https://kubernetes.io/docs/concepts/storage/volumes#hostpath
	// 主机上的目录路径
	// 如果路径是符号链接，则将遵循链接到真实路径。
	Path string `json:"path" protobuf:"bytes,1,opt,name=path"`
	// type for HostPath Volume
	// Defaults to ""
	// More info: https://kubernetes.io/docs/concepts/storage/volumes#hostpath
	// +optional
	// HostPath卷的类型
	Type *HostPathType `json:"type,omitempty" protobuf:"bytes,2,opt,name=type"`
}

// Represents an empty directory for a pod.
// Empty directory volumes support ownership management and SELinux relabeling.
// 表示pod的空目录。 空目录卷支持所有权管理和SELinux重新标记。
type EmptyDirVolumeSource struct {
	// medium represents what type of storage medium should back this directory.
	// The default is "" which means to use the node's default medium.
	// Must be an empty string (default) or Memory.
	// More info: https://kubernetes.io/docs/concepts/storage/volumes#emptydir
	// +optional
	// medium表示应该支持此目录的存储介质类型。 默认值是""，这意味着使用节点的默认介质。 必须是一个空字符串（默认值）或内存。
	Medium StorageMedium `json:"medium,omitempty" protobuf:"bytes,1,opt,name=medium,casttype=StorageMedium"`
	// sizeLimit is the total amount of local storage required for this EmptyDir volume.
	// The size limit is also applicable for memory medium.
	// The maximum usage on memory medium EmptyDir would be the minimum value between
	// the SizeLimit specified here and the sum of memory limits of all containers in a pod.
	// The default is nil which means that the limit is undefined.
	// More info: http://kubernetes.io/docs/user-guide/volumes#emptydir
	// +optional
	// sizeLimit是此EmptyDir卷所需的本地存储总量。 尺寸限制也适用于内存介质。 内存介质EmptyDir上的最大使用量将是此处指定的尺寸限制和pod中所有容器的内存限制之和之间的较小值。 默认值为nil，这意味着限制未定义。
	SizeLimit *resource.Quantity `json:"sizeLimit,omitempty" protobuf:"bytes,2,opt,name=sizeLimit"`
}

// Represents a Glusterfs mount that lasts the lifetime of a pod.
// Glusterfs volumes do not support ownership management or SELinux relabeling.
// 表示持续pod生命周期的Glusterfs挂载。 Glusterfs卷不支持所有权管理或SELinux重新标记。
type GlusterfsVolumeSource struct {
	// endpoints is the endpoint name that details Glusterfs topology.
	// More info: https://examples.k8s.io/volumes/glusterfs/README.md#create-a-pod
	// endpoints是详细描述Glusterfs拓扑的端点名称。
	EndpointsName string `json:"endpoints" protobuf:"bytes,1,opt,name=endpoints"`

	// path is the Glusterfs volume path.
	// More info: https://examples.k8s.io/volumes/glusterfs/README.md#create-a-pod
	// path是Glusterfs卷路径。
	Path string `json:"path" protobuf:"bytes,2,opt,name=path"`

	// readOnly here will force the Glusterfs volume to be mounted with read-only permissions.
	// Defaults to false.
	// More info: https://examples.k8s.io/volumes/glusterfs/README.md#create-a-pod
	// +optional
	// readOnly在这里将强制Glusterfs卷以只读权限挂载。 默认值为false。
	ReadOnly bool `json:"readOnly,omitempty" protobuf:"varint,3,opt,name=readOnly"`
}

// Represents a Glusterfs mount that lasts the lifetime of a pod.
// Glusterfs volumes do not support ownership management or SELinux relabeling.
// 表示持续pod生命周期的Glusterfs挂载。 Glusterfs卷不支持所有权管理或SELinux重新标记。
type GlusterfsPersistentVolumeSource struct {
	// endpoints is the endpoint name that details Glusterfs topology.
	// More info: https://examples.k8s.io/volumes/glusterfs/README.md#create-a-pod
	// endpoints是详细描述Glusterfs拓扑的端点名称。
	EndpointsName string `json:"endpoints" protobuf:"bytes,1,opt,name=endpoints"`

	// path is the Glusterfs volume path.
	// More info: https://examples.k8s.io/volumes/glusterfs/README.md#create-a-pod
	// path是Glusterfs卷路径。
	Path string `json:"path" protobuf:"bytes,2,opt,name=path"`

	// readOnly here will force the Glusterfs volume to be mounted with read-only permissions.
	// Defaults to false.
	// More info: https://examples.k8s.io/volumes/glusterfs/README.md#create-a-pod
	// +optional
	// readOnly在这里将强制Glusterfs卷以只读权限挂载。 默认值为false。
	ReadOnly bool `json:"readOnly,omitempty" protobuf:"varint,3,opt,name=readOnly"`

	// endpointsNamespace is the namespace that contains Glusterfs endpoint.
	// If this field is empty, the EndpointNamespace defaults to the same namespace as the bound PVC.
	// More info: https://examples.k8s.io/volumes/glusterfs/README.md#create-a-pod
	// +optional
	// endpointsNamespace 是包含Glusterfs端点的命名空间。 如果此字段为空，则EndpointNamespace默认为与绑定的PVC相同的命名空间。
	EndpointsNamespace *string `json:"endpointsNamespace,omitempty" protobuf:"bytes,4,opt,name=endpointsNamespace"`
}

// Represents a Rados Block Device mount that lasts the lifetime of a pod.
// RBD volumes support ownership management and SELinux relabeling.
// 表示持续pod生命周期的Rados块设备挂载。 RBD卷支持所有权管理和SELinux重新标记。
type RBDVolumeSource struct {
	// monitors is a collection of Ceph monitors.
	// More info: https://examples.k8s.io/volumes/rbd/README.md#how-to-use-it
	// monitors是Ceph监视器的集合。
	CephMonitors []string `json:"monitors" protobuf:"bytes,1,rep,name=monitors"`
	// image is the rados image name.
	// More info: https://examples.k8s.io/volumes/rbd/README.md#how-to-use-it
	// image是rados映像名称。
	RBDImage string `json:"image" protobuf:"bytes,2,opt,name=image"`
	// fsType is the filesystem type of the volume that you want to mount.
	// Tip: Ensure that the filesystem type is supported by the host operating system.
	// Examples: "ext4", "xfs", "ntfs". Implicitly inferred to be "ext4" if unspecified.
	// More info: https://kubernetes.io/docs/concepts/storage/volumes#rbd
	// TODO: how do we prevent errors in the filesystem from compromising the machine
	// +optional
	// fsType是您要挂载的卷的文件系统类型。 提示：确保主机操作系统支持该文件系统类型。 例如：“ext4”，“xfs”，“ntfs”。 如果未指定，则隐式推断为“ext4”。 更多信息：https://kubernetes.io/docs/concepts/storage/volumes#rbd TODO：我们如何防止文件系统中的错误损害机器
	FSType string `json:"fsType,omitempty" protobuf:"bytes,3,opt,name=fsType"`
	// pool is the rados pool name.
	// Default is rbd.
	// More info: https://examples.k8s.io/volumes/rbd/README.md#how-to-use-it
	// +optional
	// pool是rados池名称。 默认值为rbd。
	RBDPool string `json:"pool,omitempty" protobuf:"bytes,4,opt,name=pool"`
	// user is the rados user name.
	// Default is admin.
	// More info: https://examples.k8s.io/volumes/rbd/README.md#how-to-use-it
	// +optional
	// user是rados用户名。 默认值为admin。
	RadosUser string `json:"user,omitempty" protobuf:"bytes,5,opt,name=user"`
	// keyring is the path to key ring for RBDUser.
	// Default is /etc/ceph/keyring.
	// More info: https://examples.k8s.io/volumes/rbd/README.md#how-to-use-it
	// +optional
	// keyring是RBDUser的密钥环路径。 默认值为/etc/ceph/keyring。
	Keyring string `json:"keyring,omitempty" protobuf:"bytes,6,opt,name=keyring"`
	// secretRef is name of the authentication secret for RBDUser. If provided
	// overrides keyring.
	// Default is nil.
	// More info: https://examples.k8s.io/volumes/rbd/README.md#how-to-use-it
	// +optional
	// secretRef是RBDUser的身份验证密钥的名称。 如果提供，则覆盖keyring。 默认值为nil。
	SecretRef *LocalObjectReference `json:"secretRef,omitempty" protobuf:"bytes,7,opt,name=secretRef"`
	// readOnly here will force the ReadOnly setting in VolumeMounts.
	// Defaults to false.
	// More info: https://examples.k8s.io/volumes/rbd/README.md#how-to-use-it
	// +optional
	// readOnly在这里将强制VolumeMounts中的ReadOnly设置。 默认值为false。
	ReadOnly bool `json:"readOnly,omitempty" protobuf:"varint,8,opt,name=readOnly"`
}

// Represents a Rados Block Device mount that lasts the lifetime of a pod.
// RBD volumes support ownership management and SELinux relabeling.
// 表示持续pod生命周期的Rados块设备挂载。 RBD卷支持所有权管理和SELinux重新标记。
type RBDPersistentVolumeSource struct {
	// monitors is a collection of Ceph monitors.
	// More info: https://examples.k8s.io/volumes/rbd/README.md#how-to-use-it
	// monitors是Ceph监视器的集合。
	CephMonitors []string `json:"monitors" protobuf:"bytes,1,rep,name=monitors"`
	// image is the rados image name.
	// More info: https://examples.k8s.io/volumes/rbd/README.md#how-to-use-it
	// image是rados映像名称。
	RBDImage string `json:"image" protobuf:"bytes,2,opt,name=image"`
	// fsType is the filesystem type of the volume that you want to mount.
	// Tip: Ensure that the filesystem type is supported by the host operating system.
	// Examples: "ext4", "xfs", "ntfs". Implicitly inferred to be "ext4" if unspecified.
	// More info: https://kubernetes.io/docs/concepts/storage/volumes#rbd
	// TODO: how do we prevent errors in the filesystem from compromising the machine
	// +optional
	// fsType是您要挂载的卷的文件系统类型。 提示：确保主机操作系统支持该文件系统类型。 例如：“ext4”，“xfs”，“ntfs”。 如果未指定，则隐式推断为“ext4”。 更多信息：https://kubernetes.io/docs/concepts/storage/volumes#rbd TODO：我们如何防止文件系统中的错误损害机器
	FSType string `json:"fsType,omitempty" protobuf:"bytes,3,opt,name=fsType"`
	// pool is the rados pool name.
	// Default is rbd.
	// More info: https://examples.k8s.io/volumes/rbd/README.md#how-to-use-it
	// +optional
	// pool是rados池名称。 默认值为rbd。
	RBDPool string `json:"pool,omitempty" protobuf:"bytes,4,opt,name=pool"`
	// user is the rados user name.
	// Default is admin.
	// More info: https://examples.k8s.io/volumes/rbd/README.md#how-to-use-it
	// +optional
	// user是rados用户名。 默认值为admin。
	RadosUser string `json:"user,omitempty" protobuf:"bytes,5,opt,name=user"`
	// keyring is the path to key ring for RBDUser.
	// Default is /etc/ceph/keyring.
	// More info: https://examples.k8s.io/volumes/rbd/README.md#how-to-use-it
	// +optional
	// keyring是RBDUser的密钥环路径。 默认值为/etc/ceph/keyring。
	Keyring string `json:"keyring,omitempty" protobuf:"bytes,6,opt,name=keyring"`
	// secretRef is name of the authentication secret for RBDUser. If provided
	// overrides keyring.
	// Default is nil.
	// More info: https://examples.k8s.io/volumes/rbd/README.md#how-to-use-it
	// +optional
	// secretRef是RBDUser的身份验证密钥的名称。 如果提供，则覆盖keyring。 默认值为nil。
	SecretRef *SecretReference `json:"secretRef,omitempty" protobuf:"bytes,7,opt,name=secretRef"`
	// readOnly here will force the ReadOnly setting in VolumeMounts.
	// Defaults to false.
	// More info: https://examples.k8s.io/volumes/rbd/README.md#how-to-use-it
	// +optional
	// readOnly在这里将强制VolumeMounts中的ReadOnly设置。 默认值为false。
	ReadOnly bool `json:"readOnly,omitempty" protobuf:"varint,8,opt,name=readOnly"`
}

// Represents a cinder volume resource in Openstack.
// A Cinder volume must exist before mounting to a container.
// The volume must also be in the same region as the kubelet.
// Cinder volumes support ownership management and SELinux relabeling.
// 表示Openstack中的cinder卷资源。 必须在将其挂载到容器之前存在cinder卷。 卷还必须与kubelet位于同一区域。 Cinder卷支持所有权管理和SELinux重新标记。
type CinderVolumeSource struct {
	// volumeID used to identify the volume in cinder.
	// More info: https://examples.k8s.io/mysql-cinder-pd/README.md
	// 用于在cinder中识别卷的volumeID。
	VolumeID string `json:"volumeID" protobuf:"bytes,1,opt,name=volumeID"`
	// fsType is the filesystem type to mount.
	// Must be a filesystem type supported by the host operating system.
	// Examples: "ext4", "xfs", "ntfs". Implicitly inferred to be "ext4" if unspecified.
	// More info: https://examples.k8s.io/mysql-cinder-pd/README.md
	// +optional
	// fsType是要挂载的文件系统类型。 必须是主机操作系统支持的文件系统类型。 例如：“ext4”，“xfs”，“ntfs”。 如果未指定，则隐式推断为“ext4”。
	FSType string `json:"fsType,omitempty" protobuf:"bytes,2,opt,name=fsType"`
	// readOnly defaults to false (read/write). ReadOnly here will force
	// the ReadOnly setting in VolumeMounts.
	// More info: https://examples.k8s.io/mysql-cinder-pd/README.md
	// +optional
	// readOnly默认为false（读/写）。 ReadOnly在这里将强制VolumeMounts中的ReadOnly设置。
	ReadOnly bool `json:"readOnly,omitempty" protobuf:"varint,3,opt,name=readOnly"`
	// secretRef is optional: points to a secret object containing parameters used to connect
	// to OpenStack.
	// +optional
	// secretRef是可选的：指向包含用于连接到OpenStack的参数的密钥对象。
	SecretRef *LocalObjectReference `json:"secretRef,omitempty" protobuf:"bytes,4,opt,name=secretRef"`
}

// Represents a cinder volume resource in Openstack.
// A Cinder volume must exist before mounting to a container.
// The volume must also be in the same region as the kubelet.
// Cinder volumes support ownership management and SELinux relabeling.
// 表示Openstack中的cinder卷资源。 必须在将其挂载到容器之前存在cinder卷。 卷还必须与kubelet位于同一区域。 Cinder卷支持所有权管理和SELinux重新标记。
type CinderPersistentVolumeSource struct {
	// volumeID used to identify the volume in cinder.
	// More info: https://examples.k8s.io/mysql-cinder-pd/README.md
	VolumeID string `json:"volumeID" protobuf:"bytes,1,opt,name=volumeID"`
	// fsType Filesystem type to mount.
	// Must be a filesystem type supported by the host operating system.
	// Examples: "ext4", "xfs", "ntfs". Implicitly inferred to be "ext4" if unspecified.
	// More info: https://examples.k8s.io/mysql-cinder-pd/README.md
	// +optional
	FSType string `json:"fsType,omitempty" protobuf:"bytes,2,opt,name=fsType"`
	// readOnly is Optional: Defaults to false (read/write). ReadOnly here will force
	// the ReadOnly setting in VolumeMounts.
	// More info: https://examples.k8s.io/mysql-cinder-pd/README.md
	// +optional
	ReadOnly bool `json:"readOnly,omitempty" protobuf:"varint,3,opt,name=readOnly"`
	// secretRef is Optional: points to a secret object containing parameters used to connect
	// to OpenStack.
	// +optional
	SecretRef *SecretReference `json:"secretRef,omitempty" protobuf:"bytes,4,opt,name=secretRef"`
}

// Represents a Ceph Filesystem mount that lasts the lifetime of a pod
// Cephfs volumes do not support ownership management or SELinux relabeling.
// 表示持续到pod生命周期的Ceph文件系统挂载。 Cephfs卷不支持所有权管理或SELinux重新标记。
type CephFSVolumeSource struct {
	// monitors is Required: Monitors is a collection of Ceph monitors
	// More info: https://examples.k8s.io/volumes/cephfs/README.md#how-to-use-it
	Monitors []string `json:"monitors" protobuf:"bytes,1,rep,name=monitors"`
	// path is Optional: Used as the mounted root, rather than the full Ceph tree, default is /
	// +optional
	Path string `json:"path,omitempty" protobuf:"bytes,2,opt,name=path"`
	// user is optional: User is the rados user name, default is admin
	// More info: https://examples.k8s.io/volumes/cephfs/README.md#how-to-use-it
	// +optional
	User string `json:"user,omitempty" protobuf:"bytes,3,opt,name=user"`
	// secretFile is Optional: SecretFile is the path to key ring for User, default is /etc/ceph/user.secret
	// More info: https://examples.k8s.io/volumes/cephfs/README.md#how-to-use-it
	// +optional
	SecretFile string `json:"secretFile,omitempty" protobuf:"bytes,4,opt,name=secretFile"`
	// secretRef is Optional: SecretRef is reference to the authentication secret for User, default is empty.
	// More info: https://examples.k8s.io/volumes/cephfs/README.md#how-to-use-it
	// +optional
	SecretRef *LocalObjectReference `json:"secretRef,omitempty" protobuf:"bytes,5,opt,name=secretRef"`
	// readOnly is Optional: Defaults to false (read/write). ReadOnly here will force
	// the ReadOnly setting in VolumeMounts.
	// More info: https://examples.k8s.io/volumes/cephfs/README.md#how-to-use-it
	// +optional
	ReadOnly bool `json:"readOnly,omitempty" protobuf:"varint,6,opt,name=readOnly"`
}

// SecretReference represents a Secret Reference. It has enough information to retrieve secret
// in any namespace
// +structType=atomic
// SecretReference 表示一个 Secret 引用。 它有足够的信息来检索任何命名空间中的 secret
type SecretReference struct {
	// name is unique within a namespace to reference a secret resource.
	// +optional
	Name string `json:"name,omitempty" protobuf:"bytes,1,opt,name=name"`
	// namespace defines the space within which the secret name must be unique.
	// +optional
	Namespace string `json:"namespace,omitempty" protobuf:"bytes,2,opt,name=namespace"`
}

// Represents a Ceph Filesystem mount that lasts the lifetime of a pod
// Cephfs volumes do not support ownership management or SELinux relabeling.
// 表示持续到pod生命周期的Ceph文件系统挂载。 Cephfs卷不支持所有权管理或SELinux重新标记。
type CephFSPersistentVolumeSource struct {
	// monitors is Required: Monitors is a collection of Ceph monitors
	// More info: https://examples.k8s.io/volumes/cephfs/README.md#how-to-use-it
	Monitors []string `json:"monitors" protobuf:"bytes,1,rep,name=monitors"`
	// path is Optional: Used as the mounted root, rather than the full Ceph tree, default is /
	// +optional
	Path string `json:"path,omitempty" protobuf:"bytes,2,opt,name=path"`
	// user is Optional: User is the rados user name, default is admin
	// More info: https://examples.k8s.io/volumes/cephfs/README.md#how-to-use-it
	// +optional
	User string `json:"user,omitempty" protobuf:"bytes,3,opt,name=user"`
	// secretFile is Optional: SecretFile is the path to key ring for User, default is /etc/ceph/user.secret
	// More info: https://examples.k8s.io/volumes/cephfs/README.md#how-to-use-it
	// +optional
	SecretFile string `json:"secretFile,omitempty" protobuf:"bytes,4,opt,name=secretFile"`
	// secretRef is Optional: SecretRef is reference to the authentication secret for User, default is empty.
	// More info: https://examples.k8s.io/volumes/cephfs/README.md#how-to-use-it
	// +optional
	SecretRef *SecretReference `json:"secretRef,omitempty" protobuf:"bytes,5,opt,name=secretRef"`
	// readOnly is Optional: Defaults to false (read/write). ReadOnly here will force
	// the ReadOnly setting in VolumeMounts.
	// More info: https://examples.k8s.io/volumes/cephfs/README.md#how-to-use-it
	// +optional
	ReadOnly bool `json:"readOnly,omitempty" protobuf:"varint,6,opt,name=readOnly"`
}

// Represents a Flocker volume mounted by the Flocker agent.
// One and only one of datasetName and datasetUUID should be set.
// Flocker volumes do not support ownership management or SELinux relabeling.
// 表示由Flocker代理挂载的Flocker卷。 datasetName和datasetUUID中应该只设置一个。 Flocker卷不支持所有权管理或SELinux重新标记。
type FlockerVolumeSource struct {
	// datasetName is Name of the dataset stored as metadata -> name on the dataset for Flocker
	// should be considered as deprecated
	// +optional
	DatasetName string `json:"datasetName,omitempty" protobuf:"bytes,1,opt,name=datasetName"`
	// datasetUUID is the UUID of the dataset. This is unique identifier of a Flocker dataset
	// +optional
	DatasetUUID string `json:"datasetUUID,omitempty" protobuf:"bytes,2,opt,name=datasetUUID"`
}

// StorageMedium defines ways that storage can be allocated to a volume.
// StorageMedium 定义了如何将存储分配给卷。
type StorageMedium string

const (
	StorageMediumDefault         StorageMedium = ""           // use whatever the default is for the node, assume anything we don't explicitly handle is this
	StorageMediumMemory          StorageMedium = "Memory"     // use memory (e.g. tmpfs on linux)
	StorageMediumHugePages       StorageMedium = "HugePages"  // use hugepages
	StorageMediumHugePagesPrefix StorageMedium = "HugePages-" // prefix for full medium notation HugePages-<size>
)

// Protocol defines network protocols supported for things like container ports.
// +enum
// Protocol 定义了支持的网络协议，例如容器端口。
type Protocol string

const (
	// ProtocolTCP is the TCP protocol.
	// TCP协议
	ProtocolTCP Protocol = "TCP"
	// ProtocolUDP is the UDP protocol.
	// UDP协议
	ProtocolUDP Protocol = "UDP"
	// ProtocolSCTP is the SCTP protocol.
	// SCTP协议
	ProtocolSCTP Protocol = "SCTP"
)

// Represents a Persistent Disk resource in Google Compute Engine.
//
// A GCE PD must exist before mounting to a container. The disk must
// also be in the same GCE project and zone as the kubelet. A GCE PD
// can only be mounted as read/write once or read-only many times. GCE
// PDs support ownership management and SELinux relabeling.
// 表示Google Compute Engine中的持久性磁盘资源。
type GCEPersistentDiskVolumeSource struct {
	// pdName is unique name of the PD resource in GCE. Used to identify the disk in GCE.
	// More info: https://kubernetes.io/docs/concepts/storage/volumes#gcepersistentdisk
	PDName string `json:"pdName" protobuf:"bytes,1,opt,name=pdName"`
	// fsType is filesystem type of the volume that you want to mount.
	// Tip: Ensure that the filesystem type is supported by the host operating system.
	// Examples: "ext4", "xfs", "ntfs". Implicitly inferred to be "ext4" if unspecified.
	// More info: https://kubernetes.io/docs/concepts/storage/volumes#gcepersistentdisk
	// TODO: how do we prevent errors in the filesystem from compromising the machine
	// +optional
	FSType string `json:"fsType,omitempty" protobuf:"bytes,2,opt,name=fsType"`
	// partition is the partition in the volume that you want to mount.
	// If omitted, the default is to mount by volume name.
	// Examples: For volume /dev/sda1, you specify the partition as "1".
	// Similarly, the volume partition for /dev/sda is "0" (or you can leave the property empty).
	// More info: https://kubernetes.io/docs/concepts/storage/volumes#gcepersistentdisk
	// +optional
	Partition int32 `json:"partition,omitempty" protobuf:"varint,3,opt,name=partition"`
	// readOnly here will force the ReadOnly setting in VolumeMounts.
	// Defaults to false.
	// More info: https://kubernetes.io/docs/concepts/storage/volumes#gcepersistentdisk
	// +optional
	ReadOnly bool `json:"readOnly,omitempty" protobuf:"varint,4,opt,name=readOnly"`
}

// Represents a Quobyte mount that lasts the lifetime of a pod.
// Quobyte volumes do not support ownership management or SELinux relabeling.
// 表示持续到Pod生命周期的Quobyte挂载。 Quobyte卷不支持所有权管理或SELinux重新标记。
type QuobyteVolumeSource struct {
	// registry represents a single or multiple Quobyte Registry services
	// specified as a string as host:port pair (multiple entries are separated with commas)
	// which acts as the central registry for volumes
	Registry string `json:"registry" protobuf:"bytes,1,opt,name=registry"`

	// volume is a string that references an already created Quobyte volume by name.
	Volume string `json:"volume" protobuf:"bytes,2,opt,name=volume"`

	// readOnly here will force the Quobyte volume to be mounted with read-only permissions.
	// Defaults to false.
	// +optional
	ReadOnly bool `json:"readOnly,omitempty" protobuf:"varint,3,opt,name=readOnly"`

	// user to map volume access to
	// Defaults to serivceaccount user
	// +optional
	User string `json:"user,omitempty" protobuf:"bytes,4,opt,name=user"`

	// group to map volume access to
	// Default is no group
	// +optional
	Group string `json:"group,omitempty" protobuf:"bytes,5,opt,name=group"`

	// tenant owning the given Quobyte volume in the Backend
	// Used with dynamically provisioned Quobyte volumes, value is set by the plugin
	// +optional
	Tenant string `json:"tenant,omitempty" protobuf:"bytes,6,opt,name=tenant"`
}

// FlexPersistentVolumeSource represents a generic persistent volume resource that is
// provisioned/attached using an exec based plugin.
// FlexPersistentVolumeSource表示使用基于exec的插件提供/附加的通用持久卷资源。
type FlexPersistentVolumeSource struct {
	// driver is the name of the driver to use for this volume.
	// driver是用于此卷的驱动程序的名称。
	Driver string `json:"driver" protobuf:"bytes,1,opt,name=driver"`
	// fsType is the Filesystem type to mount.
	// Must be a filesystem type supported by the host operating system.
	// Ex. "ext4", "xfs", "ntfs". The default filesystem depends on FlexVolume script.
	// +optional
	// fsType是要挂载的文件系统类型。 必须是主机操作系统支持的文件系统类型。 例如“ext4”，“xfs”，“ntfs”。 默认文件系统取决于FlexVolume脚本。
	FSType string `json:"fsType,omitempty" protobuf:"bytes,2,opt,name=fsType"`
	// secretRef is Optional: SecretRef is reference to the secret object containing
	// sensitive information to pass to the plugin scripts. This may be
	// empty if no secret object is specified. If the secret object
	// contains more than one secret, all secrets are passed to the plugin
	// scripts.
	// +optional
	// secretRef是可选的：SecretRef是对包含要传递给插件脚本的敏感信息的密钥对象的引用。 如果未指定密钥对象，则可以为空。 如果密钥对象包含多个密钥，则所有密钥都传递给插件脚本。
	SecretRef *SecretReference `json:"secretRef,omitempty" protobuf:"bytes,3,opt,name=secretRef"`
	// readOnly is Optional: defaults to false (read/write). ReadOnly here will force
	// the ReadOnly setting in VolumeMounts.
	// +optional
	// readOnly是可选的：默认为false（读/写）。 ReadOnly将强制在VolumeMounts中设置为只读。
	ReadOnly bool `json:"readOnly,omitempty" protobuf:"varint,4,opt,name=readOnly"`
	// options is Optional: this field holds extra command options if any.
	// +optional
	// options是可选的：此字段保存任何额外的命令选项。
	Options map[string]string `json:"options,omitempty" protobuf:"bytes,5,rep,name=options"`
}

// FlexVolume represents a generic volume resource that is
// provisioned/attached using an exec based plugin.
// FlexVolume表示使用基于exec的插件提供/附加的通用卷资源。
type FlexVolumeSource struct {
	// driver is the name of the driver to use for this volume.
	// driver是用于此卷的驱动程序的名称。
	Driver string `json:"driver" protobuf:"bytes,1,opt,name=driver"`
	// fsType is the filesystem type to mount.
	// Must be a filesystem type supported by the host operating system.
	// Ex. "ext4", "xfs", "ntfs". The default filesystem depends on FlexVolume script.
	// +optional
	// fsType是要挂载的文件系统类型。 必须是主机操作系统支持的文件系统类型。 例如“ext4”，“xfs”，“ntfs”。 默认文件系统取决于FlexVolume脚本。
	FSType string `json:"fsType,omitempty" protobuf:"bytes,2,opt,name=fsType"`
	// secretRef is Optional: secretRef is reference to the secret object containing
	// sensitive information to pass to the plugin scripts. This may be
	// empty if no secret object is specified. If the secret object
	// contains more than one secret, all secrets are passed to the plugin
	// scripts.
	// +optional
	// secretRef是可选的：secretRef是对包含要传递给插件脚本的敏感信息的密钥对象的引用。 如果未指定密钥对象，则可以为空。 如果密钥对象包含多个密钥，则所有密钥都传递给插件脚本。
	SecretRef *LocalObjectReference `json:"secretRef,omitempty" protobuf:"bytes,3,opt,name=secretRef"`
	// readOnly is Optional: defaults to false (read/write). ReadOnly here will force
	// the ReadOnly setting in VolumeMounts.
	// +optional
	// readOnly是可选的：默认为false（读/写）。 ReadOnly将强制在VolumeMounts中设置为只读。
	ReadOnly bool `json:"readOnly,omitempty" protobuf:"varint,4,opt,name=readOnly"`
	// options is Optional: this field holds extra command options if any.
	// +optional
	// options是可选的：此字段保存任何额外的命令选项。
	Options map[string]string `json:"options,omitempty" protobuf:"bytes,5,rep,name=options"`
}

// Represents a Persistent Disk resource in AWS.
//
// An AWS EBS disk must exist before mounting to a container. The disk
// must also be in the same AWS zone as the kubelet. An AWS EBS disk
// can only be mounted as read/write once. AWS EBS volumes support
// ownership management and SELinux relabeling.
// 表示AWS中的持久磁盘资源。
type AWSElasticBlockStoreVolumeSource struct {
	// volumeID is unique ID of the persistent disk resource in AWS (Amazon EBS volume).
	// More info: https://kubernetes.io/docs/concepts/storage/volumes#awselasticblockstore
	VolumeID string `json:"volumeID" protobuf:"bytes,1,opt,name=volumeID"`
	// fsType is the filesystem type of the volume that you want to mount.
	// Tip: Ensure that the filesystem type is supported by the host operating system.
	// Examples: "ext4", "xfs", "ntfs". Implicitly inferred to be "ext4" if unspecified.
	// More info: https://kubernetes.io/docs/concepts/storage/volumes#awselasticblockstore
	// TODO: how do we prevent errors in the filesystem from compromising the machine
	// +optional
	FSType string `json:"fsType,omitempty" protobuf:"bytes,2,opt,name=fsType"`
	// partition is the partition in the volume that you want to mount.
	// If omitted, the default is to mount by volume name.
	// Examples: For volume /dev/sda1, you specify the partition as "1".
	// Similarly, the volume partition for /dev/sda is "0" (or you can leave the property empty).
	// +optional
	Partition int32 `json:"partition,omitempty" protobuf:"varint,3,opt,name=partition"`
	// readOnly value true will force the readOnly setting in VolumeMounts.
	// More info: https://kubernetes.io/docs/concepts/storage/volumes#awselasticblockstore
	// +optional
	ReadOnly bool `json:"readOnly,omitempty" protobuf:"varint,4,opt,name=readOnly"`
}

// Represents a volume that is populated with the contents of a git repository.
// Git repo volumes do not support ownership management.
// Git repo volumes support SELinux relabeling.
//
// DEPRECATED: GitRepo is deprecated. To provision a container with a git repo, mount an
// EmptyDir into an InitContainer that clones the repo using git, then mount the EmptyDir
// into the Pod's container.
// 表示使用git存储库的内容填充的卷。 Git repo卷不支持所有权管理。 Git repo卷支持SELinux relabeling。
// DEPRECATED：GitRepo已弃用。 要使用git存储库提供容器，请将EmptyDir挂载到克隆存储库的InitContainer中，然后将EmptyDir挂载到Pod的容器中。
type GitRepoVolumeSource struct {
	// repository is the URL
	Repository string `json:"repository" protobuf:"bytes,1,opt,name=repository"`
	// revision is the commit hash for the specified revision.
	// +optional
	Revision string `json:"revision,omitempty" protobuf:"bytes,2,opt,name=revision"`
	// directory is the target directory name.
	// Must not contain or start with '..'.  If '.' is supplied, the volume directory will be the
	// git repository.  Otherwise, if specified, the volume will contain the git repository in
	// the subdirectory with the given name.
	// +optional
	Directory string `json:"directory,omitempty" protobuf:"bytes,3,opt,name=directory"`
}

// Adapts a Secret into a volume.
//
// The contents of the target Secret's Data field will be presented in a volume
// as files using the keys in the Data field as the file names.
// Secret volumes support ownership management and SELinux relabeling.
// 表示将Secret适配为卷。
// 目标Secret的Data字段的内容将作为文件呈现在卷中，使用Data字段中的键作为文件名。
type SecretVolumeSource struct {
	// secretName is the name of the secret in the pod's namespace to use.
	// More info: https://kubernetes.io/docs/concepts/storage/volumes#secret
	// +optional
	// secretName是pod命名空间中要使用的secret的名称。
	SecretName string `json:"secretName,omitempty" protobuf:"bytes,1,opt,name=secretName"`
	// items If unspecified, each key-value pair in the Data field of the referenced
	// Secret will be projected into the volume as a file whose name is the
	// key and content is the value. If specified, the listed keys will be
	// projected into the specified paths, and unlisted keys will not be
	// present. If a key is specified which is not present in the Secret,
	// the volume setup will error unless it is marked optional. Paths must be
	// relative and may not contain the '..' path or start with '..'.
	// +optional
	// items如果未指定，则将引用的Secret的Data字段中的每个键值对投影到卷中，作为名称为键，内容为值的文件。
	// 如果指定，则列出的键将投影到指定的路径中，未列出的键将不会存在。
	// 如果指定了不存在于Secret中的键，则除非将其标记为可选，否则卷设置将出错。
	// 路径必须是相对路径，不能包含“..”路径或以“..”开头。
	Items []KeyToPath `json:"items,omitempty" protobuf:"bytes,2,rep,name=items"`
	// defaultMode is Optional: mode bits used to set permissions on created files by default.
	// Must be an octal value between 0000 and 0777 or a decimal value between 0 and 511.
	// YAML accepts both octal and decimal values, JSON requires decimal values
	// for mode bits. Defaults to 0644.
	// Directories within the path are not affected by this setting.
	// This might be in conflict with other options that affect the file
	// mode, like fsGroup, and the result can be other mode bits set.
	// +optional
	// defaultMode是可选的：用于默认情况下设置创建文件的权限位的模式位。
	// 必须是0000和0777之间的八进制值或0和511之间的十进制值。
	// YAML接受八进制和十进制值，JSON要求模式位的十进制值。
	// 默认为0644。 路径中的目录不受此设置的影响。
	// 这可能与影响文件模式的其他选项（如fsGroup）冲突，并且结果可能是其他模式位设置。
	DefaultMode *int32 `json:"defaultMode,omitempty" protobuf:"bytes,3,opt,name=defaultMode"`
	// optional field specify whether the Secret or its keys must be defined
	// +optional
	// optional字段指定Secret或其键是否必须定义
	Optional *bool `json:"optional,omitempty" protobuf:"varint,4,opt,name=optional"`
}

const (
	SecretVolumeSourceDefaultMode int32 = 0644
)

// Adapts a secret into a projected volume.
//
// The contents of the target Secret's Data field will be presented in a
// projected volume as files using the keys in the Data field as the file names.
// Note that this is identical to a secret volume source without the default
// mode.
// 表示将Secret适配为投影卷。
// 目标Secret的Data字段的内容将作为文件呈现在投影卷中，使用Data字段中的键作为文件名。
// 请注意，这与没有默认模式的secret卷源相同。
type SecretProjection struct {
	LocalObjectReference `json:",inline" protobuf:"bytes,1,opt,name=localObjectReference"`
	// items if unspecified, each key-value pair in the Data field of the referenced
	// Secret will be projected into the volume as a file whose name is the
	// key and content is the value. If specified, the listed keys will be
	// projected into the specified paths, and unlisted keys will not be
	// present. If a key is specified which is not present in the Secret,
	// the volume setup will error unless it is marked optional. Paths must be
	// relative and may not contain the '..' path or start with '..'.
	// +optional
	// items如果未指定，则将引用的Secret的Data字段中的每个键值对投影到卷中，作为名称为键，内容为值的文件。
	// 如果指定，则列出的键将投影到指定的路径中，未列出的键将不会存在。
	// 如果指定了不存在于Secret中的键，则除非将其标记为可选，否则卷设置将出错。
	// 路径必须是相对路径，不能包含“..”路径或以“..”开头。
	Items []KeyToPath `json:"items,omitempty" protobuf:"bytes,2,rep,name=items"`
	// optional field specify whether the Secret or its key must be defined
	// +optional
	// optional字段指定Secret或其键是否必须定义
	Optional *bool `json:"optional,omitempty" protobuf:"varint,4,opt,name=optional"`
}

// Represents an NFS mount that lasts the lifetime of a pod.
// NFS volumes do not support ownership management or SELinux relabeling.
// 表示持续到pod生命周期的NFS挂载。
// NFS卷不支持所有权管理或SELinux重新标记。
type NFSVolumeSource struct {
	// server is the hostname or IP address of the NFS server.
	// More info: https://kubernetes.io/docs/concepts/storage/volumes#nfs
	// server是NFS服务器的主机名或IP地址。
	Server string `json:"server" protobuf:"bytes,1,opt,name=server"`

	// path that is exported by the NFS server.
	// More info: https://kubernetes.io/docs/concepts/storage/volumes#nfs
	// path是NFS服务器导出的路径。
	Path string `json:"path" protobuf:"bytes,2,opt,name=path"`

	// readOnly here will force the NFS export to be mounted with read-only permissions.
	// Defaults to false.
	// More info: https://kubernetes.io/docs/concepts/storage/volumes#nfs
	// +optional
	// readOnly这里将强制NFS导出以只读权限挂载。
	ReadOnly bool `json:"readOnly,omitempty" protobuf:"varint,3,opt,name=readOnly"`
}

// Represents an ISCSI disk.
// ISCSI volumes can only be mounted as read/write once.
// ISCSI volumes support ownership management and SELinux relabeling.
// 表示ISCSI磁盘。
// ISCSI卷只能挂载为一次读/写。
// ISCSI卷支持所有权管理和SELinux重新标记。
type ISCSIVolumeSource struct {
	// targetPortal is iSCSI Target Portal. The Portal is either an IP or ip_addr:port if the port
	// is other than default (typically TCP ports 860 and 3260).
	// targetPortal是iSCSI目标门户。 门户是IP或ip_addr：port（如果端口不是默认值（通常是TCP端口860和3260））。
	TargetPortal string `json:"targetPortal" protobuf:"bytes,1,opt,name=targetPortal"`
	// iqn is the target iSCSI Qualified Name.
	// iqn是目标iSCSI限定名称。
	IQN string `json:"iqn" protobuf:"bytes,2,opt,name=iqn"`
	// lun represents iSCSI Target Lun number.
	// lun表示iSCSI目标Lun编号。
	Lun int32 `json:"lun" protobuf:"varint,3,opt,name=lun"`
	// iscsiInterface is the interface Name that uses an iSCSI transport.
	// Defaults to 'default' (tcp).
	// +optional
	// iscsiInterface是使用iSCSI传输的接口名称。
	ISCSIInterface string `json:"iscsiInterface,omitempty" protobuf:"bytes,4,opt,name=iscsiInterface"`
	// fsType is the filesystem type of the volume that you want to mount.
	// Tip: Ensure that the filesystem type is supported by the host operating system.
	// Examples: "ext4", "xfs", "ntfs". Implicitly inferred to be "ext4" if unspecified.
	// More info: https://kubernetes.io/docs/concepts/storage/volumes#iscsi
	// TODO: how do we prevent errors in the filesystem from compromising the machine
	// +optional
	// fsType是您要挂载的卷的文件系统类型。
	// 提示：确保主机操作系统支持文件系统类型。
	FSType string `json:"fsType,omitempty" protobuf:"bytes,5,opt,name=fsType"`
	// readOnly here will force the ReadOnly setting in VolumeMounts.
	// Defaults to false.
	// +optional
	// readOnly这里将强制VolumeMounts中的ReadOnly设置。
	ReadOnly bool `json:"readOnly,omitempty" protobuf:"varint,6,opt,name=readOnly"`
	// portals is the iSCSI Target Portal List. The portal is either an IP or ip_addr:port if the port
	// is other than default (typically TCP ports 860 and 3260).
	// +optional
	// portals是iSCSI目标门户列表。 门户是IP或ip_addr：port（如果端口不是默认值（通常是TCP端口860和3260））。
	Portals []string `json:"portals,omitempty" protobuf:"bytes,7,opt,name=portals"`
	// chapAuthDiscovery defines whether support iSCSI Discovery CHAP authentication
	// +optional
	// chapAuthDiscovery定义是否支持iSCSI Discovery CHAP身份验证
	DiscoveryCHAPAuth bool `json:"chapAuthDiscovery,omitempty" protobuf:"varint,8,opt,name=chapAuthDiscovery"`
	// chapAuthSession defines whether support iSCSI Session CHAP authentication
	// +optional
	// chapAuthSession定义是否支持iSCSI Session CHAP身份验证
	SessionCHAPAuth bool `json:"chapAuthSession,omitempty" protobuf:"varint,11,opt,name=chapAuthSession"`
	// secretRef is the CHAP Secret for iSCSI target and initiator authentication
	// +optional
	// secretRef是iSCSI目标和启动器身份验证的CHAP Secret
	SecretRef *LocalObjectReference `json:"secretRef,omitempty" protobuf:"bytes,10,opt,name=secretRef"`
	// initiatorName is the custom iSCSI Initiator Name.
	// If initiatorName is specified with iscsiInterface simultaneously, new iSCSI interface
	// <target portal>:<volume name> will be created for the connection.
	// +optional
	// initiatorName是自定义iSCSI启动器名称。
	// 如果同时指定了initiatorName和iscsiInterface，则将为连接创建新的iSCSI接口
	InitiatorName *string `json:"initiatorName,omitempty" protobuf:"bytes,12,opt,name=initiatorName"`
}

// ISCSIPersistentVolumeSource represents an ISCSI disk.
// ISCSI volumes can only be mounted as read/write once.
// ISCSI volumes support ownership management and SELinux relabeling.
// ISCSIPersistentVolumeSource表示ISCSI磁盘。
// ISCSI卷只能挂载为一次读/写。
// ISCSI卷支持所有权管理和SELinux重新标记。
type ISCSIPersistentVolumeSource struct {
	// targetPortal is iSCSI Target Portal. The Portal is either an IP or ip_addr:port if the port
	// is other than default (typically TCP ports 860 and 3260).
	// targetPortal是iSCSI目标门户。 门户是IP或ip_addr：port（如果端口不是默认值（通常是TCP端口860和3260））。
	TargetPortal string `json:"targetPortal" protobuf:"bytes,1,opt,name=targetPortal"`
	// iqn is Target iSCSI Qualified Name.
	// iqn是目标iSCSI限定名称。
	IQN string `json:"iqn" protobuf:"bytes,2,opt,name=iqn"`
	// lun is iSCSI Target Lun number.
	// lun是iSCSI目标Lun编号。
	Lun int32 `json:"lun" protobuf:"varint,3,opt,name=lun"`
	// iscsiInterface is the interface Name that uses an iSCSI transport.
	// Defaults to 'default' (tcp).
	// +optional
	// iscsiInterface是使用iSCSI传输的接口名称。
	ISCSIInterface string `json:"iscsiInterface,omitempty" protobuf:"bytes,4,opt,name=iscsiInterface"`
	// fsType is the filesystem type of the volume that you want to mount.
	// Tip: Ensure that the filesystem type is supported by the host operating system.
	// Examples: "ext4", "xfs", "ntfs". Implicitly inferred to be "ext4" if unspecified.
	// More info: https://kubernetes.io/docs/concepts/storage/volumes#iscsi
	// TODO: how do we prevent errors in the filesystem from compromising the machine
	// +optional
	// fsType是您要挂载的卷的文件系统类型。
	FSType string `json:"fsType,omitempty" protobuf:"bytes,5,opt,name=fsType"`
	// readOnly here will force the ReadOnly setting in VolumeMounts.
	// Defaults to false.
	// +optional
	// readOnly这里将强制VolumeMounts中的ReadOnly设置。
	ReadOnly bool `json:"readOnly,omitempty" protobuf:"varint,6,opt,name=readOnly"`
	// portals is the iSCSI Target Portal List. The Portal is either an IP or ip_addr:port if the port
	// is other than default (typically TCP ports 860 and 3260).
	// +optional
	// portals是iSCSI目标门户列表。 门户是IP或ip_addr：port（如果端口不是默认值（通常是TCP端口860和3260））。
	Portals []string `json:"portals,omitempty" protobuf:"bytes,7,opt,name=portals"`
	// chapAuthDiscovery defines whether support iSCSI Discovery CHAP authentication
	// +optional
	// chapAuthDiscovery定义是否支持iSCSI Discovery CHAP身份验证
	DiscoveryCHAPAuth bool `json:"chapAuthDiscovery,omitempty" protobuf:"varint,8,opt,name=chapAuthDiscovery"`
	// chapAuthSession defines whether support iSCSI Session CHAP authentication
	// +optional
	// chapAuthSession定义是否支持iSCSI Session CHAP身份验证
	SessionCHAPAuth bool `json:"chapAuthSession,omitempty" protobuf:"varint,11,opt,name=chapAuthSession"`
	// secretRef is the CHAP Secret for iSCSI target and initiator authentication
	// +optional
	// secretRef是iSCSI目标和启动器身份验证的CHAP Secret
	SecretRef *SecretReference `json:"secretRef,omitempty" protobuf:"bytes,10,opt,name=secretRef"`
	// initiatorName is the custom iSCSI Initiator Name.
	// If initiatorName is specified with iscsiInterface simultaneously, new iSCSI interface
	// <target portal>:<volume name> will be created for the connection.
	// +optional
	// initiatorName是自定义iSCSI启动器名称。
	// 如果同时指定了initiatorName和iscsiInterface，则将为连接创建新的iSCSI接口
	// <target portal>:<volume name>将为连接创建新的iSCSI接口
	InitiatorName *string `json:"initiatorName,omitempty" protobuf:"bytes,12,opt,name=initiatorName"`
}

// Represents a Fibre Channel volume.
// Fibre Channel volumes can only be mounted as read/write once.
// Fibre Channel volumes support ownership management and SELinux relabeling.
// 表示光纤通道卷。
// 光纤通道卷只能挂载为一次读/写。
// 光纤通道卷支持所有权管理和SELinux重新标记。
type FCVolumeSource struct {
	// targetWWNs is Optional: FC target worldwide names (WWNs)
	// +optional
	// targetWWNs是可选的：FC目标全球名称（WWNs）
	TargetWWNs []string `json:"targetWWNs,omitempty" protobuf:"bytes,1,rep,name=targetWWNs"`
	// lun is Optional: FC target lun number
	// +optional
	// lun是可选的：FC目标lun编号
	Lun *int32 `json:"lun,omitempty" protobuf:"varint,2,opt,name=lun"`
	// fsType is the filesystem type to mount.
	// Must be a filesystem type supported by the host operating system.
	// Ex. "ext4", "xfs", "ntfs". Implicitly inferred to be "ext4" if unspecified.
	// TODO: how do we prevent errors in the filesystem from compromising the machine
	// +optional
	// fsType是要挂载的文件系统类型。
	FSType string `json:"fsType,omitempty" protobuf:"bytes,3,opt,name=fsType"`
	// readOnly is Optional: Defaults to false (read/write). ReadOnly here will force
	// the ReadOnly setting in VolumeMounts.
	// +optional
	// readOnly是可选的：默认为false（读/写）。 ReadOnly这里将强制
	ReadOnly bool `json:"readOnly,omitempty" protobuf:"varint,4,opt,name=readOnly"`
	// wwids Optional: FC volume world wide identifiers (wwids)
	// Either wwids or combination of targetWWNs and lun must be set, but not both simultaneously.
	// +optional
	// wwids可选：FC卷世界范围标识符（wwids）
	// targetWWNs和lun的组合或wwids必须设置，但不能同时设置。
	WWIDs []string `json:"wwids,omitempty" protobuf:"bytes,5,rep,name=wwids"`
}

// AzureFile represents an Azure File Service mount on the host and bind mount to the pod.
type AzureFileVolumeSource struct {
	// secretName is the  name of secret that contains Azure Storage Account Name and Key
	SecretName string `json:"secretName" protobuf:"bytes,1,opt,name=secretName"`
	// shareName is the azure share Name
	ShareName string `json:"shareName" protobuf:"bytes,2,opt,name=shareName"`
	// readOnly defaults to false (read/write). ReadOnly here will force
	// the ReadOnly setting in VolumeMounts.
	// +optional
	ReadOnly bool `json:"readOnly,omitempty" protobuf:"varint,3,opt,name=readOnly"`
}

// AzureFile represents an Azure File Service mount on the host and bind mount to the pod.
type AzureFilePersistentVolumeSource struct {
	// secretName is the name of secret that contains Azure Storage Account Name and Key
	SecretName string `json:"secretName" protobuf:"bytes,1,opt,name=secretName"`
	// shareName is the azure Share Name
	ShareName string `json:"shareName" protobuf:"bytes,2,opt,name=shareName"`
	// readOnly defaults to false (read/write). ReadOnly here will force
	// the ReadOnly setting in VolumeMounts.
	// +optional
	ReadOnly bool `json:"readOnly,omitempty" protobuf:"varint,3,opt,name=readOnly"`
	// secretNamespace is the namespace of the secret that contains Azure Storage Account Name and Key
	// default is the same as the Pod
	// +optional
	SecretNamespace *string `json:"secretNamespace" protobuf:"bytes,4,opt,name=secretNamespace"`
}

// Represents a vSphere volume resource.
// 表示vSphere卷资源。
type VsphereVirtualDiskVolumeSource struct {
	// volumePath is the path that identifies vSphere volume vmdk
	VolumePath string `json:"volumePath" protobuf:"bytes,1,opt,name=volumePath"`
	// fsType is filesystem type to mount.
	// Must be a filesystem type supported by the host operating system.
	// Ex. "ext4", "xfs", "ntfs". Implicitly inferred to be "ext4" if unspecified.
	// +optional
	FSType string `json:"fsType,omitempty" protobuf:"bytes,2,opt,name=fsType"`
	// storagePolicyName is the storage Policy Based Management (SPBM) profile name.
	// +optional
	StoragePolicyName string `json:"storagePolicyName,omitempty" protobuf:"bytes,3,opt,name=storagePolicyName"`
	// storagePolicyID is the storage Policy Based Management (SPBM) profile ID associated with the StoragePolicyName.
	// +optional
	StoragePolicyID string `json:"storagePolicyID,omitempty" protobuf:"bytes,4,opt,name=storagePolicyID"`
}

// Represents a Photon Controller persistent disk resource.
// 表示Photon Controller持久性磁盘资源。
type PhotonPersistentDiskVolumeSource struct {
	// pdID is the ID that identifies Photon Controller persistent disk
	PdID string `json:"pdID" protobuf:"bytes,1,opt,name=pdID"`
	// fsType is the filesystem type to mount.
	// Must be a filesystem type supported by the host operating system.
	// Ex. "ext4", "xfs", "ntfs". Implicitly inferred to be "ext4" if unspecified.
	FSType string `json:"fsType,omitempty" protobuf:"bytes,2,opt,name=fsType"`
}

// +enum
type AzureDataDiskCachingMode string

// +enum
type AzureDataDiskKind string

const (
	AzureDataDiskCachingNone      AzureDataDiskCachingMode = "None"
	AzureDataDiskCachingReadOnly  AzureDataDiskCachingMode = "ReadOnly"
	AzureDataDiskCachingReadWrite AzureDataDiskCachingMode = "ReadWrite"

	AzureSharedBlobDisk    AzureDataDiskKind = "Shared"
	AzureDedicatedBlobDisk AzureDataDiskKind = "Dedicated"
	AzureManagedDisk       AzureDataDiskKind = "Managed"
)

// AzureDisk represents an Azure Data Disk mount on the host and bind mount to the pod.
type AzureDiskVolumeSource struct {
	// diskName is the Name of the data disk in the blob storage
	DiskName string `json:"diskName" protobuf:"bytes,1,opt,name=diskName"`
	// diskURI is the URI of data disk in the blob storage
	DataDiskURI string `json:"diskURI" protobuf:"bytes,2,opt,name=diskURI"`
	// cachingMode is the Host Caching mode: None, Read Only, Read Write.
	// +optional
	CachingMode *AzureDataDiskCachingMode `json:"cachingMode,omitempty" protobuf:"bytes,3,opt,name=cachingMode,casttype=AzureDataDiskCachingMode"`
	// fsType is Filesystem type to mount.
	// Must be a filesystem type supported by the host operating system.
	// Ex. "ext4", "xfs", "ntfs". Implicitly inferred to be "ext4" if unspecified.
	// +optional
	FSType *string `json:"fsType,omitempty" protobuf:"bytes,4,opt,name=fsType"`
	// readOnly Defaults to false (read/write). ReadOnly here will force
	// the ReadOnly setting in VolumeMounts.
	// +optional
	ReadOnly *bool `json:"readOnly,omitempty" protobuf:"varint,5,opt,name=readOnly"`
	// kind expected values are Shared: multiple blob disks per storage account  Dedicated: single blob disk per storage account  Managed: azure managed data disk (only in managed availability set). defaults to shared
	Kind *AzureDataDiskKind `json:"kind,omitempty" protobuf:"bytes,6,opt,name=kind,casttype=AzureDataDiskKind"`
}

// PortworxVolumeSource represents a Portworx volume resource.
type PortworxVolumeSource struct {
	// volumeID uniquely identifies a Portworx volume
	VolumeID string `json:"volumeID" protobuf:"bytes,1,opt,name=volumeID"`
	// fSType represents the filesystem type to mount
	// Must be a filesystem type supported by the host operating system.
	// Ex. "ext4", "xfs". Implicitly inferred to be "ext4" if unspecified.
	FSType string `json:"fsType,omitempty" protobuf:"bytes,2,opt,name=fsType"`
	// readOnly defaults to false (read/write). ReadOnly here will force
	// the ReadOnly setting in VolumeMounts.
	// +optional
	ReadOnly bool `json:"readOnly,omitempty" protobuf:"varint,3,opt,name=readOnly"`
}

// ScaleIOVolumeSource represents a persistent ScaleIO volume
type ScaleIOVolumeSource struct {
	// gateway is the host address of the ScaleIO API Gateway.
	Gateway string `json:"gateway" protobuf:"bytes,1,opt,name=gateway"`
	// system is the name of the storage system as configured in ScaleIO.
	System string `json:"system" protobuf:"bytes,2,opt,name=system"`
	// secretRef references to the secret for ScaleIO user and other
	// sensitive information. If this is not provided, Login operation will fail.
	SecretRef *LocalObjectReference `json:"secretRef" protobuf:"bytes,3,opt,name=secretRef"`
	// sslEnabled Flag enable/disable SSL communication with Gateway, default false
	// +optional
	SSLEnabled bool `json:"sslEnabled,omitempty" protobuf:"varint,4,opt,name=sslEnabled"`
	// protectionDomain is the name of the ScaleIO Protection Domain for the configured storage.
	// +optional
	ProtectionDomain string `json:"protectionDomain,omitempty" protobuf:"bytes,5,opt,name=protectionDomain"`
	// storagePool is the ScaleIO Storage Pool associated with the protection domain.
	// +optional
	StoragePool string `json:"storagePool,omitempty" protobuf:"bytes,6,opt,name=storagePool"`
	// storageMode indicates whether the storage for a volume should be ThickProvisioned or ThinProvisioned.
	// Default is ThinProvisioned.
	// +optional
	StorageMode string `json:"storageMode,omitempty" protobuf:"bytes,7,opt,name=storageMode"`
	// volumeName is the name of a volume already created in the ScaleIO system
	// that is associated with this volume source.
	VolumeName string `json:"volumeName,omitempty" protobuf:"bytes,8,opt,name=volumeName"`
	// fsType is the filesystem type to mount.
	// Must be a filesystem type supported by the host operating system.
	// Ex. "ext4", "xfs", "ntfs".
	// Default is "xfs".
	// +optional
	FSType string `json:"fsType,omitempty" protobuf:"bytes,9,opt,name=fsType"`
	// readOnly Defaults to false (read/write). ReadOnly here will force
	// the ReadOnly setting in VolumeMounts.
	// +optional
	ReadOnly bool `json:"readOnly,omitempty" protobuf:"varint,10,opt,name=readOnly"`
}

// ScaleIOPersistentVolumeSource represents a persistent ScaleIO volume
type ScaleIOPersistentVolumeSource struct {
	// gateway is the host address of the ScaleIO API Gateway.
	Gateway string `json:"gateway" protobuf:"bytes,1,opt,name=gateway"`
	// system is the name of the storage system as configured in ScaleIO.
	System string `json:"system" protobuf:"bytes,2,opt,name=system"`
	// secretRef references to the secret for ScaleIO user and other
	// sensitive information. If this is not provided, Login operation will fail.
	SecretRef *SecretReference `json:"secretRef" protobuf:"bytes,3,opt,name=secretRef"`
	// sslEnabled is the flag to enable/disable SSL communication with Gateway, default false
	// +optional
	SSLEnabled bool `json:"sslEnabled,omitempty" protobuf:"varint,4,opt,name=sslEnabled"`
	// protectionDomain is the name of the ScaleIO Protection Domain for the configured storage.
	// +optional
	ProtectionDomain string `json:"protectionDomain,omitempty" protobuf:"bytes,5,opt,name=protectionDomain"`
	// storagePool is the ScaleIO Storage Pool associated with the protection domain.
	// +optional
	StoragePool string `json:"storagePool,omitempty" protobuf:"bytes,6,opt,name=storagePool"`
	// storageMode indicates whether the storage for a volume should be ThickProvisioned or ThinProvisioned.
	// Default is ThinProvisioned.
	// +optional
	StorageMode string `json:"storageMode,omitempty" protobuf:"bytes,7,opt,name=storageMode"`
	// volumeName is the name of a volume already created in the ScaleIO system
	// that is associated with this volume source.
	VolumeName string `json:"volumeName,omitempty" protobuf:"bytes,8,opt,name=volumeName"`
	// fsType is the filesystem type to mount.
	// Must be a filesystem type supported by the host operating system.
	// Ex. "ext4", "xfs", "ntfs".
	// Default is "xfs"
	// +optional
	FSType string `json:"fsType,omitempty" protobuf:"bytes,9,opt,name=fsType"`
	// readOnly defaults to false (read/write). ReadOnly here will force
	// the ReadOnly setting in VolumeMounts.
	// +optional
	ReadOnly bool `json:"readOnly,omitempty" protobuf:"varint,10,opt,name=readOnly"`
}

// Represents a StorageOS persistent volume resource.
type StorageOSVolumeSource struct {
	// volumeName is the human-readable name of the StorageOS volume.  Volume
	// names are only unique within a namespace.
	VolumeName string `json:"volumeName,omitempty" protobuf:"bytes,1,opt,name=volumeName"`
	// volumeNamespace specifies the scope of the volume within StorageOS.  If no
	// namespace is specified then the Pod's namespace will be used.  This allows the
	// Kubernetes name scoping to be mirrored within StorageOS for tighter integration.
	// Set VolumeName to any name to override the default behaviour.
	// Set to "default" if you are not using namespaces within StorageOS.
	// Namespaces that do not pre-exist within StorageOS will be created.
	// +optional
	VolumeNamespace string `json:"volumeNamespace,omitempty" protobuf:"bytes,2,opt,name=volumeNamespace"`
	// fsType is the filesystem type to mount.
	// Must be a filesystem type supported by the host operating system.
	// Ex. "ext4", "xfs", "ntfs". Implicitly inferred to be "ext4" if unspecified.
	// +optional
	FSType string `json:"fsType,omitempty" protobuf:"bytes,3,opt,name=fsType"`
	// readOnly defaults to false (read/write). ReadOnly here will force
	// the ReadOnly setting in VolumeMounts.
	// +optional
	ReadOnly bool `json:"readOnly,omitempty" protobuf:"varint,4,opt,name=readOnly"`
	// secretRef specifies the secret to use for obtaining the StorageOS API
	// credentials.  If not specified, default values will be attempted.
	// +optional
	SecretRef *LocalObjectReference `json:"secretRef,omitempty" protobuf:"bytes,5,opt,name=secretRef"`
}

// Represents a StorageOS persistent volume resource.
type StorageOSPersistentVolumeSource struct {
	// volumeName is the human-readable name of the StorageOS volume.  Volume
	// names are only unique within a namespace.
	VolumeName string `json:"volumeName,omitempty" protobuf:"bytes,1,opt,name=volumeName"`
	// volumeNamespace specifies the scope of the volume within StorageOS.  If no
	// namespace is specified then the Pod's namespace will be used.  This allows the
	// Kubernetes name scoping to be mirrored within StorageOS for tighter integration.
	// Set VolumeName to any name to override the default behaviour.
	// Set to "default" if you are not using namespaces within StorageOS.
	// Namespaces that do not pre-exist within StorageOS will be created.
	// +optional
	VolumeNamespace string `json:"volumeNamespace,omitempty" protobuf:"bytes,2,opt,name=volumeNamespace"`
	// fsType is the filesystem type to mount.
	// Must be a filesystem type supported by the host operating system.
	// Ex. "ext4", "xfs", "ntfs". Implicitly inferred to be "ext4" if unspecified.
	// +optional
	FSType string `json:"fsType,omitempty" protobuf:"bytes,3,opt,name=fsType"`
	// readOnly defaults to false (read/write). ReadOnly here will force
	// the ReadOnly setting in VolumeMounts.
	// +optional
	ReadOnly bool `json:"readOnly,omitempty" protobuf:"varint,4,opt,name=readOnly"`
	// secretRef specifies the secret to use for obtaining the StorageOS API
	// credentials.  If not specified, default values will be attempted.
	// +optional
	SecretRef *ObjectReference `json:"secretRef,omitempty" protobuf:"bytes,5,opt,name=secretRef"`
}

// Adapts a ConfigMap into a volume.
//
// The contents of the target ConfigMap's Data field will be presented in a
// volume as files using the keys in the Data field as the file names, unless
// the items element is populated with specific mappings of keys to paths.
// ConfigMap volumes support ownership management and SELinux relabeling.
// 将ConfigMap转换为卷。
// 目标ConfigMap的Data字段的内容将以文件的形式呈现在卷中，使用Data字段中的键作为文件名，除非items元素被填充为特定的键到路径的映射。
// ConfigMap卷支持所有权管理和SELinux relabeling。
type ConfigMapVolumeSource struct {
	LocalObjectReference `json:",inline" protobuf:"bytes,1,opt,name=localObjectReference"`
	// items if unspecified, each key-value pair in the Data field of the referenced
	// ConfigMap will be projected into the volume as a file whose name is the
	// key and content is the value. If specified, the listed keys will be
	// projected into the specified paths, and unlisted keys will not be
	// present. If a key is specified which is not present in the ConfigMap,
	// the volume setup will error unless it is marked optional. Paths must be
	// relative and may not contain the '..' path or start with '..'.
	// +optional
	// items如果未指定，则引用ConfigMap中Data字段中的每个键值对将被投影到卷中，其名称是键，内容是值。
	// 如果指定，则列出的键将被投影到指定的路径中，未列出的键将不会存在。
	// 如果指定了一个不存在于ConfigMap中的键，则除非它被标记为可选，否则卷设置将出错。
	// 路径必须是相对路径，不能包含“..”路径或以“..”开头。
	Items []KeyToPath `json:"items,omitempty" protobuf:"bytes,2,rep,name=items"`
	// defaultMode is optional: mode bits used to set permissions on created files by default.
	// Must be an octal value between 0000 and 0777 or a decimal value between 0 and 511.
	// YAML accepts both octal and decimal values, JSON requires decimal values for mode bits.
	// Defaults to 0644.
	// Directories within the path are not affected by this setting.
	// This might be in conflict with other options that affect the file
	// mode, like fsGroup, and the result can be other mode bits set.
	// +optional
	// defaultMode是可选的：用于默认情况下设置创建文件的权限位的模式位。
	// 必须是0000和0777之间的八进制值或0和511之间的十进制值。
	// YAML接受八进制和十进制值，JSON要求模式位的十进制值。
	DefaultMode *int32 `json:"defaultMode,omitempty" protobuf:"varint,3,opt,name=defaultMode"`
	// optional specify whether the ConfigMap or its keys must be defined
	// +optional
	// optional指定ConfigMap或其键是否必须被定义
	Optional *bool `json:"optional,omitempty" protobuf:"varint,4,opt,name=optional"`
}

const (
	ConfigMapVolumeSourceDefaultMode int32 = 0644
)

// Adapts a ConfigMap into a projected volume.
//
// The contents of the target ConfigMap's Data field will be presented in a
// projected volume as files using the keys in the Data field as the file names,
// unless the items element is populated with specific mappings of keys to paths.
// Note that this is identical to a configmap volume source without the default
// mode.
// 将ConfigMap转换为投影卷。
// 目标ConfigMap的Data字段的内容将以文件的形式呈现在投影卷中，使用Data字段中的键作为文件名，除非items元素被填充为特定的键到路径的映射。
// 请注意，这与没有默认模式的configmap卷源相同。
type ConfigMapProjection struct {
	LocalObjectReference `json:",inline" protobuf:"bytes,1,opt,name=localObjectReference"`
	// items if unspecified, each key-value pair in the Data field of the referenced
	// ConfigMap will be projected into the volume as a file whose name is the
	// key and content is the value. If specified, the listed keys will be
	// projected into the specified paths, and unlisted keys will not be
	// present. If a key is specified which is not present in the ConfigMap,
	// the volume setup will error unless it is marked optional. Paths must be
	// relative and may not contain the '..' path or start with '..'.
	// +optional
	// items如果未指定，则引用ConfigMap中Data字段中的每个键值对将被投影到卷中，其名称是键，内容是值。
	Items []KeyToPath `json:"items,omitempty" protobuf:"bytes,2,rep,name=items"`
	// optional specify whether the ConfigMap or its keys must be defined
	// +optional
	Optional *bool `json:"optional,omitempty" protobuf:"varint,4,opt,name=optional"`
}

// ServiceAccountTokenProjection represents a projected service account token
// volume. This projection can be used to insert a service account token into
// the pods runtime filesystem for use against APIs (Kubernetes API Server or
// otherwise).
// ServiceAccountTokenProjection表示投影的服务帐户令牌卷。
// 可以使用此投影将服务帐户令牌插入到pod的运行时文件系统中，以便针对API（Kubernetes API服务器或其他）使用。
type ServiceAccountTokenProjection struct {
	// audience is the intended audience of the token. A recipient of a token
	// must identify itself with an identifier specified in the audience of the
	// token, and otherwise should reject the token. The audience defaults to the
	// identifier of the apiserver.
	// +optional
	// audience是令牌的目标受众。
	// 令牌的接收者必须使用令牌的受众中指定的标识符进行身份验证，否则应拒绝令牌。
	// 受众默认为apiserver的标识符。
	Audience string `json:"audience,omitempty" protobuf:"bytes,1,rep,name=audience"`
	// expirationSeconds is the requested duration of validity of the service
	// account token. As the token approaches expiration, the kubelet volume
	// plugin will proactively rotate the service account token. The kubelet will
	// start trying to rotate the token if the token is older than 80 percent of
	// its time to live or if the token is older than 24 hours.Defaults to 1 hour
	// and must be at least 10 minutes.
	// +optional
	// expirationSeconds是服务帐户令牌的有效期。
	// 令牌接近到期时，kubelet卷插件将主动旋转服务帐户令牌。
	// 如果令牌比其生存时间的80％旧或比24小时旧，则kubelet将开始尝试旋转令牌。
	// 默认为1小时，必须至少为10分钟。
	ExpirationSeconds *int64 `json:"expirationSeconds,omitempty" protobuf:"varint,2,opt,name=expirationSeconds"`
	// path is the path relative to the mount point of the file to project the
	// token into.
	// path是相对于挂载点的文件的路径，以将令牌投影到其中。
	Path string `json:"path" protobuf:"bytes,3,opt,name=path"`
}

// Represents a projected volume source
// 表示投影卷源
type ProjectedVolumeSource struct {
	// sources is the list of volume projections
	// +optional
	// sources是卷投影的列表
	Sources []VolumeProjection `json:"sources" protobuf:"bytes,1,rep,name=sources"`
	// defaultMode are the mode bits used to set permissions on created files by default.
	// Must be an octal value between 0000 and 0777 or a decimal value between 0 and 511.
	// YAML accepts both octal and decimal values, JSON requires decimal values for mode bits.
	// Directories within the path are not affected by this setting.
	// This might be in conflict with other options that affect the file
	// mode, like fsGroup, and the result can be other mode bits set.
	// +optional
	// defaultMode是用于默认情况下设置创建文件的权限位的模式位。
	DefaultMode *int32 `json:"defaultMode,omitempty" protobuf:"varint,2,opt,name=defaultMode"`
}

// Projection that may be projected along with other supported volume types
// 可以与其他支持的卷类型投影的投影
type VolumeProjection struct {
	// all types below are the supported types for projection into the same volume

	// secret information about the secret data to project
	// +optional
	// secret关于要投影的secret数据的信息
	Secret *SecretProjection `json:"secret,omitempty" protobuf:"bytes,1,opt,name=secret"`
	// downwardAPI information about the downwardAPI data to project
	// +optional
	// downwardAPI关于要投影的downwardAPI数据的信息
	DownwardAPI *DownwardAPIProjection `json:"downwardAPI,omitempty" protobuf:"bytes,2,opt,name=downwardAPI"`
	// configMap information about the configMap data to project
	// +optional
	// configMap关于要投影的configMap数据的信息
	ConfigMap *ConfigMapProjection `json:"configMap,omitempty" protobuf:"bytes,3,opt,name=configMap"`
	// serviceAccountToken is information about the serviceAccountToken data to project
	// +optional
	// serviceAccountToken是关于要投影的serviceAccountToken数据的信息
	ServiceAccountToken *ServiceAccountTokenProjection `json:"serviceAccountToken,omitempty" protobuf:"bytes,4,opt,name=serviceAccountToken"`
}

const (
	ProjectedVolumeSourceDefaultMode int32 = 0644
)

// Maps a string key to a path within a volume.
// 将字符串键映射到卷中的路径。
type KeyToPath struct {
	// key is the key to project.
	// key是要投影的键。
	Key string `json:"key" protobuf:"bytes,1,opt,name=key"`

	// path is the relative path of the file to map the key to.
	// May not be an absolute path.
	// May not contain the path element '..'.
	// May not start with the string '..'.
	// path是将键映射到的文件的相对路径。
	// 不可以是绝对路径。
	// 不可以包含路径元素“..”。
	// 不可以以“..”开头。
	Path string `json:"path" protobuf:"bytes,2,opt,name=path"`
	// mode is Optional: mode bits used to set permissions on this file.
	// Must be an octal value between 0000 and 0777 or a decimal value between 0 and 511.
	// YAML accepts both octal and decimal values, JSON requires decimal values for mode bits.
	// If not specified, the volume defaultMode will be used.
	// This might be in conflict with other options that affect the file
	// mode, like fsGroup, and the result can be other mode bits set.
	// +optional
	// mode是可选的：用于设置此文件权限的模式位。
	// 必须是0000和0777之间的八进制值或0和511之间的十进制值。
	// YAML接受八进制和十进制值，JSON需要十进制值。
	// 如果未指定，将使用卷defaultMode。
	// 这可能与影响文件模式的其他选项（如fsGroup）冲突，结果可能会设置其他模式位。
	Mode *int32 `json:"mode,omitempty" protobuf:"varint,3,opt,name=mode"`
}

// Local represents directly-attached storage with node affinity (Beta feature)
// Local表示具有节点亲和性的直接附加存储（Beta功能）
type LocalVolumeSource struct {
	// path of the full path to the volume on the node.
	// It can be either a directory or block device (disk, partition, ...).
	// path是节点上卷的完整路径。
	// 它可以是目录或块设备（磁盘，分区，...）。
	Path string `json:"path" protobuf:"bytes,1,opt,name=path"`

	// fsType is the filesystem type to mount.
	// It applies only when the Path is a block device.
	// Must be a filesystem type supported by the host operating system.
	// Ex. "ext4", "xfs", "ntfs". The default value is to auto-select a filesystem if unspecified.
	// +optional
	FSType *string `json:"fsType,omitempty" protobuf:"bytes,2,opt,name=fsType"`
}

// Represents storage that is managed by an external CSI volume driver (Beta feature)
// 表示由外部CSI卷驱动程序管理的存储（Beta功能）
type CSIPersistentVolumeSource struct {
	// driver is the name of the driver to use for this volume.
	// Required.
	// driver是用于此卷的驱动程序的名称。
	Driver string `json:"driver" protobuf:"bytes,1,opt,name=driver"`

	// volumeHandle is the unique volume name returned by the CSI volume
	// plugin’s CreateVolume to refer to the volume on all subsequent calls.
	// Required.
	// volumeHandle是CSI卷插件的CreateVolume返回的唯一卷名称，用于在所有后续调用中引用卷。
	VolumeHandle string `json:"volumeHandle" protobuf:"bytes,2,opt,name=volumeHandle"`

	// readOnly value to pass to ControllerPublishVolumeRequest.
	// Defaults to false (read/write).
	// +optional
	// readOnly值传递给ControllerPublishVolumeRequest。
	ReadOnly bool `json:"readOnly,omitempty" protobuf:"varint,3,opt,name=readOnly"`

	// fsType to mount. Must be a filesystem type supported by the host operating system.
	// Ex. "ext4", "xfs", "ntfs".
	// +optional
	// fsType挂载。必须是主机操作系统支持的文件系统类型。
	FSType string `json:"fsType,omitempty" protobuf:"bytes,4,opt,name=fsType"`

	// volumeAttributes of the volume to publish.
	// +optional
	VolumeAttributes map[string]string `json:"volumeAttributes,omitempty" protobuf:"bytes,5,rep,name=volumeAttributes"`

	// controllerPublishSecretRef is a reference to the secret object containing
	// sensitive information to pass to the CSI driver to complete the CSI
	// ControllerPublishVolume and ControllerUnpublishVolume calls.
	// This field is optional, and may be empty if no secret is required. If the
	// secret object contains more than one secret, all secrets are passed.
	// +optional
	// controllerPublishSecretRef是指向包含要传递给CSI驱动程序以完成CSI ControllerPublishVolume和ControllerUnpublishVolume调用的敏感信息的密钥对象的引用。
	// 此字段是可选的，如果不需要密钥，则可以为空。如果密钥对象包含多个密钥，则传递所有密钥。
	ControllerPublishSecretRef *SecretReference `json:"controllerPublishSecretRef,omitempty" protobuf:"bytes,6,opt,name=controllerPublishSecretRef"`

	// nodeStageSecretRef is a reference to the secret object containing sensitive
	// information to pass to the CSI driver to complete the CSI NodeStageVolume
	// and NodeStageVolume and NodeUnstageVolume calls.
	// This field is optional, and may be empty if no secret is required. If the
	// secret object contains more than one secret, all secrets are passed.
	// +optional
	// nodeStageSecretRef是指向包含要传递给CSI驱动程序以完成CSI NodeStageVolume和NodeStageVolume和NodeUnstageVolume调用的敏感信息的密钥对象的引用。
	// 此字段是可选的，如果不需要密钥，则可以为空。如果密钥对象包含多个密钥，则传递所有密钥。
	NodeStageSecretRef *SecretReference `json:"nodeStageSecretRef,omitempty" protobuf:"bytes,7,opt,name=nodeStageSecretRef"`

	// nodePublishSecretRef is a reference to the secret object containing
	// sensitive information to pass to the CSI driver to complete the CSI
	// NodePublishVolume and NodeUnpublishVolume calls.
	// This field is optional, and may be empty if no secret is required. If the
	// secret object contains more than one secret, all secrets are passed.
	// +optional
	// nodePublishSecretRef是指向包含要传递给CSI驱动程序以完成CSI NodePublishVolume和NodeUnpublishVolume调用的敏感信息的密钥对象的引用。
	// 此字段是可选的，如果不需要密钥，则可以为空。如果密钥对象包含多个密钥，则传递所有密钥。
	NodePublishSecretRef *SecretReference `json:"nodePublishSecretRef,omitempty" protobuf:"bytes,8,opt,name=nodePublishSecretRef"`

	// controllerExpandSecretRef is a reference to the secret object containing
	// sensitive information to pass to the CSI driver to complete the CSI
	// ControllerExpandVolume call.
	// This field is optional, and may be empty if no secret is required. If the
	// secret object contains more than one secret, all secrets are passed.
	// +optional
	// controllerExpandSecretRef是指向包含要传递给CSI驱动程序以完成CSI ControllerExpandVolume调用的敏感信息的密钥对象的引用。
	// 此字段是可选的，如果不需要密钥，则可以为空。如果密钥对象包含多个密钥，则传递所有密钥。
	ControllerExpandSecretRef *SecretReference `json:"controllerExpandSecretRef,omitempty" protobuf:"bytes,9,opt,name=controllerExpandSecretRef"`

	// nodeExpandSecretRef is a reference to the secret object containing
	// sensitive information to pass to the CSI driver to complete the CSI
	// NodeExpandVolume call.
	// This is an alpha field and requires enabling CSINodeExpandSecret feature gate.
	// This field is optional, may be omitted if no secret is required. If the
	// secret object contains more than one secret, all secrets are passed.
	// +optional
	// nodeExpandSecretRef是指向包含要传递给CSI驱动程序以完成CSI NodeExpandVolume调用的敏感信息的密钥对象的引用。
	// 这是一个alpha字段，需要启用CSINodeExpandSecret功能门。
	NodeExpandSecretRef *SecretReference `json:"nodeExpandSecretRef,omitempty" protobuf:"bytes,10,opt,name=nodeExpandSecretRef"`
}

// Represents a source location of a volume to mount, managed by an external CSI driver
// 表示要挂载的卷的源位置，由外部CSI驱动程序管理
type CSIVolumeSource struct {
	// driver is the name of the CSI driver that handles this volume.
	// Consult with your admin for the correct name as registered in the cluster.
	// driver是处理此卷的CSI驱动程序的名称。
	// 与您的管理员联系，以便获得在群集中注册的正确名称。
	Driver string `json:"driver" protobuf:"bytes,1,opt,name=driver"`

	// readOnly specifies a read-only configuration for the volume.
	// Defaults to false (read/write).
	// +optional
	// readOnly指定卷的只读配置。 默认为false（读/写）。
	ReadOnly *bool `json:"readOnly,omitempty" protobuf:"varint,2,opt,name=readOnly"`

	// fsType to mount. Ex. "ext4", "xfs", "ntfs".
	// If not provided, the empty value is passed to the associated CSI driver
	// which will determine the default filesystem to apply.
	// +optional
	// fsType挂载。 例如“ext4”，“xfs”，“ntfs”。
	FSType *string `json:"fsType,omitempty" protobuf:"bytes,3,opt,name=fsType"`

	// volumeAttributes stores driver-specific properties that are passed to the CSI
	// driver. Consult your driver's documentation for supported values.
	// +optional
	// volumeAttributes存储传递给CSI驱动程序的驱动程序特定属性。 请参阅您的驱动程序的文档以获取支持的值。
	VolumeAttributes map[string]string `json:"volumeAttributes,omitempty" protobuf:"bytes,4,rep,name=volumeAttributes"`

	// nodePublishSecretRef is a reference to the secret object containing
	// sensitive information to pass to the CSI driver to complete the CSI
	// NodePublishVolume and NodeUnpublishVolume calls.
	// This field is optional, and  may be empty if no secret is required. If the
	// secret object contains more than one secret, all secret references are passed.
	// +optional
	// nodePublishSecretRef是指向包含要传递给CSI驱动程序以完成CSI NodePublishVolume和NodeUnpublishVolume调用的敏感信息的密钥对象的引用。
	// 此字段是可选的，如果不需要密钥，则可以为空。如果密钥对象包含多个密钥，则传递所有密钥。
	NodePublishSecretRef *LocalObjectReference `json:"nodePublishSecretRef,omitempty" protobuf:"bytes,5,opt,name=nodePublishSecretRef"`
}

// Represents an ephemeral volume that is handled by a normal storage driver.
// 表示由普通存储驱动程序处理的临时卷。
type EphemeralVolumeSource struct {
	// Will be used to create a stand-alone PVC to provision the volume.
	// The pod in which this EphemeralVolumeSource is embedded will be the
	// owner of the PVC, i.e. the PVC will be deleted together with the
	// pod.  The name of the PVC will be `<pod name>-<volume name>` where
	// `<volume name>` is the name from the `PodSpec.Volumes` array
	// entry. Pod validation will reject the pod if the concatenated name
	// is not valid for a PVC (for example, too long).
	//
	// An existing PVC with that name that is not owned by the pod
	// will *not* be used for the pod to avoid using an unrelated
	// volume by mistake. Starting the pod is then blocked until
	// the unrelated PVC is removed. If such a pre-created PVC is
	// meant to be used by the pod, the PVC has to updated with an
	// owner reference to the pod once the pod exists. Normally
	// this should not be necessary, but it may be useful when
	// manually reconstructing a broken cluster.
	//
	// This field is read-only and no changes will be made by Kubernetes
	// to the PVC after it has been created.
	//
	// Required, must not be nil.
	// 将用于创建独立的PVC以提供卷。 嵌入此EphemeralVolumeSource的pod将是PVC的所有者，即PVC将与pod一起删除。
	// PVC的名称将是`<pod name>-<volume name>`，其中`<volume name>`是来自`PodSpec.Volumes`数组条目的名称。
	// Pod验证将拒绝pod，如果连接的名称对于PVC无效（例如，太长）。
	// 与pod不相关的现有PVC将不会用于pod，以避免错误地使用不相关的卷。 启动pod被阻止，直到删除不相关的PVC。
	// 如果要为pod使用这样的预先创建的PVC，则必须在pod存在时使用所有者引用更新PVC。
	// 通常不需要这样做，但在手动重建损坏的群集时可能很有用。
	// 此字段是只读的，Kubernetes在创建后不会对PVC进行任何更改。
	VolumeClaimTemplate *PersistentVolumeClaimTemplate `json:"volumeClaimTemplate,omitempty" protobuf:"bytes,1,opt,name=volumeClaimTemplate"`

	// ReadOnly is tombstoned to show why 2 is a reserved protobuf tag.
	// ReadOnly bool `json:"readOnly,omitempty" protobuf:"varint,2,opt,name=readOnly"`
}

// PersistentVolumeClaimTemplate is used to produce
// PersistentVolumeClaim objects as part of an EphemeralVolumeSource.
// PersistentVolumeClaimTemplate是用于生成PersistentVolumeClaim对象的一部分，作为EphemeralVolumeSource的一部分。
type PersistentVolumeClaimTemplate struct {
	// May contain labels and annotations that will be copied into the PVC
	// when creating it. No other fields are allowed and will be rejected during
	// validation.
	//
	// +optional
	// 可以包含在创建PVC时将复制到PVC中的标签和注释。 不允许其他字段，并且在验证期间将被拒绝。
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// The specification for the PersistentVolumeClaim. The entire content is
	// copied unchanged into the PVC that gets created from this
	// template. The same fields as in a PersistentVolumeClaim
	// are also valid here.
	// PersistentVolumeClaim的规范。 整个内容都被复制到从此模板创建的PVC中。 PersistentVolumeClaim中的相同字段也在这里有效。
	Spec PersistentVolumeClaimSpec `json:"spec" protobuf:"bytes,2,name=spec"`
}

// ContainerPort represents a network port in a single container.
// ContainerPort表示单个容器中的网络端口。
type ContainerPort struct {
	// If specified, this must be an IANA_SVC_NAME and unique within the pod. Each
	// named port in a pod must have a unique name. Name for the port that can be
	// referred to by services.
	// +optional
	// 如果指定，这必须是IANA_SVC_NAME，并且在pod中唯一。 pod中的每个命名端口必须具有唯一的名称。 可以由服务引用的端口的名称。
	Name string `json:"name,omitempty" protobuf:"bytes,1,opt,name=name"`
	// Number of port to expose on the host.
	// If specified, this must be a valid port number, 0 < x < 65536.
	// If HostNetwork is specified, this must match ContainerPort.
	// Most containers do not need this.
	// +optional
	// 在主机上公开的端口号。 如果指定，这必须是有效的端口号，0 <x <65536。 如果指定了HostNetwork，则必须与ContainerPort匹配。 大多数容器不需要这个。
	HostPort int32 `json:"hostPort,omitempty" protobuf:"varint,2,opt,name=hostPort"`
	// Number of port to expose on the pod's IP address.
	// This must be a valid port number, 0 < x < 65536.
	// 在pod的IP地址上公开的端口数。这必须是有效的端口号，0 < x < 65536。
	ContainerPort int32 `json:"containerPort" protobuf:"varint,3,opt,name=containerPort"`
	// Protocol for port. Must be UDP, TCP, or SCTP.
	// Defaults to "TCP".
	// +optional
	// +default="TCP"
	// 端口的协议。 必须是UDP，TCP或SCTP。 默认为“TCP”。
	Protocol Protocol `json:"protocol,omitempty" protobuf:"bytes,4,opt,name=protocol,casttype=Protocol"`
	// What host IP to bind the external port to.
	// +optional
	// 将外部端口绑定到哪个主机IP。
	HostIP string `json:"hostIP,omitempty" protobuf:"bytes,5,opt,name=hostIP"`
}

// VolumeMount describes a mounting of a Volume within a container.
// VolumeMount描述容器中卷的挂载。
type VolumeMount struct {
	// This must match the Name of a Volume.
	// 这必须与Volume的名称匹配。
	Name string `json:"name" protobuf:"bytes,1,opt,name=name"`
	// Mounted read-only if true, read-write otherwise (false or unspecified).
	// Defaults to false.
	// +optional
	// 如果为true，则以只读方式挂载，否则以读写方式挂载（false或未指定）。 默认为false。
	ReadOnly bool `json:"readOnly,omitempty" protobuf:"varint,2,opt,name=readOnly"`
	// Path within the container at which the volume should be mounted.  Must
	// not contain ':'.
	// 应该在其中挂载卷的容器中的路径。不能包含“:”。
	MountPath string `json:"mountPath" protobuf:"bytes,3,opt,name=mountPath"`
	// Path within the volume from which the container's volume should be mounted.
	// Defaults to "" (volume's root).
	// +optional
	// 卷中的路径，应该从中挂载容器的卷。 默认为“”（卷的根）。
	SubPath string `json:"subPath,omitempty" protobuf:"bytes,4,opt,name=subPath"`
	// mountPropagation determines how mounts are propagated from the host
	// to container and the other way around.
	// When not set, MountPropagationNone is used.
	// This field is beta in 1.10.
	// +optional
	// mountPropagation确定如何从主机传播挂载到容器和其他方式。 当未设置时，使用MountPropagationNone。 该字段在1.10中为beta版。
	MountPropagation *MountPropagationMode `json:"mountPropagation,omitempty" protobuf:"bytes,5,opt,name=mountPropagation,casttype=MountPropagationMode"`
	// Expanded path within the volume from which the container's volume should be mounted.
	// Behaves similarly to SubPath but environment variable references $(VAR_NAME) are expanded using the container's environment.
	// Defaults to "" (volume's root).
	// SubPathExpr and SubPath are mutually exclusive.
	// +optional
	// 卷中的扩展路径，应该从中挂载容器的卷。 行为类似于SubPath，但使用容器的环境对环境变量引用$(VAR_NAME)进行扩展。 默认为“”（卷的根）。 SubPathExpr和SubPath互斥。
	SubPathExpr string `json:"subPathExpr,omitempty" protobuf:"bytes,6,opt,name=subPathExpr"`
}

// MountPropagationMode describes mount propagation.
// +enum
// MountPropagationMode 描述挂载传播。
type MountPropagationMode string

const (
	// MountPropagationNone means that the volume in a container will
	// not receive new mounts from the host or other containers, and filesystems
	// mounted inside the container won't be propagated to the host or other
	// containers.
	// Note that this mode corresponds to "private" in Linux terminology.
	// MountPropagationNone表示容器中的卷不会从主机或其他容器接收新的挂载，并且在容器内部挂载的文件系统不会传播到主机或其他容器。 请注意，此模式对应于Linux术语中的“private”。
	MountPropagationNone MountPropagationMode = "None"
	// MountPropagationHostToContainer means that the volume in a container will
	// receive new mounts from the host or other containers, but filesystems
	// mounted inside the container won't be propagated to the host or other
	// containers.
	// Note that this mode is recursively applied to all mounts in the volume
	// ("rslave" in Linux terminology).
	// MountPropagationHostToContainer表示容器中的卷将从主机或其他容器接收新的挂载，但在容器内部挂载的文件系统不会传播到主机或其他容器。 请注意，此模式递归应用于卷中的所有挂载（Linux术语中的“rslave”）。
	MountPropagationHostToContainer MountPropagationMode = "HostToContainer"
	// MountPropagationBidirectional means that the volume in a container will
	// receive new mounts from the host or other containers, and its own mounts
	// will be propagated from the container to the host or other containers.
	// Note that this mode is recursively applied to all mounts in the volume
	// ("rshared" in Linux terminology).
	// MountPropagationBidirectional表示容器中的卷将从主机或其他容器接收新的挂载，并且其自己的挂载将从容器传播到主机或其他容器。 请注意，此模式递归应用于卷中的所有挂载（Linux术语中的“rshared”）。
	MountPropagationBidirectional MountPropagationMode = "Bidirectional"
)

// volumeDevice describes a mapping of a raw block device within a container.
// volumeDevice描述容器中原始块设备的映射。
type VolumeDevice struct {
	// name must match the name of a persistentVolumeClaim in the pod
	// name必须与pod中的persistentVolumeClaim的名称匹配
	Name string `json:"name" protobuf:"bytes,1,opt,name=name"`
	// devicePath is the path inside of the container that the device will be mapped to.
	// devicePath是容器内部将映射到设备的路径。
	DevicePath string `json:"devicePath" protobuf:"bytes,2,opt,name=devicePath"`
}

// EnvVar represents an environment variable present in a Container.
// EnvVar表示容器中存在的环境变量。
type EnvVar struct {
	// Name of the environment variable. Must be a C_IDENTIFIER.
	// 环境变量的名称。必须是C_IDENTIFIER。
	Name string `json:"name" protobuf:"bytes,1,opt,name=name"`

	// Optional: no more than one of the following may be specified.

	// Variable references $(VAR_NAME) are expanded
	// using the previously defined environment variables in the container and
	// any service environment variables. If a variable cannot be resolved,
	// the reference in the input string will be unchanged. Double $$ are reduced
	// to a single $, which allows for escaping the $(VAR_NAME) syntax: i.e.
	// "$$(VAR_NAME)" will produce the string literal "$(VAR_NAME)".
	// Escaped references will never be expanded, regardless of whether the variable
	// exists or not.
	// Defaults to "".
	// +optional
	// 可选：最多只能指定以下之一。
	// 变量引用$(VAR_NAME)使用容器中先前定义的环境变量和任何服务环境变量进行扩展。 如果无法解析变量，则输入字符串中的引用将保持不变。 双$$被减少为单个$，这允许对$(VAR_NAME)语法进行转义：即“$$(VAR_NAME)”将生成字符串文字“$(VAR_NAME)”。
	// 转义引用永远不会被扩展，无论变量是否存在。
	// 默认为“”。
	Value string `json:"value,omitempty" protobuf:"bytes,2,opt,name=value"`
	// Source for the environment variable's value. Cannot be used if value is not empty.
	// +optional
	// 环境变量值的来源。 如果value不为空，则不能使用。
	ValueFrom *EnvVarSource `json:"valueFrom,omitempty" protobuf:"bytes,3,opt,name=valueFrom"`
}

// EnvVarSource represents a source for the value of an EnvVar.
// EnvVarSource表示EnvVar值的来源。
type EnvVarSource struct {
	// Selects a field of the pod: supports metadata.name, metadata.namespace, `metadata.labels['<KEY>']`, `metadata.annotations['<KEY>']`,
	// spec.nodeName, spec.serviceAccountName, status.hostIP, status.podIP, status.podIPs.
	// +optional
	// 选择pod的字段：支持metadata.name，metadata.namespace，`metadata.labels['<KEY>']`，`metadata.annotations['<KEY>']`，spec.nodeName，spec.serviceAccountName，status.hostIP，status.podIP，status.podIPs。
	FieldRef *ObjectFieldSelector `json:"fieldRef,omitempty" protobuf:"bytes,1,opt,name=fieldRef"`
	// Selects a resource of the container: only resources limits and requests
	// (limits.cpu, limits.memory, limits.ephemeral-storage, requests.cpu, requests.memory and requests.ephemeral-storage) are currently supported.
	// +optional
	// 选择容器的资源：目前仅支持资源限制和请求（limits.cpu，limits.memory，limits.ephemeral-storage，requests.cpu，requests.memory和requests.ephemeral-storage）。
	ResourceFieldRef *ResourceFieldSelector `json:"resourceFieldRef,omitempty" protobuf:"bytes,2,opt,name=resourceFieldRef"`
	// Selects a key of a ConfigMap.
	// +optional
	// 选择ConfigMap的键。
	ConfigMapKeyRef *ConfigMapKeySelector `json:"configMapKeyRef,omitempty" protobuf:"bytes,3,opt,name=configMapKeyRef"`
	// Selects a key of a secret in the pod's namespace
	// +optional
	// 选择pod命名空间中的密钥
	SecretKeyRef *SecretKeySelector `json:"secretKeyRef,omitempty" protobuf:"bytes,4,opt,name=secretKeyRef"`
}

// ObjectFieldSelector selects an APIVersioned field of an object.
// +structType=atomic
// ObjectFieldSelector选择对象的APIVersioned字段。
type ObjectFieldSelector struct {
	// Version of the schema the FieldPath is written in terms of, defaults to "v1".
	// +optional
	// FieldPath的模式的模式版本，默认为“v1”。
	APIVersion string `json:"apiVersion,omitempty" protobuf:"bytes,1,opt,name=apiVersion"`
	// Path of the field to select in the specified API version.
	// 指定API版本中选择字段的路径。
	FieldPath string `json:"fieldPath" protobuf:"bytes,2,opt,name=fieldPath"`
}

// ResourceFieldSelector represents container resources (cpu, memory) and their output format
// +structType=atomic
// ResourceFieldSelector表示容器资源（cpu，内存）及其输出格式
type ResourceFieldSelector struct {
	// Container name: required for volumes, optional for env vars
	// +optional
	// 容器名称：卷需要，环境变量可选
	ContainerName string `json:"containerName,omitempty" protobuf:"bytes,1,opt,name=containerName"`
	// Required: resource to select
	// 必需：选择资源
	Resource string `json:"resource" protobuf:"bytes,2,opt,name=resource"`
	// Specifies the output format of the exposed resources, defaults to "1"
	// +optional
	// 指定公开资源的输出格式，默认为“1”
	Divisor resource.Quantity `json:"divisor,omitempty" protobuf:"bytes,3,opt,name=divisor"`
}

// Selects a key from a ConfigMap.
// +structType=atomic
// 从ConfigMap中选择一个键。
type ConfigMapKeySelector struct {
	// The ConfigMap to select from.
	// 要从中选择的ConfigMap。
	LocalObjectReference `json:",inline" protobuf:"bytes,1,opt,name=localObjectReference"`
	// The key to select.
	// 要选择的键。
	Key string `json:"key" protobuf:"bytes,2,opt,name=key"`
	// Specify whether the ConfigMap or its key must be defined
	// +optional
	// 指定ConfigMap或其键是否必须定义
	Optional *bool `json:"optional,omitempty" protobuf:"varint,3,opt,name=optional"`
}

// SecretKeySelector selects a key of a Secret.
// +structType=atomic
// SecretKeySelector选择Secret的键。
type SecretKeySelector struct {
	// The name of the secret in the pod's namespace to select from.
	// 要从pod命名空间中选择的密钥的名称。
	LocalObjectReference `json:",inline" protobuf:"bytes,1,opt,name=localObjectReference"`
	// The key of the secret to select from.  Must be a valid secret key.
	// 要从中选择的秘密的钥匙。必须是有效的密钥。
	Key string `json:"key" protobuf:"bytes,2,opt,name=key"`
	// Specify whether the Secret or its key must be defined
	// +optional
	// 指定Secret或其键是否必须定义
	Optional *bool `json:"optional,omitempty" protobuf:"varint,3,opt,name=optional"`
}

// EnvFromSource represents the source of a set of ConfigMaps
// EnvFromSource表示一组ConfigMaps的来源。
type EnvFromSource struct {
	// An optional identifier to prepend to each key in the ConfigMap. Must be a C_IDENTIFIER.
	// +optional
	// 一个可选的标识符，用于在ConfigMap中的每个键之前添加。必须是C_IDENTIFIER。
	Prefix string `json:"prefix,omitempty" protobuf:"bytes,1,opt,name=prefix"`
	// The ConfigMap to select from
	// +optional
	// 要从中选择的ConfigMap。
	ConfigMapRef *ConfigMapEnvSource `json:"configMapRef,omitempty" protobuf:"bytes,2,opt,name=configMapRef"`
	// The Secret to select from
	// +optional
	// 要从中选择的密钥。
	SecretRef *SecretEnvSource `json:"secretRef,omitempty" protobuf:"bytes,3,opt,name=secretRef"`
}

// ConfigMapEnvSource selects a ConfigMap to populate the environment
// variables with.
//
// The contents of the target ConfigMap's Data field will represent the
// key-value pairs as environment variables.
// ConfigMapEnvSource选择一个ConfigMap来填充环境变量。
// 目标ConfigMap的Data字段的内容将表示键值对作为环境变量。
type ConfigMapEnvSource struct {
	// The ConfigMap to select from.
	// 要从中选择的ConfigMap。
	LocalObjectReference `json:",inline" protobuf:"bytes,1,opt,name=localObjectReference"`
	// Specify whether the ConfigMap must be defined
	// +optional
	// 指定ConfigMap是否必须定义
	Optional *bool `json:"optional,omitempty" protobuf:"varint,2,opt,name=optional"`
}

// SecretEnvSource selects a Secret to populate the environment
// variables with.
//
// The contents of the target Secret's Data field will represent the
// key-value pairs as environment variables.
// SecretEnvSource选择一个Secret来填充环境变量。
// 目标Secret的Data字段的内容将表示键值对作为环境变量。
type SecretEnvSource struct {
	// The Secret to select from.
	LocalObjectReference `json:",inline" protobuf:"bytes,1,opt,name=localObjectReference"`
	// Specify whether the Secret must be defined
	// +optional
	// 指定Secret是否必须定义
	Optional *bool `json:"optional,omitempty" protobuf:"varint,2,opt,name=optional"`
}

// HTTPHeader describes a custom header to be used in HTTP probes
// HTTPHeader描述要在HTTP探测中使用的自定义标头
type HTTPHeader struct {
	// The header field name
	// 标头字段名称
	Name string `json:"name" protobuf:"bytes,1,opt,name=name"`
	// The header field value
	// 标头字段值
	Value string `json:"value" protobuf:"bytes,2,opt,name=value"`
}

// HTTPGetAction describes an action based on HTTP Get requests.
// HTTPGetAction描述基于HTTP Get请求的操作。
type HTTPGetAction struct {
	// Path to access on the HTTP server.
	// +optional
	// 要访问的HTTP服务器的路径。
	Path string `json:"path,omitempty" protobuf:"bytes,1,opt,name=path"`
	// Name or number of the port to access on the container.
	// Number must be in the range 1 to 65535.
	// Name must be an IANA_SVC_NAME.
	// 名称或者容器上访问的端口号。
	// 号码必须在1到65535的范围内。
	// 名称必须是IANA_SVC_NAME。
	Port intstr.IntOrString `json:"port" protobuf:"bytes,2,opt,name=port"`
	// Host name to connect to, defaults to the pod IP. You probably want to set
	// "Host" in httpHeaders instead.
	// +optional
	// 要连接的主机名，默认为pod IP。您可能想在httpHeaders中设置“Host”。
	Host string `json:"host,omitempty" protobuf:"bytes,3,opt,name=host"`
	// Scheme to use for connecting to the host.
	// Defaults to HTTP.
	// +optional
	// 用于连接到主机的方案。默认为HTTP。
	Scheme URIScheme `json:"scheme,omitempty" protobuf:"bytes,4,opt,name=scheme,casttype=URIScheme"`
	// Custom headers to set in the request. HTTP allows repeated headers.
	// +optional
	// 要在请求中设置的自定义标头。HTTP允许重复标头。
	HTTPHeaders []HTTPHeader `json:"httpHeaders,omitempty" protobuf:"bytes,5,rep,name=httpHeaders"`
}

// URIScheme identifies the scheme used for connection to a host for Get actions
// URIScheme标识用于连接到主机的方案以获取操作
// +enum
type URIScheme string

const (
	// URISchemeHTTP means that the scheme used will be http://
	URISchemeHTTP URIScheme = "HTTP"
	// URISchemeHTTPS means that the scheme used will be https://
	URISchemeHTTPS URIScheme = "HTTPS"
)

// TCPSocketAction describes an action based on opening a socket
// TCPSocketAction描述基于打开套接字的操作
type TCPSocketAction struct {
	// Number or name of the port to access on the container.
	// Number must be in the range 1 to 65535.
	// Name must be an IANA_SVC_NAME.
	// 要访问容器上的端口号或名称。
	// 号码必须在1到65535的范围内。
	// 名称必须是IANA_SVC_NAME。
	Port intstr.IntOrString `json:"port" protobuf:"bytes,1,opt,name=port"`
	// Optional: Host name to connect to, defaults to the pod IP.
	// +optional
	// 可选：要连接的主机名，默认为pod IP。
	Host string `json:"host,omitempty" protobuf:"bytes,2,opt,name=host"`
}

type GRPCAction struct {
	// Port number of the gRPC service. Number must be in the range 1 to 65535.
	// gRPC服务端口号。Number的取值范围为1 ~ 65535。
	Port int32 `json:"port" protobuf:"bytes,1,opt,name=port"`

	// Service is the name of the service to place in the gRPC HealthCheckRequest
	// (see https://github.com/grpc/grpc/blob/master/doc/health-checking.md).
	//
	// If this is not specified, the default behavior is defined by gRPC.
	// +optional
	// +default=""
	// Service是要放在gRPC HealthCheckRequest中的服务名称
	// 如果未指定此字段，则gRPC的默认行为将被定义。
	Service *string `json:"service" protobuf:"bytes,2,opt,name=service"`
}

// ExecAction describes a "run in container" action.
// ExecAction描述了“在容器中运行”的操作。
type ExecAction struct {
	// Command is the command line to execute inside the container, the working directory for the
	// command  is root ('/') in the container's filesystem. The command is simply exec'd, it is
	// not run inside a shell, so traditional shell instructions ('|', etc) won't work. To use
	// a shell, you need to explicitly call out to that shell.
	// Exit status of 0 is treated as live/healthy and non-zero is unhealthy.
	// +optional
	// Command是要在容器内执行的命令行，命令的工作目录是容器文件系统中的根目录（“/”）。
	// 命令只是exec'd，它不是在shell内运行，因此传统的shell指令（“|”等）将不起作用。
	// 要使用shell，您需要显式地调用该shell。
	// 0的退出状态被视为live/healthy，非零状态被视为unhealthy。
	Command []string `json:"command,omitempty" protobuf:"bytes,1,rep,name=command"`
}

// Probe describes a health check to be performed against a container to determine whether it is
// alive or ready to receive traffic.
// Probe描述了要对容器执行的健康检查，以确定它是否存活或准备好接收流量。
type Probe struct {
	// The action taken to determine the health of a container
	// 用于确定容器健康状况的操作
	ProbeHandler `json:",inline" protobuf:"bytes,1,opt,name=handler"`
	// Number of seconds after the container has started before liveness probes are initiated.
	// More info: https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle#container-probes
	// +optional
	// 容器启动后，开始执行存活探测的秒数。
	InitialDelaySeconds int32 `json:"initialDelaySeconds,omitempty" protobuf:"varint,2,opt,name=initialDelaySeconds"`
	// Number of seconds after which the probe times out.
	// Defaults to 1 second. Minimum value is 1.
	// More info: https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle#container-probes
	// +optional
	// 探测超时的秒数。 默认为1秒。 最小值为1。
	TimeoutSeconds int32 `json:"timeoutSeconds,omitempty" protobuf:"varint,3,opt,name=timeoutSeconds"`
	// How often (in seconds) to perform the probe.
	// Default to 10 seconds. Minimum value is 1.
	// +optional
	// 要执行探测的频率（以秒为单位）。 默认为10秒。 最小值为1。
	PeriodSeconds int32 `json:"periodSeconds,omitempty" protobuf:"varint,4,opt,name=periodSeconds"`
	// Minimum consecutive successes for the probe to be considered successful after having failed.
	// Defaults to 1. Must be 1 for liveness and startup. Minimum value is 1.
	// +optional
	// 探测失败后，要被视为成功的最小连续成功次数。 默认为1。 对于存活和启动，必须为1。 最小值为1。
	SuccessThreshold int32 `json:"successThreshold,omitempty" protobuf:"varint,5,opt,name=successThreshold"`
	// Minimum consecutive failures for the probe to be considered failed after having succeeded.
	// Defaults to 3. Minimum value is 1.
	// +optional
	// 探测成功后，要被视为失败的最小连续失败次数。 默认为3。 最小值为1。
	FailureThreshold int32 `json:"failureThreshold,omitempty" protobuf:"varint,6,opt,name=failureThreshold"`
	// Optional duration in seconds the pod needs to terminate gracefully upon probe failure.
	// The grace period is the duration in seconds after the processes running in the pod are sent
	// a termination signal and the time when the processes are forcibly halted with a kill signal.
	// Set this value longer than the expected cleanup time for your process.
	// If this value is nil, the pod's terminationGracePeriodSeconds will be used. Otherwise, this
	// value overrides the value provided by the pod spec.
	// Value must be non-negative integer. The value zero indicates stop immediately via
	// the kill signal (no opportunity to shut down).
	// This is a beta field and requires enabling ProbeTerminationGracePeriod feature gate.
	// Minimum value is 1. spec.terminationGracePeriodSeconds is used if unset.
	// +optional
	// 探测失败时，pod需要优雅终止的可选持续时间（以秒为单位）。 宽限期是在pod中运行的进程发送终止信号后的持续时间，
	// 并且是在强制使用kill信号终止进程时的时间。 将此值设置为大于进程预期清理时间的值。
	// 如果此值为nil，则将使用pod的terminationGracePeriodSeconds。 否则，此值将覆盖pod spec提供的值。
	// 值必须为非负整数。 值为零表示通过kill信号立即停止（没有机会关闭）。
	// 这是一个beta字段，需要启用ProbeTerminationGracePeriod功能门。
	// 最小值为1。 如果未设置，则使用spec.terminationGracePeriodSeconds。
	TerminationGracePeriodSeconds *int64 `json:"terminationGracePeriodSeconds,omitempty" protobuf:"varint,7,opt,name=terminationGracePeriodSeconds"`
}

// PullPolicy describes a policy for if/when to pull a container image
// +enum
// PullPolicy 描述了何时/何时拉取容器镜像的策略
type PullPolicy string

const (
	// PullAlways means that kubelet always attempts to pull the latest image. Container will fail If the pull fails.
	// PullAlways 意味着kubelet总是尝试拉取最新的镜像。 如果拉取失败，容器将失败。
	PullAlways PullPolicy = "Always"
	// PullNever means that kubelet never pulls an image, but only uses a local image. Container will fail if the image isn't present
	// PullNever 意味着kubelet永远不会拉取镜像，而只使用本地镜像。 如果镜像不存在，容器将失败
	PullNever PullPolicy = "Never"
	// PullIfNotPresent means that kubelet pulls if the image isn't present on disk. Container will fail if the image isn't present and the pull fails.
	// PullIfNotPresent 意味着kubelet会在镜像不在磁盘上时拉取。 如果镜像不存在并且拉取失败，容器将失败。
	PullIfNotPresent PullPolicy = "IfNotPresent"
)

// PreemptionPolicy describes a policy for if/when to preempt a pod.
// +enum
// PreemptionPolicy 描述了何时/何时抢占pod的策略。
type PreemptionPolicy string

const (
	// PreemptLowerPriority means that pod can preempt other pods with lower priority.
	// PreemptLowerPriority 意味着pod可以抢占其他优先级较低的pod。
	PreemptLowerPriority PreemptionPolicy = "PreemptLowerPriority"
	// PreemptNever means that pod never preempts other pods with lower priority.
	// PreemptNever 意味着pod永远不会抢占其他优先级较低的pod。
	PreemptNever PreemptionPolicy = "Never"
)

// TerminationMessagePolicy describes how termination messages are retrieved from a container.
// TerminationMessagePolicy 描述了如何从容器中检索终止消息。
// +enum
type TerminationMessagePolicy string

const (
	// TerminationMessageReadFile is the default behavior and will set the container status message to
	// the contents of the container's terminationMessagePath when the container exits.
	// TerminationMessageReadFile 是默认行为，并将将容器状态消息设置为容器的terminationMessagePath的内容，当容器退出时。
	TerminationMessageReadFile TerminationMessagePolicy = "File"
	// TerminationMessageFallbackToLogsOnError will read the most recent contents of the container logs
	// for the container status message when the container exits with an error and the
	// terminationMessagePath has no contents.
	// TerminationMessageFallbackToLogsOnError 将在容器退出时，读取容器日志的最新内容，以便容器状态消息，当容器退出时，terminationMessagePath没有内容。
	TerminationMessageFallbackToLogsOnError TerminationMessagePolicy = "FallbackToLogsOnError"
)

// Capability represent POSIX capabilities type
// Capability 表示POSIX能力类型
type Capability string

// Adds and removes POSIX capabilities from running containers.
// 添加和删除正在运行的容器的POSIX能力。
type Capabilities struct {
	// Added capabilities
	// +optional
	Add []Capability `json:"add,omitempty" protobuf:"bytes,1,rep,name=add,casttype=Capability"`
	// Removed capabilities
	// +optional
	Drop []Capability `json:"drop,omitempty" protobuf:"bytes,2,rep,name=drop,casttype=Capability"`
}

// ResourceRequirements describes the compute resource requirements.
// ResourceRequirements 描述计算资源需求
type ResourceRequirements struct {
	// Limits describes the maximum amount of compute resources allowed.
	// More info: https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/
	// +optional
	// Limits 描述允许的最大计算资源量
	Limits ResourceList `json:"limits,omitempty" protobuf:"bytes,1,rep,name=limits,casttype=ResourceList,castkey=ResourceName"`
	// Requests describes the minimum amount of compute resources required.
	// If Requests is omitted for a container, it defaults to Limits if that is explicitly specified,
	// otherwise to an implementation-defined value.
	// More info: https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/
	// +optional
	// Requests 描述所需的最小计算资源量, 如果 Requests 被省略, 它将默认为 Limits, 如果 Limits 被明确指定, 否则将默认为实现定义的值
	Requests ResourceList `json:"requests,omitempty" protobuf:"bytes,2,rep,name=requests,casttype=ResourceList,castkey=ResourceName"`
	// Claims lists the names of resources, defined in spec.resourceClaims,
	// that are used by this container.
	//
	// This is an alpha field and requires enabling the
	// DynamicResourceAllocation feature gate.
	//
	// This field is immutable.
	//
	// +listType=map
	// +listMapKey=name
	// +featureGate=DynamicResourceAllocation
	// +optional
	// Claims 列出了由此容器使用的资源的名称, 这些资源在 spec.resourceClaims 中定义
	// 这是一个 alpha 字段, 需要启用 DynamicResourceAllocation 特性开关
	Claims []ResourceClaim `json:"claims,omitempty" protobuf:"bytes,3,opt,name=claims"`
}

// ResourceClaim references one entry in PodSpec.ResourceClaims.
// ResourceClaim 引用 PodSpec.ResourceClaims 中的一个条目
type ResourceClaim struct {
	// Name must match the name of one entry in pod.spec.resourceClaims of
	// the Pod where this field is used. It makes that resource available
	// inside a container.
	// Name 必须与 Pod 中 pod.spec.resourceClaims 中的一个条目的名称匹配, 它使得该资源在容器内可用
	Name string `json:"name" protobuf:"bytes,1,opt,name=name"`
}

const (
	// TerminationMessagePathDefault means the default path to capture the application termination message running in a container
	TerminationMessagePathDefault string = "/dev/termination-log"
)

// A single application container that you want to run within a pod.
// 你想在pod中运行的单个应用程序容器
type Container struct {
	// Name of the container specified as a DNS_LABEL.
	// Each container in a pod must have a unique name (DNS_LABEL).
	// Cannot be updated.
	// 指定为DNS_LABEL的容器名称。pod中的每个容器必须有唯一的名称(DNS_LABEL)。不能更新。
	Name string `json:"name" protobuf:"bytes,1,opt,name=name"`
	// Container image name.
	// More info: https://kubernetes.io/docs/concepts/containers/images
	// This field is optional to allow higher level config management to default or override
	// container images in workload controllers like Deployments and StatefulSets.
	// +optional
	// 容器镜像名称。这个字段是可选的，以允许更高级别的配置管理默认或覆盖工作负载控制器中的容器镜像，如Deployments和StatefulSets。
	Image string `json:"image,omitempty" protobuf:"bytes,2,opt,name=image"`
	// Entrypoint array. Not executed within a shell.
	// The container image's ENTRYPOINT is used if this is not provided.
	// Variable references $(VAR_NAME) are expanded using the container's environment. If a variable
	// cannot be resolved, the reference in the input string will be unchanged. Double $$ are reduced
	// to a single $, which allows for escaping the $(VAR_NAME) syntax: i.e. "$$(VAR_NAME)" will
	// produce the string literal "$(VAR_NAME)". Escaped references will never be expanded, regardless
	// of whether the variable exists or not. Cannot be updated.
	// More info: https://kubernetes.io/docs/tasks/inject-data-application/define-command-argument-container/#running-a-command-in-a-shell
	// +optional
	// 入口点数组。不在shell中执行。如果未提供，则使用容器镜像的ENTRYPOINT。使用容器的环境变量扩展变量引用$(VAR_NAME)。如果无法解析变量，则输入字符串中的引用将保持不变。双$$被减少为单个$，这允许对$(VAR_NAME)语法进行转义：即“$$(VAR_NAME)”将产生字符串文字“$(VAR_NAME)”。转义引用永远不会被扩展，无论变量是否存在。
	Command []string `json:"command,omitempty" protobuf:"bytes,3,rep,name=command"`
	// Arguments to the entrypoint.
	// The container image's CMD is used if this is not provided.
	// Variable references $(VAR_NAME) are expanded using the container's environment. If a variable
	// cannot be resolved, the reference in the input string will be unchanged. Double $$ are reduced
	// to a single $, which allows for escaping the $(VAR_NAME) syntax: i.e. "$$(VAR_NAME)" will
	// produce the string literal "$(VAR_NAME)". Escaped references will never be expanded, regardless
	// of whether the variable exists or not. Cannot be updated.
	// More info: https://kubernetes.io/docs/tasks/inject-data-application/define-command-argument-container/#running-a-command-in-a-shell
	// +optional
	// 入口点的参数。如果未提供，则使用容器镜像的CMD。使用容器的环境变量扩展变量引用$(VAR_NAME)。如果无法解析变量，则输入字符串中的引用将保持不变。双$$被减少为单个$，这允许对$(VAR_NAME)语法进行转义：即“$$(VAR_NAME)”将产生字符串文字“$(VAR_NAME)”。转义引用永远不会被扩展，无论变量是否存在。
	Args []string `json:"args,omitempty" protobuf:"bytes,4,rep,name=args"`
	// Container's working directory.
	// If not specified, the container runtime's default will be used, which
	// might be configured in the container image.
	// Cannot be updated.
	// +optional
	// 容器的工作目录。如果未指定，则将使用容器运行时的默认值，该值可能在容器镜像中进行了配置。
	WorkingDir string `json:"workingDir,omitempty" protobuf:"bytes,5,opt,name=workingDir"`
	// List of ports to expose from the container. Not specifying a port here
	// DOES NOT prevent that port from being exposed. Any port which is
	// listening on the default "0.0.0.0" address inside a container will be
	// accessible from the network.
	// Modifying this array with strategic merge patch may corrupt the data.
	// For more information See https://github.com/kubernetes/kubernetes/issues/108255.
	// Cannot be updated.
	// +optional
	// +patchMergeKey=containerPort
	// +patchStrategy=merge
	// +listType=map
	// +listMapKey=containerPort
	// +listMapKey=protocol
	// 要从容器公开的端口列表。此处未指定端口并不会阻止该端口被暴露。在容器内监听默认“0.0.0.0”地址的任何端口都可以从网络访问。
	Ports []ContainerPort `json:"ports,omitempty" patchStrategy:"merge" patchMergeKey:"containerPort" protobuf:"bytes,6,rep,name=ports"`
	// List of sources to populate environment variables in the container.
	// The keys defined within a source must be a C_IDENTIFIER. All invalid keys
	// will be reported as an event when the container is starting. When a key exists in multiple
	// sources, the value associated with the last source will take precedence.
	// Values defined by an Env with a duplicate key will take precedence.
	// Cannot be updated.
	// +optional
	// 用于在容器中填充环境变量的源列表。 源中定义的键必须是C_IDENTIFIER。 所有无效的键都将在容器启动时报告为事件。 当键存在于多个源中时，与最后一个源关联的值将具有优先权。 由具有重复键的Env定义的值将具有优先权。
	EnvFrom []EnvFromSource `json:"envFrom,omitempty" protobuf:"bytes,19,rep,name=envFrom"`
	// List of environment variables to set in the container.
	// Cannot be updated.
	// +optional
	// +patchMergeKey=name
	// +patchStrategy=merge
	// 用于在容器中设置环境变量的列表。
	Env []EnvVar `json:"env,omitempty" patchStrategy:"merge" patchMergeKey:"name" protobuf:"bytes,7,rep,name=env"`
	// Compute Resources required by this container.
	// Cannot be updated.
	// More info: https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/
	// +optional
	// 此容器所需的计算资源。不能更新。
	Resources ResourceRequirements `json:"resources,omitempty" protobuf:"bytes,8,opt,name=resources"`
	// Pod volumes to mount into the container's filesystem.
	// Cannot be updated.
	// +optional
	// +patchMergeKey=mountPath
	// +patchStrategy=merge
	// 要挂载到容器文件系统中的Pod卷。不能更新。
	VolumeMounts []VolumeMount `json:"volumeMounts,omitempty" patchStrategy:"merge" patchMergeKey:"mountPath" protobuf:"bytes,9,rep,name=volumeMounts"`
	// volumeDevices is the list of block devices to be used by the container.
	// +patchMergeKey=devicePath
	// +patchStrategy=merge
	// +optional
	// volumeDevices是容器要使用的块设备的列表。
	VolumeDevices []VolumeDevice `json:"volumeDevices,omitempty" patchStrategy:"merge" patchMergeKey:"devicePath" protobuf:"bytes,21,rep,name=volumeDevices"`
	// Periodic probe of container liveness.
	// Container will be restarted if the probe fails.
	// Cannot be updated.
	// More info: https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle#container-probes
	// +optional
	// 容器活动性的周期性探测。如果探测失败，容器将被重新启动。不能更新。
	LivenessProbe *Probe `json:"livenessProbe,omitempty" protobuf:"bytes,10,opt,name=livenessProbe"`
	// Periodic probe of container service readiness.
	// Container will be removed from service endpoints if the probe fails.
	// Cannot be updated.
	// More info: https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle#container-probes
	// +optional
	// 容器服务就绪性的周期性探测。如果探测失败，容器将从服务端点中删除。不能更新。
	ReadinessProbe *Probe `json:"readinessProbe,omitempty" protobuf:"bytes,11,opt,name=readinessProbe"`
	// StartupProbe indicates that the Pod has successfully initialized.
	// If specified, no other probes are executed until this completes successfully.
	// If this probe fails, the Pod will be restarted, just as if the livenessProbe failed.
	// This can be used to provide different probe parameters at the beginning of a Pod's lifecycle,
	// when it might take a long time to load data or warm a cache, than during steady-state operation.
	// This cannot be updated.
	// More info: https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle#container-probes
	// +optional
	// StartupProbe指示Pod已成功初始化。如果指定，则在此成功完成之前不会执行任何其他探测。如果此探测失败，Pod将被重新启动，就像livenessProbe失败一样。这可用于在Pod生命周期的开始时提供不同的探测参数，因为在此期间可能需要很长时间来加载数据或预热缓存，而在稳定状态操作期间则不需要。这不能更新。
	StartupProbe *Probe `json:"startupProbe,omitempty" protobuf:"bytes,22,opt,name=startupProbe"`
	// Actions that the management system should take in response to container lifecycle events.
	// Cannot be updated.
	// +optional
	// 管理系统应采取的响应容器生命周期事件的操作。不能更新。
	Lifecycle *Lifecycle `json:"lifecycle,omitempty" protobuf:"bytes,12,opt,name=lifecycle"`
	// Optional: Path at which the file to which the container's termination message
	// will be written is mounted into the container's filesystem.
	// Message written is intended to be brief final status, such as an assertion failure message.
	// Will be truncated by the node if greater than 4096 bytes. The total message length across
	// all containers will be limited to 12kb.
	// Defaults to /dev/termination-log.
	// Cannot be updated.
	// +optional
	// 可选：将容器终止消息写入的文件的路径，该文件将被挂载到容器的文件系统中。写入的消息旨在是简短的最终状态，例如断言失败消息。如果大于4096字节，节点将被截断。所有容器的消息总长度将被限制为12kb。默认为/dev/termination-log。不能更新。
	TerminationMessagePath string `json:"terminationMessagePath,omitempty" protobuf:"bytes,13,opt,name=terminationMessagePath"`
	// Indicate how the termination message should be populated. File will use the contents of
	// terminationMessagePath to populate the container status message on both success and failure.
	// FallbackToLogsOnError will use the last chunk of container log output if the termination
	// message file is empty and the container exited with an error.
	// The log output is limited to 2048 bytes or 80 lines, whichever is smaller.
	// Defaults to File.
	// Cannot be updated.
	// +optional
	// 指示应如何填充终止消息。文件将使用terminationMessagePath的内容在成功和失败时填充容器状态消息。如果终止消息文件为空并且容器以错误退出，则FallbackToLogsOnError将使用容器日志输出的最后一块。日志输出限制为2048字节或80行，以较小者为准。默认为文件。不能更新。
	TerminationMessagePolicy TerminationMessagePolicy `json:"terminationMessagePolicy,omitempty" protobuf:"bytes,20,opt,name=terminationMessagePolicy,casttype=TerminationMessagePolicy"`
	// Image pull policy.
	// One of Always, Never, IfNotPresent.
	// Defaults to Always if :latest tag is specified, or IfNotPresent otherwise.
	// Cannot be updated.
	// More info: https://kubernetes.io/docs/concepts/containers/images#updating-images
	// +optional
	// 镜像拉取策略。Always，Never，IfNotPresent之一。如果指定了:latest标签，则默认为Always，否则为IfNotPresent。不能更新。
	ImagePullPolicy PullPolicy `json:"imagePullPolicy,omitempty" protobuf:"bytes,14,opt,name=imagePullPolicy,casttype=PullPolicy"`
	// SecurityContext defines the security options the container should be run with.
	// If set, the fields of SecurityContext override the equivalent fields of PodSecurityContext.
	// More info: https://kubernetes.io/docs/tasks/configure-pod-container/security-context/
	// +optional
	// SecurityContext定义容器应运行的安全选项。如果设置，则SecurityContext的字段将覆盖PodSecurityContext的相应字段。
	SecurityContext *SecurityContext `json:"securityContext,omitempty" protobuf:"bytes,15,opt,name=securityContext"`

	// Variables for interactive containers, these have very specialized use-cases (e.g. debugging)
	// and shouldn't be used for general purpose containers.

	// Whether this container should allocate a buffer for stdin in the container runtime. If this
	// is not set, reads from stdin in the container will always result in EOF.
	// Default is false.
	// +optional
	// 该容器是否应在容器运行时为stdin分配缓冲区。如果未设置，则容器中的stdin读取将始终导致EOF。默认为false。
	Stdin bool `json:"stdin,omitempty" protobuf:"varint,16,opt,name=stdin"`
	// Whether the container runtime should close the stdin channel after it has been opened by
	// a single attach. When stdin is true the stdin stream will remain open across multiple attach
	// sessions. If stdinOnce is set to true, stdin is opened on container start, is empty until the
	// first client attaches to stdin, and then remains open and accepts data until the client disconnects,
	// at which time stdin is closed and remains closed until the container is restarted. If this
	// flag is false, a container processes that reads from stdin will never receive an EOF.
	// Default is false
	// +optional
	// 容器运行时是否应在单个附加后打开stdin通道后关闭stdin通道。当stdin为true时，stdin流将在多个附加会话之间保持打开。如果将stdinOnce设置为true，则在容器启动时打开stdin，直到第一个客户端附加到stdin为空，然后保持打开并接受数据，直到客户端断开连接，此时stdin被关闭并保持关闭，直到容器重新启动。如果此标志为false，则从stdin读取的容器进程永远不会收到EOF。默认为false
	StdinOnce bool `json:"stdinOnce,omitempty" protobuf:"varint,17,opt,name=stdinOnce"`
	// Whether this container should allocate a TTY for itself, also requires 'stdin' to be true.
	// Default is false.
	// +optional
	// 该容器是否应为自身分配TTY，也需要“stdin”为true。默认为false。
	TTY bool `json:"tty,omitempty" protobuf:"varint,18,opt,name=tty"`
}

// ProbeHandler defines a specific action that should be taken in a probe.
// One and only one of the fields must be specified.
// ProbeHandler定义了探针中应该采取的特定操作。必须指定一个且只有一个字段。
type ProbeHandler struct {
	// Exec specifies the action to take.
	// +optional
	// Exec指定要采取的操作。
	Exec *ExecAction `json:"exec,omitempty" protobuf:"bytes,1,opt,name=exec"`
	// HTTPGet specifies the http request to perform.
	// +optional
	// HTTPGet指定要执行的http请求。
	HTTPGet *HTTPGetAction `json:"httpGet,omitempty" protobuf:"bytes,2,opt,name=httpGet"`
	// TCPSocket specifies an action involving a TCP port.
	// +optional
	// TCPSocket指定涉及TCP端口的操作。
	TCPSocket *TCPSocketAction `json:"tcpSocket,omitempty" protobuf:"bytes,3,opt,name=tcpSocket"`

	// GRPC specifies an action involving a GRPC port.
	// This is a beta field and requires enabling GRPCContainerProbe feature gate.
	// +featureGate=GRPCContainerProbe
	// +optional
	// GRPC 指定涉及 GRPC 端口的操作。
	GRPC *GRPCAction `json:"grpc,omitempty" protobuf:"bytes,4,opt,name=grpc"`
}

// LifecycleHandler defines a specific action that should be taken in a lifecycle
// hook. One and only one of the fields, except TCPSocket must be specified.
// LifecycleHandler定义了生命周期钩子中应该采取的特定操作。除TCPSocket之外，必须指定一个且只有一个字段。
type LifecycleHandler struct {
	// Exec specifies the action to take.
	// +optional
	// Exec指定要采取的操作。
	Exec *ExecAction `json:"exec,omitempty" protobuf:"bytes,1,opt,name=exec"`
	// HTTPGet specifies the http request to perform.
	// +optional
	// HTTPGet指定要执行的http请求。
	HTTPGet *HTTPGetAction `json:"httpGet,omitempty" protobuf:"bytes,2,opt,name=httpGet"`
	// Deprecated. TCPSocket is NOT supported as a LifecycleHandler and kept
	// for the backward compatibility. There are no validation of this field and
	// lifecycle hooks will fail in runtime when tcp handler is specified.
	// +optional
	// Deprecated. TCPSocket不支持作为LifecycleHandler，并保留用于向后兼容性。没有对此字段进行验证，当指定tcp处理程序时，生命周期钩子将在运行时失败。
	TCPSocket *TCPSocketAction `json:"tcpSocket,omitempty" protobuf:"bytes,3,opt,name=tcpSocket"`
}

// Lifecycle describes actions that the management system should take in response to container lifecycle
// events. For the PostStart and PreStop lifecycle handlers, management of the container blocks
// until the action is complete, unless the container process fails, in which case the handler is aborted.
// 生命周期描述了管理系统应采取的响应容器生命周期事件的操作。对于PostStart和PreStop生命周期处理程序，管理容器将阻塞，直到操作完成，除非容器进程失败，否则处理程序将被中止。
type Lifecycle struct {
	// PostStart is called immediately after a container is created. If the handler fails,
	// the container is terminated and restarted according to its restart policy.
	// Other management of the container blocks until the hook completes.
	// More info: https://kubernetes.io/docs/concepts/containers/container-lifecycle-hooks/#container-hooks
	// +optional
	// PostStart在创建容器后立即调用。如果处理程序失败，则根据其重启策略终止并重新启动容器。其他容器管理将阻塞，直到钩子完成。
	PostStart *LifecycleHandler `json:"postStart,omitempty" protobuf:"bytes,1,opt,name=postStart"`
	// PreStop is called immediately before a container is terminated due to an
	// API request or management event such as liveness/startup probe failure,
	// preemption, resource contention, etc. The handler is not called if the
	// container crashes or exits. The Pod's termination grace period countdown begins before the
	// PreStop hook is executed. Regardless of the outcome of the handler, the
	// container will eventually terminate within the Pod's termination grace
	// period (unless delayed by finalizers). Other management of the container blocks until the hook completes
	// or until the termination grace period is reached.
	// More info: https://kubernetes.io/docs/concepts/containers/container-lifecycle-hooks/#container-hooks
	// +optional
	// PreStop在由于API请求或管理事件（例如活动性/启动探针失败，抢占，资源争用等）而终止容器之前立即调用。如果容器崩溃或退出，则不调用处理程序。 Pod的终止宽限期倒计时在执行PreStop钩子之前开始。无论处理程序的结果如何，容器最终将在Pod的终止宽限期内终止（除非被finalizers延迟）。其他容器管理将阻塞，直到钩子完成或直到达到终止宽限期。
	PreStop *LifecycleHandler `json:"preStop,omitempty" protobuf:"bytes,2,opt,name=preStop"`
}

type ConditionStatus string

// These are valid condition statuses. "ConditionTrue" means a resource is in the condition.
// "ConditionFalse" means a resource is not in the condition. "ConditionUnknown" means kubernetes
// can't decide if a resource is in the condition or not. In the future, we could add other
// intermediate conditions, e.g. ConditionDegraded.
// 这些是有效的条件状态。 "ConditionTrue"表示资源处于条件中。 "ConditionFalse"表示资源不在条件中。 "ConditionUnknown"表示kubernetes无法确定资源是否处于条件中。将来，我们可以添加其他中间条件，例如ConditionDegraded。
const (
	ConditionTrue    ConditionStatus = "True"
	ConditionFalse   ConditionStatus = "False"
	ConditionUnknown ConditionStatus = "Unknown"
)

// ContainerStateWaiting is a waiting state of a container.
// ContainerStateWaiting是容器的等待状态。
type ContainerStateWaiting struct {
	// (brief) reason the container is not yet running.
	// +optional
	// 容器尚未运行的原因（简短）。
	Reason string `json:"reason,omitempty" protobuf:"bytes,1,opt,name=reason"`
	// Message regarding why the container is not yet running.
	// +optional
	// 关于容器尚未运行的消息。
	Message string `json:"message,omitempty" protobuf:"bytes,2,opt,name=message"`
}

// ContainerStateRunning is a running state of a container.
// ContainerStateRunning是容器的运行状态。
type ContainerStateRunning struct {
	// Time at which the container was last (re-)started
	// +optional
	// 容器上次（重新）启动的时间
	StartedAt metav1.Time `json:"startedAt,omitempty" protobuf:"bytes,1,opt,name=startedAt"`
}

// ContainerStateTerminated is a terminated state of a container.
// ContainerStateTerminated是容器的终止状态。
type ContainerStateTerminated struct {
	// Exit status from the last termination of the container
	// 容器上次终止时的退出状态
	ExitCode int32 `json:"exitCode" protobuf:"varint,1,opt,name=exitCode"`
	// Signal from the last termination of the container
	// +optional
	// 容器上次终止时的信号
	Signal int32 `json:"signal,omitempty" protobuf:"varint,2,opt,name=signal"`
	// (brief) reason from the last termination of the container
	// +optional
	// (简要)集装箱最后一次终止的原因
	Reason string `json:"reason,omitempty" protobuf:"bytes,3,opt,name=reason"`
	// Message regarding the last termination of the container
	// +optional
	// 关于容器最后一次终止的消息
	Message string `json:"message,omitempty" protobuf:"bytes,4,opt,name=message"`
	// Time at which previous execution of the container started
	// +optional
	// 容器上次执行开始的时间
	StartedAt metav1.Time `json:"startedAt,omitempty" protobuf:"bytes,5,opt,name=startedAt"`
	// Time at which the container last terminated
	// +optional
	// 容器上次终止的时间
	FinishedAt metav1.Time `json:"finishedAt,omitempty" protobuf:"bytes,6,opt,name=finishedAt"`
	// Container's ID in the format '<type>://<container_id>'
	// +optional
	// 容器的ID格式为'<type>://<container_id>'
	ContainerID string `json:"containerID,omitempty" protobuf:"bytes,7,opt,name=containerID"`
}

// ContainerState holds a possible state of container.
// Only one of its members may be specified.
// If none of them is specified, the default one is ContainerStateWaiting.
// ContainerState包含容器的可能状态。只能指定其中的一个成员。如果没有指定，那么默认值是ContainerStateWaiting。
type ContainerState struct {
	// Details about a waiting container
	// +optional
	// 关于等待容器的详细信息
	Waiting *ContainerStateWaiting `json:"waiting,omitempty" protobuf:"bytes,1,opt,name=waiting"`
	// Details about a running container
	// +optional
	// 关于正在运行的容器的详细信息
	Running *ContainerStateRunning `json:"running,omitempty" protobuf:"bytes,2,opt,name=running"`
	// Details about a terminated container
	// +optional
	// 关于终止容器的详细信息
	Terminated *ContainerStateTerminated `json:"terminated,omitempty" protobuf:"bytes,3,opt,name=terminated"`
}

// ContainerStatus contains details for the current status of this container.
// ContainerStatus包含此容器当前状态的详细信息。
type ContainerStatus struct {
	// This must be a DNS_LABEL. Each container in a pod must have a unique name.
	// Cannot be updated.
	// 必须是DNS_LABEL。每个pod中的容器必须具有唯一的名称。不能更新。
	Name string `json:"name" protobuf:"bytes,1,opt,name=name"`
	// Details about the container's current condition.
	// +optional
	// 关于容器当前条件的详细信息。
	State ContainerState `json:"state,omitempty" protobuf:"bytes,2,opt,name=state"`
	// Details about the container's last termination condition.
	// +optional
	// 关于容器上次终止条件的详细信息。
	LastTerminationState ContainerState `json:"lastState,omitempty" protobuf:"bytes,3,opt,name=lastState"`
	// Specifies whether the container has passed its readiness probe.
	// 指定容器是否已通过其就绪探测。
	Ready bool `json:"ready" protobuf:"varint,4,opt,name=ready"`
	// The number of times the container has been restarted.
	// 容器重新启动的次数
	RestartCount int32 `json:"restartCount" protobuf:"varint,5,opt,name=restartCount"`
	// The image the container is running.
	// More info: https://kubernetes.io/docs/concepts/containers/images.
	// 容器正在运行的镜像。更多信息：https://kubernetes.io/docs/concepts/containers/images。
	Image string `json:"image" protobuf:"bytes,6,opt,name=image"`
	// ImageID of the container's image.
	ImageID string `json:"imageID" protobuf:"bytes,7,opt,name=imageID"`
	// Container's ID in the format '<type>://<container_id>'.
	// +optional
	// 容器的ID格式为'<type>://<container_id>'。
	ContainerID string `json:"containerID,omitempty" protobuf:"bytes,8,opt,name=containerID"`
	// Specifies whether the container has passed its startup probe.
	// Initialized as false, becomes true after startupProbe is considered successful.
	// Resets to false when the container is restarted, or if kubelet loses state temporarily.
	// Is always true when no startupProbe is defined.
	// +optional
	// 指定容器是否已通过其启动探测。初始化为false，启动探测被认为是成功的后变为true。当容器重新启动或kubelet暂时丢失状态时，将重置为false。当未定义启动探测时，始终为true。
	Started *bool `json:"started,omitempty" protobuf:"varint,9,opt,name=started"`
}

// PodPhase is a label for the condition of a pod at the current time.
// +enum
// PodPhase是当前时间pod状态的标签。
type PodPhase string

// These are the valid statuses of pods.
// 这些是pod的有效状态。
const (
	// PodPending means the pod has been accepted by the system, but one or more of the containers
	// has not been started. This includes time before being bound to a node, as well as time spent
	// pulling images onto the host.
	// PodPending表示pod已被系统接受，但尚未启动一个或多个容器。这包括绑定到节点之前的时间，以及将镜像拉取到主机上的时间。
	PodPending PodPhase = "Pending"
	// PodRunning means the pod has been bound to a node and all of the containers have been started.
	// At least one container is still running or is in the process of being restarted.
	// PodRunning表示pod已绑定到节点，并且已启动所有容器。至少有一个容器仍在运行或正在重新启动。
	PodRunning PodPhase = "Running"
	// PodSucceeded means that all containers in the pod have voluntarily terminated
	// with a container exit code of 0, and the system is not going to restart any of these containers.
	// PodSucceeded意味着pod中的所有容器都已自愿终止，容器退出代码为0，并且系统不会重新启动这些容器。
	PodSucceeded PodPhase = "Succeeded"
	// PodFailed means that all containers in the pod have terminated, and at least one container has
	// terminated in a failure (exited with a non-zero exit code or was stopped by the system).
	// PodFailed表示pod中的所有容器都已终止，并且至少有一个容器已失败终止（退出代码为非零或被系统停止）。
	PodFailed PodPhase = "Failed"
	// PodUnknown means that for some reason the state of the pod could not be obtained, typically due
	// to an error in communicating with the host of the pod.
	// Deprecated: It isn't being set since 2015 (74da3b14b0c0f658b3bb8d2def5094686d0e9095)
	// PodUnknown表示由于某种原因无法获取pod的状态，通常是由于与pod主机通信时出错。
	PodUnknown PodPhase = "Unknown"
)

// PodConditionType is a valid value for PodCondition.Type
// PodConditionType 是 PodCondition.Type 的有效值
type PodConditionType string

// These are built-in conditions of pod. An application may use a custom condition not listed here.
// 这些是内置的pod条件。应用程序可以使用这里没有列出的自定义条件。
const (
	// ContainersReady indicates whether all containers in the pod are ready.
	// ContainersReady 表示 pod 中的所有容器是否就绪。
	ContainersReady PodConditionType = "ContainersReady"
	// PodInitialized means that all init containers in the pod have started successfully.
	// PodInitialized意味着pod中的所有init容器都已成功启动。
	PodInitialized PodConditionType = "Initialized"
	// PodReady means the pod is able to service requests and should be added to the
	// load balancing pools of all matching services.
	// PodReady表示pod能够提供服务请求，并且应该添加到所有匹配服务的负载均衡池中。
	PodReady PodConditionType = "Ready"
	// PodScheduled represents status of the scheduling process for this pod.
	// PodScheduled表示此pod的调度过程的状态。
	PodScheduled PodConditionType = "PodScheduled"
	// DisruptionTarget indicates the pod is about to be terminated due to a
	// disruption (such as preemption, eviction API or garbage-collection).
	// DisruptionTarget表示由于干扰（例如抢占，驱逐API或垃圾收集）而将终止pod。
	DisruptionTarget PodConditionType = "DisruptionTarget"
)

// These are reasons for a pod's transition to a condition.
// 这些是pod转换为条件的原因。
const (
	// PodReasonUnschedulable reason in PodScheduled PodCondition means that the scheduler
	// can't schedule the pod right now, for example due to insufficient resources in the cluster.
	// PodReasonUnschedulable PodScheduled PodCondition的原因是调度程序现在无法调度pod，例如由于集群中的资源不足。
	PodReasonUnschedulable = "Unschedulable"

	// PodReasonSchedulingGated reason in PodScheduled PodCondition means that the scheduler
	// skips scheduling the pod because one or more scheduling gates are still present.
	// PodReasonSchedulingGated PodScheduled PodCondition的原因是调度程序跳过调度pod，因为仍然存在一个或多个调度门。
	PodReasonSchedulingGated = "SchedulingGated"

	// PodReasonSchedulerError reason in PodScheduled PodCondition means that some internal error happens
	// during scheduling, for example due to nodeAffinity parsing errors.
	PodReasonSchedulerError = "SchedulerError"

	// TerminationByKubelet reason in DisruptionTarget pod condition indicates that the termination
	// is initiated by kubelet
	PodReasonTerminationByKubelet = "TerminationByKubelet"
)

// PodCondition contains details for the current condition of this pod.
// PodCondition包含此pod当前条件的详细信息。
type PodCondition struct {
	// Type is the type of the condition.
	// More info: https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle#pod-conditions
	// Type是条件的类型。
	Type PodConditionType `json:"type" protobuf:"bytes,1,opt,name=type,casttype=PodConditionType"`
	// Status is the status of the condition.
	// Can be True, False, Unknown.
	// More info: https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle#pod-conditions
	// Status是条件的状态。
	Status ConditionStatus `json:"status" protobuf:"bytes,2,opt,name=status,casttype=ConditionStatus"`
	// Last time we probed the condition.
	// +optional
	// LastProbeTime是我们探测条件的最后一次时间。
	LastProbeTime metav1.Time `json:"lastProbeTime,omitempty" protobuf:"bytes,3,opt,name=lastProbeTime"`
	// Last time the condition transitioned from one status to another.
	// +optional
	// LastTransitionTime是条件从一个状态转换为另一个状态的最后一次时间。
	LastTransitionTime metav1.Time `json:"lastTransitionTime,omitempty" protobuf:"bytes,4,opt,name=lastTransitionTime"`
	// Unique, one-word, CamelCase reason for the condition's last transition.
	// +optional
	// Reason是条件最后一次转换的唯一，单词，CamelCase原因。
	Reason string `json:"reason,omitempty" protobuf:"bytes,5,opt,name=reason"`
	// Human-readable message indicating details about last transition.
	// +optional
	// Message是指示有关最后一次转换的详细信息的人类可读消息。
	Message string `json:"message,omitempty" protobuf:"bytes,6,opt,name=message"`
}

// RestartPolicy describes how the container should be restarted.
// Only one of the following restart policies may be specified.
// If none of the following policies is specified, the default one
// is RestartPolicyAlways.
// +enum
// RestartPolicy 描述了容器应该如何重新启动。
// 只能指定以下重启策略之一。
// 如果未指定以下任何策略，则默认为 RestartPolicyAlways。
type RestartPolicy string

const (
	RestartPolicyAlways    RestartPolicy = "Always"
	RestartPolicyOnFailure RestartPolicy = "OnFailure"
	RestartPolicyNever     RestartPolicy = "Never"
)

// DNSPolicy defines how a pod's DNS will be configured.
// +enum
// DNSPolicy定义了pod的DNS将如何配置。
type DNSPolicy string

const (
	// DNSClusterFirstWithHostNet indicates that the pod should use cluster DNS
	// first, if it is available, then fall back on the default
	// (as determined by kubelet) DNS settings.
	// DNSClusterFirstWithHostNet表示如果可用，pod应该首先使用集群DNS，然后再使用默认值（由kubelet确定的DNS设置）。
	DNSClusterFirstWithHostNet DNSPolicy = "ClusterFirstWithHostNet"

	// DNSClusterFirst indicates that the pod should use cluster DNS
	// first unless hostNetwork is true, if it is available, then
	// fall back on the default (as determined by kubelet) DNS settings.
	// DNSClusterFirst表示如果可用，pod应该首先使用集群DNS，除非hostNetwork为true，然后再使用默认值（由kubelet确定的DNS设置）。
	DNSClusterFirst DNSPolicy = "ClusterFirst"

	// DNSDefault indicates that the pod should use the default (as
	// determined by kubelet) DNS settings.
	// DNSDefault表示pod应该使用默认值（由kubelet确定的DNS设置）。
	DNSDefault DNSPolicy = "Default"

	// DNSNone indicates that the pod should use empty DNS settings. DNS
	// parameters such as nameservers and search paths should be defined via
	// DNSConfig.
	// DNSNone表示pod应该使用空的DNS设置。 DNS参数（如名称服务器和搜索路径）应通过DNSConfig定义。
	DNSNone DNSPolicy = "None"
)

const (
	// DefaultTerminationGracePeriodSeconds indicates the default duration in
	// seconds a pod needs to terminate gracefully.
	DefaultTerminationGracePeriodSeconds = 30
)

// A node selector represents the union of the results of one or more label queries
// over a set of nodes; that is, it represents the OR of the selectors represented
// by the node selector terms.
// +structType=atomic
// NodeSelector表示一组节点上一个或多个标签查询的结果的并集；也就是说，它表示由节点选择器术语表示的选择器的OR。
type NodeSelector struct {
	// Required. A list of node selector terms. The terms are ORed.
	// Required。 一个节点选择器术语列表。 术语是OR的。
	NodeSelectorTerms []NodeSelectorTerm `json:"nodeSelectorTerms" protobuf:"bytes,1,rep,name=nodeSelectorTerms"`
}

// A null or empty node selector term matches no objects. The requirements of
// them are ANDed.
// The TopologySelectorTerm type implements a subset of the NodeSelectorTerm.
// +structType=atomic
// 空或空节点选择器术语不匹配任何对象。 它们的要求是AND的。
// TopologySelectorTerm类型实现了NodeSelectorTerm的子集。
type NodeSelectorTerm struct {
	// A list of node selector requirements by node's labels.
	// +optional
	// 通过节点标签的节点选择器要求列表。
	MatchExpressions []NodeSelectorRequirement `json:"matchExpressions,omitempty" protobuf:"bytes,1,rep,name=matchExpressions"`
	// A list of node selector requirements by node's fields.
	// +optional
	// 通过节点字段的节点选择器要求列表。
	MatchFields []NodeSelectorRequirement `json:"matchFields,omitempty" protobuf:"bytes,2,rep,name=matchFields"`
}

// A node selector requirement is a selector that contains values, a key, and an operator
// that relates the key and values.
// 一个节点选择器要求是一个包含值，键和与键和值相关的运算符的选择器。
type NodeSelectorRequirement struct {
	// The label key that the selector applies to.
	// 选择器适用的标签键。
	Key string `json:"key" protobuf:"bytes,1,opt,name=key"`
	// Represents a key's relationship to a set of values.
	// Valid operators are In, NotIn, Exists, DoesNotExist. Gt, and Lt.
	// 代表键与一组值的关系。 有效的运算符是In，NotIn，Exists，DoesNotExist。 Gt和Lt。
	Operator NodeSelectorOperator `json:"operator" protobuf:"bytes,2,opt,name=operator,casttype=NodeSelectorOperator"`
	// An array of string values. If the operator is In or NotIn,
	// the values array must be non-empty. If the operator is Exists or DoesNotExist,
	// the values array must be empty. If the operator is Gt or Lt, the values
	// array must have a single element, which will be interpreted as an integer.
	// This array is replaced during a strategic merge patch.
	// +optional
	// 字符串值数组。 如果运算符是In或NotIn，则values数组必须非空。 如果运算符是Exists或DoesNotExist，则values数组必须为空。 如果运算符是Gt或Lt，则values数组必须有一个元素，该元素将被解释为整数。 在战略合并补丁期间，将替换此数组。
	Values []string `json:"values,omitempty" protobuf:"bytes,3,rep,name=values"`
}

// A node selector operator is the set of operators that can be used in
// a node selector requirement.
// +enum
// 节点选择器运算符是可以在节点选择器要求中使用的运算符集。
type NodeSelectorOperator string

const (
	NodeSelectorOpIn           NodeSelectorOperator = "In"
	NodeSelectorOpNotIn        NodeSelectorOperator = "NotIn"
	NodeSelectorOpExists       NodeSelectorOperator = "Exists"
	NodeSelectorOpDoesNotExist NodeSelectorOperator = "DoesNotExist"
	NodeSelectorOpGt           NodeSelectorOperator = "Gt"
	NodeSelectorOpLt           NodeSelectorOperator = "Lt"
)

// A topology selector term represents the result of label queries.
// A null or empty topology selector term matches no objects.
// The requirements of them are ANDed.
// It provides a subset of functionality as NodeSelectorTerm.
// This is an alpha feature and may change in the future.
// +structType=atomic
// 拓扑选择器术语表示标签查询的结果。
// 空或空的拓扑选择器术语不匹配任何对象。
// 它们的要求是AND的。
// 它提供了NodeSelectorTerm的功能子集。
// 这是一个alpha功能，以后可能会改变。
type TopologySelectorTerm struct {
	// Usage: Fields of type []TopologySelectorTerm must be listType=atomic.

	// A list of topology selector requirements by labels.
	// +optional
	// 通过标签的拓扑选择器要求列表。
	MatchLabelExpressions []TopologySelectorLabelRequirement `json:"matchLabelExpressions,omitempty" protobuf:"bytes,1,rep,name=matchLabelExpressions"`
}

// A topology selector requirement is a selector that matches given label.
// This is an alpha feature and may change in the future.
// 一个拓扑选择器要求是一个选择器，它匹配给定的标签。
// 这是一个alpha功能，以后可能会改变。
type TopologySelectorLabelRequirement struct {
	// The label key that the selector applies to.
	// 选择器适用的标签键。
	Key string `json:"key" protobuf:"bytes,1,opt,name=key"`
	// An array of string values. One value must match the label to be selected.
	// Each entry in Values is ORed.
	// 一个字符串值数组。 一个值必须匹配所选标签。 Values中的每个条目都是ORed。
	Values []string `json:"values" protobuf:"bytes,2,rep,name=values"`
}

// Affinity is a group of affinity scheduling rules.
// Affinity是一组亲和性调度规则。
type Affinity struct {
	// Describes node affinity scheduling rules for the pod.
	// +optional
	// 描述Pod的节点亲和性调度规则。
	NodeAffinity *NodeAffinity `json:"nodeAffinity,omitempty" protobuf:"bytes,1,opt,name=nodeAffinity"`
	// Describes pod affinity scheduling rules (e.g. co-locate this pod in the same node, zone, etc. as some other pod(s)).
	// +optional
	// 描述Pod亲和性调度规则（例如，将此Pod与某些其他Pod（s）放在同一节点，区域等中）。
	PodAffinity *PodAffinity `json:"podAffinity,omitempty" protobuf:"bytes,2,opt,name=podAffinity"`
	// Describes pod anti-affinity scheduling rules (e.g. avoid putting this pod in the same node, zone, etc. as some other pod(s)).
	// +optional
	// 描述Pod反亲和性调度规则（例如，避免将此Pod与某些其他Pod（s）放在同一节点，区域等中）。
	PodAntiAffinity *PodAntiAffinity `json:"podAntiAffinity,omitempty" protobuf:"bytes,3,opt,name=podAntiAffinity"`
}

// Pod affinity is a group of inter pod affinity scheduling rules.
// Pod亲和性是一组Pod间亲和性调度规则。
type PodAffinity struct {
	// NOT YET IMPLEMENTED. TODO: Uncomment field once it is implemented.
	// If the affinity requirements specified by this field are not met at
	// scheduling time, the pod will not be scheduled onto the node.
	// If the affinity requirements specified by this field cease to be met
	// at some point during pod execution (e.g. due to a pod label update), the
	// system will try to eventually evict the pod from its node.
	// When there are multiple elements, the lists of nodes corresponding to each
	// podAffinityTerm are intersected, i.e. all terms must be satisfied.
	// +optional
	// RequiredDuringSchedulingRequiredDuringExecution []PodAffinityTerm  `json:"requiredDuringSchedulingRequiredDuringExecution,omitempty"`

	// If the affinity requirements specified by this field are not met at
	// scheduling time, the pod will not be scheduled onto the node.
	// If the affinity requirements specified by this field cease to be met
	// at some point during pod execution (e.g. due to a pod label update), the
	// system may or may not try to eventually evict the pod from its node.
	// When there are multiple elements, the lists of nodes corresponding to each
	// podAffinityTerm are intersected, i.e. all terms must be satisfied.
	// +optional
	// 如果在调度时间指定此字段的亲和性要求未满足，则Pod将不会被调度到节点上。
	// 如果在Pod执行期间（例如由于Pod标签更新）某个时间点不再满足此字段指定的亲和性要求，则系统将尝试最终将Pod从其节点上驱逐。
	// 当有多个元素时，对应于每个podAffinityTerm的节点列表相交，即所有条款必须满足。
	RequiredDuringSchedulingIgnoredDuringExecution []PodAffinityTerm `json:"requiredDuringSchedulingIgnoredDuringExecution,omitempty" protobuf:"bytes,1,rep,name=requiredDuringSchedulingIgnoredDuringExecution"`
	// The scheduler will prefer to schedule pods to nodes that satisfy
	// the affinity expressions specified by this field, but it may choose
	// a node that violates one or more of the expressions. The node that is
	// most preferred is the one with the greatest sum of weights, i.e.
	// for each node that meets all of the scheduling requirements (resource
	// request, requiredDuringScheduling affinity expressions, etc.),
	// compute a sum by iterating through the elements of this field and adding
	// "weight" to the sum if the node has pods which matches the corresponding podAffinityTerm; the
	// node(s) with the highest sum are the most preferred.
	// +optional
	// 调度程序将更喜欢将Pod调度到满足此字段指定的亲和性表达式的节点，但它可以选择违反一个或多个表达式的节点。
	// 最喜欢的节点是权重之和最大的节点，即对于满足所有调度要求（资源请求，requiredDuringScheduling亲和性表达式等）的每个节点，
	// 通过迭代此字段的元素并在节点具有与相应podAffinityTerm匹配的Pod时将“weight”添加到总和来计算总和来计算总和；
	// 权重之和最大的节点是最喜欢的节点。
	PreferredDuringSchedulingIgnoredDuringExecution []WeightedPodAffinityTerm `json:"preferredDuringSchedulingIgnoredDuringExecution,omitempty" protobuf:"bytes,2,rep,name=preferredDuringSchedulingIgnoredDuringExecution"`
}

// Pod anti affinity is a group of inter pod anti affinity scheduling rules.
// Pod反亲和性是一组Pod间反亲和性调度规则。
type PodAntiAffinity struct {
	// NOT YET IMPLEMENTED. TODO: Uncomment field once it is implemented.
	// If the anti-affinity requirements specified by this field are not met at
	// scheduling time, the pod will not be scheduled onto the node.
	// If the anti-affinity requirements specified by this field cease to be met
	// at some point during pod execution (e.g. due to a pod label update), the
	// system will try to eventually evict the pod from its node.
	// When there are multiple elements, the lists of nodes corresponding to each
	// podAffinityTerm are intersected, i.e. all terms must be satisfied.
	// +optional
	// RequiredDuringSchedulingRequiredDuringExecution []PodAffinityTerm  `json:"requiredDuringSchedulingRequiredDuringExecution,omitempty"`

	// If the anti-affinity requirements specified by this field are not met at
	// scheduling time, the pod will not be scheduled onto the node.
	// If the anti-affinity requirements specified by this field cease to be met
	// at some point during pod execution (e.g. due to a pod label update), the
	// system may or may not try to eventually evict the pod from its node.
	// When there are multiple elements, the lists of nodes corresponding to each
	// podAffinityTerm are intersected, i.e. all terms must be satisfied.
	// +optional
	// 如果在调度时间指定此字段的反亲和性要求未满足，则Pod将不会被调度到节点上。
	// 如果在Pod执行期间（例如由于Pod标签更新）某个时间点不再满足此字段指定的反亲和性要求，则系统将尝试最终将Pod从其节点上驱逐。
	// 当有多个元素时，对应于每个podAffinityTerm的节点列表相交，即所有条款必须满足。
	RequiredDuringSchedulingIgnoredDuringExecution []PodAffinityTerm `json:"requiredDuringSchedulingIgnoredDuringExecution,omitempty" protobuf:"bytes,1,rep,name=requiredDuringSchedulingIgnoredDuringExecution"`
	// The scheduler will prefer to schedule pods to nodes that satisfy
	// the anti-affinity expressions specified by this field, but it may choose
	// a node that violates one or more of the expressions. The node that is
	// most preferred is the one with the greatest sum of weights, i.e.
	// for each node that meets all of the scheduling requirements (resource
	// request, requiredDuringScheduling anti-affinity expressions, etc.),
	// compute a sum by iterating through the elements of this field and adding
	// "weight" to the sum if the node has pods which matches the corresponding podAffinityTerm; the
	// node(s) with the highest sum are the most preferred.
	// +optional
	// 调度程序将更喜欢将Pod调度到满足此字段指定的反亲和性表达式的节点，但它可以选择违反一个或多个表达式的节点。
	// 最喜欢的节点是权重之和最大的节点，即对于满足所有调度要求（资源请求，requiredDuringScheduling反亲和性表达式等）的每个节点，
	// 通过迭代此字段的元素并在节点具有与相应podAffinityTerm匹配的Pod时将“weight”添加到总和来计算总和来计算总和；
	// 权重之和最大的节点是最喜欢的节点。
	PreferredDuringSchedulingIgnoredDuringExecution []WeightedPodAffinityTerm `json:"preferredDuringSchedulingIgnoredDuringExecution,omitempty" protobuf:"bytes,2,rep,name=preferredDuringSchedulingIgnoredDuringExecution"`
}

// The weights of all of the matched WeightedPodAffinityTerm fields are added per-node to find the most preferred node(s)
// 为每个节点添加所有匹配的WeightedPodAffinityTerm字段的权重，以查找最优先的节点
type WeightedPodAffinityTerm struct {
	// weight associated with matching the corresponding podAffinityTerm,
	// in the range 1-100.
	// 权重与匹配相应的podAffinityTerm相关联，范围为1-100。
	Weight int32 `json:"weight" protobuf:"varint,1,opt,name=weight"`
	// Required. A pod affinity term, associated with the corresponding weight.
	// 必需。与相应权重相关联的Pod亲和性术语。
	PodAffinityTerm PodAffinityTerm `json:"podAffinityTerm" protobuf:"bytes,2,opt,name=podAffinityTerm"`
}

// Defines a set of pods (namely those matching the labelSelector
// relative to the given namespace(s)) that this pod should be
// co-located (affinity) or not co-located (anti-affinity) with,
// where co-located is defined as running on a node whose value of
// the label with key <topologyKey> matches that of any node on which
// a pod of the set of pods is running
// 定义一组Pod（即与给定命名空间（s）相对的标签选择器匹配的Pod），
// 应该与此Pod一起（亲和性）或不一起（反亲和性）共处，其中共处被定义为在其标签值为<topologyKey>的节点上运行，
// 该节点的值与运行该Pod集的任何节点上运行的Pod的值匹配。
type PodAffinityTerm struct {
	// A label query over a set of resources, in this case pods.
	// +optional
	// 对一组资源的标签查询，本例中为Pod。
	LabelSelector *metav1.LabelSelector `json:"labelSelector,omitempty" protobuf:"bytes,1,opt,name=labelSelector"`
	// namespaces specifies a static list of namespace names that the term applies to.
	// The term is applied to the union of the namespaces listed in this field
	// and the ones selected by namespaceSelector.
	// null or empty namespaces list and null namespaceSelector means "this pod's namespace".
	// +optional
	// namespaces指定了该术语适用的命名空间名称的静态列表。
	// 该术语适用于此字段中列出的命名空间的并集以及由namespaceSelector选择的命名空间。
	// 空或空的namespaces列表和null namespaceSelector表示“此Pod的命名空间”。
	Namespaces []string `json:"namespaces,omitempty" protobuf:"bytes,2,rep,name=namespaces"`
	// This pod should be co-located (affinity) or not co-located (anti-affinity) with the pods matching
	// the labelSelector in the specified namespaces, where co-located is defined as running on a node
	// whose value of the label with key topologyKey matches that of any node on which any of the
	// selected pods is running.
	// Empty topologyKey is not allowed.
	// 这个pod应该与指定名称空间中匹配labelSelector的pod共存(亲和性)或不共存(反亲和性)，其中共存被定义为运行在一个节点上，该节点上带有键topologyKey的标签的值与运行任何所选pod的任何节点的值相匹配。不允许为空topologyKey。
	TopologyKey string `json:"topologyKey" protobuf:"bytes,3,opt,name=topologyKey"`
	// A label query over the set of namespaces that the term applies to.
	// The term is applied to the union of the namespaces selected by this field
	// and the ones listed in the namespaces field.
	// null selector and null or empty namespaces list means "this pod's namespace".
	// An empty selector ({}) matches all namespaces.
	// +optional
	// 一个标签查询，用于查询该术语适用的命名空间集。
	// 该术语适用于此字段选择的命名空间的并集以及在namespaces字段中列出的命名空间。
	// null选择器和null或空的namespaces列表表示“此Pod的命名空间”。
	// 空选择器（{}）匹配所有命名空间。
	NamespaceSelector *metav1.LabelSelector `json:"namespaceSelector,omitempty" protobuf:"bytes,4,opt,name=namespaceSelector"`
}

// Node affinity is a group of node affinity scheduling rules.
// Node affinity是一组节点亲和性调度规则。
type NodeAffinity struct {
	// NOT YET IMPLEMENTED. TODO: Uncomment field once it is implemented.
	// If the affinity requirements specified by this field are not met at
	// scheduling time, the pod will not be scheduled onto the node.
	// If the affinity requirements specified by this field cease to be met
	// at some point during pod execution (e.g. due to an update), the system
	// will try to eventually evict the pod from its node.
	// +optional
	// RequiredDuringSchedulingRequiredDuringExecution *NodeSelector `json:"requiredDuringSchedulingRequiredDuringExecution,omitempty"`

	// If the affinity requirements specified by this field are not met at
	// scheduling time, the pod will not be scheduled onto the node.
	// If the affinity requirements specified by this field cease to be met
	// at some point during pod execution (e.g. due to an update), the system
	// may or may not try to eventually evict the pod from its node.
	// +optional
	// 如果在调度时间不满足此字段指定的亲和性要求，则Pod将不会被调度到该节点上。
	// 如果在Pod执行过程中的某个时间点（例如由于更新）不满足此字段指定的亲和性要求，则系统将尝试最终将Pod从其节点上驱逐。
	RequiredDuringSchedulingIgnoredDuringExecution *NodeSelector `json:"requiredDuringSchedulingIgnoredDuringExecution,omitempty" protobuf:"bytes,1,opt,name=requiredDuringSchedulingIgnoredDuringExecution"`
	// The scheduler will prefer to schedule pods to nodes that satisfy
	// the affinity expressions specified by this field, but it may choose
	// a node that violates one or more of the expressions. The node that is
	// most preferred is the one with the greatest sum of weights, i.e.
	// for each node that meets all of the scheduling requirements (resource
	// request, requiredDuringScheduling affinity expressions, etc.),
	// compute a sum by iterating through the elements of this field and adding
	// "weight" to the sum if the node matches the corresponding matchExpressions; the
	// node(s) with the highest sum are the most preferred.
	// +optional
	// 调度程序将更喜欢将Pod调度到满足此字段指定的亲和性表达式的节点上，但它可以选择违反一个或多个表达式的节点。最喜欢的节点是权重之和最大的节点，即对于满足所有调度要求（资源请求，requiredDuringScheduling亲和性表达式等）的每个节点，通过迭代此字段的元素并在节点与相应的matchExpressions匹配时将“weight”添加到总和来计算总和；具有最高总和的节点是最喜欢的。
	PreferredDuringSchedulingIgnoredDuringExecution []PreferredSchedulingTerm `json:"preferredDuringSchedulingIgnoredDuringExecution,omitempty" protobuf:"bytes,2,rep,name=preferredDuringSchedulingIgnoredDuringExecution"`
}

// An empty preferred scheduling term matches all objects with implicit weight 0
// (i.e. it's a no-op). A null preferred scheduling term matches no objects (i.e. is also a no-op).
// 一个空的首选调度术语与所有隐式权重为0的对象匹配（即它是一个无操作）。一个空的首选调度术语不匹配任何对象（即也是一个无操作）。
type PreferredSchedulingTerm struct {
	// Weight associated with matching the corresponding nodeSelectorTerm, in the range 1-100.
	// 与匹配相应的nodeSelectorTerm相关联的权重，范围为1-100。
	Weight int32 `json:"weight" protobuf:"varint,1,opt,name=weight"`
	// A node selector term, associated with the corresponding weight.
	// 与相应权重相关联的节点选择器术语。
	Preference NodeSelectorTerm `json:"preference" protobuf:"bytes,2,opt,name=preference"`
}

// The node this Taint is attached to has the "effect" on
// any pod that does not tolerate the Taint.
// 这个污点所附着的节点对任何不容忍污点的pod都有“效果”。
type Taint struct {
	// Required. The taint key to be applied to a node.
	Key string `json:"key" protobuf:"bytes,1,opt,name=key"`
	// The taint value corresponding to the taint key.
	// +optional
	// 与污点键相对应的污点值。
	Value string `json:"value,omitempty" protobuf:"bytes,2,opt,name=value"`
	// Required. The effect of the taint on pods
	// that do not tolerate the taint.
	// Valid effects are NoSchedule, PreferNoSchedule and NoExecute.
	// 需要的，污点对不容忍污点的pod的影响
	// 有效的影响是NoSchedule，PreferNoSchedule和NoExecute。
	Effect TaintEffect `json:"effect" protobuf:"bytes,3,opt,name=effect,casttype=TaintEffect"`
	// TimeAdded represents the time at which the taint was added.
	// It is only written for NoExecute taints.
	// +optional
	// TimeAdded表示添加污点的时间。 这只是为NoExecute污点写的。
	TimeAdded *metav1.Time `json:"timeAdded,omitempty" protobuf:"bytes,4,opt,name=timeAdded"`
}

// +enum
type TaintEffect string

const (
	// Do not allow new pods to schedule onto the node unless they tolerate the taint,
	// but allow all pods submitted to Kubelet without going through the scheduler
	// to start, and allow all already-running pods to continue running.
	// Enforced by the scheduler.
	// 不允许新的pod调度到节点上，除非它们容忍污点，但允许所有提交给Kubelet而不经过调度程序的pod启动，并允许所有已运行的pod继续运行。 由调度程序执行。
	TaintEffectNoSchedule TaintEffect = "NoSchedule"
	// Like TaintEffectNoSchedule, but the scheduler tries not to schedule
	// new pods onto the node, rather than prohibiting new pods from scheduling
	// onto the node entirely. Enforced by the scheduler.
	// 与TaintEffectNoSchedule类似，但调度程序尝试不要将新的pod调度到节点上，而不是完全禁止新的pod调度到节点上。 由调度程序执行。
	TaintEffectPreferNoSchedule TaintEffect = "PreferNoSchedule"
	// NOT YET IMPLEMENTED. TODO: Uncomment field once it is implemented.
	// Like TaintEffectNoSchedule, but additionally do not allow pods submitted to
	// Kubelet without going through the scheduler to start.
	// Enforced by Kubelet and the scheduler.
	// TaintEffectNoScheduleNoAdmit TaintEffect = "NoScheduleNoAdmit"

	// Evict any already-running pods that do not tolerate the taint.
	// Currently enforced by NodeController.
	// 逐出任何不容忍污点的已运行的pod。 目前由NodeController执行。
	TaintEffectNoExecute TaintEffect = "NoExecute"
)

// The pod this Toleration is attached to tolerates any taint that matches
// the triple <key,value,effect> using the matching operator <operator>.
// 此差附加到的pod允许使用匹配操作符<操作符>匹配三重<键，值，效果>的任何污染。
type Toleration struct {
	// Key is the taint key that the toleration applies to. Empty means match all taint keys.
	// If the key is empty, operator must be Exists; this combination means to match all values and all keys.
	// +optional
	// Key是污点键，容忍度适用于。 空意味着匹配所有污点键。 如果键为空，则操作符必须存在； 这种组合意味着匹配所有值和所有键。
	Key string `json:"key,omitempty" protobuf:"bytes,1,opt,name=key"`
	// Operator represents a key's relationship to the value.
	// Valid operators are Exists and Equal. Defaults to Equal.
	// Exists is equivalent to wildcard for value, so that a pod can
	// tolerate all taints of a particular category.
	// +optional
	// Operator表示键与值之间的关系。 有效的操作符是Exists和Equal。 默认为Equal。 Exists等同于值的通配符，因此pod可以容忍特定类别的所有污点。
	Operator TolerationOperator `json:"operator,omitempty" protobuf:"bytes,2,opt,name=operator,casttype=TolerationOperator"`
	// Value is the taint value the toleration matches to.
	// If the operator is Exists, the value should be empty, otherwise just a regular string.
	// +optional
	// Value是容忍度匹配的污点值。 如果操作符是Exists，则值应为空，否则只是一个常规字符串。
	Value string `json:"value,omitempty" protobuf:"bytes,3,opt,name=value"`
	// Effect indicates the taint effect to match. Empty means match all taint effects.
	// When specified, allowed values are NoSchedule, PreferNoSchedule and NoExecute.
	// +optional
	// Effect指示要匹配的污点效果。 空意味着匹配所有污点效果。 指定时，允许的值是NoSchedule，PreferNoSchedule和NoExecute。
	Effect TaintEffect `json:"effect,omitempty" protobuf:"bytes,4,opt,name=effect,casttype=TaintEffect"`
	// TolerationSeconds represents the period of time the toleration (which must be
	// of effect NoExecute, otherwise this field is ignored) tolerates the taint. By default,
	// it is not set, which means tolerate the taint forever (do not evict). Zero and
	// negative values will be treated as 0 (evict immediately) by the system.
	// +optional
	// TolerationSeconds表示容忍度（必须是NoExecute的效果，否则将忽略此字段）容忍污点的时间段。 默认情况下，它未设置，这意味着永远容忍污点（不要驱逐）。 系统将零和负值视为0（立即驱逐）。
	TolerationSeconds *int64 `json:"tolerationSeconds,omitempty" protobuf:"varint,5,opt,name=tolerationSeconds"`
}

// A toleration operator is the set of operators that can be used in a toleration.
// +enum
// 一个容忍度操作符是可以在容忍度中使用的操作符集。
type TolerationOperator string

const (
	TolerationOpExists TolerationOperator = "Exists"
	TolerationOpEqual  TolerationOperator = "Equal"
)

// PodReadinessGate contains the reference to a pod condition
// PodReadinessGate包含对pod条件的引用
type PodReadinessGate struct {
	// ConditionType refers to a condition in the pod's condition list with matching type.
	// ConditionType引用与匹配类型的pod条件列表中的条件。
	ConditionType PodConditionType `json:"conditionType" protobuf:"bytes,1,opt,name=conditionType,casttype=PodConditionType"`
}

// PodSpec is a description of a pod.
// PodSpec 是 Pod 的描述
type PodSpec struct {
	// List of volumes that can be mounted by containers belonging to the pod.
	// More info: https://kubernetes.io/docs/concepts/storage/volumes
	// +optional
	// +patchMergeKey=name
	// +patchStrategy=merge,retainKeys
	// 可以由属于pod的容器挂载的卷列表。
	Volumes []Volume `json:"volumes,omitempty" patchStrategy:"merge,retainKeys" patchMergeKey:"name" protobuf:"bytes,1,rep,name=volumes"`
	// List of initialization containers belonging to the pod.
	// Init containers are executed in order prior to containers being started. If any
	// init container fails, the pod is considered to have failed and is handled according
	// to its restartPolicy. The name for an init container or normal container must be
	// unique among all containers.
	// Init containers may not have Lifecycle actions, Readiness probes, Liveness probes, or Startup probes.
	// The resourceRequirements of an init container are taken into account during scheduling
	// by finding the highest request/limit for each resource type, and then using the max of
	// of that value or the sum of the normal containers. Limits are applied to init containers
	// in a similar fashion.
	// Init containers cannot currently be added or removed.
	// Cannot be updated.
	// More info: https://kubernetes.io/docs/concepts/workloads/pods/init-containers/
	// +patchMergeKey=name
	// +patchStrategy=merge
	// 属于pod的初始化容器列表。
	// 初始化容器按顺序执行，然后启动容器。 如果任何初始化容器失败，则pod被认为已失败，并根据其restartPolicy进行处理。 初始化容器或正常容器的名称必须在所有容器之间唯一。
	// 初始化容器可能没有生命周期操作，就绪探针，存活探针或启动探针。 在调度期间，通过查找每种资源类型的最高请求/限制，然后使用该值的最大值或正常容器的总和来考虑init容器的resourceRequirements。 限制以类似的方式应用于init容器。
	// 初始化容器目前无法添加或删除。
	InitContainers []Container `json:"initContainers,omitempty" patchStrategy:"merge" patchMergeKey:"name" protobuf:"bytes,20,rep,name=initContainers"`
	// List of containers belonging to the pod.
	// Containers cannot currently be added or removed.
	// There must be at least one container in a Pod.
	// Cannot be updated.
	// +patchMergeKey=name
	// +patchStrategy=merge
	// 属于pod的容器列表。
	// 目前无法添加或删除容器。 Pod中必须至少有一个容器。 不能更新。
	Containers []Container `json:"containers" patchStrategy:"merge" patchMergeKey:"name" protobuf:"bytes,2,rep,name=containers"`
	// List of ephemeral containers run in this pod. Ephemeral containers may be run in an existing
	// pod to perform user-initiated actions such as debugging. This list cannot be specified when
	// creating a pod, and it cannot be modified by updating the pod spec. In order to add an
	// ephemeral container to an existing pod, use the pod's ephemeralcontainers subresource.
	// +optional
	// +patchMergeKey=name
	// +patchStrategy=merge
	// 在此pod中运行的临时容器列表。 临时容器可以在现有pod中运行，以执行用户发起的操作，例如调试。 创建pod时无法指定此列表，并且无法通过更新pod spec来修改它。 为了将临时容器添加到现有pod中，请使用pod的ephemeralcontainers子资源。
	EphemeralContainers []EphemeralContainer `json:"ephemeralContainers,omitempty" patchStrategy:"merge" patchMergeKey:"name" protobuf:"bytes,34,rep,name=ephemeralContainers"`
	// Restart policy for all containers within the pod.
	// One of Always, OnFailure, Never.
	// Default to Always.
	// More info: https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle/#restart-policy
	// +optional
	// 重新启动pod中所有容器的策略。永远，失败，永远。
	RestartPolicy RestartPolicy `json:"restartPolicy,omitempty" protobuf:"bytes,3,opt,name=restartPolicy,casttype=RestartPolicy"`
	// Optional duration in seconds the pod needs to terminate gracefully. May be decreased in delete request.
	// Value must be non-negative integer. The value zero indicates stop immediately via
	// the kill signal (no opportunity to shut down).
	// If this value is nil, the default grace period will be used instead.
	// The grace period is the duration in seconds after the processes running in the pod are sent
	// a termination signal and the time when the processes are forcibly halted with a kill signal.
	// Set this value longer than the expected cleanup time for your process.
	// Defaults to 30 seconds.
	// +optional
	// 可选的秒数，pod需要优雅地终止。 可以在删除请求中减少。
	// 值必须是非负整数。 值为零表示通过立即停止信号停止（没有机会关闭）。
	// 如果此值为nil，则将使用默认的宽限期。
	// 宽限期是在pod中运行的进程发送终止信号后的持续时间，以及强制使用kill信号终止进程的时间。
	// 将此值设置为大于预期清理时间的值。
	// 默认为30秒。
	TerminationGracePeriodSeconds *int64 `json:"terminationGracePeriodSeconds,omitempty" protobuf:"varint,4,opt,name=terminationGracePeriodSeconds"`
	// Optional duration in seconds the pod may be active on the node relative to
	// StartTime before the system will actively try to mark it failed and kill associated containers.
	// Value must be a positive integer.
	// +optional
	// 可选的秒数，pod可以在节点上活动的时间相对于开始时间，系统将主动尝试将其标记为失败并杀死相关容器。
	ActiveDeadlineSeconds *int64 `json:"activeDeadlineSeconds,omitempty" protobuf:"varint,5,opt,name=activeDeadlineSeconds"`
	// Set DNS policy for the pod.
	// Defaults to "ClusterFirst".
	// Valid values are 'ClusterFirstWithHostNet', 'ClusterFirst', 'Default' or 'None'.
	// DNS parameters given in DNSConfig will be merged with the policy selected with DNSPolicy.
	// To have DNS options set along with hostNetwork, you have to specify DNS policy
	// explicitly to 'ClusterFirstWithHostNet'.
	// +optional
	// 为pod设置DNS策略。默认为“ClusterFirst”。
	DNSPolicy DNSPolicy `json:"dnsPolicy,omitempty" protobuf:"bytes,6,opt,name=dnsPolicy,casttype=DNSPolicy"`
	// NodeSelector is a selector which must be true for the pod to fit on a node.
	// Selector which must match a node's labels for the pod to be scheduled on that node.
	// More info: https://kubernetes.io/docs/concepts/configuration/assign-pod-node/
	// +optional
	// +mapType=atomic
	// NodeSelector 是一个选择器，必须为 pod 适合节点。
	NodeSelector map[string]string `json:"nodeSelector,omitempty" protobuf:"bytes,7,rep,name=nodeSelector"`

	// ServiceAccountName is the name of the ServiceAccount to use to run this pod.
	// More info: https://kubernetes.io/docs/tasks/configure-pod-container/configure-service-account/
	// +optional
	// ServiceAccountName 是用于运行此 pod 的 ServiceAccount 的名称。
	ServiceAccountName string `json:"serviceAccountName,omitempty" protobuf:"bytes,8,opt,name=serviceAccountName"`
	// DeprecatedServiceAccount is a depreciated alias for ServiceAccountName.
	// Deprecated: Use serviceAccountName instead.
	// +k8s:conversion-gen=false
	// +optional
	// DeprecatedServiceAccount 是 ServiceAccountName 的别名。
	DeprecatedServiceAccount string `json:"serviceAccount,omitempty" protobuf:"bytes,9,opt,name=serviceAccount"`
	// AutomountServiceAccountToken indicates whether a service account token should be automatically mounted.
	// +optional
	// AutomountServiceAccountToken 表示是否应自动挂载服务帐户令牌。
	AutomountServiceAccountToken *bool `json:"automountServiceAccountToken,omitempty" protobuf:"varint,21,opt,name=automountServiceAccountToken"`

	// NodeName is a request to schedule this pod onto a specific node. If it is non-empty,
	// the scheduler simply schedules this pod onto that node, assuming that it fits resource
	// requirements.
	// +optional
	// NodeName 是将此 pod 调度到特定节点的请求。 如果它不为空，则调度程序只需将此 pod 调度到该节点，假设它符合资源要求。
	NodeName string `json:"nodeName,omitempty" protobuf:"bytes,10,opt,name=nodeName"`
	// Host networking requested for this pod. Use the host's network namespace.
	// If this option is set, the ports that will be used must be specified.
	// Default to false.
	// +k8s:conversion-gen=false
	// +optional
	// HostNetwork 需要此 pod 的主机网络。 使用主机的网络名称空间。
	HostNetwork bool `json:"hostNetwork,omitempty" protobuf:"varint,11,opt,name=hostNetwork"`
	// Use the host's pid namespace.
	// Optional: Default to false.
	// +k8s:conversion-gen=false
	// +optional
	// 使用主机的 pid 名称空间。
	HostPID bool `json:"hostPID,omitempty" protobuf:"varint,12,opt,name=hostPID"`
	// Use the host's ipc namespace.
	// Optional: Default to false.
	// +k8s:conversion-gen=false
	// +optional
	// 使用主机的 ipc 名称空间。
	HostIPC bool `json:"hostIPC,omitempty" protobuf:"varint,13,opt,name=hostIPC"`
	// Share a single process namespace between all of the containers in a pod.
	// When this is set containers will be able to view and signal processes from other containers
	// in the same pod, and the first process in each container will not be assigned PID 1.
	// HostPID and ShareProcessNamespace cannot both be set.
	// Optional: Default to false.
	// +k8s:conversion-gen=false
	// +optional
	// 在 pod 中的所有容器之间共享单个进程名称空间。
	// 当设置时，容器将能够查看和信号来自同一 pod 中的其他容器的进程，并且每个容器中的第一个进程将不会被分配 PID 1。
	// HostPID 和 ShareProcessNamespace 不能同时设置。
	ShareProcessNamespace *bool `json:"shareProcessNamespace,omitempty" protobuf:"varint,27,opt,name=shareProcessNamespace"`
	// SecurityContext holds pod-level security attributes and common container settings.
	// Optional: Defaults to empty.  See type description for default values of each field.
	// +optional
	// SecurityContext 保存 pod 级别的安全属性和常用容器设置。
	SecurityContext *PodSecurityContext `json:"securityContext,omitempty" protobuf:"bytes,14,opt,name=securityContext"`
	// ImagePullSecrets is an optional list of references to secrets in the same namespace to use for pulling any of the images used by this PodSpec.
	// If specified, these secrets will be passed to individual puller implementations for them to use.
	// More info: https://kubernetes.io/docs/concepts/containers/images#specifying-imagepullsecrets-on-a-pod
	// +optional
	// +patchMergeKey=name
	// +patchStrategy=merge
	// ImagePullSecrets是一个可选的引用列表，它们引用同一名称空间中的密钥，用于拉取此 PodSpec 使用的任何图像。
	ImagePullSecrets []LocalObjectReference `json:"imagePullSecrets,omitempty" patchStrategy:"merge" patchMergeKey:"name" protobuf:"bytes,15,rep,name=imagePullSecrets"`
	// Specifies the hostname of the Pod
	// If not specified, the pod's hostname will be set to a system-defined value.
	// +optional
	// 指定 Pod 的主机名
	Hostname string `json:"hostname,omitempty" protobuf:"bytes,16,opt,name=hostname"`
	// If specified, the fully qualified Pod hostname will be "<hostname>.<subdomain>.<pod namespace>.svc.<cluster domain>".
	// If not specified, the pod will not have a domainname at all.
	// +optional
	// 如果指定，完全限定的 Pod 主机名将为 "<hostname>.<subdomain>.<pod namespace>.svc.<cluster domain>"。
	Subdomain string `json:"subdomain,omitempty" protobuf:"bytes,17,opt,name=subdomain"`
	// If specified, the pod's scheduling constraints
	// +optional
	// 如果指定，pod 的调度约束
	Affinity *Affinity `json:"affinity,omitempty" protobuf:"bytes,18,opt,name=affinity"`
	// If specified, the pod will be dispatched by specified scheduler.
	// If not specified, the pod will be dispatched by default scheduler.
	// +optional
	// 如果指定，pod 将由指定的调度程序调度。 如果未指定，pod 将由默认调度程序调度。
	SchedulerName string `json:"schedulerName,omitempty" protobuf:"bytes,19,opt,name=schedulerName"`
	// If specified, the pod's tolerations.
	// +optional
	// 如果指定，pod 的容忍度。
	Tolerations []Toleration `json:"tolerations,omitempty" protobuf:"bytes,22,opt,name=tolerations"`
	// HostAliases is an optional list of hosts and IPs that will be injected into the pod's hosts
	// file if specified. This is only valid for non-hostNetwork pods.
	// +optional
	// +patchMergeKey=ip
	// +patchStrategy=merge
	// HostAliases 是一个可选的主机和 IP 列表，如果指定，将被注入到 pod 的主机文件中。
	HostAliases []HostAlias `json:"hostAliases,omitempty" patchStrategy:"merge" patchMergeKey:"ip" protobuf:"bytes,23,rep,name=hostAliases"`
	// If specified, indicates the pod's priority. "system-node-critical" and
	// "system-cluster-critical" are two special keywords which indicate the
	// highest priorities with the former being the highest priority. Any other
	// name must be defined by creating a PriorityClass object with that name.
	// If not specified, the pod priority will be default or zero if there is no
	// default.
	// +optional
	// 如果指定，表示 pod 的优先级。 "system-node-critical" 和 "system-cluster-critical" 是两个特殊的关键字，它们表示最高优先级，前者具有最高优先级。
	// 任何其他名称必须通过创建名称为该名称的 PriorityClass 对象来定义。
	// 如果未指定，pod 优先级将为默认值或零（如果没有默认值）。
	PriorityClassName string `json:"priorityClassName,omitempty" protobuf:"bytes,24,opt,name=priorityClassName"`
	// The priority value. Various system components use this field to find the
	// priority of the pod. When Priority Admission Controller is enabled, it
	// prevents users from setting this field. The admission controller populates
	// this field from PriorityClassName.
	// The higher the value, the higher the priority.
	// +optional
	// 优先级值。 各种系统组件使用此字段来查找 pod 的优先级。 启用 Priority Admission Controller 时，它会阻止用户设置此字段。
	// 准入控制器从 PriorityClassName 填充此字段。
	Priority *int32 `json:"priority,omitempty" protobuf:"bytes,25,opt,name=priority"`
	// Specifies the DNS parameters of a pod.
	// Parameters specified here will be merged to the generated DNS
	// configuration based on DNSPolicy.
	// +optional
	// 指定 pod 的 DNS 参数。
	// 此处指定的参数将合并到基于 DNSPolicy 生成的 DNS 配置中。
	DNSConfig *PodDNSConfig `json:"dnsConfig,omitempty" protobuf:"bytes,26,opt,name=dnsConfig"`
	// If specified, all readiness gates will be evaluated for pod readiness.
	// A pod is ready when all its containers are ready AND
	// all conditions specified in the readiness gates have status equal to "True"
	// More info: https://git.k8s.io/enhancements/keps/sig-network/580-pod-readiness-gates
	// +optional
	// 如果指定，将评估所有准备就绪门闩以检查 pod 的准备情况。
	// 当 pod 中的所有容器都准备就绪并且准备就绪门闩中指定的所有条件的状态等于 "True" 时，pod 就绪。
	ReadinessGates []PodReadinessGate `json:"readinessGates,omitempty" protobuf:"bytes,28,opt,name=readinessGates"`
	// RuntimeClassName refers to a RuntimeClass object in the node.k8s.io group, which should be used
	// to run this pod.  If no RuntimeClass resource matches the named class, the pod will not be run.
	// If unset or empty, the "legacy" RuntimeClass will be used, which is an implicit class with an
	// empty definition that uses the default runtime handler.
	// More info: https://git.k8s.io/enhancements/keps/sig-node/585-runtime-class
	// +optional
	// RuntimeClassName 引用 node.k8s.io 组中的 RuntimeClass 对象，该对象应该用于运行此 pod。
	// 如果没有与命名类匹配的 RuntimeClass 资源，则不会运行 pod。
	// 如果未设置或为空，则将使用 "legacy" RuntimeClass，这是一个带有空定义的隐式类，它使用默认的运行时处理程序。
	RuntimeClassName *string `json:"runtimeClassName,omitempty" protobuf:"bytes,29,opt,name=runtimeClassName"`
	// EnableServiceLinks indicates whether information about services should be injected into pod's
	// environment variables, matching the syntax of Docker links.
	// Optional: Defaults to true.
	// +optional
	// EnableServiceLinks 指示是否应将有关服务的信息注入到 pod 的环境变量中，以匹配 Docker 链接的语法。
	EnableServiceLinks *bool `json:"enableServiceLinks,omitempty" protobuf:"varint,30,opt,name=enableServiceLinks"`
	// PreemptionPolicy is the Policy for preempting pods with lower priority.
	// One of Never, PreemptLowerPriority.
	// Defaults to PreemptLowerPriority if unset.
	// +optional
	// PreemptionPolicy 是用于抢占具有较低优先级的 pod 的策略。
	// 一个是 Never，一个是 PreemptLowerPriority。
	PreemptionPolicy *PreemptionPolicy `json:"preemptionPolicy,omitempty" protobuf:"bytes,31,opt,name=preemptionPolicy"`
	// Overhead represents the resource overhead associated with running a pod for a given RuntimeClass.
	// This field will be autopopulated at admission time by the RuntimeClass admission controller. If
	// the RuntimeClass admission controller is enabled, overhead must not be set in Pod create requests.
	// The RuntimeClass admission controller will reject Pod create requests which have the overhead already
	// set. If RuntimeClass is configured and selected in the PodSpec, Overhead will be set to the value
	// defined in the corresponding RuntimeClass, otherwise it will remain unset and treated as zero.
	// More info: https://git.k8s.io/enhancements/keps/sig-node/688-pod-overhead/README.md
	// +optional
	// Overhead 表示运行给定 RuntimeClass 的 pod 所需的资源开销。 此字段将在准入时由 RuntimeClass 准入控制器自动填充。
	// 如果启用了 RuntimeClass 准入控制器，则不得在 Pod 创建请求中设置开销。
	// RuntimeClass 准入控制器将拒绝已经设置了开销的 Pod 创建请求。
	// 如果配置了 RuntimeClass 并在 PodSpec 中选择了它，则 Overhead 将设置为相应 RuntimeClass 中定义的值，否则它将保持未设置并视为零。
	Overhead ResourceList `json:"overhead,omitempty" protobuf:"bytes,32,opt,name=overhead"`
	// TopologySpreadConstraints describes how a group of pods ought to spread across topology
	// domains. Scheduler will schedule pods in a way which abides by the constraints.
	// All topologySpreadConstraints are ANDed.
	// +optional
	// +patchMergeKey=topologyKey
	// +patchStrategy=merge
	// +listType=map
	// +listMapKey=topologyKey
	// +listMapKey=whenUnsatisfiable
	// TopologySpreadConstraints 描述了一组 pod 应该如何在拓扑域之间分布。 调度程序将以符合约束的方式调度 pod。
	TopologySpreadConstraints []TopologySpreadConstraint `json:"topologySpreadConstraints,omitempty" patchStrategy:"merge" patchMergeKey:"topologyKey" protobuf:"bytes,33,opt,name=topologySpreadConstraints"`
	// If true the pod's hostname will be configured as the pod's FQDN, rather than the leaf name (the default).
	// In Linux containers, this means setting the FQDN in the hostname field of the kernel (the nodename field of struct utsname).
	// In Windows containers, this means setting the registry value of hostname for the registry key HKEY_LOCAL_MACHINE\\SYSTEM\\CurrentControlSet\\Services\\Tcpip\\Parameters to FQDN.
	// If a pod does not have FQDN, this has no effect.
	// Default to false.
	// +optional
	// 如果为 true，则 pod 的主机名将配置为 pod 的 FQDN，而不是叶子名称（默认值）。
	// 在 Linux 容器中，这意味着在内核的主机名字段（结构体 utsname 的 nodename 字段）中设置 FQDN。
	// 在 Windows 容器中，这意味着将主机名的注册表值设置为 HKEY_LOCAL_MACHINE\\SYSTEM\\CurrentControlSet\\Services\\Tcpip\\Parameters 的注册表键的 FQDN。
	// 如果 pod 没有 FQDN，则不起作用。
	SetHostnameAsFQDN *bool `json:"setHostnameAsFQDN,omitempty" protobuf:"varint,35,opt,name=setHostnameAsFQDN"`
	// Specifies the OS of the containers in the pod.
	// Some pod and container fields are restricted if this is set.
	//
	// If the OS field is set to linux, the following fields must be unset:
	// -securityContext.windowsOptions
	//
	// If the OS field is set to windows, following fields must be unset:
	// - spec.hostPID
	// - spec.hostIPC
	// - spec.hostUsers
	// - spec.securityContext.seLinuxOptions
	// - spec.securityContext.seccompProfile
	// - spec.securityContext.fsGroup
	// - spec.securityContext.fsGroupChangePolicy
	// - spec.securityContext.sysctls
	// - spec.shareProcessNamespace
	// - spec.securityContext.runAsUser
	// - spec.securityContext.runAsGroup
	// - spec.securityContext.supplementalGroups
	// - spec.containers[*].securityContext.seLinuxOptions
	// - spec.containers[*].securityContext.seccompProfile
	// - spec.containers[*].securityContext.capabilities
	// - spec.containers[*].securityContext.readOnlyRootFilesystem
	// - spec.containers[*].securityContext.privileged
	// - spec.containers[*].securityContext.allowPrivilegeEscalation
	// - spec.containers[*].securityContext.procMount
	// - spec.containers[*].securityContext.runAsUser
	// - spec.containers[*].securityContext.runAsGroup
	// +optional
	// 指定 pod 中容器的操作系统。 如果设置了此字段，则某些 pod 和容器字段将受到限制。
	// 如果将 OS 字段设置为 linux，则必须取消以下字段设置：
	OS *PodOS `json:"os,omitempty" protobuf:"bytes,36,opt,name=os"`

	// Use the host's user namespace.
	// Optional: Default to true.
	// If set to true or not present, the pod will be run in the host user namespace, useful
	// for when the pod needs a feature only available to the host user namespace, such as
	// loading a kernel module with CAP_SYS_MODULE.
	// When set to false, a new userns is created for the pod. Setting false is useful for
	// mitigating container breakout vulnerabilities even allowing users to run their
	// containers as root without actually having root privileges on the host.
	// This field is alpha-level and is only honored by servers that enable the UserNamespacesSupport feature.
	// +k8s:conversion-gen=false
	// +optional
	// 使用主机的用户命名空间。
	// 可选：默认为 true。
	// 如果设置为 true 或未设置，则 pod 将在主机用户命名空间中运行，这对于 pod 需要仅在主机用户命名空间中可用的功能时很有用，例如使用 CAP_SYS_MODULE 加载内核模块。
	// 当设置为 false 时，将为 pod 创建一个新的 userns。 将其设置为 false 对于缓解容器突破漏洞尤其是允许用户以 root 身份运行其容器，即使主机上没有 root 权限也是如此。
	// 此字段为 alpha 级别，并且仅由启用 UserNamespacesSupport 功能的服务器支持。
	HostUsers *bool `json:"hostUsers,omitempty" protobuf:"bytes,37,opt,name=hostUsers"`

	// SchedulingGates is an opaque list of values that if specified will block scheduling the pod.
	// More info:  https://git.k8s.io/enhancements/keps/sig-scheduling/3521-pod-scheduling-readiness.
	//
	// This is an alpha-level feature enabled by PodSchedulingReadiness feature gate.
	// +optional
	// +patchMergeKey=name
	// +patchStrategy=merge
	// +listType=map
	// +listMapKey=name
	// SchedulingGates是一个不透明的值列表，如果指定将阻止调度 pod。
	SchedulingGates []PodSchedulingGate `json:"schedulingGates,omitempty" patchStrategy:"merge" patchMergeKey:"name" protobuf:"bytes,38,opt,name=schedulingGates"`
	// ResourceClaims defines which ResourceClaims must be allocated
	// and reserved before the Pod is allowed to start. The resources
	// will be made available to those containers which consume them
	// by name.
	//
	// This is an alpha field and requires enabling the
	// DynamicResourceAllocation feature gate.
	//
	// This field is immutable.
	//
	// +patchMergeKey=name
	// +patchStrategy=merge,retainKeys
	// +listType=map
	// +listMapKey=name
	// +featureGate=DynamicResourceAllocation
	// +optional
	// ResourceClaims 定义了在允许 Pod 启动之前必须分配和保留哪些 ResourceClaims。 资源将被这些容器使用名称消耗。
	ResourceClaims []PodResourceClaim `json:"resourceClaims,omitempty" patchStrategy:"merge,retainKeys" patchMergeKey:"name" protobuf:"bytes,39,rep,name=resourceClaims"`
}

// PodResourceClaim references exactly one ResourceClaim through a ClaimSource.
// It adds a name to it that uniquely identifies the ResourceClaim inside the Pod.
// Containers that need access to the ResourceClaim reference it with this name.
// PodResourceClaim 通过一个ClaimSource引用一个resourcecclaim。
// 它添加了一个名称，该名称在Pod内唯一标识resourceclaim。
// 需要访问resourceclaim的容器通过此名称引用它。
type PodResourceClaim struct {
	// Name uniquely identifies this resource claim inside the pod.
	// This must be a DNS_LABEL.
	// 名字在pod内唯一标识这个resourceclaim。
	Name string `json:"name" protobuf:"bytes,1,name=name"`

	// Source describes where to find the ResourceClaim.
	// 描述在哪里找到resourceclaim。
	Source ClaimSource `json:"source,omitempty" protobuf:"bytes,2,name=source"`
}

// ClaimSource describes a reference to a ResourceClaim.
//
// Exactly one of these fields should be set.  Consumers of this type must
// treat an empty object as if it has an unknown value.
// ClaimSource 描述了对resourceclaim的引用。
// 这些字段中的一个应该被设置。 此类型的消费者必须将空对象视为未知值。
type ClaimSource struct {
	// ResourceClaimName is the name of a ResourceClaim object in the same
	// namespace as this pod.
	// ResourceClaimName 是与此 pod 相同命名空间中的 ResourceClaim 对象的名称。
	ResourceClaimName *string `json:"resourceClaimName,omitempty" protobuf:"bytes,1,opt,name=resourceClaimName"`

	// ResourceClaimTemplateName is the name of a ResourceClaimTemplate
	// object in the same namespace as this pod.
	//
	// The template will be used to create a new ResourceClaim, which will
	// be bound to this pod. When this pod is deleted, the ResourceClaim
	// will also be deleted. The name of the ResourceClaim will be <pod
	// name>-<resource name>, where <resource name> is the
	// PodResourceClaim.Name. Pod validation will reject the pod if the
	// concatenated name is not valid for a ResourceClaim (e.g. too long).
	//
	// An existing ResourceClaim with that name that is not owned by the
	// pod will not be used for the pod to avoid using an unrelated
	// resource by mistake. Scheduling and pod startup are then blocked
	// until the unrelated ResourceClaim is removed.
	//
	// This field is immutable and no changes will be made to the
	// corresponding ResourceClaim by the control plane after creating the
	// ResourceClaim.
	// ResourceClaimTemplateName 是与此 pod 相同命名空间中的 ResourceClaimTemplate 对象的名称。
	// 模板将用于创建新的 ResourceClaim，该 ResourceClaim 将绑定到此 pod。 当删除此 pod 时，ResourceClaim 也将被删除。 ResourceClaim 的名称将为 <pod name>-<resource name>，其中 <resource name> 是 PodResourceClaim.Name。 Pod 验证将拒绝 pod，如果连接的名称对于 ResourceClaim 无效（例如太长）。
	// 具有该名称但不属于 pod 的现有 ResourceClaim 将不会用于 pod，以避免意外使用不相关的资源。 然后，调度和 pod 启动将被阻止，直到删除不相关的 ResourceClaim。
	// 此字段不可变，并且在创建 ResourceClaim 后，控制平面不会对相应的 ResourceClaim 进行任何更改。
	ResourceClaimTemplateName *string `json:"resourceClaimTemplateName,omitempty" protobuf:"bytes,2,opt,name=resourceClaimTemplateName"`
}

// OSName is the set of OS'es that can be used in OS.
// OSName 是可以在 OS 中使用的 OS 集合。
type OSName string

// These are valid values for OSName
const (
	Linux   OSName = "linux"
	Windows OSName = "windows"
)

// PodOS defines the OS parameters of a pod.
// PodOS 定义了 pod 的 OS 参数。
type PodOS struct {
	// Name is the name of the operating system. The currently supported values are linux and windows.
	// Additional value may be defined in future and can be one of:
	// https://github.com/opencontainers/runtime-spec/blob/master/config.md#platform-specific-configuration
	// Clients should expect to handle additional values and treat unrecognized values in this field as os: null
	// Name 是操作系统的名称。 目前支持的值是 linux 和 windows。 未来可能会定义其他值，可以是以下之一：
	Name OSName `json:"name" protobuf:"bytes,1,opt,name=name"`
}

// PodSchedulingGate is associated to a Pod to guard its scheduling.
// PodSchedulingGate 与 Pod 关联，以保护其调度。
type PodSchedulingGate struct {
	// Name of the scheduling gate.
	// Each scheduling gate must have a unique name field.
	// 调度门的名称 , 每个调度门必须有一个唯一的名称字段。
	Name string `json:"name" protobuf:"bytes,1,opt,name=name"`
}

// +enum
type UnsatisfiableConstraintAction string

const (
	// DoNotSchedule instructs the scheduler not to schedule the pod
	// when constraints are not satisfied.
	// DoNotSchedule 指示调度程序在不满足约束时不要调度 pod。
	DoNotSchedule UnsatisfiableConstraintAction = "DoNotSchedule"
	// ScheduleAnyway instructs the scheduler to schedule the pod
	// even if constraints are not satisfied.
	// ScheduleAnyway 指示调度程序即使不满足约束也要调度 pod。
	ScheduleAnyway UnsatisfiableConstraintAction = "ScheduleAnyway"
)

// NodeInclusionPolicy defines the type of node inclusion policy
// +enum
// NodeInclusionPolicy 定义了节点包含策略的类型
type NodeInclusionPolicy string

const (
	// NodeInclusionPolicyIgnore means ignore this scheduling directive when calculating pod topology spread skew.
	// NodeInclusionPolicyIgnore 意味着在计算 pod 拓扑扩展偏斜时忽略此调度指令。
	NodeInclusionPolicyIgnore NodeInclusionPolicy = "Ignore"
	// NodeInclusionPolicyHonor means use this scheduling directive when calculating pod topology spread skew.
	// NodeInclusionPolicyHonor 意味着在计算 pod 拓扑扩展偏斜时使用此调度指令。
	NodeInclusionPolicyHonor NodeInclusionPolicy = "Honor"
)

// TopologySpreadConstraint specifies how to spread matching pods among the given topology.
// TopologySpreadConstraint 指定如何在给定拓扑之间分散匹配的 pod。
type TopologySpreadConstraint struct {
	// MaxSkew describes the degree to which pods may be unevenly distributed.
	// When `whenUnsatisfiable=DoNotSchedule`, it is the maximum permitted difference
	// between the number of matching pods in the target topology and the global minimum.
	// The global minimum is the minimum number of matching pods in an eligible domain
	// or zero if the number of eligible domains is less than MinDomains.
	// For example, in a 3-zone cluster, MaxSkew is set to 1, and pods with the same
	// labelSelector spread as 2/2/1:
	// In this case, the global minimum is 1.
	// +-------+-------+-------+
	// | zone1 | zone2 | zone3 |
	// +-------+-------+-------+
	// |  P P  |  P P  |   P   |
	// +-------+-------+-------+
	// - if MaxSkew is 1, incoming pod can only be scheduled to zone3 to become 2/2/2;
	// scheduling it onto zone1(zone2) would make the ActualSkew(3-1) on zone1(zone2)
	// violate MaxSkew(1).
	// - if MaxSkew is 2, incoming pod can be scheduled onto any zone.
	// When `whenUnsatisfiable=ScheduleAnyway`, it is used to give higher precedence
	// to topologies that satisfy it.
	// It's a required field. Default value is 1 and 0 is not allowed.
	// MaxSkew 描述了 pod 分布不均匀的程度。 当 `whenUnsatisfiable=DoNotSchedule` 时，它是目标拓扑中匹配 pod 的数量与全局最小值之间允许的最大差异。 全局最小值是合格域中匹配 pod 的最小数量，如果合格域的数量少于 MinDomains，则为零。 例如，在 3 个区域的集群中，MaxSkew 设置为 1，并且具有相同 labelSelector 的 pod 分布为 2/2/1： 在这种情况下，全局最小值为 1。 +-------+-------+-------+ | zone1 | zone2 | zone3 | +-------+-------+-------+ | P P | P P | P | +-------+-------+-------+ - 如果 MaxSkew 为 1，则传入的 pod 只能调度到 zone3 以变为 2/2/2; 将其调度到 zone1(zone2) 将使 zone1(zone2) 上的 ActualSkew(3-1) 违反 MaxSkew(1)。 - 如果 MaxSkew 为 2，则传入的 pod 可以调度到任何区域。 当 `whenUnsatisfiable=ScheduleAnyway` 时，它用于给满足它的拓扑更高的优先级。 这是一个必填字段。 默认值为 1，不允许为 0。
	MaxSkew int32 `json:"maxSkew" protobuf:"varint,1,opt,name=maxSkew"`
	// TopologyKey is the key of node labels. Nodes that have a label with this key
	// and identical values are considered to be in the same topology.
	// We consider each <key, value> as a "bucket", and try to put balanced number
	// of pods into each bucket.
	// We define a domain as a particular instance of a topology.
	// Also, we define an eligible domain as a domain whose nodes meet the requirements of
	// nodeAffinityPolicy and nodeTaintsPolicy.
	// e.g. If TopologyKey is "kubernetes.io/hostname", each Node is a domain of that topology.
	// And, if TopologyKey is "topology.kubernetes.io/zone", each zone is a domain of that topology.
	// It's a required field.
	// TopologyKey 是节点标签的键。 具有此键和相同值的标签的节点被认为是在同一拓扑中。 我们将每个 <key, value> 视为一个“桶”，并尝试将平衡数量的 pod 放入每个桶中。 我们将域定义为拓扑的特定实例。 此外，我们将合格域定义为其节点满足 nodeAffinityPolicy 和 nodeTaintsPolicy 的要求的域。 例如。 如果 TopologyKey 为 "kubernetes.io/hostname"，则每个节点都是该拓扑的域。 并且，如果 TopologyKey 为 "topology.kubernetes.io/zone"，则每个区域都是该拓扑的域。 这是一个必填字段。
	TopologyKey string `json:"topologyKey" protobuf:"bytes,2,opt,name=topologyKey"`
	// WhenUnsatisfiable indicates how to deal with a pod if it doesn't satisfy
	// the spread constraint.
	// - DoNotSchedule (default) tells the scheduler not to schedule it.
	// - ScheduleAnyway tells the scheduler to schedule the pod in any location,
	//   but giving higher precedence to topologies that would help reduce the
	//   skew.
	// A constraint is considered "Unsatisfiable" for an incoming pod
	// if and only if every possible node assignment for that pod would violate
	// "MaxSkew" on some topology.
	// For example, in a 3-zone cluster, MaxSkew is set to 1, and pods with the same
	// labelSelector spread as 3/1/1:
	// +-------+-------+-------+
	// | zone1 | zone2 | zone3 |
	// +-------+-------+-------+
	// | P P P |   P   |   P   |
	// +-------+-------+-------+
	// If WhenUnsatisfiable is set to DoNotSchedule, incoming pod can only be scheduled
	// to zone2(zone3) to become 3/2/1(3/1/2) as ActualSkew(2-1) on zone2(zone3) satisfies
	// MaxSkew(1). In other words, the cluster can still be imbalanced, but scheduler
	// won't make it *more* imbalanced.
	// It's a required field.
	// WhenUnsatisfiable 指示如果 pod 不满足分布约束，如何处理它。
	//- DoNotSchedule(默认) 告诉调度程序不要对其进行调度。
	//- ScheduleAnyway 告诉调度程序在任何位置调度 pod，但是给予更高优先级的拓扑，以帮助减少偏斜。
	//如果仅当该 pod 的每个可能的节点分配都会违反某些拓扑上的 "MaxSkew" 时，才将约束视为传入 pod 的 "Unsatisfiable"。 例如，在 3 个区域的集群中，MaxSkew 设置为 1，并且具有相同 labelSelector 的 pod 分布为 3/1/1： +-------+-------+-------+ | zone1 | zone2 | zone3 | +-------+-------+-------+ | P P P | P | P | +-------+-------+-------+ 如果 WhenUnsatisfiable 设置为 DoNotSchedule，则传入的 pod 只能调度到 zone2(zone3) 以变为 3/2/1(3/1/2)，因为 zone2(zone3) 上的 ActualSkew(2-1) 满足 MaxSkew(1)。 换句话说，集群仍然可能不平衡，但调度程序不会使其更不平衡。 这是一个必填字段。
	WhenUnsatisfiable UnsatisfiableConstraintAction `json:"whenUnsatisfiable" protobuf:"bytes,3,opt,name=whenUnsatisfiable,casttype=UnsatisfiableConstraintAction"`
	// LabelSelector is used to find matching pods.
	// Pods that match this label selector are counted to determine the number of pods
	// in their corresponding topology domain.
	// +optional
	// LabelSelector 用于查找匹配的 pod。 匹配此标签选择器的 pod 被计数以确定其相应拓扑域中的 pod 数量。 +可选
	LabelSelector *metav1.LabelSelector `json:"labelSelector,omitempty" protobuf:"bytes,4,opt,name=labelSelector"`
	// MinDomains indicates a minimum number of eligible domains.
	// When the number of eligible domains with matching topology keys is less than minDomains,
	// Pod Topology Spread treats "global minimum" as 0, and then the calculation of Skew is performed.
	// And when the number of eligible domains with matching topology keys equals or greater than minDomains,
	// this value has no effect on scheduling.
	// As a result, when the number of eligible domains is less than minDomains,
	// scheduler won't schedule more than maxSkew Pods to those domains.
	// If value is nil, the constraint behaves as if MinDomains is equal to 1.
	// Valid values are integers greater than 0.
	// When value is not nil, WhenUnsatisfiable must be DoNotSchedule.
	//
	// For example, in a 3-zone cluster, MaxSkew is set to 2, MinDomains is set to 5 and pods with the same
	// labelSelector spread as 2/2/2:
	// +-------+-------+-------+
	// | zone1 | zone2 | zone3 |
	// +-------+-------+-------+
	// |  P P  |  P P  |  P P  |
	// +-------+-------+-------+
	// The number of domains is less than 5(MinDomains), so "global minimum" is treated as 0.
	// In this situation, new pod with the same labelSelector cannot be scheduled,
	// because computed skew will be 3(3 - 0) if new Pod is scheduled to any of the three zones,
	// it will violate MaxSkew.
	//
	// This is a beta field and requires the MinDomainsInPodTopologySpread feature gate to be enabled (enabled by default).
	// +optional
	// MinDomains 表示最少可用域的数量。 当与匹配的拓扑键的可用域的数量少于 minDomains 时，Pod 拓扑扩展将 "全局最小值" 视为 0，然后执行 Skew 计算。 当与匹配的拓扑键的可用域的数量等于或大于 minDomains 时，此值对调度没有影响。 因此，当可用域的数量少于 minDomains 时，调度程序不会将超过 maxSkew Pod 调度到这些域。 如果值为 nil，则约束的行为与 MinDomains 等于 1 时相同。 有效值为大于 0 的整数。 当值不为 nil 时，WhenUnsatisfiable 必须为 DoNotSchedule。 例如，在 3 个区域的集群中，MaxSkew 设置为 2，MinDomains 设置为 5，并且具有相同 labelSelector 的 pod 分布为 2/2/2： +-------+-------+-------+ | zone1 | zone2 | zone3 | +-------+-------+-------+ | P P | P P | P P | +-------+-------+-------+ 可用域的数量少于 5(MinDomains)，因此将 "全局最小值" 视为 0。 在这种情况下，具有相同 labelSelector 的新 pod 无法被调度，因为如果将新 Pod 调度到这三个区域的任何一个，则计算出的偏斜将为 3(3-0)，这将违反 MaxSkew。 这是一个 beta 字段，需要启用 MinDomainsInPodTopologySpread 功能门控(默认启用)。 +可选
	MinDomains *int32 `json:"minDomains,omitempty" protobuf:"varint,5,opt,name=minDomains"`
	// NodeAffinityPolicy indicates how we will treat Pod's nodeAffinity/nodeSelector
	// when calculating pod topology spread skew. Options are:
	// - Honor: only nodes matching nodeAffinity/nodeSelector are included in the calculations.
	// - Ignore: nodeAffinity/nodeSelector are ignored. All nodes are included in the calculations.
	//
	// If this value is nil, the behavior is equivalent to the Honor policy.
	// This is a beta-level feature default enabled by the NodeInclusionPolicyInPodTopologySpread feature flag.
	// +optional
	// NodeAffinityPolicy 表示我们如何处理 Pod 的 nodeAffinity/nodeSelector，以计算 pod 拓扑扩展偏斜。 选项是： - 荣誉：仅包含与 nodeAffinity/nodeSelector 匹配的节点在计算中。 - 忽略：nodeAffinity/nodeSelector 被忽略。 所有节点都包含在计算中。 如果此值为 nil，则行为等同于荣誉政策。 这是一个 beta 级功能，通过 NodeInclusionPolicyInPodTopologySpread 功能标志默认启用。 +可选
	NodeAffinityPolicy *NodeInclusionPolicy `json:"nodeAffinityPolicy,omitempty" protobuf:"bytes,6,opt,name=nodeAffinityPolicy"`
	// NodeTaintsPolicy indicates how we will treat node taints when calculating
	// pod topology spread skew. Options are:
	// - Honor: nodes without taints, along with tainted nodes for which the incoming pod
	// has a toleration, are included.
	// - Ignore: node taints are ignored. All nodes are included.
	//
	// If this value is nil, the behavior is equivalent to the Ignore policy.
	// This is a beta-level feature default enabled by the NodeInclusionPolicyInPodTopologySpread feature flag.
	// +optional
	// NodeTaintsPolicy 表示我们如何处理 node 污点，以计算 pod 拓扑扩展偏斜。 选项是： - 荣誉：没有污点的节点，以及具有传入 pod 的容忍度的污点节点，都包含在内。 - 忽略：忽略节点污点。 所有节点都包含在内。 如果此值为 nil，则行为等同于忽略策略。 这是一个 beta 级功能，通过 NodeInclusionPolicyInPodTopologySpread 功能标志默认启用。 +可选
	NodeTaintsPolicy *NodeInclusionPolicy `json:"nodeTaintsPolicy,omitempty" protobuf:"bytes,7,opt,name=nodeTaintsPolicy"`
	// MatchLabelKeys is a set of pod label keys to select the pods over which
	// spreading will be calculated. The keys are used to lookup values from the
	// incoming pod labels, those key-value labels are ANDed with labelSelector
	// to select the group of existing pods over which spreading will be calculated
	// for the incoming pod. Keys that don't exist in the incoming pod labels will
	// be ignored. A null or empty list means only match against labelSelector.
	// +listType=atomic
	// +optional
	// MatchLabelKeys是一组 pod 标签键，用于选择要计算扩展的 pod。 该键用于从传入的 pod 标签中查找值，这些键值标签与 labelSelector 进行 AND 运算，以选择要计算传入 pod 的扩展的现有 pod 组。 传入的 pod 标签中不存在的键将被忽略。 空列表或空列表表示仅与 labelSelector 匹配。 +listType=atomic +可选
	MatchLabelKeys []string `json:"matchLabelKeys,omitempty" protobuf:"bytes,8,opt,name=matchLabelKeys"`
}

const (
	// The default value for enableServiceLinks attribute.
	// enableServiceLinks属性的默认值。
	DefaultEnableServiceLinks = true
)

// HostAlias holds the mapping between IP and hostnames that will be injected as an entry in the
// pod's hosts file.
// HostAlias 处理 IP 和主机名之间的映射，该映射将作为 pod 的 hosts 文件中的条目注入。
type HostAlias struct {
	// IP address of the host file entry.
	// IP 地址的主机文件条目。
	IP string `json:"ip,omitempty" protobuf:"bytes,1,opt,name=ip"`
	// Hostnames for the above IP address.
	// Hostnames 用于上述 IP 地址。
	Hostnames []string `json:"hostnames,omitempty" protobuf:"bytes,2,rep,name=hostnames"`
}

// PodFSGroupChangePolicy holds policies that will be used for applying fsGroup to a volume
// when volume is mounted.
// +enum
// PodFSGroupChangePolicy 处理将在挂载卷时将 fsGroup 应用于卷的策略。
type PodFSGroupChangePolicy string

const (
	// FSGroupChangeOnRootMismatch indicates that volume's ownership and permissions will be changed
	// only when permission and ownership of root directory does not match with expected
	// permissions on the volume. This can help shorten the time it takes to change
	// ownership and permissions of a volume.
	// FSGroupChangeOnRootMismatch 表示仅当根目录的权限和所有权与卷上的预期权限不匹配时，才会更改卷的所有权和权限。 这可以帮助缩短更改卷的所有权和权限所需的时间。
	FSGroupChangeOnRootMismatch PodFSGroupChangePolicy = "OnRootMismatch"
	// FSGroupChangeAlways indicates that volume's ownership and permissions
	// should always be changed whenever volume is mounted inside a Pod. This the default
	// behavior.
	// FSGroupChangeAlways 表示无论何时在 Pod 内挂载卷，都应更改卷的所有权和权限。 这是默认行为。
	FSGroupChangeAlways PodFSGroupChangePolicy = "Always"
)

// PodSecurityContext holds pod-level security attributes and common container settings.
// Some fields are also present in container.securityContext.  Field values of
// container.securityContext take precedence over field values of PodSecurityContext.
// PodSecurityContext 保存 pod 级别的安全属性和常见的容器设置。 一些字段也存在于 container.securityContext 中。 container.securityContext 的字段值优先于 PodSecurityContext 的字段值。
type PodSecurityContext struct {
	// The SELinux context to be applied to all containers.
	// If unspecified, the container runtime will allocate a random SELinux context for each
	// container.  May also be set in SecurityContext.  If set in
	// both SecurityContext and PodSecurityContext, the value specified in SecurityContext
	// takes precedence for that container.
	// Note that this field cannot be set when spec.os.name is windows.
	// +optional
	SELinuxOptions *SELinuxOptions `json:"seLinuxOptions,omitempty" protobuf:"bytes,1,opt,name=seLinuxOptions"`
	// The Windows specific settings applied to all containers.
	// If unspecified, the options within a container's SecurityContext will be used.
	// If set in both SecurityContext and PodSecurityContext, the value specified in SecurityContext takes precedence.
	// Note that this field cannot be set when spec.os.name is linux.
	// +optional
	WindowsOptions *WindowsSecurityContextOptions `json:"windowsOptions,omitempty" protobuf:"bytes,8,opt,name=windowsOptions"`
	// The UID to run the entrypoint of the container process.
	// Defaults to user specified in image metadata if unspecified.
	// May also be set in SecurityContext.  If set in both SecurityContext and
	// PodSecurityContext, the value specified in SecurityContext takes precedence
	// for that container.
	// Note that this field cannot be set when spec.os.name is windows.
	// +optional
	RunAsUser *int64 `json:"runAsUser,omitempty" protobuf:"varint,2,opt,name=runAsUser"`
	// The GID to run the entrypoint of the container process.
	// Uses runtime default if unset.
	// May also be set in SecurityContext.  If set in both SecurityContext and
	// PodSecurityContext, the value specified in SecurityContext takes precedence
	// for that container.
	// Note that this field cannot be set when spec.os.name is windows.
	// +optional
	RunAsGroup *int64 `json:"runAsGroup,omitempty" protobuf:"varint,6,opt,name=runAsGroup"`
	// Indicates that the container must run as a non-root user.
	// If true, the Kubelet will validate the image at runtime to ensure that it
	// does not run as UID 0 (root) and fail to start the container if it does.
	// If unset or false, no such validation will be performed.
	// May also be set in SecurityContext.  If set in both SecurityContext and
	// PodSecurityContext, the value specified in SecurityContext takes precedence.
	// +optional
	RunAsNonRoot *bool `json:"runAsNonRoot,omitempty" protobuf:"varint,3,opt,name=runAsNonRoot"`
	// A list of groups applied to the first process run in each container, in addition
	// to the container's primary GID, the fsGroup (if specified), and group memberships
	// defined in the container image for the uid of the container process. If unspecified,
	// no additional groups are added to any container. Note that group memberships
	// defined in the container image for the uid of the container process are still effective,
	// even if they are not included in this list.
	// Note that this field cannot be set when spec.os.name is windows.
	// +optional
	SupplementalGroups []int64 `json:"supplementalGroups,omitempty" protobuf:"varint,4,rep,name=supplementalGroups"`
	// A special supplemental group that applies to all containers in a pod.
	// Some volume types allow the Kubelet to change the ownership of that volume
	// to be owned by the pod:
	//
	// 1. The owning GID will be the FSGroup
	// 2. The setgid bit is set (new files created in the volume will be owned by FSGroup)
	// 3. The permission bits are OR'd with rw-rw----
	//
	// If unset, the Kubelet will not modify the ownership and permissions of any volume.
	// Note that this field cannot be set when spec.os.name is windows.
	// +optional
	FSGroup *int64 `json:"fsGroup,omitempty" protobuf:"varint,5,opt,name=fsGroup"`
	// Sysctls hold a list of namespaced sysctls used for the pod. Pods with unsupported
	// sysctls (by the container runtime) might fail to launch.
	// Note that this field cannot be set when spec.os.name is windows.
	// +optional
	Sysctls []Sysctl `json:"sysctls,omitempty" protobuf:"bytes,7,rep,name=sysctls"`
	// fsGroupChangePolicy defines behavior of changing ownership and permission of the volume
	// before being exposed inside Pod. This field will only apply to
	// volume types which support fsGroup based ownership(and permissions).
	// It will have no effect on ephemeral volume types such as: secret, configmaps
	// and emptydir.
	// Valid values are "OnRootMismatch" and "Always". If not specified, "Always" is used.
	// Note that this field cannot be set when spec.os.name is windows.
	// +optional
	FSGroupChangePolicy *PodFSGroupChangePolicy `json:"fsGroupChangePolicy,omitempty" protobuf:"bytes,9,opt,name=fsGroupChangePolicy"`
	// The seccomp options to use by the containers in this pod.
	// Note that this field cannot be set when spec.os.name is windows.
	// +optional
	SeccompProfile *SeccompProfile `json:"seccompProfile,omitempty" protobuf:"bytes,10,opt,name=seccompProfile"`
}

// SeccompProfile defines a pod/container's seccomp profile settings.
// Only one profile source may be set.
// +union
type SeccompProfile struct {
	// type indicates which kind of seccomp profile will be applied.
	// Valid options are:
	//
	// Localhost - a profile defined in a file on the node should be used.
	// RuntimeDefault - the container runtime default profile should be used.
	// Unconfined - no profile should be applied.
	// +unionDiscriminator
	Type SeccompProfileType `json:"type" protobuf:"bytes,1,opt,name=type,casttype=SeccompProfileType"`
	// localhostProfile indicates a profile defined in a file on the node should be used.
	// The profile must be preconfigured on the node to work.
	// Must be a descending path, relative to the kubelet's configured seccomp profile location.
	// Must only be set if type is "Localhost".
	// +optional
	LocalhostProfile *string `json:"localhostProfile,omitempty" protobuf:"bytes,2,opt,name=localhostProfile"`
}

// SeccompProfileType defines the supported seccomp profile types.
// +enum
type SeccompProfileType string

const (
	// SeccompProfileTypeUnconfined indicates no seccomp profile is applied (A.K.A. unconfined).
	SeccompProfileTypeUnconfined SeccompProfileType = "Unconfined"
	// SeccompProfileTypeRuntimeDefault represents the default container runtime seccomp profile.
	SeccompProfileTypeRuntimeDefault SeccompProfileType = "RuntimeDefault"
	// SeccompProfileTypeLocalhost indicates a profile defined in a file on the node should be used.
	// The file's location relative to <kubelet-root-dir>/seccomp.
	SeccompProfileTypeLocalhost SeccompProfileType = "Localhost"
)

// PodQOSClass defines the supported qos classes of Pods.
// +enum
type PodQOSClass string

const (
	// PodQOSGuaranteed is the Guaranteed qos class.
	PodQOSGuaranteed PodQOSClass = "Guaranteed"
	// PodQOSBurstable is the Burstable qos class.
	PodQOSBurstable PodQOSClass = "Burstable"
	// PodQOSBestEffort is the BestEffort qos class.
	PodQOSBestEffort PodQOSClass = "BestEffort"
)

// PodDNSConfig defines the DNS parameters of a pod in addition to
// those generated from DNSPolicy.
type PodDNSConfig struct {
	// A list of DNS name server IP addresses.
	// This will be appended to the base nameservers generated from DNSPolicy.
	// Duplicated nameservers will be removed.
	// +optional
	Nameservers []string `json:"nameservers,omitempty" protobuf:"bytes,1,rep,name=nameservers"`
	// A list of DNS search domains for host-name lookup.
	// This will be appended to the base search paths generated from DNSPolicy.
	// Duplicated search paths will be removed.
	// +optional
	Searches []string `json:"searches,omitempty" protobuf:"bytes,2,rep,name=searches"`
	// A list of DNS resolver options.
	// This will be merged with the base options generated from DNSPolicy.
	// Duplicated entries will be removed. Resolution options given in Options
	// will override those that appear in the base DNSPolicy.
	// +optional
	Options []PodDNSConfigOption `json:"options,omitempty" protobuf:"bytes,3,rep,name=options"`
}

// PodDNSConfigOption defines DNS resolver options of a pod.
type PodDNSConfigOption struct {
	// Required.
	Name string `json:"name,omitempty" protobuf:"bytes,1,opt,name=name"`
	// +optional
	Value *string `json:"value,omitempty" protobuf:"bytes,2,opt,name=value"`
}

// IP address information for entries in the (plural) PodIPs field.
// Each entry includes:
//
//	IP: An IP address allocated to the pod. Routable at least within the cluster.
type PodIP struct {
	// ip is an IP address (IPv4 or IPv6) assigned to the pod
	IP string `json:"ip,omitempty" protobuf:"bytes,1,opt,name=ip"`
}

// EphemeralContainerCommon is a copy of all fields in Container to be inlined in
// EphemeralContainer. This separate type allows easy conversion from EphemeralContainer
// to Container and allows separate documentation for the fields of EphemeralContainer.
// When a new field is added to Container it must be added here as well.
type EphemeralContainerCommon struct {
	// Name of the ephemeral container specified as a DNS_LABEL.
	// This name must be unique among all containers, init containers and ephemeral containers.
	Name string `json:"name" protobuf:"bytes,1,opt,name=name"`
	// Container image name.
	// More info: https://kubernetes.io/docs/concepts/containers/images
	Image string `json:"image,omitempty" protobuf:"bytes,2,opt,name=image"`
	// Entrypoint array. Not executed within a shell.
	// The image's ENTRYPOINT is used if this is not provided.
	// Variable references $(VAR_NAME) are expanded using the container's environment. If a variable
	// cannot be resolved, the reference in the input string will be unchanged. Double $$ are reduced
	// to a single $, which allows for escaping the $(VAR_NAME) syntax: i.e. "$$(VAR_NAME)" will
	// produce the string literal "$(VAR_NAME)". Escaped references will never be expanded, regardless
	// of whether the variable exists or not. Cannot be updated.
	// More info: https://kubernetes.io/docs/tasks/inject-data-application/define-command-argument-container/#running-a-command-in-a-shell
	// +optional
	Command []string `json:"command,omitempty" protobuf:"bytes,3,rep,name=command"`
	// Arguments to the entrypoint.
	// The image's CMD is used if this is not provided.
	// Variable references $(VAR_NAME) are expanded using the container's environment. If a variable
	// cannot be resolved, the reference in the input string will be unchanged. Double $$ are reduced
	// to a single $, which allows for escaping the $(VAR_NAME) syntax: i.e. "$$(VAR_NAME)" will
	// produce the string literal "$(VAR_NAME)". Escaped references will never be expanded, regardless
	// of whether the variable exists or not. Cannot be updated.
	// More info: https://kubernetes.io/docs/tasks/inject-data-application/define-command-argument-container/#running-a-command-in-a-shell
	// +optional
	Args []string `json:"args,omitempty" protobuf:"bytes,4,rep,name=args"`
	// Container's working directory.
	// If not specified, the container runtime's default will be used, which
	// might be configured in the container image.
	// Cannot be updated.
	// +optional
	WorkingDir string `json:"workingDir,omitempty" protobuf:"bytes,5,opt,name=workingDir"`
	// Ports are not allowed for ephemeral containers.
	// +optional
	// +patchMergeKey=containerPort
	// +patchStrategy=merge
	// +listType=map
	// +listMapKey=containerPort
	// +listMapKey=protocol
	Ports []ContainerPort `json:"ports,omitempty" patchStrategy:"merge" patchMergeKey:"containerPort" protobuf:"bytes,6,rep,name=ports"`
	// List of sources to populate environment variables in the container.
	// The keys defined within a source must be a C_IDENTIFIER. All invalid keys
	// will be reported as an event when the container is starting. When a key exists in multiple
	// sources, the value associated with the last source will take precedence.
	// Values defined by an Env with a duplicate key will take precedence.
	// Cannot be updated.
	// +optional
	EnvFrom []EnvFromSource `json:"envFrom,omitempty" protobuf:"bytes,19,rep,name=envFrom"`
	// List of environment variables to set in the container.
	// Cannot be updated.
	// +optional
	// +patchMergeKey=name
	// +patchStrategy=merge
	Env []EnvVar `json:"env,omitempty" patchStrategy:"merge" patchMergeKey:"name" protobuf:"bytes,7,rep,name=env"`
	// Resources are not allowed for ephemeral containers. Ephemeral containers use spare resources
	// already allocated to the pod.
	// +optional
	Resources ResourceRequirements `json:"resources,omitempty" protobuf:"bytes,8,opt,name=resources"`
	// Pod volumes to mount into the container's filesystem. Subpath mounts are not allowed for ephemeral containers.
	// Cannot be updated.
	// +optional
	// +patchMergeKey=mountPath
	// +patchStrategy=merge
	VolumeMounts []VolumeMount `json:"volumeMounts,omitempty" patchStrategy:"merge" patchMergeKey:"mountPath" protobuf:"bytes,9,rep,name=volumeMounts"`
	// volumeDevices is the list of block devices to be used by the container.
	// +patchMergeKey=devicePath
	// +patchStrategy=merge
	// +optional
	VolumeDevices []VolumeDevice `json:"volumeDevices,omitempty" patchStrategy:"merge" patchMergeKey:"devicePath" protobuf:"bytes,21,rep,name=volumeDevices"`
	// Probes are not allowed for ephemeral containers.
	// +optional
	LivenessProbe *Probe `json:"livenessProbe,omitempty" protobuf:"bytes,10,opt,name=livenessProbe"`
	// Probes are not allowed for ephemeral containers.
	// +optional
	ReadinessProbe *Probe `json:"readinessProbe,omitempty" protobuf:"bytes,11,opt,name=readinessProbe"`
	// Probes are not allowed for ephemeral containers.
	// +optional
	StartupProbe *Probe `json:"startupProbe,omitempty" protobuf:"bytes,22,opt,name=startupProbe"`
	// Lifecycle is not allowed for ephemeral containers.
	// +optional
	Lifecycle *Lifecycle `json:"lifecycle,omitempty" protobuf:"bytes,12,opt,name=lifecycle"`
	// Optional: Path at which the file to which the container's termination message
	// will be written is mounted into the container's filesystem.
	// Message written is intended to be brief final status, such as an assertion failure message.
	// Will be truncated by the node if greater than 4096 bytes. The total message length across
	// all containers will be limited to 12kb.
	// Defaults to /dev/termination-log.
	// Cannot be updated.
	// +optional
	TerminationMessagePath string `json:"terminationMessagePath,omitempty" protobuf:"bytes,13,opt,name=terminationMessagePath"`
	// Indicate how the termination message should be populated. File will use the contents of
	// terminationMessagePath to populate the container status message on both success and failure.
	// FallbackToLogsOnError will use the last chunk of container log output if the termination
	// message file is empty and the container exited with an error.
	// The log output is limited to 2048 bytes or 80 lines, whichever is smaller.
	// Defaults to File.
	// Cannot be updated.
	// +optional
	TerminationMessagePolicy TerminationMessagePolicy `json:"terminationMessagePolicy,omitempty" protobuf:"bytes,20,opt,name=terminationMessagePolicy,casttype=TerminationMessagePolicy"`
	// Image pull policy.
	// One of Always, Never, IfNotPresent.
	// Defaults to Always if :latest tag is specified, or IfNotPresent otherwise.
	// Cannot be updated.
	// More info: https://kubernetes.io/docs/concepts/containers/images#updating-images
	// +optional
	ImagePullPolicy PullPolicy `json:"imagePullPolicy,omitempty" protobuf:"bytes,14,opt,name=imagePullPolicy,casttype=PullPolicy"`
	// Optional: SecurityContext defines the security options the ephemeral container should be run with.
	// If set, the fields of SecurityContext override the equivalent fields of PodSecurityContext.
	// +optional
	SecurityContext *SecurityContext `json:"securityContext,omitempty" protobuf:"bytes,15,opt,name=securityContext"`

	// Variables for interactive containers, these have very specialized use-cases (e.g. debugging)
	// and shouldn't be used for general purpose containers.

	// Whether this container should allocate a buffer for stdin in the container runtime. If this
	// is not set, reads from stdin in the container will always result in EOF.
	// Default is false.
	// +optional
	Stdin bool `json:"stdin,omitempty" protobuf:"varint,16,opt,name=stdin"`
	// Whether the container runtime should close the stdin channel after it has been opened by
	// a single attach. When stdin is true the stdin stream will remain open across multiple attach
	// sessions. If stdinOnce is set to true, stdin is opened on container start, is empty until the
	// first client attaches to stdin, and then remains open and accepts data until the client disconnects,
	// at which time stdin is closed and remains closed until the container is restarted. If this
	// flag is false, a container processes that reads from stdin will never receive an EOF.
	// Default is false
	// +optional
	StdinOnce bool `json:"stdinOnce,omitempty" protobuf:"varint,17,opt,name=stdinOnce"`
	// Whether this container should allocate a TTY for itself, also requires 'stdin' to be true.
	// Default is false.
	// +optional
	TTY bool `json:"tty,omitempty" protobuf:"varint,18,opt,name=tty"`
}

// EphemeralContainerCommon converts to Container. All fields must be kept in sync between
// these two types.
var _ = Container(EphemeralContainerCommon{})

// An EphemeralContainer is a temporary container that you may add to an existing Pod for
// user-initiated activities such as debugging. Ephemeral containers have no resource or
// scheduling guarantees, and they will not be restarted when they exit or when a Pod is
// removed or restarted. The kubelet may evict a Pod if an ephemeral container causes the
// Pod to exceed its resource allocation.
//
// To add an ephemeral container, use the ephemeralcontainers subresource of an existing
// Pod. Ephemeral containers may not be removed or restarted.
type EphemeralContainer struct {
	// Ephemeral containers have all of the fields of Container, plus additional fields
	// specific to ephemeral containers. Fields in common with Container are in the
	// following inlined struct so than an EphemeralContainer may easily be converted
	// to a Container.
	EphemeralContainerCommon `json:",inline" protobuf:"bytes,1,req"`

	// If set, the name of the container from PodSpec that this ephemeral container targets.
	// The ephemeral container will be run in the namespaces (IPC, PID, etc) of this container.
	// If not set then the ephemeral container uses the namespaces configured in the Pod spec.
	//
	// The container runtime must implement support for this feature. If the runtime does not
	// support namespace targeting then the result of setting this field is undefined.
	// +optional
	TargetContainerName string `json:"targetContainerName,omitempty" protobuf:"bytes,2,opt,name=targetContainerName"`
}

// PodStatus represents information about the status of a pod. Status may trail the actual
// state of a system, especially if the node that hosts the pod cannot contact the control
// plane.
type PodStatus struct {
	// The phase of a Pod is a simple, high-level summary of where the Pod is in its lifecycle.
	// The conditions array, the reason and message fields, and the individual container status
	// arrays contain more detail about the pod's status.
	// There are five possible phase values:
	//
	// Pending: The pod has been accepted by the Kubernetes system, but one or more of the
	// container images has not been created. This includes time before being scheduled as
	// well as time spent downloading images over the network, which could take a while.
	// Running: The pod has been bound to a node, and all of the containers have been created.
	// At least one container is still running, or is in the process of starting or restarting.
	// Succeeded: All containers in the pod have terminated in success, and will not be restarted.
	// Failed: All containers in the pod have terminated, and at least one container has
	// terminated in failure. The container either exited with non-zero status or was terminated
	// by the system.
	// Unknown: For some reason the state of the pod could not be obtained, typically due to an
	// error in communicating with the host of the pod.
	//
	// More info: https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle#pod-phase
	// +optional
	Phase PodPhase `json:"phase,omitempty" protobuf:"bytes,1,opt,name=phase,casttype=PodPhase"`
	// Current service state of pod.
	// More info: https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle#pod-conditions
	// +optional
	// +patchMergeKey=type
	// +patchStrategy=merge
	Conditions []PodCondition `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type" protobuf:"bytes,2,rep,name=conditions"`
	// A human readable message indicating details about why the pod is in this condition.
	// +optional
	Message string `json:"message,omitempty" protobuf:"bytes,3,opt,name=message"`
	// A brief CamelCase message indicating details about why the pod is in this state.
	// e.g. 'Evicted'
	// +optional
	Reason string `json:"reason,omitempty" protobuf:"bytes,4,opt,name=reason"`
	// nominatedNodeName is set only when this pod preempts other pods on the node, but it cannot be
	// scheduled right away as preemption victims receive their graceful termination periods.
	// This field does not guarantee that the pod will be scheduled on this node. Scheduler may decide
	// to place the pod elsewhere if other nodes become available sooner. Scheduler may also decide to
	// give the resources on this node to a higher priority pod that is created after preemption.
	// As a result, this field may be different than PodSpec.nodeName when the pod is
	// scheduled.
	// +optional
	NominatedNodeName string `json:"nominatedNodeName,omitempty" protobuf:"bytes,11,opt,name=nominatedNodeName"`

	// IP address of the host to which the pod is assigned. Empty if not yet scheduled.
	// +optional
	HostIP string `json:"hostIP,omitempty" protobuf:"bytes,5,opt,name=hostIP"`
	// IP address allocated to the pod. Routable at least within the cluster.
	// Empty if not yet allocated.
	// +optional
	PodIP string `json:"podIP,omitempty" protobuf:"bytes,6,opt,name=podIP"`

	// podIPs holds the IP addresses allocated to the pod. If this field is specified, the 0th entry must
	// match the podIP field. Pods may be allocated at most 1 value for each of IPv4 and IPv6. This list
	// is empty if no IPs have been allocated yet.
	// +optional
	// +patchStrategy=merge
	// +patchMergeKey=ip
	PodIPs []PodIP `json:"podIPs,omitempty" protobuf:"bytes,12,rep,name=podIPs" patchStrategy:"merge" patchMergeKey:"ip"`

	// RFC 3339 date and time at which the object was acknowledged by the Kubelet.
	// This is before the Kubelet pulled the container image(s) for the pod.
	// +optional
	StartTime *metav1.Time `json:"startTime,omitempty" protobuf:"bytes,7,opt,name=startTime"`

	// The list has one entry per init container in the manifest. The most recent successful
	// init container will have ready = true, the most recently started container will have
	// startTime set.
	// More info: https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle#pod-and-container-status
	InitContainerStatuses []ContainerStatus `json:"initContainerStatuses,omitempty" protobuf:"bytes,10,rep,name=initContainerStatuses"`

	// The list has one entry per container in the manifest.
	// More info: https://kubernetes.io/docs/concepts/workloads/pods/pod-lifecycle#pod-and-container-status
	// +optional
	ContainerStatuses []ContainerStatus `json:"containerStatuses,omitempty" protobuf:"bytes,8,rep,name=containerStatuses"`
	// The Quality of Service (QOS) classification assigned to the pod based on resource requirements
	// See PodQOSClass type for available QOS classes
	// More info: https://git.k8s.io/community/contributors/design-proposals/node/resource-qos.md
	// +optional
	QOSClass PodQOSClass `json:"qosClass,omitempty" protobuf:"bytes,9,rep,name=qosClass"`
	// Status for any ephemeral containers that have run in this pod.
	// +optional
	EphemeralContainerStatuses []ContainerStatus `json:"ephemeralContainerStatuses,omitempty" protobuf:"bytes,13,rep,name=ephemeralContainerStatuses"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// PodStatusResult is a wrapper for PodStatus returned by kubelet that can be encode/decoded
type PodStatusResult struct {
	metav1.TypeMeta `json:",inline"`
	// Standard object's metadata.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`
	// Most recently observed status of the pod.
	// This data may not be up to date.
	// Populated by the system.
	// Read-only.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#spec-and-status
	// +optional
	Status PodStatus `json:"status,omitempty" protobuf:"bytes,2,opt,name=status"`
}

// +genclient
// +genclient:method=UpdateEphemeralContainers,verb=update,subresource=ephemeralcontainers
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Pod is a collection of containers that can run on a host. This resource is created
// by clients and scheduled onto hosts.
type Pod struct {
	metav1.TypeMeta `json:",inline"`
	// Standard object's metadata.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// Specification of the desired behavior of the pod.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#spec-and-status
	// +optional
	Spec PodSpec `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`

	// Most recently observed status of the pod.
	// This data may not be up to date.
	// Populated by the system.
	// Read-only.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#spec-and-status
	// +optional
	Status PodStatus `json:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// PodList is a list of Pods.
type PodList struct {
	metav1.TypeMeta `json:",inline"`
	// Standard list metadata.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds
	// +optional
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// List of pods.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md
	Items []Pod `json:"items" protobuf:"bytes,2,rep,name=items"`
}

// PodTemplateSpec describes the data a pod should have when created from a template
// PodTemplateSpec 描述了从模板创建 Pod 时应该具有的数据
type PodTemplateSpec struct {
	// Standard object's metadata.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// Specification of the desired behavior of the pod.
	// 描述了 Pod 应该具有的行为
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#spec-and-status
	// +optional
	Spec PodSpec `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// PodTemplate describes a template for creating copies of a predefined pod.
type PodTemplate struct {
	metav1.TypeMeta `json:",inline"`
	// Standard object's metadata.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// Template defines the pods that will be created from this pod template.
	// https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#spec-and-status
	// +optional
	Template PodTemplateSpec `json:"template,omitempty" protobuf:"bytes,2,opt,name=template"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// PodTemplateList is a list of PodTemplates.
type PodTemplateList struct {
	metav1.TypeMeta `json:",inline"`
	// Standard list metadata.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds
	// +optional
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// List of pod templates
	Items []PodTemplate `json:"items" protobuf:"bytes,2,rep,name=items"`
}

// ReplicationControllerSpec is the specification of a replication controller.
type ReplicationControllerSpec struct {
	// Replicas is the number of desired replicas.
	// This is a pointer to distinguish between explicit zero and unspecified.
	// Defaults to 1.
	// More info: https://kubernetes.io/docs/concepts/workloads/controllers/replicationcontroller#what-is-a-replicationcontroller
	// +optional
	Replicas *int32 `json:"replicas,omitempty" protobuf:"varint,1,opt,name=replicas"`

	// Minimum number of seconds for which a newly created pod should be ready
	// without any of its container crashing, for it to be considered available.
	// Defaults to 0 (pod will be considered available as soon as it is ready)
	// +optional
	MinReadySeconds int32 `json:"minReadySeconds,omitempty" protobuf:"varint,4,opt,name=minReadySeconds"`

	// Selector is a label query over pods that should match the Replicas count.
	// If Selector is empty, it is defaulted to the labels present on the Pod template.
	// Label keys and values that must match in order to be controlled by this replication
	// controller, if empty defaulted to labels on Pod template.
	// More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/labels/#label-selectors
	// +optional
	// +mapType=atomic
	Selector map[string]string `json:"selector,omitempty" protobuf:"bytes,2,rep,name=selector"`

	// TemplateRef is a reference to an object that describes the pod that will be created if
	// insufficient replicas are detected.
	// Reference to an object that describes the pod that will be created if insufficient replicas are detected.
	// +optional
	// TemplateRef *ObjectReference `json:"templateRef,omitempty"`

	// Template is the object that describes the pod that will be created if
	// insufficient replicas are detected. This takes precedence over a TemplateRef.
	// More info: https://kubernetes.io/docs/concepts/workloads/controllers/replicationcontroller#pod-template
	// +optional
	Template *PodTemplateSpec `json:"template,omitempty" protobuf:"bytes,3,opt,name=template"`
}

// ReplicationControllerStatus represents the current status of a replication
// controller.
type ReplicationControllerStatus struct {
	// Replicas is the most recently observed number of replicas.
	// More info: https://kubernetes.io/docs/concepts/workloads/controllers/replicationcontroller#what-is-a-replicationcontroller
	Replicas int32 `json:"replicas" protobuf:"varint,1,opt,name=replicas"`

	// The number of pods that have labels matching the labels of the pod template of the replication controller.
	// +optional
	FullyLabeledReplicas int32 `json:"fullyLabeledReplicas,omitempty" protobuf:"varint,2,opt,name=fullyLabeledReplicas"`

	// The number of ready replicas for this replication controller.
	// +optional
	ReadyReplicas int32 `json:"readyReplicas,omitempty" protobuf:"varint,4,opt,name=readyReplicas"`

	// The number of available replicas (ready for at least minReadySeconds) for this replication controller.
	// +optional
	AvailableReplicas int32 `json:"availableReplicas,omitempty" protobuf:"varint,5,opt,name=availableReplicas"`

	// ObservedGeneration reflects the generation of the most recently observed replication controller.
	// +optional
	ObservedGeneration int64 `json:"observedGeneration,omitempty" protobuf:"varint,3,opt,name=observedGeneration"`

	// Represents the latest available observations of a replication controller's current state.
	// +optional
	// +patchMergeKey=type
	// +patchStrategy=merge
	Conditions []ReplicationControllerCondition `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type" protobuf:"bytes,6,rep,name=conditions"`
}

type ReplicationControllerConditionType string

// These are valid conditions of a replication controller.
const (
	// ReplicationControllerReplicaFailure is added in a replication controller when one of its pods
	// fails to be created due to insufficient quota, limit ranges, pod security policy, node selectors,
	// etc. or deleted due to kubelet being down or finalizers are failing.
	ReplicationControllerReplicaFailure ReplicationControllerConditionType = "ReplicaFailure"
)

// ReplicationControllerCondition describes the state of a replication controller at a certain point.
type ReplicationControllerCondition struct {
	// Type of replication controller condition.
	Type ReplicationControllerConditionType `json:"type" protobuf:"bytes,1,opt,name=type,casttype=ReplicationControllerConditionType"`
	// Status of the condition, one of True, False, Unknown.
	Status ConditionStatus `json:"status" protobuf:"bytes,2,opt,name=status,casttype=ConditionStatus"`
	// The last time the condition transitioned from one status to another.
	// +optional
	LastTransitionTime metav1.Time `json:"lastTransitionTime,omitempty" protobuf:"bytes,3,opt,name=lastTransitionTime"`
	// The reason for the condition's last transition.
	// +optional
	Reason string `json:"reason,omitempty" protobuf:"bytes,4,opt,name=reason"`
	// A human readable message indicating details about the transition.
	// +optional
	Message string `json:"message,omitempty" protobuf:"bytes,5,opt,name=message"`
}

// +genclient
// +genclient:method=GetScale,verb=get,subresource=scale,result=k8s.io/api/autoscaling/v1.Scale
// +genclient:method=UpdateScale,verb=update,subresource=scale,input=k8s.io/api/autoscaling/v1.Scale,result=k8s.io/api/autoscaling/v1.Scale
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ReplicationController represents the configuration of a replication controller.
type ReplicationController struct {
	metav1.TypeMeta `json:",inline"`

	// If the Labels of a ReplicationController are empty, they are defaulted to
	// be the same as the Pod(s) that the replication controller manages.
	// Standard object's metadata. More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// Spec defines the specification of the desired behavior of the replication controller.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#spec-and-status
	// +optional
	Spec ReplicationControllerSpec `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`

	// Status is the most recently observed status of the replication controller.
	// This data may be out of date by some window of time.
	// Populated by the system.
	// Read-only.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#spec-and-status
	// +optional
	Status ReplicationControllerStatus `json:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ReplicationControllerList is a collection of replication controllers.
type ReplicationControllerList struct {
	metav1.TypeMeta `json:",inline"`
	// Standard list metadata.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds
	// +optional
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// List of replication controllers.
	// More info: https://kubernetes.io/docs/concepts/workloads/controllers/replicationcontroller
	Items []ReplicationController `json:"items" protobuf:"bytes,2,rep,name=items"`
}

// Session Affinity Type string
// +enum
type ServiceAffinity string

const (
	// ServiceAffinityClientIP is the Client IP based.
	ServiceAffinityClientIP ServiceAffinity = "ClientIP"

	// ServiceAffinityNone - no session affinity.
	ServiceAffinityNone ServiceAffinity = "None"
)

const DefaultClientIPServiceAffinitySeconds int32 = 10800

// SessionAffinityConfig represents the configurations of session affinity.
type SessionAffinityConfig struct {
	// clientIP contains the configurations of Client IP based session affinity.
	// +optional
	ClientIP *ClientIPConfig `json:"clientIP,omitempty" protobuf:"bytes,1,opt,name=clientIP"`
}

// ClientIPConfig represents the configurations of Client IP based session affinity.
type ClientIPConfig struct {
	// timeoutSeconds specifies the seconds of ClientIP type session sticky time.
	// The value must be >0 && <=86400(for 1 day) if ServiceAffinity == "ClientIP".
	// Default value is 10800(for 3 hours).
	// +optional
	TimeoutSeconds *int32 `json:"timeoutSeconds,omitempty" protobuf:"varint,1,opt,name=timeoutSeconds"`
}

// Service Type string describes ingress methods for a service
// +enum
type ServiceType string

const (
	// ServiceTypeClusterIP means a service will only be accessible inside the
	// cluster, via the cluster IP.
	ServiceTypeClusterIP ServiceType = "ClusterIP"

	// ServiceTypeNodePort means a service will be exposed on one port of
	// every node, in addition to 'ClusterIP' type.
	ServiceTypeNodePort ServiceType = "NodePort"

	// ServiceTypeLoadBalancer means a service will be exposed via an
	// external load balancer (if the cloud provider supports it), in addition
	// to 'NodePort' type.
	ServiceTypeLoadBalancer ServiceType = "LoadBalancer"

	// ServiceTypeExternalName means a service consists of only a reference to
	// an external name that kubedns or equivalent will return as a CNAME
	// record, with no exposing or proxying of any pods involved.
	ServiceTypeExternalName ServiceType = "ExternalName"
)

// ServiceInternalTrafficPolicy describes how nodes distribute service traffic they
// receive on the ClusterIP.
// +enum
type ServiceInternalTrafficPolicy string

const (
	// ServiceInternalTrafficPolicyCluster routes traffic to all endpoints.
	ServiceInternalTrafficPolicyCluster ServiceInternalTrafficPolicy = "Cluster"

	// ServiceInternalTrafficPolicyLocal routes traffic only to endpoints on the same
	// node as the client pod (dropping the traffic if there are no local endpoints).
	ServiceInternalTrafficPolicyLocal ServiceInternalTrafficPolicy = "Local"
)

// for backwards compat
// +enum
type ServiceInternalTrafficPolicyType = ServiceInternalTrafficPolicy

// ServiceExternalTrafficPolicy describes how nodes distribute service traffic they
// receive on one of the Service's "externally-facing" addresses (NodePorts, ExternalIPs,
// and LoadBalancer IPs.
// +enum
type ServiceExternalTrafficPolicy string

const (
	// ServiceExternalTrafficPolicyCluster routes traffic to all endpoints.
	ServiceExternalTrafficPolicyCluster ServiceExternalTrafficPolicy = "Cluster"

	// ServiceExternalTrafficPolicyLocal preserves the source IP of the traffic by
	// routing only to endpoints on the same node as the traffic was received on
	// (dropping the traffic if there are no local endpoints).
	ServiceExternalTrafficPolicyLocal ServiceExternalTrafficPolicy = "Local"
)

// for backwards compat
// +enum
type ServiceExternalTrafficPolicyType = ServiceExternalTrafficPolicy

const (
	ServiceExternalTrafficPolicyTypeLocal   = ServiceExternalTrafficPolicyLocal
	ServiceExternalTrafficPolicyTypeCluster = ServiceExternalTrafficPolicyCluster
)

// These are the valid conditions of a service.
const (
	// LoadBalancerPortsError represents the condition of the requested ports
	// on the cloud load balancer instance.
	LoadBalancerPortsError = "LoadBalancerPortsError"
)

// ServiceStatus represents the current status of a service.
type ServiceStatus struct {
	// LoadBalancer contains the current status of the load-balancer,
	// if one is present.
	// +optional
	LoadBalancer LoadBalancerStatus `json:"loadBalancer,omitempty" protobuf:"bytes,1,opt,name=loadBalancer"`
	// Current service state
	// +optional
	// +patchMergeKey=type
	// +patchStrategy=merge
	// +listType=map
	// +listMapKey=type
	Conditions []metav1.Condition `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type" protobuf:"bytes,2,rep,name=conditions"`
}

// LoadBalancerStatus represents the status of a load-balancer.
type LoadBalancerStatus struct {
	// Ingress is a list containing ingress points for the load-balancer.
	// Traffic intended for the service should be sent to these ingress points.
	// +optional
	Ingress []LoadBalancerIngress `json:"ingress,omitempty" protobuf:"bytes,1,rep,name=ingress"`
}

// LoadBalancerIngress represents the status of a load-balancer ingress point:
// traffic intended for the service should be sent to an ingress point.
type LoadBalancerIngress struct {
	// IP is set for load-balancer ingress points that are IP based
	// (typically GCE or OpenStack load-balancers)
	// +optional
	IP string `json:"ip,omitempty" protobuf:"bytes,1,opt,name=ip"`

	// Hostname is set for load-balancer ingress points that are DNS based
	// (typically AWS load-balancers)
	// +optional
	Hostname string `json:"hostname,omitempty" protobuf:"bytes,2,opt,name=hostname"`

	// Ports is a list of records of service ports
	// If used, every port defined in the service should have an entry in it
	// +listType=atomic
	// +optional
	Ports []PortStatus `json:"ports,omitempty" protobuf:"bytes,4,rep,name=ports"`
}

// IPFamily represents the IP Family (IPv4 or IPv6). This type is used
// to express the family of an IP expressed by a type (e.g. service.spec.ipFamilies).
// +enum
type IPFamily string

const (
	// IPv4Protocol indicates that this IP is IPv4 protocol
	IPv4Protocol IPFamily = "IPv4"
	// IPv6Protocol indicates that this IP is IPv6 protocol
	IPv6Protocol IPFamily = "IPv6"
)

// IPFamilyPolicy represents the dual-stack-ness requested or required by a Service
// +enum
type IPFamilyPolicy string

const (
	// IPFamilyPolicySingleStack indicates that this service is required to have a single IPFamily.
	// The IPFamily assigned is based on the default IPFamily used by the cluster
	// or as identified by service.spec.ipFamilies field
	IPFamilyPolicySingleStack IPFamilyPolicy = "SingleStack"
	// IPFamilyPolicyPreferDualStack indicates that this service prefers dual-stack when
	// the cluster is configured for dual-stack. If the cluster is not configured
	// for dual-stack the service will be assigned a single IPFamily. If the IPFamily is not
	// set in service.spec.ipFamilies then the service will be assigned the default IPFamily
	// configured on the cluster
	IPFamilyPolicyPreferDualStack IPFamilyPolicy = "PreferDualStack"
	// IPFamilyPolicyRequireDualStack indicates that this service requires dual-stack. Using
	// IPFamilyPolicyRequireDualStack on a single stack cluster will result in validation errors. The
	// IPFamilies (and their order) assigned  to this service is based on service.spec.ipFamilies. If
	// service.spec.ipFamilies was not provided then it will be assigned according to how they are
	// configured on the cluster. If service.spec.ipFamilies has only one entry then the alternative
	// IPFamily will be added by apiserver
	IPFamilyPolicyRequireDualStack IPFamilyPolicy = "RequireDualStack"
)

// for backwards compat
// +enum
type IPFamilyPolicyType = IPFamilyPolicy

// ServiceSpec describes the attributes that a user creates on a service.
type ServiceSpec struct {
	// The list of ports that are exposed by this service.
	// More info: https://kubernetes.io/docs/concepts/services-networking/service/#virtual-ips-and-service-proxies
	// +patchMergeKey=port
	// +patchStrategy=merge
	// +listType=map
	// +listMapKey=port
	// +listMapKey=protocol
	Ports []ServicePort `json:"ports,omitempty" patchStrategy:"merge" patchMergeKey:"port" protobuf:"bytes,1,rep,name=ports"`

	// Route service traffic to pods with label keys and values matching this
	// selector. If empty or not present, the service is assumed to have an
	// external process managing its endpoints, which Kubernetes will not
	// modify. Only applies to types ClusterIP, NodePort, and LoadBalancer.
	// Ignored if type is ExternalName.
	// More info: https://kubernetes.io/docs/concepts/services-networking/service/
	// +optional
	// +mapType=atomic
	Selector map[string]string `json:"selector,omitempty" protobuf:"bytes,2,rep,name=selector"`

	// clusterIP is the IP address of the service and is usually assigned
	// randomly. If an address is specified manually, is in-range (as per
	// system configuration), and is not in use, it will be allocated to the
	// service; otherwise creation of the service will fail. This field may not
	// be changed through updates unless the type field is also being changed
	// to ExternalName (which requires this field to be blank) or the type
	// field is being changed from ExternalName (in which case this field may
	// optionally be specified, as describe above).  Valid values are "None",
	// empty string (""), or a valid IP address. Setting this to "None" makes a
	// "headless service" (no virtual IP), which is useful when direct endpoint
	// connections are preferred and proxying is not required.  Only applies to
	// types ClusterIP, NodePort, and LoadBalancer. If this field is specified
	// when creating a Service of type ExternalName, creation will fail. This
	// field will be wiped when updating a Service to type ExternalName.
	// More info: https://kubernetes.io/docs/concepts/services-networking/service/#virtual-ips-and-service-proxies
	// +optional
	ClusterIP string `json:"clusterIP,omitempty" protobuf:"bytes,3,opt,name=clusterIP"`

	// ClusterIPs is a list of IP addresses assigned to this service, and are
	// usually assigned randomly.  If an address is specified manually, is
	// in-range (as per system configuration), and is not in use, it will be
	// allocated to the service; otherwise creation of the service will fail.
	// This field may not be changed through updates unless the type field is
	// also being changed to ExternalName (which requires this field to be
	// empty) or the type field is being changed from ExternalName (in which
	// case this field may optionally be specified, as describe above).  Valid
	// values are "None", empty string (""), or a valid IP address.  Setting
	// this to "None" makes a "headless service" (no virtual IP), which is
	// useful when direct endpoint connections are preferred and proxying is
	// not required.  Only applies to types ClusterIP, NodePort, and
	// LoadBalancer. If this field is specified when creating a Service of type
	// ExternalName, creation will fail. This field will be wiped when updating
	// a Service to type ExternalName.  If this field is not specified, it will
	// be initialized from the clusterIP field.  If this field is specified,
	// clients must ensure that clusterIPs[0] and clusterIP have the same
	// value.
	//
	// This field may hold a maximum of two entries (dual-stack IPs, in either order).
	// These IPs must correspond to the values of the ipFamilies field. Both
	// clusterIPs and ipFamilies are governed by the ipFamilyPolicy field.
	// More info: https://kubernetes.io/docs/concepts/services-networking/service/#virtual-ips-and-service-proxies
	// +listType=atomic
	// +optional
	ClusterIPs []string `json:"clusterIPs,omitempty" protobuf:"bytes,18,opt,name=clusterIPs"`

	// type determines how the Service is exposed. Defaults to ClusterIP. Valid
	// options are ExternalName, ClusterIP, NodePort, and LoadBalancer.
	// "ClusterIP" allocates a cluster-internal IP address for load-balancing
	// to endpoints. Endpoints are determined by the selector or if that is not
	// specified, by manual construction of an Endpoints object or
	// EndpointSlice objects. If clusterIP is "None", no virtual IP is
	// allocated and the endpoints are published as a set of endpoints rather
	// than a virtual IP.
	// "NodePort" builds on ClusterIP and allocates a port on every node which
	// routes to the same endpoints as the clusterIP.
	// "LoadBalancer" builds on NodePort and creates an external load-balancer
	// (if supported in the current cloud) which routes to the same endpoints
	// as the clusterIP.
	// "ExternalName" aliases this service to the specified externalName.
	// Several other fields do not apply to ExternalName services.
	// More info: https://kubernetes.io/docs/concepts/services-networking/service/#publishing-services-service-types
	// +optional
	Type ServiceType `json:"type,omitempty" protobuf:"bytes,4,opt,name=type,casttype=ServiceType"`

	// externalIPs is a list of IP addresses for which nodes in the cluster
	// will also accept traffic for this service.  These IPs are not managed by
	// Kubernetes.  The user is responsible for ensuring that traffic arrives
	// at a node with this IP.  A common example is external load-balancers
	// that are not part of the Kubernetes system.
	// +optional
	ExternalIPs []string `json:"externalIPs,omitempty" protobuf:"bytes,5,rep,name=externalIPs"`

	// Supports "ClientIP" and "None". Used to maintain session affinity.
	// Enable client IP based session affinity.
	// Must be ClientIP or None.
	// Defaults to None.
	// More info: https://kubernetes.io/docs/concepts/services-networking/service/#virtual-ips-and-service-proxies
	// +optional
	SessionAffinity ServiceAffinity `json:"sessionAffinity,omitempty" protobuf:"bytes,7,opt,name=sessionAffinity,casttype=ServiceAffinity"`

	// Only applies to Service Type: LoadBalancer.
	// This feature depends on whether the underlying cloud-provider supports specifying
	// the loadBalancerIP when a load balancer is created.
	// This field will be ignored if the cloud-provider does not support the feature.
	// Deprecated: This field was under-specified and its meaning varies across implementations,
	// and it cannot support dual-stack.
	// As of Kubernetes v1.24, users are encouraged to use implementation-specific annotations when available.
	// This field may be removed in a future API version.
	// +optional
	LoadBalancerIP string `json:"loadBalancerIP,omitempty" protobuf:"bytes,8,opt,name=loadBalancerIP"`

	// If specified and supported by the platform, this will restrict traffic through the cloud-provider
	// load-balancer will be restricted to the specified client IPs. This field will be ignored if the
	// cloud-provider does not support the feature."
	// More info: https://kubernetes.io/docs/tasks/access-application-cluster/create-external-load-balancer/
	// +optional
	LoadBalancerSourceRanges []string `json:"loadBalancerSourceRanges,omitempty" protobuf:"bytes,9,opt,name=loadBalancerSourceRanges"`

	// externalName is the external reference that discovery mechanisms will
	// return as an alias for this service (e.g. a DNS CNAME record). No
	// proxying will be involved.  Must be a lowercase RFC-1123 hostname
	// (https://tools.ietf.org/html/rfc1123) and requires `type` to be "ExternalName".
	// +optional
	ExternalName string `json:"externalName,omitempty" protobuf:"bytes,10,opt,name=externalName"`

	// externalTrafficPolicy describes how nodes distribute service traffic they
	// receive on one of the Service's "externally-facing" addresses (NodePorts,
	// ExternalIPs, and LoadBalancer IPs). If set to "Local", the proxy will configure
	// the service in a way that assumes that external load balancers will take care
	// of balancing the service traffic between nodes, and so each node will deliver
	// traffic only to the node-local endpoints of the service, without masquerading
	// the client source IP. (Traffic mistakenly sent to a node with no endpoints will
	// be dropped.) The default value, "Cluster", uses the standard behavior of
	// routing to all endpoints evenly (possibly modified by topology and other
	// features). Note that traffic sent to an External IP or LoadBalancer IP from
	// within the cluster will always get "Cluster" semantics, but clients sending to
	// a NodePort from within the cluster may need to take traffic policy into account
	// when picking a node.
	// +optional
	ExternalTrafficPolicy ServiceExternalTrafficPolicy `json:"externalTrafficPolicy,omitempty" protobuf:"bytes,11,opt,name=externalTrafficPolicy"`

	// healthCheckNodePort specifies the healthcheck nodePort for the service.
	// This only applies when type is set to LoadBalancer and
	// externalTrafficPolicy is set to Local. If a value is specified, is
	// in-range, and is not in use, it will be used.  If not specified, a value
	// will be automatically allocated.  External systems (e.g. load-balancers)
	// can use this port to determine if a given node holds endpoints for this
	// service or not.  If this field is specified when creating a Service
	// which does not need it, creation will fail. This field will be wiped
	// when updating a Service to no longer need it (e.g. changing type).
	// This field cannot be updated once set.
	// +optional
	HealthCheckNodePort int32 `json:"healthCheckNodePort,omitempty" protobuf:"bytes,12,opt,name=healthCheckNodePort"`

	// publishNotReadyAddresses indicates that any agent which deals with endpoints for this
	// Service should disregard any indications of ready/not-ready.
	// The primary use case for setting this field is for a StatefulSet's Headless Service to
	// propagate SRV DNS records for its Pods for the purpose of peer discovery.
	// The Kubernetes controllers that generate Endpoints and EndpointSlice resources for
	// Services interpret this to mean that all endpoints are considered "ready" even if the
	// Pods themselves are not. Agents which consume only Kubernetes generated endpoints
	// through the Endpoints or EndpointSlice resources can safely assume this behavior.
	// +optional
	PublishNotReadyAddresses bool `json:"publishNotReadyAddresses,omitempty" protobuf:"varint,13,opt,name=publishNotReadyAddresses"`

	// sessionAffinityConfig contains the configurations of session affinity.
	// +optional
	SessionAffinityConfig *SessionAffinityConfig `json:"sessionAffinityConfig,omitempty" protobuf:"bytes,14,opt,name=sessionAffinityConfig"`

	// TopologyKeys is tombstoned to show why 16 is reserved protobuf tag.
	// TopologyKeys []string `json:"topologyKeys,omitempty" protobuf:"bytes,16,opt,name=topologyKeys"`

	// IPFamily is tombstoned to show why 15 is a reserved protobuf tag.
	// IPFamily *IPFamily `json:"ipFamily,omitempty" protobuf:"bytes,15,opt,name=ipFamily,Configcasttype=IPFamily"`

	// IPFamilies is a list of IP families (e.g. IPv4, IPv6) assigned to this
	// service. This field is usually assigned automatically based on cluster
	// configuration and the ipFamilyPolicy field. If this field is specified
	// manually, the requested family is available in the cluster,
	// and ipFamilyPolicy allows it, it will be used; otherwise creation of
	// the service will fail. This field is conditionally mutable: it allows
	// for adding or removing a secondary IP family, but it does not allow
	// changing the primary IP family of the Service. Valid values are "IPv4"
	// and "IPv6".  This field only applies to Services of types ClusterIP,
	// NodePort, and LoadBalancer, and does apply to "headless" services.
	// This field will be wiped when updating a Service to type ExternalName.
	//
	// This field may hold a maximum of two entries (dual-stack families, in
	// either order).  These families must correspond to the values of the
	// clusterIPs field, if specified. Both clusterIPs and ipFamilies are
	// governed by the ipFamilyPolicy field.
	// +listType=atomic
	// +optional
	IPFamilies []IPFamily `json:"ipFamilies,omitempty" protobuf:"bytes,19,opt,name=ipFamilies,casttype=IPFamily"`

	// IPFamilyPolicy represents the dual-stack-ness requested or required by
	// this Service. If there is no value provided, then this field will be set
	// to SingleStack. Services can be "SingleStack" (a single IP family),
	// "PreferDualStack" (two IP families on dual-stack configured clusters or
	// a single IP family on single-stack clusters), or "RequireDualStack"
	// (two IP families on dual-stack configured clusters, otherwise fail). The
	// ipFamilies and clusterIPs fields depend on the value of this field. This
	// field will be wiped when updating a service to type ExternalName.
	// +optional
	IPFamilyPolicy *IPFamilyPolicy `json:"ipFamilyPolicy,omitempty" protobuf:"bytes,17,opt,name=ipFamilyPolicy,casttype=IPFamilyPolicy"`

	// allocateLoadBalancerNodePorts defines if NodePorts will be automatically
	// allocated for services with type LoadBalancer.  Default is "true". It
	// may be set to "false" if the cluster load-balancer does not rely on
	// NodePorts.  If the caller requests specific NodePorts (by specifying a
	// value), those requests will be respected, regardless of this field.
	// This field may only be set for services with type LoadBalancer and will
	// be cleared if the type is changed to any other type.
	// +optional
	AllocateLoadBalancerNodePorts *bool `json:"allocateLoadBalancerNodePorts,omitempty" protobuf:"bytes,20,opt,name=allocateLoadBalancerNodePorts"`

	// loadBalancerClass is the class of the load balancer implementation this Service belongs to.
	// If specified, the value of this field must be a label-style identifier, with an optional prefix,
	// e.g. "internal-vip" or "example.com/internal-vip". Unprefixed names are reserved for end-users.
	// This field can only be set when the Service type is 'LoadBalancer'. If not set, the default load
	// balancer implementation is used, today this is typically done through the cloud provider integration,
	// but should apply for any default implementation. If set, it is assumed that a load balancer
	// implementation is watching for Services with a matching class. Any default load balancer
	// implementation (e.g. cloud providers) should ignore Services that set this field.
	// This field can only be set when creating or updating a Service to type 'LoadBalancer'.
	// Once set, it can not be changed. This field will be wiped when a service is updated to a non 'LoadBalancer' type.
	// +optional
	LoadBalancerClass *string `json:"loadBalancerClass,omitempty" protobuf:"bytes,21,opt,name=loadBalancerClass"`

	// InternalTrafficPolicy describes how nodes distribute service traffic they
	// receive on the ClusterIP. If set to "Local", the proxy will assume that pods
	// only want to talk to endpoints of the service on the same node as the pod,
	// dropping the traffic if there are no local endpoints. The default value,
	// "Cluster", uses the standard behavior of routing to all endpoints evenly
	// (possibly modified by topology and other features).
	// +optional
	InternalTrafficPolicy *ServiceInternalTrafficPolicy `json:"internalTrafficPolicy,omitempty" protobuf:"bytes,22,opt,name=internalTrafficPolicy"`
}

// ServicePort contains information on service's port.
type ServicePort struct {
	// The name of this port within the service. This must be a DNS_LABEL.
	// All ports within a ServiceSpec must have unique names. When considering
	// the endpoints for a Service, this must match the 'name' field in the
	// EndpointPort.
	// Optional if only one ServicePort is defined on this service.
	// +optional
	Name string `json:"name,omitempty" protobuf:"bytes,1,opt,name=name"`

	// The IP protocol for this port. Supports "TCP", "UDP", and "SCTP".
	// Default is TCP.
	// +default="TCP"
	// +optional
	Protocol Protocol `json:"protocol,omitempty" protobuf:"bytes,2,opt,name=protocol,casttype=Protocol"`

	// The application protocol for this port.
	// This field follows standard Kubernetes label syntax.
	// Un-prefixed names are reserved for IANA standard service names (as per
	// RFC-6335 and https://www.iana.org/assignments/service-names).
	// Non-standard protocols should use prefixed names such as
	// mycompany.com/my-custom-protocol.
	// +optional
	AppProtocol *string `json:"appProtocol,omitempty" protobuf:"bytes,6,opt,name=appProtocol"`

	// The port that will be exposed by this service.
	Port int32 `json:"port" protobuf:"varint,3,opt,name=port"`

	// Number or name of the port to access on the pods targeted by the service.
	// Number must be in the range 1 to 65535. Name must be an IANA_SVC_NAME.
	// If this is a string, it will be looked up as a named port in the
	// target Pod's container ports. If this is not specified, the value
	// of the 'port' field is used (an identity map).
	// This field is ignored for services with clusterIP=None, and should be
	// omitted or set equal to the 'port' field.
	// More info: https://kubernetes.io/docs/concepts/services-networking/service/#defining-a-service
	// +optional
	TargetPort intstr.IntOrString `json:"targetPort,omitempty" protobuf:"bytes,4,opt,name=targetPort"`

	// The port on each node on which this service is exposed when type is
	// NodePort or LoadBalancer.  Usually assigned by the system. If a value is
	// specified, in-range, and not in use it will be used, otherwise the
	// operation will fail.  If not specified, a port will be allocated if this
	// Service requires one.  If this field is specified when creating a
	// Service which does not need it, creation will fail. This field will be
	// wiped when updating a Service to no longer need it (e.g. changing type
	// from NodePort to ClusterIP).
	// More info: https://kubernetes.io/docs/concepts/services-networking/service/#type-nodeport
	// +optional
	NodePort int32 `json:"nodePort,omitempty" protobuf:"varint,5,opt,name=nodePort"`
}

// +genclient
// +genclient:skipVerbs=deleteCollection
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Service is a named abstraction of software service (for example, mysql) consisting of local port
// (for example 3306) that the proxy listens on, and the selector that determines which pods
// will answer requests sent through the proxy.
type Service struct {
	metav1.TypeMeta `json:",inline"`
	// Standard object's metadata.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// Spec defines the behavior of a service.
	// https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#spec-and-status
	// +optional
	Spec ServiceSpec `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`

	// Most recently observed status of the service.
	// Populated by the system.
	// Read-only.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#spec-and-status
	// +optional
	Status ServiceStatus `json:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}

const (
	// ClusterIPNone - do not assign a cluster IP
	// no proxying required and no environment variables should be created for pods
	ClusterIPNone = "None"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ServiceList holds a list of services.
type ServiceList struct {
	metav1.TypeMeta `json:",inline"`
	// Standard list metadata.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds
	// +optional
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// List of services
	Items []Service `json:"items" protobuf:"bytes,2,rep,name=items"`
}

// +genclient
// +genclient:method=CreateToken,verb=create,subresource=token,input=k8s.io/api/authentication/v1.TokenRequest,result=k8s.io/api/authentication/v1.TokenRequest
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ServiceAccount binds together:
// * a name, understood by users, and perhaps by peripheral systems, for an identity
// * a principal that can be authenticated and authorized
// * a set of secrets
type ServiceAccount struct {
	metav1.TypeMeta `json:",inline"`
	// Standard object's metadata.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// Secrets is a list of the secrets in the same namespace that pods running using this ServiceAccount are allowed to use.
	// Pods are only limited to this list if this service account has a "kubernetes.io/enforce-mountable-secrets" annotation set to "true".
	// This field should not be used to find auto-generated service account token secrets for use outside of pods.
	// Instead, tokens can be requested directly using the TokenRequest API, or service account token secrets can be manually created.
	// More info: https://kubernetes.io/docs/concepts/configuration/secret
	// +optional
	// +patchMergeKey=name
	// +patchStrategy=merge
	Secrets []ObjectReference `json:"secrets,omitempty" patchStrategy:"merge" patchMergeKey:"name" protobuf:"bytes,2,rep,name=secrets"`

	// ImagePullSecrets is a list of references to secrets in the same namespace to use for pulling any images
	// in pods that reference this ServiceAccount. ImagePullSecrets are distinct from Secrets because Secrets
	// can be mounted in the pod, but ImagePullSecrets are only accessed by the kubelet.
	// More info: https://kubernetes.io/docs/concepts/containers/images/#specifying-imagepullsecrets-on-a-pod
	// +optional
	ImagePullSecrets []LocalObjectReference `json:"imagePullSecrets,omitempty" protobuf:"bytes,3,rep,name=imagePullSecrets"`

	// AutomountServiceAccountToken indicates whether pods running as this service account should have an API token automatically mounted.
	// Can be overridden at the pod level.
	// +optional
	AutomountServiceAccountToken *bool `json:"automountServiceAccountToken,omitempty" protobuf:"varint,4,opt,name=automountServiceAccountToken"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ServiceAccountList is a list of ServiceAccount objects
type ServiceAccountList struct {
	metav1.TypeMeta `json:",inline"`
	// Standard list metadata.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds
	// +optional
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// List of ServiceAccounts.
	// More info: https://kubernetes.io/docs/tasks/configure-pod-container/configure-service-account/
	Items []ServiceAccount `json:"items" protobuf:"bytes,2,rep,name=items"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Endpoints is a collection of endpoints that implement the actual service. Example:
//
//	 Name: "mysvc",
//	 Subsets: [
//	   {
//	     Addresses: [{"ip": "10.10.1.1"}, {"ip": "10.10.2.2"}],
//	     Ports: [{"name": "a", "port": 8675}, {"name": "b", "port": 309}]
//	   },
//	   {
//	     Addresses: [{"ip": "10.10.3.3"}],
//	     Ports: [{"name": "a", "port": 93}, {"name": "b", "port": 76}]
//	   },
//	]
type Endpoints struct {
	metav1.TypeMeta `json:",inline"`
	// Standard object's metadata.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// The set of all endpoints is the union of all subsets. Addresses are placed into
	// subsets according to the IPs they share. A single address with multiple ports,
	// some of which are ready and some of which are not (because they come from
	// different containers) will result in the address being displayed in different
	// subsets for the different ports. No address will appear in both Addresses and
	// NotReadyAddresses in the same subset.
	// Sets of addresses and ports that comprise a service.
	// +optional
	Subsets []EndpointSubset `json:"subsets,omitempty" protobuf:"bytes,2,rep,name=subsets"`
}

// EndpointSubset is a group of addresses with a common set of ports. The
// expanded set of endpoints is the Cartesian product of Addresses x Ports.
// For example, given:
//
//	{
//	  Addresses: [{"ip": "10.10.1.1"}, {"ip": "10.10.2.2"}],
//	  Ports:     [{"name": "a", "port": 8675}, {"name": "b", "port": 309}]
//	}
//
// The resulting set of endpoints can be viewed as:
//
//	a: [ 10.10.1.1:8675, 10.10.2.2:8675 ],
//	b: [ 10.10.1.1:309, 10.10.2.2:309 ]
type EndpointSubset struct {
	// IP addresses which offer the related ports that are marked as ready. These endpoints
	// should be considered safe for load balancers and clients to utilize.
	// +optional
	Addresses []EndpointAddress `json:"addresses,omitempty" protobuf:"bytes,1,rep,name=addresses"`
	// IP addresses which offer the related ports but are not currently marked as ready
	// because they have not yet finished starting, have recently failed a readiness check,
	// or have recently failed a liveness check.
	// +optional
	NotReadyAddresses []EndpointAddress `json:"notReadyAddresses,omitempty" protobuf:"bytes,2,rep,name=notReadyAddresses"`
	// Port numbers available on the related IP addresses.
	// +optional
	Ports []EndpointPort `json:"ports,omitempty" protobuf:"bytes,3,rep,name=ports"`
}

// EndpointAddress is a tuple that describes single IP address.
// +structType=atomic
type EndpointAddress struct {
	// The IP of this endpoint.
	// May not be loopback (127.0.0.0/8 or ::1), link-local (169.254.0.0/16 or fe80::/10),
	// or link-local multicast (224.0.0.0/24 or ff02::/16).
	IP string `json:"ip" protobuf:"bytes,1,opt,name=ip"`
	// The Hostname of this endpoint
	// +optional
	Hostname string `json:"hostname,omitempty" protobuf:"bytes,3,opt,name=hostname"`
	// Optional: Node hosting this endpoint. This can be used to determine endpoints local to a node.
	// +optional
	NodeName *string `json:"nodeName,omitempty" protobuf:"bytes,4,opt,name=nodeName"`
	// Reference to object providing the endpoint.
	// +optional
	TargetRef *ObjectReference `json:"targetRef,omitempty" protobuf:"bytes,2,opt,name=targetRef"`
}

// EndpointPort is a tuple that describes a single port.
// +structType=atomic
type EndpointPort struct {
	// The name of this port.  This must match the 'name' field in the
	// corresponding ServicePort.
	// Must be a DNS_LABEL.
	// Optional only if one port is defined.
	// +optional
	Name string `json:"name,omitempty" protobuf:"bytes,1,opt,name=name"`

	// The port number of the endpoint.
	Port int32 `json:"port" protobuf:"varint,2,opt,name=port"`

	// The IP protocol for this port.
	// Must be UDP, TCP, or SCTP.
	// Default is TCP.
	// +optional
	Protocol Protocol `json:"protocol,omitempty" protobuf:"bytes,3,opt,name=protocol,casttype=Protocol"`

	// The application protocol for this port.
	// This field follows standard Kubernetes label syntax.
	// Un-prefixed names are reserved for IANA standard service names (as per
	// RFC-6335 and https://www.iana.org/assignments/service-names).
	// Non-standard protocols should use prefixed names such as
	// mycompany.com/my-custom-protocol.
	// +optional
	AppProtocol *string `json:"appProtocol,omitempty" protobuf:"bytes,4,opt,name=appProtocol"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// EndpointsList is a list of endpoints.
type EndpointsList struct {
	metav1.TypeMeta `json:",inline"`
	// Standard list metadata.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds
	// +optional
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// List of endpoints.
	Items []Endpoints `json:"items" protobuf:"bytes,2,rep,name=items"`
}

// NodeSpec describes the attributes that a node is created with.
type NodeSpec struct {
	// PodCIDR represents the pod IP range assigned to the node.
	// +optional
	PodCIDR string `json:"podCIDR,omitempty" protobuf:"bytes,1,opt,name=podCIDR"`

	// podCIDRs represents the IP ranges assigned to the node for usage by Pods on that node. If this
	// field is specified, the 0th entry must match the podCIDR field. It may contain at most 1 value for
	// each of IPv4 and IPv6.
	// +optional
	// +patchStrategy=merge
	PodCIDRs []string `json:"podCIDRs,omitempty" protobuf:"bytes,7,opt,name=podCIDRs" patchStrategy:"merge"`

	// ID of the node assigned by the cloud provider in the format: <ProviderName>://<ProviderSpecificNodeID>
	// +optional
	ProviderID string `json:"providerID,omitempty" protobuf:"bytes,3,opt,name=providerID"`
	// Unschedulable controls node schedulability of new pods. By default, node is schedulable.
	// More info: https://kubernetes.io/docs/concepts/nodes/node/#manual-node-administration
	// +optional
	Unschedulable bool `json:"unschedulable,omitempty" protobuf:"varint,4,opt,name=unschedulable"`
	// If specified, the node's taints.
	// +optional
	Taints []Taint `json:"taints,omitempty" protobuf:"bytes,5,opt,name=taints"`

	// Deprecated: Previously used to specify the source of the node's configuration for the DynamicKubeletConfig feature. This feature is removed.
	// +optional
	ConfigSource *NodeConfigSource `json:"configSource,omitempty" protobuf:"bytes,6,opt,name=configSource"`

	// Deprecated. Not all kubelets will set this field. Remove field after 1.13.
	// see: https://issues.k8s.io/61966
	// +optional
	DoNotUseExternalID string `json:"externalID,omitempty" protobuf:"bytes,2,opt,name=externalID"`
}

// NodeConfigSource specifies a source of node configuration. Exactly one subfield (excluding metadata) must be non-nil.
// This API is deprecated since 1.22
type NodeConfigSource struct {
	// For historical context, regarding the below kind, apiVersion, and configMapRef deprecation tags:
	// 1. kind/apiVersion were used by the kubelet to persist this struct to disk (they had no protobuf tags)
	// 2. configMapRef and proto tag 1 were used by the API to refer to a configmap,
	//    but used a generic ObjectReference type that didn't really have the fields we needed
	// All uses/persistence of the NodeConfigSource struct prior to 1.11 were gated by alpha feature flags,
	// so there was no persisted data for these fields that needed to be migrated/handled.

	// +k8s:deprecated=kind
	// +k8s:deprecated=apiVersion
	// +k8s:deprecated=configMapRef,protobuf=1

	// ConfigMap is a reference to a Node's ConfigMap
	ConfigMap *ConfigMapNodeConfigSource `json:"configMap,omitempty" protobuf:"bytes,2,opt,name=configMap"`
}

// ConfigMapNodeConfigSource contains the information to reference a ConfigMap as a config source for the Node.
// This API is deprecated since 1.22: https://git.k8s.io/enhancements/keps/sig-node/281-dynamic-kubelet-configuration
type ConfigMapNodeConfigSource struct {
	// Namespace is the metadata.namespace of the referenced ConfigMap.
	// This field is required in all cases.
	Namespace string `json:"namespace" protobuf:"bytes,1,opt,name=namespace"`

	// Name is the metadata.name of the referenced ConfigMap.
	// This field is required in all cases.
	Name string `json:"name" protobuf:"bytes,2,opt,name=name"`

	// UID is the metadata.UID of the referenced ConfigMap.
	// This field is forbidden in Node.Spec, and required in Node.Status.
	// +optional
	UID types.UID `json:"uid,omitempty" protobuf:"bytes,3,opt,name=uid"`

	// ResourceVersion is the metadata.ResourceVersion of the referenced ConfigMap.
	// This field is forbidden in Node.Spec, and required in Node.Status.
	// +optional
	ResourceVersion string `json:"resourceVersion,omitempty" protobuf:"bytes,4,opt,name=resourceVersion"`

	// KubeletConfigKey declares which key of the referenced ConfigMap corresponds to the KubeletConfiguration structure
	// This field is required in all cases.
	KubeletConfigKey string `json:"kubeletConfigKey" protobuf:"bytes,5,opt,name=kubeletConfigKey"`
}

// DaemonEndpoint contains information about a single Daemon endpoint.
type DaemonEndpoint struct {
	/*
		The port tag was not properly in quotes in earlier releases, so it must be
		uppercased for backwards compat (since it was falling back to var name of
		'Port').
	*/

	// Port number of the given endpoint.
	Port int32 `json:"Port" protobuf:"varint,1,opt,name=Port"`
}

// NodeDaemonEndpoints lists ports opened by daemons running on the Node.
type NodeDaemonEndpoints struct {
	// Endpoint on which Kubelet is listening.
	// +optional
	KubeletEndpoint DaemonEndpoint `json:"kubeletEndpoint,omitempty" protobuf:"bytes,1,opt,name=kubeletEndpoint"`
}

// NodeSystemInfo is a set of ids/uuids to uniquely identify the node.
type NodeSystemInfo struct {
	// MachineID reported by the node. For unique machine identification
	// in the cluster this field is preferred. Learn more from man(5)
	// machine-id: http://man7.org/linux/man-pages/man5/machine-id.5.html
	MachineID string `json:"machineID" protobuf:"bytes,1,opt,name=machineID"`
	// SystemUUID reported by the node. For unique machine identification
	// MachineID is preferred. This field is specific to Red Hat hosts
	// https://access.redhat.com/documentation/en-us/red_hat_subscription_management/1/html/rhsm/uuid
	SystemUUID string `json:"systemUUID" protobuf:"bytes,2,opt,name=systemUUID"`
	// Boot ID reported by the node.
	BootID string `json:"bootID" protobuf:"bytes,3,opt,name=bootID"`
	// Kernel Version reported by the node from 'uname -r' (e.g. 3.16.0-0.bpo.4-amd64).
	KernelVersion string `json:"kernelVersion" protobuf:"bytes,4,opt,name=kernelVersion"`
	// OS Image reported by the node from /etc/os-release (e.g. Debian GNU/Linux 7 (wheezy)).
	OSImage string `json:"osImage" protobuf:"bytes,5,opt,name=osImage"`
	// ContainerRuntime Version reported by the node through runtime remote API (e.g. containerd://1.4.2).
	ContainerRuntimeVersion string `json:"containerRuntimeVersion" protobuf:"bytes,6,opt,name=containerRuntimeVersion"`
	// Kubelet Version reported by the node.
	KubeletVersion string `json:"kubeletVersion" protobuf:"bytes,7,opt,name=kubeletVersion"`
	// KubeProxy Version reported by the node.
	KubeProxyVersion string `json:"kubeProxyVersion" protobuf:"bytes,8,opt,name=kubeProxyVersion"`
	// The Operating System reported by the node
	OperatingSystem string `json:"operatingSystem" protobuf:"bytes,9,opt,name=operatingSystem"`
	// The Architecture reported by the node
	Architecture string `json:"architecture" protobuf:"bytes,10,opt,name=architecture"`
}

// NodeConfigStatus describes the status of the config assigned by Node.Spec.ConfigSource.
type NodeConfigStatus struct {
	// Assigned reports the checkpointed config the node will try to use.
	// When Node.Spec.ConfigSource is updated, the node checkpoints the associated
	// config payload to local disk, along with a record indicating intended
	// config. The node refers to this record to choose its config checkpoint, and
	// reports this record in Assigned. Assigned only updates in the status after
	// the record has been checkpointed to disk. When the Kubelet is restarted,
	// it tries to make the Assigned config the Active config by loading and
	// validating the checkpointed payload identified by Assigned.
	// +optional
	Assigned *NodeConfigSource `json:"assigned,omitempty" protobuf:"bytes,1,opt,name=assigned"`
	// Active reports the checkpointed config the node is actively using.
	// Active will represent either the current version of the Assigned config,
	// or the current LastKnownGood config, depending on whether attempting to use the
	// Assigned config results in an error.
	// +optional
	Active *NodeConfigSource `json:"active,omitempty" protobuf:"bytes,2,opt,name=active"`
	// LastKnownGood reports the checkpointed config the node will fall back to
	// when it encounters an error attempting to use the Assigned config.
	// The Assigned config becomes the LastKnownGood config when the node determines
	// that the Assigned config is stable and correct.
	// This is currently implemented as a 10-minute soak period starting when the local
	// record of Assigned config is updated. If the Assigned config is Active at the end
	// of this period, it becomes the LastKnownGood. Note that if Spec.ConfigSource is
	// reset to nil (use local defaults), the LastKnownGood is also immediately reset to nil,
	// because the local default config is always assumed good.
	// You should not make assumptions about the node's method of determining config stability
	// and correctness, as this may change or become configurable in the future.
	// +optional
	LastKnownGood *NodeConfigSource `json:"lastKnownGood,omitempty" protobuf:"bytes,3,opt,name=lastKnownGood"`
	// Error describes any problems reconciling the Spec.ConfigSource to the Active config.
	// Errors may occur, for example, attempting to checkpoint Spec.ConfigSource to the local Assigned
	// record, attempting to checkpoint the payload associated with Spec.ConfigSource, attempting
	// to load or validate the Assigned config, etc.
	// Errors may occur at different points while syncing config. Earlier errors (e.g. download or
	// checkpointing errors) will not result in a rollback to LastKnownGood, and may resolve across
	// Kubelet retries. Later errors (e.g. loading or validating a checkpointed config) will result in
	// a rollback to LastKnownGood. In the latter case, it is usually possible to resolve the error
	// by fixing the config assigned in Spec.ConfigSource.
	// You can find additional information for debugging by searching the error message in the Kubelet log.
	// Error is a human-readable description of the error state; machines can check whether or not Error
	// is empty, but should not rely on the stability of the Error text across Kubelet versions.
	// +optional
	Error string `json:"error,omitempty" protobuf:"bytes,4,opt,name=error"`
}

// NodeStatus is information about the current status of a node.
type NodeStatus struct {
	// Capacity represents the total resources of a node.
	// More info: https://kubernetes.io/docs/concepts/storage/persistent-volumes#capacity
	// +optional
	Capacity ResourceList `json:"capacity,omitempty" protobuf:"bytes,1,rep,name=capacity,casttype=ResourceList,castkey=ResourceName"`
	// Allocatable represents the resources of a node that are available for scheduling.
	// Defaults to Capacity.
	// +optional
	Allocatable ResourceList `json:"allocatable,omitempty" protobuf:"bytes,2,rep,name=allocatable,casttype=ResourceList,castkey=ResourceName"`
	// NodePhase is the recently observed lifecycle phase of the node.
	// More info: https://kubernetes.io/docs/concepts/nodes/node/#phase
	// The field is never populated, and now is deprecated.
	// +optional
	Phase NodePhase `json:"phase,omitempty" protobuf:"bytes,3,opt,name=phase,casttype=NodePhase"`
	// Conditions is an array of current observed node conditions.
	// More info: https://kubernetes.io/docs/concepts/nodes/node/#condition
	// +optional
	// +patchMergeKey=type
	// +patchStrategy=merge
	Conditions []NodeCondition `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type" protobuf:"bytes,4,rep,name=conditions"`
	// List of addresses reachable to the node.
	// Queried from cloud provider, if available.
	// More info: https://kubernetes.io/docs/concepts/nodes/node/#addresses
	// Note: This field is declared as mergeable, but the merge key is not sufficiently
	// unique, which can cause data corruption when it is merged. Callers should instead
	// use a full-replacement patch. See https://pr.k8s.io/79391 for an example.
	// Consumers should assume that addresses can change during the
	// lifetime of a Node. However, there are some exceptions where this may not
	// be possible, such as Pods that inherit a Node's address in its own status or
	// consumers of the downward API (status.hostIP).
	// +optional
	// +patchMergeKey=type
	// +patchStrategy=merge
	Addresses []NodeAddress `json:"addresses,omitempty" patchStrategy:"merge" patchMergeKey:"type" protobuf:"bytes,5,rep,name=addresses"`
	// Endpoints of daemons running on the Node.
	// +optional
	DaemonEndpoints NodeDaemonEndpoints `json:"daemonEndpoints,omitempty" protobuf:"bytes,6,opt,name=daemonEndpoints"`
	// Set of ids/uuids to uniquely identify the node.
	// More info: https://kubernetes.io/docs/concepts/nodes/node/#info
	// +optional
	NodeInfo NodeSystemInfo `json:"nodeInfo,omitempty" protobuf:"bytes,7,opt,name=nodeInfo"`
	// List of container images on this node
	// +optional
	Images []ContainerImage `json:"images,omitempty" protobuf:"bytes,8,rep,name=images"`
	// List of attachable volumes in use (mounted) by the node.
	// +optional
	VolumesInUse []UniqueVolumeName `json:"volumesInUse,omitempty" protobuf:"bytes,9,rep,name=volumesInUse"`
	// List of volumes that are attached to the node.
	// +optional
	VolumesAttached []AttachedVolume `json:"volumesAttached,omitempty" protobuf:"bytes,10,rep,name=volumesAttached"`
	// Status of the config assigned to the node via the dynamic Kubelet config feature.
	// +optional
	Config *NodeConfigStatus `json:"config,omitempty" protobuf:"bytes,11,opt,name=config"`
}

type UniqueVolumeName string

// AttachedVolume describes a volume attached to a node
type AttachedVolume struct {
	// Name of the attached volume
	Name UniqueVolumeName `json:"name" protobuf:"bytes,1,rep,name=name"`

	// DevicePath represents the device path where the volume should be available
	DevicePath string `json:"devicePath" protobuf:"bytes,2,rep,name=devicePath"`
}

// AvoidPods describes pods that should avoid this node. This is the value for a
// Node annotation with key scheduler.alpha.kubernetes.io/preferAvoidPods and
// will eventually become a field of NodeStatus.
type AvoidPods struct {
	// Bounded-sized list of signatures of pods that should avoid this node, sorted
	// in timestamp order from oldest to newest. Size of the slice is unspecified.
	// +optional
	PreferAvoidPods []PreferAvoidPodsEntry `json:"preferAvoidPods,omitempty" protobuf:"bytes,1,rep,name=preferAvoidPods"`
}

// Describes a class of pods that should avoid this node.
type PreferAvoidPodsEntry struct {
	// The class of pods.
	PodSignature PodSignature `json:"podSignature" protobuf:"bytes,1,opt,name=podSignature"`
	// Time at which this entry was added to the list.
	// +optional
	EvictionTime metav1.Time `json:"evictionTime,omitempty" protobuf:"bytes,2,opt,name=evictionTime"`
	// (brief) reason why this entry was added to the list.
	// +optional
	Reason string `json:"reason,omitempty" protobuf:"bytes,3,opt,name=reason"`
	// Human readable message indicating why this entry was added to the list.
	// +optional
	Message string `json:"message,omitempty" protobuf:"bytes,4,opt,name=message"`
}

// Describes the class of pods that should avoid this node.
// Exactly one field should be set.
type PodSignature struct {
	// Reference to controller whose pods should avoid this node.
	// +optional
	PodController *metav1.OwnerReference `json:"podController,omitempty" protobuf:"bytes,1,opt,name=podController"`
}

// Describe a container image
type ContainerImage struct {
	// Names by which this image is known.
	// e.g. ["kubernetes.example/hyperkube:v1.0.7", "cloud-vendor.registry.example/cloud-vendor/hyperkube:v1.0.7"]
	// +optional
	Names []string `json:"names" protobuf:"bytes,1,rep,name=names"`
	// The size of the image in bytes.
	// +optional
	SizeBytes int64 `json:"sizeBytes,omitempty" protobuf:"varint,2,opt,name=sizeBytes"`
}

// +enum
type NodePhase string

// These are the valid phases of node.
const (
	// NodePending means the node has been created/added by the system, but not configured.
	NodePending NodePhase = "Pending"
	// NodeRunning means the node has been configured and has Kubernetes components running.
	NodeRunning NodePhase = "Running"
	// NodeTerminated means the node has been removed from the cluster.
	NodeTerminated NodePhase = "Terminated"
)

type NodeConditionType string

// These are valid but not exhaustive conditions of node. A cloud provider may set a condition not listed here.
// The built-in set of conditions are:
// NodeReachable, NodeLive, NodeReady, NodeSchedulable, NodeRunnable.
const (
	// NodeReady means kubelet is healthy and ready to accept pods.
	NodeReady NodeConditionType = "Ready"
	// NodeMemoryPressure means the kubelet is under pressure due to insufficient available memory.
	NodeMemoryPressure NodeConditionType = "MemoryPressure"
	// NodeDiskPressure means the kubelet is under pressure due to insufficient available disk.
	NodeDiskPressure NodeConditionType = "DiskPressure"
	// NodePIDPressure means the kubelet is under pressure due to insufficient available PID.
	NodePIDPressure NodeConditionType = "PIDPressure"
	// NodeNetworkUnavailable means that network for the node is not correctly configured.
	NodeNetworkUnavailable NodeConditionType = "NetworkUnavailable"
)

// NodeCondition contains condition information for a node.
type NodeCondition struct {
	// Type of node condition.
	Type NodeConditionType `json:"type" protobuf:"bytes,1,opt,name=type,casttype=NodeConditionType"`
	// Status of the condition, one of True, False, Unknown.
	Status ConditionStatus `json:"status" protobuf:"bytes,2,opt,name=status,casttype=ConditionStatus"`
	// Last time we got an update on a given condition.
	// +optional
	LastHeartbeatTime metav1.Time `json:"lastHeartbeatTime,omitempty" protobuf:"bytes,3,opt,name=lastHeartbeatTime"`
	// Last time the condition transit from one status to another.
	// +optional
	LastTransitionTime metav1.Time `json:"lastTransitionTime,omitempty" protobuf:"bytes,4,opt,name=lastTransitionTime"`
	// (brief) reason for the condition's last transition.
	// +optional
	Reason string `json:"reason,omitempty" protobuf:"bytes,5,opt,name=reason"`
	// Human readable message indicating details about last transition.
	// +optional
	Message string `json:"message,omitempty" protobuf:"bytes,6,opt,name=message"`
}

type NodeAddressType string

// These are built-in addresses type of node. A cloud provider may set a type not listed here.
const (
	// NodeHostName identifies a name of the node. Although every node can be assumed
	// to have a NodeAddress of this type, its exact syntax and semantics are not
	// defined, and are not consistent between different clusters.
	NodeHostName NodeAddressType = "Hostname"

	// NodeInternalIP identifies an IP address which is assigned to one of the node's
	// network interfaces. Every node should have at least one address of this type.
	//
	// An internal IP is normally expected to be reachable from every other node, but
	// may not be visible to hosts outside the cluster. By default it is assumed that
	// kube-apiserver can reach node internal IPs, though it is possible to configure
	// clusters where this is not the case.
	//
	// NodeInternalIP is the default type of node IP, and does not necessarily imply
	// that the IP is ONLY reachable internally. If a node has multiple internal IPs,
	// no specific semantics are assigned to the additional IPs.
	NodeInternalIP NodeAddressType = "InternalIP"

	// NodeExternalIP identifies an IP address which is, in some way, intended to be
	// more usable from outside the cluster then an internal IP, though no specific
	// semantics are defined. It may be a globally routable IP, though it is not
	// required to be.
	//
	// External IPs may be assigned directly to an interface on the node, like a
	// NodeInternalIP, or alternatively, packets sent to the external IP may be NAT'ed
	// to an internal node IP rather than being delivered directly (making the IP less
	// efficient for node-to-node traffic than a NodeInternalIP).
	NodeExternalIP NodeAddressType = "ExternalIP"

	// NodeInternalDNS identifies a DNS name which resolves to an IP address which has
	// the characteristics of a NodeInternalIP. The IP it resolves to may or may not
	// be a listed NodeInternalIP address.
	NodeInternalDNS NodeAddressType = "InternalDNS"

	// NodeExternalDNS identifies a DNS name which resolves to an IP address which has
	// the characteristics of a NodeExternalIP. The IP it resolves to may or may not
	// be a listed NodeExternalIP address.
	NodeExternalDNS NodeAddressType = "ExternalDNS"
)

// NodeAddress contains information for the node's address.
type NodeAddress struct {
	// Node address type, one of Hostname, ExternalIP or InternalIP.
	Type NodeAddressType `json:"type" protobuf:"bytes,1,opt,name=type,casttype=NodeAddressType"`
	// The node address.
	Address string `json:"address" protobuf:"bytes,2,opt,name=address"`
}

// ResourceName is the name identifying various resources in a ResourceList.
type ResourceName string

// Resource names must be not more than 63 characters, consisting of upper- or lower-case alphanumeric characters,
// with the -, _, and . characters allowed anywhere, except the first or last character.
// The default convention, matching that for annotations, is to use lower-case names, with dashes, rather than
// camel case, separating compound words.
// Fully-qualified resource typenames are constructed from a DNS-style subdomain, followed by a slash `/` and a name.
const (
	// CPU, in cores. (500m = .5 cores)
	ResourceCPU ResourceName = "cpu"
	// Memory, in bytes. (500Gi = 500GiB = 500 * 1024 * 1024 * 1024)
	ResourceMemory ResourceName = "memory"
	// Volume size, in bytes (e,g. 5Gi = 5GiB = 5 * 1024 * 1024 * 1024)
	ResourceStorage ResourceName = "storage"
	// Local ephemeral storage, in bytes. (500Gi = 500GiB = 500 * 1024 * 1024 * 1024)
	// The resource name for ResourceEphemeralStorage is alpha and it can change across releases.
	ResourceEphemeralStorage ResourceName = "ephemeral-storage"
)

const (
	// Default namespace prefix.
	ResourceDefaultNamespacePrefix = "kubernetes.io/"
	// Name prefix for huge page resources (alpha).
	ResourceHugePagesPrefix = "hugepages-"
	// Name prefix for storage resource limits
	ResourceAttachableVolumesPrefix = "attachable-volumes-"
)

// ResourceList is a set of (resource name, quantity) pairs.
// ResourceList 是一组（资源名称，数量）对。
type ResourceList map[ResourceName]resource.Quantity

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Node is a worker node in Kubernetes.
// Each node will have a unique identifier in the cache (i.e. in etcd).
type Node struct {
	metav1.TypeMeta `json:",inline"`
	// Standard object's metadata.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// Spec defines the behavior of a node.
	// https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#spec-and-status
	// +optional
	Spec NodeSpec `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`

	// Most recently observed status of the node.
	// Populated by the system.
	// Read-only.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#spec-and-status
	// +optional
	Status NodeStatus `json:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// NodeList is the whole list of all Nodes which have been registered with master.
type NodeList struct {
	metav1.TypeMeta `json:",inline"`
	// Standard list metadata.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds
	// +optional
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// List of nodes
	Items []Node `json:"items" protobuf:"bytes,2,rep,name=items"`
}

// FinalizerName is the name identifying a finalizer during namespace lifecycle.
type FinalizerName string

// These are internal finalizer values to Kubernetes, must be qualified name unless defined here or
// in metav1.
const (
	FinalizerKubernetes FinalizerName = "kubernetes"
)

// NamespaceSpec describes the attributes on a Namespace.
type NamespaceSpec struct {
	// Finalizers is an opaque list of values that must be empty to permanently remove object from storage.
	// More info: https://kubernetes.io/docs/tasks/administer-cluster/namespaces/
	// +optional
	Finalizers []FinalizerName `json:"finalizers,omitempty" protobuf:"bytes,1,rep,name=finalizers,casttype=FinalizerName"`
}

// NamespaceStatus is information about the current status of a Namespace.
type NamespaceStatus struct {
	// Phase is the current lifecycle phase of the namespace.
	// More info: https://kubernetes.io/docs/tasks/administer-cluster/namespaces/
	// +optional
	Phase NamespacePhase `json:"phase,omitempty" protobuf:"bytes,1,opt,name=phase,casttype=NamespacePhase"`

	// Represents the latest available observations of a namespace's current state.
	// +optional
	// +patchMergeKey=type
	// +patchStrategy=merge
	Conditions []NamespaceCondition `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type" protobuf:"bytes,2,rep,name=conditions"`
}

// +enum
type NamespacePhase string

// These are the valid phases of a namespace.
const (
	// NamespaceActive means the namespace is available for use in the system
	NamespaceActive NamespacePhase = "Active"
	// NamespaceTerminating means the namespace is undergoing graceful termination
	NamespaceTerminating NamespacePhase = "Terminating"
)

const (
	// NamespaceTerminatingCause is returned as a defaults.cause item when a change is
	// forbidden due to the namespace being terminated.
	NamespaceTerminatingCause metav1.CauseType = "NamespaceTerminating"
)

type NamespaceConditionType string

// These are built-in conditions of a namespace.
const (
	// NamespaceDeletionDiscoveryFailure contains information about namespace deleter errors during resource discovery.
	NamespaceDeletionDiscoveryFailure NamespaceConditionType = "NamespaceDeletionDiscoveryFailure"
	// NamespaceDeletionContentFailure contains information about namespace deleter errors during deletion of resources.
	NamespaceDeletionContentFailure NamespaceConditionType = "NamespaceDeletionContentFailure"
	// NamespaceDeletionGVParsingFailure contains information about namespace deleter errors parsing GV for legacy types.
	NamespaceDeletionGVParsingFailure NamespaceConditionType = "NamespaceDeletionGroupVersionParsingFailure"
	// NamespaceContentRemaining contains information about resources remaining in a namespace.
	NamespaceContentRemaining NamespaceConditionType = "NamespaceContentRemaining"
	// NamespaceFinalizersRemaining contains information about which finalizers are on resources remaining in a namespace.
	NamespaceFinalizersRemaining NamespaceConditionType = "NamespaceFinalizersRemaining"
)

// NamespaceCondition contains details about state of namespace.
type NamespaceCondition struct {
	// Type of namespace controller condition.
	Type NamespaceConditionType `json:"type" protobuf:"bytes,1,opt,name=type,casttype=NamespaceConditionType"`
	// Status of the condition, one of True, False, Unknown.
	Status ConditionStatus `json:"status" protobuf:"bytes,2,opt,name=status,casttype=ConditionStatus"`
	// +optional
	LastTransitionTime metav1.Time `json:"lastTransitionTime,omitempty" protobuf:"bytes,4,opt,name=lastTransitionTime"`
	// +optional
	Reason string `json:"reason,omitempty" protobuf:"bytes,5,opt,name=reason"`
	// +optional
	Message string `json:"message,omitempty" protobuf:"bytes,6,opt,name=message"`
}

// +genclient
// +genclient:nonNamespaced
// +genclient:skipVerbs=deleteCollection
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Namespace provides a scope for Names.
// Use of multiple namespaces is optional.
type Namespace struct {
	metav1.TypeMeta `json:",inline"`
	// Standard object's metadata.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// Spec defines the behavior of the Namespace.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#spec-and-status
	// +optional
	Spec NamespaceSpec `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`

	// Status describes the current status of a Namespace.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#spec-and-status
	// +optional
	Status NamespaceStatus `json:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// NamespaceList is a list of Namespaces.
type NamespaceList struct {
	metav1.TypeMeta `json:",inline"`
	// Standard list metadata.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds
	// +optional
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// Items is the list of Namespace objects in the list.
	// More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/namespaces/
	Items []Namespace `json:"items" protobuf:"bytes,2,rep,name=items"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Binding ties one object to another; for example, a pod is bound to a node by a scheduler.
// Deprecated in 1.7, please use the bindings subresource of pods instead.
type Binding struct {
	metav1.TypeMeta `json:",inline"`
	// Standard object's metadata.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// The target object that you want to bind to the standard object.
	Target ObjectReference `json:"target" protobuf:"bytes,2,opt,name=target"`
}

// Preconditions must be fulfilled before an operation (update, delete, etc.) is carried out.
// +k8s:openapi-gen=false
type Preconditions struct {
	// Specifies the target UID.
	// +optional
	UID *types.UID `json:"uid,omitempty" protobuf:"bytes,1,opt,name=uid,casttype=k8s.io/apimachinery/pkg/types.UID"`
}

// +k8s:conversion-gen:explicit-from=net/url.Values
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// PodLogOptions is the query options for a Pod's logs REST call.
type PodLogOptions struct {
	metav1.TypeMeta `json:",inline"`

	// The container for which to stream logs. Defaults to only container if there is one container in the pod.
	// +optional
	Container string `json:"container,omitempty" protobuf:"bytes,1,opt,name=container"`
	// Follow the log stream of the pod. Defaults to false.
	// +optional
	Follow bool `json:"follow,omitempty" protobuf:"varint,2,opt,name=follow"`
	// Return previous terminated container logs. Defaults to false.
	// +optional
	Previous bool `json:"previous,omitempty" protobuf:"varint,3,opt,name=previous"`
	// A relative time in seconds before the current time from which to show logs. If this value
	// precedes the time a pod was started, only logs since the pod start will be returned.
	// If this value is in the future, no logs will be returned.
	// Only one of sinceSeconds or sinceTime may be specified.
	// +optional
	SinceSeconds *int64 `json:"sinceSeconds,omitempty" protobuf:"varint,4,opt,name=sinceSeconds"`
	// An RFC3339 timestamp from which to show logs. If this value
	// precedes the time a pod was started, only logs since the pod start will be returned.
	// If this value is in the future, no logs will be returned.
	// Only one of sinceSeconds or sinceTime may be specified.
	// +optional
	SinceTime *metav1.Time `json:"sinceTime,omitempty" protobuf:"bytes,5,opt,name=sinceTime"`
	// If true, add an RFC3339 or RFC3339Nano timestamp at the beginning of every line
	// of log output. Defaults to false.
	// +optional
	Timestamps bool `json:"timestamps,omitempty" protobuf:"varint,6,opt,name=timestamps"`
	// If set, the number of lines from the end of the logs to show. If not specified,
	// logs are shown from the creation of the container or sinceSeconds or sinceTime
	// +optional
	TailLines *int64 `json:"tailLines,omitempty" protobuf:"varint,7,opt,name=tailLines"`
	// If set, the number of bytes to read from the server before terminating the
	// log output. This may not display a complete final line of logging, and may return
	// slightly more or slightly less than the specified limit.
	// +optional
	LimitBytes *int64 `json:"limitBytes,omitempty" protobuf:"varint,8,opt,name=limitBytes"`

	// insecureSkipTLSVerifyBackend indicates that the apiserver should not confirm the validity of the
	// serving certificate of the backend it is connecting to.  This will make the HTTPS connection between the apiserver
	// and the backend insecure. This means the apiserver cannot verify the log data it is receiving came from the real
	// kubelet.  If the kubelet is configured to verify the apiserver's TLS credentials, it does not mean the
	// connection to the real kubelet is vulnerable to a man in the middle attack (e.g. an attacker could not intercept
	// the actual log data coming from the real kubelet).
	// +optional
	InsecureSkipTLSVerifyBackend bool `json:"insecureSkipTLSVerifyBackend,omitempty" protobuf:"varint,9,opt,name=insecureSkipTLSVerifyBackend"`
}

// +k8s:conversion-gen:explicit-from=net/url.Values
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// PodAttachOptions is the query options to a Pod's remote attach call.
// ---
// TODO: merge w/ PodExecOptions below for stdin, stdout, etc
// and also when we cut V2, we should export a "StreamOptions" or somesuch that contains Stdin, Stdout, Stder and TTY
type PodAttachOptions struct {
	metav1.TypeMeta `json:",inline"`

	// Stdin if true, redirects the standard input stream of the pod for this call.
	// Defaults to false.
	// +optional
	Stdin bool `json:"stdin,omitempty" protobuf:"varint,1,opt,name=stdin"`

	// Stdout if true indicates that stdout is to be redirected for the attach call.
	// Defaults to true.
	// +optional
	Stdout bool `json:"stdout,omitempty" protobuf:"varint,2,opt,name=stdout"`

	// Stderr if true indicates that stderr is to be redirected for the attach call.
	// Defaults to true.
	// +optional
	Stderr bool `json:"stderr,omitempty" protobuf:"varint,3,opt,name=stderr"`

	// TTY if true indicates that a tty will be allocated for the attach call.
	// This is passed through the container runtime so the tty
	// is allocated on the worker node by the container runtime.
	// Defaults to false.
	// +optional
	TTY bool `json:"tty,omitempty" protobuf:"varint,4,opt,name=tty"`

	// The container in which to execute the command.
	// Defaults to only container if there is only one container in the pod.
	// +optional
	Container string `json:"container,omitempty" protobuf:"bytes,5,opt,name=container"`
}

// +k8s:conversion-gen:explicit-from=net/url.Values
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// PodExecOptions is the query options to a Pod's remote exec call.
// ---
// TODO: This is largely identical to PodAttachOptions above, make sure they stay in sync and see about merging
// and also when we cut V2, we should export a "StreamOptions" or somesuch that contains Stdin, Stdout, Stder and TTY
type PodExecOptions struct {
	metav1.TypeMeta `json:",inline"`

	// Redirect the standard input stream of the pod for this call.
	// Defaults to false.
	// +optional
	Stdin bool `json:"stdin,omitempty" protobuf:"varint,1,opt,name=stdin"`

	// Redirect the standard output stream of the pod for this call.
	// +optional
	Stdout bool `json:"stdout,omitempty" protobuf:"varint,2,opt,name=stdout"`

	// Redirect the standard error stream of the pod for this call.
	// +optional
	Stderr bool `json:"stderr,omitempty" protobuf:"varint,3,opt,name=stderr"`

	// TTY if true indicates that a tty will be allocated for the exec call.
	// Defaults to false.
	// +optional
	TTY bool `json:"tty,omitempty" protobuf:"varint,4,opt,name=tty"`

	// Container in which to execute the command.
	// Defaults to only container if there is only one container in the pod.
	// +optional
	Container string `json:"container,omitempty" protobuf:"bytes,5,opt,name=container"`

	// Command is the remote command to execute. argv array. Not executed within a shell.
	Command []string `json:"command" protobuf:"bytes,6,rep,name=command"`
}

// +k8s:conversion-gen:explicit-from=net/url.Values
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// PodPortForwardOptions is the query options to a Pod's port forward call
// when using WebSockets.
// The `port` query parameter must specify the port or
// ports (comma separated) to forward over.
// Port forwarding over SPDY does not use these options. It requires the port
// to be passed in the `port` header as part of request.
type PodPortForwardOptions struct {
	metav1.TypeMeta `json:",inline"`

	// List of ports to forward
	// Required when using WebSockets
	// +optional
	Ports []int32 `json:"ports,omitempty" protobuf:"varint,1,rep,name=ports"`
}

// +k8s:conversion-gen:explicit-from=net/url.Values
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// PodProxyOptions is the query options to a Pod's proxy call.
type PodProxyOptions struct {
	metav1.TypeMeta `json:",inline"`

	// Path is the URL path to use for the current proxy request to pod.
	// +optional
	Path string `json:"path,omitempty" protobuf:"bytes,1,opt,name=path"`
}

// +k8s:conversion-gen:explicit-from=net/url.Values
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// NodeProxyOptions is the query options to a Node's proxy call.
type NodeProxyOptions struct {
	metav1.TypeMeta `json:",inline"`

	// Path is the URL path to use for the current proxy request to node.
	// +optional
	Path string `json:"path,omitempty" protobuf:"bytes,1,opt,name=path"`
}

// +k8s:conversion-gen:explicit-from=net/url.Values
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ServiceProxyOptions is the query options to a Service's proxy call.
type ServiceProxyOptions struct {
	metav1.TypeMeta `json:",inline"`

	// Path is the part of URLs that include service endpoints, suffixes,
	// and parameters to use for the current proxy request to service.
	// For example, the whole request URL is
	// http://localhost/api/v1/namespaces/kube-system/services/elasticsearch-logging/_search?q=user:kimchy.
	// Path is _search?q=user:kimchy.
	// +optional
	Path string `json:"path,omitempty" protobuf:"bytes,1,opt,name=path"`
}

// ObjectReference contains enough information to let you inspect or modify the referred object.
// ---
// New uses of this type are discouraged because of difficulty describing its usage when embedded in APIs.
//  1. Ignored fields.  It includes many fields which are not generally honored.  For instance, ResourceVersion and FieldPath are both very rarely valid in actual usage.
//  2. Invalid usage help.  It is impossible to add specific help for individual usage.  In most embedded usages, there are particular
//     restrictions like, "must refer only to types A and B" or "UID not honored" or "name must be restricted".
//     Those cannot be well described when embedded.
//  3. Inconsistent validation.  Because the usages are different, the validation rules are different by usage, which makes it hard for users to predict what will happen.
//  4. The fields are both imprecise and overly precise.  Kind is not a precise mapping to a URL. This can produce ambiguity
//     during interpretation and require a REST mapping.  In most cases, the dependency is on the group,resource tuple
//     and the version of the actual struct is irrelevant.
//  5. We cannot easily change it.  Because this type is embedded in many locations, updates to this type
//     will affect numerous schemas.  Don't make new APIs embed an underspecified API type they do not control.
//
// Instead of using this type, create a locally provided and used type that is well-focused on your reference.
// For example, ServiceReferences for admission registration: https://github.com/kubernetes/api/blob/release-1.17/admissionregistration/v1/types.go#L533 .
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
// +structType=atomic
// ObjectReference 包含足够的信息让你检查或修改引用的对象。
// ---
// 不鼓励使用此类型的新用法，因为在嵌入 API 中难以描述其用法。
//  1. 忽略的字段。 它包含许多通常不被一般使用的字段。 例如，ResourceVersion 和 FieldPath 都很少在实际使用中有效。
//  2. 无效的使用帮助。 无法为单个用法添加特定的帮助。 在大多数嵌入用法中，有特定的限制，例如，“必须仅引用类型 A 和 B”或“UID 不受支持”或“名称必须受限制”。
//     这些在嵌入时无法很好地描述。
//  3. 不一致的验证。 由于用法不同，因此验证规则因用法而异，这使得用户难以预测将发生什么。
//  4. 字段既不精确又过于精确。 Kind 不是 URL 的精确映射。 这可能会导致解释和要求 REST 映射的歧义。
//     在大多数情况下，依赖关系是在 group,resource 元组和实际结构的版本是不相关的。
//  5. 我们无法轻松更改它。 因为此类型嵌入在许多位置中，对此类型的更新将影响许多模式。 不要使新 API 嵌入他们不控制的不明确的 API 类型。
//
// 请改用本地提供和使用的类型，该类型专注于您的引用。 例如，用于准入注册的 ServiceReferences：
type ObjectReference struct {
	// Kind of the referent.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds
	// +optional
	// Kind的引用。
	Kind string `json:"kind,omitempty" protobuf:"bytes,1,opt,name=kind"`
	// Namespace of the referent.
	// More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/namespaces/
	// +optional
	// 引用的命名空间。
	Namespace string `json:"namespace,omitempty" protobuf:"bytes,2,opt,name=namespace"`
	// Name of the referent.
	// More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names
	// +optional
	// 引用的名称。
	Name string `json:"name,omitempty" protobuf:"bytes,3,opt,name=name"`
	// UID of the referent.
	// More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#uids
	// +optional
	// 引用的UID。
	UID types.UID `json:"uid,omitempty" protobuf:"bytes,4,opt,name=uid,casttype=k8s.io/apimachinery/pkg/types.UID"`
	// API version of the referent.
	// +optional
	// 引用的API版本。
	APIVersion string `json:"apiVersion,omitempty" protobuf:"bytes,5,opt,name=apiVersion"`
	// Specific resourceVersion to which this reference is made, if any.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#concurrency-control-and-consistency
	// +optional
	// 对此引用所做的特定 resourceVersion（如果有）。
	ResourceVersion string `json:"resourceVersion,omitempty" protobuf:"bytes,6,opt,name=resourceVersion"`

	// If referring to a piece of an object instead of an entire object, this string
	// should contain a valid JSON/Go field access statement, such as desiredState.manifest.containers[2].
	// For example, if the object reference is to a container within a pod, this would take on a value like:
	// "spec.containers{name}" (where "name" refers to the name of the container that triggered
	// the event) or if no container name is specified "spec.containers[2]" (container with
	// index 2 in this pod). This syntax is chosen only to have some well-defined way of
	// referencing a part of an object.
	// TODO: this design is not final and this field is subject to change in the future.
	// +optional
	// 如果引用对象的一部分而不是整个对象，则此字符串应包含有效的 JSON/Go 字段访问语句，例如 desiredState.manifest.containers[2]。
	// 例如，如果对象引用是对 pod 中的容器的引用，则其值将为： "spec.containers{name}"（其中 "name" 指的是触发事件的容器的名称），
	// 或者如果未指定容器名称，则为 "spec.containers[2]"（此 pod 中索引为 2 的容器）。
	// 仅选择此语法是为了有一种明确定义的方式来引用对象的一部分。
	FieldPath string `json:"fieldPath,omitempty" protobuf:"bytes,7,opt,name=fieldPath"`
}

// LocalObjectReference contains enough information to let you locate the
// referenced object inside the same namespace.
// +structType=atomic
// LocalObjectReference 包含足够的信息让你在同一个命名空间中定位到对象
type LocalObjectReference struct {
	// Name of the referent.
	// More info: https://kubernetes.io/docs/concepts/overview/working-with-objects/names/#names
	// TODO: Add other useful fields. apiVersion, kind, uid?
	// +optional
	// Name 是对象的名称
	Name string `json:"name,omitempty" protobuf:"bytes,1,opt,name=name"`
}

// TypedLocalObjectReference contains enough information to let you locate the
// typed referenced object inside the same namespace.
// +structType=atomic
// TypedLocalObjectReference 包含足够的信息让你在同一个命名空间中定位到指定类型的对象
type TypedLocalObjectReference struct {
	// APIGroup is the group for the resource being referenced.
	// If APIGroup is not specified, the specified Kind must be in the core API group.
	// For any other third-party types, APIGroup is required.
	// +optional
	// APIGroup 是资源所在的组，如果没有指定，Kind 必须是 core 组的资源
	APIGroup *string `json:"apiGroup" protobuf:"bytes,1,opt,name=apiGroup"`
	// Kind is the type of resource being referenced
	Kind string `json:"kind" protobuf:"bytes,2,opt,name=kind"`
	// Name is the name of resource being referenced
	Name string `json:"name" protobuf:"bytes,3,opt,name=name"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// SerializedReference is a reference to serialized object.
type SerializedReference struct {
	metav1.TypeMeta `json:",inline"`
	// The reference to an object in the system.
	// +optional
	Reference ObjectReference `json:"reference,omitempty" protobuf:"bytes,1,opt,name=reference"`
}

// EventSource contains information for an event.
type EventSource struct {
	// Component from which the event is generated.
	// +optional
	Component string `json:"component,omitempty" protobuf:"bytes,1,opt,name=component"`
	// Node name on which the event is generated.
	// +optional
	Host string `json:"host,omitempty" protobuf:"bytes,2,opt,name=host"`
}

// Valid values for event types (new types could be added in future)
const (
	// Information only and will not cause any problems
	EventTypeNormal string = "Normal"
	// These events are to warn that something might go wrong
	EventTypeWarning string = "Warning"
)

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Event is a report of an event somewhere in the cluster.  Events
// have a limited retention time and triggers and messages may evolve
// with time.  Event consumers should not rely on the timing of an event
// with a given Reason reflecting a consistent underlying trigger, or the
// continued existence of events with that Reason.  Events should be
// treated as informative, best-effort, supplemental data.
type Event struct {
	metav1.TypeMeta `json:",inline"`
	// Standard object's metadata.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
	metav1.ObjectMeta `json:"metadata" protobuf:"bytes,1,opt,name=metadata"`

	// The object that this event is about.
	InvolvedObject ObjectReference `json:"involvedObject" protobuf:"bytes,2,opt,name=involvedObject"`

	// This should be a short, machine understandable string that gives the reason
	// for the transition into the object's current status.
	// TODO: provide exact specification for format.
	// +optional
	Reason string `json:"reason,omitempty" protobuf:"bytes,3,opt,name=reason"`

	// A human-readable description of the status of this operation.
	// TODO: decide on maximum length.
	// +optional
	Message string `json:"message,omitempty" protobuf:"bytes,4,opt,name=message"`

	// The component reporting this event. Should be a short machine understandable string.
	// +optional
	Source EventSource `json:"source,omitempty" protobuf:"bytes,5,opt,name=source"`

	// The time at which the event was first recorded. (Time of server receipt is in TypeMeta.)
	// +optional
	FirstTimestamp metav1.Time `json:"firstTimestamp,omitempty" protobuf:"bytes,6,opt,name=firstTimestamp"`

	// The time at which the most recent occurrence of this event was recorded.
	// +optional
	LastTimestamp metav1.Time `json:"lastTimestamp,omitempty" protobuf:"bytes,7,opt,name=lastTimestamp"`

	// The number of times this event has occurred.
	// +optional
	Count int32 `json:"count,omitempty" protobuf:"varint,8,opt,name=count"`

	// Type of this event (Normal, Warning), new types could be added in the future
	// +optional
	Type string `json:"type,omitempty" protobuf:"bytes,9,opt,name=type"`

	// Time when this Event was first observed.
	// +optional
	EventTime metav1.MicroTime `json:"eventTime,omitempty" protobuf:"bytes,10,opt,name=eventTime"`

	// Data about the Event series this event represents or nil if it's a singleton Event.
	// +optional
	Series *EventSeries `json:"series,omitempty" protobuf:"bytes,11,opt,name=series"`

	// What action was taken/failed regarding to the Regarding object.
	// +optional
	Action string `json:"action,omitempty" protobuf:"bytes,12,opt,name=action"`

	// Optional secondary object for more complex actions.
	// +optional
	Related *ObjectReference `json:"related,omitempty" protobuf:"bytes,13,opt,name=related"`

	// Name of the controller that emitted this Event, e.g. `kubernetes.io/kubelet`.
	// +optional
	ReportingController string `json:"reportingComponent" protobuf:"bytes,14,opt,name=reportingComponent"`

	// ID of the controller instance, e.g. `kubelet-xyzf`.
	// +optional
	ReportingInstance string `json:"reportingInstance" protobuf:"bytes,15,opt,name=reportingInstance"`
}

// EventSeries contain information on series of events, i.e. thing that was/is happening
// continuously for some time.
type EventSeries struct {
	// Number of occurrences in this series up to the last heartbeat time
	Count int32 `json:"count,omitempty" protobuf:"varint,1,name=count"`
	// Time of the last occurrence observed
	LastObservedTime metav1.MicroTime `json:"lastObservedTime,omitempty" protobuf:"bytes,2,name=lastObservedTime"`

	// +k8s:deprecated=state,protobuf=3
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// EventList is a list of events.
type EventList struct {
	metav1.TypeMeta `json:",inline"`
	// Standard list metadata.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds
	// +optional
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// List of events
	Items []Event `json:"items" protobuf:"bytes,2,rep,name=items"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// List holds a list of objects, which may not be known by the server.
type List metav1.List

// LimitType is a type of object that is limited. It can be Pod, Container, PersistentVolumeClaim or
// a fully qualified resource name.
type LimitType string

const (
	// Limit that applies to all pods in a namespace
	LimitTypePod LimitType = "Pod"
	// Limit that applies to all containers in a namespace
	LimitTypeContainer LimitType = "Container"
	// Limit that applies to all persistent volume claims in a namespace
	LimitTypePersistentVolumeClaim LimitType = "PersistentVolumeClaim"
)

// LimitRangeItem defines a min/max usage limit for any resource that matches on kind.
type LimitRangeItem struct {
	// Type of resource that this limit applies to.
	Type LimitType `json:"type" protobuf:"bytes,1,opt,name=type,casttype=LimitType"`
	// Max usage constraints on this kind by resource name.
	// +optional
	Max ResourceList `json:"max,omitempty" protobuf:"bytes,2,rep,name=max,casttype=ResourceList,castkey=ResourceName"`
	// Min usage constraints on this kind by resource name.
	// +optional
	Min ResourceList `json:"min,omitempty" protobuf:"bytes,3,rep,name=min,casttype=ResourceList,castkey=ResourceName"`
	// Default resource requirement limit value by resource name if resource limit is omitted.
	// +optional
	Default ResourceList `json:"default,omitempty" protobuf:"bytes,4,rep,name=default,casttype=ResourceList,castkey=ResourceName"`
	// DefaultRequest is the default resource requirement request value by resource name if resource request is omitted.
	// +optional
	DefaultRequest ResourceList `json:"defaultRequest,omitempty" protobuf:"bytes,5,rep,name=defaultRequest,casttype=ResourceList,castkey=ResourceName"`
	// MaxLimitRequestRatio if specified, the named resource must have a request and limit that are both non-zero where limit divided by request is less than or equal to the enumerated value; this represents the max burst for the named resource.
	// +optional
	MaxLimitRequestRatio ResourceList `json:"maxLimitRequestRatio,omitempty" protobuf:"bytes,6,rep,name=maxLimitRequestRatio,casttype=ResourceList,castkey=ResourceName"`
}

// LimitRangeSpec defines a min/max usage limit for resources that match on kind.
type LimitRangeSpec struct {
	// Limits is the list of LimitRangeItem objects that are enforced.
	Limits []LimitRangeItem `json:"limits" protobuf:"bytes,1,rep,name=limits"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// LimitRange sets resource usage limits for each kind of resource in a Namespace.
type LimitRange struct {
	metav1.TypeMeta `json:",inline"`
	// Standard object's metadata.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// Spec defines the limits enforced.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#spec-and-status
	// +optional
	Spec LimitRangeSpec `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// LimitRangeList is a list of LimitRange items.
type LimitRangeList struct {
	metav1.TypeMeta `json:",inline"`
	// Standard list metadata.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds
	// +optional
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// Items is a list of LimitRange objects.
	// More info: https://kubernetes.io/docs/concepts/configuration/manage-resources-containers/
	Items []LimitRange `json:"items" protobuf:"bytes,2,rep,name=items"`
}

// The following identify resource constants for Kubernetes object types
const (
	// Pods, number
	ResourcePods ResourceName = "pods"
	// Services, number
	ResourceServices ResourceName = "services"
	// ReplicationControllers, number
	ResourceReplicationControllers ResourceName = "replicationcontrollers"
	// ResourceQuotas, number
	ResourceQuotas ResourceName = "resourcequotas"
	// ResourceSecrets, number
	ResourceSecrets ResourceName = "secrets"
	// ResourceConfigMaps, number
	ResourceConfigMaps ResourceName = "configmaps"
	// ResourcePersistentVolumeClaims, number
	ResourcePersistentVolumeClaims ResourceName = "persistentvolumeclaims"
	// ResourceServicesNodePorts, number
	ResourceServicesNodePorts ResourceName = "services.nodeports"
	// ResourceServicesLoadBalancers, number
	ResourceServicesLoadBalancers ResourceName = "services.loadbalancers"
	// CPU request, in cores. (500m = .5 cores)
	ResourceRequestsCPU ResourceName = "requests.cpu"
	// Memory request, in bytes. (500Gi = 500GiB = 500 * 1024 * 1024 * 1024)
	ResourceRequestsMemory ResourceName = "requests.memory"
	// Storage request, in bytes
	ResourceRequestsStorage ResourceName = "requests.storage"
	// Local ephemeral storage request, in bytes. (500Gi = 500GiB = 500 * 1024 * 1024 * 1024)
	ResourceRequestsEphemeralStorage ResourceName = "requests.ephemeral-storage"
	// CPU limit, in cores. (500m = .5 cores)
	ResourceLimitsCPU ResourceName = "limits.cpu"
	// Memory limit, in bytes. (500Gi = 500GiB = 500 * 1024 * 1024 * 1024)
	ResourceLimitsMemory ResourceName = "limits.memory"
	// Local ephemeral storage limit, in bytes. (500Gi = 500GiB = 500 * 1024 * 1024 * 1024)
	ResourceLimitsEphemeralStorage ResourceName = "limits.ephemeral-storage"
)

// The following identify resource prefix for Kubernetes object types
const (
	// HugePages request, in bytes. (500Gi = 500GiB = 500 * 1024 * 1024 * 1024)
	// As burst is not supported for HugePages, we would only quota its request, and ignore the limit.
	ResourceRequestsHugePagesPrefix = "requests.hugepages-"
	// Default resource requests prefix
	DefaultResourceRequestsPrefix = "requests."
)

// A ResourceQuotaScope defines a filter that must match each object tracked by a quota
// +enum
type ResourceQuotaScope string

const (
	// Match all pod objects where spec.activeDeadlineSeconds >=0
	ResourceQuotaScopeTerminating ResourceQuotaScope = "Terminating"
	// Match all pod objects where spec.activeDeadlineSeconds is nil
	ResourceQuotaScopeNotTerminating ResourceQuotaScope = "NotTerminating"
	// Match all pod objects that have best effort quality of service
	ResourceQuotaScopeBestEffort ResourceQuotaScope = "BestEffort"
	// Match all pod objects that do not have best effort quality of service
	ResourceQuotaScopeNotBestEffort ResourceQuotaScope = "NotBestEffort"
	// Match all pod objects that have priority class mentioned
	ResourceQuotaScopePriorityClass ResourceQuotaScope = "PriorityClass"
	// Match all pod objects that have cross-namespace pod (anti)affinity mentioned.
	ResourceQuotaScopeCrossNamespacePodAffinity ResourceQuotaScope = "CrossNamespacePodAffinity"
)

// ResourceQuotaSpec defines the desired hard limits to enforce for Quota.
type ResourceQuotaSpec struct {
	// hard is the set of desired hard limits for each named resource.
	// More info: https://kubernetes.io/docs/concepts/policy/resource-quotas/
	// +optional
	Hard ResourceList `json:"hard,omitempty" protobuf:"bytes,1,rep,name=hard,casttype=ResourceList,castkey=ResourceName"`
	// A collection of filters that must match each object tracked by a quota.
	// If not specified, the quota matches all objects.
	// +optional
	Scopes []ResourceQuotaScope `json:"scopes,omitempty" protobuf:"bytes,2,rep,name=scopes,casttype=ResourceQuotaScope"`
	// scopeSelector is also a collection of filters like scopes that must match each object tracked by a quota
	// but expressed using ScopeSelectorOperator in combination with possible values.
	// For a resource to match, both scopes AND scopeSelector (if specified in spec), must be matched.
	// +optional
	ScopeSelector *ScopeSelector `json:"scopeSelector,omitempty" protobuf:"bytes,3,opt,name=scopeSelector"`
}

// A scope selector represents the AND of the selectors represented
// by the scoped-resource selector requirements.
// +structType=atomic
type ScopeSelector struct {
	// A list of scope selector requirements by scope of the resources.
	// +optional
	MatchExpressions []ScopedResourceSelectorRequirement `json:"matchExpressions,omitempty" protobuf:"bytes,1,rep,name=matchExpressions"`
}

// A scoped-resource selector requirement is a selector that contains values, a scope name, and an operator
// that relates the scope name and values.
type ScopedResourceSelectorRequirement struct {
	// The name of the scope that the selector applies to.
	ScopeName ResourceQuotaScope `json:"scopeName" protobuf:"bytes,1,opt,name=scopeName"`
	// Represents a scope's relationship to a set of values.
	// Valid operators are In, NotIn, Exists, DoesNotExist.
	Operator ScopeSelectorOperator `json:"operator" protobuf:"bytes,2,opt,name=operator,casttype=ScopedResourceSelectorOperator"`
	// An array of string values. If the operator is In or NotIn,
	// the values array must be non-empty. If the operator is Exists or DoesNotExist,
	// the values array must be empty.
	// This array is replaced during a strategic merge patch.
	// +optional
	Values []string `json:"values,omitempty" protobuf:"bytes,3,rep,name=values"`
}

// A scope selector operator is the set of operators that can be used in
// a scope selector requirement.
// +enum
type ScopeSelectorOperator string

const (
	ScopeSelectorOpIn           ScopeSelectorOperator = "In"
	ScopeSelectorOpNotIn        ScopeSelectorOperator = "NotIn"
	ScopeSelectorOpExists       ScopeSelectorOperator = "Exists"
	ScopeSelectorOpDoesNotExist ScopeSelectorOperator = "DoesNotExist"
)

// ResourceQuotaStatus defines the enforced hard limits and observed use.
type ResourceQuotaStatus struct {
	// Hard is the set of enforced hard limits for each named resource.
	// More info: https://kubernetes.io/docs/concepts/policy/resource-quotas/
	// +optional
	Hard ResourceList `json:"hard,omitempty" protobuf:"bytes,1,rep,name=hard,casttype=ResourceList,castkey=ResourceName"`
	// Used is the current observed total usage of the resource in the namespace.
	// +optional
	Used ResourceList `json:"used,omitempty" protobuf:"bytes,2,rep,name=used,casttype=ResourceList,castkey=ResourceName"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ResourceQuota sets aggregate quota restrictions enforced per namespace
type ResourceQuota struct {
	metav1.TypeMeta `json:",inline"`
	// Standard object's metadata.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// Spec defines the desired quota.
	// https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#spec-and-status
	// +optional
	Spec ResourceQuotaSpec `json:"spec,omitempty" protobuf:"bytes,2,opt,name=spec"`

	// Status defines the actual enforced quota and its current usage.
	// https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#spec-and-status
	// +optional
	Status ResourceQuotaStatus `json:"status,omitempty" protobuf:"bytes,3,opt,name=status"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ResourceQuotaList is a list of ResourceQuota items.
type ResourceQuotaList struct {
	metav1.TypeMeta `json:",inline"`
	// Standard list metadata.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds
	// +optional
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// Items is a list of ResourceQuota objects.
	// More info: https://kubernetes.io/docs/concepts/policy/resource-quotas/
	Items []ResourceQuota `json:"items" protobuf:"bytes,2,rep,name=items"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Secret holds secret data of a certain type. The total bytes of the values in
// the Data field must be less than MaxSecretSize bytes.
type Secret struct {
	metav1.TypeMeta `json:",inline"`
	// Standard object's metadata.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// Immutable, if set to true, ensures that data stored in the Secret cannot
	// be updated (only object metadata can be modified).
	// If not set to true, the field can be modified at any time.
	// Defaulted to nil.
	// +optional
	Immutable *bool `json:"immutable,omitempty" protobuf:"varint,5,opt,name=immutable"`

	// Data contains the secret data. Each key must consist of alphanumeric
	// characters, '-', '_' or '.'. The serialized form of the secret data is a
	// base64 encoded string, representing the arbitrary (possibly non-string)
	// data value here. Described in https://tools.ietf.org/html/rfc4648#section-4
	// +optional
	Data map[string][]byte `json:"data,omitempty" protobuf:"bytes,2,rep,name=data"`

	// stringData allows specifying non-binary secret data in string form.
	// It is provided as a write-only input field for convenience.
	// All keys and values are merged into the data field on write, overwriting any existing values.
	// The stringData field is never output when reading from the API.
	// +k8s:conversion-gen=false
	// +optional
	StringData map[string]string `json:"stringData,omitempty" protobuf:"bytes,4,rep,name=stringData"`

	// Used to facilitate programmatic handling of secret data.
	// More info: https://kubernetes.io/docs/concepts/configuration/secret/#secret-types
	// +optional
	Type SecretType `json:"type,omitempty" protobuf:"bytes,3,opt,name=type,casttype=SecretType"`
}

const MaxSecretSize = 1 * 1024 * 1024

type SecretType string

const (
	// SecretTypeOpaque is the default. Arbitrary user-defined data
	SecretTypeOpaque SecretType = "Opaque"

	// SecretTypeServiceAccountToken contains a token that identifies a service account to the API
	//
	// Required fields:
	// - Secret.Annotations["kubernetes.io/service-account.name"] - the name of the ServiceAccount the token identifies
	// - Secret.Annotations["kubernetes.io/service-account.uid"] - the UID of the ServiceAccount the token identifies
	// - Secret.Data["token"] - a token that identifies the service account to the API
	SecretTypeServiceAccountToken SecretType = "kubernetes.io/service-account-token"

	// ServiceAccountNameKey is the key of the required annotation for SecretTypeServiceAccountToken secrets
	ServiceAccountNameKey = "kubernetes.io/service-account.name"
	// ServiceAccountUIDKey is the key of the required annotation for SecretTypeServiceAccountToken secrets
	ServiceAccountUIDKey = "kubernetes.io/service-account.uid"
	// ServiceAccountTokenKey is the key of the required data for SecretTypeServiceAccountToken secrets
	ServiceAccountTokenKey = "token"
	// ServiceAccountKubeconfigKey is the key of the optional kubeconfig data for SecretTypeServiceAccountToken secrets
	ServiceAccountKubeconfigKey = "kubernetes.kubeconfig"
	// ServiceAccountRootCAKey is the key of the optional root certificate authority for SecretTypeServiceAccountToken secrets
	ServiceAccountRootCAKey = "ca.crt"
	// ServiceAccountNamespaceKey is the key of the optional namespace to use as the default for namespaced API calls
	ServiceAccountNamespaceKey = "namespace"

	// SecretTypeDockercfg contains a dockercfg file that follows the same format rules as ~/.dockercfg
	//
	// Required fields:
	// - Secret.Data[".dockercfg"] - a serialized ~/.dockercfg file
	SecretTypeDockercfg SecretType = "kubernetes.io/dockercfg"

	// DockerConfigKey is the key of the required data for SecretTypeDockercfg secrets
	DockerConfigKey = ".dockercfg"

	// SecretTypeDockerConfigJson contains a dockercfg file that follows the same format rules as ~/.docker/config.json
	//
	// Required fields:
	// - Secret.Data[".dockerconfigjson"] - a serialized ~/.docker/config.json file
	SecretTypeDockerConfigJson SecretType = "kubernetes.io/dockerconfigjson"

	// DockerConfigJsonKey is the key of the required data for SecretTypeDockerConfigJson secrets
	DockerConfigJsonKey = ".dockerconfigjson"

	// SecretTypeBasicAuth contains data needed for basic authentication.
	//
	// Required at least one of fields:
	// - Secret.Data["username"] - username used for authentication
	// - Secret.Data["password"] - password or token needed for authentication
	SecretTypeBasicAuth SecretType = "kubernetes.io/basic-auth"

	// BasicAuthUsernameKey is the key of the username for SecretTypeBasicAuth secrets
	BasicAuthUsernameKey = "username"
	// BasicAuthPasswordKey is the key of the password or token for SecretTypeBasicAuth secrets
	BasicAuthPasswordKey = "password"

	// SecretTypeSSHAuth contains data needed for SSH authetication.
	//
	// Required field:
	// - Secret.Data["ssh-privatekey"] - private SSH key needed for authentication
	SecretTypeSSHAuth SecretType = "kubernetes.io/ssh-auth"

	// SSHAuthPrivateKey is the key of the required SSH private key for SecretTypeSSHAuth secrets
	SSHAuthPrivateKey = "ssh-privatekey"
	// SecretTypeTLS contains information about a TLS client or server secret. It
	// is primarily used with TLS termination of the Ingress resource, but may be
	// used in other types.
	//
	// Required fields:
	// - Secret.Data["tls.key"] - TLS private key.
	//   Secret.Data["tls.crt"] - TLS certificate.
	// TODO: Consider supporting different formats, specifying CA/destinationCA.
	SecretTypeTLS SecretType = "kubernetes.io/tls"

	// TLSCertKey is the key for tls certificates in a TLS secret.
	TLSCertKey = "tls.crt"
	// TLSPrivateKeyKey is the key for the private key field in a TLS secret.
	TLSPrivateKeyKey = "tls.key"
	// SecretTypeBootstrapToken is used during the automated bootstrap process (first
	// implemented by kubeadm). It stores tokens that are used to sign well known
	// ConfigMaps. They are used for authn.
	SecretTypeBootstrapToken SecretType = "bootstrap.kubernetes.io/token"
)

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// SecretList is a list of Secret.
type SecretList struct {
	metav1.TypeMeta `json:",inline"`
	// Standard list metadata.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds
	// +optional
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// Items is a list of secret objects.
	// More info: https://kubernetes.io/docs/concepts/configuration/secret
	Items []Secret `json:"items" protobuf:"bytes,2,rep,name=items"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ConfigMap holds configuration data for pods to consume.
type ConfigMap struct {
	metav1.TypeMeta `json:",inline"`
	// Standard object's metadata.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// Immutable, if set to true, ensures that data stored in the ConfigMap cannot
	// be updated (only object metadata can be modified).
	// If not set to true, the field can be modified at any time.
	// Defaulted to nil.
	// +optional
	Immutable *bool `json:"immutable,omitempty" protobuf:"varint,4,opt,name=immutable"`

	// Data contains the configuration data.
	// Each key must consist of alphanumeric characters, '-', '_' or '.'.
	// Values with non-UTF-8 byte sequences must use the BinaryData field.
	// The keys stored in Data must not overlap with the keys in
	// the BinaryData field, this is enforced during validation process.
	// +optional
	Data map[string]string `json:"data,omitempty" protobuf:"bytes,2,rep,name=data"`

	// BinaryData contains the binary data.
	// Each key must consist of alphanumeric characters, '-', '_' or '.'.
	// BinaryData can contain byte sequences that are not in the UTF-8 range.
	// The keys stored in BinaryData must not overlap with the ones in
	// the Data field, this is enforced during validation process.
	// Using this field will require 1.10+ apiserver and
	// kubelet.
	// +optional
	BinaryData map[string][]byte `json:"binaryData,omitempty" protobuf:"bytes,3,rep,name=binaryData"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ConfigMapList is a resource containing a list of ConfigMap objects.
type ConfigMapList struct {
	metav1.TypeMeta `json:",inline"`

	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
	// +optional
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// Items is the list of ConfigMaps.
	Items []ConfigMap `json:"items" protobuf:"bytes,2,rep,name=items"`
}

// Type and constants for component health validation.
type ComponentConditionType string

// These are the valid conditions for the component.
const (
	ComponentHealthy ComponentConditionType = "Healthy"
)

// Information about the condition of a component.
type ComponentCondition struct {
	// Type of condition for a component.
	// Valid value: "Healthy"
	Type ComponentConditionType `json:"type" protobuf:"bytes,1,opt,name=type,casttype=ComponentConditionType"`
	// Status of the condition for a component.
	// Valid values for "Healthy": "True", "False", or "Unknown".
	Status ConditionStatus `json:"status" protobuf:"bytes,2,opt,name=status,casttype=ConditionStatus"`
	// Message about the condition for a component.
	// For example, information about a health check.
	// +optional
	Message string `json:"message,omitempty" protobuf:"bytes,3,opt,name=message"`
	// Condition error code for a component.
	// For example, a health check error code.
	// +optional
	Error string `json:"error,omitempty" protobuf:"bytes,4,opt,name=error"`
}

// +genclient
// +genclient:nonNamespaced
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// ComponentStatus (and ComponentStatusList) holds the cluster validation info.
// Deprecated: This API is deprecated in v1.19+
type ComponentStatus struct {
	metav1.TypeMeta `json:",inline"`
	// Standard object's metadata.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// List of component conditions observed
	// +optional
	// +patchMergeKey=type
	// +patchStrategy=merge
	Conditions []ComponentCondition `json:"conditions,omitempty" patchStrategy:"merge" patchMergeKey:"type" protobuf:"bytes,2,rep,name=conditions"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// Status of all the conditions for the component as a list of ComponentStatus objects.
// Deprecated: This API is deprecated in v1.19+
type ComponentStatusList struct {
	metav1.TypeMeta `json:",inline"`
	// Standard list metadata.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#types-kinds
	// +optional
	metav1.ListMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// List of ComponentStatus objects.
	Items []ComponentStatus `json:"items" protobuf:"bytes,2,rep,name=items"`
}

// DownwardAPIVolumeSource represents a volume containing downward API info.
// Downward API volumes support ownership management and SELinux relabeling.
type DownwardAPIVolumeSource struct {
	// Items is a list of downward API volume file
	// +optional
	Items []DownwardAPIVolumeFile `json:"items,omitempty" protobuf:"bytes,1,rep,name=items"`
	// Optional: mode bits to use on created files by default. Must be a
	// Optional: mode bits used to set permissions on created files by default.
	// Must be an octal value between 0000 and 0777 or a decimal value between 0 and 511.
	// YAML accepts both octal and decimal values, JSON requires decimal values for mode bits.
	// Defaults to 0644.
	// Directories within the path are not affected by this setting.
	// This might be in conflict with other options that affect the file
	// mode, like fsGroup, and the result can be other mode bits set.
	// +optional
	DefaultMode *int32 `json:"defaultMode,omitempty" protobuf:"varint,2,opt,name=defaultMode"`
}

const (
	DownwardAPIVolumeSourceDefaultMode int32 = 0644
)

// DownwardAPIVolumeFile represents information to create the file containing the pod field
type DownwardAPIVolumeFile struct {
	// Required: Path is  the relative path name of the file to be created. Must not be absolute or contain the '..' path. Must be utf-8 encoded. The first item of the relative path must not start with '..'
	Path string `json:"path" protobuf:"bytes,1,opt,name=path"`
	// Required: Selects a field of the pod: only annotations, labels, name and namespace are supported.
	// +optional
	FieldRef *ObjectFieldSelector `json:"fieldRef,omitempty" protobuf:"bytes,2,opt,name=fieldRef"`
	// Selects a resource of the container: only resources limits and requests
	// (limits.cpu, limits.memory, requests.cpu and requests.memory) are currently supported.
	// +optional
	ResourceFieldRef *ResourceFieldSelector `json:"resourceFieldRef,omitempty" protobuf:"bytes,3,opt,name=resourceFieldRef"`
	// Optional: mode bits used to set permissions on this file, must be an octal value
	// between 0000 and 0777 or a decimal value between 0 and 511.
	// YAML accepts both octal and decimal values, JSON requires decimal values for mode bits.
	// If not specified, the volume defaultMode will be used.
	// This might be in conflict with other options that affect the file
	// mode, like fsGroup, and the result can be other mode bits set.
	// +optional
	Mode *int32 `json:"mode,omitempty" protobuf:"varint,4,opt,name=mode"`
}

// Represents downward API info for projecting into a projected volume.
// Note that this is identical to a downwardAPI volume source without the default
// mode.
type DownwardAPIProjection struct {
	// Items is a list of DownwardAPIVolume file
	// +optional
	Items []DownwardAPIVolumeFile `json:"items,omitempty" protobuf:"bytes,1,rep,name=items"`
}

// SecurityContext holds security configuration that will be applied to a container.
// Some fields are present in both SecurityContext and PodSecurityContext.  When both
// are set, the values in SecurityContext take precedence.
type SecurityContext struct {
	// The capabilities to add/drop when running containers.
	// Defaults to the default set of capabilities granted by the container runtime.
	// Note that this field cannot be set when spec.os.name is windows.
	// +optional
	Capabilities *Capabilities `json:"capabilities,omitempty" protobuf:"bytes,1,opt,name=capabilities"`
	// Run container in privileged mode.
	// Processes in privileged containers are essentially equivalent to root on the host.
	// Defaults to false.
	// Note that this field cannot be set when spec.os.name is windows.
	// +optional
	Privileged *bool `json:"privileged,omitempty" protobuf:"varint,2,opt,name=privileged"`
	// The SELinux context to be applied to the container.
	// If unspecified, the container runtime will allocate a random SELinux context for each
	// container.  May also be set in PodSecurityContext.  If set in both SecurityContext and
	// PodSecurityContext, the value specified in SecurityContext takes precedence.
	// Note that this field cannot be set when spec.os.name is windows.
	// +optional
	SELinuxOptions *SELinuxOptions `json:"seLinuxOptions,omitempty" protobuf:"bytes,3,opt,name=seLinuxOptions"`
	// The Windows specific settings applied to all containers.
	// If unspecified, the options from the PodSecurityContext will be used.
	// If set in both SecurityContext and PodSecurityContext, the value specified in SecurityContext takes precedence.
	// Note that this field cannot be set when spec.os.name is linux.
	// +optional
	WindowsOptions *WindowsSecurityContextOptions `json:"windowsOptions,omitempty" protobuf:"bytes,10,opt,name=windowsOptions"`
	// The UID to run the entrypoint of the container process.
	// Defaults to user specified in image metadata if unspecified.
	// May also be set in PodSecurityContext.  If set in both SecurityContext and
	// PodSecurityContext, the value specified in SecurityContext takes precedence.
	// Note that this field cannot be set when spec.os.name is windows.
	// +optional
	RunAsUser *int64 `json:"runAsUser,omitempty" protobuf:"varint,4,opt,name=runAsUser"`
	// The GID to run the entrypoint of the container process.
	// Uses runtime default if unset.
	// May also be set in PodSecurityContext.  If set in both SecurityContext and
	// PodSecurityContext, the value specified in SecurityContext takes precedence.
	// Note that this field cannot be set when spec.os.name is windows.
	// +optional
	RunAsGroup *int64 `json:"runAsGroup,omitempty" protobuf:"varint,8,opt,name=runAsGroup"`
	// Indicates that the container must run as a non-root user.
	// If true, the Kubelet will validate the image at runtime to ensure that it
	// does not run as UID 0 (root) and fail to start the container if it does.
	// If unset or false, no such validation will be performed.
	// May also be set in PodSecurityContext.  If set in both SecurityContext and
	// PodSecurityContext, the value specified in SecurityContext takes precedence.
	// +optional
	RunAsNonRoot *bool `json:"runAsNonRoot,omitempty" protobuf:"varint,5,opt,name=runAsNonRoot"`
	// Whether this container has a read-only root filesystem.
	// Default is false.
	// Note that this field cannot be set when spec.os.name is windows.
	// +optional
	ReadOnlyRootFilesystem *bool `json:"readOnlyRootFilesystem,omitempty" protobuf:"varint,6,opt,name=readOnlyRootFilesystem"`
	// AllowPrivilegeEscalation controls whether a process can gain more
	// privileges than its parent process. This bool directly controls if
	// the no_new_privs flag will be set on the container process.
	// AllowPrivilegeEscalation is true always when the container is:
	// 1) run as Privileged
	// 2) has CAP_SYS_ADMIN
	// Note that this field cannot be set when spec.os.name is windows.
	// +optional
	AllowPrivilegeEscalation *bool `json:"allowPrivilegeEscalation,omitempty" protobuf:"varint,7,opt,name=allowPrivilegeEscalation"`
	// procMount denotes the type of proc mount to use for the containers.
	// The default is DefaultProcMount which uses the container runtime defaults for
	// readonly paths and masked paths.
	// This requires the ProcMountType feature flag to be enabled.
	// Note that this field cannot be set when spec.os.name is windows.
	// +optional
	ProcMount *ProcMountType `json:"procMount,omitempty" protobuf:"bytes,9,opt,name=procMount"`
	// The seccomp options to use by this container. If seccomp options are
	// provided at both the pod & container level, the container options
	// override the pod options.
	// Note that this field cannot be set when spec.os.name is windows.
	// +optional
	SeccompProfile *SeccompProfile `json:"seccompProfile,omitempty" protobuf:"bytes,11,opt,name=seccompProfile"`
}

// +enum
type ProcMountType string

const (
	// DefaultProcMount uses the container runtime defaults for readonly and masked
	// paths for /proc.  Most container runtimes mask certain paths in /proc to avoid
	// accidental security exposure of special devices or information.
	DefaultProcMount ProcMountType = "Default"

	// UnmaskedProcMount bypasses the default masking behavior of the container
	// runtime and ensures the newly created /proc the container stays in tact with
	// no modifications.
	UnmaskedProcMount ProcMountType = "Unmasked"
)

// SELinuxOptions are the labels to be applied to the container
type SELinuxOptions struct {
	// User is a SELinux user label that applies to the container.
	// +optional
	User string `json:"user,omitempty" protobuf:"bytes,1,opt,name=user"`
	// Role is a SELinux role label that applies to the container.
	// +optional
	Role string `json:"role,omitempty" protobuf:"bytes,2,opt,name=role"`
	// Type is a SELinux type label that applies to the container.
	// +optional
	Type string `json:"type,omitempty" protobuf:"bytes,3,opt,name=type"`
	// Level is SELinux level label that applies to the container.
	// +optional
	Level string `json:"level,omitempty" protobuf:"bytes,4,opt,name=level"`
}

// WindowsSecurityContextOptions contain Windows-specific options and credentials.
type WindowsSecurityContextOptions struct {
	// GMSACredentialSpecName is the name of the GMSA credential spec to use.
	// +optional
	GMSACredentialSpecName *string `json:"gmsaCredentialSpecName,omitempty" protobuf:"bytes,1,opt,name=gmsaCredentialSpecName"`

	// GMSACredentialSpec is where the GMSA admission webhook
	// (https://github.com/kubernetes-sigs/windows-gmsa) inlines the contents of the
	// GMSA credential spec named by the GMSACredentialSpecName field.
	// +optional
	GMSACredentialSpec *string `json:"gmsaCredentialSpec,omitempty" protobuf:"bytes,2,opt,name=gmsaCredentialSpec"`

	// The UserName in Windows to run the entrypoint of the container process.
	// Defaults to the user specified in image metadata if unspecified.
	// May also be set in PodSecurityContext. If set in both SecurityContext and
	// PodSecurityContext, the value specified in SecurityContext takes precedence.
	// +optional
	RunAsUserName *string `json:"runAsUserName,omitempty" protobuf:"bytes,3,opt,name=runAsUserName"`

	// HostProcess determines if a container should be run as a 'Host Process' container.
	// This field is alpha-level and will only be honored by components that enable the
	// WindowsHostProcessContainers feature flag. Setting this field without the feature
	// flag will result in errors when validating the Pod. All of a Pod's containers must
	// have the same effective HostProcess value (it is not allowed to have a mix of HostProcess
	// containers and non-HostProcess containers).  In addition, if HostProcess is true
	// then HostNetwork must also be set to true.
	// +optional
	HostProcess *bool `json:"hostProcess,omitempty" protobuf:"bytes,4,opt,name=hostProcess"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// RangeAllocation is not a public type.
type RangeAllocation struct {
	metav1.TypeMeta `json:",inline"`
	// Standard object's metadata.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty" protobuf:"bytes,1,opt,name=metadata"`

	// Range is string that identifies the range represented by 'data'.
	Range string `json:"range" protobuf:"bytes,2,opt,name=range"`
	// Data is a bit array containing all allocated addresses in the previous segment.
	Data []byte `json:"data" protobuf:"bytes,3,opt,name=data"`
}

const (
	// DefaultSchedulerName defines the name of default scheduler.
	DefaultSchedulerName = "default-scheduler"

	// RequiredDuringScheduling affinity is not symmetric, but there is an implicit PreferredDuringScheduling affinity rule
	// corresponding to every RequiredDuringScheduling affinity rule.
	// When the --hard-pod-affinity-weight scheduler flag is not specified,
	// DefaultHardPodAffinityWeight defines the weight of the implicit PreferredDuringScheduling affinity rule.
	DefaultHardPodAffinitySymmetricWeight int32 = 1
)

// Sysctl defines a kernel parameter to be set
type Sysctl struct {
	// Name of a property to set
	Name string `json:"name" protobuf:"bytes,1,opt,name=name"`
	// Value of a property to set
	Value string `json:"value" protobuf:"bytes,2,opt,name=value"`
}

// NodeResources is an object for conveying resource information about a node.
// see https://kubernetes.io/docs/concepts/architecture/nodes/#capacity for more details.
type NodeResources struct {
	// Capacity represents the available resources of a node
	Capacity ResourceList `protobuf:"bytes,1,rep,name=capacity,casttype=ResourceList,castkey=ResourceName"`
}

const (
	// Enable stdin for remote command execution
	ExecStdinParam = "input"
	// Enable stdout for remote command execution
	ExecStdoutParam = "output"
	// Enable stderr for remote command execution
	ExecStderrParam = "error"
	// Enable TTY for remote command execution
	ExecTTYParam = "tty"
	// Command to run for remote command execution
	ExecCommandParam = "command"

	// Name of header that specifies stream type
	StreamType = "streamType"
	// Value for streamType header for stdin stream
	StreamTypeStdin = "stdin"
	// Value for streamType header for stdout stream
	StreamTypeStdout = "stdout"
	// Value for streamType header for stderr stream
	StreamTypeStderr = "stderr"
	// Value for streamType header for data stream
	StreamTypeData = "data"
	// Value for streamType header for error stream
	StreamTypeError = "error"
	// Value for streamType header for terminal resize stream
	StreamTypeResize = "resize"

	// Name of header that specifies the port being forwarded
	PortHeader = "port"
	// Name of header that specifies a request ID used to associate the error
	// and data streams for a single forwarded connection
	PortForwardRequestIDHeader = "requestID"
)

// PortStatus represents the error condition of a service port

type PortStatus struct {
	// Port is the port number of the service port of which status is recorded here
	Port int32 `json:"port" protobuf:"varint,1,opt,name=port"`
	// Protocol is the protocol of the service port of which status is recorded here
	// The supported values are: "TCP", "UDP", "SCTP"
	Protocol Protocol `json:"protocol" protobuf:"bytes,2,opt,name=protocol,casttype=Protocol"`
	// Error is to record the problem with the service port
	// The format of the error shall comply with the following rules:
	// - built-in error values shall be specified in this file and those shall use
	//   CamelCase names
	// - cloud provider specific error values must have names that comply with the
	//   format foo.example.com/CamelCase.
	// ---
	// The regex it matches is (dns1123SubdomainFmt/)?(qualifiedNameFmt)
	// +optional
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Pattern=`^([a-z0-9]([-a-z0-9]*[a-z0-9])?(\.[a-z0-9]([-a-z0-9]*[a-z0-9])?)*/)?(([A-Za-z0-9][-A-Za-z0-9_.]*)?[A-Za-z0-9])$`
	// +kubebuilder:validation:MaxLength=316
	Error *string `json:"error,omitempty" protobuf:"bytes,3,opt,name=error"`
}
