FROM alpine:3.10

RUN apk --no-cache add curl  
# download SAP root certificates
RUN curl -sSL -f -k http://aia.pki.co.sap.com/aia/SAPNetCA_G2.crt -o /usr/local/share/ca-certificates/SAPNetCA_G2.crt && \
    curl -sSL -f -k http://aia.pki.co.sap.com/aia/SAP%20Global%20Root%20CA.crt -o /usr/local/share/ca-certificates/SAP%20Global%20Root%20CA.crt && \
    update-ca-certificates

# Install Cloud Foundry cli
RUN curl -sSL -f -k 'https://cli.run.pivotal.io/stable?release=linux64-binary&source=github' -o /tmp/cf-cli.tgz
RUN mkdir -p /usr/local/bin && \
  tar -xzf /tmp/cf-cli.tgz -C /usr/local/bin && \
  cf --version && \
  rm -fv /tmp/cf-cli.tgz

COPY multiapps-plugin.linux64 multiapps-plugin.linux64
RUN cf install-plugin ./multiapps-plugin.linux64 -f
ENV I_ACCEPT_THE_RISK_OF=CONFLICT 
ENV I_WANT_TO_STREAM_THE_MTA_FROM=FILE
CMD [  
  echo 'USAGE: mkfifo mtaPipe.mtar; I_ACCEPT_THE_RISK_OF=CONFLICT I_WANT_TO_STREAM_THE_MTA_FROM=FILE cf deploy mtaPipe.mtar & '
  echo '...... curl https://github.wdf.sap.corp/raw/xs2ds/XSOQTests/master/test_resources/health-check/anatz-staticfile.mtar >> mtaPipe.mtar'
  ]
