# go-sbrfdll
Обертка над sbrf.dll для работы с банковскими терминалами Сбера через USB под Windows с использованием OLE Automation.

## Установка
Перед запуском необходимо установить драйвера для работы с терминалом и зарегистрировать sbrf.dll в системе.
Если у вас 64-битная система, то:
<br>– sbrf.dll регистрировать через Windows/SysWOW64/regsvr32.exe
<br>– компилировать проект в 32-битном Go.

## Реализованные операции
| Операция       | Аргументы                                     |
|----------------|:----------------------------------------------|
| проверка связи | <id операции> connect                         |
| оплата         | <id операции> pay <сумма в копейках>          |
| возврат        | <id операции> return <сумма в копейках> <RRN> |
| отмена         | <id операции> cancel <сумма в копейках> <RRN> |
| сверка итогов  | <id операции> close                           |

## Результат выполнения
Результат выполнения операции выводится в файл result.bin в текущей директории, а лог в файл sber.log.
<br>Структура выходного файла result.bin (в байтах):
```
[0-3] - id операции (берется из входных параметров)
[4-7] - результат операции (0 если успешно)
[8-19] - RRN операции
[20-...] - банковский слип (чек) если операция успешна, либо текстовое описание ошибки
```