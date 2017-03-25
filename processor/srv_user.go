package processor

import "sync"

var (
	oneGetUserService    sync.Once
	oneCreateUserService sync.Once
	oneDeleteUserService sync.Once
	oneQueryUserService  sync.Once
	oneUpdateUserService sync.Once
	onePatchUserService  sync.Once

	getUserServiceProcessor    Processor
	createUserServiceProcessor Processor
	deleteUserServiceProcessor Processor
	queryUserServiceProcessor  Processor
	updateUserServiceProcessor Processor
	patchUserServiceProcessor  Processor
)

func DeleteUserServiceProcessor() Processor {
	oneDeleteUserService.Do(func() {
		deleteUserServiceProcessor = &ErrorHandlingProcessor{
			Op: []Processor{
				GetWorkerBean(ParamUserDelete),
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
	return deleteUserServiceProcessor
}

func PatchUserServiceProcessor() Processor {
	onePatchUserService.Do(func() {
		patchUserServiceProcessor = &ErrorHandlingProcessor{
			Op: []Processor{
				GetWorkerBean(ParamUserPatch),
				GetWorkerBean(DbUserGetToResource),
				GetWorkerBean(DbUserGetToReference),
				GetWorkerBean(Modification),
				GetWorkerBean(ValidateType),
				GetWorkerBean(ValidateRequired),
				GetWorkerBean(ValidateMutability),
				GetWorkerBean(UpdateMeta),
				GetWorkerBean(DbUserReplace),
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
	return patchUserServiceProcessor
}

func UpdateUserServiceProcessor() Processor {
	oneUpdateUserService.Do(func() {
		updateUserServiceProcessor = &ErrorHandlingProcessor{
			Op: []Processor{
				GetWorkerBean(ParamUserReplace),
				GetWorkerBean(DbUserGetToReference),
				GetWorkerBean(ValidateType),
				GetWorkerBean(ValidateRequired),
				GetWorkerBean(ValidateMutability),
				GetWorkerBean(UpdateMeta),
				GetWorkerBean(DbUserReplace),
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
	return updateUserServiceProcessor
}

func CreateUserServiceProcessor() Processor {
	oneCreateUserService.Do(func() {
		createUserServiceProcessor = &ErrorHandlingProcessor{
			Op: []Processor{
				GetWorkerBean(ParamUserCreate),
				GetWorkerBean(ValidateType),
				GetWorkerBean(ValidateRequired),
				GetWorkerBean(GenerateId),
				GetWorkerBean(GenerateUserMeta),
				GetWorkerBean(DbUserCreate),
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
	return createUserServiceProcessor
}

func GetUserServiceProcessor() Processor {
	oneGetUserService.Do(func() {
		getUserServiceProcessor = &ErrorHandlingProcessor{
			Op: []Processor{
				GetWorkerBean(ParamUserGet),
				GetWorkerBean(DbUserGetToSingleResult),
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
	return getUserServiceProcessor
}

func QueryUserServiceProcessor() Processor {
	oneQueryUserService.Do(func() {
		queryUserServiceProcessor = &ErrorHandlingProcessor{
			Op: []Processor{
				GetWorkerBean(ParamUserQuery),
				GetWorkerBean(ParseFilter),
				GetWorkerBean(DbUserQuery),
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
	return queryUserServiceProcessor
}
