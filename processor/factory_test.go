package processor

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetWorkerBean(t *testing.T) {
	for _, test := range []struct {
		bn BeanName
	}{
		{DbUserCreate},
		{DbGroupCreate},
		{DbUserDelete},
		{DbGroupDelete},
		{DbUserGetToSingleResult},
		{DbUserGetToReference},
		{DbGroupGetToSingleResult},
		{DbGroupGetToReference},
		{DbUserQuery},
		{DbGroupQuery},
		//{DbRootQuery},
		{DbUserReplace},
		{DbGroupReplace},
		{FormatCase},
		{GenerateId},
		{GenerateUserMeta},
		{GenerateGroupMeta},
		{UpdateMeta},
		{JsonSimple},
		{JsonAssisted},
		{JsonHybridList},
		{SetJsonToSingle},
		{SetJsonToMultiple},
		{ValidateType},
		{ValidateRequired},
		{ValidateMutability},
		{TranslateError},
		{ParseFilter},
		{ParamUserGet},
		{ParamGroupGet},
		{ParamUserCreate},
		{ParamGroupCreate},
		{ParamUserDelete},
		{ParamGroupDelete},
		{ParamUserQuery},
		{ParamGroupQuery},
		{ParamRootQuery},
		{ParamUserReplace},
		{ParamGroupReplace},
	} {
		assert.NotNil(t, GetWorkerBean(test.bn))
	}
}
