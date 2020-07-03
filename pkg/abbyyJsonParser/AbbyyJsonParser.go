package abbyyJsonParser

import (
	"fmt"
	"github.com/mailru/easyjson/jlexer"
	"github.com/satori/go.uuid"
	slRedis "kallaur.ru/libs/abbyyservice/pkg/redis"
	"strconv"
	"strings"
)

type AbbyyJsonData struct {
	UserLexem  string // желательно перевести в один регистр, то что пришло от пользователя
	Lang       int
	CountWords int
	HasUL      bool          // файл содержит пользовательскую лексему
	Lexems     []*AbbyyLexem // массив лексем в файле
}

// Слова по лексеме
// все слова переводим в нижний регистр
// Lexem and ParadigmName тоже сохраняем в нижнем регистре
type AbbyyLexem struct {
	Lexem        string
	Uuid         string
	PaoS         string
	ParadigmName string
	Grammar      []string
	Groups       []*AbbyyGroup       // отсюда будем доставать слова
	Words        map[string][]string // ключ слово в нижнем регистре, значение токены которые найдены по слову
}

type AbbyyGroup struct {
	Name    string // имя группы
	Columns int
	Rows    int
	Idx     int        // индекс группы в массиве групп
	Data    [][][]byte // массив массивов строк таблицы
}

// ****** AbbyyJsonData methods ******************

func (ajd *AbbyyJsonData) AddLexem(lexemData *AbbyyLexem) {
	count := len(ajd.Lexems)
	if count == 0 {
		ajd.Lexems = make([]*AbbyyLexem, 2)
	}
	ajd.Lexems = append(ajd.Lexems, lexemData)
}

// Возвращает идентификатор hash table в редис
func (ajd *AbbyyJsonData) SaveToRedis() (string, error) {
	uuidAJD := uuid.NewV4().String()
	uuidLexemsList := uuid.NewV4().String()
	lexemUuidsList := make([]string, 0, len(ajd.Lexems))

	hashTableAJDHeader := make(map[string]string, 1)
	hashTableAJDHeader["userlexem"] = ajd.UserLexem
	hashTableAJDHeader["lang"] = strconv.Itoa(ajd.Lang)
	hashTableAJDHeader["countwords"] = strconv.Itoa(ajd.CountWords)
	hashTableAJDHeader["hasul"] = strconv.FormatBool(ajd.HasUL)

	appErr := slRedis.InitRedisPool()
	if appErr != nil {
		return "", fmt.Errorf("code = %s. message = %s", appErr.Code, appErr.Message)
	}
	for _, lexem := range ajd.Lexems {
		err := lexem.saveToRedis()
		if err != nil {
			return "", fmt.Errorf("not_save_lexem: %s", lexem.Lexem)
		}
		lexemUuidsList = append(lexemUuidsList, lexem.Uuid)
	}
	hashTableAJDHeader["lexems"] = uuidLexemsList
	err := slRedis.HMSetMap(uuidAJD, hashTableAJDHeader)
	if err != nil {
		return "", err
	}

	return uuidAJD, err
}

// Если объект полностью создан из данных редиса ставим true во втором значении
// Если вместо указателя на структуру возвращаем nil значит произошла ошибка получения данных
// Если данные получены полностью то возвращаем указатель на заполненную структуру
func (ajd *AbbyyJsonData) MakeFromRedis() *AbbyyJsonData {
	var keyAjdObject string
	appError := slRedis.InitRedisPool()
	if appError != nil {
		return nil
	}
	keyAjd, err := slRedis.GetAjdUuid()
	if err != nil {
		return nil
	}
	err = slRedis.LPop(keyAjd, &keyAjdObject)
	if err != nil {
		return nil
	}
	return nil
}

// ****** AbbyyLexem methods ******************

func (al *AbbyyLexem) AddGrammar(grammar string, sep string) int {
	elements := strings.Split(grammar, sep)
	al.Grammar = make([]string, len(elements), len(elements))
	for _, element := range elements {
		al.Grammar = append(al.Grammar, element)
	}
	return len(elements)
}

func (al *AbbyyLexem) AddParadigmName(name string) {
	al.ParadigmName = strings.ToLower(name)
}

// сохраняем слова из лексемы
func (al *AbbyyLexem) saveWords(uuidHT string) error {
	htmap := make(map[string]string, 1)
	for word, tokents := range al.Words {
		htmap[word] = strings.Join(tokents, ";")
	}
	err := slRedis.HMSetMap(uuidHT, htmap)

	return err
}

func (al *AbbyyLexem) saveToRedis() error {
	ht := make(map[string]string, 1)
	uuidWords := uuid.NewV4().String()

	ht["lexem"] = al.Lexem
	ht["paos"] = al.PaoS
	ht["paradigmname"] = al.ParadigmName
	ht["grammar"] = strings.Join(al.Grammar, ";")
	ht["words"] = uuidWords
	err := al.saveWords(uuidWords)
	if err != nil {
		return fmt.Errorf("words_not_saved")
	}
	err = slRedis.HMSetMap(al.Uuid, ht)
	return err
}

// Добавляем слово и токены по нему
// al.Words уже должны быть инициализированы
func (al *AbbyyLexem) AddWord(word string, tokens string, sep string) {
	elements := strings.Split(tokens, sep)
	_, ok := al.Words[word]
	if !ok {
		al.Words[word] = make([]string, 0, 8)
	}
	for _, element := range elements {
		al.Words[word] = append(al.Words[word], element)
	}
}

// Количество слов собранных из json файла
func (al *AbbyyLexem) CountWords() int {
	return len(al.Words)
}

// ********************* ************************

// возвращаем мапу и uuid объекта

// создаем указатель на новый объект AbbyyLexem
func NewAbbyyLexem(lexem string, paos string) *AbbyyLexem {
	if len(lexem) == 0 || len(paos) == 0 {
		return nil // не валидные данные для инициализации
	}
	uuidObj := uuid.NewV4()
	al := new(AbbyyLexem)
	al.Lexem = strings.ToLower(lexem)
	al.PaoS = strings.ToLower(paos)
	al.Grammar = []string{}
	al.Words = map[string][]string{}
	al.Uuid = uuidObj.String()

	return al
}

// Точка входа для разбора json-а от Abbyy
func MarshalAbbyyJsonData(json string, userLexem string, lang int) (*AbbyyJsonData, error) {
	lexems, err := getLexems(json)
	if err != nil {
		return nil, fmt.Errorf("lexems not found")
	}
	ajd := AbbyyJsonData{
		UserLexem:  strings.ToLower(userLexem),
		Lang:       lang,
		HasUL:      false,
		CountWords: 0,
		Lexems:     lexems,
	}
	for _, lexem := range ajd.Lexems {
		if lexem.CountWords() == 0 {

		}
	}
	fmt.Printf("In Lexems find - %d words", lexems[0].CountWords())
	return &ajd, nil
}

// достаем слова из таблиц групп и добавляем токены
func FetchWords(ajd *AbbyyJsonData) []*jlexer.LexerError {
	var isVerb bool
	for _, lexem := range ajd.Lexems {
		isVerb = lexem.PaoS == "глагол"
		for _, group := range lexem.Groups {
			lErr := parseTable(group, isVerb, &lexem.Words)
			if lErr != nil {
				lexem.Words = map[string][]string{} // обнуляем мапу со словами
				ajd.CountWords = 0
				ajd.HasUL = false
				return lErr
			}
		}
		ajd.CountWords += len(lexem.Words)
		if !ajd.HasUL {
			_, ok := lexem.Words[ajd.UserLexem]
			ajd.HasUL = ok
		}
	}
	return nil
}
