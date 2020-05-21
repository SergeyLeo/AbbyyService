package abbyyJsonParser

import (
	"fmt"
	"github.com/mailru/easyjson/jlexer"
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

// создаем указатель на новый объект AbbyyLexem
func NewAbbyyLexem(lexem string, paos string) *AbbyyLexem {
	if len(lexem) == 0 || len(paos) == 0 {
		return nil // не валидные данные для инициализации
	}

	al := new(AbbyyLexem)
	al.Lexem = strings.ToLower(lexem)
	al.PaoS = strings.ToLower(paos)
	al.Grammar = []string{}
	al.Words = map[string][]string{}

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