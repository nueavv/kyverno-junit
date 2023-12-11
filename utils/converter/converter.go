package converter

import (
	"os"
	"encoding/xml"
	"gopkg.in/yaml.v3"
	"github.com/nueavv/kyverno-junit/utils/junit"
	kyverno "github.com/kyverno/kyverno/api/policyreport/v1alpha2"
)

func readClusterPolicyReport(data string) (kyverno.ClusterPolicyReport, error) {
	var clusterPolicyReport kyverno.ClusterPolicyReport
	err := yaml.Unmarshal([]byte(data), &clusterPolicyReport)
	if err != nil {
		return clusterPolicyReport, err
	}
	return clusterPolicyReport, nil
}

func readPolicyReport(data string) (kyverno.PolicyReport, error) {
	var policyReport kyverno.PolicyReport
	err := yaml.Unmarshal([]byte(data), &policyReport)
	if err != nil {
		return policyReport, err
	}
	return policyReport, nil
}


func MakeClusterJunitReport(report kyverno.ClusterPolicyReport, output string) error {
	var testsuite junit.TestSuite
	testsuite.Name = "kyverno cluster policy report analyze"

	results := kyverno.GetResults(report)
	for _, result := range results {
		testcase := &junit.TestCase{
			Name: result.Policy,
		}

		switch result.Result() {
		case kyverno.StatusError:
			testcase.Errors = append(testcase.Errors, &junit.Error{
				Message: report.Message,
				Type:    report.Rule,
			})
			// if report.IsFileAnalze() {
			// 	testcase.Errors[0].File = report.GetErrorFile()
			// 	testcase.Errors[0].Line = report.GetErrorLine()
			// }
		case kyverno.StatusFail:
			testcase.Failures = append(testcase.Failures, &junit.Failure{
				Message: report.Message,
				Type:    report.Rule,
			})
		// case kyverno.StatusWarn:
		// 	testcase.Failures = append(testcase.Failures, &junit.Failure{
		// 		Message: report.GetMessage(),
		// 		Type:    report.GetCode(),
		// 	})
		// case kyverno.StatusPass:
		// 	testcase.Errors = append(testcase.Errors, &junit.Error{
		// 		Message: report.GetMessage(),
		// 		Type:    report.GetCode(),
		// 	})
		// 	if report.IsFileAnalze() {
		// 		testcase.Errors[0].File = report.GetErrorFile()
		// 		testcase.Errors[0].Line = report.GetErrorLine()
		// 	}
		// case kyverno.StatusSkip:
		// 	testcase.Failures = append(testcase.Failures, &junit.Failure{
		// 		Message: report.GetMessage(),
		// 		Type:    report.GetCode(),
		// 	})
		}

		testsuite.TestCases = append(testsuite.TestCases, testcase)
	}

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
	var testsuite junit.TestSuite
	// for _, report := range reports {
	// 	testcase := &junit.TestCase{
	// 		Name: report.GetOrigin(),
	// 	}

	// 	switch report.GetLevel() {
	// 	case converter.StatusError:
	// 		testcase.Errors = append(testcase.Errors, &junit.Error{
	// 			Message: report.GetMessage(),
	// 			Type:    report.GetCode(),
	// 		})
	// 		if report.IsFileAnalze() {
	// 			testcase.Errors[0].File = report.GetErrorFile()
	// 			testcase.Errors[0].Line = report.GetErrorLine()
	// 		}
	// 	case converter.StatusFailed:
	// 		testcase.Failures = append(testcase.Failures, &junit.Failure{
	// 			Message: report.GetMessage(),
	// 			Type:    report.GetCode(),
	// 		})
	// 	}

	// 	testsuite.TestCases = append(testsuite.TestCases, testcase)
	// }

	testsuite.Name = "kyverno policy report analyze"
	// testsuite.Errors = converter.GetErrorCount(reports)
	// testsuite.Failures = converter.GetWarningCount(reports)

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