package main

import (
	"context"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"database/sql"
	"encoding/pem"
	"fmt"
	"log"
	"os"
	"time"

	_ "github.com/go-sql-driver/mysql"

	//"crypto/md5"
	_ "os/exec"
)

const (
	username = "root"
	password = "temppassword"
	hostname = "localhost:3306" //DON'T USE A PUBLIC IP UNLESS YOU KNOW WHAT YOU'RE DOING
	dbname   = "basicsilenceserver"
)

func dsn(dbName string) string {
	return fmt.Sprintf("%s:%s@tcp(%s)/%s", username, password, hostname, dbName)
}

func main() {
	log.Printf("[ Silence Server ] - Generating servers keys..")
	privateKey, err := rsa.GenerateKey(rand.Reader, 4096)
	if err != nil {
		fmt.Printf("Error generating RSA Keys..\n")
		os.Exit(1)
	}

	publicKey := &privateKey.PublicKey
	var privateKeyBytes []byte = x509.MarshalPKCS1PrivateKey(privateKey)
	privateKeyBlock := &pem.Block{
		Type:  "RSA PRIVATE KEY",
		Bytes: privateKeyBytes,
	}

	pemFilePrivate, err := os.Create("private.pem")
	if err != nil {
		fmt.Printf("Error creating private.pem: %s \n", err)
		os.Exit(1)
	}

	err = pem.Encode(pemFilePrivate, privateKeyBlock)
	if err != nil {
		fmt.Printf("Error writing private.pem: %s \n", err)
		os.Exit(1)
	}

	publicKeyBytes, err := x509.MarshalPKIXPublicKey(publicKey)
	if err != nil {
		fmt.Printf("Error creating publickey: %s \n", err)
		os.Exit(1)
	}

	publicKeyBlock := &pem.Block{
		Type:  "RSA PUBLIC KEY",
		Bytes: publicKeyBytes,
	}

	pemFilePublic, err := os.Create("public.pem")
	if err != nil {
		fmt.Printf("Error creating public.pem: %s \n", err)
		os.Exit(1)
	}
	err = pem.Encode(pemFilePublic, publicKeyBlock)
	if err != nil {
		fmt.Printf("Error writing public.pem: %s \n", err)
		os.Exit(1)
	}

	log.Printf("[ Silence Server ] - Finished setting up keys")
	log.Printf("[ Silence Server ] - Setting up MySQL Database")

	db, err := sql.Open("mysql", dsn(""))
	if err != nil {
		log.Printf("Error %s when opening DB\n", err)
		return
	}
	defer db.Close()

	ctx, cancelfunc := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()
	res, err := db.ExecContext(ctx, "CREATE DATABASE IF NOT EXISTS "+dbname)
	if err != nil {
		log.Printf("Error %s when creating DB\n", err)
		return
	}
	no, err := res.RowsAffected()
	if err != nil {
		log.Printf("Error %s when fetching rows", err)
		return
	}
	log.Printf("[ Silence Server ] - rows affected %d\n", no)

	db.Close()
	db, err = sql.Open("mysql", dsn(dbname))
	if err != nil {
		log.Printf("Error %s when opening DB", err)
		return
	}
	defer db.Close()

	db.SetMaxOpenConns(20)
	db.SetMaxIdleConns(20)
	db.SetConnMaxLifetime(time.Minute * 5)

	ctx, cancelfunc = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()
	err = db.PingContext(ctx)
	if err != nil {
		log.Printf("Errors %s pinging DB", err)
		return
	}
	log.Printf("[ Silence Server ] - Connected to DB %s successfully\n", dbname)

	ctx, cancelfunc = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()
	res, err = db.ExecContext(ctx, "CREATE TABLE IF NOT EXISTS Users(id int primary key auto_increment, username text, password text, date datetime default CURRENT_TIMESTAMP);")
	if err != nil {
		log.Printf("%s  when creating DB table\n", err)
		return
	}
	rows, err := res.RowsAffected()
	if err != nil {
		log.Printf("Error %s when getting rows affected", err)
		return
	}
	log.Printf("[ Silence Server ] - Rows affected when creating table: %d", rows)

	ctx, cancelfunc = context.WithTimeout(context.Background(), 5*time.Second)
	defer cancelfunc()
	res, err = db.ExecContext(ctx, "INSERT INTO Users(username, password) VALUES ('default_user', 'd86e9c0f8f252fd60ac057a29cf8c814');")
	if err != nil {
		log.Printf("Error %s when adding default user\n", err)
		return
	}

}
