package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	sdkArgs "github.com/newrelic/infra-integrations-sdk/args"
	"github.com/newrelic/infra-integrations-sdk/log"
	"github.com/newrelic/infra-integrations-sdk/metric"
	"github.com/newrelic/infra-integrations-sdk/sdk"
)

type argumentList struct {
	sdkArgs.DefaultArgumentList
}

const (
	integrationName    = "com.digivizer.newrelic-smith"
	integrationVersion = "0.1.0"
)

var (
	args argumentList
)

func extractStatus(smithAgent string) (string, error) {
	splitLine := strings.Split(string(smithAgent), " ")
	metricValue := strings.TrimSpace(splitLine[0])
	return metricValue, nil
}

func getAgentListString() (string, error) {
	agentPath := fmt.Sprintf("%s/%s/*.rb", os.Getenv("SMITH_AGENT_DIRECTORIES"), os.Getenv("SMITH_AGENTS_GROUP"))
	cmd := exec.Command("/bin/bash", "-c", "cat "+agentPath+"| grep class | grep ' < ' | awk -F ' ' '{print $2}'")

	output, err := cmd.CombinedOutput()

	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(output)), nil
}

func isAgentRunning(agentName string) (bool, error) {
	cmd := exec.Command("/bin/pgrep -f ", agentName, " && echo OK || echo Failed")

	output, err := cmd.CombinedOutput()

	if err != nil {
		return false, err
	}

	if string(output) == "OK" {
		return true, nil
	}

	return false, nil
}

func populateInventory(inventory sdk.Inventory) error {
	smithAgent, _ := getAgentListString()
	agentsList, _ := extractAgentList(smithAgent)
	runningAgentsList, _ := extractRunningAgentList(smithAgent)

	inventory.SetItem("smithAgentsGroup", "value", os.Getenv("SMITH_AGENTS_GROUP"))
	inventory.SetItem("smithAgentDirectories", "value", os.Getenv("SMITH_AGENT_DIRECTORIES"))
	inventory.SetItem("smithRunningAgents", "value", runningAgentsList)
	inventory.SetItem("smithAgents", "value", agentsList)

	return nil
}

func extractAgentTotalCount(smithAgent string) (int, error) {
	splitLine, _ := extractAgentList(smithAgent)
	metricValue := len(splitLine)
	return metricValue, nil
}

func extractRunningAgentList(smithAgent string) ([]string, error) {
	agentList, _ := extractAgentList(smithAgent)

	var runningAgents []string

	for _, agentName := range agentList {
		running, _ := isAgentRunning(agentName)
		if running {
			runningAgents = append(runningAgents, agentName)
		}
	}

	return runningAgents, nil
}

func extractAgentList(smithAgent string) ([]string, error) {
	splitLine := strings.Split(string(smithAgent), "\n")
	return splitLine, nil
}

func populateMetrics(ms *metric.MetricSet) error {
	list, err := getAgentListString()

	if err != nil {
		return err
	}

	agentTotal, _ := extractAgentTotalCount(list)
	// agentList, _ := extractAgentList(list)
	agentRunningList, _ := extractRunningAgentList(list)

	if err != nil {
		return err
	}

	ms.SetMetric("agents.total", agentTotal, metric.GAUGE)
	ms.SetMetric("agents.running", len(agentRunningList), metric.GAUGE)

	return nil
}

func main() {
	integration, err := sdk.NewIntegration(integrationName, integrationVersion, &args)
	fatalIfErr(err)

	if args.All || args.Inventory {
		fatalIfErr(populateInventory(integration.Inventory))
	}

	if args.All || args.Metrics {
		ms := integration.NewMetricSet("agents")
		fatalIfErr(populateMetrics(ms))
	}
	fatalIfErr(integration.Publish())
}

func fatalIfErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
