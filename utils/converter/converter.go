package converter

import (
	"encoding/xml"
	"fmt"
	kyverno "github.com/kyverno/kyverno/api/policyreport/v1alpha2"
	"github.com/nueavv/kyverno-junit/utils/junit"
	"gopkg.in/yaml.v3"
	"os"
)

func ReadClusterPolicyReport(data []byte) (kyverno.ClusterPolicyReport, error) {
	var clusterPolicyReport kyverno.ClusterPolicyReport
	err := yaml.Unmarshal(data, &clusterPolicyReport)
	if err != nil {
		return clusterPolicyReport, err
	}
	return clusterPolicyReport, nil
}

func ReadPolicyReport(data []byte) (kyverno.PolicyReport, error) {
	var policyReport kyverno.PolicyReport
	err := yaml.Unmarshal(data, &policyReport)
	if err != nil {
		return policyReport, err
	}
	return policyReport, nil
}

func MakeClusterJunitReport(report kyverno.ClusterPolicyReport, output string) error {
	testsuite := makeReport(report.GetResults())
	testsuite.Name = "kyverno cluster policy report analyze"
	xmlBytes, err := xml.Marshal(testsuite)
	if err != nil {
		return err
	}

	error := os.WriteFile(output, xmlBytes, 0660)
	if error != nil {
		return error
	}

	return nil
}

func MakeJunitReport(report kyverno.PolicyReport, output string) error {
	testsuite := makeReport(report.GetResults())
	testsuite.Name = "kyverno policy report analyze"
	xmlBytes, err := xml.Marshal(testsuite)
	if err != nil {
		return err
	}

	error := os.WriteFile(output, xmlBytes, 0660)
	if error != nil {
		return error
	}

	return nil
}

func makeReport(results []kyverno.PolicyReportResult) junit.TestSuite {
	var testsuite junit.TestSuite
	for _, result := range results {
		testcase := &junit.TestCase{
			Name: result.Policy,
		}

		switch result.Result {
		case kyverno.StatusError:
			testcase.Errors = append(testcase.Errors, &junit.Error{
				Message: result.Message,
				Type:    result.Rule,
			})
		case kyverno.StatusFail, kyverno.StatusWarn:
			testcase.Failures = append(testcase.Failures, &junit.Failure{
				Message: result.Message,
				Type:    result.Rule,
			})
		// case kyverno.StatusPass:
		// 	testcase.Status =
		case kyverno.StatusSkip:
			testcase.Skipped = fmt.Sprintf("Policy: %s, Rule: %s", result.Policy, result.Rule)
		}
		testsuite.TestCases = append(testsuite.TestCases, testcase)
	}
	return testsuite
}
