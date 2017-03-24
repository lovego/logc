FROM 192.168.202.12:5000/xiaomei/baseimage

USER ubuntu
RUN mkdir /home/ubuntu/logs
RUN mkdir /home/ubuntu/logc /home/ubuntu/logs/logc
WORKDIR /home/ubuntu/logc

VOLUME api_logs
#VOLUME recommend_logs

COPY logc ./
RUN ln -s ../logs/logc log

EXPOSE 3000
ENTRYPOINT [ "./logc" ]
