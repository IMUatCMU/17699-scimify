package processor

import "sync"

var (
	oneGetGroupService    sync.Once
	oneCreateGroupService sync.Once
	oneDeleteGroupService sync.Once
	oneQueryGroupService  sync.Once
	oneUpdateGroupService sync.Once
	onePatchGroupService  sync.Once

	getGroupServiceProcessor    Processor
	createGroupServiceProcessor Processor
	deleteGroupServiceProcessor Processor
	queryGroupServiceProcessor  Processor
	updateGroupServiceProcessor Processor
	patchGroupServiceProcessor  Processor
)

func GetGroupServiceProcessor() Processor {
	oneGetGroupService.Do(func() {
		getGroupServiceProcessor = &ErrorHandlingProcessor{
			Op: []Processor{
				GetWorkerBean(ParamGroupGet),
				GetWorkerBean(DbGroupGetToSingleResult),
				GetWorkerBean(SetJsonToSingle),
				GetWorkerBean(SetAllHeader),
				GetWorkerBean(JsonAssisted),
				GetWorkerBean(SetStatusToOk),
			},
			ErrOp: []Processor{
				GetWorkerBean(TranslateError),
				GetWorkerBean(SetJsonToError),
				GetWorkerBean(JsonSimple),
				GetWorkerBean(SetStatusToError),
			},
		}
	})
	return getGroupServiceProcessor
}

func CreateGroupServiceProcessor() Processor {
	oneCreateGroupService.Do(func() {
		createGroupServiceProcessor = &ErrorHandlingProcessor{
			Op: []Processor{
				GetWorkerBean(ParamGroupCreate),
				GetWorkerBean(ValidateType),
				GetWorkerBean(ValidateRequired),
				GetWorkerBean(GenerateId),
				GetWorkerBean(GenerateGroupMeta),
				GetWorkerBean(DbGroupCreate),
				GetWorkerBean(SetJsonToResource),
				GetWorkerBean(SetAllHeader),
				GetWorkerBean(JsonAssisted),
				GetWorkerBean(SetStatusToCreated),
			},
			ErrOp: []Processor{
				GetWorkerBean(TranslateError),
				GetWorkerBean(SetJsonToError),
				GetWorkerBean(JsonSimple),
				GetWorkerBean(SetStatusToError),
			},
		}
	})
	return createGroupServiceProcessor
}

func DeleteGroupServiceProcessor() Processor {
	oneDeleteGroupService.Do(func() {
		deleteGroupServiceProcessor = &ErrorHandlingProcessor{
			Op: []Processor{
				GetWorkerBean(ParamGroupDelete),
				GetWorkerBean(DbUserDelete),
				GetWorkerBean(SetStatusToNoContent),
			},
			ErrOp: []Processor{
				GetWorkerBean(TranslateError),
				GetWorkerBean(SetJsonToError),
				GetWorkerBean(JsonSimple),
				GetWorkerBean(SetStatusToError),
			},
		}
	})
	return deleteGroupServiceProcessor
}

func QueryGroupServiceProcessor() Processor {
	oneQueryGroupService.Do(func() {
		queryGroupServiceProcessor = &ErrorHandlingProcessor{
			Op: []Processor{
				GetWorkerBean(ParamGroupQuery),
				GetWorkerBean(ParseFilter),
				GetWorkerBean(DbGroupQuery),
				GetWorkerBean(SetJsonToMultiple),
				GetWorkerBean(JsonHybridList),
				GetWorkerBean(SetStatusToOk),
			},
			ErrOp: []Processor{
				GetWorkerBean(TranslateError),
				GetWorkerBean(SetJsonToError),
				GetWorkerBean(JsonSimple),
				GetWorkerBean(SetStatusToError),
			},
		}
	})
	return queryGroupServiceProcessor
}

func UpdateGroupServiceProcessor() Processor {
	oneUpdateGroupService.Do(func() {
		updateGroupServiceProcessor = &ErrorHandlingProcessor{
			Op: []Processor{
				GetWorkerBean(ParamGroupReplace),
				GetWorkerBean(DbGroupGetToReference),
				GetWorkerBean(ValidateType),
				GetWorkerBean(ValidateRequired),
				GetWorkerBean(ValidateMutability),
				GetWorkerBean(UpdateMeta),
				GetWorkerBean(DbGroupReplace),
				GetWorkerBean(SetJsonToResource),
				GetWorkerBean(SetAllHeader),
				GetWorkerBean(JsonAssisted),
				GetWorkerBean(SetStatusToOk),
			},
			ErrOp: []Processor{
				GetWorkerBean(TranslateError),
				GetWorkerBean(SetJsonToError),
				GetWorkerBean(JsonSimple),
				GetWorkerBean(SetStatusToError),
			},
		}
	})
	return updateGroupServiceProcessor
}

func PatchGroupServiceProcessor() Processor {
	onePatchGroupService.Do(func() {
		patchGroupServiceProcessor = &ErrorHandlingProcessor{
			Op: []Processor{
				GetWorkerBean(ParamGroupPatch),
				GetWorkerBean(DbGroupGetToResource),
				GetWorkerBean(DbGroupGetToReference),
				GetWorkerBean(Modification),
				GetWorkerBean(ValidateType),
				GetWorkerBean(ValidateRequired),
				GetWorkerBean(ValidateMutability),
				GetWorkerBean(UpdateMeta),
				GetWorkerBean(DbGroupReplace),
				GetWorkerBean(SetJsonToResource),
				GetWorkerBean(SetAllHeader),
				GetWorkerBean(JsonAssisted),
				GetWorkerBean(SetStatusToOk),
			},
			ErrOp: []Processor{
				GetWorkerBean(TranslateError),
				GetWorkerBean(SetJsonToError),
				GetWorkerBean(JsonSimple),
				GetWorkerBean(SetStatusToError),
			},
		}
	})
	return patchGroupServiceProcessor
}