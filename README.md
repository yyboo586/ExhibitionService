# ExhibitionService 展会服务系统设计文档

## 📋 项目概述
ExhibitionService 是展会管理平台的核心服务，负责提供所有展会相关的业务功能，包括移动端用户接口和管理后台接口。系统采用微服务架构，与身份认证服务、消息推送服务、文件服务等独立服务进行交互，提供从展会预告到直播互动的完整业务闭环。

- 展会平台: 技术平台提供商。 
- 展会主办方: 展会的组织者，通常是政府或者行业协会。
    - 确定展会主题、定位和受众。
    - 指定整体策略、规划展会内容和形式。
    - 负责招商（吸引商户参展）和招观（吸引访客）。
    - 指定展会规则、政策、时间表。
- 商户: 参展商/卖家/内容提供者。
    - 在分配的虚拟展位空间内搭建和布置自己的展示内容（产品资料、公司介绍、宣传视频等）。
    - 安排人员在线值守，通过聊天、直播、视频会议等方式与访客进行实时互动、答疑解惑、上午洽谈。
    - 发布直播、研讨会、新品发布会等活动。
- 访客: 观众/买家/学习者/行业人士。
    - 注册、登录线上展会
    - 浏览展会整体信息、主题内容
    - 搜索和观看感兴趣的商户展位、查看产品、资料、视频
    - 参与观看直播、发送聊天内容




## 🏗️ 系统架构

### 服务职责划分

#### ExhibitionService (展会服务)
- **核心职责**：提供所有展会相关的业务功能
- **服务范围**：
  - 移动端用户接口（首页、搜索、展会、个人中心、消息中心）
  - 管理后台接口（IUQT官方、展会公司、商户后台）
  - 展会业务逻辑处理
  - 直播系统管理
  - 审核流程管理

#### 独立服务
- **身份认证服务 (AuthService)**
  - 用户登录、注册
  - 角色权限管理
  - JWT Token管理
- **消息推送服务 (NotificationService)**
  - 实时消息推送
  - 消息模板管理
  - 推送渠道管理
- **文件服务 (FileService)**
  - 文件上传下载
  - 文件存储管理
  - 文件访问控制
- **直播服务 (LiveServer)**
  - 直播流推送
  - 那直播间由谁管理呢

### 角色划分
- **普通用户**：浏览展会、观看直播、参与互动
- **展会公司**：创建展会、管理商户、控制直播
- **商户**：参与展会、进行直播展示
- **IUQT官方**：平台管理、审核、监控

### 服务交互方式
```
ExhibitionService ←→ AuthService (身份验证、权限检查)
ExhibitionService ←→ NotificationService (消息推送)
ExhibitionService ←→ FileService (文件上传下载)
ExhibitionService ←→ LiveServer (直播流)
```

## 🎯 ExhibitionService API 接口设计
- API前缀为 /api/v1/exhibition-service

### 📱 移动端用户接口

#### 1. 展会首页模块

##### 功能特性
- **推荐展示**
  - Banner轮播图
  - 消息推送中心
  - 直播频道推荐

##### ExhibitionService 接口
```
GET /api/v1/exhibition-service/home/banner                    # 获取Banner列表
GET /api/v1/exhibition-service/home/recommendations           # 获取推荐内容
GET /api/v1/exhibition-service/home/live-channels             # 获取直播频道
```

##### 与独立服务交互
```
POST /api/v1/exhibition-service/home/notifications/send       # 发送消息推送 (调用NotificationService)
GET /api/v1/exhibition-service/home/banner/{id}/image         # 获取Banner图片 (调用FileService)
```

##### 消息推送类型
- 预约展会开始通知
- 关注公司开展会通知
- 直播推送通知

#### 2. 搜索模块

##### 功能特性
- **展会搜索**
  - 按展会名称/关键词搜索
  - 按直播间关键词/类型搜索
  - 搜索结果分类展示

##### ExhibitionService 接口
```
GET /api/v1/exhibition-service/search/exhibitions             # 搜索展会
GET /api/v1/exhibition-service/search/live-streams            # 搜索直播间
GET /api/v1/exhibition-service/search/suggestions              # 搜索建议
```

#### 3. 展会模块

##### 功能特性
- **展会列表**
  - 展会基础信息展示
  - 主办方信息管理
  - 直播列表管理
  - 展会预告功能

##### ExhibitionService 接口
```
GET /api/v1/exhibition-service/exhibitions/list               # 获取展会列表
GET /api/v1/exhibition-service/exhibitions/{id}               # 获取展会详情
GET /api/v1/exhibition-service/exhibitions/{id}/organizer     # 获取主办方信息
GET /api/v1/exhibition-service/exhibitions/{id}/live-streams   # 获取直播列表
POST /api/v1/exhibition-service/exhibitions/{id}/reserve      # 预约展会
POST /api/v1/exhibition-service/exhibitions/{id}/favorite     # 收藏展会
```

##### 直播功能接口
```
GET /api/v1/exhibition-service/live/{id}/info                 # 获取直播信息
GET /api/v1/exhibition-service/live/{id}/stats                # 获取实时数据
POST /api/v1/exhibition-service/live/{id}/like                # 点赞
POST /api/v1/exhibition-service/live/{id}/comment             # 发送弹幕
POST /api/v1/exhibition-service/live/{id}/share               # 分享直播
POST /api/v1/exhibition-service/live/{id}/connect             # 连麦申请
```

##### 与独立服务交互
```
GET /api/v1/exhibition-service/exhibitions/{id}/images        # 获取展会图片 (调用FileService)
POST /api/v1/exhibition-service/exhibitions/{id}/notify        # 发送展会通知 (调用NotificationService)
```

#### 4. 个人中心模块

##### 功能特性
- **个人信息管理**
- **展会参与记录**
- **收藏管理**
- **关注管理**
- **黑名单管理**

##### ExhibitionService 接口
```
GET /api/v1/exhibition-service/user/profile                   # 获取用户信息
PUT /api/v1/exhibition-service/user/profile                   # 更新用户信息
GET /api/v1/exhibition-service/user/exhibitions               # 获取参与的展会
GET /api/v1/exhibition-service/user/favorites                 # 获取收藏列表
GET /api/v1/exhibition-service/user/follows                   # 获取关注列表
GET /api/v1/exhibition-service/user/blacklist                 # 获取黑名单
POST /api/v1/exhibition-service/user/follow/{id}              # 关注用户/公司
POST /api/v1/exhibition-service/user/blacklist/{id}           # 拉黑用户/公司
DELETE /api/v1/exhibition-service/user/blacklist/{id}         # 取消拉黑
```

##### 与独立服务交互
```
GET /api/v1/exhibition-service/user/profile/avatar           # 获取用户头像 (调用FileService)
PUT /api/v1/exhibition-service/user/profile/avatar            # 更新用户头像 (调用FileService)
POST /api/v1/exhibition-service/user/notifications/subscribe # 订阅通知 (调用NotificationService)
```

#### 5. 消息中心模块

##### 功能特性
- **系统通知**
- **审核结果通知**
- **封禁通知**

##### ExhibitionService 接口
```
GET /api/v1/exhibition-service/messages/list                   # 获取消息列表
GET /api/v1/exhibition-service/messages/{id}                  # 获取消息详情
PUT /api/v1/exhibition-service/messages/{id}/read             # 标记已读
DELETE /api/v1/exhibition-service/messages/{id}               # 删除消息
```

##### 与独立服务交互
```
POST /api/v1/exhibition-service/messages/send                  # 发送消息 (调用NotificationService)
GET /api/v1/exhibition-service/messages/templates              # 获取消息模板 (调用NotificationService)
```

### 🖥️ 管理后台接口

#### 1. IUQT官方后台模块

##### 功能特性
- **展会公司管理**
- **展会管理**
- **商户管理**
- **用户管理**
- **申请审核**

# 展会公司
- **统一社会信用代码**: 是中国大陆境内依法注册的法人和其他组织的唯一、终身不变的身份识别码。
## 数据表设计
```sql
-- 主办方必须是依法注册的实体
CREATE TABLE IF NOT EXISTS `t_company` (
    `id` VARCHAR(40) NOT NULL COMMENT '展会公司ID',
    `name` VARCHAR(32) NOT NULL COMMENT '展会公司名称',
    `country` VARCHAR(255) NOT NULL COMMENT '国家',
    `status` TINYINT(1) NOT NULL COMMENT '公司状态(0:待审核、1:审核通过、2:禁用、3:注销)',
    `phone` VARCHAR(32) COMMENT '手机',
    `email` VARCHAR(32) COMMENT '邮箱',
    `address` VARCHAR(255) COMMENT '地址',
    `description` TEXT COMMENT '公司描述',

    `business_license` VARCHAR(255) DEFAULT '' COMMENT '营业执照',
    `social_credit_code` VARCHAR(255) NOT NULL COMMENT '统一社会信用代码',
    `legal_person_name` VARCHAR(32) DEFAULT '' COMMENT '法人姓名',
    `legal_person_card_number` VARCHAR(255) DEFAULT '' COMMENT '法人证件号',
    `legal_person_photo_url` TEXT COMMENT '法人证件照',
    `legal_person_phone` VARCHAR(32) COMMENT '法人手机号',

    `apply_time` BIGINT(20) NOT NULL COMMENT '申请时间',
    `approve_time` BIGINT(20) COMMENT '入驻时间',
    `create_time` BIGINT(20) NOT NULL COMMENT '创建时间',
    `update_time` BIGINT(20) NOT NULL COMMENT '更新时间',
    PRIMARY KEY(`id`),
    UNIQUE KEY `idx_name` (`name`)
)ENGINE = InnoDB COMMENT '主办方';
```

## 接口设计
- 创建展会公司
```curl
POST   /api/v1/exhibition-service/company               # 申请创建展会公司
DELETE /api/v1/exhibition-service/company/:id           # 删除展会公司
PATCH  /api/v1/exhibition-service/company/:id           # 修改展会公司信息
GET    /api/v1/exhibition-service/company/:id           # 获取展会公司详情
GET    /api/v1/exhibition-service/company?name=xx       # 展会公司列表（如果查询参数不为空，则为搜索）

PATCH  /api/v1/exhibition-service/company/:id/approve   # 申请通过
PATCH  /api/v1/exhibition-service/company/:id/reject    # 申请拒绝
PATCH  /api/v1/exhibition-service/company/:id/ban       # 禁用展会公司
PATCH  /api/v1/exhibition-service/company/:id/unban     # 禁用展会公司
GET    /api/v1/exhibition-service/company/applications  # 获取创建申请
```

## 业务逻辑
### 展会公司创建
- 登录Web后台后，才可申请入驻。在此之前，需要先创建账户。
- 两种注册账户的方式
  - 手机号注册(独立账户，可选择绑定IUQT账户)。
  - 通过IUQT账户一键登录，后台创建关联用户。
- 账户登录
  - 手机号登录/独立账户登录。
  - 通过IUQT账户一键登录。



```
# 展会管理
GET /api/v1/exhibition-service/admin/exhibitions              # 获取展会列表
GET /api/v1/exhibition-service/admin/exhibitions/{id}         # 获取展会详情
POST /api/v1/exhibition-service/admin/exhibitions/{id}/shutdown # 关停展会
POST /api/v1/exhibition-service/admin/exhibitions/{id}/restart # 重启展会

# 商户管理
GET /api/v1/exhibition-service/admin/merchants                # 获取商户列表
GET /api/v1/exhibition-service/admin/merchants/{id}           # 获取商户详情
POST /api/v1/exhibition-service/admin/merchants/{id}/ban      # 封禁商户
POST /api/v1/exhibition-service/admin/merchants/{id}/unban    # 解封商户

# 用户管理
GET /api/v1/exhibition-service/admin/users                    # 获取用户列表
POST /api/v1/exhibition-service/admin/users/{id}/ban          # 封禁用户
POST /api/v1/exhibition-service/admin/users/{id}/unban        # 解封用户
DELETE /api/v1/exhibition-service/admin/users/{id}            # 删除用户

# 申请审核
GET /api/v1/exhibition-service/admin/applications             # 获取申请列表
POST /api/v1/exhibition-service/admin/applications/{id}/approve # 审核通过
POST /api/v1/exhibition-service/admin/applications/{id}/reject # 审核驳回
```

##### 与独立服务交互
```
GET /api/v1/exhibition-service/admin/users/{id}/auth-info     # 获取用户认证信息 (调用AuthService)
POST /api/v1/exhibition-service/admin/notifications/send       # 发送管理通知 (调用NotificationService)
GET /api/v1/exhibition-service/admin/files/{id}                # 获取审核文件 (调用FileService)
```

#### 2. 展会公司后台模块

##### 功能特性
- **展会管理**
- **商户管理**
- **人员管理**
- **消息中心**

##### ExhibitionService 接口
```
# 展会管理
GET /api/v1/exhibition-service/company/exhibitions            # 获取展会列表
POST /api/v1/exhibition-service/company/exhibitions           # 创建展会
PUT /api/v1/exhibition-service/company/exhibitions/{id}       # 更新展会
GET /api/v1/exhibition-service/company/exhibitions/{id}/stats # 获取展会统计

# 商户管理
GET /api/v1/exhibition-service/company/merchants              # 获取商户列表
POST /api/v1/exhibition-service/company/merchants/{id}/approve # 审核商户申请
POST /api/v1/exhibition-service/company/merchants/{id}/ban    # 封禁商户
POST /api/v1/exhibition-service/company/merchants/{id}/unban  # 解封商户

# 直播管理
GET /api/v1/exhibition-service/company/live-streams           # 获取直播列表
POST /api/v1/exhibition-service/company/live-streams/{id}/control # 直播控制
GET /api/v1/exhibition-service/company/live-streams/{id}/stats # 获取直播统计

# 人员管理
GET /api/v1/exhibition-service/company/users                  # 获取用户列表
POST /api/v1/exhibition-service/company/users                 # 添加用户
PUT /api/v1/exhibition-service/company/users/{id}/permissions # 设置用户权限
DELETE /api/v1/exhibition-service/company/users/{id}          # 删除用户
```

##### 与独立服务交互
```
POST /api/v1/exhibition-service/company/exhibitions/{id}/files # 上传展会文件 (调用FileService)
POST /api/v1/exhibition-service/company/notifications/send     # 发送公司通知 (调用NotificationService)
GET /api/v1/exhibition-service/company/users/{id}/auth-check   # 验证用户权限 (调用AuthService)
```

#### 3. 商户后台模块

##### 功能特性
- **展会参与管理**
- **直播管理**
- **商户信息管理**
- **人员管理**

##### ExhibitionService 接口
```
# 展会参与
GET /api/v1/exhibition-service/merchant/exhibitions           # 获取参与的展会
POST /api/v1/exhibition-service/merchant/exhibitions/{id}/apply # 申请参与展会
POST /api/v1/exhibition-service/merchant/exhibitions/{id}/exit # 退出展会

# 直播管理
GET /api/v1/exhibition-service/merchant/live-streams          # 获取直播列表
POST /api/v1/exhibition-service/merchant/live-streams         # 创建直播
PUT /api/v1/exhibition-service/merchant/live-streams/{id}     # 更新直播信息
POST /api/v1/exhibition-service/merchant/live-streams/{id}/start # 开始直播
POST /api/v1/exhibition-service/merchant/live-streams/{id}/stop # 结束直播
GET /api/v1/exhibition-service/merchant/live-streams/{id}/stats # 获取直播统计

# 商户信息
GET /api/v1/exhibition-service/merchant/profile               # 获取商户信息
PUT /api/v1/exhibition-service/merchant/profile               # 更新商户信息
POST /api/v1/exhibition-service/merchant/profile/submit       # 提交审核

# 人员管理
GET /api/v1/exhibition-service/merchant/users                 # 获取用户列表
POST /api/v1/exhibition-service/merchant/users                # 添加用户
PUT /api/v1/exhibition-service/merchant/users/{id}/permissions # 设置用户权限
```

##### 与独立服务交互
```
POST /api/v1/exhibition-service/merchant/profile/files        # 上传商户文件 (调用FileService)
POST /api/v1/exhibition-service/merchant/notifications/send   # 发送商户通知 (调用NotificationService)
GET /api/v1/exhibition-service/merchant/users/{id}/auth-verify # 验证用户身份 (调用AuthService)
```

## 🔧 ExhibitionService 核心功能模块

### 1. 直播系统

##### 功能特性
- **实时直播**
- **弹幕系统**
- **连麦功能**
- **数据统计**

##### ExhibitionService 接口
```
# 直播控制
POST /api/v1/exhibition-service/live/start                    # 开始直播
POST /api/v1/exhibition-service/live/stop                     # 结束直播
GET /api/v1/exhibition-service/live/{id}/status               # 获取直播状态

# 互动功能
POST /api/v1/exhibition-service/live/{id}/like                # 点赞
POST /api/v1/exhibition-service/live/{id}/comment             # 发送弹幕
GET /api/v1/exhibition-service/live/{id}/comments             # 获取弹幕列表
POST /api/v1/exhibition-service/live/{id}/connect             # 连麦申请
POST /api/v1/exhibition-service/live/{id}/disconnect          # 断开连麦

# 数据统计
GET /api/v1/exhibition-service/live/{id}/viewers              # 获取观看人数
GET /api/v1/exhibition-service/live/{id}/likes                # 获取点赞数
GET /api/v1/exhibition-service/live/{id}/comments-count       # 获取弹幕数
GET /api/v1/exhibition-service/live/{id}/shares               # 获取分享数
```

##### 与独立服务交互
```
POST /api/v1/exhibition-service/live/{id}/notifications       # 发送直播通知 (调用NotificationService)
GET /api/v1/exhibition-service/live/{id}/stream-url           # 获取直播流地址 (调用FileService)
```

### 2. 审核系统

##### 功能特性
- **入驻审核**
- **展会申请审核**
- **内容审核**

##### ExhibitionService 接口
```
# 申请管理
GET /api/v1/exhibition-service/applications                   # 获取申请列表
POST /api/v1/exhibition-service/applications                  # 提交申请
GET /api/v1/exhibition-service/applications/{id}              # 获取申请详情
PUT /api/v1/exhibition-service/applications/{id}              # 更新申请

# 审核流程
POST /api/v1/exhibition-service/applications/{id}/approve     # 审核通过
POST /api/v1/exhibition-service/applications/{id}/reject      # 审核驳回
GET /api/v1/exhibition-service/applications/{id}/history      # 获取审核历史
```

##### 与独立服务交互
```
GET /api/v1/exhibition-service/applications/{id}/files        # 获取申请文件 (调用FileService)
POST /api/v1/exhibition-service/applications/{id}/notify      # 发送审核通知 (调用NotificationService)
GET /api/v1/exhibition-service/applications/{id}/auth-check   # 验证申请人身份 (调用AuthService)
```

### 3. 消息推送系统

##### 功能特性
- **实时通知**
- **消息分类**
- **推送管理**

##### ExhibitionService 接口
```
# 消息管理
GET /api/v1/exhibition-service/notifications/list              # 获取消息列表
POST /api/v1/exhibition-service/notifications/send             # 发送通知
GET /api/v1/exhibition-service/notifications/templates         # 获取推送模板
POST /api/v1/exhibition-service/notifications/subscribe        # 订阅通知
POST /api/v1/exhibition-service/notifications/unsubscribe     # 取消订阅
```

##### 与独立服务交互
```
POST /api/v1/exhibition-service/notifications/push             # 推送消息 (调用NotificationService)
GET /api/v1/exhibition-service/notifications/channels          # 获取推送渠道 (调用NotificationService)
```

### 4. 文件管理系统

##### 功能特性
- **文件上传**
- **文件存储**
- **文件管理**

##### ExhibitionService 接口
```
# 文件管理
GET /api/v1/exhibition-service/files/list                     # 获取文件列表
GET /api/v1/exhibition-service/files/{id}                     # 获取文件信息
DELETE /api/v1/exhibition-service/files/{id}                  # 删除文件
```

##### 与独立服务交互
```
POST /api/v1/exhibition-service/files/upload                  # 上传文件 (调用FileService)
GET /api/v1/exhibition-service/files/{id}/download            # 下载文件 (调用FileService)
GET /api/v1/exhibition-service/files/{id}/preview             # 预览文件 (调用FileService)
```

## 🌐 多语言支持

##### 功能特性
- **双语切换**：本国语言 + 英语
- **动态语言包**
- **用户语言偏好**

##### ExhibitionService 接口
```
GET /api/v1/exhibition-service/i18n/languages                 # 获取支持的语言列表
GET /api/v1/exhibition-service/i18n/translations              # 获取翻译内容
PUT /api/v1/exhibition-service/user/language                  # 设置用户语言偏好
```

## 🔒 权限管理

##### 功能特性
- **角色权限**
- **功能权限**
- **数据权限**

##### ExhibitionService 接口
```
GET /api/v1/exhibition-service/permissions/roles              # 获取角色列表
GET /api/v1/exhibition-service/permissions/{role}/functions   # 获取角色功能权限
PUT /api/v1/exhibition-service/permissions/{role}/functions   # 设置角色功能权限
GET /api/v1/exhibition-service/user/permissions               # 获取用户权限
```

##### 与独立服务交互
```
GET /api/v1/exhibition-service/auth/verify                    # 验证用户身份 (调用AuthService)
POST /api/v1/exhibition-service/auth/check-permission         # 检查用户权限 (调用AuthService)
```

## 📊 数据统计

##### 功能特性
- **实时数据**
- **历史统计**
- **报表生成**

##### ExhibitionService 接口
```
GET /api/v1/exhibition-service/stats/exhibitions              # 展会统计
GET /api/v1/exhibition-service/stats/live-streams            # 直播统计
GET /api/v1/exhibition-service/stats/users                    # 用户统计
GET /api/v1/exhibition-service/stats/merchants                # 商户统计
GET /api/v1/exhibition-service/reports/generate               # 生成报表
```

## 🚀 ExhibitionService 部署说明

### 环境要求
- Node.js 16+
- MongoDB 4.4+
- Redis 6.0+
- Nginx 1.18+

### 服务依赖
- **AuthService** (身份认证服务)
- **NotificationService** (消息推送服务)
- **FileService** (文件服务)

### 安装步骤
1. 克隆项目代码
2. 安装依赖包
3. 配置环境变量
4. 启动数据库服务
5. 配置独立服务连接
6. 运行ExhibitionService

### 配置说明
- 数据库连接配置
- Redis缓存配置
- 独立服务连接配置
- API网关配置

## 📝 开发规范

### 代码规范
- 使用ESLint进行代码检查
- 遵循RESTful API设计规范
- 统一的错误处理机制
- 完善的日志记录

### 测试规范
- 单元测试覆盖率 > 80%
- 集成测试
- 端到端测试

## 🔄 版本历史

### v1.0.0 (2024-10-14)
- 初始版本发布
- 基础功能实现
- 多角色权限管理
- 直播系统集成
- 微服务架构设计

---

*本文档基于XMind思维导图分析整理，详细功能实现请参考具体代码文档。*
