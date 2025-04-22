package webscan

import (
	"context"
	"fmt"
	"slack-wails/lib/util"
	"strings"
	"testing"

	nuclei "github.com/projectdiscovery/nuclei/v3/lib"
	"github.com/projectdiscovery/nuclei/v3/pkg/output"
	syncutil "github.com/projectdiscovery/utils/sync"
)

func TestNucleiCaller(t *testing.T) {
	// proxys := []string{"http://127.0.0.1:8080"}
	ne, err := nuclei.NewNucleiEngineCtx(context.Background(),
		nuclei.WithTemplatesOrWorkflows(nuclei.TemplateSources{
			Templates: []string{util.HomeDir() + "/slack/config/pocs"},
		}), // -t
		nuclei.DisableUpdateCheck(), // -duc
		// nuclei.WithProxy(proxys, false), // -proxy
	)
	if err != nil {
		panic(err)
	}
	// load targets and optionally probe non http/https targets
	ne.LoadTargets([]string{}, false)
	err = ne.ExecuteWithCallback(func(event *output.ResultEvent) {
		fmt.Printf("[%s] [%s] %s\n", event.TemplateID, event.Info.SeverityHolder.Severity.String(), event.Matched)
		if event.Info.Reference != nil && !event.Info.Reference.IsEmpty() {
			fmt.Printf("Reference: %s\n", event.Info.Reference.ToSlice())
		}
		fmt.Printf("ExtractedResults: %s\n", strings.Join(event.ExtractedResults, ","))
	})
	if err != nil {
		panic(err)
	}
	defer ne.Close()
}

func TestThreadSafeNucleiCaller(t *testing.T) {
	ctx := context.Background()
	// when running nuclei in parallel for first time it is a good practice to make sure
	// templates exists first

	// create nuclei engine with options
	ne, err := nuclei.NewThreadSafeNucleiEngineCtx(ctx)
	if err != nil {
		panic(err)
	}
	// setup sizedWaitgroup to handle concurrency
	sg, err := syncutil.New(syncutil.WithSize(10))
	if err != nil {
		panic(err)
	}

	// scan 1 = run dns templates on scanme.sh
	sg.Add()
	go func() {
		defer sg.Done()
		err = ne.ExecuteNucleiWithOpts([]string{"scanme.sh"},
			nuclei.WithTemplateFilters(nuclei.TemplateFilters{ProtocolTypes: "dns"}),
			nuclei.WithHeaders([]string{"X-Bug-Bounty: pdteam"}),
			nuclei.EnablePassiveMode(),
		)
		if err != nil {
			panic(err)
		}
	}()

	// scan 2 = run templates with oast tags on honey.scanme.sh
	sg.Add()
	go func() {
		defer sg.Done()
		err = ne.ExecuteNucleiWithOpts([]string{"https://202.88.229.90/"}, nuclei.WithTemplatesOrWorkflows(nuclei.TemplateSources{Templates: []string{"/Users/qwtd/slack/config/pocs/dss-download-fileread.yaml"}}))
		if err != nil {
			panic(err)
		}
	}()

	// wait for all scans to finish
	sg.Wait()
	defer ne.Close()

	// Output:
	// [dns-saas-service-detection] scanme.sh
	// [nameserver-fingerprint] scanme.sh
	// [dns-saas-service-detection] honey.scanme.sh
}
