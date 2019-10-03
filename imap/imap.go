package imap

import (
	"bytes"
	"fmt"
	"github.com/bilfash/jingu/config"
	"github.com/emersion/go-imap"
	"github.com/emersion/go-imap/client"
	"github.com/DusanKasan/parsemail"
	"github.com/mholt/archiver"
	"io/ioutil"
	"log"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
	"time"
)

func Download(config config.Config ) error {
	hostAddr := fmt.Sprintf("%s:%s", config.Host(), config.Port())
	c, err := client.DialTLS(hostAddr, nil)
	if err != nil {
		log.Fatal(err)
	}

	if err := c.Login(config.Username(), config.Password()); err != nil {
		log.Fatal(err)
	}

	mailboxes := make(chan *imap.MailboxInfo, 10)
	done := make(chan error, 1)
	go func() {
		done <- c.List("", "*", mailboxes)
	}()

	for range mailboxes {
		continue
	}

	if err := <-done; err != nil {
		log.Fatal(err)
	}

	_, err = c.Select(config.Mailbox(), false)
	if err != nil {
		log.Fatal(err)
	}

	criteria := imap.NewSearchCriteria()
	criteria.Text = config.Subjects()

	seqNums, err := c.Search(criteria)
	if err != nil {
		log.Fatal(err)
	}

	seqset := new(imap.SeqSet)
	seqset.AddNum(seqNums...)

	messages := make(chan *imap.Message, 15)
	var section imap.BodySectionName
	go func() {
		done <- c.Fetch(seqset, []imap.FetchItem{section.FetchItem(), imap.FetchEnvelope}, messages)
	}()

	sourceFolder := config.SinkFolder()
	for msg := range messages {
		literal := msg.GetBody(&section)
		email, err := parsemail.Parse(literal)
		if err != nil {
			log.Fatal(err)
		}

		for _, a := range (email.Attachments) {
			buf := new(bytes.Buffer)
			buf.ReadFrom(a.Data)

			isZip, _ := regexp.MatchString(".zip$", a.Filename)
			if isZip {
				processZip(a.Filename, buf, email.Subject, a, config)
				continue
			}

			now := time.Now()
			folderPath := fmt.Sprintf("%s/%d-%d/%d", sourceFolder, now.Year(), now.Month(), now.Day())
			os.MkdirAll(folderPath, os.ModePerm)
			fileName := fmt.Sprintf("%s/%s_%s", folderPath, email.Subject, a.Filename)
			_, err := os.Stat(fileName)
			if os.IsExist(err) {
				continue
			}

			isMatch, _ := regexp.MatchString(config.FilePattern(), a.Filename)
			if isMatch {
				err := ioutil.WriteFile(fileName, buf.Bytes(), 0644)
				if err != nil {
					log.Fatal(err)
				}
			}
		}
	}

	if err := <-done; err != nil {
		log.Fatal(err)
	}

	return nil
}

func processZip(fileName string, buf *bytes.Buffer, subject string, a parsemail.Attachment, config config.Config) {
	tmpFolder := "./tmp"
	os.MkdirAll(tmpFolder, os.ModePerm)
	defer os.RemoveAll(tmpFolder)
	tmpZipPath := fmt.Sprintf("%s/%s", tmpFolder, fileName)
	err := ioutil.WriteFile(tmpZipPath, buf.Bytes(), 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove(tmpZipPath)

	zipFolder := fmt.Sprintf("%s/%s",tmpFolder, strings.ReplaceAll(a.Filename, ".zip", "") + strconv.FormatInt(time.Now().Unix(), 10))
	err = archiver.Unarchive(tmpZipPath, zipFolder)
	if err != nil {
		log.Fatal(err)
	}
	defer os.RemoveAll(zipFolder)

	files := make([]string, 0)
	paths := make([]string, 0)
	filepath.Walk(zipFolder, func(path string, info os.FileInfo, err error) error {
		if !info.IsDir() {
			files = append(files, info.Name())
			paths = append(paths, path)
		}
		if err != nil {
			log.Fatal(err)
		}
		return nil
	})

	now := time.Now()
	for k, p := range paths {
		folderPath := fmt.Sprintf("%s/%d-%d/%d", config.SinkFolder(), now.Year(), now.Month(), now.Day())
		os.MkdirAll(folderPath, os.ModePerm)
		fileName := fmt.Sprintf("%s/%s_%s", folderPath, subject, files[k])
		_, err := os.Stat(fileName)
		if os.IsExist(err) {
			continue
		}

		isMatch, _ := regexp.MatchString(config.FilePattern(), p)
		if isMatch {
			input, err := ioutil.ReadFile(p)
			if err != nil {
				log.Fatal(err)
			}

			err = ioutil.WriteFile(fileName, input, 0644)
			if err != nil {
				log.Fatal(err)
			}
		}
	}
}