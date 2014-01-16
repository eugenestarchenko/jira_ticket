package main

//TODO show errors if they happen

import (
    "bytes"
    "encoding/json"
    "flag"
    "fmt"
    "io/ioutil"
    "net/http"
    "os"
)

type Reply struct {
    Id, Key, Self string
}

func main() {

    var description string
    var issuetype string
    var project string
    var should_error int
    var r Reply
    var summary string
    var url string

    // Jira base configuration to be overidden via ENV variables
    jira_url := os.Getenv("JIRA_URL")
    if jira_url == "" {
        url = "https://jira-webhook-test.puppetlabs.com/rest/api/2/issue"
    } else {
        url = os.Getenv("JIRA_URL")
    }

    jira_username := os.Getenv("JIRA_USERNAME")
    if jira_username == "" {
        fmt.Println("You must supply a username for usage via JIRA_USERNAME environment variable")
        os.Exit(1)
    }

    jira_password := os.Getenv("JIRA_PASSWORD")
    if jira_password == "" {
        fmt.Println("You must supply a password for usage via JIRA_PASSSWORD environment variable")
        os.Exit(1)
    }

    flag.StringVar(&project, "project", "RE", "Jira project prefix")
    flag.StringVar(&summary, "summary", "Please do the needful", "Issue summary")
    flag.StringVar(&description, "desc", "More information about your ticket.", "Long form description and details for the issue.")
    flag.StringVar(&issuetype, "issuetype", "Bug", "Bug|Task|Epic|Improvement|Story|Support escalation|New feature")
    flag.Parse()

    if summary == "Please do the needful" {
        fmt.Println("Please provide your own summary.")
        should_error = 1
    }

    if description == "More information about your ticket." {
        fmt.Println("Please provide your own description.")
        should_error = 1
    }

    if should_error == 1 {
        os.Exit(1)
    }

    // A map to then jsonify, perhaps struct would have been better
    ticket := map[string]interface{}{
        "fields": map[string]interface{}{
            "project": map[string]interface{}{
                "key": project,
            },
            "summary":     summary,
            "description": description,
            "issuetype": map[string]interface{}{
                "name": issuetype,
            },
        },
    }

    f, err := json.Marshal(ticket)
    if err != nil {
        fmt.Println("error:", err)
    }

    input_json := []byte(f)
    body := bytes.NewBuffer(input_json)
    client := &http.Client{}

    req, err := http.NewRequest("POST", url, body)
    if err != nil {
        panic("Error while building jira request")
    }
    req.SetBasicAuth(jira_username, jira_password)
    req.Header.Add("Content-Type", "application/json")

    resp, err := client.Do(req)
    defer resp.Body.Close()
    contents, err := ioutil.ReadAll(resp.Body)

    err = json.Unmarshal(contents, &r)
    if err != nil {
        fmt.Println("error:", err)
    }

    fmt.Println()
    fmt.Println("Your new issue is: " + r.Key)
    fmt.Println("https://jira-webhook-test.puppetlabs.com/browse/" + r.Key)
    fmt.Println()
}
