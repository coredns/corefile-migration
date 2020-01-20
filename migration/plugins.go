package migration

import "github.com/coredns/corefile-migration/migration/corefile"

var plugins = map[string]map[string]plugin{
	"kubernetes": {
		"v1": plugin{
			namedOptions: map[string]option{
				"resyncperiod":       {},
				"endpoint":           {},
				"tls":                {},
				"namespaces":         {},
				"labels":             {},
				"pods":               {},
				"endpoint_pod_names": {},
				"upstream":           {},
				"ttl":                {},
				"noendpoints":        {},
				"transfer":           {},
				"fallthrough":        {},
				"ignore":             {},
			},
		},
		"v2": plugin{
			namedOptions: map[string]option{
				"resyncperiod":       {},
				"endpoint":           {},
				"tls":                {},
				"namespaces":         {},
				"labels":             {},
				"pods":               {},
				"endpoint_pod_names": {},
				"upstream":           {},
				"ttl":                {},
				"noendpoints":        {},
				"transfer":           {},
				"fallthrough":        {},
				"ignore":             {},
				"kubeconfig":         {}, // new option
			},
		},
		"v3": plugin{
			namedOptions: map[string]option{
				"resyncperiod": {},
				"endpoint": { // new deprecation
					status: deprecated,
					action: useFirstArgumentOnly,
				},
				"tls":                {},
				"kubeconfig":         {},
				"namespaces":         {},
				"labels":             {},
				"pods":               {},
				"endpoint_pod_names": {},
				"upstream":           {},
				"ttl":                {},
				"noendpoints":        {},
				"transfer":           {},
				"fallthrough":        {},
				"ignore":             {},
			},
		},
		"v4": plugin{
			namedOptions: map[string]option{
				"resyncperiod": {},
				"endpoint": {
					status: ignored,
					action: useFirstArgumentOnly,
				},
				"tls":                {},
				"kubeconfig":         {},
				"namespaces":         {},
				"labels":             {},
				"pods":               {},
				"endpoint_pod_names": {},
				"upstream": { // new deprecation
					status: deprecated,
					action: removeOption,
				},
				"ttl":         {},
				"noendpoints": {},
				"transfer":    {},
				"fallthrough": {},
				"ignore":      {},
			},
		},
		"v5": plugin{
			namedOptions: map[string]option{
				"resyncperiod": { // new deprecation
					status: deprecated,
					action: removeOption,
				},
				"endpoint": {
					status: ignored,
					action: useFirstArgumentOnly,
				},
				"tls":                {},
				"kubeconfig":         {},
				"namespaces":         {},
				"labels":             {},
				"pods":               {},
				"endpoint_pod_names": {},
				"upstream": {
					status: ignored,
					action: removeOption,
				},
				"ttl":         {},
				"noendpoints": {},
				"transfer":    {},
				"fallthrough": {},
				"ignore":      {},
			},
		},
		"v6": plugin{
			namedOptions: map[string]option{
				"resyncperiod": { // new removal
					status: removed,
					action: removeOption,
				},
				"endpoint": {
					status: ignored,
					action: useFirstArgumentOnly,
				},
				"tls":                {},
				"kubeconfig":         {},
				"namespaces":         {},
				"labels":             {},
				"pods":               {},
				"endpoint_pod_names": {},
				"upstream": {
					status: ignored,
					action: removeOption,
				},
				"ttl":         {},
				"noendpoints": {},
				"transfer":    {},
				"fallthrough": {},
				"ignore":      {},
			},
		},
		"v7": plugin{
			namedOptions: map[string]option{
				// resyncperiod removed
				"endpoint": {
					status: ignored,
					action: useFirstArgumentOnly,
				},
				"tls":                {},
				"kubeconfig":         {},
				"namespaces":         {},
				"labels":             {},
				"pods":               {},
				"endpoint_pod_names": {},
				"upstream": {
					status: ignored,
					action: removeOption,
				},
				"ttl":         {},
				"noendpoints": {},
				"transfer":    {},
				"fallthrough": {},
				"ignore":      {},
			},
		},
	},

	"errors": {
		"v1": plugin{},
		"v2": plugin{
			namedOptions: map[string]option{
				"consolidate": {},
			},
		},
	},

	"health": {
		"v1": plugin{
			namedOptions: map[string]option{
				"lameduck": {},
			},
		},
		"v1 add lameduck": plugin{
			namedOptions: map[string]option{
				"lameduck": {
					status: newdefault,
					add: func(c *corefile.Plugin) (*corefile.Plugin, error) {
						return addOptionToPlugin(c, &corefile.Option{Name: "lameduck 5s"})
					},
					downAction: removeOption,
				},
			},
		},
	},

	"hosts": {
		"v1": plugin{
			namedOptions: map[string]option{
				"ttl":         {},
				"no_reverse":  {},
				"reload":      {},
				"fallthrough": {},
			},
			patternOptions: map[string]option{
				`\d{1,3}\.\d{1,3}\.\d{1,3}\.\d{1,3}`:              {}, // close enough
				`[0-9A-Fa-f]{1,4}:[:0-9A-Fa-f]+:[0-9A-Fa-f]{1,4}`: {}, // less close, but still close enough
			},
		},
	},

	"log": {
		"v1": plugin{
			namedOptions: map[string]option{
				"class": {},
			},
		},
	},

	"cache": {
		"v1": plugin{
			namedOptions: map[string]option{
				"success":  {},
				"denial":   {},
				"prefetch": {},
			},
		},
		"v2": plugin{
			namedOptions: map[string]option{
				"success":     {},
				"denial":      {},
				"prefetch":    {},
				"serve_stale": {}, // new option
			},
		},
	},

	"forward": {
		"v1": plugin{
			namedOptions: map[string]option{
				"except":         {},
				"force_tcp":      {},
				"expire":         {},
				"max_fails":      {},
				"tls":            {},
				"tls_servername": {},
				"policy":         {},
				"health_check":   {},
			},
		},
		"v2": plugin{
			namedOptions: map[string]option{
				"except":         {},
				"force_tcp":      {},
				"prefer_udp":     {},
				"expire":         {},
				"max_fails":      {},
				"tls":            {},
				"tls_servername": {},
				"policy":         {},
				"health_check":   {},
			},
		},
	},

	"k8s_external": {
		"v1": plugin{
			namedOptions: map[string]option{
				"apex": {},
				"ttl":  {},
			},
		},
	},
}
