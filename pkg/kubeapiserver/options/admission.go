/*
Copyright 2018 The Kubernetes Authors.

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

package options

import (
	"fmt"
	"strings"

	"github.com/spf13/pflag"

	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/apiserver/pkg/admission"
	"k8s.io/apiserver/pkg/server"
	genericoptions "k8s.io/apiserver/pkg/server/options"
	"k8s.io/client-go/informers"
	"k8s.io/client-go/rest"
	"k8s.io/component-base/featuregate"
)

// AdmissionOptions holds the admission options.
// It is a wrap of generic AdmissionOptions.
// AdmissionOptions 处理 admission 控制器的选项, 它是通用 AdmissionOptions 的封装
type AdmissionOptions struct {
	// GenericAdmission holds the generic admission options.
	// GenericAdmission 处理通用 admission 控制器的选项
	GenericAdmission *genericoptions.AdmissionOptions
	// DEPRECATED flag, should use EnabledAdmissionPlugins and DisabledAdmissionPlugins.
	// They are mutually exclusive, specify both will lead to an error.
	PluginNames []string
}

// NewAdmissionOptions creates a new instance of AdmissionOptions
// NewAdmissionOptions 创建一个 AdmissionOptions 实例
// Note:
//
//		In addition it calls RegisterAllAdmissionPlugins to register
//		all kube-apiserver admission plugins.
//	 除此之外, 它还调用了 RegisterAllAdmissionPlugins 来注册所有 kube-apiserver admission 控制器
//
//		Provides the list of RecommendedPluginOrder that holds sane values
//		that can be used by servers that don't care about admission chain.
//		Servers that do care can overwrite/append that field after creation.
//	 提供了一个 RecommendedPluginOrder 列表, 该列表包含了一些可以被不关心 admission 链的服务器使用的值
//	 关心的服务器可以在创建之后覆盖/追加该字段
func NewAdmissionOptions() *AdmissionOptions {
	options := genericoptions.NewAdmissionOptions()
	// register all admission plugins
	RegisterAllAdmissionPlugins(options.Plugins)
	// set RecommendedPluginOrder
	options.RecommendedPluginOrder = AllOrderedPlugins
	// set DefaultOffPlugins
	options.DefaultOffPlugins = DefaultOffAdmissionPlugins()

	return &AdmissionOptions{
		GenericAdmission: options,
	}
}

// AddFlags adds flags related to admission for kube-apiserver to the specified FlagSet
// AddFlags 将与 kube-apiserver 相关的 admission 控制器的 flag 添加到指定的 FlagSet 中
func (a *AdmissionOptions) AddFlags(fs *pflag.FlagSet) {
	fs.StringSliceVar(&a.PluginNames, "admission-control", a.PluginNames, ""+
		"Admission is divided into two phases. "+
		"In the first phase, only mutating admission plugins run. "+
		"In the second phase, only validating admission plugins run. "+
		"The names in the below list may represent a validating plugin, a mutating plugin, or both. "+
		"The order of plugins in which they are passed to this flag does not matter. "+
		"Comma-delimited list of: "+strings.Join(a.GenericAdmission.Plugins.Registered(), ", ")+".")
	fs.MarkDeprecated("admission-control", "Use --enable-admission-plugins or --disable-admission-plugins instead. Will be removed in a future version.")
	fs.Lookup("admission-control").Hidden = false

	a.GenericAdmission.AddFlags(fs)
}

// Validate verifies flags passed to kube-apiserver AdmissionOptions.
// Kube-apiserver verifies PluginNames and then call generic AdmissionOptions.Validate.
func (a *AdmissionOptions) Validate() []error {
	if a == nil {
		return nil
	}
	var errs []error
	if a.PluginNames != nil &&
		(a.GenericAdmission.EnablePlugins != nil || a.GenericAdmission.DisablePlugins != nil) {
		errs = append(errs, fmt.Errorf("admission-control and enable-admission-plugins/disable-admission-plugins flags are mutually exclusive"))
	}

	registeredPlugins := sets.NewString(a.GenericAdmission.Plugins.Registered()...)
	for _, name := range a.PluginNames {
		if !registeredPlugins.Has(name) {
			errs = append(errs, fmt.Errorf("admission-control plugin %q is unknown", name))
		}
	}

	errs = append(errs, a.GenericAdmission.Validate()...)

	return errs
}

// ApplyTo adds the admission chain to the server configuration.
// Kube-apiserver just call generic AdmissionOptions.ApplyTo.
// ApplyTo 将 admission 链添加到服务器配置中
func (a *AdmissionOptions) ApplyTo(
	c *server.Config,
	informers informers.SharedInformerFactory,
	kubeAPIServerClientConfig *rest.Config,
	features featuregate.FeatureGate,
	pluginInitializers ...admission.PluginInitializer,
) error {
	if a == nil {
		return nil
	}

	if a.PluginNames != nil {
		// pass PluginNames to generic AdmissionOptions
		a.GenericAdmission.EnablePlugins, a.GenericAdmission.DisablePlugins = computePluginNames(a.PluginNames, a.GenericAdmission.RecommendedPluginOrder)
	}

	return a.GenericAdmission.ApplyTo(c, informers, kubeAPIServerClientConfig, features, pluginInitializers...)
}

// explicitly disable all plugins that are not in the enabled list
func computePluginNames(explicitlyEnabled []string, all []string) (enabled []string, disabled []string) {
	return explicitlyEnabled, sets.NewString(all...).Difference(sets.NewString(explicitlyEnabled...)).List()
}
