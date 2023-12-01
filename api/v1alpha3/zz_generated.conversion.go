//go:build !ignore_autogenerated
// +build !ignore_autogenerated

/*
Copyright 2021.

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
// Code generated by conversion-gen. DO NOT EDIT.

package v1alpha3

import (
	unsafe "unsafe"

	unversioned "github.com/eraser-dev/eraser/api/unversioned"
	conversion "k8s.io/apimachinery/pkg/conversion"
	runtime "k8s.io/apimachinery/pkg/runtime"
)

func init() {
	localSchemeBuilder.Register(RegisterConversions)
}

// RegisterConversions adds conversion functions to the given scheme.
// Public to allow building arbitrary schemes.
func RegisterConversions(s *runtime.Scheme) error {
	if err := s.AddGeneratedConversionFunc((*Components)(nil), (*unversioned.Components)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha3_Components_To_unversioned_Components(a.(*Components), b.(*unversioned.Components), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*unversioned.Components)(nil), (*Components)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_unversioned_Components_To_v1alpha3_Components(a.(*unversioned.Components), b.(*Components), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*ContainerConfig)(nil), (*unversioned.ContainerConfig)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha3_ContainerConfig_To_unversioned_ContainerConfig(a.(*ContainerConfig), b.(*unversioned.ContainerConfig), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*unversioned.ContainerConfig)(nil), (*ContainerConfig)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_unversioned_ContainerConfig_To_v1alpha3_ContainerConfig(a.(*unversioned.ContainerConfig), b.(*ContainerConfig), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*EraserConfig)(nil), (*unversioned.EraserConfig)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha3_EraserConfig_To_unversioned_EraserConfig(a.(*EraserConfig), b.(*unversioned.EraserConfig), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*unversioned.EraserConfig)(nil), (*EraserConfig)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_unversioned_EraserConfig_To_v1alpha3_EraserConfig(a.(*unversioned.EraserConfig), b.(*EraserConfig), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*ImageJobCleanupConfig)(nil), (*unversioned.ImageJobCleanupConfig)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha3_ImageJobCleanupConfig_To_unversioned_ImageJobCleanupConfig(a.(*ImageJobCleanupConfig), b.(*unversioned.ImageJobCleanupConfig), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*unversioned.ImageJobCleanupConfig)(nil), (*ImageJobCleanupConfig)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_unversioned_ImageJobCleanupConfig_To_v1alpha3_ImageJobCleanupConfig(a.(*unversioned.ImageJobCleanupConfig), b.(*ImageJobCleanupConfig), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*ImageJobConfig)(nil), (*unversioned.ImageJobConfig)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha3_ImageJobConfig_To_unversioned_ImageJobConfig(a.(*ImageJobConfig), b.(*unversioned.ImageJobConfig), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*unversioned.ImageJobConfig)(nil), (*ImageJobConfig)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_unversioned_ImageJobConfig_To_v1alpha3_ImageJobConfig(a.(*unversioned.ImageJobConfig), b.(*ImageJobConfig), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*ManagerConfig)(nil), (*unversioned.ManagerConfig)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha3_ManagerConfig_To_unversioned_ManagerConfig(a.(*ManagerConfig), b.(*unversioned.ManagerConfig), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*unversioned.ManagerConfig)(nil), (*ManagerConfig)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_unversioned_ManagerConfig_To_v1alpha3_ManagerConfig(a.(*unversioned.ManagerConfig), b.(*ManagerConfig), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*NodeFilterConfig)(nil), (*unversioned.NodeFilterConfig)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha3_NodeFilterConfig_To_unversioned_NodeFilterConfig(a.(*NodeFilterConfig), b.(*unversioned.NodeFilterConfig), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*unversioned.NodeFilterConfig)(nil), (*NodeFilterConfig)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_unversioned_NodeFilterConfig_To_v1alpha3_NodeFilterConfig(a.(*unversioned.NodeFilterConfig), b.(*NodeFilterConfig), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*OptionalContainerConfig)(nil), (*unversioned.OptionalContainerConfig)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha3_OptionalContainerConfig_To_unversioned_OptionalContainerConfig(a.(*OptionalContainerConfig), b.(*unversioned.OptionalContainerConfig), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*unversioned.OptionalContainerConfig)(nil), (*OptionalContainerConfig)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_unversioned_OptionalContainerConfig_To_v1alpha3_OptionalContainerConfig(a.(*unversioned.OptionalContainerConfig), b.(*OptionalContainerConfig), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*ProfileConfig)(nil), (*unversioned.ProfileConfig)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha3_ProfileConfig_To_unversioned_ProfileConfig(a.(*ProfileConfig), b.(*unversioned.ProfileConfig), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*unversioned.ProfileConfig)(nil), (*ProfileConfig)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_unversioned_ProfileConfig_To_v1alpha3_ProfileConfig(a.(*unversioned.ProfileConfig), b.(*ProfileConfig), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*RepoTag)(nil), (*unversioned.RepoTag)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha3_RepoTag_To_unversioned_RepoTag(a.(*RepoTag), b.(*unversioned.RepoTag), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*unversioned.RepoTag)(nil), (*RepoTag)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_unversioned_RepoTag_To_v1alpha3_RepoTag(a.(*unversioned.RepoTag), b.(*RepoTag), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*ResourceRequirements)(nil), (*unversioned.ResourceRequirements)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha3_ResourceRequirements_To_unversioned_ResourceRequirements(a.(*ResourceRequirements), b.(*unversioned.ResourceRequirements), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*unversioned.ResourceRequirements)(nil), (*ResourceRequirements)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_unversioned_ResourceRequirements_To_v1alpha3_ResourceRequirements(a.(*unversioned.ResourceRequirements), b.(*ResourceRequirements), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*ScheduleConfig)(nil), (*unversioned.ScheduleConfig)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_v1alpha3_ScheduleConfig_To_unversioned_ScheduleConfig(a.(*ScheduleConfig), b.(*unversioned.ScheduleConfig), scope)
	}); err != nil {
		return err
	}
	if err := s.AddGeneratedConversionFunc((*unversioned.ScheduleConfig)(nil), (*ScheduleConfig)(nil), func(a, b interface{}, scope conversion.Scope) error {
		return Convert_unversioned_ScheduleConfig_To_v1alpha3_ScheduleConfig(a.(*unversioned.ScheduleConfig), b.(*ScheduleConfig), scope)
	}); err != nil {
		return err
	}
	return nil
}

func autoConvert_v1alpha3_Components_To_unversioned_Components(in *Components, out *unversioned.Components, s conversion.Scope) error {
	if err := Convert_v1alpha3_OptionalContainerConfig_To_unversioned_OptionalContainerConfig(&in.Collector, &out.Collector, s); err != nil {
		return err
	}
	if err := Convert_v1alpha3_OptionalContainerConfig_To_unversioned_OptionalContainerConfig(&in.Scanner, &out.Scanner, s); err != nil {
		return err
	}
	if err := Convert_v1alpha3_ContainerConfig_To_unversioned_ContainerConfig(&in.Remover, &out.Remover, s); err != nil {
		return err
	}
	return nil
}

// Convert_v1alpha3_Components_To_unversioned_Components is an autogenerated conversion function.
func Convert_v1alpha3_Components_To_unversioned_Components(in *Components, out *unversioned.Components, s conversion.Scope) error {
	return autoConvert_v1alpha3_Components_To_unversioned_Components(in, out, s)
}

func autoConvert_unversioned_Components_To_v1alpha3_Components(in *unversioned.Components, out *Components, s conversion.Scope) error {
	if err := Convert_unversioned_OptionalContainerConfig_To_v1alpha3_OptionalContainerConfig(&in.Collector, &out.Collector, s); err != nil {
		return err
	}
	if err := Convert_unversioned_OptionalContainerConfig_To_v1alpha3_OptionalContainerConfig(&in.Scanner, &out.Scanner, s); err != nil {
		return err
	}
	if err := Convert_unversioned_ContainerConfig_To_v1alpha3_ContainerConfig(&in.Remover, &out.Remover, s); err != nil {
		return err
	}
	return nil
}

// Convert_unversioned_Components_To_v1alpha3_Components is an autogenerated conversion function.
func Convert_unversioned_Components_To_v1alpha3_Components(in *unversioned.Components, out *Components, s conversion.Scope) error {
	return autoConvert_unversioned_Components_To_v1alpha3_Components(in, out, s)
}

func autoConvert_v1alpha3_ContainerConfig_To_unversioned_ContainerConfig(in *ContainerConfig, out *unversioned.ContainerConfig, s conversion.Scope) error {
	if err := Convert_v1alpha3_RepoTag_To_unversioned_RepoTag(&in.Image, &out.Image, s); err != nil {
		return err
	}
	if err := Convert_v1alpha3_ResourceRequirements_To_unversioned_ResourceRequirements(&in.Request, &out.Request, s); err != nil {
		return err
	}
	if err := Convert_v1alpha3_ResourceRequirements_To_unversioned_ResourceRequirements(&in.Limit, &out.Limit, s); err != nil {
		return err
	}
	out.Config = (*string)(unsafe.Pointer(in.Config))
	return nil
}

// Convert_v1alpha3_ContainerConfig_To_unversioned_ContainerConfig is an autogenerated conversion function.
func Convert_v1alpha3_ContainerConfig_To_unversioned_ContainerConfig(in *ContainerConfig, out *unversioned.ContainerConfig, s conversion.Scope) error {
	return autoConvert_v1alpha3_ContainerConfig_To_unversioned_ContainerConfig(in, out, s)
}

func autoConvert_unversioned_ContainerConfig_To_v1alpha3_ContainerConfig(in *unversioned.ContainerConfig, out *ContainerConfig, s conversion.Scope) error {
	if err := Convert_unversioned_RepoTag_To_v1alpha3_RepoTag(&in.Image, &out.Image, s); err != nil {
		return err
	}
	if err := Convert_unversioned_ResourceRequirements_To_v1alpha3_ResourceRequirements(&in.Request, &out.Request, s); err != nil {
		return err
	}
	if err := Convert_unversioned_ResourceRequirements_To_v1alpha3_ResourceRequirements(&in.Limit, &out.Limit, s); err != nil {
		return err
	}
	out.Config = (*string)(unsafe.Pointer(in.Config))
	return nil
}

// Convert_unversioned_ContainerConfig_To_v1alpha3_ContainerConfig is an autogenerated conversion function.
func Convert_unversioned_ContainerConfig_To_v1alpha3_ContainerConfig(in *unversioned.ContainerConfig, out *ContainerConfig, s conversion.Scope) error {
	return autoConvert_unversioned_ContainerConfig_To_v1alpha3_ContainerConfig(in, out, s)
}

func autoConvert_v1alpha3_EraserConfig_To_unversioned_EraserConfig(in *EraserConfig, out *unversioned.EraserConfig, s conversion.Scope) error {
	if err := Convert_v1alpha3_ManagerConfig_To_unversioned_ManagerConfig(&in.Manager, &out.Manager, s); err != nil {
		return err
	}
	if err := Convert_v1alpha3_Components_To_unversioned_Components(&in.Components, &out.Components, s); err != nil {
		return err
	}
	return nil
}

// Convert_v1alpha3_EraserConfig_To_unversioned_EraserConfig is an autogenerated conversion function.
func Convert_v1alpha3_EraserConfig_To_unversioned_EraserConfig(in *EraserConfig, out *unversioned.EraserConfig, s conversion.Scope) error {
	return autoConvert_v1alpha3_EraserConfig_To_unversioned_EraserConfig(in, out, s)
}

func autoConvert_unversioned_EraserConfig_To_v1alpha3_EraserConfig(in *unversioned.EraserConfig, out *EraserConfig, s conversion.Scope) error {
	if err := Convert_unversioned_ManagerConfig_To_v1alpha3_ManagerConfig(&in.Manager, &out.Manager, s); err != nil {
		return err
	}
	if err := Convert_unversioned_Components_To_v1alpha3_Components(&in.Components, &out.Components, s); err != nil {
		return err
	}
	return nil
}

// Convert_unversioned_EraserConfig_To_v1alpha3_EraserConfig is an autogenerated conversion function.
func Convert_unversioned_EraserConfig_To_v1alpha3_EraserConfig(in *unversioned.EraserConfig, out *EraserConfig, s conversion.Scope) error {
	return autoConvert_unversioned_EraserConfig_To_v1alpha3_EraserConfig(in, out, s)
}

func autoConvert_v1alpha3_ImageJobCleanupConfig_To_unversioned_ImageJobCleanupConfig(in *ImageJobCleanupConfig, out *unversioned.ImageJobCleanupConfig, s conversion.Scope) error {
	out.DelayOnSuccess = unversioned.Duration(in.DelayOnSuccess)
	out.DelayOnFailure = unversioned.Duration(in.DelayOnFailure)
	return nil
}

// Convert_v1alpha3_ImageJobCleanupConfig_To_unversioned_ImageJobCleanupConfig is an autogenerated conversion function.
func Convert_v1alpha3_ImageJobCleanupConfig_To_unversioned_ImageJobCleanupConfig(in *ImageJobCleanupConfig, out *unversioned.ImageJobCleanupConfig, s conversion.Scope) error {
	return autoConvert_v1alpha3_ImageJobCleanupConfig_To_unversioned_ImageJobCleanupConfig(in, out, s)
}

func autoConvert_unversioned_ImageJobCleanupConfig_To_v1alpha3_ImageJobCleanupConfig(in *unversioned.ImageJobCleanupConfig, out *ImageJobCleanupConfig, s conversion.Scope) error {
	out.DelayOnSuccess = Duration(in.DelayOnSuccess)
	out.DelayOnFailure = Duration(in.DelayOnFailure)
	return nil
}

// Convert_unversioned_ImageJobCleanupConfig_To_v1alpha3_ImageJobCleanupConfig is an autogenerated conversion function.
func Convert_unversioned_ImageJobCleanupConfig_To_v1alpha3_ImageJobCleanupConfig(in *unversioned.ImageJobCleanupConfig, out *ImageJobCleanupConfig, s conversion.Scope) error {
	return autoConvert_unversioned_ImageJobCleanupConfig_To_v1alpha3_ImageJobCleanupConfig(in, out, s)
}

func autoConvert_v1alpha3_ImageJobConfig_To_unversioned_ImageJobConfig(in *ImageJobConfig, out *unversioned.ImageJobConfig, s conversion.Scope) error {
	out.SuccessRatio = in.SuccessRatio
	if err := Convert_v1alpha3_ImageJobCleanupConfig_To_unversioned_ImageJobCleanupConfig(&in.Cleanup, &out.Cleanup, s); err != nil {
		return err
	}
	return nil
}

// Convert_v1alpha3_ImageJobConfig_To_unversioned_ImageJobConfig is an autogenerated conversion function.
func Convert_v1alpha3_ImageJobConfig_To_unversioned_ImageJobConfig(in *ImageJobConfig, out *unversioned.ImageJobConfig, s conversion.Scope) error {
	return autoConvert_v1alpha3_ImageJobConfig_To_unversioned_ImageJobConfig(in, out, s)
}

func autoConvert_unversioned_ImageJobConfig_To_v1alpha3_ImageJobConfig(in *unversioned.ImageJobConfig, out *ImageJobConfig, s conversion.Scope) error {
	out.SuccessRatio = in.SuccessRatio
	if err := Convert_unversioned_ImageJobCleanupConfig_To_v1alpha3_ImageJobCleanupConfig(&in.Cleanup, &out.Cleanup, s); err != nil {
		return err
	}
	return nil
}

// Convert_unversioned_ImageJobConfig_To_v1alpha3_ImageJobConfig is an autogenerated conversion function.
func Convert_unversioned_ImageJobConfig_To_v1alpha3_ImageJobConfig(in *unversioned.ImageJobConfig, out *ImageJobConfig, s conversion.Scope) error {
	return autoConvert_unversioned_ImageJobConfig_To_v1alpha3_ImageJobConfig(in, out, s)
}

func autoConvert_v1alpha3_ManagerConfig_To_unversioned_ManagerConfig(in *ManagerConfig, out *unversioned.ManagerConfig, s conversion.Scope) error {
	// WARNING: in.Runtime requires manual conversion: inconvertible types (./api/v1alpha3.RuntimeSpec vs github.com/eraser-dev/eraser/api/unversioned.Runtime)
	if err := Convert_v1alpha3_RuntimeSpec_To_unversioned_Runtime(&in.Runtime, &out.Runtime, s); err != nil {
		return err
	}
	out.OTLPEndpoint = in.OTLPEndpoint
	out.LogLevel = in.LogLevel
	if err := Convert_v1alpha3_ScheduleConfig_To_unversioned_ScheduleConfig(&in.Scheduling, &out.Scheduling, s); err != nil {
		return err
	}
	if err := Convert_v1alpha3_ProfileConfig_To_unversioned_ProfileConfig(&in.Profile, &out.Profile, s); err != nil {
		return err
	}
	if err := Convert_v1alpha3_ImageJobConfig_To_unversioned_ImageJobConfig(&in.ImageJob, &out.ImageJob, s); err != nil {
		return err
	}
	out.PullSecrets = *(*[]string)(unsafe.Pointer(&in.PullSecrets))
	if err := Convert_v1alpha3_NodeFilterConfig_To_unversioned_NodeFilterConfig(&in.NodeFilter, &out.NodeFilter, s); err != nil {
		return err
	}
	out.PriorityClassName = in.PriorityClassName
	return nil
}

func autoConvert_unversioned_ManagerConfig_To_v1alpha3_ManagerConfig(in *unversioned.ManagerConfig, out *ManagerConfig, s conversion.Scope) error {
	// WARNING: in.Runtime requires manual conversion: inconvertible types (github.com/eraser-dev/eraser/api/unversioned.Runtime vs ./api/v1alpha3.RuntimeSpec)
	if err := Convert_unversioned_Runtime_To_v1alpha3_RuntimeSpec(&in.Runtime, &out.Runtime, s); err != nil {
		return err
	}
	out.OTLPEndpoint = in.OTLPEndpoint
	out.LogLevel = in.LogLevel
	if err := Convert_unversioned_ScheduleConfig_To_v1alpha3_ScheduleConfig(&in.Scheduling, &out.Scheduling, s); err != nil {
		return err
	}
	if err := Convert_unversioned_ProfileConfig_To_v1alpha3_ProfileConfig(&in.Profile, &out.Profile, s); err != nil {
		return err
	}
	if err := Convert_unversioned_ImageJobConfig_To_v1alpha3_ImageJobConfig(&in.ImageJob, &out.ImageJob, s); err != nil {
		return err
	}
	out.PullSecrets = *(*[]string)(unsafe.Pointer(&in.PullSecrets))
	if err := Convert_unversioned_NodeFilterConfig_To_v1alpha3_NodeFilterConfig(&in.NodeFilter, &out.NodeFilter, s); err != nil {
		return err
	}
	out.PriorityClassName = in.PriorityClassName
	return nil
}

func autoConvert_v1alpha3_NodeFilterConfig_To_unversioned_NodeFilterConfig(in *NodeFilterConfig, out *unversioned.NodeFilterConfig, s conversion.Scope) error {
	out.Type = in.Type
	out.Selectors = *(*[]string)(unsafe.Pointer(&in.Selectors))
	return nil
}

// Convert_v1alpha3_NodeFilterConfig_To_unversioned_NodeFilterConfig is an autogenerated conversion function.
func Convert_v1alpha3_NodeFilterConfig_To_unversioned_NodeFilterConfig(in *NodeFilterConfig, out *unversioned.NodeFilterConfig, s conversion.Scope) error {
	return autoConvert_v1alpha3_NodeFilterConfig_To_unversioned_NodeFilterConfig(in, out, s)
}

func autoConvert_unversioned_NodeFilterConfig_To_v1alpha3_NodeFilterConfig(in *unversioned.NodeFilterConfig, out *NodeFilterConfig, s conversion.Scope) error {
	out.Type = in.Type
	out.Selectors = *(*[]string)(unsafe.Pointer(&in.Selectors))
	return nil
}

// Convert_unversioned_NodeFilterConfig_To_v1alpha3_NodeFilterConfig is an autogenerated conversion function.
func Convert_unversioned_NodeFilterConfig_To_v1alpha3_NodeFilterConfig(in *unversioned.NodeFilterConfig, out *NodeFilterConfig, s conversion.Scope) error {
	return autoConvert_unversioned_NodeFilterConfig_To_v1alpha3_NodeFilterConfig(in, out, s)
}

func autoConvert_v1alpha3_OptionalContainerConfig_To_unversioned_OptionalContainerConfig(in *OptionalContainerConfig, out *unversioned.OptionalContainerConfig, s conversion.Scope) error {
	out.Enabled = in.Enabled
	if err := Convert_v1alpha3_ContainerConfig_To_unversioned_ContainerConfig(&in.ContainerConfig, &out.ContainerConfig, s); err != nil {
		return err
	}
	return nil
}

// Convert_v1alpha3_OptionalContainerConfig_To_unversioned_OptionalContainerConfig is an autogenerated conversion function.
func Convert_v1alpha3_OptionalContainerConfig_To_unversioned_OptionalContainerConfig(in *OptionalContainerConfig, out *unversioned.OptionalContainerConfig, s conversion.Scope) error {
	return autoConvert_v1alpha3_OptionalContainerConfig_To_unversioned_OptionalContainerConfig(in, out, s)
}

func autoConvert_unversioned_OptionalContainerConfig_To_v1alpha3_OptionalContainerConfig(in *unversioned.OptionalContainerConfig, out *OptionalContainerConfig, s conversion.Scope) error {
	out.Enabled = in.Enabled
	if err := Convert_unversioned_ContainerConfig_To_v1alpha3_ContainerConfig(&in.ContainerConfig, &out.ContainerConfig, s); err != nil {
		return err
	}
	return nil
}

// Convert_unversioned_OptionalContainerConfig_To_v1alpha3_OptionalContainerConfig is an autogenerated conversion function.
func Convert_unversioned_OptionalContainerConfig_To_v1alpha3_OptionalContainerConfig(in *unversioned.OptionalContainerConfig, out *OptionalContainerConfig, s conversion.Scope) error {
	return autoConvert_unversioned_OptionalContainerConfig_To_v1alpha3_OptionalContainerConfig(in, out, s)
}

func autoConvert_v1alpha3_ProfileConfig_To_unversioned_ProfileConfig(in *ProfileConfig, out *unversioned.ProfileConfig, s conversion.Scope) error {
	out.Enabled = in.Enabled
	out.Port = in.Port
	return nil
}

// Convert_v1alpha3_ProfileConfig_To_unversioned_ProfileConfig is an autogenerated conversion function.
func Convert_v1alpha3_ProfileConfig_To_unversioned_ProfileConfig(in *ProfileConfig, out *unversioned.ProfileConfig, s conversion.Scope) error {
	return autoConvert_v1alpha3_ProfileConfig_To_unversioned_ProfileConfig(in, out, s)
}

func autoConvert_unversioned_ProfileConfig_To_v1alpha3_ProfileConfig(in *unversioned.ProfileConfig, out *ProfileConfig, s conversion.Scope) error {
	out.Enabled = in.Enabled
	out.Port = in.Port
	return nil
}

// Convert_unversioned_ProfileConfig_To_v1alpha3_ProfileConfig is an autogenerated conversion function.
func Convert_unversioned_ProfileConfig_To_v1alpha3_ProfileConfig(in *unversioned.ProfileConfig, out *ProfileConfig, s conversion.Scope) error {
	return autoConvert_unversioned_ProfileConfig_To_v1alpha3_ProfileConfig(in, out, s)
}

func autoConvert_v1alpha3_RepoTag_To_unversioned_RepoTag(in *RepoTag, out *unversioned.RepoTag, s conversion.Scope) error {
	out.Repo = in.Repo
	out.Tag = in.Tag
	return nil
}

// Convert_v1alpha3_RepoTag_To_unversioned_RepoTag is an autogenerated conversion function.
func Convert_v1alpha3_RepoTag_To_unversioned_RepoTag(in *RepoTag, out *unversioned.RepoTag, s conversion.Scope) error {
	return autoConvert_v1alpha3_RepoTag_To_unversioned_RepoTag(in, out, s)
}

func autoConvert_unversioned_RepoTag_To_v1alpha3_RepoTag(in *unversioned.RepoTag, out *RepoTag, s conversion.Scope) error {
	out.Repo = in.Repo
	out.Tag = in.Tag
	return nil
}

// Convert_unversioned_RepoTag_To_v1alpha3_RepoTag is an autogenerated conversion function.
func Convert_unversioned_RepoTag_To_v1alpha3_RepoTag(in *unversioned.RepoTag, out *RepoTag, s conversion.Scope) error {
	return autoConvert_unversioned_RepoTag_To_v1alpha3_RepoTag(in, out, s)
}

func autoConvert_v1alpha3_ResourceRequirements_To_unversioned_ResourceRequirements(in *ResourceRequirements, out *unversioned.ResourceRequirements, s conversion.Scope) error {
	out.Mem = in.Mem
	out.CPU = in.CPU
	return nil
}

// Convert_v1alpha3_ResourceRequirements_To_unversioned_ResourceRequirements is an autogenerated conversion function.
func Convert_v1alpha3_ResourceRequirements_To_unversioned_ResourceRequirements(in *ResourceRequirements, out *unversioned.ResourceRequirements, s conversion.Scope) error {
	return autoConvert_v1alpha3_ResourceRequirements_To_unversioned_ResourceRequirements(in, out, s)
}

func autoConvert_unversioned_ResourceRequirements_To_v1alpha3_ResourceRequirements(in *unversioned.ResourceRequirements, out *ResourceRequirements, s conversion.Scope) error {
	out.Mem = in.Mem
	out.CPU = in.CPU
	return nil
}

// Convert_unversioned_ResourceRequirements_To_v1alpha3_ResourceRequirements is an autogenerated conversion function.
func Convert_unversioned_ResourceRequirements_To_v1alpha3_ResourceRequirements(in *unversioned.ResourceRequirements, out *ResourceRequirements, s conversion.Scope) error {
	return autoConvert_unversioned_ResourceRequirements_To_v1alpha3_ResourceRequirements(in, out, s)
}

func autoConvert_v1alpha3_ScheduleConfig_To_unversioned_ScheduleConfig(in *ScheduleConfig, out *unversioned.ScheduleConfig, s conversion.Scope) error {
	out.RepeatInterval = unversioned.Duration(in.RepeatInterval)
	out.BeginImmediately = in.BeginImmediately
	return nil
}

// Convert_v1alpha3_ScheduleConfig_To_unversioned_ScheduleConfig is an autogenerated conversion function.
func Convert_v1alpha3_ScheduleConfig_To_unversioned_ScheduleConfig(in *ScheduleConfig, out *unversioned.ScheduleConfig, s conversion.Scope) error {
	return autoConvert_v1alpha3_ScheduleConfig_To_unversioned_ScheduleConfig(in, out, s)
}

func autoConvert_unversioned_ScheduleConfig_To_v1alpha3_ScheduleConfig(in *unversioned.ScheduleConfig, out *ScheduleConfig, s conversion.Scope) error {
	out.RepeatInterval = Duration(in.RepeatInterval)
	out.BeginImmediately = in.BeginImmediately
	return nil
}

// Convert_unversioned_ScheduleConfig_To_v1alpha3_ScheduleConfig is an autogenerated conversion function.
func Convert_unversioned_ScheduleConfig_To_v1alpha3_ScheduleConfig(in *unversioned.ScheduleConfig, out *ScheduleConfig, s conversion.Scope) error {
	return autoConvert_unversioned_ScheduleConfig_To_v1alpha3_ScheduleConfig(in, out, s)
}

// The following were added manually for ManagerConfig RuntimeSpec conversions:

func Convert_v1alpha3_ManagerConfig_To_unversioned_ManagerConfig(in *ManagerConfig, out *unversioned.ManagerConfig, s conversion.Scope) error{
	return autoConvert_v1alpha3_ManagerConfig_To_unversioned_ManagerConfig(in, out, s)
}

// TODO: change this to use unversioned.RuntimeSpec when unversioned is updated
func autoConvert_v1alpha3_RuntimeSpec_To_unversioned_Runtime(in *RuntimeSpec, out *unversioned.Runtime, s conversion.Scope) error {
	*out = unversioned.Runtime(in.Name)
	return nil
}

func Convert_v1alpha3_RuntimeSpec_To_unversioned_Runtime(in *RuntimeSpec, out *unversioned.Runtime, s conversion.Scope) error {
	return autoConvert_v1alpha3_RuntimeSpec_To_unversioned_Runtime(in, out, s)
}

func Convert_unversioned_ManagerConfig_To_v1alpha3_ManagerConfig(in *unversioned.ManagerConfig, out *ManagerConfig, s conversion.Scope) error{
	return autoConvert_unversioned_ManagerConfig_To_v1alpha3_ManagerConfig(in, out, s)
}

// TODO: change this to use unversioned.RuntimeSpec when unversioned is updated
func autoConvert_unversioned_Runtime_To_v1alpha3_RuntimeSpec(in *unversioned.Runtime, out *RuntimeSpec, s conversion.Scope) error {
	out.Name = Runtime(string(*in))
	out.Address = RuntimeAddress("")
	return nil
}

func Convert_unversioned_Runtime_To_v1alpha3_RuntimeSpec(in *unversioned.Runtime, out *RuntimeSpec, s conversion.Scope) error {
	return autoConvert_unversioned_Runtime_To_v1alpha3_RuntimeSpec(in, out, s)
}
