/*
 Navicat Premium Data Transfer

 Source Server         : 127
 Source Server Type    : MySQL
 Source Server Version : 100133
 Source Host           : localhost:3306
 Source Schema         : sgs_nuyan

 Target Server Type    : MySQL
 Target Server Version : 100133
 File Encoding         : 65001

 Date: 24/03/2020 14:55:54
*/

SET NAMES utf8mb4;
SET FOREIGN_KEY_CHECKS = 0;

-- ----------------------------
-- 账号表
-- ----------------------------
DROP TABLE IF EXISTS `account`;
CREATE TABLE `account` (
  `userid` bigint(20) UNSIGNED NOT NULL COMMENT '用户id',
  `account_type` tinyint(4) NOT NULL DEFAULT 1 COMMENT '账号类型:1.测试账户,4.微信账户',
  `account` varchar(128) CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL DEFAULT '' COMMENT '账号',
  `loginid` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin NOT NULL DEFAULT '' COMMENT '渠道唯一id',
  `head_img_url` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin NOT NULL DEFAULT '' COMMENT '头像地址',
  `headframe_img_url` varchar(255) CHARACTER SET utf8 COLLATE utf8_bin NOT NULL DEFAULT '' COMMENT '头像框地址',
  `nickname` varchar(128) CHARACTER SET utf8 COLLATE utf8_general_ci NULL DEFAULT NULL COMMENT '昵称',
  `login_ip` varchar(32) CHARACTER SET utf8 COLLATE utf8_unicode_ci NOT NULL DEFAULT '' COMMENT '最近一次登录IP',
  `login_time` timestamp NULL DEFAULT NULL COMMENT '最近一次登录时间',
  `sex` tinyint(4) NOT NULL DEFAULT 1 COMMENT '性别',
  `status` tinyint(4) NOT NULL DEFAULT 0 COMMENT '状态标志位:0.正常,1.封号',
  `unblock_time` bigint(20) NOT NULL DEFAULT 0 COMMENT '账号解封时间',
  `born_time` bigint(20) NOT NULL DEFAULT 0 COMMENT '出生时间（创角时间）',
  `created_ip` varchar(32) CHARACTER SET utf8 COLLATE utf8_unicode_ci NOT NULL DEFAULT '' COMMENT '注册IP地址',
  `created_at` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP COMMENT '注册时间',
  `origin_data` text CHARACTER SET utf8 COLLATE utf8_general_ci NOT NULL COMMENT '注册时初始数据',
  PRIMARY KEY (`userid`) USING BTREE,
  UNIQUE INDEX `account`(`account`) USING BTREE,
  UNIQUE INDEX `account_type_unionid`(`account_type`,`loginid`) USING BTREE
) ENGINE = InnoDB CHARACTER SET = utf8 COLLATE = utf8_general_ci ROW_FORMAT = Dynamic;


-- ----------------------------
-- 用户信息表
-- ----------------------------
DROP TABLE IF EXISTS `user`;
CREATE TABLE `user`  (
  `userid` bigint(20) UNSIGNED NOT NULL AUTO_INCREMENT COMMENT '用户id',
  `level` int(11) NOT NULL DEFAULT 0 COMMENT '等级',
  `exp` int(11) NOT NULL DEFAULT 0 COMMENT '经验',
  `char_info` blob NOT NULL COMMENT '角色信息',
  `client_data` blob NOT NULL COMMENT '客户端自用的数据',
  `prop` longblob NOT NULL COMMENT '道具',
  `goods` longblob NOT NULL COMMENT '物品',
  `task` longblob NOT NULL COMMENT '任务',
  `extend` blob NOT NULL COMMENT '玩家扩展信息',
  `login_time` timestamp NOT NULL  COMMENT '登录时间',
  `logout_time` timestamp NOT NULL  COMMENT '登出时间',
  PRIMARY KEY (`userid`) USING BTREE
) ENGINE = InnoDB AUTO_INCREMENT = 1000000664 CHARACTER SET = utf8 COLLATE = utf8_general_ci ROW_FORMAT = Dynamic;


SET FOREIGN_KEY_CHECKS = 1;
