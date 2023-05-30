package main

import "fmt"

// codeToString возвращает текстовое описание ошибки банковского терминала
func codeToString(code uint32) string {
	codes := map[uint32]string{
		12:   "Эта версия не поддерживает режим РС-3",
		36:   "В пинпаде нет ключа в ячейке 9",
		99:   "Пинпад не подключен",
		101:  "Операция не поддерживается.",
		115:  "Библиотека занята другим процессом, требуется подождать его завершения.",
		233:  "Пинпад не подключен",
		238:  "Пинпад отключился",
		248:  "Динамическая память закончилась",
		249:  "На терминал передана команда не содержащая обязательные параметры",
		250:  "Внутренняя ошибка: операция отменена Internal error",
		252:  "Внутренняя ошибка: операция не поддерживается Internal error",
		253:  "Аппаратный сбой. Устройство ещё не готово. Hardware failure",
		361:  "Нарушился контакт с чипом карты.",
		362:  "Карта не читается Card is not responding",
		363:  "Карта не читается. Попробуйте магн.ридер",
		364:  "Карта не читается",
		402:  "Карта не была выдана. Изымите карту!",
		403:  "ПИН неверен",
		405:  "ПИН блокирован",
		500:  "Карта терминала старой версии!",
		505:  "Карта терминала заполнена. Выполните инкассацию",
		507:  "Срок действия карты истек",
		514:  "На терминале установлена неверная дата",
		521:  "На карте недостаточно средств",
		561:  "Нарушен список операций на карте. Изымите карту!",
		579:  "Карта блокирована. Изымите карту!",
		584:  "Сегодня по этой карте больше операций делать нельзя",
		585:  "Период обслуживания истек",
		586:  "Превышен лимит, разрешенный без связи с банком",
		705:  "Карта блокирована. Изымите карту!",
		706:  "Карта блокирована",
		708:  "ПИН зачисления заблокирован",
		709:  "ПИН списания заблокирован",
		2000: "Операция прервана клиентом",
		2002: "Превышено время ожидания",
		2008: "Для этой карты операция запрещена",
		3162: "Срок действия карты СБЕРКАРТ окончен. Эта ошибка не должна вылезти в проме.",
		4073: "Биосканер не подключен",
		4100: "Нет связи с банком",
		4101: "На терминале нет стоп-листа. Выполните инкассацию",
		4102: "На терминале нет таблицы комиссий",
		4104: "Неверный ответ на команду",
		4106: "ПИН неверен",
		4107: "ПИН блокирован",
		4108: "Номер карты неверен",
		4110: "Карта терм.заполнена.Выполните инкассацию",
		4111: "Стоп-лист устарел. Выполните инкассацию",
		4112: "Неверный стоп-лист. Выполните инкассацию",
		4113: "Превышен лимит операций за сутки",
		4115: "Для таких карт ручной ввод запрещен",
		4116: "Цифры не совпадают!",
		4117: "Клиент отказался от ввода ПИНа",
		4118: "Операции не найдены",
		4119: "Нет связи с банком",
		4120: "Пинпад не подключен или не загружены ключи",
		4121: "Терминал неисправен!",
		4122: "Ошибка смены ключей!",
		4123: "Сначала выполните сверку итогов",
		4124: "Не загружены ключи",
		4125: "На карте есть чип. Вставьте карту чипом",
		4128: "Ошибка настройки терминала",
		4130: "Память заполнена. Сделайте сверку ито- гов или инкассацию.",
		4131: "Пинпад был заменен. Выполните загрузку параметров",
		4132: "Операция отклонена картой Transaction declined by card",
		4133: "Неверный код ответа по протоколу VISA2",
		4134: "Сначала выполните сверку итогов Totals required",
		4135: "Неверно настроены отделы в терминале",
		4136: "Требуется более свежая версия прошивки в пинпаде",
		4137: "ПИНы не совпадают. Попробуйте еще раз.",
		4138: "Карта отправителя и получателя не могут совпадать.",
		4139: "Нет адреса для связи.",
		4148: "Карта в стоп-листе!",
		4149: "На карте нет имени держателя",
		4150: "Превышен лимит операций",
		4151: "Срок действия карты истек",
		4157: "Превышена максимальная сумма операции.",
		4159: "Валюта операции не поддерживается бесконтактным ридером",
		4174: "Файл не найден",
		4175: "Слишком большой файл",
		4176: "Неизвестная версия Vivopay",
		4185: "Неверная карта администратора",
		4186: "Ключ уже введен!",
		4187: "Неверный номер карты",
		4188: "Неверный срок действия карты",
		4189: "Недопустимое значение!",
		4190: "Карта не читается. Попробуйте магн.ридер",
		4203: "Терминал не зарегистрирован",
		4204: "Внутренняя ошибка сервера",
		4205: "Ошибка связи с хостом",
		4206: "Нарушение протокола",
		4207: "Нарушение формата сообщений",
		4208: "Ошибка базы данных",
		4209: "Некорректные данные",
		4210: "Ошибка шифрования",
		4211: "Отсутствует ключ",
		4213: "Сервер PSDB слишком нагружен. Повторите позже.",
		4220: "Не указан код региона для удаленной загрузки",
		4221: "Не удалось восстановить связь с ККМ после удаленной загрузки",
		4222: "Память заполнена. Необходимо отправить чеки на сервер",
		4300: "От ККМ поступило недостаточно параметров",
		4303: "Мы принимаем только Visa",
		4311: "Операция не найдена",
		4313: "Номер карты не соответствует исходному",
		4314: "Это не карта СБЕРКАРТ",
		4315: "Разрешены только отмены в текущей смене",
		4319: "Сумма не должна превышать 42 млн.",
		4323: "Номер карты не совпадает с исходным",
		4325: "Сумма не указана!",
		4326: "Карта прочитана не полностью. Повторите считывание карты.",
		4327: "Нет товаров для отображения",
		4328: "Информация о товаре отсутствует или неполна.",
		4329: "Справочник товаров переполнен. Выполните сверку.",
		4330: "Товар не найден.",
		4334: "Карта не считана. Либо цикл ожидания карты прерван нажатием клавиши ESC, либо истек таймаут.",
		4336: "Валюта указана неверно.",
		4337: "Из кассовой программы передан неверный тип карты.",
		4342: "Ошибка: невозможно запустить диалоговое окно UPOS.",
		4351: "Настроечные файлы *.tlv не найдены",
		4355: "Этот палец уже зарегистрирован в базе",
		4358: "Палец не опознан!",
		4362: "Пинпад временно заблокирован. Повторите операцию через 15 сек.",
		4363: "Превышена сумма оригинальной операции",
		4365: "Режим электронного захвата подписи не поддерживается",
		4366: "Рассчитанная скидка меньше минимально допустимой.",
		4367: "RKL: неверный формат запроса",
		4368: "RKL: не создана ключевая пара СА",
		4369: "RKL: не загружен сертификат хоста",
		4370: "RKL: не загружен публичный ключ СА",
		4371: "Текущая версия ОС не поддерживает RKL",
		4372: "RKL: хост CA дает некорректный ответ. Необходимо перезагрузить терминал",
		4380: "Штатная сверка итогов не выполнена.",
		4381: "Неверный формат QR-кода",
		4382: "Количество товара не должно превышать 4 млн. 200 тыс. единиц",
		4383: "Не удалось открыть сканирующее устройство",
		4384: "Считаны не все данные",
		4385: "Неверный номер пользователя",
		4388: "Место закончилось. Передайте чеки в банк.",
		4389: "Чек уже успешно передан",
		4400: "Возможно, карта преждевременно вынута",
		4401: "Позвоните в банк по т.(800)775-55-55 (495)544-45-46 (495)788-92-74",
		4402: "Позвоните в банк",
		4403: "Терминал заблокирован. Обратитесь в банк.",
		4404: "Изымите карту",
		4405: "Отказано",
		4406: "Общая ошибка",
		4407: "Изымите карту",
		4408: "Отказано",
		4410: "Позвоните в Амекс по т. 8(800)2006203 или 8(495)6443054",
		4411: "Отказано",
		4412: "Транзакция неверна",
		4413: "Сумма неверна",
		4414: "Карта неверна",
		4419: "Повторите позже",
		4433: "Изымите карту",
		4438: "Изымите карту",
		4441: "Изымите карту",
		4443: "Изымите карту",
		4450: "Отказано",
		4451: "Недостаточно средств",
		4454: "Срок действия карты истек",
		4455: "ПИН неверен",
		4457: "Транзакция не разрешена картой",
		4458: "Транзакция не разрешена терминалом",
		4461: "Исчерпан лимит",
		4462: "Карта ограничена",
		4465: "Исчерпан лимит",
		4468: "Повторите позже",
		4475: "ПИН заблокирован",
		4476: "Нет исходной операции",
		4478: "Счет неверен",
		4481: "Повторите позже",
		4482: "Отказано",
		4483: "Ошибка обработки ПИНа",
		4486: "Ошибка обработки ПИНа",
		4488: "Ошибка обработки ПИНа",
		4489: "МАС-код неверен",
		4490: "Неверная контрольная информация",
		4491: "Эмитент недоступен",
		4493: "Транзакция запрещена",
		4494: "Повторная транзакция",
		4495: "Отказано",
		4496: "Ошибка системы",
		4497: "Повторите операцию позже",
		4498: "МАС-код неверен",
		4499: "Ошибка формата",
		4710: "Такая карта не обслуживается",
		5001: "Отказ карты при выборе приложения Error application selection",
		5002: "Отказ карты. Некорректный ответ Chip error",
		5003: "Отказ карты. Некорректный ответ Chip error",
		5015: "Операция отменена клиентом",
		5029: "Мы принимаем только Visa",
		5042: "Ключ удаленной загрузки неверен",
		5044: "Нужно позвонить в банк Call issuer",
		5053: "На карте неверные данные Data integrity error",
		5055: "Карта отклонила операцию Transaction declined by card",
		5063: "Карта не ведет историю операций",
		5075: "Необходимо вставить карту в чиповый ридер",
		5084: "Введите пароль на телефоне",
		5100: "Подлинность данных не проверена",
		5101: "Ошибка проверки SDA Integrity check error",
		5102: "На карте нет нужных данных",
		5103: "Карта в стоп-листе",
		5104: "Ошибка проверки DDA Integrity check error",
		5105: "Ошибка проверки CDA Integrity check error",
		5108: "Неверная версия приложения EMV",
		5109: "Срок действия карты истек",
		5110: "Срок действия карты еще не настал",
		5111: "Для этой карты такая операция запрещена Operation is prohibited",
		5112: "Карта только что выдана",
		5116: "Личность клиента не проверена Cardholder verification error",
		5117: "Неизвестный код CVM Cardholder verification error",
		5118: "ПИН блокирован",
		5119: "Пин-пад неисправен",
		5120: "Клиент не ввел ПИН",
		5124: "Такая сумма требует связи с банком",
		5125: "Превышен нижний лимит карты",
		5126: "Превышен верхний лимит карты",
		5133: "Операция отклонена картой Transaction declined by card",
	}
	if v, ok := codes[code]; ok {
		return v
	}
	return fmt.Sprintf("Неизвестная ошибка (%v)", code)
}
