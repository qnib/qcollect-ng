FROM qnib/uplain-init

COPY qcollect-ng.yml /etc/
COPY lib/handlers /opt/qcollect-ng/lib/handlers
COPY lib/collectors /opt/qcollect-ng/lib/collectors
COPY lib/filters /opt/qcollect-ng/lib/filters

ADD ./bin/* /usr/local/bin/

CMD ["qcollect-ng", "--config=/etc/qcollect-ng.yml", "--ld-path=/opt/qcollect-ng/lib/"]

