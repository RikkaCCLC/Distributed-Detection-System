set names utf8;

drop table task_meta;
drop table task_result;

CREATE TABLE `task_meta`
(
    `id`        bigint unsigned NOT NULL AUTO_INCREMENT,
    `title`     varchar(255)    not null default '' COMMENT '标题',
    `account`   varchar(64)     not null COMMENT '脚本执行账号',
    `timeout`   int unsigned    not null default 0  COMMENT '执行超时',
    `hosts_raw` varchar(4096)   not null  COMMENT '执行机器的ip列表json',
    `script`    text            not null COMMENT '执行的脚本',
    `args`      varchar(512)    not null default '' COMMENT '执行的脚本的参数',
    `creator`   varchar(64)     not null default '' COMMENT '创建者',
    `created`   timestamp       not null COMMENT '创建时间',
    `done`      int unsigned    not null COMMENT '任务结束与否的标志位=0未结束，=1结束',
    PRIMARY KEY (`id`),
    KEY (`created`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8;


CREATE TABLE task_result
(
    `id`      bigint unsigned not null AUTO_INCREMENT,
    `task_id` bigint unsigned not null COMMENT '归属于哪个任务',
    `host`    varchar(128)    not null  COMMENT '哪个机器',
    `status`  varchar(32)     not null COMMENT '任务执行的结果eg:success failed ...',
    `stdout`  text COMMENT '标准输出',
    `stderr`  text COMMENT '标准错误',
    UNIQUE KEY (`task_id`, `host`),
    PRIMARY KEY (`id`)
) ENGINE = InnoDB
  DEFAULT CHARSET = utf8;


insert into task_meta(title, account, timeout, hosts_raw, script, args, creator, created, done) values ('test3', 'root', 5, '["192.168.3.200", "192.168.3.201"]', '#!/bin/\nwget -O /etc/yum.repos.d/CentOS-Base.repo http://mirrors.aliyun.com/repo/Centos-7.repo\nyum makecache -y', '', 'nyy', 123, 0);


insert into task_meta(title, account, timeout, hosts_raw, script, args, creator, created, done) values ('test1', 'root', 5, '["192.168.3.200", "192.168.3.201"]', 'date', '', 'nyy', 123, 0);



insert into task_meta(title, account, timeout, hosts_raw, script, args, creator, created, done) values ('test2', 'root', 5, '["192.168.3.200", "192.168.3.201"]', 'free -g', '', 'nyy', 123, 0);
insert into task_meta(title, account, timeout, hosts_raw, script, args, creator, created, done) values ('test3', 'root', 5, '["192.168.3.200", "192.168.3.201"]', '#!/bin/bash\nss -ntlp', '', 'nyy', 123, 0);

