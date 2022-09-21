## road-server pm2 script 
pm2 stop app
rm -rf ./app
go build -o app
pm2 start app
