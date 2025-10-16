CREATE DATABASE IF NOT EXISTS `exhibition_service` CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;

USE `exhibition_service`;

CREATE TABLE IF NOT EXISTS `t_file` (
    `id` BIGINT(20) NOT NULL AUTO_INCREMENT COMMENT '文件ID',
    `module` TINYINT(1) DEFAULT 0 COMMENT '模块类型(1:公司,2:服务提供商,3:商户,4:展会)',
    `custom_id` VARCHAR(40) DEFAULT '' COMMENT '自定义ID',
    `type` TINYINT(1) DEFAULT 0 COMMENT '文件类型',
    `file_id` VARCHAR(40) NOT NULL COMMENT '文件ID',
    `file_name` VARCHAR(255) NOT NULL COMMENT '文件名称',
    `file_link` TEXT COMMENT '媒体URL',
    `status` TINYINT(1) DEFAULT 0 COMMENT '文件上传状态(1:上传成功)',
    `create_time` BIGINT(20) NOT NULL COMMENT '创建时间',
    `update_time` BIGINT(20) COMMENT '更新时间',
    PRIMARY KEY (`id`),
    KEY `idx_module_id_type` (`module`, `custom_id`, `type`),
    UNIQUE KEY `idx_file_id` (`file_id`)
) ENGINE=InnoDB COMMENT='文件表';

CREATE TABLE IF NOT EXISTS `t_company` (
    `id` VARCHAR(40) NOT NULL COMMENT '公司ID',
    `name` VARCHAR(32) NOT NULL COMMENT '公司名称',
    `type` TINYINT(1) NOT NULL COMMENT '公司类型(1:服务提供商、2:商户)',
    `country` VARCHAR(255) NOT NULL COMMENT '国家',
    `city` VARCHAR(255) NOT NULL COMMENT '城市',
    `address` VARCHAR(255) NOT NULL COMMENT '地址',
    `email` VARCHAR(100) NOT NULL COMMENT '邮箱',
    `description` TEXT NOT NULL COMMENT '公司简介',
    `version` BIGINT(20) NOT NULL DEFAULT 0 COMMENT '版本号',

    `social_credit_code` VARCHAR(255) NOT NULL COMMENT '统一社会信用代码',
    `legal_person_name` VARCHAR(32) NOT NULL COMMENT '法人姓名',
    `legal_person_card_number` VARCHAR(32) NOT NULL COMMENT '法人证件号',
    
    `create_time` BIGINT(20) NOT NULL COMMENT '创建时间',
    `update_time` BIGINT(20) NOT NULL COMMENT '更新时间',
    PRIMARY KEY(`id`),
    UNIQUE KEY `idx_code` (`social_credit_code`)
)ENGINE=InnoDB COMMENT='公司表';

CREATE TABLE IF NOT EXISTS `t_service_provider` (
    `id` VARCHAR(40) NOT NULL COMMENT '服务提供商ID',
    `company_id` VARCHAR(40) NOT NULL COMMENT '关联的公司ID',
    `name` VARCHAR(100) NOT NULL COMMENT '服务提供商名称',
    `status` TINYINT(1) NOT NULL DEFAULT 0 COMMENT '状态(0:待审核、1:已审核、2:已禁用)',  
    `website` VARCHAR(255) COMMENT '官网',      
    `contact_person_name` VARCHAR(50) COMMENT '联系人姓名',
    `contact_person_phone` VARCHAR(32) COMMENT '联系人电话',
    `contact_person_email` VARCHAR(100) COMMENT '联系人邮箱',
    `description` TEXT COMMENT '服务提供商描述',
    `version` BIGINT(20) NOT NULL DEFAULT 0 COMMENT '版本号',

    `create_time` BIGINT(20) NOT NULL COMMENT '创建时间',
    `submit_for_review_time` BIGINT(20) COMMENT '提交审核时间',
    `approve_time` BIGINT(20) COMMENT '审核通过时间',
    `update_time` BIGINT(20) NOT NULL COMMENT '更新时间',
    PRIMARY KEY(`id`),
    UNIQUE KEY `idx_company_id` (`company_id`)
) ENGINE=InnoDB COMMENT='服务提供商表';

CREATE TABLE IF NOT EXISTS `t_merchant` (
    `id` VARCHAR(40) NOT NULL COMMENT '商户ID',
    `company_id` VARCHAR(40) NOT NULL COMMENT '关联的公司ID',
    `name` VARCHAR(100) NOT NULL COMMENT '商户名称',
    `status` TINYINT(1) NOT NULL DEFAULT 0 COMMENT '状态(0:待审核、1:已审核、2:已禁用)',
    `contact_person_name` VARCHAR(50) COMMENT '联系人姓名',
    `contact_person_phone` VARCHAR(32) COMMENT '联系人电话',
    `contact_person_email` VARCHAR(100) COMMENT '联系人邮箱',    
    `website` VARCHAR(255) COMMENT '商户官网',
    `description` TEXT COMMENT '商户描述',
    `version` BIGINT(20) NOT NULL DEFAULT 0 COMMENT '版本号',

    `create_time` BIGINT(20) NOT NULL COMMENT '创建时间',
    `submit_for_review_time` BIGINT(20) COMMENT '提交审核时间',
    `approve_time` BIGINT(20) COMMENT '审核通过时间',
    `update_time` BIGINT(20) NOT NULL COMMENT '更新时间',
    PRIMARY KEY(`id`),
    KEY `idx_company_id` (`company_id`)
) ENGINE=InnoDB COMMENT='商户表';

CREATE TABLE IF NOT EXISTS `t_exhibition` (
    `id` VARCHAR(40) NOT NULL COMMENT '展会ID',
    `title` VARCHAR(100) NOT NULL COMMENT '展会名称',
    `status` TINYINT(1) NOT NULL DEFAULT 0 COMMENT '展会状态(0:筹备中、1:待审核、2:已批准、3:报名中、4:进行中、5:已结束、6:已取消)',
    `industry` VARCHAR(50) COMMENT '所属行业',
    `tags` TEXT COMMENT '展会标签',
    `website` VARCHAR(255) COMMENT '展会官网',
    `venue` VARCHAR(255) COMMENT '展会地点',
    `venue_address` VARCHAR(255) COMMENT '展会详细地址',
    `country` VARCHAR(255) COMMENT '国家',
    `city` VARCHAR(255) COMMENT '城市',
    `description` TEXT COMMENT '展会详细描述',
    `version` BIGINT(20) NOT NULL DEFAULT 0 COMMENT '版本号',

    `registration_start` BIGINT(20) COMMENT '报名开始时间',
    `registration_end` BIGINT(20) COMMENT '报名结束时间',
    `start_time` BIGINT(20) NOT NULL COMMENT '展会开始时间',
    `end_time` BIGINT(20) NOT NULL COMMENT '展会结束时间',

    `create_time` BIGINT(20) NOT NULL COMMENT '创建时间',
    `submit_for_review_time` BIGINT(20) COMMENT '提交审核时间',
    `approve_time` BIGINT(20) COMMENT '审核通过时间',
    `update_time` BIGINT(20) NOT NULL COMMENT '更新时间',
    PRIMARY KEY(`id`)
) ENGINE=InnoDB COMMENT='展会表';

CREATE TABLE IF NOT EXISTS `t_exhibition_service_provider` (
    `id` BIGINT(20) NOT NULL AUTO_INCREMENT COMMENT '主键',
    `exhibition_id` VARCHAR(40) NOT NULL COMMENT '展会ID',
    `service_provider_id` VARCHAR(40) NOT NULL COMMENT '服务提供商ID',
    `role_type` TINYINT(1) NOT NULL COMMENT '角色类型(1:主办方、2:联合主办方、3:协办方)',
    `create_time` BIGINT(20) NOT NULL COMMENT '创建时间',
    `update_time` BIGINT(20) NOT NULL COMMENT '更新时间',
    PRIMARY KEY (`id`)
) ENGINE=InnoDB COMMENT='展会与服务提供商关联表';


CREATE TABLE IF NOT EXISTS `t_exhibition_stats` (
    `id` BIGINT(20) NOT NULL AUTO_INCREMENT COMMENT '统计ID',
    `exhibition_id` VARCHAR(40) NOT NULL COMMENT '展会ID',
    `stat_date` DATE NOT NULL COMMENT '统计日期',
    `stat_type` TINYINT(1) NOT NULL COMMENT '统计类型(1:日统计、2:周统计、3:月统计)',
    `view_count` INT(11) DEFAULT 0 COMMENT '浏览量',
    `reservation_count` INT(11) DEFAULT 0 COMMENT '预约数',
    `favorite_count` INT(11) DEFAULT 0 COMMENT '收藏数',
    `merchant_count` INT(11) DEFAULT 0 COMMENT '商户数',
    `live_count` INT(11) DEFAULT 0 COMMENT '直播数',
    `total_viewer_count` INT(11) DEFAULT 0 COMMENT '总观看人数',
    `create_time` BIGINT(20) NOT NULL COMMENT '创建时间',
    `update_time` BIGINT(20) NOT NULL COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `idx_exhibition_date_type` (`exhibition_id`, `stat_date`, `stat_type`),
    KEY `idx_stat_date` (`stat_date`),
    KEY `idx_stat_type` (`stat_type`),
    FOREIGN KEY (`exhibition_id`) REFERENCES `t_exhibition`(`id`) ON DELETE CASCADE
) ENGINE=InnoDB COMMENT='展会统计表';