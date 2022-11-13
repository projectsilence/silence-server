# Silence

Silence - an E2E encrypted private messaging service

Regaining your silence
  

## Setting Up a Server
### Initial Setup

```console

# clone repo

$ git clone https://github.com/projectsilence/silence-server

  

# cd into directory

$ cd silence-server

  

# start mysql service and change root password

$ sudo apt install mariadb-server mariadb-client # if not already installed

$ sudo service mysql start

$ sudo mysql -u root -p

MariaDB [(none)]> ALTER USER 'root'@'localhost' IDENTIFIED BY 'your_desired_password';

  

# edit configuration of setup.go

$ nano setup.go

# CHANGE THIS to your desired database

  

const (

username = "root" # Only change if you know what you're doing

password = "temppassword" # Password for root account setup above

hostname = "localhost:3306" # Host for database to be created

dbname = "silenceserver" # Name you want for the database

)

  

# run the script

$ go run setup.go

```

### Running the server (server.go still a work in progress)

```console

# ( Optional ) Build the server - not recommended due to "allowed IPs"

$ go build server.go

# Run the serevr

$ go run server.go

```

## WARNING - Only run a server on a secure, dedicated server/machine for maximum security.
