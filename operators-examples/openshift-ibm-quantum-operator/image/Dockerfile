FROM centos/python-36-centos7:latest



# User setup for OpenShift
USER root
ENV HOME=/ibmq
RUN mkdir -p ${HOME}/ibmq_operator && \
    useradd -u 9000 -r -g 0 -d ${HOME} -s /sbin/nologin \
    -c "IBMQ User" ibmq-user

WORKDIR ${HOME}

# Install dependencies
RUN pip install --upgrade pip && \
    pip install qiskit jupyter matplotlib numpy requests Flask gunicorn

ADD ./ ${HOME}/ibmq_operator
RUN chown -R 9000:0 ${HOME} && \
    find ${HOME} -type d -exec chmod g+ws {} \;
WORKDIR ${HOME}/ibmq_operator

# Run the notebook
USER 1001
RUN ls -la

CMD gunicorn wsgi:application -b 0.0.0.0:5000
EXPOSE 5000
