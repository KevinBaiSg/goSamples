# cobra

* 初始化项目 `cobra init --pkg-name github.com/KevinBaiSg/cobra`
* 添加command `cobra add config`

## cobra 提供的功能
- 简易的子命令行模式，如 `app server`， `app fetch` 等.
- 完全兼容 posix 命令行模式(包括 short & long 版本)
- 嵌套子命令
- 支持全局，局部，串联 flags
- 支持生成应用和命令，示例：`cobra init appname`， `cobra add cmdname`.
- 智能建议(app srver... did you mean app server?)
- 自动生成 commands 和 flags 的帮助信息
- 自动生成详细的 `-h, --help` 信息.
- 自动生成应用程序在 bash 下命令自动完成功能
- 自动生成应用程序的man手册
- 命令行别名
- 灵活定义help和usage信息
- 可选的紧密集成的 viper apps

## 概念
Cobra is built on a structure of commands, arguments & flags.
Commands represent actions, Args are things and Flags are modifiers for those actions.

### Commands
### Args
### Flags


Field name   | Mandatory? | Allowed values  | Allowed special characters
----------   | ---------- | --------------  | --------------------------
Minutes      | Yes        | 0-59            | * / , -
Hours        | Yes        | 0-23            | * / , -
Day of month | Yes        | 1-31            | * / , - ?
Month        | Yes        | 1-12 or JAN-DEC | * / , -
Day of week  | Yes        | 0-6 or SUN-SAT  | * / , - ?