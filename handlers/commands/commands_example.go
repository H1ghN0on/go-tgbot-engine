package commands

import "github.com/H1ghN0on/go-tgbot-engine/bot/bottypes"

var NothingnessCommand bottypes.Command = bottypes.Command{Command: "/nothingness", Description: "Пустышка", SkipOnBack: true}

// Back Handler
var BackCommandCommand bottypes.Command = bottypes.Command{Command: "/back_command", Description: "Вернуться к предыдущей команде"}
var BackStateCommand bottypes.Command = bottypes.Command{Command: "/back_state", Description: "Вернуться к предыдущему состоянию"}

// Start Handler
var StartCommand bottypes.Command = bottypes.Command{Command: "/start", Description: "Команда Start"}
var ShowCommandsCommand bottypes.Command = bottypes.Command{Command: "/show_commands", Description: "Главное меню"}
var LevelOneCommand bottypes.Command = bottypes.Command{Command: "/level_one", Description: "Уровень 1"}
var LevelTwoCommand bottypes.Command = bottypes.Command{Command: "/level_two", Description: "Уровень 2"}
var LevelThreeCommand bottypes.Command = bottypes.Command{Command: "/level_three", Description: "Уровень 3"}
var BigMessagesCommand bottypes.Command = bottypes.Command{Command: "/big_messages", Description: "Большие сообщения"}

// Checkbox Handler
var CheckboxStartCommand bottypes.Command = bottypes.Command{Command: "/checkboxes_start", Description: "Запуск чекбоксов"}
var CheckboxFirstCommand bottypes.Command = bottypes.Command{Command: "/checkboxes_first", Description: "Установка первого чекбокса"}
var CheckboxSecondCommand bottypes.Command = bottypes.Command{Command: "/checkboxes_second", Description: "Установка второго чекбокса"}
var CheckboxThirdCommand bottypes.Command = bottypes.Command{Command: "/checkboxes_third", Description: "Установка третьего чекбокса"}
var CheckboxFourthCommand bottypes.Command = bottypes.Command{Command: "/checkboxes_fourth", Description: "Установка четвертого чекбокса"}
var CheckboxAcceptCommand bottypes.Command = bottypes.Command{Command: "/checkboxes_accept", Description: "Завершение работы с чекбоксами"}

// Info Handler
var SetInfoStartCommand bottypes.Command = bottypes.Command{Command: "/set_info_start", Description: "Запуск установки информации"}
var SetNameCommand bottypes.Command = bottypes.Command{Command: "/set_name", Description: "Установка имени"}
var SetSurnameCommand bottypes.Command = bottypes.Command{Command: "/set_surname", Description: "Установка фамилии"}
var SetAgeCommand bottypes.Command = bottypes.Command{Command: "/set_age", Description: "Установка возраста"}
var SetInfoEndCommand bottypes.Command = bottypes.Command{Command: "/set_info_end", Description: "Завершение установки информации"}

// Keyboard Handler
var KeyboardStartCommand bottypes.Command = bottypes.Command{Command: "/keyboard_start", Description: "Запуск клавиатуры"}
var KeyboardOneCommand bottypes.Command = bottypes.Command{Command: "/keyboard_one", Description: "Первая клавиатура"}
var KeyboardTwoCommand bottypes.Command = bottypes.Command{Command: "/keyboard_two", Description: "Вторая клавиатура"}
var KeyboardThreeCommand bottypes.Command = bottypes.Command{Command: "/keyboard_three", Description: "Третья клавиатура"}
var KeyboardFinishCommand bottypes.Command = bottypes.Command{Command: "/keyboard_finish", Description: "Завершение клавиатуры"}

// Level Four Handler
var LevelFourStartCommand bottypes.Command = bottypes.Command{Command: "/level_four_start", Description: "Твой последний выход"}
var LevelFourOneCommand bottypes.Command = bottypes.Command{Command: "/level_four_one", Description: "РАЗ"}
var LevelFourTwoCommand bottypes.Command = bottypes.Command{Command: "/level_four_two", Description: "ДВА"}
var LevelFourThreeCommand bottypes.Command = bottypes.Command{Command: "/level_four_three", Description: "ТРИ"}
var LevelFourFourCommand bottypes.Command = bottypes.Command{Command: "/level_four_four", Description: "ЧЕТЫРЕ"}

// Dynamic Keyboard Handler
var DynamicKeyboardStartCommand bottypes.Command = bottypes.Command{Command: "/dynamic_keyboard_start", Description: "Запуск динамической клавиатуры"}
var DynamicKeyboardFirstStageCommand bottypes.Command = bottypes.Command{Command: "/dynamic_keyboard_first_stage", Description: "Первая фаза динамической клавиатуры"}
var DynamicKeyboardSecondStageCommand bottypes.Command = bottypes.Command{Command: "/dynamic_keyboard_second_stage", Description: "Вторая фаза динамической клавиатуры"}
var DynamicKeyboardFinishCommand bottypes.Command = bottypes.Command{Command: "/dynamic_keyboard_finish", Description: "Завершение динамической клавиатуры"}

// Calendar Handler
var CalendarStartCommand bottypes.Command = bottypes.Command{Command: "/calendar_start", Description: "Запуск календаря"}
var CalendarChooseCommand bottypes.Command = bottypes.Command{Command: "/calendar_choose", Description: "Выбор календаря"}
var CalendarChooseFirstCommand bottypes.Command = bottypes.Command{Command: "/calendar_choose_first", Description: "Расписание 1", SkipOnBack: true}
var CalendarChooseSecondCommand bottypes.Command = bottypes.Command{Command: "/calendar_choose_second", Description: "Расписание 2", SkipOnBack: true}
var CalendarLaunchCommand bottypes.Command = bottypes.Command{Command: "/calendar_launch", Description: "Предлагает кнопку запуска календаря"}
var CalendarNextMonthCommand bottypes.Command = bottypes.Command{Command: "/calendar_next_month", Description: "Следующий месяц", SkipOnBack: true}
var CalendarPrevMonthCommand bottypes.Command = bottypes.Command{Command: "/calendar_prev_month", Description: "Предыдущий месяц", SkipOnBack: true}
var CalendarNextYearCommand bottypes.Command = bottypes.Command{Command: "/calendar_next_year", Description: "Следующий год", SkipOnBack: true}
var CalendarPrevYearCommand bottypes.Command = bottypes.Command{Command: "/calendar_prev_year", Description: "Следующий год", SkipOnBack: true}
var CalendarSetDayCommand bottypes.Command = bottypes.Command{Command: "/calendar_set_day", Description: "Выбрать день"}
var CalendarSetTimeCommand bottypes.Command = bottypes.Command{Command: "/calendar_set_time", Description: "Выбрать время"}
var CalendarFinishCommand bottypes.Command = bottypes.Command{Command: "/calendar_finish", Description: "Закончить с календарем"}

var Commands = []bottypes.Command{
	NothingnessCommand,
	BackCommandCommand,
	BackStateCommand,

	StartCommand,
	ShowCommandsCommand,
	LevelOneCommand,
	LevelTwoCommand,
	LevelThreeCommand,
	BigMessagesCommand,

	CheckboxStartCommand,
	CheckboxFirstCommand,
	CheckboxSecondCommand,
	CheckboxThirdCommand,
	CheckboxFourthCommand,
	CheckboxAcceptCommand,

	SetInfoStartCommand,
	SetNameCommand,
	SetSurnameCommand,
	SetAgeCommand,
	SetInfoEndCommand,

	KeyboardStartCommand,
	KeyboardOneCommand,
	KeyboardTwoCommand,
	KeyboardThreeCommand,
	KeyboardFinishCommand,

	LevelFourStartCommand,
	LevelFourOneCommand,
	LevelFourTwoCommand,
	LevelFourThreeCommand,
	LevelFourFourCommand,

	DynamicKeyboardStartCommand,
	DynamicKeyboardFirstStageCommand,
	DynamicKeyboardSecondStageCommand,
	DynamicKeyboardFinishCommand,

	CalendarStartCommand,
	CalendarChooseCommand,
	CalendarChooseFirstCommand,
	CalendarChooseSecondCommand,
	CalendarLaunchCommand,
	CalendarNextMonthCommand,
	CalendarPrevMonthCommand,
	CalendarNextYearCommand,
	CalendarPrevYearCommand,
	CalendarSetDayCommand,
	CalendarSetTimeCommand,
	CalendarFinishCommand,
}
