package main

import (
	"io"
	"net"
	"fmt"
	"time"
	"slices"
	"errors"
	"strconv"
	"net/http"
	"github.com/charmbracelet/log"
	"github.com/Supraboy981322/gomn"
)

var (
	blkDB bool
	db gomn.Map
	port = 8944
	clDBSec int
	clToDef bool
	dbPath string
	logLvl string
	config gomn.Map
	adminPermIP []string
	bin map[string][]byte
	configPath = "conf.gomn"
)

func init() {
	var err error

	log.Info("initializing server...")

	//temporarilly use debug mode
	//  (changed when parsing config)
	log.SetLevel(log.DebugLevel)

	//configure the server (clearly)
	if err = configure(); err != nil {
		if err.Error() == "parsing config:  open conf.gomn: no such file or directory" {
			log.Error("config does not exist") 
			defConf()
		} else {
			log.Fatalf("failed to configure:  %v", err)
		}
	}

	//initialize database (clearly)
	if err = initDB(); err != nil {
		log.Fatalf("failed to initialize database:  %v", err)
	}

	//generate database map in memory
	if err = mapDB(); err != nil {
		log.Fatalf("failed to map database:  %v", err)
	}
}

func main() {
	log.Info("started")

	http.HandleFunc("/get", getHan)
	http.HandleFunc("/set", setHan)
	http.HandleFunc("/db", dbHan)
	http.HandleFunc("/del", delHan)
	http.HandleFunc("/dbADMIN", dbAdminHan)
	
	log.Infof("listening on port %d", port)

	if clDBSec > 0 {
		go clDB()
	}

	//looks cleaner 
	portStr := ":"+strconv.Itoa(port)
	log.Fatal(http.ListenAndServe(portStr, nil))
}

func getHan(w http.ResponseWriter, r *http.Request) {
	var ok bool //avoids golang quirk (is bug?)

	//get the key
	var key string //used later, so init outside block
	if key = getKey(r); key == "" {
		//read the req body if no key header
		bod, err := io.ReadAll(r.Body)
		if err != nil {
			eror(w, "reading req body", err)
			return
		};key = string(bod)
		//return err if still no key 
		if key == "" {
			w.Write([]byte("need key\n"))
			return
		}
	}

	//get val from db 
	var val []byte //used later, so init outside block
	if val, ok = db[key].([]byte); !ok {
		//check problem
		if db[key] != nil {//if not []byte, it's invalid 
			err := errors.New(key)
			eror(w, "invalid value in db. key", err)
		}	else {//if empty, assume missing
			err := errors.New(key)
			eror(w, "key does not exist", err)
		}
		return
	}

	//check if header specifies to not use newline
	newline := chkHeaders(r, []string{"n", "newline", `\n`, "\n"})
	if newline == "" || newline == "true" {
		val = append(val, []byte("\n")...)
	}
	
	//self-explainitory
	logReq("/get", r.RemoteAddr, "key="+key)

	//send value
	w.Write(val)
}

func setHan(w http.ResponseWriter, r *http.Request) {
	//get the key header
	var key string//used later, so init outside block
	if key = getKey(r); key == "" {
		if key == "" {
			w.Write([]byte("need key\n"))
			return
		}
	}
	
	//get the value header
	var val string//used later, so init outside block
	if val = getVal(r); val == "" {
		//if no val header, use body
		bod, err := io.ReadAll(r.Body)
		if err != nil {
			eror(w, "reading req body", err)
			return
		};val = string(bod)
		//if body is also empty, return err
		if val == "" {
			w.Write([]byte("need value\n"))
			return
		}
	}
	
	//self-explainitory
	logReq("/set", r.RemoteAddr, "key="+key)

	//set value in db (in memory)
	db[key] = []byte(val)
	w.Write([]byte("added to db\n"))

	updateDB(w) //save in-memory db to disk

	//confirm completion
	//  (updateDB sends progress to client)
	w.Write([]byte("done\n"))
}

func delHan(w http.ResponseWriter, r *http.Request) {
	//get the key
	var key string//used later, so init outside block
	if key = getKey(r); key == "" {
		//use body if header empty
		bod, err := io.ReadAll(r.Body)
		if err != nil {
			eror(w, "reading req body", err)
			return
		};key = string(bod)
		//if no body, err
		if key == "" {
			w.Write([]byte("need key\n"))
			return
		}
	}

	//self-explainitory
	logReq("/del", r.RemoteAddr, "key="+key)

	delete(db, key) //delete from in-memory db
	updateDB(w) //save changes to disk

	//confirm completion
	msg := fmt.Sprintf("deleted:  %s\n", key)
	w.Write([]byte(msg))
}

func dbHan(w http.ResponseWriter, r *http.Request) {
	//get the type header 
	typ := chkHeaders(r, []string{"type", "t", "typ"})
	if typ == "" {
		//check body
		bod, err := io.ReadAll(r.Body)
		if err != nil {
			eror(w, "reading req body", err)
			return
		};typ = string(bod)
		//if no body, err
		if typ == "" {
			typ = "bin"
		}
	}

	//self-explainitory
	logReq("/db", r.RemoteAddr, "typ="+typ)
	
	//stream db to client
	dlBin(w, typ)
}

func dbAdminHan(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	flusher, ok := w.(http.Flusher)
	if !ok { /* send no response */ return }
	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		log.Error(err)
		/* send no response */
		return
	}
	if slices.Contains(adminPermIP, ip) && len(adminPermIP) != 0 {
		flusher.Flush()
		w.Write([]byte("is admin\n"))
		log.Warn("validated admin request, omitting ip for security")
	} else if len(adminPermIP) == 0 {
		w.Write([]byte("no admins set; refusing requested action\n"))
		return
	} else { /* don't even respond */ return }

	//only one header checked for the action
	action := r.Header.Get("A")

	//is case-sensitive
	switch action {
   case "deleteProd()":
		//log it
		//  (there should be something between warn and fatal that isn't an err) 
		log.Warn("ADMIN REQUESTED deleteProd()... stalling")

		//mk list of lines
		lines := []string{
			"\033[1mWELL THAT'S EXTREME\033[0m",
    	"but who am \033[1mI\033[0m to question an \033[1;4;32madmin\033[0m?",
			"maybe this is some last-stitch effort to save the company from a bad db entry",
			"  (probably malware or something illegal)",
			"perhaps your boss (while drunk) ordered you to \033[1;4;5;41mdelete prod\033[0m or get fired",
			"could be that you're trying to sabotage the company",
			"it's possible that you're a hacker who's trying to destroy the server",
			"or, this was an \033[4maccident\033[0m",
			"either way...",
			"you have \033[1;4;31m10 seconds\033[0m before \033[1;4;5;41mdeleteProd()\033[0m is run",
			"i hope you don't take too long to read this.",
			"that would be humorous, and also highly unfortunate",
			"waiting \033[1;4;31m10 seconds\033[0m...",
		}//for each line, print and wait 4 seconds
		for _, line := range lines {
			w.Write([]byte(line+"\n"))
			flusher.Flush()
			time.Sleep(4 * time.Second)
		}
		
		//in order to give them a
		//  mild heart attack, the server
		//    doesn't check that all params
		//      are present and valid until
		//        after the 10 second warning
		mkDefault, err := strconv.ParseBool(r.Header.Get("mkDefault"))
		if r.Header.Get("mkDefault") == "" || err != nil {
			log.Warn("nevermind, admin request invalid; making them sweat instead")
			//it proceeds to print
			//  a blank line, to feel
			//    like something is happening
			w.Write([]byte("\n"))

			//then it waits 8 seconds
			//  to maximize the chances
			//    of a mild heart-attack
			flusher.Flush()
			time.Sleep(8 * time.Second)

			//and lets them know that
			//  their mild heart-attack
			//    was for nothing
			w.Write([]byte("wait, nevermind, your request \033[31misn't valid\033[0m\n"))
  		flusher.Flush()
			w.Write([]byte("\033[1;32maborting action...\033[0m\n"))
		  flusher.Flush()

			//only for them to return 
			//  to their terminal's cursor,
			//    so they can contemplate what
			//      they almost did
			return
		}

		//finally...
		//  do the deed.
		log.Warn("running out of time... can't stall for much longer")
		deleteProd(mkDefault)

		//make them feel bad
		lines = []string{
			"",//start with a blank line 
			"\033[32mcongrats\033[0m",
			"you just \033[1;31mdeleted\033[0m the database",
			"i hope you feel good about yourself.",
			"if you thought this was a joke or an easter-egg...",
			"you have \033[1;4;31mseverely\033[0m \033[1;31mmistaken\033[0m, and you should recover from your backups.",
			"wait, you did make backups, \033[1;4mright?\033[0m",
		}//print and wait 4 seconds for each line
		for _, line := range lines {
			w.Write([]byte(line+"\n"))
			flusher.Flush()
			time.Sleep(4 * time.Second)
		}
	 default:
		log.Error("VALIDATED ADMIN ATTEMPTED INVALID ACTION")
		w.Write([]byte("invalid action\n"))
		return
	}
}
