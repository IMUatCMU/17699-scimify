package defaults

import "sync"

var (
	oneIdGeneration,
	oneMetaGeneration,
	oneMetaUpdate,
	oneReadOnlyCopy,
	oneFormatCase sync.Once

	idGenerationInstance,
	metaGenerationInstance,
	metaUpdateInstance,
	readOnlyCopyInstance,
	formatCaseInstance ValueDefaulter
)

func GetIdGenerationValueDefaulter() ValueDefaulter {
	oneIdGeneration.Do(func() {
		idGenerationInstance = &idGenerationValueDefaulter{}
	})
	return idGenerationInstance
}

func GetMetaGenerationValueDefaulter() ValueDefaulter {
	oneMetaGeneration.Do(func() {
		metaGenerationInstance = &metaGenerationValueDefaulter{}
	})
	return metaGenerationInstance
}

func GetMetaUpdateValueDefaulter() ValueDefaulter {
	oneMetaUpdate.Do(func() {
		metaUpdateInstance = &metaUpdateValueDefaulter{}
	})
	return metaUpdateInstance
}

func GetReadOnlyCopyValueDefaulter() ValueDefaulter {
	oneReadOnlyCopy.Do(func() {
		readOnlyCopyInstance = &copyReadOnlyValueDefaulter{}
	})
	return readOnlyCopyInstance
}

func GetFormatCaseInstanceValueDefaulter() ValueDefaulter {
	oneFormatCase.Do(func() {
		formatCaseInstance = &caseFormatValueDefaulter{}
	})
	return formatCaseInstance
}

var (
	oneResourceCreationChain,
	oneResourceUpdateChain sync.Once

	resourceCreationValueDefaulter,
	resourceUpdateValueDefaulter ValueDefaulter
)

func GetResourceCreationValueDefaulter() ValueDefaulter {
	oneResourceCreationChain.Do(func() {
		resourceCreationValueDefaulter = &delegateValueDefaulter{
			Defaulters: []ValueDefaulter{
				GetFormatCaseInstanceValueDefaulter(),
				GetIdGenerationValueDefaulter(),
				GetMetaGenerationValueDefaulter(),
			},
		}
	})
	return resourceCreationValueDefaulter
}

func GetResourceUpdateValueDefaulter() ValueDefaulter {
	oneResourceUpdateChain.Do(func() {
		resourceUpdateValueDefaulter = &delegateValueDefaulter{
			Defaulters: []ValueDefaulter{
				GetFormatCaseInstanceValueDefaulter(),
				GetReadOnlyCopyValueDefaulter(),
				GetMetaUpdateValueDefaulter(),
			},
		}
	})
	return resourceUpdateValueDefaulter
}
