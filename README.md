# PrivateBlockchain-sample-ByGo
[![license](https://img.shields.io/github/license/pure-admin/vue-pure-admin.svg)](LICENSE)

**中文** | [English](./README.en-US.md)
## 介绍
一款基于go语言开发的简易仿BTC区块链私有链模板，基于[blockchain_go](https://github.com/Jeiwan/blockchain_go "一个纯go语言编写的blockchain公链项目")进行提取和修改，包含一个区块链项目本地所需的基本模块，同时具备**创建账户**、**查询本机所有账户**、**创建创世区块**、**产生交易**、**输出区块链上所有交易**等区块链基本功能。满足最基本的区块链学习和测试需求。
## 参考教学视频
- [点我跳转至教程视频详情页（P1-P100）](https://www.bilibili.com/video/BV15T4y1B7TW/?vd_source=7ac88985bb2e529383ca0a4c99f675aa "区块链实战 | 基于Golang公链开发实战")
## 使用方法
本demo使用需要用到go及相关IDE，且go版本需小于等于1.18.8（IDE推荐使用[Goland](https://www.jetbrains.com/go/ "GoLand by JetBrains: More than just a Go IDE")，可以一键安装go并配置GOROOT和GOPATH，同时支持[学生免费获取](https://www.jetbrains.com/shop/eform/students "JetBrains Products for Learning")）
### ① 进入文件所在路径（windows环境）
首先通过进入项目所在包含**PrivateBlockchain-ByGo-WithoutNetworkServices**文件夹的目录，然后在其路径搜索栏输入**cmd**打开命令行工具；
在**命令行工具**中运行如下命令，进入**PrivateBlockchain-ByGo-WithoutNetworkServices**文件夹中；
`cd PrivateBlockchain-ByGo-WithoutNetworkServices`
进入成功后，你的IDE的terminal界面应如下图所示：

![成功图](https://raw.githubusercontent.com/HelloHaoWu/PrivateBlockchain-sample-ByGo/main/PrivateBlockchain-ByGo-WithoutNetworkServices/Images/%E8%BF%9B%E5%85%A5%E5%90%8E%E7%9A%84cmd%E7%8A%B6%E6%80%81.png)
### ② 运行程序
程序通过`./main`方法进行调用。在上述界面输入`./main`，会显示该程序的命令行提示。

![第二步成功图](https://raw.githubusercontent.com/HelloHaoWu/PrivateBlockchain-sample-ByGo/main/PrivateBlockchain-ByGo-WithoutNetworkServices/Images/%E8%BF%90%E8%A1%8Cmain%E5%90%8E%E7%8A%B6%E6%80%81.png)
### ③ 创建交易钱包
运行`./main CreateWallets`，即可创建你的区块交易钱包。

![第3步成功](https://raw.githubusercontent.com/HelloHaoWu/PrivateBlockchain-sample-ByGo/main/PrivateBlockchain-ByGo-WithoutNetworkServices/Images/%E5%88%9B%E5%BB%BAwallet.png)
### ④ 创建创世区块
运行`./main CreateBlockchain -address AUfDgH7UYzs67r16YFJgVVSYuekHmBJbrp`（将钱包地址替换成你自己的钱包地址）即可创建该区块链的创世区块。

![第4步成功](https://raw.githubusercontent.com/HelloHaoWu/PrivateBlockchain-sample-ByGo/main/PrivateBlockchain-ByGo-WithoutNetworkServices/Images/%E5%88%9B%E5%BB%BA%E5%88%9B%E4%B8%96%E5%8C%BA%E5%9D%97.png)
### ④ 检查对应钱包的货币剩余
运行`./main getbalance -address AUfDgH7UYzs67r16YFJgVVSYuekHmBJbrp`（将钱包地址替换成你自己的钱包地址），即可检查对应钱包的货币剩余，生成创世区块会使得其对应钱包内包含10枚货币。

![第5步成功](https://raw.githubusercontent.com/HelloHaoWu/PrivateBlockchain-sample-ByGo/main/PrivateBlockchain-ByGo-WithoutNetworkServices/Images/%E8%BE%93%E5%87%BA%E8%B4%A7%E5%B8%81.png)
### ⑤ 实现交易发送
再次**创建一个交易钱包地址**，然后输入`./main send -from '[\"AUfDgH7UYzs67r16YFJgVVSYuekHmBJbrp\"]' -to '[\"AFnVAZzHm98B2wvV8ZKVSaoeGbxa6yZ2Xx\"]' -amount '[\"2\"]'`（将对应的发送和接收钱包Hash值替换为你自己的钱包Hash值），即可进行价值2个代币的转账。同时，由于**挖矿奖励**，挖矿方（**在该私有链中为发送方**）会获得1个代币的奖励。

![第6步成功](https://github.com/HelloHaoWu/PrivateBlockchain-sample-ByGo/blob/main/PrivateBlockchain-ByGo-WithoutNetworkServices/Images/%E5%8F%91%E9%80%81%E6%88%90%E5%8A%9F.png)
### ⑥ 完整区块链输出
运行`./main printchain`，即可显示当前完整区块链数据的情况。
## 扩展
该项目已有**具有网络服务**的多节点部署版本，如有合作可联系我。
## 联系方式
网易邮箱：y4782266@163.com
