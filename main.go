package main

import (
	"context"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"gopkg.in/yaml.v3"
	"os"
)

type Config struct {
	NoteDir  string `json:"note_dir" yaml:"note_dir"`
	APIKey   string `json:"api_key" yaml:"api_key"`
	APIBase  string `json:"api_base" yaml:"api_base"`
	NotionSK string `json:"notion_sk" yaml:"notion_sk"`
	PageID   string `json:"page_id" yaml:"page_id"`
}

func LoadConfig(cfg interface{}, filename string) error {
	var raw []byte
	buf, err := os.ReadFile(filename)
	if err != nil {
		return errors.Wrap(err, "ioutil.ReadFile")
	}
	raw = buf
	if err := yaml.Unmarshal(raw, cfg); err != nil {
		return errors.Wrap(err, "yaml.Unmarshal")
	}
	return nil
}

func main() {
	config := Config{}
	if err := LoadConfig(&config, "config.yaml"); err != nil {
		log.Errorf("load config error: %v", err)
		return
	}
	agent := NewJournalAgent(config.APIKey, config.APIBase, config.NoteDir)
	review, err := agent.NotesReview()
	if err != nil {
		log.Errorf("review error: %v", err)
		return
	}
	notion := NewNotionAPI(config.NotionSK, config.PageID)
	err = notion.InsertNote(context.Background(), review)
	if err != nil {
		log.Errorf("insert note error: %v", err)
		return
	}
	log.Infof("success journal agent")
}
