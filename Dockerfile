FROM debian:jessie

#Install compiler
RUN apt-get update \
	&& apt-get install -y --no-install-recommends git ca-certificates curl \
	&& apt-get clean

ENV GO_VERSION 1.7.1
RUN curl -fsSL "https://storage.googleapis.com/golang/go${GO_VERSION}.linux-amd64.tar.gz" \
	| tar -xzC /usr/local


#Create target dirs
RUN mkdir -p /opt/draftTerm/etc /opt/draftTerm/bin /opt/draftTerm/src


#Compile
ENV PATH /usr/local/go/bin:$PATH
ENV GOPATH /opt/draftTerm

ADD . /opt/draftTerm/src/github.com/mysinmyc/draftTerm

RUN go get -t github.com/mysinmyc/draftTerm/cmd/draftTermd

#Compile generate_cert utility 
RUN cd $GOPATH/bin \
	&& go build /usr/local/go/src/crypto/tls/generate_cert.go


#Configure for container runtime
ADD ./sh/runContainer.sh /opt/draftTerm/bin

RUN useradd -m guest -s /bin/bash \
	&& echo guest:changeIt1! | chpasswd 

EXPOSE 8443

ENTRYPOINT ["/opt/draftTerm/bin/runContainer.sh"]

