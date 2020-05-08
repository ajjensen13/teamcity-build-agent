package buildtools

import (
	"bytes"
	"context"
	"encoding/csv"
	"fmt"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"
)

const (
	use   = `scrapbook`
	short = `Generate value.yaml files from local docker images.`
	long  = `Generate value.yaml files from local docker images.

The following example generates a values.yaml file that maps
image.tag to the digest value of the most recently build image
named docker.io/example-image:12345

	teamcity tool scrapbook --value image.tag=docker.io/example-image:12345
`
)

var scrapbookCmd = &cobra.Command{
	Use:   use,
	Short: short,
	Long:  long,
	Run:   run,
}

const (
	valuesFlagName     = "value"
	valuesUsage        = `A value to include in the scrapbook. Example: yaml.key=docker/repository:tag`
	labelsFlagName     = "label"
	labelsUsage        = `Label used to filter images by. Example: build=12345`
	defaultHandlerName = "digest"
)

func init() {
	rootCmd.AddCommand(scrapbookCmd)
	scrapbookCmd.Flags().StringArrayP(valuesFlagName, string(valuesFlagName[0]), nil, valuesUsage)
	scrapbookCmd.Flags().StringArrayP(labelsFlagName, string(labelsFlagName[0]), nil, labelsUsage)
}

func run(cmd *cobra.Command, _ []string) {
	values, err := cmd.Flags().GetStringArray(valuesFlagName)
	if err != nil {
		log.Panic(fmt.Errorf("error getting values from arguments: %w", err))
	}

	if len(values) == 0 {
		log.Fatalf("error getting values from arguments: did you forget to specify any with --%s", valuesFlagName)
	}

	labels, err := cmd.Flags().GetStringArray(labelsFlagName)
	if err != nil {
		log.Panic(fmt.Errorf("error getting labels from arguments: %w", err))
	}

	var filters []imageFilter
	for _, label := range labels {
		filters = append(filters, labelFilter(label))
	}

	builder := yamlBuilder{}
	for ndx, spec := range values {
		s := strings.Split(spec, "=")
		var key, repository, handler, tag string
		switch len(s) {
		case 0:
			log.Fatalf("missing key in value %d: %q", ndx, spec)
		case 1:
			log.Fatalf("missing repository in value %d: %q", ndx, spec)
		case 2:
			key = s[0]
			repository, tag = splitRepoAndTag(s[1])
			handler = defaultHandlerName
		case 3:
			key = s[0]
			repository, tag = splitRepoAndTag(s[1])
			handler = s[2]
		}

		err := builder.provideHandler(cmd.Context(), handler)(cmd.Context(), key, repository, tag, filters)
		if err != nil {
			log.Fatal(fmt.Errorf("error while building value for spec %s [%d]: %w", spec, ndx, err))
		}
	}

	err = yaml.NewEncoder(os.Stdout).Encode(builder)
	if err != nil {
		log.Fatal(fmt.Errorf("error while generating values: %w", err))
	}
}

func splitRepoAndTag(repoAndTag string) (repo, tag string) {
	s := strings.Split(repoAndTag, ":")
	switch len(s) {
	case 1:
		return s[0], ""
	case 2:
		return s[0], s[1]
	default:
		log.Fatal(fmt.Errorf("failed to parse repository string: %v", repoAndTag))
	}
	//noinspection GoUnreachableCode
	panic("unreachable")
}

type yamlBuilder map[string]interface{}

func (v yamlBuilder) add(key, value string) {
	s := strings.SplitN(key, ".", 2)
	root := s[0]

	switch len(s) {
	case 1:
		v[root] = value
	default:
		if _, ok := v[root]; !ok {
			v[root] = yamlBuilder{}
		}
		v[root].(yamlBuilder).add(s[1], value)
	}
}

func (v yamlBuilder) provideHandler(ctx context.Context, handler string) func(ctx context.Context, key, repository, tag string, filters []imageFilter) error {
	switch handler {
	case "digest":
		return func(ctx context.Context, key, repository, tag string, filters []imageFilter) error {
			dockerImages, err := images(ctx, repository)
			if err != nil {
				return err
			}

			// This extra filtering is necessary because tag filtering messes up digest reporting
			// https://github.com/moby/moby/issues/29901
			dockerImages = filterByTag(tag, dockerImages)

			digest, err := latestDigest(ctx, dockerImages)
			if err != nil {
				return err
			}

			v.add(key, digest)
			return nil
		}
	default:
		log.Fatalf("unknown handler %q", handler)
	}
	//noinspection GoUnreachableCode
	panic("unreachable")
}

func filterByTag(tag string, dockerImages []dockerImage) []dockerImage {
	if tag == "" {
		return dockerImages
	}

	var result []dockerImage
	for _, image := range dockerImages {
		if image.Tag == tag {
			result = append(result, image)
		}
	}
	return result
}

type dockerImage struct {
	ID           string
	Repository   string
	Tag          string
	Digest       string
	CreatedSince string
	CreatedAt    time.Time
	Size         string
}

type imageFilter string

func labelFilter(filter string) imageFilter {
	return imageFilter(fmt.Sprintf(`label=%s`, filter))
}

func images(ctx context.Context, repository string, filters ...imageFilter) ([]dockerImage, error) {
	const formatTemplate = `{{ .ID }},{{ .Repository }},{{ .Tag }},{{ .Digest }},{{ .CreatedSince }},{{ .CreatedAt }},{{ .Size }},`
	arguments := []string{"images", repository, "--format", formatTemplate, "--digests"}
	for _, filter := range filters {
		arguments = append(arguments, "--filter", string(filter))
	}

	cmd := exec.CommandContext(ctx, "docker", arguments...)
	// fmt.Fprintf(os.Stderr, "executing command: %v\n", cmd)

	var output []byte
	output, err := cmd.Output()
	if err != nil {
		return nil, err
	}

	cr := csv.NewReader(bytes.NewReader(output))
	rs, err := cr.ReadAll()
	if err != nil {
		return nil, err
	}

	result := make([]dockerImage, len(rs))

	for i, r := range rs {
		_ = r[6] // compiler hint to avoid bounds checks below

		createdAt, err := time.Parse(`2006-01-02 15:04:05 -0700 MST`, r[5])
		if err != nil {
			return nil, err
		}

		result[i] = dockerImage{
			ID:           r[0],
			Repository:   r[1],
			Tag:          r[2],
			Digest:       r[3],
			CreatedSince: r[4],
			CreatedAt:    createdAt,
			Size:         r[6],
		}
	}

	return result, nil
}

func latestDigest(ctx context.Context, dockerImages []dockerImage) (string, error) {
	switch len(dockerImages) {
	case 1:
		return dockerImages[0].Digest, nil
	case 0:
		return "", fmt.Errorf("image not found")
	default:
		result := &dockerImages[0]
		for _, image := range dockerImages {
			if result.CreatedAt.Before(image.CreatedAt) {
				result = &image
			}
		}
		return result.Digest, nil
	}
}
