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
		ajd.Lexems = make([]*AbbyyLexem, 0, 1)
	}
	ajd.Lexems = append(ajd.Lexems, lexemData)
}

// Возвращает идентификатор hash table в редис
func (ajd *AbbyyJsonData) SaveToRedis() (string, error) {
	uuidAJD := uuid.NewV4().String()
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
	hashTableAJDHeader["lexems"] = strings.Join(lexemUuidsList, ";")
	err := slRedis.HMSetMap(uuidAJD, hashTableAJDHeader)
	if err != nil {
		return "", err
	}

	return uuidAJD, err
}

func (ajd *AbbyyJsonData) ToServiceMap() *map[string]string {
	properties := make(map[string]string, 1)
	properties["user_lexem"] = ajd.UserLexem
	properties["lang"] = string(ajd.Lang)
	properties["has_ul"] = strconv.FormatBool(ajd.HasUL)
	properties["count_words"] = string(ajd.CountWords)

	return &properties
}

func (ajd *AbbyyJsonData) validateHeader(data *map[string]string) error {
	fields := []string{
		"userlexem", "lang", "countwords", "hasul", "lexems",
	}

	for _, field := range fields {
		_, ok := (*data)[field]
		if !ok {
			return fmt.Errorf("не найдено поле %s", field)
		}
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
	for word, tokens := range al.Words {
		previos, ok := htmap[word]
		if ok {
			allTokens := fmt.Sprintf("%s;%s", previos, strings.Join(tokens, ";"))
			htmap[word] = allTokens
		} else {
			htmap[word] = strings.Join(tokens, ";")
		}
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

func (al *AbbyyLexem) validate(data *map[string]string) error {
	fields := []string{
		"lexem", "paos", "paradigmname", "grammar", "words",
	}

	for _, field := range fields {
		_, ok := (*data)[field]
		if !ok {
			return fmt.Errorf("не найдено поле %s", field)
		}
	}
	return nil
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

func MakeAjdFromRedis(uuids ...string) (*AbbyyJsonData, error) {
	var uuidAjd string

	header := make(map[string]string)

	if len(uuids) > 0 {
		uuidAjd = uuids[0]
	} else {
		uuidAjd = ""
	}
	appErr := slRedis.InitRedisPool()
	if appErr != nil {
		return nil, fmt.Errorf(appErr.Error())
	}
	// если есть uuid ищем сразу по нему, иначе первый ключ в списке
	if len(uuidAjd) == 0 {
		tmpValue, err := slRedis.GetAjdUuid()
		if err != nil {
			return nil, err
		}
		uuidAjd = tmpValue
	}
	// получаем заголовок ajd
	err := slRedis.HGetAll(uuidAjd, &header)
	if err != nil {
		return nil, err
	}
	ajd := new(AbbyyJsonData)
	err = ajd.validateHeader(&header)
	if err != nil {
		return nil, fmt.Errorf("ошибка структуры Abbyy Json Data")
	}
	ajd.UserLexem = header["userlexem"]
	ajd.Lang, err = strconv.Atoi(header["lang"])
	if err != nil {
		return nil, fmt.Errorf("ошибка конвертации в число параметра %s", "Lang")
	}
	ajd.CountWords, err = strconv.Atoi(header["countwords"])
	if err != nil {
		return nil, fmt.Errorf("ошибка конвертации в число параметра %s", "CountWords")
	}
	ajd.HasUL, err = strconv.ParseBool(header["hasul"])
	if err != nil {
		return nil, fmt.Errorf("ошибка конвертации в булево параметра %s", "HasUl")
	}
	for _, key := range strings.Split(header["lexems"], ";") {
		lexem, err := makeLexemFromRedis(key)
		if err != nil {
			return nil, fmt.Errorf("ошибка создания лексемы uuid = %s", key)
		}
		ajd.AddLexem(lexem)
	}

	return ajd, nil
}

// создаем лексему из редиса
func makeLexemFromRedis(uuid string) (*AbbyyLexem, error) {
	appErr := slRedis.InitRedisPool()
	if appErr != nil {
		return nil, fmt.Errorf(appErr.Error())
	}
	lexem := make(map[string]string, 1)
	err := slRedis.HGetAll(uuid, &lexem)
	if err != nil {
		return nil, fmt.Errorf("ошибка загрузки лексемы")
	}
	wordsKey, ok := lexem["words"]
	if !ok {
		return nil, fmt.Errorf("в лексеме не найден ключ %s", "words")
	}

	words, err := makeWordsMapFromRedis(wordsKey)
	if err != nil {
		return nil, fmt.Errorf("ошибка загрузки слов по лексеме %s", uuid)
	}
	// собираем лексему
	al := new(AbbyyLexem)
	err = al.validate(&lexem)
	if err != nil {
		return nil, err
	}
	al.Uuid = uuid
	al.Words = words
	al.Grammar = strings.Split(lexem["grammar"], ";")
	al.ParadigmName = lexem["paradigmname"]
	al.PaoS = lexem["paos"]
	al.Lexem = lexem["lexem"]
	al.Groups = nil

	return al, nil
}

// загружаем слова из редиса
func makeWordsMapFromRedis(uuid string) (map[string][]string, error) {
	appErr := slRedis.InitRedisPool()
	if appErr != nil {
		return nil, fmt.Errorf(appErr.Error())
	}
	tmpWords := make(map[string]string, 1)
	words := make(map[string][]string, 1)

	err := slRedis.HGetAll(uuid, &tmpWords)
	if err != nil {
		return nil, fmt.Errorf("ошибка загрузки слов по ключу %s", uuid)
	}
	for key, wordLine := range tmpWords {
		wordsList := strings.Split(wordLine, ";")
		words[key] = wordsList
	}
	return words, nil
}
