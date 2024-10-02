drop table stree_path;


CREATE TABLE `stree_path` (
  `id` int(11) NOT NULL AUTO_INCREMENT,
  `level` tinyint(4) NOT NULL,
  `path` varchar(200) DEFAULT NULL,
  `node_name` varchar(200) DEFAULT NULL,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_unique_key` (`level`,`path`,`node_name`) USING BTREE COMMENT '唯一索引'
) ENGINE=InnoDB  DEFAULT CHARSET=utf8;


