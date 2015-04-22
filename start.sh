#!/bin/bash

cat <<EOF > /etc/nginx/conf.d/default.conf
server {
	listen 80;
	listen [::]:80 default_server ipv6only=on;

	root $IMGSERV_HOME/htdocs/;
	index index.html;
    error_log /tmp/nginx_debug.log debug;

	# Make site accessible from http://localhost/
	server_name localhost;

	location / {
		# First attempt to serve request as file, then
		# as directory, then fall back to displaying a 404.
		# try_files \$uri \$uri/ /index.html;
		# Uncomment to enable naxsi on this location
		# include /etc/nginx/naxsi.rules

        location ~* ^/?$ {
            try_files /index.html =404;
        }

        #FIXME: this actually allows /root/anything/view/test.jpg
        #where root is the location directive we're in, and anything can be anything
        #it should be subdirectories of root only
        location ~* /view/(.*\.(png|jpg|gif))$ {
            try_files   /../images/\$1 =404;
        }

        location ~* (/[^\/]*\.(html|css|js))$ {
            try_files   \$1 =404;
        }

        proxy_pass http://127.0.0.1:5000/;    
    }
}
EOF

#imgserv > /var/log/imgserv.log
#nginx -g "daemon off;"
nginx
imgserv
