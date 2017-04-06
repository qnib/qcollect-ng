FROM qnib/uplain-init

COPY qcollect-ng.yml /etc/
COPY ./lib /opt/qcollect-ng/
ADD ./bin/* /usr/local/bin/

CMD ["qcollect-ng", "--config=/etc/qcollect-ng.yml", "--ld-path=/opt/qcollect-ng/lib/"]

