package validation

import "sync"

var typeRuleValidator Validator
var requiredRuleValidator Validator
var mutabilityRuleValidator Validator
var delegateRuleValidator Validator

var oneTypeRule, oneRequiredRule, oneMutabilityRule, oneDelegate sync.Once

func GetValidator() Validator {
	if nil == delegateRuleValidator {
		if nil == typeRuleValidator {
			oneTypeRule.Do(func() {
				typeRuleValidator = &typeRulesValidator{}
			})
		}
		if nil == requiredRuleValidator {
			oneRequiredRule.Do(func() {
				requiredRuleValidator = &requiredRulesValidator{}
			})
		}
		if nil == mutabilityRuleValidator {
			oneMutabilityRule.Do(func() {
				mutabilityRuleValidator = &mutabilityRulesValidator{}
			})
		}
		oneDelegate.Do(func() {
			delegateRuleValidator = &delegateValidator{
				Concurrent:false,
				Delegates:[]interface{}{
					typeRuleValidator,
					requiredRuleValidator,
					mutabilityRuleValidator,
				},
			}
		})
	}
	return delegateRuleValidator
}