package github

import (
	"context"
	"encoding/json"
	"fmt"
	"io"

	githubv56 "github.com/google/go-github/v56/github"
	v2 "github.com/phoban01/ocm-v2/api/v2"
	"github.com/phoban01/ocm-v2/api/v2/configr"
	"github.com/phoban01/ocm-v2/api/v2/mutate"
	"golang.org/x/oauth2"
)

func (r *repository) WriteBlob(_ context.Context, acc v2.Access) (v2.Access, error) {
	return nil, nil
}

func (r *repository) Write(ctx context.Context, component v2.Component) error {
	auth, err := r.auth.Authorization()
	if err != nil {
		return err
	}

	ts := oauth2.StaticTokenSource(&oauth2.Token{AccessToken: auth.AccessToken})

	tc := oauth2.NewClient(context.Background(), ts)

	client := githubv56.NewClient(tc)

	desc, err := component.Descriptor()
	if err != nil {
		return err
	}

	path := fmt.Sprintf("%s/%s", desc.Name, desc.Version)

	resources, err := component.Resources()
	if err != nil {
		return err
	}

	commitMsg, err := r.config.Get("commit.message")
	if err == configr.ErrNotFound {
		commitMsg = "[auto] add blob"
	} else if err != nil {
		return err
	}

	visitedResources := make([]v2.Resource, len(resources))
	for i, item := range resources {
		if item.Deferrable() {
			visitedResources[i] = item
			continue
		}

		acc, err := item.Access()
		if err != nil {
			return err
		}

		dig, err := acc.Digest()
		if err != nil {
			return err
		}

		data, err := acc.Data()
		if err != nil {
			return err
		}

		content, err := io.ReadAll(data)
		if err != nil {
			return err
		}

		fpath := fmt.Sprintf("%s/blobs/%s", path, dig.Value)

		// Check if file already exists
		fileContent, _, _, _ := client.Repositories.GetContents(ctx, r.owner, r.repo, fpath,
			&githubv56.RepositoryContentGetOptions{},
		)

		if fileContent == nil {
			_, _, err = client.Repositories.CreateFile(ctx, r.owner, r.repo, fpath,
				&githubv56.RepositoryContentFileOptions{
					Message: &commitMsg,
					Content: content,
				},
			)
			if err != nil {
				return err
			}
		} else {
			_, _, err = client.Repositories.UpdateFile(ctx, r.owner, r.repo, fpath,
				&githubv56.RepositoryContentFileOptions{
					Message: &commitMsg,
					Content: content,
					SHA:     fileContent.SHA,
				},
			)
			if err != nil {
				return err
			}
		}

		visitedResources[i] = item
	}

	component = mutate.WithResources(component, visitedResources...)

	desc, err = component.Descriptor()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(desc, " ", " ")
	if err != nil {
		return err
	}

	cdpath := fmt.Sprintf("%s/component-descriptor.json", path)

	fileContent, _, _, _ := client.Repositories.GetContents(ctx, r.owner, r.repo, cdpath,
		&githubv56.RepositoryContentGetOptions{},
	)

	if fileContent == nil {
		_, _, err = client.Repositories.CreateFile(ctx, r.owner, r.repo, cdpath,
			&githubv56.RepositoryContentFileOptions{
				Message: &commitMsg,
				Content: data,
			},
		)
		if err != nil {
			return err
		}
	} else {
		_, _, err = client.Repositories.UpdateFile(ctx, r.owner, r.repo, cdpath,
			&githubv56.RepositoryContentFileOptions{
				Message: &commitMsg,
				Content: data,
				SHA:     fileContent.SHA,
			},
		)
	}

	return nil
}
