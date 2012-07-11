/*
 * Copyright (c) 2012 Jason McVetta.  This is Free Software, released under the
 * terms of the WTFPL v2.  It comes without any warranty, express or implied.
 * See http://sam.zoy.org/wtfpl/COPYING for more details.
 * 
 *
 * Static files used this example are derived from the example included with 
 * the jQuery File Upload plugin:
 * https://github.com/blueimp/jQuery-File-Upload
 *
 * Copyright 2011, Sebastian Tschan
 * https://blueimp.net
 *
 * Original software by Tschan licensed under the MIT license:
 * http://www.opensource.org/licenses/MIT
 */


package main

import (
	"github.com/bmizerany/mc"
	"github.com/jmcvetta/jfu"
	"github.com/jmcvetta/mgourl"
	"labix.org/v2/mgo"
	"log"
	"net/http"
	"os"
	"path/filepath"
)

var (
	pathRoot  = os.Getenv("PWD")                // Root of the application install
	pathHtml  = filepath.Join(pathRoot, "html") // HTML documents
	mcServers = os.Getenv("MEMCACHIER_SERVERS")
	mcUser    = os.Getenv("MEMCACHIER_USERNAME")
	mcPasswd  = os.Getenv("MEMCACHIER_PASSWORD")
)

func main() {
	log.SetFlags(log.Ltime | log.Lshortfile)
	//
	// Initialize MongoDB
	//
	mongoUri := os.Getenv("MONGOLAB_URI")
	dbName := "photoshare" // If no DB name specified, use "photoshare"
	if mongoUri != "" {
		switch _, auth, _, err := mgourl.ParseURL(mongoUri); true {
		case err != nil:
			log.Fatal("Could not parse MongoDB URL:", err)
		case auth.Db != "":
			dbName = auth.Db
		}
	} else {
		mongoUri = "localhost"
	}
	//
	// Initialize MongoDB connection
	//
	log.Println("Connecting to MongoDB on", mongoUri)
	conn, err := mgo.Dial(mongoUri)
	if err != nil {
		log.Fatalln(err)
	}
	db := conn.DB(dbName)
	gfs := db.GridFS("test_foobar")
	store := jfu.NewMongoStore(gfs)
	//
	// Initialize Memcache connection
	//
	if mcServers == "" {
		mcServers = "localhost:11211"
		}
	client, err := mc.Dial("tcp", mcServers)
	if err != nil {
		log.Panic(err)
	}
	if mcUser != "" && mcPasswd != "" {
		client.Auth(mcUser, mcPasswd)
		if err != nil {
			log.Panic(err)
		}
	}
	//
	// Initialize UploadHandler
	//
	conf := jfu.DefaultConfig
	conf.MaxFileSize = 100 << 10 // 100kb
	uh := jfu.UploadHandler{
		Prefix: "/jfu",
		Store:  &store,
		Conf:   &conf,
		Cache:  client,
	}
	//
	// Register Handlers
	//
	http.Handle("/jfu", &uh)
	http.Handle("/jfu/", &uh)
	log.Println("Serve JFU")
	path := http.Dir(filepath.Join(pathRoot, "static"))
	http.Handle("/", http.FileServer(path))
	log.Println("Serve files from ", path)
	//
	// Start the webserver
	//
	port := os.Getenv("PORT")
	if port == "" {
		port = "5000"
		}
	addr := "0.0.0.0:" + port
	log.Println("Starting webserver on", addr, "...")
	log.Fatal(http.ListenAndServe(addr, nil))
}
