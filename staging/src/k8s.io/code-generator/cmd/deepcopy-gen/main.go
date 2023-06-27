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

// deepcopy-gen is a tool for auto-generating DeepCopy functions.
// deepcopy-gen是一个自动生成DeepCopy函数的工具。
//
// Given a list of input directories, it will generate functions that
// efficiently perform a full deep-copy of each type.  For any type that
// offers a `.DeepCopy()` method, it will simply call that.  Otherwise it will
// use standard value assignment whenever possible.  If that is not possible it
// will try to call its own generated copy function for the type, if the type is
// within the allowed root packages.  Failing that, it will fall back on
// `conversion.Cloner.DeepCopy(val)` to make the copy.  The resulting file will
// be stored in the same directory as the processed source package.
//
// 给定一个输入目录的列表，它将生成函数，这些函数可以有效地执行每种类型的完全深度复制。
// 对于任何提供`.DeepCopy()`方法的类型，它将简单地调用该方法。
// 否则，它将尽可能使用标准值分配。 如果不可能，它将尝试调用其自己生成的副本函数，如果类型在允许的根包中。
// 如果失败，它将回退到`conversion.Cloner.DeepCopy(val)`来进行复制。
// 生成的文件将存储在处理的源包目录中。
//
// Generation is governed by comment tags in the source.  Any package may
// request DeepCopy generation by including a comment in the file-comments of
// one file, of the form:
//
// 生成由源中的注释标记控制。 任何包都可以通过在一个文件的文件注释中包含一个注释来请求DeepCopy生成，形式如下：
//
//	// +k8s:deepcopy-gen=package
//
// DeepCopy functions can be generated for individual types, rather than the
// entire package by specifying a comment on the type definion of the form:
//
// DeepCopy函数可以为单个类型生成，而不是整个包，方法是在类型定义上指定一个注释，形式如下：
//
//	// +k8s:deepcopy-gen=true
//
// When generating for a whole package, individual types may opt out of
// DeepCopy generation by specifying a comment on the of the form:
//
// 当为整个包生成时，单个类型可以通过在形式上指定注释来选择退出DeepCopy生成：
//
//	// +k8s:deepcopy-gen=false
//
// Note that registration is a whole-package option, and is not available for
// individual types.
// 注意，注册是整个包的选项，不适用于单个类型。
package main

import (
	"flag"

	"github.com/spf13/pflag"
	"k8s.io/gengo/examples/deepcopy-gen/generators"
	"k8s.io/klog/v2"

	generatorargs "k8s.io/code-generator/cmd/deepcopy-gen/args"
)

func main() {
	klog.InitFlags(nil)
	genericArgs, customArgs := generatorargs.NewDefaults()

	genericArgs.AddFlags(pflag.CommandLine)
	customArgs.AddFlags(pflag.CommandLine)
	flag.Set("logtostderr", "true")
	pflag.CommandLine.AddGoFlagSet(flag.CommandLine)
	pflag.Parse()

	if err := generatorargs.Validate(genericArgs); err != nil {
		klog.Fatalf("Error: %v", err)
	}

	// Run it.
	if err := genericArgs.Execute(
		generators.NameSystems(),
		generators.DefaultNameSystem(),
		generators.Packages,
	); err != nil {
		klog.Fatalf("Error: %v", err)
	}
	klog.V(2).Info("Completed successfully.")
}
