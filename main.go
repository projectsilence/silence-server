package main

import (
	"fmt"
	"bufio"
	"log"
	"net"
	//"strings"
)

//Function to see if IP adress attempting connection is in the allowed IP list.
func Find(slice []string, val string) (int, bool) {
    for i, item := range slice {
        if item == val {
            return i, true
        }
    }
    return -1, false
}

//Function to handle connections
func connectionHandler(conn net.Conn) {
	defer conn.Close()

    //List of allowed IPs...... EDIT TO INCLUDE MORE IPS
    allowedIps := []string{"127.0.0.1"}

    //Setup message scanner and finds IP address
	scanner := bufio.NewScanner(conn)
    if addr, ok := conn.RemoteAddr().(*net.TCPAddr); ok {
        fmt.Println("Connection attempt from: ", addr.IP.String())

        //Checks if IP is in allowed IP list using Find function from above
        _, found := Find(allowedIps, addr.IP.String())
        if !found {
            fmt.Println("Connection from ", addr.IP.String(), " dissalowed..")
            conn.Write([]byte("You aren't in the allowed IPs list"))
            conn.Close()
        } else {
            conn.Write([]byte("Enter desired username..\n> "))
		    username := scanner.Text()
            conn.Write([]byte("Enter desired password..\n> "))
            password := scanner.Text()
            fmt.Printf(username, password)
	        }
        }
    

    //Error handling
	if err := scanner.Err(); err != nil {
        fmt.Println("error:", err)
    }
}


//Main function which starts server
func main() {
    connection, err := net.Listen("tcp", "127.0.0.1:8081")
    if err != nil {
        log.Fatal(err)
    }
	fmt.Println("[ Silence Server ] - Serving on port 8081")
    for {
        conn, err := connection.Accept()
        if err != nil {
            log.Fatal(err)
        }
        go connectionHandler(conn)
    }
}
