package smsservice_test

import (
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"

	"testing"
)

func TestSmsService(t *testing.T) {
	RegisterFailHandler(Fail)
	RunSpecs(t, "Sms Service Suite")
}
