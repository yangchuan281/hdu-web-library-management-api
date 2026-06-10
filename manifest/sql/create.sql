CREATE DATABASE IF NOT EXISTS library_management DEFAULT CHARSET utf8mb4 COLLATE utf8mb4_unicode_ci;

USE library_management;

-- 用户表
CREATE TABLE `users` (
    `id`         BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '用户ID',
    `name`       VARCHAR(50)     NOT NULL COMMENT '用户姓名',
    `email`      VARCHAR(100)    DEFAULT NULL COMMENT '邮箱',
    `phone`      VARCHAR(20)     DEFAULT NULL COMMENT '手机号（邮箱注册可为空）',
    `password`   VARCHAR(64)     NOT NULL COMMENT '登录密码',
    `role`       VARCHAR(20)     NOT NULL DEFAULT 'user' COMMENT '角色：user-普通用户 admin-管理员',
    `status`     TINYINT         NOT NULL DEFAULT 1 COMMENT '状态：1-正常，0-禁用',
    `created_at` DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_email` (`email`),
    UNIQUE KEY `uk_phone` (`phone`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='用户表';

-- 图书表
CREATE TABLE `books` (
    `id`           BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '图书ID',
    `title`        VARCHAR(200)    NOT NULL COMMENT '图书名称',
    `author`       VARCHAR(100)    NOT NULL COMMENT '作者',
    `isbn`         VARCHAR(20)     NOT NULL COMMENT 'ISBN号',
    `publish_date` DATE            DEFAULT NULL COMMENT '出版日期',
    `status`       TINYINT         NOT NULL DEFAULT 1 COMMENT '状态：1-可借阅，0-已借出，2-下架',
    `created_at`   DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at`   DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    UNIQUE KEY `uk_isbn` (`isbn`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='图书表';

-- 借阅记录表
CREATE TABLE `borrows` (
    `id`         BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '借阅记录ID',
    `user_id`    BIGINT UNSIGNED NOT NULL COMMENT '用户ID',
    `book_id`    BIGINT UNSIGNED NOT NULL COMMENT '图书ID',
    `borrow_at`  DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '借阅时间',
    `return_at`  DATETIME        DEFAULT NULL COMMENT '归还时间',
    `created_at` DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    `updated_at` DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP ON UPDATE CURRENT_TIMESTAMP COMMENT '更新时间',
    PRIMARY KEY (`id`),
    KEY `idx_user_id` (`user_id`),
    KEY `idx_book_id` (`book_id`),
    KEY `idx_borrow_at` (`borrow_at`),
    CONSTRAINT `fk_borrows_user` FOREIGN KEY (`user_id`) REFERENCES `users` (`id`) ON DELETE CASCADE,
    CONSTRAINT `fk_borrows_book` FOREIGN KEY (`book_id`) REFERENCES `books` (`id`) ON DELETE CASCADE
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='借阅记录表';

-- 验证码表（邮箱注册、密码重置）
CREATE TABLE `verification_codes` (
    `id`         BIGINT UNSIGNED NOT NULL AUTO_INCREMENT COMMENT 'ID',
    `email`      VARCHAR(100)    NOT NULL COMMENT '邮箱',
    `code`       VARCHAR(4)      NOT NULL COMMENT '4位数字验证码',
    `type`       VARCHAR(20)     NOT NULL COMMENT '类型：register-注册, reset-重置密码',
    `used`       TINYINT         NOT NULL DEFAULT 0 COMMENT '是否已使用：0-未使用，1-已使用',
    `expires_at` DATETIME        NOT NULL COMMENT '过期时间',
    `created_at` DATETIME        NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '创建时间',
    PRIMARY KEY (`id`),
    KEY `idx_email_type` (`email`, `type`),
    KEY `idx_expires_at` (`expires_at`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COMMENT='验证码表';
