# Start from Jupyter image
FROM jupyter/scipy-notebook:584e9ab39d22
# Change permissions on data directory (for OpenShift).
# See https://developers.redhat.com/blog/2020/10/26/adapting-docker-and-kubernetes-containers-to-run-on-red-hat-openshift-container-platform/
USER root
RUN chgrp -R 0 /home/jovyan && \
    chmod -R 777 /home/jovyan
# Install Qiskit
USER $NB_UID
RUN pip install qiskit[visualization]
