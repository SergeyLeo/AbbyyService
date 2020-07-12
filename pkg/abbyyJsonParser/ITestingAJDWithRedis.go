// тестируем сохранение и чтение данных структуры Abbyy Json Data
package abbyyJsonParser

import (
	"fmt"
	"github.com/icrowley/fake"
	"github.com/satori/go.uuid"
	"github.com/stretchr/testify/assert"
	"testing"
)

func SaveAjdRedisTest(t *testing.T) {
	al, err := factoryAbbyyLexemSubstantiv()
	if err != nil {
		t.Errorf(err.Error())
		t.Fail()
	}
	ajd := AbbyyJsonData{
		UserLexem:  "очки",
		Lang:       1049,
		CountWords: 4,
		HasUL:      true,
	}
	ajd.AddLexem(al)
	key, err := ajd.SaveToRedis()
	if err != nil {
		t.Errorf(err.Error())
		t.Fail()
	}
	ajdFromRedis, err := MakeAjdFromRedis(key)
	if err != nil {
		t.Errorf(err.Error())
		t.Fail()
	}

	expectedMap := ajd.ToServiceMap()
	actualMap := ajdFromRedis.ToServiceMap()
	assertAjd(t, expectedMap, actualMap)
}

func factoryAbbyyLexemSubstantiv() (*AbbyyLexem, error) {
	err := fake.SetLang("ru")
	if err != nil {
		return nil, fmt.Errorf("error language selecting")
	}
	uuidAl := uuid.NewV4()
	al := AbbyyLexem{
		Lexem:        "aкция",
		Uuid:         uuidAl.String(),
		PaoS:         "существительное",
		ParadigmName: "акция",
		Grammar: []string{
			"существительное",
			"неодушевленное",
			"женский род",
		},
		Groups: nil,
		Words: map[string][]string{
			"aкция": {"и.п", "ед.ч"},
			"акции": {"и.п.", "ед.ч"},
			"aкции": {"р.п", "ед.ч"},
			"акций": {"р.п.", "мн.ч"},
		},
	}

	return &al, nil
}

func assertAjd(t *testing.T, input, output *map[string]string) bool {
	assertObj := assert.New(t)
	return assertObj.ElementsMatch(input, output, "AJD not equal")
}
