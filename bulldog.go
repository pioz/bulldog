package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	"github.com/vharitonsky/iniflags"
)

// Implement flag.Value interface
type urls []string

func (u *urls) String() string {
	return fmt.Sprint(*u)
}

func (u *urls) Set(value string) error {
	if len(*u) > 0 {
		return errors.New("urls flag already set")
	}
	for _, url := range strings.Split(value, ",") {
		if url != "" {
			*u = append(*u, url)
		}
	}
	return nil
}

// Config is a struct that contains config information
type Config struct {
	urlFile               string
	sleep, sleepWithError int
	timeout               int
	oneCheck              bool
	logFile               string
	quiet                 bool
	gmail, pass, to       string
	urls                  urls
}

func configure(config *Config) {
	flag.StringVar(&config.urlFile, "f", "", "File that contains the urls to check each per line.")

	flag.IntVar(&config.sleep, "s", 60, "Seconds to sleep between a complete URLs check to another.")
	flag.IntVar(&config.sleepWithError, "se", 600, "Seconds to sleep between a complete URLs check to another if a check has failed.")
	flag.IntVar(&config.timeout, "t", 10, "Check request timeout in seconds.")
	flag.BoolVar(&config.oneCheck, "1", false, "Perform URLs checks only one time.")

	flag.StringVar(&config.logFile, "logfile", "", "Logfile path.")
	flag.BoolVar(&config.quiet, "quiet", false, "Disable logging.")

	flag.StringVar(&config.gmail, "gmail", "", "Gmail account. If this is present send email using the Gmail smtp server. Use -pass flag to specify the Gmail account password. If this flag is empty send email using `mail` command line program.")
	flag.StringVar(&config.pass, "pass", "", "Gmail account password. Only relevant when using -gmail flag.")
	flag.StringVar(&config.to, "to", "", "When a check fails send an email on this email address. If is empty the email alert is disabled.")

	flag.Var(&config.urls, "urls", "Comma-separated list of URLs to check.")

	version := flag.Bool("v", false, "Print version.")
	iniflags.Parse()

	if *version {
		fmt.Println("0.2.0")
		os.Exit(2)
	}

	if config.quiet {
		log.SetOutput(ioutil.Discard)
	} else if config.logFile != "" {
		file, err := os.OpenFile(config.logFile, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0644)
		if err == nil {
			log.SetOutput(file)
		}
	}

	if config.urlFile != "" {
		file, err := os.Open(config.urlFile)
		defer file.Close()
		if err != nil {
			panic("Can not open urlFile '" + config.urlFile + "'")
		}
		scanner := bufio.NewScanner(file)
		for scanner.Scan() {
			config.urls = append(config.urls, scanner.Text())
		}
		if err := scanner.Err(); err != nil {
			panic(err)
		}
	}
}

func check(client *http.Client, u string) error {
	resp, err := client.Get(u)
	if err != nil || resp.StatusCode != 200 {
		var checkError error
		if err != nil {
			checkError = err
		} else {
			checkError = &url.Error{Op: "Get", URL: u, Err: fmt.Errorf("status code is %d", resp.StatusCode)}
		}
		return checkError
	}
	return nil
}

func main() {
	var (
		config Config
		sleep  time.Duration
		mailer *Mailer
		client *http.Client
	)
	configure(&config)
	if len(config.urls) == 0 {
		log.Println("Nothing to check. Exiting...")
		return
	}
	client = &http.Client{Timeout: time.Second * time.Duration(config.timeout)}
	mailer = &Mailer{gmail: config.gmail, pass: config.pass, to: config.to}
	log.Printf("Starting to check these urls => %v...\n", config.urls)
	for {
		unreachable := make([]error, 0)
		for _, url := range config.urls {
			err := check(client, url)
			if err != nil {
				unreachable = append(unreachable, err)
				log.Printf("Error for '%s': %s\n", url, err.Error())
			}
		}
		if len(unreachable) == 0 {
			sleep = time.Second * time.Duration(config.sleep)
		} else {
			mailer.BuildAndSendEmail(unreachable)
			sleep = time.Second * time.Duration(config.sleepWithError)
		}
		if config.oneCheck {
			break
		}
		time.Sleep(sleep)
	}
}
