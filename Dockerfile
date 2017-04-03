FROM qnib/uplain-golang

COPY bin/qcollect-ng /usr/local/bin/
COPY lib/ /opt/qcollect-ng/lib/
COPY resources/qcollect-ng.yml /etc/

CMD ["qcollect-ng", "--config=/etc/qcollect-ng.yml", "--ld-path=/opt/qcollect-ng/lib/"]

