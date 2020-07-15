# twitter-graphql

Installation
------------

1) Install Mysql in your system. You can find the installation guide here ```https://dev.mysql.com/doc/mysql-osx-excerpt/5.7/en/osx-installation-pkg.html```

2) Set the ```MYSQL_PASSWORD``` and ```MYSQL_ROOT_PASSWORD``` under ```environment:``` in ```docker-compose.yaml``` file 

```
    database:
            
            environment:
                  MYSQL_DATABASE: twitter
                  MYSQL_USER: new-user
                  MYSQL_PASSWORD: <your root password>
                  MYSQL_ROOT_PASSWORD: <your root password>
                  
```

Running The Application Using Docker Container
-----------------------------------------------

1) Go to root directory. i.e ```twitter-graphql```
2) Execute command ```docker-compose up --build``` in terminal

The command will build the application from Dockerfile. You can eleminate the ```--build``` option to avoid building everytime you run your application.

Testing the GraphQl API
-----------------------

1) Go to ``` http://localhost:8080/``` and enter the Graphql queries and mutations you want to test.


NOTES:
------

1) The above application will bind the GraphQl server to port ```:8080```. Make sure the port is not already binded.
2) Do not run the mySQL server manually. The docker container will automatically start the server and bind the application to it.
