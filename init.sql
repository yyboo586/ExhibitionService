CREATE DATABASE IF NOT EXISTS `exhibition_service` CHARACTER SET utf8mb4 COLLATE utf8mb4_general_ci;

USE `exhibition_service`;

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

CREATE TABLE IF NOT EXISTS `t_file` (
    `id` BIGINT(20) NOT NULL AUTO_INCREMENT COMMENT '文件ID',
    `company_id` VARCHAR(40) DEFAULT '' COMMENT '展会公司ID',
    `type` TINYINT(1) COMMENT '文件类型(1:法人证件照)',
    `file_id` VARCHAR(40) NOT NULL COMMENT '文件ID',
    `file_name` VARCHAR(255) NOT NULL COMMENT '文件名称',
    `file_link` TEXT COMMENT '媒体URL',
    `status` TINYINT(1) NOT NULL DEFAULT 0 COMMENT '文件上传状态(1:上传成功)',
    `create_time` BIGINT(20) NOT NULL COMMENT '创建时间',
    `update_time` BIGINT(20) COMMENT '更新时间',
    PRIMARY KEY (`id`),
    KEY `idx_company_id_type` (`company_id`, `type`),
    UNIQUE KEY `idx_file_id` (`file_id`)
) ENGINE=InnoDB COMMENT='文件表';