## DISCLAIMER
- Sig-Node社区已经达成了一个普遍的共识，作为一种最佳实践，避免引入任何新的检查点支持。
避免引入任何新的检查点支持。我们达成了这个共识
在与生产环境中一些难以调试的问题作斗争之后
导致的一些难以调试的问题后，我们达成了这个共识。
- 对检查点数据结构的任何改变都会被认为是不兼容的，如果一个组件需要确保阅读旧格式的检查点文件的向后兼容性，它应该添加自己的处理方法。

## Introduction
这个文件夹包含了一个框架和基元，即检查点管理器，它被其他几个Kubelet子模块所使用。
它被其他几个Kubelet子模块使用，`dockershim`，`devicemanager`，`pods`和`cpumanager`。
和 "cpumanager"，以实现每个子模块层面的检查点。正如
在上面的 "免责声明 "部分已经解释过，在Kubelet中引入任何进一步的
在Kubelet中引入任何进一步的检查点。如果仍然需要检查点，那么这个文件夹
提供了实现检查点的通用API和框架。
在所有的子模块中使用相同的API将有助于在
Kubelet级别的一致性。

Below is the history of checkpointing support in Kubelet.

| Package | First checkpointing support merged on | PR link |
| ------- | --------------------------------------| ------- |
|kubelet/dockershim | Feb 3, 2017 | [[CRI] Implement Dockershim Checkpoint](https://github.com/kubernetes/kubernetes/pull/39903)
|devicemanager| Sep 6, 2017 | [Deviceplugin checkpoint](https://github.com/kubernetes/kubernetes/pull/51744)
| kubelet/pod | Nov 22, 2017 | [Initial basic bootstrap-checkpoint support](https://github.com/kubernetes/kubernetes/pull/50984)
|cpumanager| Oct 27, 2017 |[Add file backed state to cpu manager ](https://github.com/kubernetes/kubernetes/pull/54408)
