package main

import (
    "flag"
    "log"
    "net/http"
    "os"

    "github.com/google/go-github/github"
    "github.com/parkr/auto-reply/ctx"
    "github.com/parkr/auto-reply/hooks"
)

var context *ctx.Context

func main() {
    var port string
    flag.StringVar(&port, "port", "8080", "The port to serve to")
    flag.Parse()
    context = ctx.NewDefaultContext()

    http.HandleFunc("/_ping", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("Content-Type", "text/plain")
        w.Write([]byte("ok\n"))
    }))

    // Add your event handlers. Check out the documentation for the
    // github.com/parkr/auto-reply/hooks package to see all supported events.
    eventHandlers := hooks.EventHandlerMap{}

    // Build the affinity handler.
    aff := &affinity.Handler{}
    aff.AddRepo("myorg", "myproject")
    aff.AddTeam(context, 123, "Performance", "@myorg/performance")
    aff.AddTeam(context, 456, "Documentation", "@myorg/documentation")

    // Add the affinity handler's various event handlers to the event handlers map :)
    eventHandlers.AddHandler(hooks.IssuesEvent, aff.AssignIssueToAffinityTeamCaptain)
    eventHandlers.AddHandler(hooks.IssueCommentEvent, aff.AssignIssueToAffinityTeamCaptainFromComment)
    eventHandlers.AddHandler(hooks.PullRequestEvent, aff.AssignPRToAffinityTeamCaptain)

    // Create the webhook handler. GlobalHandler takes the list of event handlers from
    // its configuration and fires each of them based on the X-GitHub-Event header from
    // the webhook payload.
    myOrgHandler := hooks.GlobalHandler{
        Context:       context,
        EventHandlers: eventHandlers,
    }
    http.HandleFunc("/_github/myproject", myOrgHandler)

    log.Printf("Listening on :%s", port)
    log.Fatal(http.ListenAndServe(":"+port, nil))
}
