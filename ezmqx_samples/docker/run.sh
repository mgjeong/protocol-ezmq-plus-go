#!/bin/sh

export LD_LIBRARY_PATH=../ezmqx_extlibs
cd ezmqx_samples

if [ "topicdiscovery" = $1 ]; then
    echo "start topicdiscovery: $2"
	./topicdiscovery -t $2
elif [ "publisher" = $1 ]; then
	echo "start publisher: $2"
	./publisher -t $2
elif [ "amlsubscriber" = $1 ]; then
    echo "start subscriber with topic: $2"
	./amlsubscriber -t $2 -h true
elif [ "xmlsubscriber" = $1 ]; then
    echo "start subscriber with topic: $3"
	./xmlsubscriber -t $2 -h true
else
	echo "Wrong arguments!!!"
	echo "Examples:"
	echo " publisher topic"
	echo " amlsubscriber topic"
	echo " xmlsubscriber topic"
	echo " topicdiscovery topic"
fi



