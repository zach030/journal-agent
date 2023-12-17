package main

import (
	"context"
	"fmt"
	"github.com/jomei/notionapi"
	log "github.com/sirupsen/logrus"
	"time"
)

type Notion struct {
	client *notionapi.Client
	PageID string
}

func NewNotionAPI(sk, pid string) *Notion {
	return &Notion{
		client: notionapi.NewClient(notionapi.Token(sk)),
		PageID: pid,
	}
}

func (n *Notion) InsertNote(ctx context.Context, review string) error {
	to := time.Now().Format("2006-01-02")
	from := time.Now().AddDate(0, 0, -7).Format("2006-01-02")
	var emoji notionapi.Emoji = "📌"
	_, err := n.client.Page.Create(ctx, &notionapi.PageCreateRequest{
		Icon: &notionapi.Icon{
			Type:  "emoji",
			Emoji: &emoji,
		},
		Cover: &notionapi.Image{
			Type: notionapi.FileTypeExternal,
			External: &notionapi.FileObject{
				URL: "https://upload.wikimedia.org/wikipedia/commons/6/62/Tuscankale.jpg",
			},
		},
		Parent: notionapi.Parent{
			Type:   notionapi.ParentTypePageID,
			PageID: notionapi.PageID(n.PageID),
		},
		Properties: notionapi.Properties{
			"title": notionapi.TitleProperty{
				ID:   "title",
				Type: "title",
				Title: []notionapi.RichText{
					{
						Type: notionapi.ObjectTypeText,
						Text: &notionapi.Text{
							Content: fmt.Sprintf("From %s to %s", from, to),
						},
					},
				},
			},
		},
		Children: []notionapi.Block{
			notionapi.ParagraphBlock{
				BasicBlock: notionapi.BasicBlock{
					Object: notionapi.ObjectTypeBlock,
					Type:   notionapi.BlockTypeParagraph,
				},
				Paragraph: notionapi.Paragraph{
					RichText: []notionapi.RichText{
						{
							Text: &notionapi.Text{
								Content: review,
							},
						},
					},
					Children: nil,
				},
			},
		},
	})
	if err != nil {
		log.Errorf("create page error: %v", err)
		return err
	}
	return nil
}

func (n *Notion) GetPage(pageID string) {
	p, err := n.client.Page.Get(context.Background(), notionapi.PageID(pageID))
	if err != nil {
		log.Errorf("get page error: %v", err)
		return
	}
	log.Infof("page: %v", p)
}
