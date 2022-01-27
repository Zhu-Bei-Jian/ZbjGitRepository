/*
Navicat MySQL Data Transfer

Source Server         : 测试机
Source Server Version : 50622
Source Host           : 10.225.10.34:3306
Source Database       : sgs_10year

Target Server Type    : MYSQL
Target Server Version : 50622
File Encoding         : 65001

Date: 2017-02-23 12:21:32
*/

SET FOREIGN_KEY_CHECKS=0;

-- ----------------------------
-- Table structure for sgs_admin_accounts
-- ----------------------------
DROP TABLE IF EXISTS `admin_account`;
CREATE TABLE `admin_account` (
    `account_id` bigint(20) unsigned NOT NULL AUTO_INCREMENT COMMENT '后台账号ID',
    `account` varchar(60) COLLATE utf8_unicode_ci NOT NULL DEFAULT '0' COMMENT '后台账号',
    `password` varchar(60) COLLATE utf8_unicode_ci NOT NULL DEFAULT '' COMMENT '后台密码',
    `salt` varchar(10) CHARACTER SET utf8 NOT NULL DEFAULT '' COMMENT '盐',
    `rights` text CHARACTER SET utf8 NOT NULL DEFAULT '',
    PRIMARY KEY (`account_id`),
    UNIQUE KEY `account` (`account`) USING BTREE
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci COMMENT='用户后台账号表';

-- ----------------------------
-- Table structure for sgs_admin_users
-- ----------------------------
DROP TABLE IF EXISTS `admin_server_info`;
CREATE TABLE `admin_server_info` (
    `key` char(64) CHARACTER SET utf8 NOT NULL,
    `value` text CHARACTER SET utf8 DEFAULT '',
    PRIMARY KEY (`key`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8 COLLATE=utf8_unicode_ci COMMENT='后台服务参数';

INSERT INTO `sgs_spy`.`admin_account`(`account_id`, `account`, `password`, `salt`, `rights`) VALUES (1, 'a', '23b7451be9c421eb6b0908df6458a609', 'yrvpil6e', '');


SET FOREIGN_KEY_CHECKS=1;

