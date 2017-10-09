# Kraken

Package Deps....

checkout be41f5093e2b05c7a0befe35b04b715eb325ab43 of apiextensions-apiserver
removed the apiextensions-apiserver vendor directory
removed the examples directory to avoid conflict

checkout v4.0.0 of client-go

checkout release-1.7 of apimachinery

go get github.com/lib/pq
go get github.com/fatih/color
go get github.com/Sirupsen/logrus
go get github.com/evanphx/json-patch
go get github.com/gorilla/websocket
go get github.com/gorilla/mux
go get github.com/spf13/cobra
go get github.com/spf13/viper

cd src/github.com/spf13/cobra
git checkout a3cd8ab85aeba3522b9b59242f3b86ddbc67f8bd

# below not longer needed

#checkout 4b8fc5be9b77d91bbb6525d18591c43699a2b4e5 of k8s.io/api
#
#git clone https://github.com/kubernetes/api.git


