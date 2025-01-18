FROM python:3.13-alpine

RUN apk add --no-cache git

RUN pip install hapm==0.3.0

CMD ["hapm"]
