FROM nginx

ENV IMGSERV_HOME /imgserv
ENV IMGSERV_MAX_IMAGES 2000

ADD ./imgserv /bin/imgserv
ADD . $IMGSERV_HOME/
ADD http://code.jquery.com/jquery-1.11.2.min.js $IMGSERV_HOME/htdocs/jquery.min.js
RUN (chmod 755 $IMGSERV_HOME/htdocs/*)

RUN ${IMGSERV_HOME}/setup.sh

EXPOSE 80
WORKDIR $IMGSERV_HOME
CMD ${IMGSERV_HOME}/start.sh
