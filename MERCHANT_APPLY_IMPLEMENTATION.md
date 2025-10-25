# 商户申请参加展会功能实现总结

## 功能概述
基于当前数据表设计，参考现有代码实现风格，成功实现了商户申请参加展会的完整功能，包括状态检查和状态流转管理。

## 实现内容

### 1. 数据模型层 (Model Layer)
- **实体定义**: `internal/model/entity/t_exhibition_merchant.go`
  - 定义了展会与商户关联表的实体结构
- **业务模型**: `internal/model/exhibition.go`
  - 新增了 `ExhibitionMerchantStatus` 状态枚举
  - 新增了 `ExhibitionMerchantEvent` 事件枚举
  - 新增了 `ExhibitionMerchant` 业务模型
  - 实现了状态和事件的文本转换函数

### 2. 数据访问层 (DAO Layer)
- **DAO实现**: `internal/dao/internal/t_exhibition_merchant.go`
  - 实现了展会与商户关联表的DAO操作
  - 支持事务操作和上下文管理
- **DAO注册**: `internal/dao/t_exhibition_merchant.go`
  - 注册了全局DAO实例

### 3. 业务逻辑层 (Logic Layer)
- **业务逻辑**: `internal/logics/exhibition_merchant.go`
  - 实现了商户申请参加展会的核心业务逻辑
  - 包含完整的状态检查和验证
  - 实现了状态流转机制
  - 支持以下操作：
    - `ApplyForExhibition`: 商户申请参加展会
    - `GetMerchantApplication`: 获取商户申请状态
    - `ListExhibitionApplications`: 获取展会申请列表
    - `ListMerchantApplications`: 获取商户申请列表
    - `HandleEvent`: 处理状态事件

### 4. 接口层 (Interface Layer)
- **接口定义**: `internal/interfaces/logics.go`
  - 新增了 `IExhibitionMerchant` 接口
  - 定义了完整的业务方法签名

### 5. 服务层 (Service Layer)
- **服务实现**: `internal/service/exhibition.go`
  - 实现了商户申请参加展会的HTTP服务接口
  - 包含数据转换和错误处理
  - 支持以下API：
    - `ApplyForExhibition`: 商户申请参加展会
    - `GetMerchantApplication`: 获取商户申请状态
    - `ListExhibitionApplications`: 获取展会申请列表
    - `ListMerchantApplications`: 获取商户申请列表

### 6. API层 (API Layer)
- **API定义**: `api/v1/system/exhibition.go`
  - 定义了完整的请求和响应结构
  - 包含参数验证和文档注释
  - 支持RESTful API设计

### 7. 主程序集成 (Main Integration)
- **服务注册**: `main.go`
  - 集成了新的展会与商户关联服务
  - 更新了服务依赖注入

## 状态管理

### 申请状态
- `ExhibitionMerchantStatusPending` (0): 待审核
- `ExhibitionMerchantStatusApproved` (1): 审核通过
- `ExhibitionMerchantStatusRejected` (2): 审核拒绝
- `ExhibitionMerchantStatusWithdrawn` (3): 已退出

### 状态事件
- `ExhibitionMerchantEventApply` (0): 申请参加
- `ExhibitionMerchantEventApprove` (1): 审核通过
- `ExhibitionMerchantEventReject` (2): 审核拒绝
- `ExhibitionMerchantEventWithdraw` (3): 退出展会

### 状态流转规则
- 待审核 → 审核通过/审核拒绝/退出申请
- 审核通过 → 退出申请
- 审核拒绝 → 重新申请
- 已退出 → 重新申请

## 状态检查逻辑

### 申请前检查
1. **展会状态检查**: 展会必须处于"报名中"状态
2. **商户状态检查**: 商户必须处于"已审核"状态
3. **重复申请检查**: 防止重复申请（待审核/已通过状态）
4. **重新申请支持**: 被拒绝或退出的申请可以重新提交

### 状态转换检查
- 实现了完整的状态转换映射
- 确保状态转换的合法性
- 支持状态事件的原子性操作

## API接口

### 1. 商户申请参加展会
```
POST /api/v1/exhibition-service/exhibitions/{exhibition_id}/apply
```
**请求参数**:
- `exhibition_id`: 展会ID (路径参数)
- `merchant_id`: 商户ID (请求体)

**响应**:
- 成功: 返回申请结果消息
- 失败: 返回具体错误信息

### 2. 获取商户申请状态
```
GET /api/v1/exhibition-service/exhibitions/{exhibition_id}/applications/{merchant_id}
```

### 3. 获取展会申请列表
```
GET /api/v1/exhibition-service/exhibitions/{exhibition_id}/applications
```

### 4. 获取商户申请列表
```
GET /api/v1/exhibition-service/merchants/{merchant_id}/applications
```

## 技术特点

1. **遵循现有架构**: 完全按照现有代码风格和架构模式实现
2. **完整的状态管理**: 实现了完整的状态流转和事件处理机制
3. **事务支持**: 所有数据库操作都支持事务
4. **错误处理**: 完善的错误处理和错误码管理
5. **参数验证**: 完整的请求参数验证
6. **文档注释**: 详细的API文档注释

## 编译验证
- 所有代码已通过Go编译器验证
- 无语法错误和类型错误
- 符合Go语言最佳实践

## 总结
成功实现了商户申请参加展会的完整功能，包括数据模型、业务逻辑、API接口和状态管理。实现严格遵循了现有代码的架构模式和编码风格，确保了代码的一致性和可维护性。
