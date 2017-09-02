package main

//Packages required to run the script. Install the gin before running the script since it is a additional package and it won't come with go.
//Command to install gin - go get github.com/codegangsta/gin
import (
	"bytes"
	"fmt"
	"net/http"
	"os"
	"os/exec"
	"time"
	"strconv"
	"github.com/gin-gonic/gin"
)

func main() {
	//This is the default router provided for gin, which helps us to route the requests to different methods from the url.
	router := gin.Default()

	//healthcheck webservice which returns a 200 status and a blank page.
	router.GET("/healthcheck")

	//create webservice to create an instance in google cloud and return the ip address of the new instance.
	//It uses the username and password given in the request data.
	router.POST("/v1/instances/create", func(c *gin.Context) {
		username := c.PostForm("username")
		password := c.PostForm("password")

		//Content for the temproary file, which is used for poviding the content of startup-script option in gloud command.
		//This enables  PasswordAuthentication and PermitRootLogin which is disabled by default
		//It also creates a user with the given username and assigns it with the given password passed as part of query parameter
		content := "#!/bin/bash\nuseradd -m -s /bin/bash "
		content += username
		content += "\necho "
		content += username
		content += ":"
		content += password
		content += "| chpasswd\nsed -re 's/^(PasswordAuthentication)([[:space:]]+)no/\\1\\2yes/' -i.`date -I` /etc/ssh/sshd_config\nsed -re 's/^(PermitRootLogin)([[:space:]]+)(.*)/\\1\\2yes/' /etc/ssh/sshd_config -i\nservice ssh reload"
		data := []byte(content)

		//Creating the file
		f, err := os.Create("update_sshd.sh")
		//Error handling
		if err != nil {
			fmt.Print(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": fmt.Sprintf("Error occured while creating instance for the user %s", err),
			})
		    return
		}
		//Writing the data
    	output, err := f.Write(data)
    	//Error handling
    	if err != nil {
			fmt.Print(err.Error())
			c.JSON(http.StatusInternalServerError, gin.H{
				"message": fmt.Sprintf("Error occured while creating instance for the user %s", err),
			})
		    return
		}
		fmt.Printf("wrote %d bytes\n", output)

		//Taking the current time in unix format ( Using it for dynamic name of instance. Can be altered if instance name is given as part of query parameter )
  		time_now := strconv.FormatInt(time.Now().UTC().UnixNano(), 10)

  		//Name of the instance to be created
  		name := username
  		name += time_now

  		//Command for creating an instance in google cloud 
    	cmd := "gcloud compute instances create "
    	cmd += name
    	cmd += " --metadata-from-file startup-script='update_sshd.sh' --zone asia-southeast1-a --machine-type=f1-micro"

    	//Executing the command for creating instance
    	command := exec.Command("bash", "-c", cmd)

    	//Storing the result and error if any
    	var out bytes.Buffer
		var stderr bytes.Buffer
		command.Stdout = &out
		command.Stderr = &stderr
		err = command.Run()

		//Deleting the temproary file created
		defer os.Remove("update_sshd.sh") 

		//Error handling
		if err != nil {
		    fmt.Println(fmt.Sprint(err) + ": " + stderr.String())
		    c.JSON(http.StatusUnprocessableEntity, gin.H{
				"message": fmt.Sprintf("Error occured while creating instance for the user %s", stderr.String()),
			})
		    return
		}

		//Command for retrieving the ip address of the newly created isntance
		get_ip := "gcloud compute instances describe "
    	get_ip += name
    	get_ip += " --zone asia-southeast1-a --format='value(networkInterfaces[0].accessConfigs[0].natIP)'"

    	//Executing the command for retrieving the ip
    	ip_out, err := exec.Command("bash", "-c", get_ip).Output()

		//Error handling
		if err != nil {
			fmt.Print("Error occured while retrieving IP: " + err.Error())
		    return
		}

		//Returning the success message along with the newly created instance's ip address
		c.JSON(http.StatusOK, gin.H{
			"message": fmt.Sprintf("Instance IP - %s", ip_out),
		})
	})
	//The service will run in the port specified here
	router.Run(":3000")
}