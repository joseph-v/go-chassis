package cse

import (
	"strconv"
	"strings"

	"github.com/ServiceComb/go-chassis/core/router/model"
)

// DarkLaunchRule2RouteRule translates dark launch rule to route rule
func DarkLaunchRule2RouteRule(rule *model.DarkLaunchRule) []*model.RouteRule {
	if rule.Type == DarkLaunchTypeRate {
		routes := make([]*model.RouteTag, len(rule.Items))
		for i, v := range rule.Items {
			weight, _ := strconv.Atoi(v.PolicyCondition)
			version := strings.Replace(v.GroupCondition, "version=", "", 1)
			routes[i] = &model.RouteTag{
				Weight: weight,
				Tags:   map[string]string{"version": version},
			}
		}
		return []*model.RouteRule{{
			Routes:     routes,
			Precedence: 1,
		}}
	}
	if rule.Type == DarkLaunchTypeRule {
		rules := make([]*model.RouteRule, len(rule.Items))
		for i, v := range rule.Items {
			con := v.PolicyCondition
			version := strings.Replace(v.GroupCondition, "version=", "", 1)
			match := model.Match{
				HTTPHeaders: map[string]map[string]string{},
				Headers:     map[string]map[string]string{},
			}
			if strings.Contains(con, "!=") {
				match.HTTPHeaders[strings.Split(con, "!=")[0]] = map[string]string{"noEqu": strings.Split(con, "!=")[1]}
				match.Headers[strings.Split(con, "!=")[0]] = map[string]string{"noEqu": strings.Split(con, "!=")[1]}
			} else if strings.Contains(con, ">=") {
				match.HTTPHeaders[strings.Split(con, ">=")[0]] = map[string]string{"noLess": strings.Split(con, ">=")[1]}
				match.Headers[strings.Split(con, ">=")[0]] = map[string]string{"noLess": strings.Split(con, ">=")[1]}
			} else if strings.Contains(con, "<=") {
				match.HTTPHeaders[strings.Split(con, "<=")[0]] = map[string]string{"noGreater": strings.Split(con, "<=")[1]}
				match.Headers[strings.Split(con, "<=")[0]] = map[string]string{"noGreater": strings.Split(con, "<=")[1]}
			} else if strings.Contains(con, "=") {
				match.HTTPHeaders[strings.Split(con, "=")[0]] = map[string]string{"exact": strings.Split(con, "=")[1]}
				match.Headers[strings.Split(con, "=")[0]] = map[string]string{"exact": strings.Split(con, "=")[1]}
			} else if strings.Contains(con, ">") {
				match.HTTPHeaders[strings.Split(con, ">")[0]] = map[string]string{"greater": strings.Split(con, ">")[1]}
				match.Headers[strings.Split(con, ">")[0]] = map[string]string{"greater": strings.Split(con, ">")[1]}
			} else if strings.Contains(con, "<") {
				match.HTTPHeaders[strings.Split(con, "<")[0]] = map[string]string{"less": strings.Split(con, "<")[1]}
				match.Headers[strings.Split(con, "<")[0]] = map[string]string{"less": strings.Split(con, "<")[1]}
			} else if strings.Contains(con, "~") {
				match.HTTPHeaders[strings.Split(con, "~")[0]] = map[string]string{"regex": strings.Split(con, "~")[1]}
				match.Headers[strings.Split(con, "~")[0]] = map[string]string{"regex": strings.Split(con, "~")[1]}
			}
			newRule := &model.RouteRule{
				Routes:     generateRouteTags(strings.Split(version, ",")),
				Match:      match,
				Precedence: 1,
			}
			rules[i] = newRule
		}
		return rules
	}
	return nil
}

// generateRouteTags generate route tags
func generateRouteTags(versions []string) []*model.RouteTag {
	length := len(versions)
	if length == 1 {
		return []*model.RouteTag{{
			Weight: 100,
			Tags:   map[string]string{"version": versions[0]},
		}}
	}

	tags := make([]*model.RouteTag, length)
	for i, v := range versions {
		tags[i] = &model.RouteTag{
			Weight: 100 / length,
			Tags:   map[string]string{"version": v},
		}
	}
	return tags
}
