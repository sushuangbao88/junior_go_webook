create database webook;
use webook;

#添加[用户]表
CREATE TABLE `users` (
  `id` bigint NOT NULL AUTO_INCREMENT,
  `email` varchar(255) DEFAULT NULL COMMENT '电子邮箱',
  `password` varchar(255) NOT NULL DEFAULT '' COMMENT '密码',
  `nickname` varchar(63) NOT NULL DEFAULT '' COMMENT '昵称',
  `birthday` varchar(15) NOT NULL DEFAULT '' COMMENT '生日',
  `gender` enum('1','2') DEFAULT '1' COMMENT '性别：1，男性；2，女性',
  `phone` varchar(31) NOT NULL DEFAULT '' COMMENT '手机号',
  `create_dt` bigint NOT NULL DEFAULT '0' COMMENT '创建时间，毫秒时间戳',
  `update_dt` bigint NOT NULL DEFAULT '0' COMMENT '更新时间，毫秒时间戳',
  `delete_dt` bigint NOT NULL DEFAULT '0' COMMENT '删除时间，毫秒时间戳',
  PRIMARY KEY (`id`),
  UNIQUE KEY `uni_users_email` (`email`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_0900_ai_ci;