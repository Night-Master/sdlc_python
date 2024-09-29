# sdlc_python
sdlc_python 是一个基于python语言构建的devsecops平台，旨在促进 DevSecOps 和安全开发生命周期 (SDLC) 实践。它通过模拟常见漏洞来增强开发人员的安全意识(对应sdlc中对开发人员的安全培训)，并且使用了大模型进行代码安全审计（对应sdlc中代码审计阶段），帮助企业进行安全左移。除了用于 DevSecOps 实践外，sdlc_python 还可以用于学习漏洞知识、渗透测试和代码审计。本项目采用了前后端分离的设计模式，其中后端利用了轻量级框架 Flask，而前端则使用了 Vue 3。

## 使用样例
sast扫描
![image](https://github.com/Night-Master/sdlc_python/blob/main/data/use3.png)
漏洞示范
![image](https://github.com/Night-Master/sdlc_python/blob/main/data/use2.png)

## 主要特性

- **前后端分离**：后端处理业务逻辑，前端负责用户交互。
- **轻量高效**：使用 flask 框架。
- **代码扫描**：使用大模型对代码进行安全检测，为devsecops中漏洞左移做保障。

## 技术栈

- **后端**: python (flask 框架)
- **前端**: Vue 3

## 安装与运行

1. 下载最新版本的发布包
2. 运行发布包中的start.bat
3. 在web登录界面输入账号：user1，密码：hello
### 系统需求

- 支持的操作系统: Windows x86
- 计划支持: Linux, ARM 架构
