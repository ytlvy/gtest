package http

import (
	"fmt"
	"net/http"

	"gopkg.in/go-playground/webhooks.v5/github"
)

type WebHook struct {
}

func (w *WebHook) Start() {
	hook, _ := github.New(github.Options.Secret("MyGitHubSuperSecretSecrect...?"))

	http.HandleFunc("/webhook", func(w http.ResponseWriter, r *http.Request) {
		payload, err := hook.Parse(r, github.ReleaseEvent, github.PullRequestEvent)
		if err != nil {
			if err == github.ErrEventNotFound {

			}
		}
		switch payload.(type) {

		case github.ReleasePayload:
			release := payload.(github.ReleasePayload)

			fmt.Printf("%+v", release)

		case github.PushPayload:
			pullRequest := payload.(github.PushPayload)

			fmt.Printf("%+v", pullRequest)
		}
	})
	http.ListenAndServe(":3000", nil)
}
