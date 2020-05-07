package buildtools

import (
	"bytes"
	"context"
	"fmt"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
	"log"
	"os"
	"os/exec"
	"strings"
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
	defaultHandlerName = "digest"
)

func init() {
	rootCmd.AddCommand(scrapbookCmd)
	scrapbookCmd.Flags().StringArrayP(valuesFlagName, string(valuesFlagName[0]), nil, valuesUsage)
}

func run(cmd *cobra.Command, _ []string) {
	values, err := cmd.Flags().GetStringArray(valuesFlagName)
	if err != nil {
		log.Panic(fmt.Errorf("error getting values from arguments: %w", err))
	}

	if len(values) == 0 {
		log.Fatalf("error getting values from arguments: did you forget to specify any with --%s", valuesFlagName)
	}

	builder := new(yamlBuilder)
	for ndx, spec := range values {
		s := strings.Split(spec, "=")
		var key, repository, handler string
		switch len(s) {
		case 0:
			log.Fatalf("missing key in value %d: %q", ndx, spec)
		case 1:
			log.Fatalf("missing repository in value %d: %q", ndx, spec)
		case 2:
			key = s[0]
			repository = s[1]
			handler = defaultHandlerName
		case 3:
			key = s[0]
			repository = s[1]
			handler = s[2]
		}

		err := builder.provideHandler(cmd.Context(), handler)(cmd.Context(), key, repository)
		if err != nil {
			log.Fatal(fmt.Errorf("error while building value for spec %s [%d]: %w", spec, ndx, err))
		}
	}

	err = yaml.NewEncoder(os.Stdout).Encode(builder)
	if err != nil {
		log.Fatal(fmt.Errorf("error while generating values: %w", err))
	}
}

type yamlBuilder yaml.MapSlice

func (v *yamlBuilder) provideHandler(ctx context.Context, handler string) func(ctx context.Context, key, repository string) error {
	switch handler {
	case "digest":
		return v.latestDigest
	default:
		log.Fatalf("unknown handler %q", handler)
	}
	//noinspection GoUnreachableCode
	panic("unreachable")
}

func (v *yamlBuilder) latestDigest(ctx context.Context, key, repository string) error {
	cmd := exec.CommandContext(ctx, "docker", "images", "--format", "{{ .Digest }}", "--digests", repository)

	var output []byte
	output, err := cmd.Output()
	if err != nil {
		return err
	}

	var digest string
	_, err = fmt.Fscanln(bytes.NewReader(output), &digest)
	if err != nil {
		return err
	}

	*v = append(*v, yaml.MapItem{Key: key, Value: digest})

	return nil
}
