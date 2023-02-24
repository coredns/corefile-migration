package migration

import (
	"testing"
)

func TestMigrate(t *testing.T) {
	testCases := []struct {
		name             string
		fromVersion      string
		toVersion        string
		deprecations     bool
		startCorefile    string
		expectedCorefile string
	}{
		{

			name:         "Add max_concurrent option to forward plugin",
			fromVersion:  "1.6.2",
			toVersion:    "1.7.0",
			deprecations: true,
			startCorefile: `.:53 {
    errors
    health {
        lameduck 5s
    }
    ready
    kubernetes cluster.local in-addr.arpa ip6.arpa {
        pods insecure
        fallthrough in-addr.arpa ip6.arpa
    }
    prometheus :9153
    forward . /etc/resolv.conf
    cache 30
    loop
    reload
    loadbalance
}
`,
			expectedCorefile: `.:53 {
    errors
    health {
        lameduck 5s
    }
    ready
    kubernetes cluster.local in-addr.arpa ip6.arpa {
        pods insecure
        fallthrough in-addr.arpa ip6.arpa
    }
    prometheus :9153
    forward . /etc/resolv.conf {
        max_concurrent 1000
    }
    cache 30
    loop
    reload
    loadbalance
}
`,
		},
		{
			name:         "replace/remove proxy options",
			fromVersion:  "1.3.1",
			toVersion:    "1.5.0",
			deprecations: true,
			startCorefile: `.:53 {
    errors
    health
    loop
    kubernetes cluster.local in-addr.arpa ip6.arpa {
        pods insecure
        fallthrough in-addr.arpa ip6.arpa
    }
    prometheus :9153
    proxy . /etc/resolv.conf {
       except
       fail_timeout
       max_fails
       health_check
       spray
       policy least_conn
    }
    cache 30
    reload
    loadbalance
}
`,
			expectedCorefile: `.:53 {
    errors
    health
    loop
    kubernetes cluster.local in-addr.arpa ip6.arpa {
        pods insecure
        fallthrough in-addr.arpa ip6.arpa
    }
    prometheus :9153
    forward . /etc/resolv.conf {
        except
        force_tcp
    }
    cache 30
    reload
    loadbalance
    ready
}
`,
		},
		{
			name:         "replace proxy options",
			fromVersion:  "1.3.1",
			toVersion:    "1.5.0",
			deprecations: true,
			startCorefile: `.:53 {
    errors
    health
    loop
    kubernetes cluster.local in-addr.arpa ip6.arpa {
        pods insecure
        fallthrough in-addr.arpa ip6.arpa
    }
    prometheus :9153
    proxy . /etc/resolv.conf {
       protocol force_tcp insecure
    }
    cache 30
    reload
    loadbalance
}
`,
			expectedCorefile: `.:53 {
    errors
    health
    loop
    kubernetes cluster.local in-addr.arpa ip6.arpa {
        pods insecure
        fallthrough in-addr.arpa ip6.arpa
    }
    prometheus :9153
    forward . /etc/resolv.conf {
        force_tcp
    }
    cache 30
    reload
    loadbalance
    ready
}
`,
		},
		{

			name:         "Add lameduck option to health plugin",
			fromVersion:  "1.6.2",
			toVersion:    "1.6.6",
			deprecations: true,
			startCorefile: `.:53 {
    errors
    health
    ready
    kubernetes cluster.local in-addr.arpa ip6.arpa {
        pods insecure
        fallthrough in-addr.arpa ip6.arpa
    }
    prometheus :9153
    forward . /etc/resolv.conf
    cache 30
    loop
    reload
    loadbalance
}
`,
			expectedCorefile: `.:53 {
    errors
    health {
        lameduck 5s
    }
    ready
    kubernetes cluster.local in-addr.arpa ip6.arpa {
        pods insecure
        fallthrough in-addr.arpa ip6.arpa
    }
    prometheus :9153
    forward . /etc/resolv.conf
    cache 30
    loop
    reload
    loadbalance
}
`,
		},
		{
			name:         "Handle Kubernetes plugin deprecations",
			fromVersion:  "1.4.0",
			toVersion:    "1.7.1",
			deprecations: true,
			startCorefile: `.:53 {
    errors
    health
    kubernetes cluster.local in-addr.arpa ip6.arpa {
        resyncperiod 60s
        pods insecure
        upstream
        fallthrough in-addr.arpa ip6.arpa
    }
    prometheus :9153
    proxy . /etc/resolv.conf
    cache 30
    loop
    reload
    loadbalance
}
`,
			expectedCorefile: `.:53 {
    errors
    health {
        lameduck 5s
    }
    kubernetes cluster.local in-addr.arpa ip6.arpa {
        pods insecure
        fallthrough in-addr.arpa ip6.arpa
    }
    prometheus :9153
    forward . /etc/resolv.conf {
        max_concurrent 1000
    }
    cache 30
    loop
    reload
    loadbalance
    ready
}
`,
		},
		{
			name:         "Remove invalid proxy option",
			fromVersion:  "1.1.3",
			toVersion:    "1.2.6",
			deprecations: true,
			startCorefile: `.:53 {
    errors
    health
    kubernetes cluster.local in-addr.arpa ip6.arpa {
        endpoint thing1 thing2
        pods insecure
        upstream
        fallthrough in-addr.arpa ip6.arpa
    }
    prometheus :9153
    proxy example.org 1.2.3.4:53 {
        protocol https_google
    }
    cache 30
    loop
    reload
    loadbalance
}
`,
			expectedCorefile: `.:53 {
    errors
    health
    kubernetes cluster.local in-addr.arpa ip6.arpa {
        endpoint thing1 thing2
        pods insecure
        upstream
        fallthrough in-addr.arpa ip6.arpa
    }
    prometheus :9153
    proxy example.org 1.2.3.4:53
    cache 30
    loop
    reload
    loadbalance
}
`,
		},
		{
			name:         "Migrate from proxy to forward and handle Kubernetes deprecations",
			fromVersion:  "1.3.1",
			toVersion:    "1.5.2",
			deprecations: true,
			startCorefile: `.:53 {
    errors
    health
    kubernetes cluster.local in-addr.arpa ip6.arpa {
        endpoint thing1 thing2
        pods insecure
        upstream
        fallthrough in-addr.arpa ip6.arpa
    }
    prometheus :9153
    proxy . /etc/resolv.conf
    cache 30
    loop
    reload
    loadbalance
}
`,
			expectedCorefile: `.:53 {
    errors
    health
    kubernetes cluster.local in-addr.arpa ip6.arpa {
        endpoint thing1
        pods insecure
        fallthrough in-addr.arpa ip6.arpa
    }
    prometheus :9153
    forward . /etc/resolv.conf
    cache 30
    loop
    reload
    loadbalance
    ready
}
`,
		},
		{
			name:         "add missing loop and ready plugins",
			fromVersion:  "1.1.3",
			toVersion:    "1.5.0",
			deprecations: true,
			startCorefile: `.:53 {
    errors
    health
    kubernetes cluster.local in-addr.arpa ip6.arpa {
        pods insecure
        upstream
        fallthrough in-addr.arpa ip6.arpa
    }
    prometheus :9153
    proxy . /etc/resolv.conf
    cache 30
    reload
    loadbalance
}
`,
			expectedCorefile: `.:53 {
    errors
    health
    kubernetes cluster.local in-addr.arpa ip6.arpa {
        pods insecure
        fallthrough in-addr.arpa ip6.arpa
    }
    prometheus :9153
    forward . /etc/resolv.conf
    cache 30
    reload
    loadbalance
    loop
    ready
}
`,
		},
		{
			name:         "handle multiple proxy plugins",
			fromVersion:  "1.1.3",
			toVersion:    "1.5.0",
			deprecations: true,
			startCorefile: `.:53 {
    errors
    health
    kubernetes cluster.local in-addr.arpa ip6.arpa {
        pods insecure
        upstream
        fallthrough in-addr.arpa ip6.arpa
    }
    prometheus :9153
    proxy mystub-1.example.org 1.2.3.4
    proxy mystub-2.example.org 5.6.7.8
    proxy . /etc/resolv.conf
    cache 30
    reload
    loadbalance
}
`,
			expectedCorefile: `.:53 {
    errors
    health
    kubernetes cluster.local in-addr.arpa ip6.arpa {
        pods insecure
        fallthrough in-addr.arpa ip6.arpa
    }
    prometheus :9153
    forward . /etc/resolv.conf
    cache 30
    reload
    loadbalance
    loop
    ready
}

mystub-1.example.org {
    forward . 1.2.3.4
    loop
    errors
    cache 30
}

mystub-2.example.org {
    forward . 5.6.7.8
    loop
    errors
    cache 30
}
`,
		},
		{
			name:         "no-op same version migration",
			fromVersion:  "1.3.1",
			toVersion:    "1.3.1",
			deprecations: true,
			startCorefile: `.:53 {
    errors
    health
    kubernetes cluster.local in-addr.arpa ip6.arpa {
        pods insecure
        upstream
        fallthrough in-addr.arpa ip6.arpa
    }
    prometheus :9153
    proxy . /etc/resolv.conf
    cache 30
    loop
    reload
    loadbalance
}
`,
			expectedCorefile: `.:53 {
    errors
    health
    kubernetes cluster.local in-addr.arpa ip6.arpa {
        pods insecure
        upstream
        fallthrough in-addr.arpa ip6.arpa
    }
    prometheus :9153
    proxy . /etc/resolv.conf
    cache 30
    loop
    reload
    loadbalance
}
`,
		},
		{
			name:         "handle supported non-default plugins (hosts)",
			fromVersion:  "1.3.1",
			toVersion:    "1.6.6",
			deprecations: true,
			startCorefile: `.:53 {
    errors
    health
    ready
    autopath {
        @kubernetes
    }
    kubernetes cluster.local in-addr.arpa ip6.arpa {
        pods insecure
        fallthrough in-addr.arpa ip6.arpa
    }
    prometheus :9153
    forward . /etc/resolv.conf
    cache 30
    loop
    reload
    loadbalance
    hosts {
        1.2.3.4 hello
        1:2:3::4 goodbye
        ttl 30
    }
}
`,
			expectedCorefile: `.:53 {
    errors
    health {
        lameduck 5s
    }
    ready
    autopath {
        @kubernetes
    }
    kubernetes cluster.local in-addr.arpa ip6.arpa {
        pods insecure
        fallthrough in-addr.arpa ip6.arpa
    }
    prometheus :9153
    forward . /etc/resolv.conf
    cache 30
    loop
    reload
    loadbalance
    hosts {
        1.2.3.4 hello
        1:2:3::4 goodbye
        ttl 30
    }
}
`,
		},
		{
			name:         "handle supported non-default plugins (rewrite)",
			fromVersion:  "1.2.6",
			toVersion:    "1.5.0",
			deprecations: true,
			startCorefile: `.:53 {
    errors
    health
    loop
    kubernetes cluster.local in-addr.arpa ip6.arpa {
        pods insecure
        fallthrough in-addr.arpa ip6.arpa
    }
    prometheus :9153
    proxy . /etc/resolv.conf
    cache 30
    reload
    loadbalance
    rewrite continue {
        ttl regex (.*)\.coredns\.rocks 15
    }
}
`,
			expectedCorefile: `.:53 {
    errors
    health
    loop
    kubernetes cluster.local in-addr.arpa ip6.arpa {
        pods insecure
        fallthrough in-addr.arpa ip6.arpa
    }
    prometheus :9153
    forward . /etc/resolv.conf
    cache 30
    reload
    loadbalance
    rewrite continue {
        ttl regex (.*)\.coredns\.rocks 15
    }
    ready
}
`,
		},
		{
			name:         "k8s transfer moved to plugin",
			fromVersion:  "1.7.1",
			toVersion:    "1.8.0",
			deprecations: true,
			startCorefile: `.:53 {
    kubernetes cluster.local in-addr.arpa ip6.arpa {
        transfer to 1.2.3.4 5.6.7.8
    }
}
`,
			expectedCorefile: `.:53 {
    kubernetes cluster.local in-addr.arpa ip6.arpa
    transfer cluster.local {
        to 1.2.3.4 5.6.7.8
    }
}
`,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {

			result, err := Migrate(testCase.fromVersion, testCase.toVersion, testCase.startCorefile, testCase.deprecations)

			if err != nil {
				t.Errorf("%v", err)
			}

			if result != testCase.expectedCorefile {
				t.Errorf("expected != result\n%v\n%v", testCase.expectedCorefile, result)
			}
		})
	}
}

func TestMigrateDown(t *testing.T) {
	testCases := []struct {
		name             string
		fromVersion      string
		toVersion        string
		deprecations     bool
		startCorefile    string
		expectedCorefile string
	}{
		{
			name:        "from 1.6.2 to 1.1.3",
			fromVersion: "1.6.2",
			toVersion:   "1.1.3",
			startCorefile: `.:53 {
    errors
    health
    ready
    kubernetes cluster.local in-addr.arpa ip6.arpa {
        pods insecure
        upstream
        fallthrough in-addr.arpa ip6.arpa
    }
    prometheus :9153
    forward . /etc/resolv.conf
    cache 30
    loop
    reload
    loadbalance
}
`,
			expectedCorefile: `.:53 {
    errors
    health
    kubernetes cluster.local in-addr.arpa ip6.arpa {
        pods insecure
        upstream
        fallthrough in-addr.arpa ip6.arpa
    }
    prometheus :9153
    forward . /etc/resolv.conf
    cache 30
    reload
    loadbalance
}
`,
		},
		{
			name:        "from 1.5.0 to 1.3.1",
			fromVersion: "1.5.0",
			toVersion:   "1.3.1",
			startCorefile: `.:53 {
    errors
    health
    ready
    kubernetes cluster.local in-addr.arpa ip6.arpa {
        pods insecure
        upstream
        fallthrough in-addr.arpa ip6.arpa
    }
    prometheus :9153
    forward . /etc/resolv.conf
    cache 30
    loop
    reload
    loadbalance
}
`,
			expectedCorefile: `.:53 {
    errors
    health
    kubernetes cluster.local in-addr.arpa ip6.arpa {
        pods insecure
        upstream
        fallthrough in-addr.arpa ip6.arpa
    }
    prometheus :9153
    forward . /etc/resolv.conf
    cache 30
    loop
    reload
    loadbalance
}
`,
		},
		{
			name:         "no-op same version migration",
			fromVersion:  "1.3.1",
			toVersion:    "1.3.1",
			deprecations: true,
			startCorefile: `.:53 {
    errors
    health
    kubernetes cluster.local in-addr.arpa ip6.arpa {
        pods insecure
        upstream
        fallthrough in-addr.arpa ip6.arpa
    }
    prometheus :9153
    proxy . /etc/resolv.conf
    cache 30
    loop
    reload
    loadbalance
}
`,
			expectedCorefile: `.:53 {
    errors
    health
    kubernetes cluster.local in-addr.arpa ip6.arpa {
        pods insecure
        upstream
        fallthrough in-addr.arpa ip6.arpa
    }
    prometheus :9153
    proxy . /etc/resolv.conf
    cache 30
    loop
    reload
    loadbalance
}
`,
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {

			result, err := MigrateDown(testCase.fromVersion, testCase.toVersion, testCase.startCorefile)

			if err != nil {
				t.Errorf("%v", err)
			}

			if result != testCase.expectedCorefile {
				t.Errorf("expected != result:\n%v\n%v", testCase.expectedCorefile, result)
			}
		})
	}
}

func TestDeprecated(t *testing.T) {
	startCorefile := `.:53 {
    errors
    health
    kubernetes cluster.local in-addr.arpa ip6.arpa {
        pods insecure
        upstream
        fallthrough in-addr.arpa ip6.arpa
        resyncperiod
    }
    prometheus :9153
    proxy . /etc/resolv.conf
    cache 30
    reload
    loop
    loadbalance
}
`

	expected := []Notice{
		{Plugin: "kubernetes", Option: "upstream", Severity: SevDeprecated, Version: "1.4.0"},
		{Plugin: "proxy", Severity: SevDeprecated, ReplacedBy: "forward", Version: "1.4.0"},
		{Option: "upstream", Plugin: "kubernetes", Severity: SevIgnored, Version: "1.5.0"},
		{Plugin: "kubernetes", Option: "resyncperiod", Severity: SevDeprecated, Version: "1.5.0"},
		{Plugin: "proxy", Severity: SevRemoved, ReplacedBy: "forward", Version: "1.5.0"},
		{Plugin: "ready", Severity: SevNewDefault, Version: "1.5.0"},
		{Option: "upstream", Plugin: "kubernetes", Severity: SevIgnored, Version: "1.5.1"},
		{Plugin: "kubernetes", Option: "resyncperiod", Severity: SevDeprecated, Version: "1.5.1"},
		{Option: "upstream", Plugin: "kubernetes", Severity: SevIgnored, Version: "1.5.2"},
		{Plugin: "kubernetes", Option: "resyncperiod", Severity: SevDeprecated, Version: "1.5.2"},
		{Option: "upstream", Plugin: "kubernetes", Severity: SevIgnored, Version: "1.6.0"},
		{Plugin: "kubernetes", Option: "resyncperiod", Severity: SevIgnored, Version: "1.6.0"},
		{Option: "upstream", Plugin: "kubernetes", Severity: SevIgnored, Version: "1.6.1"},
		{Plugin: "kubernetes", Option: "resyncperiod", Severity: SevIgnored, Version: "1.6.1"},
		{Option: "upstream", Plugin: "kubernetes", Severity: SevIgnored, Version: "1.6.2"},
		{Plugin: "kubernetes", Option: "resyncperiod", Severity: SevIgnored, Version: "1.6.2"},
		{Option: "upstream", Plugin: "kubernetes", Severity: SevIgnored, Version: "1.6.3"},
		{Plugin: "kubernetes", Option: "resyncperiod", Severity: SevIgnored, Version: "1.6.3"},
		{Option: "upstream", Plugin: "kubernetes", Severity: SevIgnored, Version: "1.6.4"},
		{Plugin: "kubernetes", Option: "resyncperiod", Severity: SevIgnored, Version: "1.6.4"},
		{Plugin: "health", Option: "lameduck", Severity: SevNewDefault, Version: "1.6.5"},
		{Option: "upstream", Plugin: "kubernetes", Severity: SevIgnored, Version: "1.6.5"},
		{Plugin: "kubernetes", Option: "resyncperiod", Severity: SevIgnored, Version: "1.6.5"},
		{Option: "upstream", Plugin: "kubernetes", Severity: SevIgnored, Version: "1.6.6"},
		{Plugin: "kubernetes", Option: "resyncperiod", Severity: SevIgnored, Version: "1.6.6"},
		{Option: "upstream", Plugin: "kubernetes", Severity: SevIgnored, Version: "1.6.7"},
		{Plugin: "kubernetes", Option: "resyncperiod", Severity: SevIgnored, Version: "1.6.7"},
		{Option: "upstream", Plugin: "kubernetes", Severity: SevIgnored, Version: "1.6.9"},
		{Plugin: "kubernetes", Option: "resyncperiod", Severity: SevIgnored, Version: "1.6.9"},
		{Option: "upstream", Plugin: "kubernetes", Severity: SevRemoved, Version: "1.7.0"},
		{Plugin: "kubernetes", Option: "resyncperiod", Severity: SevRemoved, Version: "1.7.0"},
	}

	result, err := Deprecated("1.1.3", "1.7.0", startCorefile)

	if err != nil {
		t.Fatal(err)
	}

	if len(result) != len(expected) {
		t.Fatalf("expected to find %v notifications; got %v", len(expected), len(result))
	}

	for i, dep := range expected {
		if result[i].ToString() != dep.ToString() {
			t.Errorf("expected to get '%v'; got '%v'", dep.ToString(), result[i].ToString())
		}
	}

	result, err = Deprecated("1.3.1", "1.3.1", startCorefile)
	if err != nil {
		t.Fatal(err)
	}
	expected = []Notice{}
	if len(result) != len(expected) {
		t.Fatalf("expected to find %v notifications in no-op upgrade; got %v", len(expected), len(result))
	}
}

func TestUnsupported(t *testing.T) {
	testCases := []struct {
		name          string
		fromVersion   string
		toVersion     string
		startCorefile string
		expected      []Notice
	}{
		{
			name: "Unsupported route53",
			startCorefile: `.:53 {
    errors {
        consolidate
    }
    health {
        lameduck
    }
    kubernetes cluster.local in-addr.arpa ip6.arpa {
        pods insecure
        upstream
        fallthrough in-addr.arpa ip6.arpa
    }
    route53 example.org.:Z1Z2Z3Z4DZ5Z6Z7
    prometheus :9153
    proxy . /etc/resolv.conf
    cache 30
    loop
    reload
    loadbalance
}
`,
			fromVersion: "1.3.1",
			toVersion:   "1.5.0",
			expected: []Notice{
				{Plugin: "route53", Severity: SevUnsupported, Version: "1.4.0"},
				{Plugin: "route53", Severity: SevUnsupported, Version: "1.5.0"},
			},
		},
		{
			name: "Unsupported route53 - same coredns version",
			startCorefile: `.:53 {
    errors {
        consolidate
    }
    health {
        lameduck
    }
    kubernetes cluster.local in-addr.arpa ip6.arpa {
        pods insecure
        upstream
        fallthrough in-addr.arpa ip6.arpa
    }
    route53 example.org.:Z1Z2Z3Z4DZ5Z6Z7
    prometheus :9153
    proxy . /etc/resolv.conf
    cache 30
    loop
    reload
    loadbalance
}
`,
			fromVersion: "1.3.1",
			toVersion:   "1.3.1",
			expected:    []Notice{},
		},
		{
			name: "Wrong plugin option",
			startCorefile: `.:53 {
    errors
    health
    kubernetes cluster.local in-addr.arpa ip6.arpa {
        moo insecure
        upstream
        fallthrough in-addr.arpa ip6.arpa
    }
    prometheus :9153
    proxy . /etc/resolv.conf
    cache 30
    loop
    reload
    loadbalance
}
`,
			fromVersion: "1.3.1",
			toVersion:   "1.5.0",
			expected: []Notice{
				{Option: "moo", Plugin: "kubernetes", Severity: SevUnsupported, Version: "1.4.0"},
				{Option: "moo", Plugin: "kubernetes", Severity: SevUnsupported, Version: "1.5.0"},
			},
		},
		{
			name: "Invalid Plugin",
			startCorefile: `.:53 {
    errors
    health
    kubernetes cluster.local in-addr.arpa ip6.arpa {
        pods insecure
        upstream
        fallthrough in-addr.arpa ip6.arpa
    }
    invalid
    prometheus :9153
    forward . /etc/resolv.conf
    cache 30
    loop
    reload
    loadbalance
}
`,
			fromVersion: "1.6.6",
			toVersion:   "1.6.7",
			expected: []Notice{
				{Plugin: "invalid", Severity: SevUnsupported, Version: "1.6.7"},
			},
		},
	}

	for _, testCase := range testCases {
		t.Run(testCase.name, func(t *testing.T) {
			result, err := Unsupported(testCase.fromVersion, testCase.toVersion, testCase.startCorefile)

			if err != nil {
				t.Fatal(err)
			}

			if len(result) != len(testCase.expected) {
				t.Fatalf("expected to find %v deprecations; got %v", len(testCase.expected), len(result))
			}

			for i, dep := range testCase.expected {
				if result[i].ToString() != dep.ToString() {
					t.Errorf("expected to get '%v'; got '%v'", dep.ToString(), result[i].ToString())
				}
			}
		})
	}

}

func TestDefault(t *testing.T) {
	defaultCorefiles := []string{`.:53 {
    errors
    health
    kubernetes cluster.local in-addr.arpa ip6.arpa {
        pods insecure
        upstream
        fallthrough in-addr.arpa ip6.arpa
        ttl 30
    }
    prometheus :9153
    forward . /etc/resolv.conf
    cache 30
    loop
    reload
    loadbalance
}
`,
		`.:53 {
    errors
    health
    kubernetes myzone.org in-addr.arpa ip6.arpa {
        pods insecure
        upstream
        fallthrough in-addr.arpa ip6.arpa
        ttl 30
    }
    prometheus :9153
    forward . /etc/resolv.conf
    cache 30
    loop
    reload
    loadbalance
}
`,
		`.:53 {
    errors
    health {
        lameduck 5s
    }
    ready
    kubernetes myzone.org in-addr.arpa ip6.arpa {
        pods insecure
        fallthrough in-addr.arpa ip6.arpa
        ttl 30
    }
    prometheus :9153
    forward . /etc/resolv.conf
    cache 30
    loop
    reload
    loadbalance
}
`,
		`.:53 {
    errors
    health
    ready
    kubernetes myzone.org in-addr.arpa ip6.arpa {
        pods insecure
        fallthrough in-addr.arpa ip6.arpa
        ttl 30
    }
    prometheus :9153
    forward . /etc/resolv.conf
    cache 30
    loop
    reload
    loadbalance
}
`}

	nonDefaultCorefiles := []string{`.:53 {
    errors
    health
    rewrite name suffix myzone.org cluster.local
    kubernetes cluster.local in-addr.arpa ip6.arpa {
        pods insecure
        upstream
        fallthrough in-addr.arpa ip6.arpa
    }
    prometheus :9153
    forward . /etc/resolv.conf
    cache 30
    loop
    reload
    loadbalance
}
`,
		`.:53 {
    errors
    health
    kubernetes cluster.local in-addr.arpa ip6.arpa {
        pods insecure
        upstream
        fallthrough in-addr.arpa ip6.arpa
    }
    prometheus :9153
    forward . /etc/resolv.conf
    cache 30
    loop
    reload
    loadbalance
}
stubzone.org:53 {
    forward . 1.2.3.4
}
`}

	for _, d := range defaultCorefiles {
		if !Default("", d) {
			t.Errorf("expected config to be identified as a default: %v", d)
		}
	}
	for _, d := range nonDefaultCorefiles {
		if Default("", d) {
			t.Errorf("expected config to NOT be identified as a default: %v", d)
		}
	}
}

func TestValidUpMigration(t *testing.T) {
	testCases := []struct {
		from      string
		to        string
		shouldErr bool
	}{
		{"1.3.1", "1.3.1", false},
		{"1.3.1", "1.5.0", false},
		{"1.5.0", "1.3.1", true},
		{"banana", "1.5.0", true},
		{"1.3.1", "apple", true},
		{"banana", "apple", true},
	}

	for _, tc := range testCases {
		err := ValidUpMigration(tc.from, tc.to)

		if !tc.shouldErr && err != nil {
			t.Errorf("expected '%v' to '%v' to be valid versions.", tc.from, tc.to)
		}
		if tc.shouldErr && err == nil {
			t.Errorf("expected '%v' to '%v' to be invalid versions.", tc.from, tc.to)
		}
	}
}

func TestValidDownMigration(t *testing.T) {
	testCases := []struct {
		from      string
		to        string
		shouldErr bool
	}{
		{"1.3.1", "1.3.1", true},
		{"1.3.1", "1.5.0", true},
		{"1.5.0", "1.3.1", false},
		{"banana", "1.5.0", true},
		{"1.3.1", "apple", true},
		{"banana", "apple", true},
	}

	for _, tc := range testCases {
		err := validDownMigration(tc.from, tc.to)

		if !tc.shouldErr && err != nil {
			t.Errorf("expected '%v' to '%v' to be valid versions.", tc.from, tc.to)
		}
		if tc.shouldErr && err == nil {
			t.Errorf("expected '%v' to '%v' to be invalid versions.", tc.from, tc.to)
		}
	}
}

func TestMatchOption(t *testing.T) {
	o := option{}
	p := plugin{
		namedOptions:   map[string]option{"named-option": o},
		patternOptions: map[string]option{"pattern-option-[A-Z]+[0-9]+": o},
	}

	tests := []struct {
		name    string
		matched bool
	}{
		{"named-option", true},
		{"qwerty", false},
		{"pattern-option-A10", true},
		{"pattern-option-a10", false},
	}
	for _, test := range tests {
		gotopt, matched := matchOption(test.name, p)
		if matched != test.matched {
			t.Fatalf("expected %v to match plugin option", test.name)
		}
		if matched == test.matched && !test.matched {
			continue
		}
		if gotopt == nil {
			t.Fatal("expected non-nil returned option")
		}
		if gotopt.name != test.name {
			t.Fatalf("expected returned option name == '%v' got '%v'", test.name, gotopt.name)
		}
	}

}

func TestVersionFromSHA(t *testing.T) {
	testCases := []struct {
		sha       string
		version   string
		shouldErr bool
	}{
		{"2c8d61c46f484d881db43b34d13ca47a269336e576c81cf007ca740fa9ec0800", "1.6.7", false},
		{"blah", "", true},
	}

	for _, tc := range testCases {
		ver, err := VersionFromSHA(tc.sha)

		if !tc.shouldErr && err != nil {
			t.Fatalf("expected '%v' to not error.", tc.sha)
		}
		if tc.shouldErr && err == nil {
			t.Fatalf("expected '%v' to error.", tc.sha)
		}
		if tc.version != ver {
			t.Fatalf("expected '%v' to result in '%v', got '%v'.", tc.sha, tc.version, ver)
		}
	}
}
