package app

import (
	"context"
	"fmt"
	"log"

	dlp "cloud.google.com/go/dlp/apiv2"
	dlppb "google.golang.org/genproto/googleapis/privacy/dlp/v2"
)

var err error

type DlpClass struct {
	Client *dlp.Client
}

func (d *DlpClass) getInfoType() []*dlppb.InfoType {
	return []*dlppb.InfoType{
		{
			Name: "PERSON_NAME",
		},
		{
			Name: "US_STATE",
		},
		{
			Name: "STREET_ADDRESS",
		},
		{
			Name: "LOCATION",
		},
	}
}

func (d *DlpClass) getResult(resp *dlppb.InspectContentResponse, includeQuote bool) {
	findings := resp.GetResult().GetFindings()
	if len(findings) == 0 {
		log.Println("No findings.")
	} else {
		log.Println("Findings:")
		for _, f := range findings {
			if includeQuote {
				log.Println("\tQuote: ", f.GetQuote())
			}
			log.Println("\tInfo type: ", f.GetInfoType().GetName())
			log.Println("\tLikelihood: ", f.GetLikelihood())
		}
	}
}

func (d *DlpClass) Scan(input string, projectID string) {
	log.Printf("Scanning keyword: %v in project: %v ...\n", input, projectID)
	ctx := context.Background()

	// Creates a DLP client.
	d.Client, err = dlp.NewClient(ctx)
	if err != nil {
		log.Fatalf("error creating DLP client: %v", err)
	}
	defer d.Client.Close()

	// The minimum likelihood required before returning a match.
	minLikelihood := dlppb.Likelihood_POSSIBLE
	// The maximum number of findings to report (0 = server maximum).
	maxFindings := int32(0)
	// Whether to include the matching string.
	includeQuote := true
	// The infoTypes of information to match.
	infoTypes := d.getInfoType()

	// Construct item to inspect.
	item := &dlppb.ContentItem{
		DataItem: &dlppb.ContentItem_Value{
			Value: input,
		},
	}

	// Construct request.
	req := &dlppb.InspectContentRequest{
		Parent: fmt.Sprintf("projects/%s/locations/global", projectID),
		InspectConfig: &dlppb.InspectConfig{
			InfoTypes:     infoTypes,
			MinLikelihood: minLikelihood,
			Limits: &dlppb.InspectConfig_FindingLimits{
				MaxFindingsPerRequest: maxFindings,
			},
			IncludeQuote: includeQuote,
		},
		Item: item,
	}

	// Run request.
	resp, err := d.Client.InspectContent(ctx, req)
	if err != nil {
		log.Fatal(err)
	}

	d.getResult(resp, includeQuote)
}
