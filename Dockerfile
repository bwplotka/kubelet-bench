FROM ubuntu:21.10

ADD ./kubelet /bin/kubelet

ENTRYPOINT [ "/bin/kubelet" ]