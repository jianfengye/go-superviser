go-superviser
=============

go-superviser for restart a go service

you can use the command:

go-superviser -project=httpserver -run=true

params:

project: project to be supervised, it will find to $GOPATH/src/{$project}

run: whether to run after build the project. 
if it is false or empty, the superviser will only build the project.
if it is true, the superviser will build and run the project.