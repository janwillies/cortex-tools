package analyse

import (
	"sort"

	"github.com/cortexproject/cortex-tools/pkg/rules/rwrulefmt"
	"github.com/pkg/errors"
	"github.com/prometheus/prometheus/promql/parser"
	log "github.com/sirupsen/logrus"
)

type MetricsInRuler struct {
	MetricsUsed    []string            `json:"metricsUsed"`
	OverallMetrics map[string]struct{} `json:"-"`
	RuleGroups     []RuleGroupMetrics  `json:"ruleGroups"`
}

type RuleGroupMetrics struct {
	Namespace   string   `json:"namspace"`
	GroupName   string   `json:"name"`
	Metrics     []string `json:"metrics"`
	ParseErrors []string `json:"parse_errors"`
}

func ParseMetricsInRuleGroup(mir *MetricsInRuler, group rwrulefmt.RuleGroup, ns string) error {
	var (
		ruleMetrics = make(map[string]struct{})
		refMetrics  = make(map[string]struct{})
		parseErrors []error
	)

	for _, rule := range group.Rules {
		if rule.Record.Value != "" {
			ruleMetrics[rule.Record.Value] = struct{}{}
		}

		query := rule.Expr.Value
		expr, err := parser.ParseExpr(query)
		if err != nil {
			parseErrors = append(parseErrors, errors.Wrapf(err, "query=%v", query))
			log.Debugln("msg", "promql parse error", "err", err, "query", query)
			continue
		}

		parser.Inspect(expr, func(node parser.Node, path []parser.Node) error {
			if n, ok := node.(*parser.VectorSelector); ok {
				refMetrics[n.Name] = struct{}{}
			}

			return nil
		})
	}

	// remove defined recording rule metrics in same RG
	for ruleMetric := range ruleMetrics {
		delete(refMetrics, ruleMetric)
	}

	var metricsInGroup []string
	var parseErrs []string

	for metric := range refMetrics {
		if metric == "" {
			continue
		}
		metricsInGroup = append(metricsInGroup, metric)
		mir.OverallMetrics[metric] = struct{}{}
	}
	sort.Strings(metricsInGroup)

	for _, err := range parseErrors {
		parseErrs = append(parseErrs, err.Error())
	}

	mir.RuleGroups = append(mir.RuleGroups, RuleGroupMetrics{
		Namespace:   ns,
		GroupName:   group.Name,
		Metrics:     metricsInGroup,
		ParseErrors: parseErrs,
	})

	return nil
}
