package datadog

import (
	"fmt"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
	"github.com/zorkian/go-datadog-api"
)

func TestAccDatadogMonitor_Basic(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDatadogMonitorDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckDatadogMonitorConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDatadogMonitorExists("datadog_monitor.foo"),
					resource.TestCheckResourceAttr(
						"datadog_monitor.foo", "name", "name for monitor foo"),
					resource.TestCheckResourceAttr(
						"datadog_monitor.foo", "message", "some message Notify: @hipchat-channel"),
					resource.TestCheckResourceAttr(
						"datadog_monitor.foo", "type", "metric alert"),
					resource.TestCheckResourceAttr(
						"datadog_monitor.foo", "query", "avg(last_1h):avg:aws.ec2.cpu{environment:foo,host:foo} by {host} > 2"),
					resource.TestCheckResourceAttr(
						"datadog_monitor.foo", "notify_no_data", "false"),
					resource.TestCheckResourceAttr(
						"datadog_monitor.foo", "renotify_interval", "60"),
					resource.TestCheckResourceAttr(
						"datadog_monitor.foo", "thresholds.warning", "1.0"),
					resource.TestCheckResourceAttr(
						"datadog_monitor.foo", "thresholds.critical", "2.0"),
					resource.TestCheckResourceAttr(
						"datadog_monitor.foo", "require_full_window", "true"),
					resource.TestCheckResourceAttr(
						"datadog_monitor.foo", "locked", "false"),
					resource.TestCheckResourceAttr(
						"datadog_monitor.foo", "tags.foo", "bar"),
					resource.TestCheckResourceAttr(
						"datadog_monitor.foo", "tags.bar", "baz"),
				),
			},
		},
	})
}

func TestAccDatadogMonitor_Updated(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDatadogMonitorDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckDatadogMonitorConfig,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDatadogMonitorExists("datadog_monitor.foo"),
					resource.TestCheckResourceAttr(
						"datadog_monitor.foo", "name", "name for monitor foo"),
					resource.TestCheckResourceAttr(
						"datadog_monitor.foo", "message", "some message Notify: @hipchat-channel"),
					resource.TestCheckResourceAttr(
						"datadog_monitor.foo", "escalation_message", "the situation has escalated @pagerduty"),
					resource.TestCheckResourceAttr(
						"datadog_monitor.foo", "query", "avg(last_1h):avg:aws.ec2.cpu{environment:foo,host:foo} by {host} > 2"),
					resource.TestCheckResourceAttr(
						"datadog_monitor.foo", "type", "metric alert"),
					resource.TestCheckResourceAttr(
						"datadog_monitor.foo", "notify_no_data", "false"),
					resource.TestCheckResourceAttr(
						"datadog_monitor.foo", "renotify_interval", "60"),
					resource.TestCheckResourceAttr(
						"datadog_monitor.foo", "thresholds.warning", "1.0"),
					resource.TestCheckResourceAttr(
						"datadog_monitor.foo", "thresholds.critical", "2.0"),
					resource.TestCheckResourceAttr(
						"datadog_monitor.foo", "notify_audit", "false"),
					resource.TestCheckResourceAttr(
						"datadog_monitor.foo", "timeout_h", "60"),
					resource.TestCheckResourceAttr(
						"datadog_monitor.foo", "include_tags", "true"),
					resource.TestCheckResourceAttr(
						"datadog_monitor.foo", "require_full_window", "true"),
					resource.TestCheckResourceAttr(
						"datadog_monitor.foo", "locked", "false"),
					resource.TestCheckResourceAttr(
						"datadog_monitor.foo", "tags.foo", "bar"),
					resource.TestCheckResourceAttr(
						"datadog_monitor.foo", "tags.bar", "baz"),
				),
			},
			resource.TestStep{
				Config: testAccCheckDatadogMonitorConfigUpdated,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDatadogMonitorExists("datadog_monitor.foo"),
					resource.TestCheckResourceAttr(
						"datadog_monitor.foo", "name", "name for monitor bar"),
					resource.TestCheckResourceAttr(
						"datadog_monitor.foo", "message", "a different message Notify: @hipchat-channel"),
					resource.TestCheckResourceAttr(
						"datadog_monitor.foo", "query", "avg(last_1h):avg:aws.ec2.cpu{environment:bar,host:bar} by {host} > 3"),
					resource.TestCheckResourceAttr(
						"datadog_monitor.foo", "escalation_message", "the situation has escalated! @pagerduty"),
					resource.TestCheckResourceAttr(
						"datadog_monitor.foo", "type", "metric alert"),
					resource.TestCheckResourceAttr(
						"datadog_monitor.foo", "notify_no_data", "true"),
					resource.TestCheckResourceAttr(
						"datadog_monitor.foo", "no_data_timeframe", "20"),
					resource.TestCheckResourceAttr(
						"datadog_monitor.foo", "renotify_interval", "40"),
					resource.TestCheckResourceAttr(
						"datadog_monitor.foo", "thresholds.ok", "0.0"),
					resource.TestCheckResourceAttr(
						"datadog_monitor.foo", "thresholds.warning", "1.0"),
					resource.TestCheckResourceAttr(
						"datadog_monitor.foo", "thresholds.critical", "3.0"),
					resource.TestCheckResourceAttr(
						"datadog_monitor.foo", "notify_audit", "true"),
					resource.TestCheckResourceAttr(
						"datadog_monitor.foo", "timeout_h", "70"),
					resource.TestCheckResourceAttr(
						"datadog_monitor.foo", "include_tags", "false"),
					resource.TestCheckResourceAttr(
						"datadog_monitor.foo", "silenced.*", "0"),
					resource.TestCheckResourceAttr(
						"datadog_monitor.foo", "require_full_window", "false"),
					resource.TestCheckResourceAttr(
						"datadog_monitor.foo", "locked", "true"),
					resource.TestCheckResourceAttr(
						"datadog_monitor.foo", "tags.baz", "qux"),
					resource.TestCheckResourceAttr(
						"datadog_monitor.foo", "tags.quux", "corge"),
				),
			},
		},
	})
}

func TestAccDatadogMonitor_TrimWhitespace(t *testing.T) {
	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckDatadogMonitorDestroy,
		Steps: []resource.TestStep{
			resource.TestStep{
				Config: testAccCheckDatadogMonitorConfigWhitespace,
				Check: resource.ComposeTestCheckFunc(
					testAccCheckDatadogMonitorExists("datadog_monitor.foo"),
					resource.TestCheckResourceAttr(
						"datadog_monitor.foo", "name", "name for monitor foo"),
					resource.TestCheckResourceAttr(
						"datadog_monitor.foo", "message", "some message Notify: @hipchat-channel"),
					resource.TestCheckResourceAttr(
						"datadog_monitor.foo", "type", "metric alert"),
					resource.TestCheckResourceAttr(
						"datadog_monitor.foo", "query", "avg(last_1h):avg:aws.ec2.cpu{environment:foo,host:foo} by {host} > 2"),
					resource.TestCheckResourceAttr(
						"datadog_monitor.foo", "notify_no_data", "false"),
					resource.TestCheckResourceAttr(
						"datadog_monitor.foo", "renotify_interval", "60"),
					resource.TestCheckResourceAttr(
						"datadog_monitor.foo", "thresholds.ok", "0.0"),
					resource.TestCheckResourceAttr(
						"datadog_monitor.foo", "thresholds.warning", "1.0"),
					resource.TestCheckResourceAttr(
						"datadog_monitor.foo", "thresholds.critical", "2.0"),
				),
			},
		},
	})
}

func testAccCheckDatadogMonitorDestroy(s *terraform.State) error {
	client := testAccProvider.Meta().(*datadog.Client)

	if err := destroyHelper(s, client); err != nil {
		return err
	}
	return nil
}

func testAccCheckDatadogMonitorExists(n string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*datadog.Client)
		if err := existsHelper(s, client); err != nil {
			return err
		}
		return nil
	}
}

const testAccCheckDatadogMonitorConfig = `
resource "datadog_monitor" "foo" {
  name = "name for monitor foo"
  type = "metric alert"
  message = "some message Notify: @hipchat-channel"
  escalation_message = "the situation has escalated @pagerduty"

  query = "avg(last_1h):avg:aws.ec2.cpu{environment:foo,host:foo} by {host} > 2"

  thresholds {
	warning = "1.0"
	critical = "2.0"
  }

  notify_no_data = false
  renotify_interval = 60

  notify_audit = false
  timeout_h = 60
  include_tags = true
  require_full_window = true
  locked = false
  tags {
	"foo" = "bar"
	"bar" = "baz"
  }
}
`

const testAccCheckDatadogMonitorConfigUpdated = `
resource "datadog_monitor" "foo" {
  name = "name for monitor bar"
  type = "metric alert"
  message = "a different message Notify: @hipchat-channel"
  escalation_message = "the situation has escalated @pagerduty"

  query = "avg(last_1h):avg:aws.ec2.cpu{environment:bar,host:bar} by {host} > 3"

  thresholds {
	ok = "0.0"
	warning = "1.0"
	critical = "3.0"
  }

  notify_no_data = true
  no_data_timeframe = 20
  renotify_interval = 40
  escalation_message = "the situation has escalated! @pagerduty"
  notify_audit = true
  timeout_h = 70
  include_tags = false
  require_full_window = false
  locked = true
  silenced {
	"*" = 0
  }
  tags {
	"baz"  = "qux"
	"quux" = "corge"
  }
}
`

const testAccCheckDatadogMonitorConfigWhitespace = `
resource "datadog_monitor" "foo" {
  name = "name for monitor foo"
  type = "metric alert"
  message = <<EOF
some message Notify: @hipchat-channel
EOF
  escalation_message = <<EOF
the situation has escalated @pagerduty
EOF
  query = <<EOF
avg(last_1h):avg:aws.ec2.cpu{environment:foo,host:foo} by {host} > 2
EOF
  thresholds {
	ok = "0.0"
	warning = "1.0"
	critical = "2.0"
  }

  notify_no_data = false
  renotify_interval = 60

  notify_audit = false
  timeout_h = 60
  include_tags = true
}
`

func destroyHelper(s *terraform.State, client *datadog.Client) error {
	for _, r := range s.RootModule().Resources {
		i, _ := strconv.Atoi(r.Primary.ID)
		if _, err := client.GetMonitor(i); err != nil {
			if strings.Contains(err.Error(), "404 Not Found") {
				continue
			}
			return fmt.Errorf("Received an error retrieving monitor %s", err)
		}
		return fmt.Errorf("Monitor still exists")
	}
	return nil
}

func existsHelper(s *terraform.State, client *datadog.Client) error {
	for _, r := range s.RootModule().Resources {
		i, _ := strconv.Atoi(r.Primary.ID)
		if _, err := client.GetMonitor(i); err != nil {
			return fmt.Errorf("Received an error retrieving monitor %s", err)
		}
	}
	return nil
}
