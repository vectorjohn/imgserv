#!/bin/bash

#mkdir $IMGSERV_HOME
#mkdir $IMGSERV_HOME/images

cat <<EOF > $IMGSERV_HOME/config.json
{
    "root": "$IMGSERV_HOME/images",
    "max_images": $IMGSERV_MAX_IMAGES,
    "port": 5000
}
EOF
