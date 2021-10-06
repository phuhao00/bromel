



# CH



docker run -d --name ch-server --ulimit nofile=262144:262144 -p 8123:8123 -p 9000:9000 -p 9009:9009 yandex/clickhouse-server





1.docker exec -it  docker-clickhouse /bin/bash  进入容器

2.clickhouse-client 进入clickhouse命令行

3.show databases 查看所有的数据库

4.clickhouse 允许远程访问，将clickhouse的配置文件拷贝出来
 docker cp clickhouse-server:/etc/clickhouse-server/ /etc/clickhouse-server/

5.修改 /etc/clickhouse-server/config.xml 中 65行 注释去掉<listen_host>::</listen_host>

6.用自定义配置文件启动容器
 docker run -d --name docker-clickhouse --ulimit nofile=262144:262144 -p 8123:8123 -p 9000:9000 -p 9009:9009 -v /etc/clickhouse-server/config.xml:/etc/clickhouse-server/config.xml yandex/clickhouse-server
 端口必须映射出来，不然阿里云的远程访问不到端口

