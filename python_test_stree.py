import json

import requests
import curlify

JSON_H = {'Content-Type': 'application/json'}


def resource_mount():
    data = {
        "target_path": "waimai.ditu.es",
        "resource_type": "resource_host",
        "resource_ids": [1],

    }
    print(data)
    uri = 'http://localhost:8082/api/v1/resource-mount'
    res = requests.post(uri, json=data, headers=JSON_H)
    print(curlify.to_curl(res.request))
    print(res.status_code)
    print(res.text)


def resource_unmount():
    data = {
        "target_path": "waimai.ditu.es",
        "resource_type": "resource_host",
        "resource_ids": [1],

    }
    print(data)
    uri = 'http://localhost:8082/api/v1/resource-unmount'
    res = requests.delete(uri, json=data, headers=JSON_H)
    print(curlify.to_curl(res.request))
    print(res.status_code)
    print(res.text)


def resource_query():
    matcher1 = {

        "key": "stree_app",
        "value": "kafka",
        "type": 1
    }

    matcher2 = {
        "key": "name",
        "value": "genMockResourceHost_host_3",
        "type": 1
    }
    matcher3 = {
        "key": "private_ip",
        "value": "8.*.5.*",
        "type": 3
    }
    matcher4 = {
        "key": "os",
        "value": "amd64",
        "type": 2
    }

    matcher5 = {

        "key": "stree_app",
        "value": "kafka|es",
        "type": 3
    }
    matcher6 = {

        "key": "stree_group",
        "value": "inf",
        "type": 1
    }
    matcherall = {
        "key": "private_ip",
        "value": "8.*",
        "type": 3
    }

    data = {
        "resource_type": "resource_host",
        "labels":
        # [matcher1],
        # [matcher1,matcher4],
            [matcherall],
        # [matcher1,matcher3],
        #     [matcher5, matcher6],
        'target_label': 'cluster'  # 查询分布时才需要

    }
    print(data)
    g_parms = {
        "page_size": 1200,
    }
    uri = 'http://localhost:8082/api/v1/resource-query'
    res = requests.post(uri, json=data, params=g_parms, headers=JSON_H)
    print(curlify.to_curl(res.request))
    print(res.status_code)
    # print(res.text)
    data = res.json().get("result")
    print(len(data))
    for i in data:
        print(i)


def resource_dis():
    matcher1 = {

        "key": "stree_app",
        "value": "kafka",
        "type": 1
    }

    matcher2 = {
        "key": "name",
        "value": "genMockResourceHost_host_3",
        "type": 1
    }
    matcher3 = {
        "key": "private_ip",
        "value": "8.*.5.*",
        "type": 3
    }
    matcher4 = {
        "key": "os",
        "value": "amd64",
        "type": 2
    }

    matcher5 = {

        "key": "stree_app",
        "value": "kafka|es",
        "type": 3
    }
    matcher6 = {

        "key": "stree_group",
        "value": "inf",
        "type": 1
    }

    data = {
        "resource_type": "resource_host",
        "labels":
        # [matcher1],
            [matcher1, matcher3],
        # [matcher5,matcher6],
        # 'target_label': 'cluster'  # 查询分布时才需要
        'target_label': 'region'  # 查询分布时才需要

    }
    print(data)
    g_parms = {
        "page_size": 1200,
    }
    uri = 'http://localhost:8082/api/v1/resource-distribution'
    res = requests.post(uri, json=data, params=g_parms, headers=JSON_H)
    print(curlify.to_curl(res.request))
    print(res.status_code)
    print(res.text)
    # data = res.json().get("result")
    # print(len(data))
    # for i in data:
    #     print(i)


def resource_group():
    data = {
        # "label": "cluster",
        # "label": "stree_app",
        # "label": "stree_product",
        # "label": "stree_group",
        # "label": "private_ip",
        "resource_type": "resource_host",

    }
    print(data)
    uri = 'http://localhost:8082/api/v1/resource-group'
    res = requests.get(uri, params=data, headers=JSON_H)
    print(curlify.to_curl(res.request))
    print(res.status_code)
    print(res.text)


# resource_unmount()
def node_path_add():
    data = {
        "node": "a1.b1.c1"

    }
    print(data)
    uri = 'http://localhost:8082/api/v1/node-path'
    res = requests.post(uri, json=data)
    print(res.status_code)
    print(res.text)


# def node_path_query():
#     data = {
#         "node": "waimai",
#         "query_type":2,
#
#     }
#     print(data)
#     uri = 'http://localhost:8082/api/v1/node-path'
#     res = requests.get(uri, json=data, headers=JSON_H)
#     print(curlify.to_curl(res.request))
#     print(res.status_code)
#     print(res.text)


def node_path_query():
    data = {
        "target_path": "waimai.ditu.es",
        "resource_type": "resource_host",
        "resource_ids": [1],

    }
    data = {"node": "a1", "query_type": 5}
    print(data)
    uri = 'http://localhost:8082/api/v1/node-path'
    res = requests.get(uri, json=data, headers=JSON_H)
    print(curlify.to_curl(res.request))
    print(res.status_code)
    print(res.text)


def log_job_query():
    data = {
        "target_path": "waimai.ditu.es",
        "resource_type": "resource_host",
        "resource_ids": [1],

    }
    uri = 'http://localhost:8082/api/v1/log-job'
    res = requests.get(uri, headers=JSON_H)
    print(curlify.to_curl(res.request))
    print(res.status_code)
    print(res.text)


def log_job_add():
    """


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
    :return:

      - metric_name: log_containerd_total
    metric_help: /var/log/messages 中的 containerd日志 total
    file_path: /var/log/messages
    pattern:  ".*containerd.*"
    func: cnt
    tags:
      level: ".*level=(.*?) .*"
    """
    tag = {
        "level": '''.*level=(.*?) .*'''
    }
    data = {
        "metric_name": "log_nginx_code_max_1",
        "metric_help": "log_abc",
        "file_path": "/var/log/nginx/access.log",
        "pattern": ".*\[code=(.*?)\].*",
        "func": "max",
        "creator": "xiaoyi",
        # "tag_json": json.dumps(tag),

    }
    print(data)
    uri = 'http://192.168.3.200:8082/api/v1/log-job'
    res = requests.post(uri, json=data, headers=JSON_H)
    print(curlify.to_curl(res.request))
    print(res.status_code)
    print(res.text)


def task_add():
    """

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
    """

    data = {
        "script": '''ping  qq.com''',
        # "script": '''#!/bin/bash\nss -ntlp''',
        "account": "root",
        "title": "test04",
        "timeout": 600,
        "hosts": json.dumps(["192.168.3.200", "192.168.3.201"]),
        "creator": "xiaoyi",
        # "action": "kill",
        # "created": 123,
        # "done": "0",
        # "tag_json": json.dumps(tag),

    }
    print(data)
    uri = 'http://192.168.3.200:8082/api/v1/task'
    # uri = 'http://localhost:8082/api/v1/task'
    res = requests.post(uri, json=data, headers=JSON_H)
    print(curlify.to_curl(res.request))
    print(res.status_code)
    print(res.text)


def task_query():
    uri = 'http://localhost:8082/api/v1/task1'
    res = requests.get(uri, headers=JSON_H)
    print(curlify.to_curl(res.request))
    print(res.status_code)
    print(res.text)


def task_kill():
    uri = 'http://localhost:8082/api/v1/kill-task'
    uri = 'http://192.168.3.200:8082/api/v1/kill-task'
    params = {
        "task_id": 4
    }
    res = requests.post(uri, params=params, headers=JSON_H)
    print(curlify.to_curl(res.request))
    print(res.status_code)
    print(res.text)


# resource_dis()
# resource_query()
# resource_group()
# log_job_query()
# task_add()
# task_query()

# node_path_add()
# node_path_query()
# resource_mount()
# resource_unmount()
# log_job_add()
# task_add()
task_kill()