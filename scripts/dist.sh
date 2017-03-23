go build -ldflags "-X main.BuildVersion=`hg id -i` -X main.BuildTime=`date -u '+%Y-%m-%d_%I:%M:%S%p'` -X main.BuildType=`echo $AJBUILD`" -o dist/ajournal .

