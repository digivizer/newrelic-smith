package main

import (
	"testing"
)

func TestExtractAgentTotalCount(t *testing.T) {
	input := `DashboardFacebookIngressAgent
DashboardFacebookPostsSchedulerAgent
DashboardFacebookProfileSchedulerAgent
DashboardFilterAgent
DashboardHandleHistoriesAgent
DashboardIngressAgent
DashboardInstagramIngressAgent
DashboardInstagramPostsSchedulerAgent
DashboardInstagramUserSchedulerAgent
DashboardPostPersistAgent
DashboardTwitterIngressAgent
DashboardTwitterTweetsSchedulerAgent
DashboardYoutubeChannelSchedulerAgent`

	actual, _ := extractAgentTotalCount(input)
	expected := 13

	if actual != expected {
		t.Errorf("extractStatus was incorrect, got: %d, expected: %d", actual, expected)
	}
}

func TestExtractStatus(t *testing.T) {
	input := "running  b0d48aa6-e2d8-4990-a41a-3b8b689ca3ad  13917  2017/09/20 16:38:55    FacebookCommentsDownloadAgent"
	actual, _ := extractStatus(input)
	expected := "running"

	if actual != expected {
		t.Errorf("extractStatus was incorrect, got: %s, expected: %s", actual, expected)
	}
}

func TestPopulateInventory(t *testing.T) {
	// Insert here the logic for your tests
	actual := 2
	expected := 2
	if actual != expected {
		t.Errorf("PopulateInventory was incorrect, got: %d, expected: %d", actual, expected)
	}
}

func TestPopulateMetrics(t *testing.T) {
	// Insert here the logic for your tests
	actual := "foo"
	expected := "foo"
	if actual != expected {
		t.Errorf("PopulateMetrics was incorrect, got: %s, expected: %s", actual, expected)
	}
}
