package importer

import (
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"time"

	"cgt.name/pkg/go-mwclient"
	"github.com/sirupsen/logrus"
)

// BATCH_SIZE is the number of wiki pages to retrieve at one time.
const BATCH_SIZE = 20

func Import(wikidata string, dryrun bool, username, password, url string) error {
	w, err := mwclient.New(url, "")
	if err != nil {
		return err
	}

	if dryrun {
		logrus.Info("doing dry run, will not make alterations")
	}

	err = w.Login(username, password)
	if err != nil {
		logrus.Warnf("error logging in: %s", err)
		// if !dryrun {
		// return fmt.Errorf("error logging in: %s", err)
		// }
	}

	ids := GetExistingPages(w)
	logrus.Infof("existing pages: %d", len(ids))

	wikifiles, err := ioutil.ReadDir(wikidata)
	if err != nil {
		return err
	}

	logrus.Infof("Loading %d wiki files", len(wikifiles))

	var (
		creations = map[string]struct{}{}
		updates   = map[string]struct{}{}
		unchanged = map[string]struct{}{}
	)
	doBatch := func(wikifiles []os.FileInfo) {
		pages := []string{}
		pageData := map[string]mwclient.BriefRevision{}

		for _, fileinfo := range wikifiles {
			if !strings.HasSuffix(fileinfo.Name(), ".wiki") {
				logrus.Warnf("Skipping non-wiki file %s")
				continue
			}
			pageTitle := strings.TrimSuffix(fileinfo.Name(), ".wiki")
			ids[pageTitle] = true
			// pageName is the pageTitle with the namespace included.
			pageName := fmt.Sprintf("RawData:%s", pageTitle)
			pages = append(pages, pageName)
		}

		pageData, err := w.GetPagesByName(pages...)
		for err != nil {
			if strings.Contains(err.Error(), "not found") {
				break
			}
			logrus.Warnf("Error getting pages, retrying: %s", err)
			time.Sleep(time.Second)
			pageData, err = w.GetPagesByName(pages...)
		}

		for _, fileinfo := range wikifiles {
			if !strings.HasSuffix(fileinfo.Name(), ".wiki") {
				logrus.Warnf("Skipping non-wiki file %s")
				continue
			}
			pageTitle := strings.TrimSuffix(fileinfo.Name(), ".wiki")
			pageName := fmt.Sprintf("RawData:%s", pageTitle)

			// check if there is an old page
			pageRev, ok := pageData[pageName]
			if !ok {
				logrus.Warnf("page %s not in set", pageName)
				continue
			}
			content := pageRev.Content

			if content != "" {
				logrus.Debugf("Page %s already has content", pageName)
			} else {
				creations[pageTitle] = struct{}{}
				logrus.Infof("CREATE %s", pageName)
			}
			ids[pageName] = true

			fileBytes, err := ioutil.ReadFile(filepath.Join(wikidata, fileinfo.Name()))
			fileContent := string(fileBytes)
			if err != nil {
				logrus.Errorf("Error reading %s: %s", fileinfo.Name(), err)
				continue
			}

			if strings.TrimSpace(fileContent) == strings.TrimSpace(content) {
				logrus.Debugf("UNCHANGED %s", pageName)
				unchanged[pageTitle] = struct{}{}
			} else {
				if content != "" {
					logrus.Infof("UPDATE %s", pageName)
					updates[pageTitle] = struct{}{}
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
	}

	// do batches of pages, so we can pull down page info all at once.
	for i := 0; i < len(wikifiles); i = i + BATCH_SIZE {
		j := i + BATCH_SIZE
		if j > len(wikifiles) {
			j = len(wikifiles)
		}
		logrus.Infof(
			"doing batch from %s (%d) to %s (%d) (%d total)",
			wikifiles[i].Name(), i, wikifiles[j-1].Name(), j, len(wikifiles),
		)
		doBatch(wikifiles[i:j])
	}

	deletions := map[string]struct{}{}

	logrus.Info("updated pages, deleting unused, NOT ACTUALLY DELETING RIGHT NOW")
	//for id, included := range ids {
	//	pageName := fmt.Sprintf("RawData:%s", id)
	//	if !included {
	//		token, err := w.GetToken(mwclient.CSRFToken)
	//		if err != nil {
	//			logrus.Warnf("%s", err)
	//			continue
	//		}
	//		if !dryrun {
	//			logrus.Infof("DELETE %s", pageName)
	//			// delete the page
	//			_, err = w.Post(map[string]string{
	//				"action": "delete",
	//				"reason": "updater determined page no longer in use",
	//				"title":  pageName,
	//				"token":  token,
	//			})
	//			for err != nil {
	//				logrus.Errorf("error deleting page %s (retrying): %s", pageName, err)
	//				time.Sleep(1 * time.Second)
	//				_, err = w.Post(map[string]string{
	//					"action": "delete",
	//					"reason": "updater determined page no longer in use",
	//					"title":  pageName,
	//					"token":  token,
	//				})
	//			}
	//		} else {
	//			logrus.Infof("DELETE %s", pageName)
	//		}
	//		deletions[id] = struct{}{}
	//	}
	//}

	if dryrun {
		logrus.Infof(
			"dry run, would have created %d, updated %d, deleted %d, and left %d unchanged",
			len(creations), len(updates), len(deletions), len(unchanged),
		)
	} else {
		logrus.Infof(
			"created %d, updated %d, deleted %d, and left %d unchanged",
			len(creations), len(updates), len(deletions), len(unchanged),
		)
	}

	return nil
}

func GetExistingPages(w *mwclient.Client) map[string]bool {
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
			logrus.Errorf(err.Error())
		}

		cargoquery, err := resp.GetObjectArray("cargoquery")
		if err != nil {
			logrus.Errorf(err.Error())
		}

		if len(cargoquery) == 0 {
			break
		}
		for _, row := range cargoquery {
			id, err := row.GetString("title", "Id")
			if err != nil {
				logrus.Errorf(err.Error())
			}
			ids[id] = false
		}
		offset = offset + limit
	}

	offset = 0
	for {
		parameters := map[string]string{
			"action": "cargoquery",
			"tables": "Chassis",
			"fields": "VariantName,Name",
			"format": "json",
			"limit":  strconv.Itoa(limit),
			"offset": strconv.Itoa(offset),
		}

		resp, err := w.Get(parameters)
		if err != nil {
			logrus.Errorf(err.Error())
		}

		cargoquery, err := resp.GetObjectArray("cargoquery")
		if err != nil {
			logrus.Errorf(err.Error())
		}

		if len(cargoquery) == 0 {
			break
		}
		for _, row := range cargoquery {
			variant, err := row.GetString("title", "VariantName")
			if err != nil {
				logrus.Errorf(err.Error())
			}
			name, err := row.GetString("title", "Name")
			if err != nil {
				logrus.Errorf(err.Error())
			}

			ids[fmt.Sprintf("MechDef_%s_%s", name, variant)] = false
		}
		offset = offset + limit
	}

	return ids
}
