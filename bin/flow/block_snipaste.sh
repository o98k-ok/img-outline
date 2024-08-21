#!/bin/bash
#
#
#
#
#
# check log config
close_code="6"
limit=10000 # 1000 seconds


real_code=$(cat ~/.snipaste/config.ini| grep -i '\[Log\]' -A 1 | grep 'level=' | cut -c 7-)
if [ "${real_code}" == ${close_code} ]
then
	echo "log code is not match"
	exit 1
fi


# record start ts
before=$(date +%s)


# start sync shot
/Applications/Snipaste.app/Contents/MacOS/Snipaste snip


# wait until logs show finish 
# please check enable logs of snip
while true
do
	sleep .1

	let a++
	if [ $a -gt $limit ]
	then
		echo "shot timeout"
		exit 1
	fi

  after=$(tail -n 200 ~/.snipaste/splog.txt | grep 'about to quit' | tail -n1 | sed 's/\[//g;s/\]//g'  | cut -c 1-19 | xargs -I tt date -j -f "%Y-%m-%d %H:%M:%S" tt +%s)
	if [ $after -gt $before ]
  then
		echo "shot ok"
		exit 0
  fi
done
