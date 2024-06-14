package TestService

import (
	"testing"

	servant "github.com/TAULargeScaleWorkshop/HANA/large-scale-workshop/services/test-service/servant"
)

func TestExtractLinksFromURLServant(t *testing.T) {
	links, err := servant.ExtractLinksFromURL("https://www.microsoft.com", 1)
	if err != nil {
		t.Fatalf("ExtractLinksFromURL failed: %v", err)
	}
	if len(links) == 0 {
		t.Fatalf("ExtractLinksFromURL returned no links")
	}
	t.Logf("Returned links: %v\n", links)
}
