package importer

import (
	"fmt"
	"io/ioutil"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"cgt.name/pkg/go-mwclient"
	"github.com/sirupsen/logrus"
)

const URL = "https://www.bta3062.com/api.php"

func Import(wikidata string, dryrun bool, username, password string) error {
	w, err := mwclient.New(URL, "")
	if err != nil {
		return err
	}

	if dryrun {
		logrus.Info("doing dry run, will not make alterations")
	}

	err = w.Login(username, password)
	if err != nil {
		if !dryrun {
			return fmt.Errorf("error logging in: %s", err)
		}
	}

	ids := GetExistingGear(w)
	logrus.Infof("existing gear: %d", len(ids))

	wikifiles, err := ioutil.ReadDir(wikidata)
	if err != nil {
		return err
	}

	logrus.Infof("Loading %d wiki files", len(wikifiles))

	var (
		creations = 0
		updates   = 0
		unchanged = 0
	)
	for _, fileinfo := range wikifiles {
		if !strings.HasSuffix(fileinfo.Name(), ".wiki") {
			logrus.Warnf("Skipping non-wiki file %s")
			continue
		}
		pageName := strings.TrimSuffix(fileinfo.Name(), ".wiki")
		ids[pageName] = true

		// check if there is an old page
		content, _, err := w.GetPageByName(pageName)
		for err != nil {
			if strings.Contains(err.Error(), "not found") {
				break
			}
			logrus.Warnf("Error getting page %s, retrying: %s", pageName, err)
			time.Sleep(time.Second)
			content, _, err = w.GetPageByName(pageName)
		}

		if content != "" {
			logrus.Debugf("Page %s already has content", pageName)
		} else {
			creations = creations + 1
			logrus.Infof("Page %s does not yet exist", pageName)
		}
		ids[pageName] = true

		fileBytes, err := ioutil.ReadFile(filepath.Join(wikidata, fileinfo.Name()))
		fileContent := string(fileBytes)
		if err != nil {
			logrus.Errorf("Error reading %s: %s", fileinfo.Name(), err)
			continue
		}

		if strings.TrimSpace(fileContent) == strings.TrimSpace(content) {
			logrus.Infof("%s file content matches existing page content", pageName)
			unchanged = unchanged + 1
		} else {
			if content != "" {
				logrus.Infof("%s page content does not match", pageName)
				updates = updates + 1
			}
			if !dryrun {
				if err := w.Edit(map[string]string{
					"title":   pageName,
					"text":    fileContent,
					"summary": "automated page update",
				}); err != nil {
					logrus.Errorf("Error writing page %s to wiki: %s\n", pageName, err)
				}
			}
		}
	}

	deletions := 0

	logrus.Info("updated pages, deleting unused")
	for id, included := range ids {
		if !included {
			token, err := w.GetToken(mwclient.CSRFToken)
			if err != nil {
				logrus.Warnf("%s", err)
				continue
			}
			if !dryrun {
				logrus.Infof("Deleting page %s", id)
				// delete the page
				_, err = w.Post(map[string]string{
					"action": "delete",
					"reason": "updater determined page no longer in use",
					"title":  id,
					"token":  token,
				})
				for err != nil {
					logrus.Errorf("error deleting page %s (retrying): %s", id, err)
					time.Sleep(1 * time.Second)
					_, err = w.Post(map[string]string{
						"action": "delete",
						"reason": "updater determined page no longer in use",
						"title":  id,
						"token":  token,
					})
				}
			} else {
				logrus.Infof("would delete %s", id)
			}
			deletions = deletions + 1
		}
	}

	if dryrun {
		logrus.Infof(
			"dry run, would have created %d, updated %d, deleted %d, and left %d unchanged",
			creations, updates, deletions, unchanged,
		)
	} else {
		logrus.Infof(
			"created %d, updated %d, deleted %d, and left %d unchanged",
			creations, updates, deletions, unchanged,
		)
	}

	return nil
}

func GetExistingGear(w *mwclient.Client) map[string]bool {
	limit := 200
	offset := 0

	ids := map[string]bool{}

	for {
		parameters := map[string]string{
			"action": "cargoquery",
			"tables": "Gear",
			"fields": "Id",
			"format": "json",
			"limit":  strconv.Itoa(limit),
			"offset": strconv.Itoa(offset),
		}

		resp, err := w.Get(parameters)
		if err != nil {
			fmt.Println(err)
		}

		cargoquery, err := resp.GetObjectArray("cargoquery")
		if err != nil {
			fmt.Println(err)
		}

		if len(cargoquery) == 0 {
			break
		}
		for _, row := range cargoquery {
			id, err := row.GetString("title", "Id")
			if err != nil {
				fmt.Println(err)
			}
			ids[id] = false
		}
		offset = offset + limit
	}

	return ids
}
