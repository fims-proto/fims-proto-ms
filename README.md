# Backup for tenant module

为了架构的简洁和一些未解决的安全性问题, 决定暂时不支持多租户.

Golang runtime 不会有太大的 footprint, 对于多租户的需求, 选择使用多部署来实现.

假如今后有需要一份部署支持多租户的需求, 可以借鉴这个 branch 中的代码.