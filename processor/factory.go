package processor

import "sync"

var (
	oneValidateType,
	oneValidateRequired,
	oneValidateMutability,
	oneGenerateId,
	oneGenerateUserMeta,
	oneGenerateGroupMeta,
	oneUpdateMeta,
	oneFormatCase	sync.Once

	validateTypeInstance,
	validateRequiredInstance,
	validateMutabilityInstance,
	generateIdInstance,
	generateUserMetaInstance,
	generateGroupMetaInstance,
	updateMetaInstance,
	formatCaseInstance	Processor
)

func GetValidateTypeProcessor() Processor {
	oneValidateType.Do(func() {
		validateTypeInstance = &typeValidationProcessor{}
	})
	return validateTypeInstance
}

func GetValidateRequiredProcessor() Processor {
	oneValidateRequired.Do(func() {
		validateRequiredInstance = &requiredValidationProcessor{}
	})
	return validateRequiredInstance
}

func GetValidateMutabilityInstance() Processor {
	oneValidateMutability.Do(func() {
		validateMutabilityInstance = &mutabilityValidationProcessor{}
	})
	return validateMutabilityInstance
}

func GetGenerateIdInstance() Processor {
	oneGenerateId.Do(func() {
		generateIdInstance = &generateIdProcessor{}
	})
	return generateIdInstance
}

func GetGenerateUserMetaInstance() Processor {
	oneGenerateUserMeta.Do(func() {
		generateUserMetaInstance = &generateMetaProcessor{
			ResourceType:"User",
			ResourceTypeUri:"/Users",
		}
	})
	return generateUserMetaInstance
}

func GetGenerateGroupMetaInstance() Processor {
	oneGenerateGroupMeta.Do(func() {
		generateGroupMetaInstance = &generateMetaProcessor{
			ResourceType:"Group",
			ResourceTypeUri:"/Groups",
		}
	})
	return generateGroupMetaInstance
}

func GetUpdateMetaInstance() Processor {
	oneUpdateMeta.Do(func() {
		updateMetaInstance = &updateMetaProcessor{}
	})
	return updateMetaInstance
}

func GetFormatCaseInstance() Processor {
	oneFormatCase.Do(func() {
		formatCaseInstance = &formatCaseProcessor{}
	})
	return formatCaseInstance
}