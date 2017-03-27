package processor

import (
	"fmt"
	"github.com/spf13/viper"
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
			SrvUserCreate:   {processor: CreateUserServiceProcessor(), num: poolSize(SrvUserCreate)},
			SrvUserReplace:  {processor: UpdateUserServiceProcessor(), num: poolSize(SrvUserReplace)},
			SrvUserPatch:    {processor: PatchUserServiceProcessor(), num: poolSize(SrvUserPatch)},
			SrvUserGet:      {processor: GetUserServiceProcessor(), num: poolSize(SrvUserGet)},
			SrvUserQuery:    {processor: QueryUserServiceProcessor(), num: poolSize(SrvUserQuery)},
			SrvUserDelete:   {processor: DeleteUserServiceProcessor(), num: poolSize(SrvUserDelete)},
			SrvGroupCreate:  {processor: CreateGroupServiceProcessor(), num: poolSize(SrvGroupCreate)},
			SrvGroupReplace: {processor: UpdateGroupServiceProcessor(), num: poolSize(SrvGroupReplace)},
			SrvGroupPatch:   {processor: PatchGroupServiceProcessor(), num: poolSize(SrvGroupPatch)},
			SrvGroupGet:     {processor: GetGroupServiceProcessor(), num: poolSize(SrvGroupGet)},
			SrvGroupQuery:   {processor: QueryGroupServiceProcessor(), num: poolSize(SrvGroupQuery)},
			SrvGroupDelete:  {processor: DeleteGroupServiceProcessor(), num: poolSize(SrvGroupDelete)},
			BulkDispatch:    {processor: BulkDispatchProcessor(), num: poolSize(BulkDispatch)},
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

func poolSize(name BeanName) int { return viper.GetInt(fmt.Sprintf("scim.threadPool.%s", name)) }

func GetWorkerBean(bn BeanName) Worker {
	oneWorkerBeanMap.Do(func() {
		workerBeanMap = map[BeanName]bean{
			DbUserCreate:             {processor: DBUserCreateProcessor(), num: poolSize(DbUserCreate)},
			DbGroupCreate:            {processor: DBGroupCreateProcessor(), num: poolSize(DbGroupCreate)},
			DbUserDelete:             {processor: DBUserDeleteProcessor(), num: poolSize(DbUserDelete)},
			DbGroupDelete:            {processor: DBGroupDeleteProcessor(), num: poolSize(DbGroupDelete)},
			DbUserGetToSingleResult:  {processor: DBUserGetToSingleResultProcessor(), num: poolSize(DbUserGetToSingleResult)},
			DbUserGetToResource:      {processor: DBUserGetToResourceProcessor(), num: poolSize(DbUserGetToResource)},
			DbUserGetToReference:     {processor: DBUserGetToReferenceProcessor(), num: poolSize(DbUserGetToReference)},
			DbGroupGetToSingleResult: {processor: DBGroupGetToSingleResultProcessor(), num: poolSize(DbGroupGetToSingleResult)},
			DbGroupGetToResource:     {processor: DBUserGetToResourceProcessor(), num: poolSize(DbGroupGetToResource)},
			DbGroupGetToReference:    {processor: DBGroupGetToReferenceProcessor(), num: poolSize(DbGroupGetToReference)},
			DbUserQuery:              {processor: DBUserQueryProcessor(), num: poolSize(DbUserQuery)},
			DbGroupQuery:             {processor: DBGroupQueryProcessor(), num: poolSize(DbGroupQuery)},
			DbRootQuery:              {processor: DBRootQueryProcessor(), num: poolSize(DbRootQuery)},
			DbUserReplace:            {processor: DBUserReplaceProcessor(), num: poolSize(DbUserReplace)},
			DbGroupReplace:           {processor: DBGroupReplaceProcessor(), num: poolSize(DbGroupReplace)},
			DbSPConfigGet:            {processor: DBSPConfigGetProcessor(), num: poolSize(DbSPConfigGet)},
			DbSchemaGet:              {processor: DBSchemaGetProcessor(), num: poolSize(DbSchemaGet)},
			DbResourceTypeGetAll:     {processor: DbGetAllResourceTypesProcessor(), num: poolSize(DbResourceTypeGetAll)},
			DbSchemaGetAll:           {processor: DbGetAllSchemasProcessor(), num: poolSize(DbSchemaGetAll)},
			FormatCase:               {processor: FormatCaseProcessor(), num: poolSize(FormatCase)},
			GenerateId:               {processor: GenerateIdProcessor(), num: poolSize(GenerateId)},
			GenerateUserMeta:         {processor: GenerateUserMetaProcessor(), num: poolSize(GenerateUserMeta)},
			GenerateGroupMeta:        {processor: GenerateGroupMetaProcessor(), num: poolSize(GenerateGroupMeta)},
			UpdateMeta:               {processor: UpdateMetaProcessor(), num: poolSize(UpdateMeta)},
			Modification:             {processor: ModificationProcessor(), num: poolSize(Modification)},
			JsonSimple:               {processor: SimpleJsonSerializationProcessor(), num: poolSize(JsonSimple)},
			JsonAssisted:             {processor: AssistedJsonSerializationProcessor(), num: poolSize(JsonAssisted)},
			JsonHybridList:           {processor: ListResponseJsonSerializationProcessor(), num: poolSize(JsonHybridList)},
			SetJsonToSingle:          {processor: SingleResultAsJsonTargetProcessor(), num: poolSize(SetJsonToSingle)},
			SetJsonToMultiple:        {processor: MultipleResultAsJsonTargetProcessor(), num: poolSize(SetJsonToMultiple)},
			SetJsonToError:           {processor: ErrorAsJsonTargetProcessor(), num: poolSize(SetJsonToError)},
			SetJsonToResource:        {processor: ResourceAsJsonTargetProcessor(), num: poolSize(SetJsonToResource)},
			ValidateType:             {processor: TypeValidationProcessor(), num: poolSize(ValidateType)},
			ValidateRequired:         {processor: RequiredValidationProcessor(), num: poolSize(ValidateRequired)},
			ValidateMutability:       {processor: MutabilityValidationProcessor(), num: poolSize(ValidateMutability)},
			TranslateError:           {processor: ErrorTranslatingProcessor(), num: poolSize(TranslateError)},
			ParseFilter:              {processor: ParseFilterProcessor(), num: poolSize(ParseFilter)},
			ParamSchemaGet:           {processor: ParseParamForSchemaGetEndpointProcessor(), num: poolSize(ParamSchemaGet)},
			ParamUserGet:             {processor: ParseParamForUserGetEndpointProcessor(), num: poolSize(ParamUserGet)},
			ParamGroupGet:            {processor: ParseParamForGroupGetEndpointProcessor(), num: poolSize(ParamGroupGet)},
			ParamUserCreate:          {processor: ParseParamForUserCreateEndpointProcessor(), num: poolSize(ParamUserCreate)},
			ParamGroupCreate:         {processor: ParseParamForGroupCreateEndpointProcessor(), num: poolSize(ParamGroupCreate)},
			ParamUserDelete:          {processor: ParseParamForUserDeleteEndpointProcessor(), num: poolSize(ParamUserDelete)},
			ParamGroupDelete:         {processor: ParseParamForGroupDeleteEndpointProcessor(), num: poolSize(ParamGroupDelete)},
			ParamUserQuery:           {processor: ParseParamForUserQueryEndpointProcessor(), num: poolSize(ParamUserQuery)},
			ParamGroupQuery:          {processor: ParseParamForGroupQueryEndpointProcessor(), num: poolSize(ParamGroupQuery)},
			ParamRootQuery:           {processor: ParseParamForRootQueryEndpointProcessor(), num: poolSize(ParamRootQuery)},
			ParamUserReplace:         {processor: ParseParamForUserReplaceEndpointProcessor(), num: poolSize(ParamUserReplace)},
			ParamGroupReplace:        {processor: ParseParamForGroupReplaceEndpointProcessor(), num: poolSize(ParamGroupReplace)},
			ParamUserPatch:           {processor: ParseParamForUserPatchEndpointProcessor(), num: poolSize(ParamUserPatch)},
			ParamGroupPatch:          {processor: ParseParamForGroupPatchEndpointProcessor(), num: poolSize(ParamGroupPatch)},
			ParamBulk:                {processor: ParseParamForBulkEndpointProcessor(), num: poolSize(ParamBulk)},
			SetAllHeader:             {processor: SetAllHeaderProcessor(), num: poolSize(SetAllHeader)},
			SetStatusToError:         {processor: SetStatusToErrorProcessor(), num: poolSize(SetStatusToError)},
			SetStatusToOk:            {processor: SetStatusToOKProcessor(), num: poolSize(SetStatusToOk)},
			SetStatusToCreated:       {processor: SetStatusToCreatedProcessor(), num: poolSize(SetStatusToCreated)},
			SetStatusToNoContent:     {processor: SetStatusToNoContentProcessor(), num: poolSize(SetStatusToNoContent)},
		}
	})
	if b, ok := workerBeanMap[bn]; !ok {
		log.Panicf("No bean by the name %s", bn)
		return nil
	} else {
		b.once.Do(func() {
			if b.num == 1 {
				b.worker = &SimpleWorker{processor: b.processor}
			} else {
				b.worker = &WorkerWrapper{processor: b.processor}
				b.worker.initialize(b.num)
			}
		})
		return b.worker
	}
}
