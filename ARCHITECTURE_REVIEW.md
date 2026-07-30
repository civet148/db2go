# db2go v3 架构分析与优化建议报告

> 分析日期: 2026-07-29 | 代码版本: v3.8.3

---

## 一、当前架构总览

### 1.1 包依赖关系

```
main.go
  ├── schema/        (数据模型、代码生成、工具函数)
  │   ├── types.go          类型映射表 & CmdFlags
  │   ├── table_schema.go   数据结构 & Exporter 注册/工厂
  │   ├── gen_go.go         Go 文件生成逻辑
  │   ├── gen_proto.go      Proto 文件生成逻辑
  │   └── utils.go          工具函数 & 命名转换
  ├── mysql/         (MySQL 查询 + ExportGo/ExportProto 实现)
  ├── postgres/      (PostgreSQL 查询 + ExportGo/ExportProto 实现)
  ├── mssql/         (SQL Server 查询 + ExportGo/ExportProto 实现)
  ├── opengauss/     (OpenGauss 查询 + ExportGo/ExportProto 实现)
  └── parser/        (Go AST 解析 + 文件合并)
```

### 1.2 核心架构模式

- **工厂 + 注册模式**: `schema.Register()` / `schema.NewExporter()` 实现多数据库调度
- **策略模式雏形**: 各数据库包通过 `Exporter` 接口统一导出入口
- **模板方法模式缺失**: `ExportGo()` / `ExportProto()` 在每个 dialect 中重复实现

---

## 二、关键问题分析

### 2.1 严重代码重复 —— Dialect 间 Export 逻辑高度雷同

**`ExportProto()` 方法在 4 个包中几乎完全一致**（mysql:93-113, postgres:88-125, mssql:66-103, opengauss:85-122），差异仅有 `queryTableSchemas()` 调用方式。

**`ExportGo()` 方法结构相同**，差异仅在查询建库DDL（仅MySQL支持）。

```go
// 以 ExportProto 为例，4个文件几乎完全相同的模式：
func (m *ExporterXxx) ExportProto() (err error) {
    schemas, err = m.queryTableSchemas(...)  // ← 仅有这行不同
    strHead := schema.MakeProtoHead(cmd)
    for i, v := range schemas {
        m.queryTableColumns(v)
        strBody := schema.MakeProtoBody(cmd, v)
        file, _ = schema.CreateOutputFile(cmd, v, "proto", append)
        // 文件写入逻辑...
    }
}
```

**建议**: 将 ExportGo/ExportProto 提升到 `schema` 公共层，只在 dialect 中提供查询接口。

### 2.2 字符串拼接式代码生成 —— 难以维护和扩展

`gen_go.go`、`gen_proto.go` 大量使用 `fmt.Sprintf` + 字符串拼接生成代码：

- `makeTableStructure()` (`gen_go.go:302-387`) — 手拼 struct 定义
- `makeObjectMethods()` (`gen_go.go:176-198`) — 手拼 Getter/Setter
- `makeColumnConsts()` (`gen_go.go:285-293`) — 手拼常量
- `MakeProtoBody()` (`gen_proto.go:32-48`) — 手拼 message 定义

**问题**:
- 格式错误难以定位 (无语法校验)
- 调整输出格式需要修改多处拼接逻辑
- 无法利用 Go `text/template` 的模板复用能力

**建议**: 引入 `text/template` 引擎管理代码模板。

### 2.3 数据库查询缺乏统一抽象

每个 dialect 自行定义查询方法且签名不一致：

| 数据库 | 查询表方法 | 查询列方法 | 是否支持索引 | 是否支持 DDL |
|--------|-----------|-----------|------------|------------|
| MySQL | `queryTableSchemas(cmd, e)` | `queryTableColumns(v)` | 是 | 是 |
| PostgreSQL | `queryTableSchemas()` | `queryTableColumns(v)` | 是 | 否 |
| MSSQL | `queryTableSchemas()` | `queryTableColumns(v)` | 否 | 否 |
| OpenGauss | `queryTableSchemas()` | `queryTableColumns(v)` | 否 | 否 |

**建议**: 定义 `SchemaProvider` 接口统一查询契约。

### 2.4 Exporter 结构体职责过重

每个 dialect 的 Exporter 同时承担:
1. 数据库连接持有
2. SQL 查询执行
3. 结果映射
4. 文件生成调度

**建议**: 分离为 `Queryer`（数据查询） + `Generator`（代码生成） + `Writer`（文件输出）。

### 2.5 package 命名与导入不一致

- `opengauss/open_gauss.go` 包名 `opengauss` 与目录名一致但文件名 `open_gauss.go` 用了下划线
- `schema` 包名与目录一致，但内部同时拥有 `types.go` 中的类型映射数据和 `gen_*.go` 的生成逻辑
- 建议将 `schema/` 按职责拆分为 `model/` + `generator/` + `writer/`

### 2.6 parser 包过于复杂

`parser/go_parser.go:776` 行，大量基于正则的 AST 解析逻辑。

```go
// 多个正则编译在热路径中反复执行:
reLineTailComment := regexp.MustCompile(...)
reBlockComment := regexp.MustCompile(...)
reImport := regexp.MustCompile(...)
reIdent := regexp.MustCompile(...)
reAssign := regexp.MustCompile(...)
reField := regexp.MustCompile(...)
```

**建议**: 将 `regexp.MustCompile` 提升为包级别变量避免重复编译；简化合并策略。

### 2.7 错误处理风格不统一

交叉三种模式：
1. 返回 error（正确做法）
2. `log.Error()` + return nil（丢失错误链）
3. 直接 `os.Exit(1)`（`main.go:196`）

```go
// mssql.go:58-62 — 返回 nil 掩盖错误
if err = m.queryTableColumns(v); err != nil {
    log.Error(err.Error())
    return   // ← 应该 return err
}
```

### 2.8 全局可变状态

```go
// table_schema.go:131
var instances = make(map[string]Instance, 1)
```

多线程并发注册或使用寄存器存在竞态。同时 `CmdFlags` 对象作为巨型配置对象传递。

### 2.9 测试覆盖严重不足

整个项目只有一个测试文件 `parser/parser_test.go`，仅测试 AST 合并，**核心代码生成逻辑零测试**。

### 2.10 类型映射表孤立维护

`schema/types.go` 中的 `db2goTypes`、`db2protoTypes`、`db2goTypesUnsigned`、`db2protoTypesUnsigned` 四张表需要手动同步维护，难以保证一致性。

---

## 三、架构优化方案

### 3.1 目标架构（建议）

```
main.go
├── model/              (纯数据模型)
│   ├── table.go            TableSchema, TableColumn, TableIndex
│   ├── config.go           CmdFlags → 拆分为 Config / Options
│   └── types.go            类型映射表
├── dialect/            (数据库方言实现)
│   ├── interface.go        SchemaProvider 接口定义
│   ├── mysql.go            MySQL 查询实现
│   ├── postgres.go         PostgreSQL 查询实现
│   ├── mssql.go            MSSQL 查询实现
│   └── opengauss.go        OpenGauss 查询实现
├── generator/          (代码生成器)
│   ├── go_generator.go     Go 代码生成
│   ├── proto_generator.go  Proto 代码生成
│   ├── templates/          模板文件
│   │   ├── go_struct.tmpl
│   │   ├── go_getter.tmpl
│   │   ├── go_setter.tmpl
│   │   └── proto_message.tmpl
│   └── converter.go        命名转换工具
├── writer/             (文件输出 + Merge)
│   ├── file_writer.go      文件写入
│   └── merger.go           代码合并
└── engine/             (调度引擎)
    └── engine.go           ExportGo / ExportProto 统一流程
```

### 3.2 核心重构点

#### 3.2.1 抽取公共 Export 流程（消除 80% 重复代码）

```go
// dialect/interface.go
type SchemaProvider interface {
    QueryTableSchemas(cmd *model.Config) ([]*model.TableSchema, error)
    QueryTableColumns(table *model.TableSchema) error
    QueryTableIndexes(table *model.TableSchema) error
    QueryCreateDatabaseDDL(cmd *model.Config) (*model.CreateDatabaseDDL, error) // optional
}
```

```go
// engine/engine.go — 统一的导出调度
func ExportGo(provider SchemaProvider, cmd *model.Config) error {
    schemas, _ := provider.QueryTableSchemas(cmd)
    for _, s := range schemas {
        provider.QueryTableColumns(s)
        provider.QueryTableIndexes(s)
        generator.GenerateGoFile(cmd, s)
    }
}
```

每个 dialect 只需实现 `SchemaProvider` 接口，无需再写 ExportGo/ExportProto。

#### 3.2.2 引入模板引擎

将 `gen_go.go` 和 `gen_proto.go` 中的字符串拼接改为 `text/template`：

```go
// 模板: go_struct.tmpl
type {{.StructName}} struct {
{{- range .Columns }}
    {{.GoName}} {{.GoType}} `json:"{{.JSONName}}" db:"{{.DBName}}" {{.Tags}}` // {{.Comment}}
{{- end }}
}
```

#### 3.2.3 拆分 schema 包

- `model/` — 纯数据结构（TableSchema, TableColumn, CmdFlags → Config）
- `generator/` — 代码生成逻辑（Go 和 Proto）
- `writer/` — 文件输出 + merge

#### 3.2.4 类型映射表自动化

```go
// 用单一源表自动推导 signed/unsigned 和 proto 映射
var db2goTypeMap = map[string]TypeMapping{
    DB_COLUMN_TYPE_BIGINT: {Signed: "int64", Unsigned: "uint64", ProtoSigned: "sint64", ProtoUnsigned: "uint64"},
    DB_COLUMN_TYPE_INT:    {Signed: "int32", Unsigned: "uint32", ProtoSigned: "sint32", ProtoUnsigned: "uint32"},
    // ...
}
```

#### 3.2.5 错误处理统一

- 所有数据库查询方法必须返回 `error`
- 引入 `errors.Wrap` 链式错误
- 移除无意义的 `log.Error()` + `return nil` 模式

#### 3.2.6 测试策略建议

| 层级 | 测试内容 | 建议框架 |
|------|---------|---------|
| 单元测试 | 类型映射、命名转换、标签生成 | go test |
| 快照测试 | 代码生成输出与预期比对 | go test + golden files |
| 集成测试 | 连接真实 DB 验证查询 | docker-compose + testcontainers |

### 3.3 速赢优化（低风险高收益）

| 优先级 | 优化项 | 文件 | 工作量 |
|--------|--------|------|--------|
| P0 | `regexp.MustCompile` 提升为包级变量 | `parser/go_parser.go` | ~5min |
| P0 | 修复 `mssql.go:58` / `opengauss.go:78` 等 return nil 掩盖错误 | 各 dialect | ~10min |
| P1 | 提取公有 ExportProto 到 schema 层（消除 4 份拷贝） | 所有 dialect | ~1h |
| P1 | 提取公有 ExportGo 中的公共循环逻辑 | 所有 dialect | ~1h |
| P2 | 将 db2goTypes/db2protoTypes 合并为单一映射表 | `schema/types.go` | ~30min |
| P2 | 添加生成代码快照测试 | `gen_go.go` / `gen_proto.go` | ~2h |

---

## 四、总结

db2go v3 的核心问题在于 **dialect 间 Export 流程的高度重复** 和 **代码生成的字符串拼接模式**。当前架构在 v2 → v3 的演进中引入了 `Exporter` 接口和注册机制，但未将公共流程抽象到位，导致新增数据库方言需要拷贝大量样板代码。

建议优先抽取公共 Export 流程（消除 ~70% 重复代码），再逐步将代码生成迁移至模板引擎，最后补充测试覆盖。整体重构风险可控，可增量实施。
