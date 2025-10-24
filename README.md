# TODO List
- 入参的参数校验没做，比如长度、格式等。
- Company某些字段，应该是唯一的。
- 系统架构部分，补充服务交互逻辑。最好详细点。

# 项目分析
- ExhibitionService 是展会管理平台的核心服务，负责提供所有展会相关的业务功能.
- 包括移动端用户接口和管理后台接口。系统采用微服务架构，与身份认证服务、消息推送服务、文件服务等独立服务进行交互，提供从展会预告到直播互动的完整业务闭环。

## 业务逻辑
- 登录Web后台后，才可选择入驻服务类型。

### 登录逻辑
- 两种登录方式
  - 手机号注册(后续可选择绑定IUQT账户)。
  - 通过IUQT账户一键登录，后台自动创建关联(登录后需要补充手机号信息)。

### 入驻逻辑
```mermaid
flowchart TD
    A[登录平台] --> B[选择入驻服务类型]
    
    B --> C1[商户（展商）]
    B --> C2[服务商（展会组织方）]
    
    C1 --> D[填写申请材料，提交审核]
    C2 --> D
    
    D --> E[等待审核结果]
    
    E --> F{审核结果}
    
    F -->|通过| G[签署合同]
    F -->|拒绝| H[根据拒绝原因修改资质信息]
    
    H --> I[重新提交审核]
    I --> E
    
    G --> J[完成入驻]
    
    %% 注释说明
    classDef default fill:#e1f5fe,stroke:#01579b,stroke-width:2px
    classDef decision fill:#fff3e0,stroke:#ef6c00,stroke-width:2px
    class F decision
```

## 业务模型
- <span style="color:red">**展会平台(Platform)**</span>
  - 技术平台提供商，也就是展会的承办方。
- <span style="color:red">**公司(Company)**</span>
  - 维护核心资质
    - 营业执照
    - 统一社会信用代码
    - 法人姓名
    - 法人证件号
    - 法人证件照
- <span style="color:red">**服务提供商(ServiceProvider)**</span>
  - 依赖公司主体存在(平台运营规则)
  - 继承公司基础属性 + 服务提供商专属属性

- <span style="color:red">**商户(Merchant)**</span>
  - 依赖公司主体存在(平台运营规则)
  - 继承公司基础属性 + 扩展商户专属属性
  - <span style="color:red">同一个公司，可以创建多个商户身份</span>。但独立运营。

- <span style="color:red">**展会(Exhibition)**</span>


- <span style="color:red">访客</span>
  - 匿名用户（最低权限）
    - 仅能浏览公开展会信息
  - 注册用户（标准权限）
    - 可收藏展会、预约参观。
  - 认证买家（高级权限）
    - 可联系展商、发起采购询盘
    
# 系统架构

## 服务职责划分

### ExhibitionService (展会服务)
- **核心职责**：提供所有展会相关的业务功能
- **服务范围**：
  - 提供移动端用户接口(首页、搜索、展会、个人中心、消息中心)
  - 提供管理后台接口(IUQT官方、展会公司、商户后台)
  - 展会业务逻辑处理(微服务内权限管理)
  - 直播间管理(创建、删除等，推流由其他服务处理)

### 身份认证服务
- 用户登录、注册
- 基础角色、权限管理
- JWT Token管理

### 消息推送服务
- 实时消息推送
- 消息模板管理
- 推送渠道管理

### 文件服务
- 文件上传下载
- 文件存储管理
- 文件访问控制

### 直播服务
- 直播流推送

# 服务提供商

## 状态流转
- Pending(0): 待审核。
- Approved(1): 审核通过。此时服务商可正常开展业务。
- Rejected(2): 审核驳回。此时服务商可修改资料信息，然后再次提交审核。
- Disabled(3): 禁用。服务商被平台禁用（因违规等），不可开展业务（区别于注销）。
- UnRegisted(4): 注销。服务商主动退出/运营强制退出。
```mermaid
stateDiagram-v2
  [*] --> Pending: 填写入驻材料，提交审核

  Pending --> Approved: 审核通过
  Pending --> Rejected: 审核驳回
  
  Rejected --> Pending: 修改后再次提交审核
  
  Approved --> Disabled: 禁用
  Approved --> UnRegister: 主动注销/强制注销

  Disabled --> Approved: 解禁
```

## 数据表设计
```sql
```

## 接口设计
```curl
```

# 商户

## 状态流转
### 状态流转
- Pending(0): 信息已录入，等待平台运营审核。
- Approved(1): 审核通过，此时商户可正常开展业务。
- Rejected(2): 审核驳回，需修改后重新提交审核。
- Disabled(3): 禁用，商户被平台禁用（因违规等），不可开展业务（区别于注销）。
- UnRegisted(4): 注销，商户主动退出，账号永久失效。
```mermaid
stateDiagram-v2
  [*] --> Pending: 填写入驻材料，提交审核

  Pending --> Approved: 审核通过
  Pending --> Rejected: 审核驳回
  
  Rejected --> Pending: 修改后再次提交审核
  
  Approved --> Disabled: 禁用
  Approved --> UnRegister: 主动注销/强制注销

  Disabled --> Approved: 解禁
```

## 数据表设计
```sql
```

## 接口设计
```curl

```


# 展会
- 可以由单个服务提供商创建，也可以由多个服务提供商联合创建，也可以有协办商。
- 一个服务提供商可以同时创建多个展会。
- 一个商户可以同时参加多个展会。

## 业务逻辑
### 展会创建
- 创建展会时必须指定所有主办方。

## 状态流转
```mermaid
stateDiagram-v2
[*] --> Preparing: 创建展会

Preparing --> Pending: 提交审核
Preparing --> Cancelled: 主动取消

Pending --> Approved: 审核通过
Pending --> Preparing: 审核驳回

Approved --> Enrolling: 到达展会报名时间
Approved --> Cancelled: 主动取消

Enrolling --> Running: 到达展会开始时间
Enrolling --> Cancelled: 主动取消

Running --> Ended: 到达展会结束时间
Running --> Cancelled: 强制取消
```
### 状态定义
- Preparing(0): 筹备中。展会初始状态，进行基础信息配置、展位规划等准备工作。
- Pending(1): 待审核。提交审核后，等待运营人员审核（可退回修改）。
- Approved(2): 已批准。审核通过但是未到报名时间。
- Enrolling(3): 报名中。商家可以申请报名。访客可以预约展会。
- Running(4): 进行中。展会正式开放，参展商和观众可线上互动。
- Ended(5): 已结束。（自动归档数据）
- Cancelled(6): 已取消。主动终止展会（违规行为、运营调整等）

## 数据表设计
```sql
```


## 接口设计
```curl
POST  /api/v1/exhibition-service/exhibition             # 服务商 -- 创建展会(创建时指定所有主办方)
POST  /api/v1/exhibition-service/exhibition/application # 商户 -- 申请参加展会

