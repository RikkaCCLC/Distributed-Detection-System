drop table log_job;


CREATE TABLE `log_job` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `metric_name` varchar(200) DEFAULT NULL,
  `metric_help` varchar(200) DEFAULT NULL,
  `file_path` varchar(200) DEFAULT NULL,
  `pattern` varchar(200) DEFAULT NULL,
  `func` varchar(20) DEFAULT NULL,
  `creator` varchar(100) DEFAULT NULL,
  `tag_json` varchar(200) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_unique_key` (`metric_name`,`tag_json`) USING BTREE COMMENT '唯一索引'
) ENGINE=InnoDB  DEFAULT CHARSET=utf8;


