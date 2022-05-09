FROM alpine:latest

# 环境变量 KEY0 的内容为 VAL
ENV KEY0 VAL

# 环境变量 KEY1 的内容为 this is key1 env（没有引号）
# 暂不支持处理双引号内容中的命令或环境变量
ENV KEY1 "this is key1 env"

# 环境变量 KEY2 的内容为 this is key2 env（没有引号）
ENV KEY2 'this is key2 env'

# RUN 命令可以使用 \ 和 &&
RUN touch file-from-touch \
&& echo "file content 1" > file-from-touch

# 也可以重复使用 RUN 来声明多个指令
RUN echo "file content 2" >> file-from-touch

# 第一个参数是宿主机目录，第二个参数是容器目录
COPY . /cbox-src

# 支持三个 option，默认值同 docker
# CMD 后面的语法与 RUN 命令相同
HEALTHCHECK --interval=5 --timeout=10 --retries=3 CMD \
echo healthy >> health-check

# ENTRYPOINT 不支持使用 && 来连接多个命令
ENTRYPOINT cat file-from-touch

# 这是一个扩展语句，用于直接声明构建的镜像的名字
NAME my-image