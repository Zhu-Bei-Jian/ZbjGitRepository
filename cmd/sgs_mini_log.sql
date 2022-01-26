
-- ----------------------------
-- Table structure for sgs_logs
-- ----------------------------
DROP TABLE IF EXISTS `sgs_logs`;
CREATE TABLE `sgs_logs` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `game_id` int(11) NOT NULL,
  `server_id` int(11) NOT NULL,
  `area_id` int(11) NOT NULL,
  `login_from` int(11) NOT NULL,
  `log_time` timestamp NOT NULL DEFAULT CURRENT_TIMESTAMP,
  `user_account` varchar(64) NOT NULL,
  `user_level` int(11) DEFAULT '0',
  `op_type` int(11) NOT NULL DEFAULT '0',
  `param1` bigint(20) NOT NULL DEFAULT '0',
  `param2` bigint(20) NOT NULL DEFAULT '0',
  `log_info` text DEFAULT NULL,
  `opmark` int(11) NOT NULL DEFAULT '0',
  PRIMARY KEY (`id`)
) ENGINE=InnoDB AUTO_INCREMENT=1 DEFAULT CHARSET=utf8;
