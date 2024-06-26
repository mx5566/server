From golang:1.21 as BUILD
# 设置工作目录
WORKDIR /app

# 拷贝源代码
COPY ./ ./

ARG BINNAME=server
RUN sh build.sh ${BINNAME}

EXPOSE 13000

# VOLUME /app/release
# WORKDIR /app/release

# binaries
FROM alpine:latest

WORKDIR /release

ARG BINNAME=server
ENV TYPE=gate

COPY --from=BUILD /app/release/${BINNAME}   ./

CMD if [ "${TYPE}" = "gate" ]; then sh start_gate.sh; elif [ "${TYPE}" = "game" ]; then sh start_game.sh; elif [ "${TYPE}" = "world"]; then sh start_world.sh; fi
