go build -ldflags "-X main.BuildVersion=`git log --pretty=format:"%H" -n 1` -X main.BuildTime=`date -u '+%Y-%m-%d_%I:%M:%S%p'` -X main.BuildType=`echo $AJBUILD`" -o dist/ajournal .
mkdir -p dist/www
cp -r www/css dist/www/css
cp -r www/fonts dist/www/fonts
cp -r www/images dist/www/images
cp -r www/js dist/www/js
riot www/tags www/js/dist.js
cp www/index.html dist/www/index.html
cp www/app.html dist/www/app.html
cp -r db dist/db
cp -r i18n dist/www/i18n
