FROM redis:7-alpine

ENV REDIS_PASSWORD=""

EXPOSE 6379

CMD ["sh", "-c", "exec redis-server --requirepass \"$REDIS_PASSWORD\""]
