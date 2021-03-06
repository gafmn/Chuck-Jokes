# Chuck-Jokes

Простое CLI-приложение для работы с https://api.chucknorris.io/

## Функционал
1) Печать на экран случайной шутки
2) Выгружать n случайных шуток со всех категорий

## Как запускать

Для начала следует сбилдить приложение

```go build -o joker cmd/joker/main.go```

Для того, чтобы напечатать случайную шутку, надо ввести команду

`./joker random`

Для того, чтобы выгрузить n случайных шуток со всех категорий, надо ввести команду

`./joker dump -n <number>`

где `n` задается пользователем.

**Note**: Если пользователь не ввел `n`, то возьмется дефолтное значение (5)

## Заметка
Не все уникальные шутки категорий из API (https://api.chucknorris.io/) имеют количество больше заданного `n`. Возникает **WARNING**, если количество шуток вероятно меньше `n`, потому что нельзя получить количество всех шуток из категории. В таком случае в файл запишется вероятное максимальное количество шуток из заданной категории. 
