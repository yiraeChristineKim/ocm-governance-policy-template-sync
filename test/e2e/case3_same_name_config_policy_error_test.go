// Copyright (c) 2021 Red Hat, Inc.
// Copyright Contributors to the Open Cluster Management project

package e2e

import (
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	"open-cluster-management.io/governance-policy-propagator/test/utils"
)

const (
	case3PolicyName       string = "same-policy"
	case3PolicyYaml       string = "../resources/case3_same_name_config_policy_error/case3-policy.yaml"
	case3ConfigPolicyName string = "policy-config-same"
)

var _ = Describe("Test same name of configurationpolicies error", func() {
	AfterEach(func() {
		_, err := utils.KubectlWithOutput("delete", "policies", "--all", "-n", testNamespace)
		Expect(err).ShouldNot(HaveOccurred())
		_, err = utils.KubectlWithOutput("delete", "events", "--field-selector=involvedObject.name="+case3PolicyName,
			"--ignore-not-found", "-n", testNamespace)
		Expect(err).ShouldNot(HaveOccurred())
	})
	It("should not reconcile the policy when there are same names in policy-template", func() {
		By("Creating policy")
		_, err := utils.KubectlWithOutput("apply", "-f", case3PolicyYaml, "-n", testNamespace)
		Expect(err).ShouldNot(HaveOccurred())

		By("Should generate warning events")
		Eventually(
			checkForEvent(case3PolicyName,
				"There are duplicate name in configurationpolicies, please check the policy"),
			defaultTimeoutSeconds,
			1,
		).Should(BeTrue())

		By("Should not create any configuration policy")
		Consistently(func() interface{} {
			return utils.GetWithTimeout(clientManagedDynamic, gvrConfigurationPolicy, case3ConfigPolicyName,
				testNamespace, false, defaultTimeoutSeconds)
		}, 10, 1).Should(BeNil())
	})
})
