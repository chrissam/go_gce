Webservice to create GCE in Go
=====================
This is a basic go script to create a small Google Compute Engine instance by passing in username and password as the query parameter in the POST request. 

The script will spin a new GCE instance by passing a startup script during boot. The startup script will create a new user and assigns a password from the query parameter and enables root and password authentication by modifying the sshd_config file.


Requirement
==========

 - The script requires go to be installed and the right paths set as part of the $PATH environment variable. Refer [<i class="icon-upload"></i> Go Installation doc](https://golang.org/doc/install) to know more on this
 - Google SDK to be installed and configured in the instance where the script will be run. The user should have necessary permission to create GCE instance.


End Points
=========

The script has two endpoints

 - GET /healthcheck
 - POST /v1/instances/create

The GET end point just responds with a 200 response code and display a blank page to indicate the service is running

The POST end point takes two query parameters **username** and **password**, which will be assigned to the instance for login. 


> **Note:**

> - The script creates a default instance of *f1-micro* family type. Make necessary changes in script if you want to customise instance type, image type, region, zone etc. ( Basically you can enhance it by creating additional query parameters and using it in the script )


Usage
=====

####  **Build**
The first step is the build process where the go script will be compiled and an executable file gets created. 

    cd ~/go/src/go_gce
    
    go build spin_gce.go


####  **Run** 
Start the script by running the below command

    go run spin_gce &


By default the script runs in port 3000. 


####  **GET**

    curl -X GET http://localhost:3000/healthcheck
Replace *localhost* with the corresponding Domain name or IP


####  **POST**

    curl -X POST -F 'username=test' -F 'password=password123' http://localhost:3000/v1/instances/create
Replace *localhost* with the corresponding Domain name or IP



