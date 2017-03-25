package processor

import (
	"log"
	"sync"
)

type BeanName string

const (
	SrvUserCreate            = BeanName("SrvUserCreate")
	SrvUserReplace           = BeanName("SrvUserReplace")
	SrvUserPatch             = BeanName("SrvUserPatch")
	SrvUserGet               = BeanName("SrvUserGet")
	SrvUserQuery             = BeanName("SrvUserQuery")
	SrvUserDelete            = BeanName("SrvUserDelete")
	SrvGroupCreate           = BeanName("SrvGroupCreate")
	SrvGroupReplace          = BeanName("SrvGroupReplace")
	SrvGroupPatch            = BeanName("SrvGroupPatch")
	SrvGroupGet              = BeanName("SrvGroupGet")
	SrvGroupQuery            = BeanName("SrvGroupQuery")
	SrvGroupDelete           = BeanName("SrvGroupDelete")
	DbUserCreate             = BeanName("DbUserCreate")
	DbGroupCreate            = BeanName("DbGroupCreate")
	DbUserDelete             = BeanName("DbUserDelete")
	DbGroupDelete            = BeanName("DbGroupDelete")
	DbUserGetToSingleResult  = BeanName("DbUserGetToSingleResult")
	DbUserGetToResource      = BeanName("DbUserGetToResource")
	DbUserGetToReference     = BeanName("DbUserGetToReference")
	DbGroupGetToSingleResult = BeanName("DbGroupGetToSingleResult")
	DbGroupGetToResource     = BeanName("DbGroupGetToResource")
	DbGroupGetToReference    = BeanName("DbGroupGetToReference")
	DbUserQuery              = BeanName("DbUserQuery")
	DbGroupQuery             = BeanName("DbGroupQuery")
	DbRootQuery              = BeanName("DbRootQuery")
	DbUserReplace            = BeanName("DbUserReplace")
	DbGroupReplace           = BeanName("DbGroupReplace")
	DbSPConfigGet            = BeanName("DbSPConfigGet")
	DbSchemaGet              = BeanName("DbSchemaGet")
	DbResourceTypeGetAll     = BeanName("DbResourceTypeGetAll")
	DbSchemaGetAll           = BeanName("DbSchemaGetAll")
	FormatCase               = BeanName("FormatCase")
	GenerateId               = BeanName("GenerateId")
	GenerateUserMeta         = BeanName("GenerateUserMeta")
	GenerateGroupMeta        = BeanName("GenerateGroupMeta")
	UpdateMeta               = BeanName("UpdateMeta")
	Modification             = BeanName("Modification")
	JsonSimple               = BeanName("JsonSimple")
	JsonAssisted             = BeanName("JsonAssisted")
	JsonHybridList           = BeanName("JsonHybridList")
	SetJsonToSingle          = BeanName("SetJsonToSingle")
	SetJsonToMultiple        = BeanName("SetJsonToMultiple")
	SetJsonToError           = BeanName("SetJsonToError")
	SetJsonToResource        = BeanName("SetJsonToResource")
	ValidateType             = BeanName("ValidateType")
	ValidateRequired         = BeanName("ValidateRequired")
	ValidateMutability       = BeanName("ValidateMutability")
	TranslateError           = BeanName("TranslateError")
	ParseFilter              = BeanName("ParseFilter")
	ParamSchemaGet           = BeanName("ParamSchemaGet")
	ParamUserGet             = BeanName("ParamUserGet")
	ParamGroupGet            = BeanName("ParamGroupGet")
	ParamUserCreate          = BeanName("ParamUserCreate")
	ParamGroupCreate         = BeanName("ParamGroupCreate")
	ParamUserDelete          = BeanName("ParamUserDelete")
	ParamGroupDelete         = BeanName("ParamGroupDelete")
	ParamUserQuery           = BeanName("ParamUserQuery")
	ParamGroupQuery          = BeanName("ParamGroupQuery")
	ParamRootQuery           = BeanName("ParamRootQuery")
	ParamUserReplace         = BeanName("ParamUserReplace")
	ParamGroupReplace        = BeanName("ParamGroupReplace")
	ParamUserPatch           = BeanName("ParamUserPatch")
	ParamGroupPatch          = BeanName("ParamGroupPatch")
	ParamBulk                = BeanName("ParamBulk")
	BulkDispatch             = BeanName("BulkDispatch")
	SetAllHeader             = BeanName("SetAllHeader")
	SetStatusToError         = BeanName("SetStatusToError")
	SetStatusToOk            = BeanName("SetStatusToOk")
	SetStatusToCreated       = BeanName("SetStatusToCreated")
	SetStatusToNoContent     = BeanName("SetStatusToNoContent")
)

type bean struct {
	processor Processor
	num       int
	once      sync.Once
	worker    Worker
}

var (
	workerBeanMap    map[BeanName]bean
	oneWorkerBeanMap sync.Once

	serviceBeanMap    map[BeanName]bean
	oneServiceBeanMap sync.Once
)

func GetServiceBean(bn BeanName) Worker {
	oneServiceBeanMap.Do(func() {
		serviceBeanMap = map[BeanName]bean{
			SrvUserCreate:   {processor: CreateUserServiceProcessor(), num: 2},
			SrvUserReplace:  {processor: UpdateUserServiceProcessor(), num: 2},
			SrvUserPatch:    {processor: PatchUserServiceProcessor(), num: 2},
			SrvUserGet:      {processor: GetUserServiceProcessor(), num: 2},
			SrvUserQuery:    {processor: QueryUserServiceProcessor(), num: 2},
			SrvUserDelete:   {processor: DeleteUserServiceProcessor(), num: 2},
			SrvGroupCreate:  {processor: CreateGroupServiceProcessor(), num: 2},
			SrvGroupReplace: {processor: UpdateGroupServiceProcessor(), num: 2},
			SrvGroupPatch:   {processor: PatchGroupServiceProcessor(), num: 2},
			SrvGroupGet:     {processor: GetGroupServiceProcessor(), num: 2},
			SrvGroupQuery:   {processor: QueryGroupServiceProcessor(), num: 2},
			SrvGroupDelete:  {processor: DeleteGroupServiceProcessor(), num: 2},
			BulkDispatch:    {processor: BulkDispatchProcessor(), num: 2},
		}
	})
	if b, ok := serviceBeanMap[bn]; !ok {
		log.Panicf("No bean by the name %s", bn)
		return nil
	} else {
		b.once.Do(func() {
			b.worker = &WorkerWrapper{processor: b.processor}
			b.worker.initialize(b.num)
		})
		return b.worker
	}
}

func GetWorkerBean(bn BeanName) Worker {
	oneWorkerBeanMap.Do(func() {
		workerBeanMap = map[BeanName]bean{
			DbUserCreate:             {processor: DBUserCreateProcessor(), num: 2},
			DbGroupCreate:            {processor: DBGroupCreateProcessor(), num: 2},
			DbUserDelete:             {processor: DBUserDeleteProcessor(), num: 2},
			DbGroupDelete:            {processor: DBGroupDeleteProcessor(), num: 2},
			DbUserGetToSingleResult:  {processor: DBUserGetToSingleResultProcessor(), num: 2},
			DbUserGetToResource:      {processor: DBUserGetToResourceProcessor(), num: 2},
			DbUserGetToReference:     {processor: DBUserGetToReferenceProcessor(), num: 2},
			DbGroupGetToSingleResult: {processor: DBGroupGetToSingleResultProcessor(), num: 2},
			DbGroupGetToResource:     {processor: DBUserGetToResourceProcessor(), num: 2},
			DbGroupGetToReference:    {processor: DBGroupGetToReferenceProcessor(), num: 2},
			DbUserQuery:              {processor: DBUserQueryProcessor(), num: 2},
			DbGroupQuery:             {processor: DBGroupQueryProcessor(), num: 2},
			DbRootQuery:              {processor: DBRootQueryProcessor(), num: 2},
			DbUserReplace:            {processor: DBUserReplaceProcessor(), num: 2},
			DbGroupReplace:           {processor: DBGroupReplaceProcessor(), num: 2},
			DbSPConfigGet:            {processor: DBSPConfigGetProcessor(), num: 1},
			DbSchemaGet:              {processor: DBSchemaGetProcessor(), num: 1},
			DbResourceTypeGetAll:     {processor: DbGetAllResourceTypesProcessor(), num: 1},
			DbSchemaGetAll:           {processor: DbGetAllSchemasProcessor(), num: 1},
			FormatCase:               {processor: FormatCaseProcessor(), num: 2},
			GenerateId:               {processor: GenerateIdProcessor(), num: 2},
			GenerateUserMeta:         {processor: GenerateUserMetaProcessor(), num: 2},
			GenerateGroupMeta:        {processor: GenerateGroupMetaProcessor(), num: 2},
			UpdateMeta:               {processor: UpdateMetaProcessor(), num: 2},
			Modification:             {processor: ModificationProcessor(), num: 2},
			JsonSimple:               {processor: SimpleJsonSerializationProcessor(), num: 2},
			JsonAssisted:             {processor: AssistedJsonSerializationProcessor(), num: 2},
			JsonHybridList:           {processor: ListResponseJsonSerializationProcessor(), num: 2},
			SetJsonToSingle:          {processor: SingleResultAsJsonTargetProcessor(), num: 2},
			SetJsonToMultiple:        {processor: MultipleResultAsJsonTargetProcessor(), num: 2},
			SetJsonToError:           {processor: ErrorAsJsonTargetProcessor(), num: 2},
			SetJsonToResource:        {processor: ResourceAsJsonTargetProcessor(), num: 2},
			ValidateType:             {processor: TypeValidationProcessor(), num: 2},
			ValidateRequired:         {processor: RequiredValidationProcessor(), num: 2},
			ValidateMutability:       {processor: MutabilityValidationProcessor(), num: 2},
			TranslateError:           {processor: ErrorTranslatingProcessor(), num: 2},
			ParseFilter:              {processor: ParseFilterProcessor(), num: 2},
			ParamSchemaGet:           {processor: ParseParamForSchemaGetEndpointProcessor(), num: 1},
			ParamUserGet:             {processor: ParseParamForUserGetEndpointProcessor(), num: 2},
			ParamGroupGet:            {processor: ParseParamForGroupGetEndpointProcessor(), num: 2},
			ParamUserCreate:          {processor: ParseParamForUserCreateEndpointProcessor(), num: 2},
			ParamGroupCreate:         {processor: ParseParamForGroupCreateEndpointProcessor(), num: 2},
			ParamUserDelete:          {processor: ParseParamForUserDeleteEndpointProcessor(), num: 2},
			ParamGroupDelete:         {processor: ParseParamForGroupDeleteEndpointProcessor(), num: 2},
			ParamUserQuery:           {processor: ParseParamForUserQueryEndpointProcessor(), num: 2},
			ParamGroupQuery:          {processor: ParseParamForGroupQueryEndpointProcessor(), num: 2},
			ParamRootQuery:           {processor: ParseParamForRootQueryEndpointProcessor(), num: 2},
			ParamUserReplace:         {processor: ParseParamForUserReplaceEndpointProcessor(), num: 2},
			ParamGroupReplace:        {processor: ParseParamForGroupReplaceEndpointProcessor(), num: 2},
			ParamUserPatch:           {processor: ParseParamForUserPatchEndpointProcessor(), num: 2},
			ParamGroupPatch:          {processor: ParseParamForGroupPatchEndpointProcessor(), num: 2},
			ParamBulk:                {processor: ParseParamForBulkEndpointProcessor(), num: 2},
			SetAllHeader:             {processor: SetAllHeaderProcessor(), num: 2},
			SetStatusToError:         {processor: SetStatusToErrorProcessor(), num: 2},
			SetStatusToOk:            {processor: SetStatusToOKProcessor(), num: 2},
			SetStatusToCreated:       {processor: SetStatusToCreatedProcessor(), num: 2},
			SetStatusToNoContent:     {processor: SetStatusToNoContentProcessor(), num: 2},
		}
	})
	if b, ok := workerBeanMap[bn]; !ok {
		log.Panicf("No bean by the name %s", bn)
		return nil
	} else {
		b.once.Do(func() {
			b.worker = &WorkerWrapper{processor: b.processor}
			b.worker.initialize(b.num)
		})
		return b.worker
	}
}
