package main

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/jenkins-x/go-scm/scm"
	"github.com/jenkins-x/jx-helpers/v3/pkg/scmhelpers"
	"github.com/jenkins-x/jx-logging/v3/pkg/log"
	"github.com/pkg/errors"
)

type Options struct {
	httpClient *http.Client
	scmhelpers.PullRequestOptions
	Result *scm.PullRequest
}

func newOptions() (*Options, error) {
	return &Options{httpClient: &http.Client{}}, nil
}

func main() {
	o, err := newOptions()
	if err != nil {
		log.Logger().Fatalf("failed to validate options: %v", err)
	}

	err = o.PullRequestOptions.Validate()
	if err != nil {
		log.Logger().Fatalf("failed to validate: %v", err)
	}

	o.Result, err = o.DiscoverPullRequest()
	if err != nil {
		log.Logger().Fatalf("failed to discover pull request: %v", err)
	}

	joke, err := o.getChuckNorrisJoke()
	if err != nil {
		log.Logger().Fatalf("failed to get a chuck norris joke: %v", err)
	}
	log.Logger().Infof("about to comment on pr joke: %s", joke)

	err = o.commentPullRequest(joke)
	if err != nil {
		log.Logger().Fatalf("failed to comment on pull request: %v", err)
	}
	log.Logger().Infof("successfully commented")
}

func (o *Options) commentPullRequest(joke string) error {
	ctx := context.Background()
	comment := &scm.CommentInput{Body: joke}
	_, _, err := o.ScmClient.PullRequests.CreateComment(ctx, o.FullRepositoryName, o.Number, comment)
	prName := "#" + strconv.Itoa(o.Number)
	if err != nil {
		return errors.Wrapf(err, "failed to comment on pull request %s on repository %s", prName, o.FullRepositoryName)
	}
	log.Logger().Infof("commented on pull request %s on repository %s", prName, o.FullRepositoryName)

	return nil

}

func (o *Options) getChuckNorrisJoke() (string, error) {
	// hardcoding dev else explicit jokes can come through and they are not appropriate
	resp, err := o.httpClient.Get("https://api.chucknorris.io/jokes/random?category=dev")
	if err != nil {
		return "", errors.Wrapf(err, "failed to query chuck norris api")
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", errors.Wrapf(err, "failed to read body")
	}
	r := Response{}
	err = json.Unmarshal(body, &r)
	if err != nil {
		return "", errors.Wrapf(err, "failed to unmarshal response")
	}
	return r.Joke, nil
}

type Response struct {
	Joke string `json:"value"`
}
