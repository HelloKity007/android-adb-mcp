# Android ADB MCP

融合 [Android-MCP](https://github.com/HelloKity007/Android-MCP)、[android-action-mcp](https://github.com/HelloKity007/android-action-mcp)、[android-mcp-server](https://github.com/HelloKity007/android-mcp-server)、[adb-tui](https://github.com/HelloKity007/adb-tui) 四个项目的完整 Android ADB MCP 解决方案。

## 功能特性

### MCP Server 模式
- 85+ 个 MCP 工具，覆盖 Android 设备操作全场景
- UI 树解析 + 标注截图（融合 Android-MCP）
- 智能表单填充、滚动、选择（融合 android-action-mcp）
- scrcpy 集成（融合 android-mcp-server）
- stdio / HTTP 双传输模式

### TUI 模式
- 基于 bubbletea 的交互式终端界面
- 设备管理、Shell、日志、文件浏览、包管理等 11 个视图
- 多主题支持

### ADB 客户端库
- 22 个模块，覆盖所有 ADB 操作
- 完整的单元测试覆盖

## 快速开始

```bash
# 编译
go build -o android-adb-mcp.exe ./cmd/server/

# MCP Server 模式 (stdio)
./android-adb-mcp.exe

# MCP Server 模式 (HTTP)
./android-adb-mcp.exe mcp --transport http --addr :8080
```

## Claude Code 配置

```json
{
  "mcpServers": {
    "android-adb-mcp": {
      "command": "E:\\code_local\\python\\android-adb-mcp\\android-adb-mcp.exe",
      "args": []
    }
  }
}
```

## 项目结构

```
cmd/
  server/          -- MCP Server 入口
internal/
  adb/             -- 核心 ADB 客户端库 (来源: adb-tui)
  mcp/             -- MCP 服务器实现 (来源: adb-tui)
    tools/         -- MCP 工具注册
    transport/     -- stdio/HTTP 传输层
  snapshot/        -- UI 快照引擎 (融合 Android-MCP)
  config/          -- 配置管理
pkg/
  jsonrpc/         -- JSON-RPC 2.0 实现
```

## 来源项目

| 项目 | 语言 | 贡献 |
|------|------|------|
| [adb-tui](https://github.com/HelloKity007/adb-tui) | Go | ADB 客户端库、MCP 协议、TUI 界面 |
| [android-action-mcp](https://github.com/HelloKity007/android-action-mcp) | Go | 智能表单、滚动、选择、对话框处理 |
| [Android-MCP](https://github.com/HelloKity007/Android-MCP) | Python | UI 树解析、标注截图、选择器点击 |
| [android-mcp-server](https://github.com/HelloKity007/android-mcp-server) | TypeScript | scrcpy 集成、ADB 自动下载 |
