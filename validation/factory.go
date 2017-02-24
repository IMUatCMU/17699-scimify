package validation

import "sync"

var typeRuleValidator Validator
var requiredRuleValidator Validator
var mutabilityRuleValidator Validator
var delegateRuleValidator Validator

var oneTypeRule, oneRequiredRule, oneMutabilityRule, oneDelegate sync.Once

func getTypeRuleValidator() Validator {
	oneTypeRule.Do(func() {
		typeRuleValidator = &typeRulesValidator{}
	})
	return typeRuleValidator
}

func getRequiredRuleValidator() Validator {
	oneRequiredRule.Do(func() {
		requiredRuleValidator = &requiredRulesValidator{}
	})
	return requiredRuleValidator
}

func getMutabilityRuleValidator() Validator {
	oneMutabilityRule.Do(func() {
		mutabilityRuleValidator = &mutabilityRulesValidator{}
	})
	return mutabilityRuleValidator
}

func GetValidator() Validator {
	oneDelegate.Do(func() {
		delegateRuleValidator = &delegateValidator{
			Concurrent: false,
			Delegates: []interface{}{
				getTypeRuleValidator(),
				getRequiredRuleValidator(),
				getMutabilityRuleValidator(),
			},
		}
	})
	return delegateRuleValidator
}

var resourceCreationValidator Validator
var oneResourceCreation sync.Once

func GetResourceCreationValidator() Validator {
	oneResourceCreation.Do(func() {
		resourceCreationValidator = &delegateValidator{
			Concurrent: false,
			Delegates: []interface{}{
				getTypeRuleValidator(),
				getRequiredRuleValidator(),
			},
		}
	})
	return resourceCreationValidator
}

var resourceUpdateValidator Validator
var oneResourceUpdate sync.Once

func GetResourceUpdateValidator() Validator {
	oneResourceUpdate.Do(func() {
		resourceUpdateValidator = &delegateValidator{
			Concurrent: false,
			Delegates: []interface{}{
				getTypeRuleValidator(),
				getRequiredRuleValidator(),
				getMutabilityRuleValidator(),
			},
		}
	})
	return resourceUpdateValidator
}
